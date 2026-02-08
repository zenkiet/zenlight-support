package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/minio/selfupdate"
)

var GH_TOKEN = "" // replace with your GitHub token

const repoURL = "https://api.github.com/repos/zenkiet/zenlight-support/releases/latest"

type UpdateInfo struct {
	Available    bool   `json:"available"`
	CurrentVer   string `json:"currentVersion"`
	Build        string `json:"build,omitempty"`
	LatestVer    string `json:"latestVersion"`
	ReleaseNotes string `json:"releaseNotes"`
	DownloadURL  string `json:"downloadUrl"`
	Error        string `json:"error,omitempty"`
}

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type release struct {
	TagName   string         `json:"tag_name"`
	CreatedAt string         `json:"created_at"`
	Body      string         `json:"body"`
	Assets    []releaseAsset `json:"assets"`
}

func (a *App) CheckUpdate() UpdateInfo {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", repoURL, nil)
	if err != nil {
		return errorUpdateInfo(err, a.appVer)
	}
	req.Header.Set("Authorization", "token "+GH_TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errorUpdateInfo(err, a.appVer)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errorUpdateInfo(fmt.Errorf("unexpected status code: %d", resp.StatusCode), a.appVer)
	}

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return errorUpdateInfo(fmt.Errorf("failed to decode release info: %w", err), a.appVer)
	}

	current, _ := version.NewVersion(a.appVer)
	latest, err := version.NewVersion(rel.TagName)
	if err != nil {
		return errorUpdateInfo(fmt.Errorf("invalid tag format in release: %s", rel.TagName), a.appVer)
	}

	if latest.LessThanOrEqual(current) {
		return UpdateInfo{Available: false, CurrentVer: a.appVer, LatestVer: rel.TagName, Build: rel.CreatedAt}
	}

	downloadURL := findAsset(rel.Assets, "zenlight-support.exe")
	if downloadURL == "" {
		return errorUpdateInfo(fmt.Errorf("no suitable asset found for download"), a.appVer)
	}

	return UpdateInfo{
		Available:    true,
		CurrentVer:   a.appVer,
		Build:        rel.CreatedAt,
		LatestVer:    rel.TagName,
		ReleaseNotes: rel.Body,
		DownloadURL:  downloadURL,
	}
}

func (a *App) DoUpdate(downloadURL string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("only supported on windows")
	}

	resp, err := http.Get(downloadURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if err := selfupdate.Apply(resp.Body, selfupdate.Options{}); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

func (a *App) RestartApp() {
	selfExecutable, _ := os.Executable()
	cmd := exec.Command(selfExecutable)
	cmd.Start()
	os.Exit(0)
}

func errorUpdateInfo(err error, appVer string) UpdateInfo {
	return UpdateInfo{Available: false, Error: err.Error(), CurrentVer: appVer}
}

func findAsset(assets []releaseAsset, appName string) string {
	targetExt := ".exe"
	targetArch := runtime.GOARCH

	for _, a := range assets {
		name := strings.ToLower(a.Name)
		if !strings.Contains(name, strings.ToLower(appName)) {
			continue
		}
		if !strings.HasSuffix(name, targetExt) {
			continue
		}
		if strings.Contains(name, "amd64") && targetArch != "amd64" {
			continue
		}
		if strings.Contains(name, "arm64") && targetArch != "arm64" {
			continue
		}
		return a.BrowserDownloadURL
	}
	return ""
}
