package domain

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
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type ServiceMetrics struct {
	PID        uint32  `json:"pid"`
	CreateTime int64   `json:"createTime"`
	CPUUsage   float64 `json:"cpu"`
	MemUsage   uint64  `json:"mem"`
}

type ServiceStatus struct {
	ID      string          `json:"id"`
	Status  Status          `json:"status"`
	Metrics *ServiceMetrics `json:"metrics,omitempty"`
}

type ServiceManager interface {
	Connect() error
	Disconnect() error

	GetServiceState(serviceName string) (Status, error)
	GetServiceMetrics(serviceName string) (*ServiceMetrics, error)

	StartService(serviceName string) error
	StopService(serviceName string) error

	// StartLogWatcher(filePath string, onLog func(string), onError func(error))
	// StopLogWatcher()

}
