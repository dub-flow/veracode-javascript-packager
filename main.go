package main

import (
	"archive/zip"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// flag to make sure a message is only logged once
var didPrintNodeModulesMsg bool = false
var didPrintTestsMsg bool = false
var didPrintStylesheetsMsg bool = false
var didPrintImagesMsg bool = false
var didPrintDocumentsMsg bool = false
var didPrintGitFolderMsg bool = false

func main() {
	// parse all the command line flags
	sourcePtr := flag.String("source", "sample-node-project", "The path of the Node.js app you want to package")
	targetPtr := flag.String("target", ".", "The path where you want the output.zip to be stored to")
	testsPtr := flag.String("tests", "test", "The path that contains your Node.js test files (relative to the source)")
	flag.Parse()

	outputZipPath := *targetPtr + "/output.zip"
	testsPath := *sourcePtr + "/" + *testsPtr + "/"

	log.Println("Veracode Node Packager - Started")
	log.Println("Source directory to zip up:", *sourcePtr)
	log.Println("Test directory (its content will be omitted):", testsPath)
	log.Println("Zip Process - Started...")

	// NOTE: Search for `package-lock.json` file

	if err := zipSource(*sourcePtr, outputZipPath, testsPath); err != nil {
		log.Fatal(err)
	}

	log.Println("Zip Process - Finished")
	log.Println("Wrote output to:", outputZipPath)
}

func zipSource(source, target string, testsPath string) error {
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
	// check for the `node_modules` folder
	if isNodeModules(path) {
		return false
	}

	// check if it is a `test` path (i.e., a file that e.g. contains unit tests)
	if isTestFile(path, testsPath) {
		return false
	}

	// check for style sheets (.css and .scss)
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

	// check for the `.git` folder
	if isGitFolder(path) {
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
	if strings.Contains(path, "node_modules") {
		if !didPrintNodeModulesMsg {
			log.Println("Ignoring the entire `node_modules` folder")
			didPrintNodeModulesMsg = true
		}

		return true
	}

	return false
}

func isTestFile(path string, testsPath string) bool {
	if strings.Contains(path, testsPath) {
		if !didPrintTestsMsg {
			log.Printf("Ignoring the entire content of the `%s` folder (contains test files)", testsPath)
			didPrintTestsMsg = true
		}

		return true
	}

	return false
}

func isStyleSheet(path string) bool {
	if strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".scss") {
		if !didPrintStylesheetsMsg {
			log.Println("Ignoring style sheets (such as `.css`)")
			didPrintStylesheetsMsg = true
		}

		return true
	}

	return false
}

func isImage(path string) bool {
	imageExtensions := [6]string{".jpg", ".png", ".jpeg", ".gif", ".svg", ".bmp"}

	for _, element := range imageExtensions {
		if strings.HasSuffix(path, element) {
			if !didPrintImagesMsg {
				log.Println("Ignoring images (such as `.jpg`)")
				didPrintImagesMsg = true
			}

			return true
		}
	}

	return false
}

func isDocument(path string) bool {
	documentExtensions := [2]string{".pdf", ".md"}

	for _, element := range documentExtensions {
		if strings.HasSuffix(path, element) {
			if !didPrintDocumentsMsg {
				log.Println("Ignoring documents (such as `.pdf`)")
				didPrintDocumentsMsg = true
			}

			return true
		}
	}

	return false
}

func isGitFolder(path string) bool {
	if strings.Contains(path, ".git") {
		if !didPrintGitFolderMsg {
			log.Println("Ignoring `.git`")
			didPrintGitFolderMsg = true
		}

		return true
	}

	return false
}

func isMiscNotRequiredFile(path string) bool {
	notRequiredSuffices := [1]string{".DS_Store"}

	for _, element := range notRequiredSuffices {
		if strings.HasSuffix(path, element) {
			// NOTE: At the moment, these "misc" files aren't logged to avoid logging too much
			return true
		}
	}

	return false
}
