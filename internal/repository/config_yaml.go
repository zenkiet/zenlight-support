package repository

import (
	"bytes"
	"os"
	"sync"
	"zenlight-support/internal/domain"

	"go.yaml.in/yaml/v3"
)

type YamlConfigRepository struct {
	Path string
	mu   sync.RWMutex
}

func NewYamlConfigRepository(path string) *YamlConfigRepository {
	return &YamlConfigRepository{Path: path}
}

func (r *YamlConfigRepository) Load() (*domain.Config, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile(r.Path)
	if err != nil {
		return nil, err
	}

	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	var cfg domain.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (r *YamlConfigRepository) Save(cfg *domain.Config) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(r.Path, data, 0644)
}
