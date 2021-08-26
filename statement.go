package rds

import (
	"context"
	"database/sql/driver"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

// NewStatement for the provided connection
func NewStatement(_ context.Context, connection *Connection, sql string) *Statement {
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
	if s.conn == nil {
		return ErrClosed
	}
	s.conn = nil
	return nil
}

// NumInput returns the number of placeholder parameters.
func (s *Statement) NumInput() int {
	return -1
}

// Exec executes a query that doesn't return rows, such as an INSERT or UPDATE.
func (s *Statement) Exec(values []driver.Value) (driver.Result, error) {
	args := s.ConvertOrdinal(values)
	out, err := s.executeStatement(context.Background(), s.query, args)
	if err != nil {
		return nil, err
	}
	return NewResult(out), nil
}

// ExecContext executes a query that doesn't return rows, such as an INSERT or UPDATE.
func (s *Statement) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	out, err := s.executeStatement(ctx, s.query, args)
	if err != nil {
		return nil, err
	}
	return NewResult(out), nil
}

// Query executes a query that may return rows, such as a SELECT.
func (s *Statement) Query(values []driver.Value) (driver.Rows, error) {
	// We're trying to execute this as an ordinal query, so convert it.
	args := s.ConvertOrdinal(values)

	out, err := s.executeStatement(context.Background(), s.query, args)
	if err != nil {
		return nil, err
	}
	return NewRows(s.conn.dialect, out), nil
}

// QueryContext executes a query that may return rows, such as a SELECT.
func (s *Statement) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	out, err := s.executeStatement(ctx, s.query, args)
	if err != nil {
		return nil, err
	}
	return NewRows(s.conn.dialect, out), nil
}

// ConvertOrdinal converts a list of Values to Ordinal NamedValues
func (s *Statement) ConvertOrdinal(values []driver.Value) []driver.NamedValue {
	// Start with the MySQL separator as a default
	namedValues := make([]driver.NamedValue, len(values))
	for i, v := range values {
		namedValues[i] = driver.NamedValue{
			Name:    "",
			Ordinal: i + 1,
			Value:   v,
		}
	}
	return namedValues
}

func (s *Statement) executeStatement(ctx context.Context, query string, values []driver.NamedValue) (*rdsdataservice.ExecuteStatementOutput, error) {
	input, err := s.conn.dialect.MigrateQuery(query, values)
	if err != nil {
		return nil, err
	}

	if s.conn.tx != nil {
		input.TransactionId = s.conn.tx.TransactionID
	}

	input.IncludeResultMetadata = aws.Bool(true)
	input.ResourceArn = aws.String(s.conn.resourceARN)
	input.SecretArn = aws.String(s.conn.secretARN)
	input.Database = aws.String(s.conn.database)
	return s.conn.rds.ExecuteStatementWithContext(ctx, input)
}
