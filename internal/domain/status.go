package domain

type Status int8

const (
	STOPPED Status = 1 // Running
	RUNNING Status = 4 // Stopped
)

type ResourceStatus struct {
	ID      string           `json:"id"`
	Status  Status           `json:"status"`
	Metrics *ResourceMetrics `json:"metrics,omitempty"`
}

type ResourceMetrics struct {
	// --- Service ---
	PID        uint32  `json:"pid"`
	CreateTime int64   `json:"createTime"`
	CPUUsage   float64 `json:"cpu"`
	MemUsage   uint64  `json:"mem"`

	// --- Directory ---
	TotalSize    int64 `json:"totalSize,omitempty"`
	LastModified int64 `json:"lastModified,omitempty"`
}
