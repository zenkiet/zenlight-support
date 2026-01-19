//go:build !windows

package service

import (
	"context"
	"sync"

	"window-service-watcher/internal/domain"

	"github.com/nxadm/tail"
)

type MockManager struct {
	logCancel context.CancelFunc
	logMutex  sync.Mutex
}

// StartLogWatcher implements [domain.ServiceManager].
func (m *MockManager) StartLogWatcher(filePath string, onLog func(string), onError func(error)) {
	if filePath == "" {
		return
	}

	m.StopLogWatcher()
	ctx, cancel := context.WithCancel(context.Background())

	m.logMutex.Lock()
	m.logCancel = cancel
	m.logMutex.Unlock()

	go func(ctx context.Context) {
		t, err := tail.TailFile(filePath, tail.Config{
			Follow: true,
			ReOpen: true,
			Poll:   true, // window often use polling
		})
		if err != nil {
			onError(err)
			return
		}
		defer func() {
			t.Cleanup()
			t.Stop()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case line, ok := <-t.Lines:
				if !ok {
					return
				}
				if line.Err != nil {
					onError(line.Err)
					continue
				}
				onLog(line.Text)
			}
		}
	}(ctx)
}

// CheckStatus implements [domain.ServiceManager].
func (m *MockManager) CheckStatus() (domain.ServiceStatus, error) {
	return domain.ServiceStatus{
		Name:      "MockService",
		Status:    "Running",
		IsHealthy: true,
	}, nil
}

// Connect implements [domain.ServiceManager].
func (m *MockManager) Connect() error {
	return nil
}

// Disconnect implements [domain.ServiceManager].
func (m *MockManager) Disconnect() error {
	return nil
}

// StopLogWatcher implements [domain.ServiceManager].
func (m *MockManager) StopLogWatcher() {
	m.logMutex.Lock()
	defer m.logMutex.Unlock()

	if m.logCancel != nil {
		m.logCancel()
		m.logCancel = nil
	}
}

func NewManager() domain.ServiceManager {
	return &MockManager{}
}
