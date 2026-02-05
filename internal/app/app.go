package app

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"zenlight-support/internal/domain"
	"zenlight-support/pkg/sql"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	Ctx     context.Context
	cfg     domain.Config
	mgr     domain.ResourceManager
	itemMap map[string]domain.ResourceConfig
	watcher *ServiceWatcher
	appVer  string
}

func NewApp(cfg domain.Config, mgr domain.ResourceManager, appVer string) *App {
	itemMap := make(map[string]domain.ResourceConfig)
	for _, item := range cfg.Resources {
		itemMap[item.ID] = item
	}

	return &App{
		cfg:     cfg,
		mgr:     mgr,
		itemMap: itemMap,
		watcher: NewServiceWatcher(cfg, mgr),
		appVer:  appVer,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.Ctx = ctx

	if err := a.mgr.Connect(); err != nil {
		wailsRuntime.LogError(a.Ctx, "Service Manager Connect Error: "+err.Error())
		return
	}

	go a.watcher.Start(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case updates := <-a.watcher.Updates():
				wailsRuntime.EventsEmit(a.Ctx, "services-update", updates)
			}
		}
	}()
}

func (a *App) Shutdown(ctx context.Context) {
	a.mgr.Disconnect()
}

func (a *App) GetServices() []domain.ResourceConfig {
	return a.filterByType(domain.ServiceType)
}

func (a *App) GetDirectories() []domain.ResourceConfig {
	return a.filterByType(domain.DirectoryType)
}

func (a *App) filterByType(rType domain.ResourceType) []domain.ResourceConfig {
	var result []domain.ResourceConfig
	for _, item := range a.cfg.Resources {
		if item.Type == rType {
			result = append(result, item)
		}
	}
	return result
}

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

func (a *App) GetResourceMetrics(id string) (*domain.ResourceMetrics, error) {
	cfg, ok := a.itemMap[id]
	if !ok {
		return nil, fmt.Errorf("service config not found for ID: %s", id)
	}

	if cfg.Type == domain.ServiceType {
		return a.mgr.GetServiceMetrics(cfg.ServiceName)
	}

	if cfg.Type == domain.DirectoryType {
		return a.mgr.GetDirectoryMetrics(cfg.Path)
	}

	return nil, fmt.Errorf("unsupported resource type for metrics: %s", id)
}

func (a *App) GetSqlConfig() (*domain.SQLConfig, error) {
	return a.cfg.SQLConfig, nil
}

func (a *App) ExecuteSQLScript(id string, script string) (*sql.Result, error) {
	return a.mgr.ExecuteSQLScript(a.cfg.SQLConfig.Server, a.cfg.SQLConfig.Database, script)
}
