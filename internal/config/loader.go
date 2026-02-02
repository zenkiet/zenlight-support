package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"window-service-watcher/internal/domain"
	"window-service-watcher/internal/repository"

	"github.com/google/uuid"
	"github.com/hashicorp/go-version"
)

const (
	fileName = "config.yaml"
)

func defaultConfig(appVer string) domain.Config {
	return domain.Config{
		Version: strings.TrimSpace(appVer),
		Resources: []domain.ResourceConfig{
			{
				ID:          uuid.NewString(),
				Name:        "Blogic Report Service",
				Description: "Generates and manages reports for Blogic View",
				Type:        domain.ServiceType,
				ServiceName: "BlogicReportService",
				Installable: true,
			},
			{
				ID:          uuid.NewString(),
				Name:        "POS Server",
				Description: "Point of Sale server for Blogic",
				Type:        domain.DirectoryType,
				Path:        "C:\\inetpub\\wwwroot\\BLogicService\\bin",
				Installable: true,
			},
			{
				ID:          uuid.NewString(),
				Name:        "Kiosk Web",
				Description: "Sources static files for Blogic Kiosk Web application.",
				Type:        domain.DirectoryType,
				Path:        "C:\\Program Files (x86)\\BLogic Systems\\BLogicKioskWeb",
				Installable: true,
			},
			{
				ID:          uuid.NewString(),
				Name:        "BlogicViewTaskScheduleHandler",
				Description: "Schedule send report for Blogic View",
				Type:        domain.DirectoryType,
				Path:        "C:\\Program Files (x86)\\BLogic Systems\\BLogic Service\\BLogicViewTaskScheduleHandler",
				Installable: true,
			},
		},
	}
}

func LoadConfig(appVer string) (*domain.Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, fileName)
	repo := repository.NewYamlConfigRepository(configPath)

	return ensureConfig(repo, appVer)
}

func ensureConfig(repo *repository.YamlConfigRepository, appVer string) (*domain.Config, error) {
	cfg, err := repo.Load()

	if os.IsNotExist(err) {
		newCfg := defaultConfig(appVer)
		if err := repo.Save(&newCfg); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return &newCfg, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	isOutdated, err := checkIsOutdated(cfg.Version, appVer)
	if err != nil {
		return nil, fmt.Errorf("failed to compare config versions: %w", err)
	}
	if isOutdated {
		if err := backupConfigFile(repo.Path, cfg.Version); err != nil {
			return nil, fmt.Errorf("failed to backup old config file: %w", err)
		}
		newCfg := defaultConfig(appVer)
		if err := repo.Save(&newCfg); err != nil {
			return nil, fmt.Errorf("failed to save new config: %w", err)
		}
		return &newCfg, nil
	}

	return cfg, nil
}

func checkIsOutdated(currentVerStr, appVerStr string) (bool, error) {
	if currentVerStr == "" {
		return true, nil
	}

	vCurrent, err := version.NewVersion(strings.TrimSpace(currentVerStr))
	if err != nil {
		return false, err
	}

	vApp, err := version.NewVersion(strings.TrimSpace(appVerStr))
	if err != nil {
		return false, err
	}

	return vApp.GreaterThan(vCurrent), nil
}

func backupConfigFile(path string, oldVer string) error {
	backupPath := fmt.Sprintf("%s.%s.bak", path, oldVer)
	return os.Rename(path, backupPath)
}
