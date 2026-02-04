package domain

import "window-service-watcher/pkg/sql"

type Status int8

// don't care other states
const (
	// ERROR Status = iota
	STOPPED = 1
	// START_PENDING
	// STOP_PENDING
	RUNNING = 4
	// CONTINUE_PENDING
	// PAUSE_PENDING
	// PAUSED
)

type InstallFileDTO struct {
	Name      string `json:"name"`
	Data      []byte `json:"data"`
	Extension string `json:"extension"`
}

type ResourceMetrics struct {
	// --- Service ---
	PID        uint32  `json:"pid"`
	CreateTime int64   `json:"createTime"`
	CPUUsage   float64 `json:"cpu"`
	MemUsage   uint64  `json:"mem"`

	// --- Directory ---
	TotalSize int64 `json:"totalSize,omitempty"`
	// FileCount    int   `json:"fileCount,omitempty"`
	LastModified int64 `json:"lastModified,omitempty"`
}

type ResourceStatus struct {
	ID      string           `json:"id"`
	Status  Status           `json:"status"`
	Metrics *ResourceMetrics `json:"metrics,omitempty"`
}

type ResourceManager interface {
	Connect() error
	Disconnect() error

	GetResourceState(resourceName string) (Status, error)
	GetServiceMetrics(resourceName string) (*ResourceMetrics, error)
	GetDirectoryMetrics(path string) (*ResourceMetrics, error)

	StartService(serviceName string) error
	StopService(serviceName string) error

	ExecuteSQLScript(server, database, script string) (*sql.Result, error)

	// StartLogWatcher(filePath string, onLog func(string), onError func(error))
	// StopLogWatcher()

}
