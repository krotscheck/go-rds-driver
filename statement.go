package rds

import (
	"context"
	"database/sql/driver"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

// NewStatement for the provided connection
func NewStatement(_ context.Context, connection *Connection, sql string) *Statement {
	// TODO: Determine if this is an ordinal or named query.
	// TODO: Convert an ordinal query into a named query
	// TODO: Implement the ordinal methods below into name conversion queries.
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
	out, err := s.executeStatement(ctx, args)
	if err != nil {
		return nil, err
	}
	return NewResult(out), nil
}

// Query executes a query that may return rows, such as a SELECT.
func (s *Statement) Query(args []driver.Value) (driver.Rows, error) {
	return nil, ErrNoPositional
}

// QueryContext executes a query that may return rows, such as a SELECT.
func (s *Statement) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	out, err := s.executeStatement(ctx, args)
	if err != nil {
		return nil, err
	}
	return NewRows(out), nil
}

func (s *Statement) executeStatement(ctx context.Context, args []driver.NamedValue) (*rdsdataservice.ExecuteStatementOutput, error) {
	var txID *string
	if s.conn.tx != nil {
		txID = s.conn.tx.TransactionID
	}
	params, err := ConvertNamedValues(args)
	if err != nil {
		return nil, err
	}

	req := &rdsdataservice.ExecuteStatementInput{
		Database:              aws.String(s.conn.database),
		ResourceArn:           aws.String(s.conn.resourceARN),
		SecretArn:             aws.String(s.conn.secretARN),
		TransactionId:         txID,
		IncludeResultMetadata: aws.Bool(true),
		Parameters:            params,
		Sql:                   aws.String(s.query),
	}

	return s.conn.rds.ExecuteStatementWithContext(ctx, req)
}
