package rds

import (
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata/types"
)

// NewResult for the executed statement
func NewResult(out *rdsdata.ExecuteStatementOutput) driver.Result {
	return &Result{out: out}
}

// Result from a query
type Result struct {
	out *rdsdata.ExecuteStatementOutput
}

// LastInsertId from the executed statement.
func (r *Result) LastInsertId() (int64, error) {
	if l := len(r.out.GeneratedFields); l == 0 {
		return 0, fmt.Errorf("no generated fields in result")
	} else if l != 1 {
		return 0, fmt.Errorf("%d generated fields in result: %v", l, r.out.GeneratedFields)
	}

	field := r.out.GeneratedFields[0]

	switch fv := field.(type) {
	case *types.FieldMemberLongValue:
		return fv.Value, nil
	default:
		break
	}

	return 0, fmt.Errorf("unhandled generated field type: %v", field)
}

// RowsAffected count
func (r *Result) RowsAffected() (int64, error) {
	return r.out.NumberOfRecordsUpdated, nil
}
