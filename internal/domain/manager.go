package domain

import "zenlight-support/pkg/sql"

type ResourceManager interface {
	Connect() error
	Disconnect() error

	GetResourceState(resourceName string) (Status, error)
	GetServiceMetrics(resourceName string) (*ResourceMetrics, error)
	GetDirectoryMetrics(path string) (*ResourceMetrics, error)

	StartService(serviceName string) error
	StopService(serviceName string) error

	ExecuteSQLScript(server, database, script string) (*sql.Result, error)
}
