package rds

import (
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata/types"
	"reflect"
	"time"
)

// FieldConverter is a function that converts the passed result row field into the expected type.
type FieldConverter func(field types.Field) (interface{}, error)

// Dialect is an interface that encapsulates a particular languages' eccentricities
type Dialect interface {
	// MigrateQuery from the dialect to RDS
	MigrateQuery(string, []driver.NamedValue) (*rdsdata.ExecuteStatementInput, error)
	// GetFieldConverter for a given ColumnMetadata.TypeName field.
	GetFieldConverter(columnType string) FieldConverter
	// IsIsolationLevelSupported for this dialect?
	IsIsolationLevelSupported(level driver.IsolationLevel) bool
}

// ConvertNamedValues converts passed driver.NamedValue instances into RDS SQLParameters
func ConvertNamedValues(args []driver.NamedValue) ([]types.SqlParameter, error) {
	var params = make([]types.SqlParameter, len(args))
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
func ConvertNamedValue(arg driver.NamedValue) (value types.SqlParameter, err error) {
	name := arg.Name

	if isNil(arg.Value) {
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberIsNull{Value: true},
		}
		return
	}
	switch t := arg.Value.(type) {
	case string:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberStringValue{Value: t},
		}
	case []byte:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberBlobValue{Value: t},
		}
	case bool:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberBooleanValue{Value: t},
		}
	case float32:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberDoubleValue{Value: float64(t)},
		}
	case float64:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberDoubleValue{Value: float64(t)},
		}
	case int:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberLongValue{Value: int64(t)},
		}
	case int8:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberLongValue{Value: int64(t)},
		}
	case int16:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberLongValue{Value: int64(t)},
		}
	case int32:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberLongValue{Value: int64(t)},
		}
	case int64:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberLongValue{Value: t},
		}
	case uint:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberLongValue{Value: int64(t)},
		}
	case uint8:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberBlobValue{Value: []byte{t}},
		}
	case uint16:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberLongValue{Value: int64(t)},
		}
	case uint32:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberLongValue{Value: int64(t)},
		}
	case uint64:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberLongValue{Value: int64(t)},
		}
	case time.Time:
		value = types.SqlParameter{
			Name: &name,
			Value: &types.FieldMemberStringValue{
				Value: t.Format("2006-01-02 15:04:05.999"),
			},
		}
	case nil:
		value = types.SqlParameter{
			Name:  &name,
			Value: &types.FieldMemberIsNull{Value: true},
		}
	default:
		err = fmt.Errorf("%s is unsupported type: %#v", name, arg.Value)
		return
	}
	return
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
