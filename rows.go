package rds

import (
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"io"
)

// NewRows instance for the provided statement output
func NewRows(out *rdsdataservice.ExecuteStatementOutput) driver.Rows {
	return &Rows{
		out:            out,
		recordPosition: 0,
	}
}

// Rows implementation for the RDS Driver
type Rows struct {
	out            *rdsdataservice.ExecuteStatementOutput
	recordPosition int
}

// Columns returns the column names in order
func (r *Rows) Columns() []string {
	// First see if there's column metadata that we can check for.
	cols := make([]string, len(r.out.ColumnMetadata))
	for i, c := range r.out.ColumnMetadata {
		cols[i] = aws.StringValue(c.Name)
	}
	return cols
}

// Close the result set
func (r *Rows) Close() error {
	// The API is stateless, so there's no connection to close
	return nil
}

// Next row in the result set
func (r *Rows) Next(dest []driver.Value) error {
	if r.recordPosition == len(r.out.Records) {
		return io.EOF
	}
	row := r.out.Records[r.recordPosition]
	r.recordPosition++

	for i, field := range row {
		coerced, err := r.convertField(field)
		if err != nil {
			return fmt.Errorf("convertValue(col=%d): %v", i, err)
		}
		dest[i] = coerced
	}

	return nil
}

func (r *Rows) convertField(field *rdsdataservice.Field) (interface{}, error) {
	switch {
	case field.BlobValue != nil:
		return field.BlobValue, nil
	case field.BooleanValue != nil:
		return *field.BooleanValue, nil
	case field.DoubleValue != nil:
		return *field.DoubleValue, nil
	case field.IsNull != nil:
		return nil, nil
	case field.LongValue != nil:
		return *field.LongValue, nil
	case field.StringValue != nil:
		return *field.StringValue, nil
	default:
		return nil, fmt.Errorf("no part of Field non-nil")
	}
}
