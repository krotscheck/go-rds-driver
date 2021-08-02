package rds

import (
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

// NewResult for the executed statement
func NewResult(out *rdsdataservice.ExecuteStatementOutput) driver.Result {
	return &Result{out: out}
}

// Result from a query
type Result struct {
	out *rdsdataservice.ExecuteStatementOutput
}

// LastInsertId from the executed statement.
func (r *Result) LastInsertId() (int64, error) {
	if l := len(r.out.GeneratedFields); l == 0 {
		return 0, fmt.Errorf("no generated fields in result")
	} else if l != 1 {
		return 0, fmt.Errorf("%d generated fields in result: %v", l, r.out.GeneratedFields)
	}
	f := r.out.GeneratedFields[0]
	if f.LongValue != nil {
		return aws.Int64Value(f.LongValue), nil
	}
	return 0, fmt.Errorf("unhandled generated field type: %v", f)
}

// RowsAffected count
func (r *Result) RowsAffected() (int64, error) {
	return aws.Int64Value(r.out.NumberOfRecordsUpdated), nil
}
