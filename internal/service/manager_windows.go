//go:build windows

package service

import (
	"context"
	"fmt"
	"sync"

	"window-service-watcher/internal/domain"

	"github.com/nxadm/tail"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"
)

type WindowsManager struct {
	ServiceName string
	mgr         *mgr.Mgr
	logCancel   context.CancelFunc
	logMutex    sync.Mutex
}

// CheckStatus implements [domain.ServiceManager].
func (w *WindowsManager) CheckStatus() (domain.ServiceStatus, error) {
	if w.mgr == nil {
		return domain.ServiceStatus{Status: "Manager Disconnected"}, fmt.Errorf("service manager not connected")
	}

	s, err := w.mgr.OpenService(w.ServiceName)
	if err != nil {
		return domain.ServiceStatus{
			Name:      w.ServiceName,
			Status:    "Not Installed",
			IsHealthy: false,
		}, nil
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return domain.ServiceStatus{
			Name:      w.ServiceName,
			Status:    "Error Querying Status",
			IsHealthy: false,
		}, nil
	}

	var state string
	isHealthy := false
	switch status.State {
	case windows.SERVICE_STOPPED:
		state = "Stopped"
	case windows.SERVICE_START_PENDING:
		state = "Starting"
	case windows.SERVICE_STOP_PENDING:
		state = "Stopping"
	case windows.SERVICE_RUNNING:
		state = "Running"
		isHealthy = true
	case windows.SERVICE_CONTINUE_PENDING:
		state = "Resuming"
	case windows.SERVICE_PAUSE_PENDING:
		state = "Pausing"
	case windows.SERVICE_PAUSED:
		state = "Paused"
	default:
		state = "Unknown"
	}

	return domain.ServiceStatus{
		Name:      w.ServiceName,
		Status:    state,
		IsHealthy: isHealthy,
	}, nil
}

// Connect implements [domain.ServiceManager].
func (w *WindowsManager) Connect() error {
	m, error := mgr.Connect()
	if error != nil {
		return error
	}
	w.mgr = m
	return nil
}

// Disconnect implements [domain.ServiceManager].
func (w *WindowsManager) Disconnect() error {
	w.StopLogWatcher()
	if w.mgr != nil {
		return w.mgr.Disconnect()
	}
	return nil
}

// StartLogWatcher implements [domain.ServiceManager].
func (w *WindowsManager) StartLogWatcher(filePath string, onLog func(string), onError func(error)) {
	if filePath == "" {
		return
	}

	w.StopLogWatcher()

	ctx, cancel := context.WithCancel(context.Background())

	w.logMutex.Lock()
	w.logCancel = cancel
	w.logMutex.Unlock()

	go func(ctx context.Context) {
		t, err := tail.TailFile(filePath, tail.Config{
			Follow: true,
			ReOpen: true,
			Poll:   true, // window often use polling
			Logger: tail.DiscardingLogger,
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

// StopLogWatcher implements [domain.ServiceManager].
func (w *WindowsManager) StopLogWatcher() {
	w.logMutex.Lock()
	defer w.logMutex.Unlock()

	if w.logCancel != nil {
		w.logCancel()
		w.logCancel = nil
	}
}

func NewManager() domain.ServiceManager {
	return &WindowsManager{
		ServiceName: "BlogicReportService",
	}
}
