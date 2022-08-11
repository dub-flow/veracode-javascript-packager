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
var didPrintStylesheetsMsg bool = false

func main() {
	log.Println("Node Packager - Started")

	// parse all the command line flags
	sourcePtr := flag.String("source", "../test-projects/my-node-test", "the path of the Node.js app you want to package")
	targetPtr := flag.String("target", ".", "the path where you want the output.zip to be stored to")
	flag.Parse()

	outputZipPath := *targetPtr + "/output.zip"
	// folderToZip := "test-projects/my-node-test" // "../test-projects/the-example-app.nodejs-master"
	// outputZip := "output.zip"

	log.Println("Source directory to zip up:", *sourcePtr)
	log.Println("Zip Process - Started")

	if err := zipSource(*sourcePtr, outputZipPath); err != nil {
		log.Fatal(err)
	}

	log.Println("Zip Process - Finished")
	log.Println("Wrote output to:", outputZipPath)
}

func zipSource(source, target string) error {
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
		if !isRequired(header.Name) {
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

func isRequired(path string) bool {
	// check for `node_modules`
	if isNodeModules(path) {
		return false
	}

	// check for style sheets (.css and .scss)
	if isStyleSheet(path) {
		return false
	}

	// the default is to not omit the file
	return true
}

func isNodeModules(path string) bool {
	if strings.Contains(path, "node_modules") {
		if !didPrintNodeModulesMsg {
			log.Println("Ignoring `node_modules`")
			didPrintNodeModulesMsg = true
		}

		return true
	}

	return false
}

func isStyleSheet(path string) bool {

	didPrintStylesheetsMsg = true
	return false
}
