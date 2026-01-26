package config

import (
	"fmt"
	"os"
	"path/filepath"
	"window-service-watcher/internal/domain"

	"github.com/google/uuid"
	"go.yaml.in/yaml/v4"
)

const fileName = "config.yaml"

func DefaultConfig() domain.Config {
	return domain.Config{
		Services: []domain.ServiceConfig{
			{
				ID:          uuid.New().String(),
				Name:        "Blogic Report Service",
				Description: "Generates and manages reports for Blogic applications.",
				ServiceName: "BlogicReportService",
				Path:        "C:\\Program Files (x86)\\BLogic Systems\\BLogic Service\\BLogicReportService",
				Installable: true,
			},
		},
	}
}

func LoadConfig() (*domain.Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, fileName)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return createConfig(configPath)
	}

	return loadConfig(configPath)
}

func createConfig(path string) (*domain.Config, error) {
	cfg := DefaultConfig()

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write default config file: %w", err)
	}

	return &cfg, nil
}

func loadConfig(path string) (*domain.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg domain.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
