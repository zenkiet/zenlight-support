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

const RepoURL = "https://api.github.com/repos/zenkiet/zenlight-support/releases/latest"

type UpdateInfo struct {
	Available    bool   `json:"available"`
	CurrentVer   string `json:"currentVersion"`
	Build        string `json:"build,omitempty"`
	LatestVer    string `json:"latestVersion"`
	ReleaseNotes string `json:"releaseNotes"`
	DownloadURL  string `json:"downloadUrl"`
	Error        string `json:"error,omitempty"`
}

type asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type Release struct {
	TagName   string  `json:"tag_name"`
	CreatedAt string  `json:"created_at"`
	Body      string  `json:"body"`
	Assets    []asset `json:"assets"`
}

func (a *App) CheckUpdate() UpdateInfo {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", RepoURL, nil)
	req.Header.Set("Authorization", "token "+GH_TOKEN)
	if err != nil {
		return errorUpdateInfo(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errorUpdateInfo(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errorUpdateInfo(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return errorUpdateInfo(fmt.Errorf("failed to decode release info: %w", err))
	}

	current, _ := version.NewVersion(a.appVer)
	latest, err := version.NewVersion(release.TagName)
	if err != nil {
		return errorUpdateInfo(fmt.Errorf("invalid tag format in release: %s", release.TagName))
	}

	if latest.LessThanOrEqual(current) {
		return UpdateInfo{Available: false, CurrentVer: a.appVer, LatestVer: release.TagName, Build: release.CreatedAt}
	}

	downloadURL := findAsset(release.Assets, "zenlight-support.exe")
	if downloadURL == "" {
		return errorUpdateInfo(fmt.Errorf("no suitable asset found for download"))
	}

	return UpdateInfo{
		Available:    true,
		CurrentVer:   a.appVer,
		Build:        release.CreatedAt,
		LatestVer:    release.TagName,
		ReleaseNotes: release.Body,
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

	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
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

func errorUpdateInfo(err error) UpdateInfo {
	return UpdateInfo{Available: false, Error: err.Error(), CurrentVer: "v0.0.0"}
}

func findAsset(assets []asset, appName string) string {
	targetExt := ".exe"
	targetArch := runtime.GOARCH

	for _, asset := range assets {
		name := strings.ToLower(asset.Name)
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
		return asset.BrowserDownloadURL
	}
	return ""
}
