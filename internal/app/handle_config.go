package app

import (
	"fmt"
	"os"
	"time"
	"zenlight-support/internal/domain"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"go.yaml.in/yaml/v3"
)

func (a *App) GetSqlConfig() (*domain.SQLConfig, error) {
	return a.cfg.SQLConfig, nil
}

func (a *App) GetConfigPath() string {
	return a.repo.Path
}

func (a *App) ExportConfig() (string, error) {
	cfg, err := a.repo.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}

	savePath, err := wailsRuntime.SaveFileDialog(a.Ctx, wailsRuntime.SaveDialogOptions{
		DefaultFilename: fmt.Sprintf("config-backup-%s.yaml", time.Now().Format(time.DateOnly)),
		Title:           "Save Config Backup",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "YAML Files", Pattern: "*.yaml;*.yml"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("dialog error: %w", err)
	}
	if savePath == "" {
		return "", nil
	}

	if err := os.WriteFile(savePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write backup: %w", err)
	}

	return savePath, nil
}

func (a *App) ImportConfig() error {
	openPath, err := wailsRuntime.OpenFileDialog(a.Ctx, wailsRuntime.OpenDialogOptions{
		Title: "Select Config File to Restore",
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "YAML Files", Pattern: "*.yaml;*.yml"},
		},
	})
	if err != nil {
		return fmt.Errorf("dialog error: %w", err)
	}
	if openPath == "" {
		return nil
	}

	content, err := os.ReadFile(openPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var cfg domain.Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return fmt.Errorf("invalid config format: %w", err)
	}

	backupPath := fmt.Sprintf("%s.%s.bak", a.repo.Path, time.Now().Format(time.DateTime))
	currentData, err := os.ReadFile(a.repo.Path)
	if err == nil {
		_ = os.WriteFile(backupPath, currentData, 0644)
	}

	if err := a.repo.Save(&cfg); err != nil {
		return fmt.Errorf("failed to save imported config: %w", err)
	}

	a.cfg = cfg
	a.itemMap = make(map[string]domain.ResourceConfig)
	for _, item := range cfg.Resources {
		a.itemMap[item.ID] = item
	}

	return nil
}
