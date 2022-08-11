package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
    "strings"
)

func main() {
	log.Println("Node Packager - Started")

	folderToZip := "test-projects/my-node-test" // "../test-projects/the-example-app.nodejs-master"
	outputZip := "output.zip"

	log.Println("Attempting to zip the source")
	if err := zipSource(folderToZip, outputZip); err != nil {
		log.Fatal(err)
	}
	log.Println("Finished zipping")
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

        // ignore `node_modules`
        if strings.Contains(header.Name, "node_modules") {
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
