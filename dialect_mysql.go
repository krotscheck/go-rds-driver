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
	"time"
)

var ordinalRegex = regexp.MustCompile("\\?{1}")

// NewMySQL dialect from our configuration
func NewMySQL(config *Config) Dialect {
	return &DialectMySQL{parseTime: config.ParseTime}
}

// DialectMySQL for version 5.7
type DialectMySQL struct {
	parseTime bool
}

// MigrateQuery converts a mysql query into an RDS stateement.
func (d *DialectMySQL) MigrateQuery(query string, args []driver.NamedValue) (*rdsdata.ExecuteStatementInput, error) {
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

		idx := 0
		query = ordinalRegex.ReplaceAllStringFunc(query, func(s string) string {
			idx = idx + 1 // ordinal regex are one-indexed
			return fmt.Sprintf(":%d", idx)
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

// GetFieldConverter knows how to parse column results.
func (d *DialectMySQL) GetFieldConverter(columnType string) FieldConverter {
	switch columnType {
	case "TINYINT UNSIGNED":
		fallthrough
	case "SMALLINT UNSIGNED":
		fallthrough
	case "MEDIUMINT UNSIGNED":
		fallthrough
	case "INT UNSIGNED":
		fallthrough
	case "BIGINT UNSIGNED":
		return func(field types.Field) (interface{}, error) {
			return uint64(field.(*types.FieldMemberLongValue).Value), nil
		}
	case "TINYINT":
		fallthrough
	case "SMALLINT":
		fallthrough
	case "MEDIUMINT":
		fallthrough
	case "INT":
		fallthrough
	case "BIGINT":
		return func(field types.Field) (interface{}, error) {
			return field.(*types.FieldMemberLongValue).Value, nil
		}
	case "DECIMAL":
		return func(field types.Field) (interface{}, error) {
			return strconv.ParseFloat(field.(*types.FieldMemberStringValue).Value, 64)
		}
	case "FLOAT":
		fallthrough
	case "DOUBLE":
		return func(field types.Field) (interface{}, error) {
			return aws.Float64(field.(*types.FieldMemberDoubleValue).Value), nil
		}
	case "BIT":
		// Bit values appear to be returned as boolean values
		return func(field types.Field) (interface{}, error) {
			if field.(*types.FieldMemberBooleanValue).Value {
				return 1, nil
			}
			return 0, nil
		}
	case "TINYTEXT":
		fallthrough
	case "TEXT":
		fallthrough
	case "MEDIUMTEXT":
		fallthrough
	case "LONGTEXT":
		fallthrough
	case "CHAR":
		fallthrough
	case "VARCHAR":
		return func(field types.Field) (interface{}, error) {
			return field.(*types.FieldMemberStringValue).Value, nil
		}
	case "DATE":
		return func(field types.Field) (interface{}, error) {
			date_str_val := field.(*types.FieldMemberStringValue).Value
			if d.parseTime {
				return time.Parse("2006-01-02", date_str_val)
			}
			return date_str_val, nil
		}
	case "TIME":
		return func(field types.Field) (interface{}, error) {
			time_str_val := field.(*types.FieldMemberStringValue).Value
			if d.parseTime {
				return time.Parse("15:04:05", time_str_val)
			}
			return time_str_val, nil
		}
	case "DATETIME":
		return func(field types.Field) (interface{}, error) {
			dt_str_val := field.(*types.FieldMemberStringValue).Value
			if d.parseTime {
				return time.Parse("2006-01-02 15:04:05", dt_str_val)
			}
			return dt_str_val, nil
		}
	case "TIMESTAMP":
		return func(field types.Field) (interface{}, error) {
			ts_str_val := field.(*types.FieldMemberStringValue).Value
			if d.parseTime {
				return time.Parse("2006-01-02 15:04:05", ts_str_val)
			}
			return ts_str_val, nil
		}
	case "YEAR":
		// RDS sends a full date string. MySQL only returns the year.
		return func(field types.Field) (interface{}, error) {
			t, err := time.Parse("2006-01-02", field.(*types.FieldMemberStringValue).Value)
			if err != nil {
				return nil, err
			}
			if d.parseTime {
				return t, nil
			}
			return strconv.Itoa(t.Year()), nil
		}
	case "BINARY":
		fallthrough
	case "VARBINARY":
		fallthrough
	case "TINYBLOB":
		fallthrough
	case "BLOB":
		fallthrough
	case "MEDIUMBLOB":
		fallthrough
	case "LONGBLOB":
		return func(field types.Field) (interface{}, error) {
			return field.(*types.FieldMemberBlobValue).Value, nil
		}
	}

	return ConversionFallback()
}

// IsIsolationLevelSupported for mysql?
func (d *DialectMySQL) IsIsolationLevelSupported(level driver.IsolationLevel) bool {
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
