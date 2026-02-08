package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"zenlight-support/internal/domain"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) Install(id string, files []domain.InstallFileDTO) error {
	cfg, ok := a.itemMap[id]
	if !ok {
		return fmt.Errorf("service config not found for ID: %s", id)
	}

	serviceName := cfg.ServiceName
	targetPath := filepath.Clean(os.ExpandEnv(cfg.Path))

	wailsRuntime.LogInfo(a.Ctx, "Installing service files to: "+targetPath)

	// Stop service if running
	if cfg.Type == domain.ServiceType {
		if err := a.stopAndWait(serviceName); err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}
	}

	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	for _, file := range files {
		filePath := filepath.Join(targetPath, file.Name)
		wailsRuntime.LogInfo(a.Ctx, "Writing file: "+filePath)
		if err := os.WriteFile(filePath, file.Data, 0755); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file.Name, err)
		}
	}

	wailsRuntime.LogInfo(a.Ctx, "Service files installed for: "+serviceName)

	// Start service after installation
	if cfg.Type == domain.ServiceType {
		if err := a.startAndWait(serviceName); err != nil {
			return fmt.Errorf("failed to start service: %w", err)
		}
	}

	return nil
}

func (a *App) startAndWait(serviceName string) error {
	state, err := a.mgr.GetResourceState(serviceName)
	if err != nil {
		return err
	}

	if state == domain.RUNNING {
		return nil
	}

	if err := a.mgr.StartService(serviceName); err != nil {
		return err
	}

	timeout := time.After(30 * time.Second) // timeout after 30 seconds
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for service to start: %s", serviceName)
		case <-ticker.C:
			currentState, err := a.mgr.GetResourceState(serviceName)
			if err != nil {
				continue
			}
			if currentState == domain.RUNNING {
				return nil
			}
		}
	}
}

func (a *App) stopAndWait(serviceName string) error {
	state, err := a.mgr.GetResourceState(serviceName)
	if err != nil {
		return err
	}

	if state == domain.STOPPED {
		return nil
	}

	if err := a.mgr.StopService(serviceName); err != nil {
		return err
	}

	timeout := time.After(30 * time.Second) // timeout after 30 seconds
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for service to stop: %s", serviceName)
		case <-ticker.C:
			currentState, err := a.mgr.GetResourceState(serviceName)
			if err != nil {
				continue
			}
			if currentState == domain.STOPPED {
				return nil
			}
		}
	}
}
