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
)

// flags to make sure a message is only logged once
var didPrintNodeModulesMsg bool = false
var didPrintTestsMsg bool = false
var didPrintDefaultTestExtensionsMsg bool = false
var didPrintDefaultTestFoldersMsg bool = false
var didPrintStylesheetsMsg bool = false
var didPrintImagesMsg bool = false
var didPrintDocumentsMsg bool = false
var didPrintFontsMsg bool = false
var didPrintIdesMsg bool = false
var didPrintBuildMsg bool = false
var didPrintPublicMsg bool = false
var didPrintDistMsg bool = false
var didPrintDbsMsg bool = false
var didPrintGitFolderMsg bool = false
var didPrintBowerComponentsMsg bool = false
var didPrintVideoMsg bool = false

func main() {
	// parse all the command line flags
	sourcePtr := flag.String("source", "./sample-projects/sample-node-project", "The path of the JavaScript app you want to package")
	targetPtr := flag.String("target", ".", "The path where you want the vc-output.zip to be stored to")
	testsPtr := flag.String("tests", "", "The path that contains your test files (relative to the source). Uses a heuristic to identifiy tests automatically in case no path is provided")
	flag.Parse()

	log.Info("#################################################")
	log.Info("#                                               #")
	log.Info("#   Veracode JavaScript Packager (Unofficial)   #")
	log.Info("#                                               #")
	log.Info("#################################################" + "\n\n")

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
		// combine that last segment of the `sourcePtr`` with the value provided via `-test`.
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
		if !strings.Contains(path, "node_modules") && !strings.Contains(path, "bower_components") {
			// check for the `package-lock.json`, `yarn.lock` or `bower.json` (required for SCA)
			if !doesSCAFileExist {
				packageManagerFiles := [3]string{"package-lock.json", "yarn.lock", "bower.json"}
				for _, element := range packageManagerFiles {
					if strings.HasSuffix(path, element) {
						doesSCAFileExist = true
					}
				}
			}

			// check for `.map` files
			if strings.HasSuffix(path, ".map") {
				doesMapFileExist = true
			}
		}

		return nil
	})

	if err != nil {
		log.Error(err)
	}

	if !doesSCAFileExist {
		log.Warn("\tNo `package-lock.json` or `yarn.lock` or `bower.json` file found.. (This file is required for Veracode SCA)")
		log.Warn("\tYou may not receive Veracode SCA results")
	}

	if doesMapFileExist {
		log.Warn("\tThe 1st party code contains `.map` files (which indicates minified JavaScript)...")
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

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}

		// check if the path is required for the upload (otherwise, it will be omitted)
		if !isRequired(header.Name, testsPath) {
			return nil
		}

		if info.IsDir() {
			header.Name += "/"
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
	// check for the `node_modules` and `bower_components` folders
	if isNodeModules(path) || isBowerComponents(path) {
		return false
	}

	// check if it is a `test` path (i.e., a file that e.g. contains unit tests)
	if isInTestFolder(path, testsPath) {
		return false
	}

	// check for common test files (like .spec.js)
	if isTestFile(path) {
		return false
	}

	// check for style sheets (like .css and .scss)
	if isStyleSheet(path) {
		return false
	}

	// check for images (like .jpg, .png, .jpeg)
	if isImage(path) {
		return false
	}

	// check for documents (like .pdf, .md)
	if isDocument(path) {
		return false
	}

	// check for fonts (like .woff)
	if isFont(path) {
		return false
	}

	// check for the `.git` folder
	if isGitFolder(path) {
		return false
	}

	// check for the dbs (like .db, .sqlite3)
	if isDb(path) {
		return false
	}

	// check for the `build`, `dist` and `public` folders
	if isBuildFolder(path) || isDistFolder(path) || isPublicFolder(path) {
		return false
	}

	// check for IDE folder (like .code, .idea)
	if isIdeFolder(path) {
		return false
	}

	// check for video files
	if isVideo(path) {
		return false
	}

	// check for the "misc" not required stuff
	if isMiscNotRequiredFile(path) {
		return false
	}

	// the default is to not omit the file
	return true
}

func isNodeModules(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"node_modules") {
		if !didPrintNodeModulesMsg {
			log.Info("\tIgnoring the entire `node_modules` folder")
			didPrintNodeModulesMsg = true
		}

		return true
	}

	return false
}

func isBowerComponents(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"bower_components") {
		if !didPrintBowerComponentsMsg {
			log.Info("\tIgnoring the entire `bower_components` folder")
			didPrintBowerComponentsMsg = true
		}

		return true
	}

	return false
}

func isInTestFolder(path string, testsPath string) bool {
	// Test folders are treated as follows:
	// 	- if `-tests` is provided, then only the provided path will be treated as a test directory (and thus, excluded)
	// 	- if `-tests` is not provided, then `isCommonTest()` will be called to exclude common test folders
	if testsPath == "" {
		return isCommonTestFolder(path)
	}

	if strings.Contains(path, string(os.PathSeparator)+testsPath) {
		if !didPrintTestsMsg {
			log.Info("\tIgnoring the entire content of the `" + testsPath + "` folder (contains test files)")
			didPrintTestsMsg = true
		}

		return true
	}

	return false
}

func isCommonTestFolder(path string) bool {
	testPaths := [3]string{"test", "e2e", "__tests__"}

	for _, element := range testPaths {
		if strings.Contains(path, string(os.PathSeparator)+element) {
			if !didPrintDefaultTestFoldersMsg {
				log.Info("\tIgnoring common test folders (such as `e2e`)")
				didPrintDefaultTestFoldersMsg = true
			}

			return true
		}
	}

	return false
}

func isTestFile(path string) bool {
	testExtensions := [3]string{".spec.ts", ".test.tsx", ".spec.js"}

	for _, element := range testExtensions {
		if strings.HasSuffix(path, element) {
			if !didPrintDefaultTestExtensionsMsg {
				log.Info("\tIgnoring common test extensions (such as `.spec.ts`)")
				didPrintDefaultTestExtensionsMsg = true
			}

			return true
		}
	}

	return false
}

func isStyleSheet(path string) bool {
	if strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".scss") {
		if !didPrintStylesheetsMsg {
			log.Info("\tIgnoring style sheets (such as `.css`)")
			didPrintStylesheetsMsg = true
		}

		return true
	}

	return false
}

func isImage(path string) bool {
	imageExtensions := [8]string{".jpg", ".png", ".jpeg", ".gif", ".svg", ".bmp", ".ico", ".icns"}

	for _, element := range imageExtensions {
		if strings.HasSuffix(path, element) {
			if !didPrintImagesMsg {
				log.Info("\tIgnoring images (such as `.jpg`)")
				didPrintImagesMsg = true
			}

			return true
		}
	}

	return false
}

func isDocument(path string) bool {
	// inspired by https://en.wikipedia.org/wiki/List_of_Microsoft_Office_filename_extensions (and additionally `.md`)
	documentExtensions := [38]string{
		".pdf",
		".md",
		".doc", ".dot", ".wbk", ".docx", ".docm", ".dotx", ".dotm", ".docb", ".wll", ".wwl",
		".xls", ".xlt", ".xlm", ".xll_", ".xla_", ".xla5", ".xla8",
		".xlsx", ".xlsm", ".xltx", ".xltm",
		".ppt", ".pot", ".pps", ".pptx", ".pptm", ".potx", ".potm",
		".one", ".ecf",
		".ACCDA", ".ACCDB", ".ACCDE", ".ACCDT", ".MDA", ".MDE",
	}

	for _, element := range documentExtensions {
		if strings.HasSuffix(path, element) {
			if !didPrintDocumentsMsg {
				log.Info("\tIgnoring documents (such as `.pdf`, `.docx`)")
				didPrintDocumentsMsg = true
			}

			return true
		}
	}

	return false
}

func isVideo(path string) bool {
	// inspired by this list: https://en.wikipedia.org/wiki/Video_file_format
	videoExtensions := [18]string{
		".mp4", ".webm", ".mkv", ".flv", ".vob", ".ogv", ".drc", ".gifv", ".mng", ".avi", ".mov", ".qt", ".mts", ".wmv", ".amv",
		".svi", ".m4v", ".mpg",
	}

	for _, element := range videoExtensions {
		if strings.HasSuffix(path, element) {
			if !didPrintVideoMsg {
				log.Info("\tIgnoring videos (such as `.mp4`)")
				didPrintVideoMsg = true
			}

			return true
		}
	}

	return false
}

func isFont(path string) bool {
	fontExtensions := [4]string{".ttf", ".otf", ".woff", ".woff2"}

	for _, element := range fontExtensions {
		if strings.HasSuffix(path, element) {
			if !didPrintFontsMsg {
				log.Info("\tIgnoring fonts (such as `.woff`)")
				didPrintFontsMsg = true
			}

			return true
		}
	}

	return false
}

func isGitFolder(path string) bool {
	if strings.HasSuffix(path, string(os.PathSeparator)+".git") {
		if !didPrintGitFolderMsg {
			log.Info("\tIgnoring `.git`")
			didPrintGitFolderMsg = true
		}

		return true
	}

	return false
}

func isDb(path string) bool {
	documentExtensions := [6]string{".db", ".db3", ".sdb", ".sqlite", ".sqlite2", ".sqlite3"}

	for _, element := range documentExtensions {
		if strings.HasSuffix(path, element) {
			if !didPrintDbsMsg {
				log.Info("\tIgnoring dbs (such as `.sqlite3`)")
				didPrintDbsMsg = true
			}

			return true
		}
	}

	return false
}

func isBuildFolder(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"build") {
		if !didPrintBuildMsg {
			log.Info("\tIgnoring `build` folder")
			didPrintBuildMsg = true
		}

		return true
	}

	return false
}

func isDistFolder(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"dist") {
		if !didPrintDistMsg {
			log.Info("\tIgnoring `dist` folder")
			didPrintDistMsg = true
		}

		return true
	}

	return false
}

func isPublicFolder(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"public") {
		if !didPrintPublicMsg {
			log.Info("\tIgnoring `build` folder")
			didPrintPublicMsg = true
		}

		return true
	}

	return false
}

func isIdeFolder(path string) bool {
	idePaths := [2]string{".vscode", ".idea"}

	for _, element := range idePaths {
		if strings.Contains(path, string(os.PathSeparator)+element) {
			if !didPrintIdesMsg {
				log.Info("\tIgnoring IDE folder (such as .code, .idea)")
				didPrintIdesMsg = true
			}

			return true
		}
	}

	return false
}

func isMiscNotRequiredFile(path string) bool {
	notRequiredSuffices := [3]string{".DS_Store", "__MACOSX", ".gitignore"}

	for _, element := range notRequiredSuffices {
		if strings.HasSuffix(path, element) {
			// NOTE: At the moment, these "misc" files aren't logged to avoid logging too much
			return true
		}
	}

	return false
}
