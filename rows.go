package rds

import (
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"io"
)

// NewRows instance for the provided statement output
func NewRows(dialect Dialect, out *rdsdataservice.ExecuteStatementOutput) driver.Rows {
	converters := make([]FieldConverter, len(out.ColumnMetadata))
	names := make([]string, len(out.ColumnMetadata))
	for i, col := range out.ColumnMetadata {
		converters[i] = dialect.GetFieldConverter(*col.TypeName)
		names[i] = *col.Name
	}

	return &Rows{
		columnNames:    names,
		converters:     converters,
		out:            out,
		recordPosition: 0,
	}
}

// Rows implementation for the RDS Driver
type Rows struct {
	out            *rdsdataservice.ExecuteStatementOutput
	columnNames    []string
	converters     []FieldConverter
	recordPosition int
}

// Columns returns the column names in order
func (r *Rows) Columns() []string {
	return r.columnNames
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
		if aws.BoolValue(field.IsNull) {
			dest[i] = nil
			continue
		}
		converter := r.converters[i]
		coerced, err := converter(field)
		if err != nil {
			return fmt.Errorf("convertValue(col=%d): %v", i, err)
		}
		dest[i] = coerced
	}

	return nil
}
