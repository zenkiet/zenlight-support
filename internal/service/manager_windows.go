//go:build windows

package service

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
	"window-service-watcher/internal/domain"

	"github.com/shirou/gopsutil/v4/process"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/mgr"
)

type processHandle struct {
	proc       *process.Process
	lastPID    int32
	lastUpdate time.Time
}

type WindowsManager struct {
	mgr          *mgr.Mgr
	mu           sync.RWMutex
	processCache map[string]*processHandle
}

// Connect implements [domain.ServiceManager].
func (w *WindowsManager) Connect() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.mgr != nil {
		return nil
	}
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	w.mgr = m
	return nil
}

// Disconnect implements [domain.ServiceManager].
func (w *WindowsManager) Disconnect() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.mgr == nil {
		return nil
	}
	err := w.mgr.Disconnect()
	w.mgr = nil
	w.processCache = make(map[string]*processHandle)
	return err
}

// GetServiceMetrics implements [domain.ServiceManager].
func (w *WindowsManager) GetServiceMetrics(serviceName string) (*domain.ServiceMetrics, error) {
	w.mu.RLock()
	m := w.mgr
	w.mu.RUnlock()

	if m == nil {
		return nil, fmt.Errorf("not connected")
	}

	s, err := m.OpenService(serviceName)
	if err != nil {
		return nil, err
	}
	defer s.Close()

	var status windows.SERVICE_STATUS_PROCESS
	var bytesReturned uint32
	err = windows.QueryServiceStatusEx(s.Handle, windows.SC_STATUS_PROCESS_INFO, (*byte)(unsafe.Pointer(&status)), uint32(unsafe.Sizeof(status)), &bytesReturned)
	if err != nil {
		return nil, err
	}

	if status.ProcessId == 0 {
		w.mu.Lock()
		delete(w.processCache, serviceName)
		w.mu.Unlock()
		return nil, nil
	}

	currentPID := int32(status.ProcessId)
	metrics := &domain.ServiceMetrics{PID: status.ProcessId}

	w.mu.Lock()
	handle, exists := w.processCache[serviceName]

	if !exists || handle.lastPID != currentPID {
		p, err := process.NewProcess(currentPID)
		if err != nil {
			w.mu.Unlock()
			return nil, err
		}
		handle = &processHandle{proc: p, lastPID: currentPID}
		w.processCache[serviceName] = handle
	}
	w.mu.Unlock()

	if cpu, err := handle.proc.CPUPercent(); err == nil {
		metrics.CPUUsage = cpu
	}

	if mem, err := handle.proc.MemoryInfo(); err == nil {
		metrics.MemUsage = mem.RSS
	}

	if metrics.CreateTime == 0 {
		if createTime, err := handle.proc.CreateTime(); err == nil {
			metrics.CreateTime = createTime
		}
	}

	return metrics, nil
}

// GetServiceState implements [domain.ServiceManager].
func (w *WindowsManager) GetServiceState(serviceName string) (domain.Status, error) {
	w.mu.RLock()
	m := w.mgr
	w.mu.RUnlock()

	if m == nil {
		return domain.STOPPED, fmt.Errorf("not connected")
	}

	s, err := m.OpenService(serviceName)
	if err != nil {
		return domain.STOPPED, err
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return domain.STOPPED, err
	}

	return domain.Status(status.State), nil
}

// StartService implements [domain.ServiceManager].
func (w *WindowsManager) StartService(serviceName string) error {
	w.mu.RLock()
	m := w.mgr
	w.mu.RUnlock()
	if m == nil {
		return fmt.Errorf("not connected")
	}
	s, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer s.Close()
	return s.Start()
}

// StopService implements [domain.ServiceManager].
func (w *WindowsManager) StopService(serviceName string) error {
	w.mu.RLock()
	m := w.mgr
	w.mu.RUnlock()
	if m == nil {
		return fmt.Errorf("not connected")
	}
	s, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer s.Close()
	_, err = s.Control(windows.SERVICE_CONTROL_STOP)
	return err
}

func NewManager() domain.ServiceManager {
	return &WindowsManager{
		processCache: make(map[string]*processHandle),
	}
}
