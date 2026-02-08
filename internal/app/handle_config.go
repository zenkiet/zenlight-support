package app

import (
	"zenlight-support/internal/domain"
)

func (a *App) GetSqlConfig() (*domain.SQLConfig, error) {
	return a.cfg.SQLConfig, nil
}
