package rds

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata/types"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var postgresRegex = regexp.MustCompile("\\$([0-9]+)")

// NewPostgres dialect from our configuration
func NewPostgres(config *Config) Dialect {
	return &DialectPostgres{parseTime: config.ParseTime}
}

// DialectPostgres is for postgres 10.14 as supported by aurora serverless
type DialectPostgres struct {
	parseTime bool
}

// MigrateQuery from Postgres to RDS.
func (d *DialectPostgres) MigrateQuery(query string, args []driver.NamedValue) (*rdsdata.ExecuteStatementInput, error) {
	// Make sure we're not mixing and matching.
	ordinal := false
	named := false
	for _, arg := range args {
		if arg.Name != "" {
			named = true
		}
		if arg.Ordinal > 0 {
			ordinal = true
		}
		if named && ordinal {
			return nil, ErrNoMixedParams
		}
	}

	// If we're ordinal, convert to named
	if ordinal {
		namedArgs := make([]driver.NamedValue, len(args))
		for i, v := range args {
			namedArgs[i] = driver.NamedValue{
				Name:  strconv.Itoa(v.Ordinal),
				Value: v.Value,
			}
		}
		args = namedArgs

		query = postgresRegex.ReplaceAllStringFunc(query, func(s string) string {
			return strings.Replace(s, "$", ":", 1)
		})

		params, err := ConvertNamedValues(namedArgs)
		return &rdsdata.ExecuteStatementInput{
			Parameters: params,
			Sql:        aws.String(query),
		}, err
	}
	params, err := ConvertNamedValues(args)
	return &rdsdata.ExecuteStatementInput{
		Parameters: params,
		Sql:        aws.String(query),
	}, err
}

// GetFieldConverter knows how to parse response data.
func (d *DialectPostgres) GetFieldConverter(columnType string) FieldConverter {
	switch strings.ToLower(columnType) {
	case "serial":
		return func(field types.Field) (interface{}, error) {
			return field.(*types.FieldMemberLongValue).Value, nil
		}
	case "bool":
		return func(field types.Field) (interface{}, error) {
			return field.(*types.FieldMemberBooleanValue).Value, nil
		}
	case "bpchar":
		fallthrough
	case "varchar":
		fallthrough
	case "text":
		return func(field types.Field) (interface{}, error) {
			return field.(*types.FieldMemberStringValue).Value, nil
		}
	case "int2":
		fallthrough
	case "int4":
		fallthrough
	case "int8":
		return func(field types.Field) (interface{}, error) {
			return field.(*types.FieldMemberLongValue).Value, nil
		}
	case "numeric":
		return func(field types.Field) (interface{}, error) {
			return strconv.ParseFloat(field.(*types.FieldMemberStringValue).Value, 64)
		}
	case "float4":
		return func(field types.Field) (interface{}, error) {
			return field.(*types.FieldMemberDoubleValue).Value, nil
		}
	case "date":
		return func(field types.Field) (interface{}, error) {
			t, err := time.Parse("2006-01-02", field.(*types.FieldMemberStringValue).Value)
			if err != nil {
				return nil, err
			}
			if d.parseTime {
				return t, nil
			}
			return t.Format(time.RFC3339), nil
		}
	case "time":
		return func(field types.Field) (interface{}, error) {
			time_str_val := field.(*types.FieldMemberStringValue).Value
			if d.parseTime {
				return time.Parse("15:04:05", time_str_val)
			}
			return time_str_val, nil
		}
	case "timestamp":
		return func(field types.Field) (interface{}, error) {
			t, err := time.Parse("2006-01-02 15:04:05", field.(*types.FieldMemberStringValue).Value)
			if err != nil {
				return nil, err
			}
			if d.parseTime {
				return t, nil
			}
			return t.Format(time.RFC3339), nil
		}
	}

	// For unknown field types (like custom ENUM types) we fall back to the
	// value type of the field.
	return func(field types.Field) (value interface{}, err error) {
		switch v := field.(type) {
		case *types.FieldMemberArrayValue:
			value = v.Value;
		case *types.FieldMemberBlobValue:
			value = v.Value
		case *types.FieldMemberBooleanValue:
			value = v.Value
		case *types.FieldMemberDoubleValue:
			value = v.Value
		case *types.FieldMemberLongValue:
			value = v.Value
		case *types.FieldMemberStringValue:
			value = v.Value
		case *types.FieldMemberIsNull:
			return nil, nil
		default:
			return nil, fmt.Errorf("Unrecognized RDS field type: %#v", field)
		}

		return
	}
}

// IsIsolationLevelSupported for postgres?
func (d *DialectPostgres) IsIsolationLevelSupported(level driver.IsolationLevel) bool {
	// SupportedIsolationLevels for the dialect
	var SupportedIsolationLevels = map[driver.IsolationLevel]bool{
		driver.IsolationLevel(sql.LevelDefault):         true,
		driver.IsolationLevel(sql.LevelRepeatableRead):  true,
		driver.IsolationLevel(sql.LevelReadCommitted):   true,
		driver.IsolationLevel(sql.LevelReadUncommitted): true,
		driver.IsolationLevel(sql.LevelSerializable):    true,
	}
	_, ok := SupportedIsolationLevels[level]
	return ok
}
