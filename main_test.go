package main

import (
	"archive/zip"
	"reflect"
	"regexp"
	"sort"

	"testing"

	log "github.com/sirupsen/logrus"
)

// Integration test for `zipSource()` with `../sample-projects/sample-node-project`
func TestZipSourceWithNodeSample(t *testing.T) {
	sourcePath := "./sample-projects/sample-node-project"
	targetPath := "./test-output/test-output.zip"

	// generate the zip file and return a list of all its file names
	zipFileContents := generateZipAndReturnItsFiles(sourcePath, targetPath, "")

	// check if the output conforms with what we expected. To do this, we sort both the expected output and the actual output
	// and then compare them.
	expectedFilesInOutputZip := []string{
		"app.js", "package.json", "package-lock.json", "testimonials-no-tests/should-be-included.js",
		"distance/should-be-included.js", "building/something.js", "bower_components/bower.json",
		"bower_components/some-thing.js", "styles/blub.css2",
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
	sourcePath := "./sample-projects/sample-node-project/"
	targetPath := "./test-output/test-output.zip"

	// generate the zip file and return a list of all its file names
	zipFileContents := generateZipAndReturnItsFiles(sourcePath, targetPath, "")

	// check if the output conforms with what we expected. To do this, we sort both the expected output and the actual output
	// and then compare them.
	expectedFilesInOutputZip := []string{
		"app.js", "package.json", "package-lock.json", "testimonials-no-tests/should-be-included.js",
		"distance/should-be-included.js", "building/something.js", "bower_components/bower.json",
		"bower_components/some-thing.js", "styles/blub.css2",
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
	sourcePath := "./sample-projects/sample-node-project"
	targetPath := "./test-output/test-output.zip"
	testsPath := "test"

	// generate the zip file and return a list of all its file names
	zipFileContents := generateZipAndReturnItsFiles(sourcePath, targetPath, testsPath)

	// check if the output conforms with what we expected. To do this, we sort both the expected output and the actual output
	// and then compare them.
	expectedFilesInOutputZip := []string{
		"app.js", "package.json", "package-lock.json", "testimonials-no-tests/should-be-included.js",
		"distance/should-be-included.js", "building/something.js", "bower_components/bower.json",
		"bower_components/some-thing.js", "styles/blub.css2", "e2e/some-more-test.js",
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
	sourcePath := "./sample-projects/sample-angular-project"
	targetPath := "./test-output/test-output.zip"

	// generate the zip file and return a list of all its file names
	zipFileContents := generateZipAndReturnItsFiles(sourcePath, targetPath, "")

	// check if the output conforms with what we expected. To do this, we sort both the expected output and the actual output
	// and then compare them.
	expectedFilesInOutputZip := []string{
		"package.json", "package-lock.json", "src/main.ts", "src/test.ts", "src/index.html",
		"src/environments/environment.prod.ts", "src/environments/environment.ts",
		"src/app/app.component.html", "src/app/app-routing.module.ts", "src/app/settings/settings-routing.module.ts",
		"src/app/settings/settings.component.ts", "src/app/settings/settings.module.ts", "src/app/settings/settings.component.html",
		"src/app/home/home-auth-resolver.service.ts", "src/app/home/home.component.ts", "src/app/home/home.module.ts",
		"src/app/home/home-routing.module.ts", "src/app/home/home.component.html",
		"src/app/core/interceptors/http.token.interceptor.ts", "src/app/core/interceptors/index.ts",
		"src/app/core/models/user.model.ts", "src/app/core/models/comment.model.ts",
		"src/app/core/models/article-list-config.model.ts", "src/app/core/models/profile.model.ts",
		"src/app/core/models/index.ts", "src/app/core/models/errors.model.ts",
		"src/app/core/models/article.model.ts", "src/app/core/core.module.ts",
		"src/app/core/index.ts", "src/app/core/services/api.service.ts",
		"src/app/core/services/comments.service.ts", "src/app/core/services/profiles.service.ts",
		"src/app/core/services/tags.service.ts", "src/app/core/services/jwt.service.ts",
		"src/app/core/services/auth-guard.service.ts", "src/app/core/services/user.service.ts",
		"src/app/core/services/index.ts", "src/app/core/services/articles.service.ts",
		"src/app/auth/auth.component.ts", "src/app/auth/no-auth-guard.service.ts",
		"src/app/auth/auth-routing.module.ts", "src/app/auth/auth.module.ts",
		"src/app/auth/auth.component.html", "src/app/shared/list-errors.component.html",
		"src/app/shared/buttons/follow-button.component.ts", "src/app/shared/buttons/follow-button.component.html",
		"src/app/shared/buttons/favorite-button.component.html", "src/app/shared/buttons/index.ts",
		"src/app/shared/buttons/favorite-button.component.ts", "src/app/shared/layout/header.component.html",
		"src/app/shared/layout/header.component.ts", "src/app/shared/layout/footer.component.ts",
		"src/app/shared/layout/index.ts", "src/app/shared/layout/footer.component.html",
		"src/app/shared/article-helpers/article-list.component.ts", "src/app/shared/article-helpers/article-preview.component.ts",
		"src/app/shared/article-helpers/article-meta.component.ts", "src/app/shared/article-helpers/index.ts",
		"src/app/shared/article-helpers/article-meta.component.html", "src/app/shared/article-helpers/article-preview.component.html",
		"src/app/shared/article-helpers/article-list.component.html", "src/app/shared/show-authed.directive.ts",
		"src/app/shared/shared.module.ts", "src/app/shared/index.ts",
		"src/app/shared/list-errors.component.ts", "src/app/app.module.ts", "src/app/app.component.ts",
		"src/app/profile/profile-favorites.component.ts", "src/app/profile/profile.component.html",
		"src/app/profile/profile-resolver.service.ts", "src/app/profile/profile-articles.component.html",
		"src/app/profile/profile.module.ts", "src/app/profile/profile.component.ts",
		"src/app/profile/profile-routing.module.ts", "src/app/profile/profile-favorites.component.html",
		"src/app/profile/profile-articles.component.ts", "src/app/index.ts", "src/app/article/article.component.html",
		"src/app/article/article-comment.component.ts", "src/app/article/article-comment.component.html",
		"src/app/article/article.component.ts", "src/app/article/article.module.ts",
		"src/app/article/markdown.pipe.ts", "src/app/article/article-resolver.service.ts",
		"src/app/article/article-routing.module.ts", "src/app/editor/editor.component.html",
		"src/app/editor/editor-routing.module.ts", "src/app/editor/editable-article-resolver.service.ts",
		"src/app/editor/editor.module.ts", "src/app/editor/editor.component.ts",
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
