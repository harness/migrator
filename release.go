package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"strings"
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
	// The newest release includes both release & pre-release
	var latest GithubRelease
	// Latest stable stable
	var latestStableRelease GithubRelease
	isPreRelease := strings.Contains(Version, "beta")

	for _, v := range releases {
		if len(latest.TagName) == 0 {
			latest = v
		}
		if !v.Prerelease && len(latestStableRelease.TagName) == 0 {
			latestStableRelease = v
		}
	}
	if Version == latest.TagName {
		return
	}
	if !isPreRelease && latestStableRelease.TagName == Version {
		return
	}
	if !latest.Prerelease {
		printUpgradeMessage(Version, latest.TagName)
		return
	}
	if isPreRelease {
		printUpgradeMessage(Version, latest.TagName)
		return
	}
	if !isPreRelease && latestStableRelease.TagName != Version {
		printUpgradeMessage(Version, latestStableRelease.TagName)
	}
}

func printUpgradeMessage(from string, to string) {
	blue := color.New(color.FgHiBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("[%s] A new release of harness-upgrade available: %s -> %s\n", blue("notice"), red(from), green(to))
	fmt.Printf("%s\n", yellow("https://github.com/harness/migrator/releases/tag/"+to))
}
