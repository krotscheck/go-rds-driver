package rds

import (
	"context"
	"database/sql/driver"
)

// NewStatement for the provided connection
func NewStatement(_ context.Context, connection *Connection, sql string) driver.Stmt {
	return &Statement{
		conn:  connection,
		query: sql,
	}
}

// Statement encapsulates a single RDS query statement
type Statement struct {
	conn  *Connection
	query string
}

// Close closes the statement.
func (s *Statement) Close() error {
	// The HTTP protocol is stateless, so nothing to do here.
	return nil
}

// NumInput returns the number of placeholder parameters.
func (s *Statement) NumInput() int {
	return -1
}

// Exec executes a query that doesn't return rows, such as an INSERT or UPDATE.
func (s *Statement) Exec(_ []driver.Value) (driver.Result, error) {
	return nil, ErrNoPositional
}

// ExecContext executes a query that doesn't return rows, such as an INSERT or UPDATE.
func (s *Statement) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	return s.conn.ExecContext(ctx, s.query, args)
}

// Query executes a query that may return rows, such as a SELECT.
func (s *Statement) Query(args []driver.Value) (driver.Rows, error) {
	return nil, ErrNoPositional
}

// QueryContext executes a query that may return rows, such as a SELECT.
func (s *Statement) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	return s.conn.QueryContext(ctx, s.query, args)
}
