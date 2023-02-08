package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
)

var AppVersion string = "0.0.0"
var latestRelease string = "https://github.com/fw10/veracode-javascript-packager/releases/latest"

func notifyOfUpdates() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", latestRelease, nil)
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var response map[string]interface{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	vCurrent, err := version.NewVersion(AppVersion)
	vLatest, err := version.NewVersion(response["tag_name"].(string))

	// check if a newer version exists
	if vCurrent.LessThan(vLatest) {
		color.HiYellow(fmt.Sprintf("Please upgrade to the latest version of this tool (%s) by visiting %s\n", response["tag_name"], latestRelease))
	}
}
