package rds

import (
	"database/sql"
	"database/sql/driver"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

// NewTx creates a new transaction
func NewTx(transactionID *string, conn *Connection) driver.Tx {
	return &Tx{
		done:          false,
		transactionID: transactionID,
		conn:          conn,
	}
}

// Tx is a transaction
type Tx struct {
	done          bool
	transactionID *string
	conn          *Connection
}

// Commit the transaction
func (r *Tx) Commit() error {
	if r.done {
		return sql.ErrTxDone
	}

	_, err := r.conn.rds.CommitTransaction(&rdsdataservice.CommitTransactionInput{
		ResourceArn:   aws.String(r.conn.resourceARN),
		SecretArn:     aws.String(r.conn.secretARN),
		TransactionId: r.transactionID,
	})
	if err != nil {
		return err
	}
	if r.conn.tx == r {
		r.conn.tx = nil
	}
	r.done = true
	return nil
}

// Rollback the transaction
func (r *Tx) Rollback() error {
	if r.done {
		return sql.ErrTxDone
	}
	_, err := r.conn.rds.RollbackTransaction(&rdsdataservice.RollbackTransactionInput{
		ResourceArn:   aws.String(r.conn.resourceARN),
		SecretArn:     aws.String(r.conn.secretARN),
		TransactionId: r.transactionID,
	})
	if err != nil {
		return err
	}
	if r.conn.tx == r {
		r.conn.tx = nil
	}
	r.done = true
	return nil
}
