package rds

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"
	"strings"
)

// NewConnection that can make transaction and statement requests against RDS
func NewConnection(ctx context.Context, rds rdsdataserviceiface.RDSDataServiceAPI, resourceARN string, secretARN string, database string, dialect Dialect) *Connection {
	return &Connection{
		ctx:         ctx,
		rds:         rds,
		resourceARN: resourceARN,
		secretARN:   secretARN,
		database:    database,
		closed:      false,
		dialect:     dialect,
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
	closed      bool
	dialect     Dialect
}

// Ping the database
func (r *Connection) Ping(ctx context.Context) (err error) {
	_, err = r.rds.ExecuteStatementWithContext(ctx, &rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(r.resourceARN),
		Database:    aws.String(r.database),
		SecretArn:   aws.String(r.secretARN),
		Sql:         aws.String("/* ping */ SELECT 1"), // This works for all databases, I think.
		Parameters:  []*rdsdataservice.SqlParameter{},
	})
	return
}

// Prepare returns a prepared statement, bound to this connection.
func (r *Connection) Prepare(query string) (driver.Stmt, error) {
	return r.PrepareContext(context.Background(), query)
}

// PrepareContext returns a prepared statement, bound to this connection.
func (r *Connection) PrepareContext(ctx context.Context, query string) (*Statement, error) {
	return NewStatement(ctx, r, query), nil
}

// Close the connection
func (r *Connection) Close() error {
	if r.closed {
		return ErrClosed
	}
	if r.tx != nil {
		if err := r.tx.Rollback(); err != nil {
			return err
		}
	}
	r.rds = nil
	r.closed = true
	return nil
}

// Begin starts and returns a new transaction.
func (r *Connection) Begin() (driver.Tx, error) {
	return r.BeginTx(r.ctx, driver.TxOptions{
		Isolation: driver.IsolationLevel(sql.LevelDefault),
		ReadOnly:  false,
	})
}

// BeginTx starts and returns a new transaction.
func (r *Connection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	// Assume that the underlying database supports all isolation levels.
	if _, ok := SupportedIsolationLevels[opts.Isolation]; !ok {
		return nil, fmt.Errorf("isolation level %d not supported", opts.Isolation)
	}
	var clause []string
	if sql.IsolationLevel(opts.Isolation) != sql.LevelDefault {
		clause = append(clause, fmt.Sprintf("ISOLATION LEVEL %s", sql.IsolationLevel(opts.Isolation).String()))
	}
	if opts.ReadOnly {
		clause = append(clause, "READ ONLY")
	} else {
		clause = append(clause, "READ WRITE")
	}
	query := fmt.Sprintf("SET TRANSACTION %s", strings.Join(clause, ", "))

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
		Done:          false,
		TransactionID: output.TransactionId,
		conn:          r,
	}

	if _, err := r.ExecContext(ctx, query, nil); err != nil {
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
		if err := r.tx.Rollback(); err != nil {
			return err
		}
		return driver.ErrBadConn
	}
	return nil
}

// IsValid is called prior to placing the connection into the connection pool. The connection will be discarded if false is returned.
func (r *Connection) IsValid() bool {
	if r.closed {
		return false
	}
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
	stmt, err := r.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return stmt.QueryContext(ctx, args)
}

// ExecContext executes a query that would normally not return a result.
func (r *Connection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	stmt, err := r.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return stmt.ExecContext(ctx, args)
}
