package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"zenlight-support/internal/domain"
)

func (a *App) StartService(id string) error {
	cfg, ok := a.itemMap[id]
	if !ok {
		return fmt.Errorf("service config not found for ID: %s", id)
	}
	return a.mgr.StartService(cfg.ServiceName)
}

func (a *App) StopService(id string) error {
	cfg, ok := a.itemMap[id]
	if !ok {
		return fmt.Errorf("service config not found for ID: %s", id)
	}
	if cfg.Type != domain.ServiceType {
		return fmt.Errorf("resource is not a service: %s", id)
	}
	return a.mgr.StopService(cfg.ServiceName)
}

func (a *App) GetServiceStatus(id string) (domain.Status, error) {
	cfg, ok := a.itemMap[id]
	if !ok {
		return domain.STOPPED, fmt.Errorf("service config not found for ID: %s", id)
	}
	if cfg.Type != domain.ServiceType {
		return domain.STOPPED, fmt.Errorf("resource is not a service: %s", id)
	}

	return a.mgr.GetResourceState(cfg.ServiceName)
}

func (a *App) OpenExplorer(id string) error {
	cfg, ok := a.itemMap[id]
	if !ok {
		return fmt.Errorf("resource config not found for ID: %s", id)
	}

	targetPath := filepath.Clean(os.ExpandEnv(cfg.Path))
	info, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("resource path error: %w", err)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		if !info.IsDir() {
			cmd = exec.Command("explorer", "/select,", targetPath)
		} else {
			cmd = exec.Command("explorer", targetPath)
		}
	case "darwin":
		cmd = exec.Command("open", targetPath)
	default:
		dir := targetPath
		if !info.IsDir() {
			dir = filepath.Dir(targetPath)
		}
		cmd = exec.Command("xdg-open", dir)
	}

	return cmd.Start()
}
