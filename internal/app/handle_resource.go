package app

import (
	"fmt"
	"zenlight-support/internal/domain"

	"github.com/google/uuid"
)

func (a *App) GetServices() []domain.ResourceConfig {
	return a.filterByType(domain.ServiceType)
}

func (a *App) GetDirectories() []domain.ResourceConfig {
	return a.filterByType(domain.DirectoryType)
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

func (a *App) filterByType(rType domain.ResourceType) []domain.ResourceConfig {
	var result []domain.ResourceConfig
	for _, item := range a.cfg.Resources {
		if item.Type == rType {
			result = append(result, item)
		}
	}
	return result
}

func (a *App) SaveResource(resource domain.ResourceConfig) (*domain.ResourceConfig, error) {
	if resource.ID == "" {
		resource.ID = uuid.NewString()
		a.cfg.Resources = append(a.cfg.Resources, resource)
	} else {
		found := false
		for i, r := range a.cfg.Resources {
			if r.ID == resource.ID {
				a.cfg.Resources[i] = resource
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("resource not found for ID: %s", resource.ID)
		}
	}

	a.itemMap[resource.ID] = resource

	if err := a.repo.Save(&a.cfg); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	return &resource, nil
}

func (a *App) DeleteResource(id string) error {
	if _, ok := a.itemMap[id]; !ok {
		return fmt.Errorf("resource not found for ID: %s", id)
	}

	for i, r := range a.cfg.Resources {
		if r.ID == id {
			a.cfg.Resources = append(a.cfg.Resources[:i], a.cfg.Resources[i+1:]...)
			break
		}
	}

	delete(a.itemMap, id)

	if err := a.repo.Save(&a.cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
