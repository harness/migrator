package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
)

type GithubRelease struct {
	Prerelease bool   `json:"prerelease"`
	TagName    string `json:"tag_name"`
}

func CheckGithubForReleases() {
	if Version == "development" {
		return
	}
	resp, err := http.Get("https://api.github.com/repos/harness/migrator/releases")
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return
	}

	releases := []GithubRelease{}
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return
	}
	version := Version
	for _, v := range releases {
		if !v.Prerelease {
			version = v.TagName
			break
		}
	}
	if version != Version {
		blue := color.New(color.FgBlue).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		fmt.Printf("[%s] A new release of harness-upgrade available: %s -> %s\n", blue("notice"), red(Version), green(version))
	}
}
