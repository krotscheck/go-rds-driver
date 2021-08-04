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
		Done:          false,
		TransactionID: transactionID,
		conn:          conn,
	}
}

// Tx is a transaction
type Tx struct {
	Done          bool
	TransactionID *string
	conn          *Connection
}

// Commit the transaction
func (r *Tx) Commit() error {
	if r.Done {
		return sql.ErrTxDone
	}

	_, err := r.conn.rds.CommitTransaction(&rdsdataservice.CommitTransactionInput{
		ResourceArn:   aws.String(r.conn.resourceARN),
		SecretArn:     aws.String(r.conn.secretARN),
		TransactionId: r.TransactionID,
	})
	if err != nil {
		return err
	}
	if r.conn.tx == r {
		r.conn.tx = nil
	}
	r.Done = true
	return nil
}

// Rollback the transaction
func (r *Tx) Rollback() error {
	if r.Done {
		return sql.ErrTxDone
	}
	_, err := r.conn.rds.RollbackTransaction(&rdsdataservice.RollbackTransactionInput{
		ResourceArn:   aws.String(r.conn.resourceARN),
		SecretArn:     aws.String(r.conn.secretARN),
		TransactionId: r.TransactionID,
	})
	if err != nil {
		return err
	}
	if r.conn.tx == r {
		r.conn.tx = nil
	}
	r.Done = true
	return nil
}
