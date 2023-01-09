package main

import (
	"fmt"
	"os"
	"strings"

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

// check for the `package-lock.json`, `yarn.lock` or `bower.json` (required for SCA)
func CheckIfSCAFileExists(path string) bool {
	// we don't want to look for `package-lock.json` and `yarn.lock` within `bower_components`
	if !strings.Contains(path, "bower_components") {
		packageManagerFiles := [2]string{"package-lock.json", "yarn.lock"}

		for _, element := range packageManagerFiles {
			if strings.HasSuffix(path, element) {
				return true
			}
		}
	}

	// NOTE: It looks like the `bower.json` file would be in `bower_components`? (tbh, I am not 100% sure how Bower
	// works exactly, but it's been depreacted like forever and I can't really be bothered looking into how exactly it works)
	bowerFile := "bower.json"
	return strings.HasSuffix(path, bowerFile)
}

// check for the `node_modules` folder
func IsNodeModules(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"node_modules") {
		if !didPrintNodeModulesMsg {
			log.Info("\tIgnoring the entire `node_modules` folder")
			didPrintNodeModulesMsg = true
		}

		return true
	}

	return false
}

// check for the `bower_components` folder
func IsBowerComponents(path string) bool {
	// NOTE: Turns out, we don't want to omit the `bower_components` folder since it's required for Bower.
	// Thus, we don't do anything here for now (this may change in the future)

	// if strings.Contains(path, string(os.PathSeparator)+"bower_components") {
	// 	if !didPrintBowerComponentsMsg {
	// 		log.Info("\tIgnoring the entire `bower_components` folder")
	// 		didPrintBowerComponentsMsg = true
	// 	}

	// 	return true
	// }

	return false
}

// check if it is a `test` path (i.e., a file that e.g. contains unit tests)
func IsInTestFolder(path string, testsPath string) bool {
	// Test folders are treated as follows:
	// 	- if `-tests` is provided, then only the provided path will be treated as a test directory (and thus, excluded)
	// 	- if `-tests` is not provided, then `IsCommonTest()` will be called to exclude common test folders
	if testsPath == "" {
		return IsCommonTestFolder(path)
	}

	// At this point, `testsPath` may have a value like this: "sample-node-project/test".
	// Thus, we want to do 2 things:
	//	  - Exclude this folder itself, i.e. check for a path that ends with "sample-node-project/test"
	// 	  - Exclude any file in it, i.e. check for path that contains "sample-node-project/test/"
	fileInTestFolderPath := testsPath + string(os.PathSeparator)

	if strings.HasSuffix(path, testsPath) || strings.Contains(path, fileInTestFolderPath) {
		if !didPrintTestsMsg {
			log.Info("\tIgnoring the entire content of the `" + testsPath + "` folder (contains test files)")
			didPrintTestsMsg = true
		}

		return true
	}

	return false
}

func IsCommonTestFolder(path string) bool {
	testPaths := [4]string{"test", "tests", "e2e", "__tests__"}

	for _, testPath := range testPaths {
		// Here, we want to do 2 things:
		//	  - Exclude the `testPath` itself, i.e. check for a path that ends e.g. with "/e2e"
		// 	  - Exclude any file in it, i.e. check for path that contains e.g. "/e2e/"
		testFolderPath := string(os.PathSeparator) + testPath
		fileInTestFolderPath := testFolderPath + string(os.PathSeparator)

		if strings.HasSuffix(path, testFolderPath) || strings.Contains(path, fileInTestFolderPath) {
			if !didPrintDefaultTestFoldersMsg {
				log.Info("\tIgnoring common test folders (such as `e2e`)")
				didPrintDefaultTestFoldersMsg = true
			}

			return true
		}
	}

	return false
}

// check for common test files (like .spec.js)
func IsTestFile(path string) bool {
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

// check for style sheets (like .css and .scss)
func IsStyleSheet(path string) bool {
	if strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".scss") {
		if !didPrintStylesheetsMsg {
			log.Info("\tIgnoring style sheets (such as `.css`)")
			didPrintStylesheetsMsg = true
		}

		return true
	}

	return false
}

// check for images (like .jpg, .png, .jpeg)
func IsImage(path string) bool {
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

// check for documents (like .pdf, .md)
func IsDocument(path string) bool {
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

// check for video files
func IsVideo(path string) bool {
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

// check for fonts (like .woff)
func IsFont(path string) bool {
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

// check for the `.git` folder
func IsGitFolder(path string) bool {
	// checks for the ".git" folder itself, i.e. for a path that ends with `/.git`
	// ... or for files within the ".git" folder, i.e. for a path that contains `/.git/`
	gitFolderPath := string(os.PathSeparator) + ".git"
	fileInGitFolderPath := gitFolderPath + string(os.PathSeparator)

	if strings.HasSuffix(path, gitFolderPath) || strings.Contains(path, fileInGitFolderPath) {
		if !didPrintGitFolderMsg {
			log.Info("\tIgnoring `.git`")
			didPrintGitFolderMsg = true
		}

		return true
	}

	return false
}

// check for the dbs (like .db, .sqlite3)
func IsDb(path string) bool {
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

// check for the `build` folder
func IsBuildFolder(path string) bool {
	// checks for the "build" folder itself, i.e. for a path that ends with `/build`
	// ... or for files within the "build" folder, i.e. for a path that contains `/build/`
	buildFolderPath := string(os.PathSeparator) + "build"
	fileInBuildFolderPath := buildFolderPath + string(os.PathSeparator)
	log.Debug(fmt.Sprintf("\tCheck for '%s' and '%s' paths (would be omitted)", buildFolderPath, fileInBuildFolderPath))
	log.Debug("Path: ", path)

	if strings.HasSuffix(path, buildFolderPath) || strings.Contains(path, fileInBuildFolderPath) {
		if !didPrintBuildMsg {
			log.Info("\tIgnoring `build` folder")
			didPrintBuildMsg = true
		}

		return true
	}

	return false
}

// check for the `dist` folder
func IsDistFolder(path string) bool {
	// checks for the "dist" folder itself, i.e. for a path that ends with `/dist`
	// ... or for files within the "dist" folder, i.e. for a path that contains `/dist/`
	distFolderPath := string(os.PathSeparator) + "dist"
	fileInDistFolderPath := distFolderPath + string(os.PathSeparator)

	if strings.HasSuffix(path, distFolderPath) || strings.Contains(path, fileInDistFolderPath) {
		if !didPrintDistMsg {
			log.Info("\tIgnoring `dist` folder")
			didPrintDistMsg = true
		}

		return true
	}

	return false
}

// check for the `public` folder
func IsPublicFolder(path string) bool {
	// checks for the "dist" folder itself, i.e. for a path that ends with `/dist`
	// ... or for files within the "dist" folder, i.e. for a path that contains `/dist/`
	publicFolderPath := string(os.PathSeparator) + "public"
	fileInPublicFolderPath := publicFolderPath + string(os.PathSeparator)

	if strings.HasSuffix(path, publicFolderPath) || strings.Contains(path, fileInPublicFolderPath) {
		if !didPrintPublicMsg {
			log.Info("\tIgnoring `build` folder")
			didPrintPublicMsg = true
		}

		return true
	}

	return false
}

// check for IDE folder (like .code, .idea)
func IsIdeFolder(path string) bool {
	idePaths := [2]string{".vscode", ".idea"}

	for _, element := range idePaths {
		// NOTE: This check should be fine using "Contains" because I don't expect these IDE folder names to be used for
		// 	anything useful (i.e., I don't anticipate a folder name with ".idea" in its name to contain anything useful)
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

// check for the "misc" not required stuff
func IsMiscNotRequiredFile(path string) bool {
	notRequiredSuffices := []string{
		".DS_Store", "__MACOSX", ".gitignore", ".gitkeep", ".gitattributes", ".npmignore", "CNAME", "tsconfig.json",
		"tslint.json", "karma.conf.js", "angular.json", ".travis.yml", ".browserslistrc", ".editorconfig",
		".d.ts", "protractor.conf.js", ".spec.json", "tsconfig.app.json", "polyfills.ts", "LICENSE", "LICENSE.md",
	}

	for _, element := range notRequiredSuffices {
		if strings.HasSuffix(path, element) {
			// NOTE: At the moment, these "misc" files aren't logged to avoid logging too much
			return true
		}
	}

	return false
}
