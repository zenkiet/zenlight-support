package app

import (
	"context"
	"sync"
	"time"
	"zenlight-support/internal/domain"
)

const (
	ScanInterval = 5 * time.Second
)

type ServiceWatcher struct {
	cfg        domain.Config
	mgr        domain.ResourceManager
	lastStatus sync.Map
	updates    chan []domain.ResourceStatus
}

func NewServiceWatcher(cfg domain.Config, mgr domain.ResourceManager) *ServiceWatcher {
	return &ServiceWatcher{
		cfg:     cfg,
		mgr:     mgr,
		updates: make(chan []domain.ResourceStatus, 10),
	}
}

func (sw *ServiceWatcher) Updates() <-chan []domain.ResourceStatus {
	return sw.updates
}

func (sw *ServiceWatcher) Start(ctx context.Context) {
	ticket := time.NewTicker(ScanInterval)
	defer ticket.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticket.C:
			sw.tick()
		}
	}
}

func (sw *ServiceWatcher) tick() {
	var changed []domain.ResourceStatus
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, svcCfg := range sw.cfg.Resources {
		if svcCfg.Type != domain.ServiceType {
			continue
		}

		wg.Add(1)
		go func(cfg domain.ResourceConfig) {
			defer wg.Done()

			status, err := sw.mgr.GetResourceState(cfg.ServiceName)
			if err != nil {
				status = domain.STOPPED
			}

			current := domain.ResourceStatus{
				ID:     cfg.ID,
				Status: status,
			}

			if status == domain.RUNNING {
				if metrics, err := sw.mgr.GetServiceMetrics(cfg.ServiceName); err == nil && metrics != nil {
					current.Metrics = metrics
				}
			}

			if sw.hasChanged(cfg.ID, current) {
				mu.Lock()
				changed = append(changed, current)
				mu.Unlock()
				sw.lastStatus.Store(cfg.ID, current)
			}
		}(svcCfg)
	}

	wg.Wait()

	if len(changed) > 0 {
		select {
		case sw.updates <- changed:
		default:
			// Skip if channel is full to avoid blocking
		}
	}
}

func (sw *ServiceWatcher) hasChanged(id string, current domain.ResourceStatus) bool {
	val, ok := sw.lastStatus.Load(id)
	if !ok {
		return true
	}
	old := val.(domain.ResourceStatus)

	// Check status
	if old.Status != current.Status {
		return true
	}

	// Check metrics
	if old.Metrics == nil && current.Metrics == nil {
		return false
	}
	if (old.Metrics == nil) != (current.Metrics == nil) {
		return true
	}

	// Check cpu and memory usage thresholds
	const cpuThreshold = 0.5             // 0.5%
	const memThreshold = 2 * 1024 * 1024 // 2 MB
	if absDiffFloat(old.Metrics.CPUUsage, current.Metrics.CPUUsage) > cpuThreshold {
		return true
	}

	if absDiffUint64(old.Metrics.MemUsage, current.Metrics.MemUsage) > memThreshold {
		return true
	}

	// Check PID
	if old.Metrics.PID != current.Metrics.PID {
		return true
	}

	return false
}

func absDiffFloat(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}

func absDiffUint64(a, b uint64) uint64 {
	if a > b {
		return a - b
	}
	return b - a
}
