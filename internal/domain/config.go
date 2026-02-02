package domain

type ResourceType string

const (
	ServiceType   ResourceType = "service"
	DirectoryType ResourceType = "directory"
)

type ResourceConfig struct {
	ID          string       `json:"id" yaml:"id"`
	Name        string       `json:"name" yaml:"name"`
	Description string       `json:"description" yaml:"description"`
	Type        ResourceType `json:"type" yaml:"type"`
	Path        string       `json:"path" yaml:"path"`
	Installable bool         `json:"installable" yaml:"installable"`

	// For services
	ServiceName string `json:"serviceName,omitempty" yaml:"service_name,omitempty"`
}

type Config struct {
	Version   string           `json:"version" yaml:"version"`
	Resources []ResourceConfig `json:"resources" yaml:"resources"`
}
