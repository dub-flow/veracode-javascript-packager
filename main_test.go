package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"sort"

	"testing"

	log "github.com/sirupsen/logrus"
)

// Integration test for `zipSource()` with one of our sample apps `../sample-projects/sample-node-project`
func TestZipSourceWithNodeSample(t *testing.T) {
	sourcePath := "./sample-projects/sample-node-project"
	targetPath := "./test-output/test-output.zip"

	// generate the zip file and return a list of all its file names
	zipFileContents := generateZipAndReturnItsFiles(sourcePath, targetPath, "")

	// check if the output conforms with what we expected. To do this, we sort both the expected output and the actual output
	// and then compare them.
	expectedFilesInOutputZip := []string{
		"/app.js", "/package.json", "/package-lock.json", "/testimonials-no-tests/should-be-included.js",
		"/distance/should-be-included.js", "/building/something.js", "/bower_components/bower.json",
		"/bower_components/some-thing.js", "/styles/blub.css2",
	}
	sort.Strings(expectedFilesInOutputZip)
	sort.Strings(zipFileContents)

	if !reflect.DeepEqual(zipFileContents, expectedFilesInOutputZip) {
		t.Error("Test failed!")
		t.Errorf("Got: %v", zipFileContents)
		t.Errorf("Expected: %v", expectedFilesInOutputZip)
	}
}

func generateZipAndReturnItsFiles(sourcePath string, targetPath string, testsPath string) []string {
	// generate the zip file, and omit all non-required files
	if err := zipSource(sourcePath, targetPath, testsPath); err != nil {
		log.Fatal(err)
	}

	// read the output zip file (`./test-output/test-output.zip`) into memory
	zipReader, err := readZipFileIntoMemory(targetPath)
	if err != nil {
		log.Error(err)
	}

	// iterate over all the files from the zip archive and get all a list of all files (similar to the output of the `tree` command)
	var zipFileContents []string
	for _, zipFile := range zipReader.File {
		// `zipFile.Name` contains all the paths within the zip file. It contains an element for each file (e.g. `/src/app.js`),
		// but also an element for each folder (e.g. `/src/`). We only care about files and thus, omit every path that does not
		// belong to a file (i.e., does not end in `.<something>`)
		if !isPathAFile(zipFile.Name) {
			log.Info("Omitted path: ", zipFile.Name)
			continue
		}

		// we keep everything that is a file
		log.Info("Keeping file:", zipFile.Name)
		zipFileContents = append(zipFileContents, zipFile.Name)
	}

	return zipFileContents
}

// reads the zip file for the provided `zipPath` into memory and returns it
func readZipFileIntoMemory(zipPath string) (*zip.Reader, error) {
	outputZipFile, err := os.ReadFile(zipPath)
	if err != nil {
		return nil, fmt.Errorf("Unable to read file: %v", err)
	}

	body, err := ioutil.ReadAll(bytes.NewReader(outputZipFile))
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, err
	}

	return zipReader, nil
}

// check if a provided path (e.g. `/some/thing/app.js` or `/some/thing/`) belongs to a file
func isPathAFile(path string) bool {
	isFileRegex := ".+\\.[a-zA-Z0-9]{1,8}$"
	doesMatch, err := regexp.MatchString(isFileRegex, path)
	if err != nil {
		log.Error(err)
	}

	return doesMatch
}
