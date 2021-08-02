package rds

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"
)

// NewConnection that can make transaction and statement requests against RDS
func NewConnection(ctx context.Context, rds rdsdataserviceiface.RDSDataServiceAPI, resourceARN string, secretARN string, database string) *Connection {
	return &Connection{
		ctx:         ctx,
		rds:         rds,
		resourceARN: resourceARN,
		secretARN:   secretARN,
		database:    database,
	}
}

// Connection to RDS's Aurora Serverless Data API
type Connection struct {
	ctx         context.Context
	rds         rdsdataserviceiface.RDSDataServiceAPI
	resourceARN string
	secretARN   string
	database    string
	tx          *Tx // The current transaction, if set
}

// Ping the database
func (r *Connection) Ping(ctx context.Context) (err error) {
	_, err = r.rds.ExecuteStatementWithContext(ctx, &rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(r.resourceARN),
		Database:    aws.String(r.database),
		SecretArn:   aws.String(r.secretARN),
		Sql:         aws.String("/* ping */ SELECT 1"), // This works for all databases, I think.
	})
	return
}

// Prepare returns a prepared statement, bound to this connection.
func (r *Connection) Prepare(query string) (driver.Stmt, error) {
	return r.PrepareContext(context.Background(), query)
}

// PrepareContext returns a prepared statement, bound to this connection.
func (r *Connection) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return NewStatement(ctx, r, query), nil
}

// Close the connection
func (r *Connection) Close() error {
	if r.tx != nil {
		if err := r.tx.Rollback(); err != nil {
			return err
		}
	}
	r.rds = nil
	return nil
}

// Begin starts and returns a new transaction.
func (r *Connection) Begin() (driver.Tx, error) {
	return r.BeginTx(r.ctx, driver.TxOptions{})
}

// BeginTx starts and returns a new transaction.
// If the ctx is canceled by the user the sql package will
// call Tx.Rollback before discarding and closing the connection.
//
// This must check opts.Isolation to determine if there is a set
// isolation level. If the driver does not support a non-default
// level and one is set or if there is a non-default isolation level
// that is not supported, an error must be returned.
//
// This must also check opts.ReadOnly to determine if the read-only
// value is true to either set the read-only transaction property if supported
// or return an error if it is not supported.
func (r *Connection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	// Assume that the underlying database supports all isolation levels.
	if _, ok := supportedIsolationLevels[opts.Isolation]; !ok {
		return nil, fmt.Errorf("isolation level %d not supported", opts.Isolation)
	}
	rw := "READ WRITE"
	if opts.ReadOnly {
		rw = "READ ONLY"
	}

	// Start the transaction
	output, err := r.rds.BeginTransactionWithContext(ctx, &rdsdataservice.BeginTransactionInput{
		Database:    aws.String(r.database),
		ResourceArn: aws.String(r.resourceARN),
		SecretArn:   aws.String(r.secretARN),
	})
	if err != nil {
		return nil, err
	}
	r.tx = &Tx{
		done:          false,
		transactionID: output.TransactionId,
		conn:          r,
	}

	// Set the isolation level and rw stats
	args := []driver.NamedValue{
		{Name: "isolation", Value: sql.IsolationLevel(opts.Isolation).String()},
		{Name: "readonly", Value: rw},
	}
	if _, err := r.ExecContext(ctx, "SET TRANSACTION ISOLATION LEVEL @isolation, @readonly", args); err != nil {
		defer func() {
			_ = r.tx.Rollback()
		}()
		return nil, err
	}
	return r.tx, nil
}

// ResetSession is called prior to executing a query on the connection
// if the connection has been used before. If the driver returns ErrBadConn
// the connection is discarded.
func (r *Connection) ResetSession(_ context.Context) error {
	if r.tx != nil {
		return driver.ErrBadConn
	}
	return nil
}

// IsValid is called prior to placing the connection into the
// connection pool. The connection will be discarded if false is returned.
func (r *Connection) IsValid() bool {
	if r.database == "" {
		return false
	}
	if r.resourceARN == "" {
		return false
	}
	if r.secretARN == "" {
		return false
	}
	return true
}

// QueryContext executes a statement that would return some kind of result.
func (r *Connection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	output, err := r.executeStatement(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return NewRows(output), nil
}

// ExecContext executes a query that would normally not return a result.
func (r *Connection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	output, err := r.executeStatement(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return NewResult(output), nil
}

func (r *Connection) executeStatement(ctx context.Context, query string, args []driver.NamedValue) (*rdsdataservice.ExecuteStatementOutput, error) {
	var txID *string
	if r.tx != nil {
		txID = r.tx.transactionID
	}
	params, err := convertNamedValues(args)
	if err != nil {
		return nil, err
	}

	req := &rdsdataservice.ExecuteStatementInput{
		Database:              aws.String(r.database),
		ResourceArn:           aws.String(r.resourceARN),
		SecretArn:             aws.String(r.secretARN),
		TransactionId:         txID,
		IncludeResultMetadata: aws.Bool(true),
		Parameters:            params,
		Sql:                   aws.String(query),
	}

	return r.rds.ExecuteStatementWithContext(ctx, req)
}
