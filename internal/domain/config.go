package domain

type ResourceType string

const (
	ServiceType   ResourceType = "service"
	DirectoryType ResourceType = "directory"
	SQLScriptType ResourceType = "sqlscript"
)

type SQLConfig struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Server      string `json:"server" yaml:"server"`
	Database    string `json:"database" yaml:"database"`
}

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
	SQLConfig *SQLConfig       `json:"sqlConfig,omitempty" yaml:"sql_config,omitempty"`
}
