package app

import "zenlight-support/pkg/sql"

func (a *App) ExecuteSQLScript(id string, script string) (*sql.Result, error) {
	return a.mgr.ExecuteSQLScript(a.cfg.SQLConfig.Server, a.cfg.SQLConfig.Database, script)
}
