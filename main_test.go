package main

import (
	"archive/zip"
	"os"
	"reflect"
	"regexp"
	"sort"

	"testing"

	log "github.com/sirupsen/logrus"
)

func TestBefore(t *testing.T) {
	// check if the `./test-output` folder exists (otherwise, create it)
	if _, err := os.Stat("." + string(os.PathSeparator) + "test-output"); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir("."+string(os.PathSeparator)+"test-output", os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}
	}
}

// Integration test for `zipSource()` with `../sample-projects/sample-node-project`
func TestZipSourceWithNodeSample(t *testing.T) {
	sourcePath := "." + string(os.PathSeparator) + "sample-projects" + string(os.PathSeparator) + "sample-node-project"
	targetPath := "." + string(os.PathSeparator) + "test-output" + string(os.PathSeparator) + "test-output.zip"

	// generate the zip file and return a list of all its file names
	zipFileContents := generateZipAndReturnItsFiles(sourcePath, targetPath, "")

	// check if the output conforms with what we expected. To do this, we sort both the expected output and the actual output
	// and then compare them.
	expectedFilesInOutputZip := []string{
		"app.js", "package.json", "package-lock.json", "testimonials-no-tests" + string(os.PathSeparator) + "should-be-included.js",
		"distance" + string(os.PathSeparator) + "should-be-included.js", "building" + string(os.PathSeparator) + "something.js",
		"bower_components" + string(os.PathSeparator) + "bower.json", "bower_components" + string(os.PathSeparator) + "some-thing.js",
		"styles" + string(os.PathSeparator) + "blub.css2",
	}
	sort.Strings(expectedFilesInOutputZip)
	sort.Strings(zipFileContents)

	if !reflect.DeepEqual(zipFileContents, expectedFilesInOutputZip) {
		t.Error("Test failed!")
		t.Errorf("Got: %v", zipFileContents)
		t.Errorf("Expected: %v", expectedFilesInOutputZip)
	}
}

// Integration test for `zipSource()` with `../sample-projects/sample-node-project/` (note the trailing slash!). The reason
// for this test is that a trailing slash in the `-source` had lead to a bug that gave me quite some headache to figure out.
func TestZipSourceWithNodeSampleAndTrailingSlash(t *testing.T) {
	sourcePath := "." + string(os.PathSeparator) + "sample-projects" + string(os.PathSeparator) + "sample-node-project/"
	targetPath := "." + string(os.PathSeparator) + "test-output" + string(os.PathSeparator) + "test-output.zip"

	// generate the zip file and return a list of all its file names
	zipFileContents := generateZipAndReturnItsFiles(sourcePath, targetPath, "")

	// check if the output conforms with what we expected. To do this, we sort both the expected output and the actual output
	// and then compare them.
	expectedFilesInOutputZip := []string{
		"app.js", "package.json", "package-lock.json", "testimonials-no-tests" + string(os.PathSeparator) + "should-be-included.js",
		"distance" + string(os.PathSeparator) + "should-be-included.js", "building" + string(os.PathSeparator) + "something.js",
		"bower_components" + string(os.PathSeparator) + "bower.json", "bower_components" + string(os.PathSeparator) + "some-thing.js",
		"styles" + string(os.PathSeparator) + "blub.css2",
	}
	sort.Strings(expectedFilesInOutputZip)
	sort.Strings(zipFileContents)

	if !reflect.DeepEqual(zipFileContents, expectedFilesInOutputZip) {
		t.Error("Test failed!")
		t.Errorf("Got: %v", zipFileContents)
		t.Errorf("Expected: %v", expectedFilesInOutputZip)
	}
}

// Integration test for `zipSource()` with `../sample-projects/sample-node-project` and `-tests` provided
func TestZipSourceWithNodeSampleWithTestsFlag(t *testing.T) {
	sourcePath := "." + string(os.PathSeparator) + "sample-projects" + string(os.PathSeparator) + "sample-node-project"
	targetPath := "." + string(os.PathSeparator) + "test-output" + string(os.PathSeparator) + "test-output.zip"
	testsPath := "test"

	// generate the zip file and return a list of all its file names
	zipFileContents := generateZipAndReturnItsFiles(sourcePath, targetPath, testsPath)

	// check if the output conforms with what we expected. To do this, we sort both the expected output and the actual output
	// and then compare them.
	expectedFilesInOutputZip := []string{
		"app.js", "package.json", "package-lock.json", "testimonials-no-tests" + string(os.PathSeparator) + "should-be-included.js",
		"distance" + string(os.PathSeparator) + "should-be-included.js", "building" + string(os.PathSeparator) + "something.js",
		"bower_components" + string(os.PathSeparator) + "bower.json", "bower_components" + string(os.PathSeparator) + "some-thing.js",
		"styles" + string(os.PathSeparator) + "blub.css2", "e2e" + string(os.PathSeparator) + "some-more-test.js",
	}
	sort.Strings(expectedFilesInOutputZip)
	sort.Strings(zipFileContents)

	if !reflect.DeepEqual(zipFileContents, expectedFilesInOutputZip) {
		t.Error("Test failed!")
		t.Errorf("Got: %v", zipFileContents)
		t.Errorf("Expected: %v", expectedFilesInOutputZip)
	}
}

// Integration test for `zipSource()` with `../sample-projects/sample-angular-project`
func TestZipSourceWithAngularSample(t *testing.T) {
	sourcePath := "." + string(os.PathSeparator) + "sample-projects" + string(os.PathSeparator) + "sample-angular-project"
	targetPath := "." + string(os.PathSeparator) + "test-output" + string(os.PathSeparator) + "test-output.zip"

	// generate the zip file and return a list of all its file names
	zipFileContents := generateZipAndReturnItsFiles(sourcePath, targetPath, "")

	// check if the output conforms with what we expected. To do this, we sort both the expected output and the actual output
	// and then compare them.
	expectedFilesInOutputZip := []string{
		"package.json", "package-lock.json", "src" + string(os.PathSeparator) + "main.ts",
		"src" + string(os.PathSeparator) + "test.ts",
		"src" + string(os.PathSeparator) + "index.html",
		"src" + string(os.PathSeparator) + "environments" + string(os.PathSeparator) + "environment.prod.ts",
		"src" + string(os.PathSeparator) + "environments" + string(os.PathSeparator) + "environment.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "app.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "app-routing.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "settings" + string(os.PathSeparator) + "settings-routing.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "settings" + string(os.PathSeparator) + "settings.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "settings" + string(os.PathSeparator) + "settings.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "settings" + string(os.PathSeparator) + "settings.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "home" + string(os.PathSeparator) + "home-auth-resolver.service.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "home" + string(os.PathSeparator) + "home.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "home" + string(os.PathSeparator) + "home.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "home" + string(os.PathSeparator) + "home-routing.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "home" + string(os.PathSeparator) + "home.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "interceptors" + string(os.PathSeparator) + "http.token.interceptor.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "interceptors" + string(os.PathSeparator) + "index.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "models" + string(os.PathSeparator) + "user.model.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "models" + string(os.PathSeparator) + "comment.model.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "models" + string(os.PathSeparator) + "article-list-config.model.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "models" + string(os.PathSeparator) + "profile.model.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "models" + string(os.PathSeparator) + "index.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "models" + string(os.PathSeparator) + "errors.model.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "models" + string(os.PathSeparator) + "article.model.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "core.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "index.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "api.service.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "comments.service.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "profiles.service.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "tags.service.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "jwt.service.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "auth-guard.service.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "user.service.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "index.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "core" + string(os.PathSeparator) + "services" + string(os.PathSeparator) + "articles.service.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "auth" + string(os.PathSeparator) + "auth.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "auth" + string(os.PathSeparator) + "no-auth-guard.service.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "auth" + string(os.PathSeparator) + "auth-routing.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "auth" + string(os.PathSeparator) + "auth.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "auth" + string(os.PathSeparator) + "auth.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "list-errors.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "buttons" + string(os.PathSeparator) + "follow-button.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "buttons" + string(os.PathSeparator) + "follow-button.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "buttons" + string(os.PathSeparator) + "favorite-button.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "buttons" + string(os.PathSeparator) + "index.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "buttons" + string(os.PathSeparator) + "favorite-button.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "layout" + string(os.PathSeparator) + "header.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "layout" + string(os.PathSeparator) + "header.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "layout" + string(os.PathSeparator) + "footer.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "layout" + string(os.PathSeparator) + "index.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "layout" + string(os.PathSeparator) + "footer.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "article-helpers" + string(os.PathSeparator) + "article-list.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "article-helpers" + string(os.PathSeparator) + "article-preview.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "article-helpers" + string(os.PathSeparator) + "article-meta.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "article-helpers" + string(os.PathSeparator) + "index.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "article-helpers" + string(os.PathSeparator) + "article-meta.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "article-helpers" + string(os.PathSeparator) + "article-preview.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "article-helpers" + string(os.PathSeparator) + "article-list.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "show-authed.directive.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "shared.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "index.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "shared" + string(os.PathSeparator) + "list-errors.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "app.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "app.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "profile" + string(os.PathSeparator) + "profile-favorites.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "profile" + string(os.PathSeparator) + "profile.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "profile" + string(os.PathSeparator) + "profile-resolver.service.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "profile" + string(os.PathSeparator) + "profile-articles.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "profile" + string(os.PathSeparator) + "profile.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "profile" + string(os.PathSeparator) + "profile.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "profile" + string(os.PathSeparator) + "profile-routing.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "profile" + string(os.PathSeparator) + "profile-favorites.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "profile" + string(os.PathSeparator) + "profile-articles.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "index.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "article" + string(os.PathSeparator) + "article.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "article" + string(os.PathSeparator) + "article-comment.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "article" + string(os.PathSeparator) + "article-comment.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "article" + string(os.PathSeparator) + "article.component.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "article" + string(os.PathSeparator) + "article.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "article" + string(os.PathSeparator) + "markdown.pipe.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "article" + string(os.PathSeparator) + "article-resolver.service.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "article" + string(os.PathSeparator) + "article-routing.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "editor" + string(os.PathSeparator) + "editor.component.html",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "editor" + string(os.PathSeparator) + "editor-routing.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "editor" + string(os.PathSeparator) + "editable-article-resolver.service.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "editor" + string(os.PathSeparator) + "editor.module.ts",
		"src" + string(os.PathSeparator) + "app" + string(os.PathSeparator) + "editor" + string(os.PathSeparator) + "editor.component.ts",
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

	// read the output zip file (e.g. `./test-output/test-output.zip`) into memory
	zipReader := readZip(targetPath)

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
func readZip(zipPath string) *zip.ReadCloser {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	return r
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
