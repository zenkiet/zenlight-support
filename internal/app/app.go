package app

import (
	"context"
	"zenlight-support/internal/domain"
	"zenlight-support/internal/repository"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	Ctx     context.Context
	cfg     domain.Config
	mgr     domain.ResourceManager
	repo    *repository.YamlConfigRepository
	itemMap map[string]domain.ResourceConfig
	watcher *ServiceWatcher
	appVer  string
}

func NewApp(cfg domain.Config, mgr domain.ResourceManager, repo *repository.YamlConfigRepository, appVer string) *App {
	itemMap := make(map[string]domain.ResourceConfig)
	for _, item := range cfg.Resources {
		itemMap[item.ID] = item
	}

	return &App{
		cfg:     cfg,
		mgr:     mgr,
		repo:    repo,
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
