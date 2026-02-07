package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

type Result struct {
	RowsAffected  int64           `json:"rowsAffected"`
	LastInsertID  int64           `json:"lastInsertId"`
	Columns       []string        `json:"columns,omitempty"`
	Data          [][]interface{} `json:"data,omitempty"`
	ExecutionTime string          `json:"executionTime,omitempty"`
}

type Executor struct {
	server   string
	database string
}

func NewExecutor(server, database string) *Executor {
	return &Executor{
		server:   server,
		database: database,
	}
}

func (e *Executor) Execute(ctx context.Context, script string) (*Result, error) {
	connStr := fmt.Sprintf("server=%s;database=%s;trusted_connection=yes;encrypt=disable", e.server, e.database)

	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	start := time.Now()

	isQuery := strings.HasPrefix(strings.TrimSpace(strings.ToUpper(script)), "SELECT")

	if isQuery {
		return e.executeQuery(ctx, db, script, start)
	}

	return e.executeNonQuery(ctx, db, script, start)
}

func (e *Executor) executeQuery(ctx context.Context, db *sql.DB, script string, start time.Time) (*Result, error) {
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	const maxRows = 1000

	var data [][]interface{}
	for rows.Next() {
		if len(data) >= maxRows {
			return nil, fmt.Errorf("query result exceeds maximum row limit of %d. Please narrow your query with a WHERE or TOP clause", maxRows)
		}

		row := make([]interface{}, len(columns))
		rowPtrs := make([]interface{}, len(columns))
		for i := range row {
			rowPtrs[i] = &row[i]
		}

		if err := rows.Scan(rowPtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		data = append(data, row)
	}

	execTime := time.Since(start).String()

	return &Result{
		Columns:       columns,
		Data:          data,
		ExecutionTime: execTime,
	}, nil
}

func (e *Executor) executeNonQuery(ctx context.Context, db *sql.DB, script string, start time.Time) (*Result, error) {
	res, err := db.ExecContext(ctx, script)
	if err != nil {
		return nil, fmt.Errorf("failed to execute non-query: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		lastInsertID = 0
	}

	execTime := time.Since(start).String()

	return &Result{
		RowsAffected:  rowsAffected,
		LastInsertID:  lastInsertID,
		ExecutionTime: execTime,
	}, nil
}
