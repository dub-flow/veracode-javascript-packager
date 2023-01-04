package main

import (
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

func IsBowerComponents(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"bower_components") {
		if !didPrintBowerComponentsMsg {
			log.Info("\tIgnoring the entire `bower_components` folder")
			didPrintBowerComponentsMsg = true
		}

		return true
	}

	return false
}

func IsInTestFolder(path string, testsPath string) bool {
	// Test folders are treated as follows:
	// 	- if `-tests` is provided, then only the provided path will be treated as a test directory (and thus, excluded)
	// 	- if `-tests` is not provided, then `IsCommonTest()` will be called to exclude common test folders
	if testsPath == "" {
		return IsCommonTestFolder(path)
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

func IsCommonTestFolder(path string) bool {
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

func IsGitFolder(path string) bool {
	if strings.HasSuffix(path, string(os.PathSeparator)+".git") {
		if !didPrintGitFolderMsg {
			log.Info("\tIgnoring `.git`")
			didPrintGitFolderMsg = true
		}

		return true
	}

	return false
}

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

func IsBuildFolder(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"build") {
		if !didPrintBuildMsg {
			log.Info("\tIgnoring `build` folder")
			didPrintBuildMsg = true
		}

		return true
	}

	return false
}

func IsDistFolder(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"dist") {
		if !didPrintDistMsg {
			log.Info("\tIgnoring `dist` folder")
			didPrintDistMsg = true
		}

		return true
	}

	return false
}

func IsPublicFolder(path string) bool {
	if strings.Contains(path, string(os.PathSeparator)+"public") {
		if !didPrintPublicMsg {
			log.Info("\tIgnoring `build` folder")
			didPrintPublicMsg = true
		}

		return true
	}

	return false
}

func IsIdeFolder(path string) bool {
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

func IsMiscNotRequiredFile(path string) bool {
	notRequiredSuffices := [3]string{".DS_Store", "__MACOSX", ".gitignore"}

	for _, element := range notRequiredSuffices {
		if strings.HasSuffix(path, element) {
			// NOTE: At the moment, these "misc" files aren't logged to avoid logging too much
			return true
		}
	}

	return false
}
