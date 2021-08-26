package rds

import (
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

// FieldConverter is a function that converts the passed result row field into the expected type.
type FieldConverter func(field *rdsdataservice.Field) (interface{}, error)

// Dialect is an interface that encapsulates a particular languages' eccentricities
type Dialect interface {
	// MigrateQuery from the dialect to RDS
	MigrateQuery(string, []driver.NamedValue) (*rdsdataservice.ExecuteStatementInput, error)
	// GetFieldConverter for a given ColumnMetadata.TypeName field.
	GetFieldConverter(columnType string) FieldConverter
	// IsIsolationLevelSupported for this dialect?
	IsIsolationLevelSupported(level driver.IsolationLevel) bool
}

// ConvertNamedValues converts passed driver.NamedValue instances into RDS SQLParameters
func ConvertNamedValues(args []driver.NamedValue) ([]*rdsdataservice.SqlParameter, error) {
	var params = make([]*rdsdataservice.SqlParameter, len(args))
	for i, arg := range args {
		sqlParam, err := ConvertNamedValue(arg)
		if err != nil {
			return nil, err
		}
		params[i] = sqlParam
	}
	return params, nil
}

// ConvertNamedValue from a NamedValue to an SqlParameter
func ConvertNamedValue(arg driver.NamedValue) (value *rdsdataservice.SqlParameter, err error) {
	name := arg.Name
	if arg.Value == nil {
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{IsNull: aws.Bool(true)},
		}
		return
	}
	switch t := arg.Value.(type) {
	case string:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{StringValue: aws.String(t)},
		}
	case []byte:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{BlobValue: t},
		}
	case bool:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{BooleanValue: &t},
		}
	case float32:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{DoubleValue: aws.Float64(float64(t))},
		}
	case float64:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{DoubleValue: &t},
		}
	case int:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{LongValue: aws.Int64(int64(t))},
		}
	case int8:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{LongValue: aws.Int64(int64(t))},
		}
	case int16:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{LongValue: aws.Int64(int64(t))},
		}
	case int32:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{LongValue: aws.Int64(int64(t))},
		}
	case int64:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{LongValue: aws.Int64(t)},
		}
	case uint:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{LongValue: aws.Int64(int64(t))},
		}
	case uint8:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{BlobValue: []byte{t}},
		}
	case uint16:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{LongValue: aws.Int64(int64(t))},
		}
	case uint32:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{LongValue: aws.Int64(int64(t))},
		}
	case uint64:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{LongValue: aws.Int64(int64(t))},
		}
	case nil:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{IsNull: aws.Bool(true)},
		}
	default:
		err = fmt.Errorf("%s is unsupported type: %#v", name, arg.Value)
		return
	}
	return
}
