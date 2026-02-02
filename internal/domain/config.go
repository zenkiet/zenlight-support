package domain

type ServiceConfig struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	ServiceName string `json:"serviceName" yaml:"service_name"`
	Path        string `json:"path" yaml:"path"`
	Installable bool   `json:"installable" yaml:"installable"`
}

type Config struct {
	Version  string          `json:"version" yaml:"version"`
	Services []ServiceConfig `json:"services" yaml:"services"`
}
