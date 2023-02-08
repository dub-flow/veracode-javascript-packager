package main

import (
	"archive/zip"
	"flag"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/fatih/color"
)

func main() {
	// parse all the command line flags
	sourcePtr := flag.String("source", "", "The path of the JavaScript app you want to package (required)")
	targetPtr := flag.String("target", ".", "The path where you want the vc-output.zip to be stored to")
	testsPtr := flag.String("tests", "", "The path that contains your test files (relative to the source). Uses a heuristic to identifiy tests automatically in case no path is provided")

	flag.Parse()

	color.Green("#################################################")
	color.Green("#                                               #")
	color.Green("#   Veracode JavaScript Packager (Unofficial)   #")
	color.Green("#                                               #")
	color.Green("#################################################" + "\n\n")

	color.Yellow("Current version: %s\n\n", AppVersion)

	// check if a later version of this tool exists
	notifyOfUpdates()

	// fail if `--source` was not provided
	if *sourcePtr == "" {
		color.Red("No `-source` was provided. Run `--help` for the built-in help.")
		return
	}

	// add the current date to the output zip name, like e.g. "2023-Jan-04"
	currentTime := time.Now()
	outputZipPath := filepath.Join(*targetPtr, "vc-output_"+currentTime.Format("2006-Jan-02")+".zip")

	// echo the provided flags
	var testsPath string
	log.Info("Provided Flags:")
	log.Info("\t`-source` directory to zip up: ", *sourcePtr)
	log.Info("\t`-target` directory for the output: ", *targetPtr)

	if *testsPtr == "" {
		log.Info("\tNo `-test` directory was provided... Heuristics will be used to identify (and omit) common test directory names" + "\n\n")
		testsPath = ""
	} else {
		// combine that last segment of the `sourcePtr` with the value provided via `-test`.
		// Example: If `-test mytests` and `-source /some/node-project`, then `testsPath` will be: "node-project/mytests"
		testsPath = filepath.Join(path.Base(*sourcePtr), *testsPtr)
		log.Info("\tProvided `-test` directory (its content will be omitted): ", testsPath, "\n\n")
	}

	// check for some "smells" (e.g. the `package-lock.json` file is missing), and print corresponding warnings/errors
	log.Info("Checking for 'smells' that indicate packaging issues - Started...")
	checkForPotentialSmells(*sourcePtr)
	log.Info("'Smells' Check - Done\n\n")

	log.Info("Creating a Zip while omitting non-required files - Started...")
	// generate the zip file, and omit all non-required files
	if err := zipSource(*sourcePtr, outputZipPath, testsPath); err != nil {
		log.Error(err)
	}

	log.Info("Zip Process - Done")
	log.Info("Wrote archive to: ", outputZipPath)
	log.Info("Please upload this archive to the Veracode Platform")
}

func checkForPotentialSmells(source string) {
	doesSCAFileExist := false
	doesMapFileExist := false

	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// only do checks for first party code
		if !strings.Contains(path, "node_modules") {
			// check if one of the files required for SCA exists... Note that `bower.json` may be part of `bower_components`
			if !doesSCAFileExist {
				doesSCAFileExist = CheckIfSCAFileExists(path)
			}

			// check for `.map` files (only in non-3rd party and "non-build" code)
			if !strings.Contains(path, "bower_components") && !strings.Contains(path, "build") &&
				!strings.Contains(path, "dist") && !strings.Contains(path, "public") {
				if strings.HasSuffix(path, ".map") {
					doesMapFileExist = true
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Error(err)
	}

	if !doesSCAFileExist {
		log.Warn("\tNo `package-lock.json` or `yarn.lock` or `bower.json` file found.. (This file is required for Veracode SCA)..." +
			" You may not receive Veracode SCA results!")
	}

	if doesMapFileExist {
		log.Warn("\tThe 1st party code contains `.map` files outside of `/build`, `/dist` or `/public` (which indicates minified JavaScript)...")
		log.Warn("\tPlease pass a directory to this tool that contains the unminified/unbundled/unconcatenated JavaScript (or TypeScript)")
	}
}

func zipSource(source string, target string, testsPath string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// avoids processing the created zip...
		// 	- Say the tool is finished and an `/vc-output_2023-Jan-05.zip` is created...
		//  - In this case, the analysis may restart with this zip as `path`
		// 		- This edge case was observed when running the tool within a sample JS app..
		//		- ... i.e., `veracode-js-packager -source . -target .`
		if strings.HasSuffix(path, ".zip") {
			return nil
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		// 	-> We want the following:
		//		- Say `-source some/path/my-js-project` is provided...
		//			- Now, say we have a path `some/path/my-js-project/build/some.js`....
		//		- In this scenario, we want `header.Name` to be `build/some.js`
		header.Name, err = filepath.Rel(source, path)
		if err != nil {
			return err
		}

		// avoids the `./` folder in the root of the output zip
		if header.Name == "." {
			return nil
		}

		// prepends the `/` we want before e.g. `build/some.js`
		headerNameWithSlash := string(os.PathSeparator) + header.Name

		// check if the path is required for the upload (otherwise, it will be omitted)
		if !isRequired(headerNameWithSlash, testsPath) {
			return nil
		}

		if info.IsDir() {
			// add e.g. a `/` if the current path is a directory
			header.Name += string(os.PathSeparator)
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func isRequired(path string, testsPath string) bool {
	return !IsNodeModules(path) &&
		!IsAngularCacheFolder(path) &&
		!IsBowerComponents(path) &&
		!IsGitFolder(path) &&
		!IsInTestFolder(path, testsPath) &&
		!IsTestFile(path) &&
		!IsStyleSheet(path) &&
		!IsImage(path) &&
		!IsVideo(path) &&
		!IsDocument(path) &&
		!IsFont(path) &&
		!IsDb(path) &&
		!IsBuildFolder(path) &&
		!IsDistFolder(path) &&
		!IsPublicFolder(path) &&
		!IsIdeFolder(path) &&
		!IsMinified(path) &&
		!IsMiscNotRequiredFile(path)

	// the default is to not omit the file
}
