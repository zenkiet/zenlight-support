package domain

type SQLConfig struct {
	ID          string `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Server      string `json:"server" yaml:"server"`
	Database    string `json:"database" yaml:"database"`
}

type Config struct {
	Version   string           `json:"version" yaml:"version"`
	Resources []ResourceConfig `json:"resources" yaml:"resources"`
	SQLConfig *SQLConfig       `json:"sqlConfig,omitempty" yaml:"sql_config,omitempty"`
}
