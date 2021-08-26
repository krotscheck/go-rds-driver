package rds

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ordinalRegex = regexp.MustCompile("\\?{1}")

func MigrateQuery(query string, args []driver.NamedValue) (*rdsdataservice.ExecuteStatementInput, error) {
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
		return &rdsdataservice.ExecuteStatementInput{
			Parameters: params,
			Sql:        aws.String(query),
		}, err
	} else {
		params, err := ConvertNamedValues(args)
		return &rdsdataservice.ExecuteStatementInput{
			Parameters: params,
			Sql:        aws.String(query),
		}, err
	}
}

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
	case int16:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{LongValue: aws.Int64(int64(t))},
		}
	case int64:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{LongValue: aws.Int64(t)},
		}
	case nil:
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{IsNull: aws.Bool(true)},
		}
	case time.Time:
		value = &rdsdataservice.SqlParameter{
			Name:     &name,
			TypeHint: aws.String("TIMESTAMP"),
			Value: &rdsdataservice.Field{
				StringValue: aws.String(t.In(time.UTC).Format("2006-01-02 15:04:05.999")),
			},
		}
	default:
		err = fmt.Errorf("%s is unsupported type: %#v", name, arg.Value)
		return
	}
	return
}

var SupportedIsolationLevels = map[driver.IsolationLevel]bool{}

// Wakeup the cluster if it's dormant
func Wakeup(r rdsdataserviceiface.RDSDataServiceAPI, resourceARN string, secretARN string, database string) (dialect Dialect, err error) {
	request := &rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(resourceARN),
		Database:    aws.String(database),
		SecretArn:   aws.String(secretARN),
		Sql:         aws.String("/* wakeup */ SELECT VERSION()"), // This works for all databases, I think.
		Parameters:  []*rdsdataservice.SqlParameter{},
	}

	err = retry(10, time.Second, func() error {
		out, err := r.ExecuteStatement(request)
		if err != nil {
			return err
		}
		rows := NewRows(out)

		values := make([]driver.Value, 1)
		if err := rows.Next(values); err != nil {
			return err
		}
		if len(values) < 1 {
			return fmt.Errorf("invalid response to version request")
		}
		version, ok := values[0].(string)
		if !ok {
			return fmt.Errorf("invalid response to version request")
		}
		if strings.Contains(strings.ToLower(version), "postgres") {
			dialect = DialectPostgres
		} else {
			dialect = DialectMySQL
		}
		return err
	})
	return
}

func retry(attempts int, sleep time.Duration, callback func() error) (err error) {
	for i := 0; ; i++ {
		err = callback()
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(sleep)
		log.Println("retrying after error:", err)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

func init() {
	// List of supported isolation levels for both Postgres and mysql
	SupportedIsolationLevels[driver.IsolationLevel(sql.LevelDefault)] = true
	SupportedIsolationLevels[driver.IsolationLevel(sql.LevelRepeatableRead)] = true
	SupportedIsolationLevels[driver.IsolationLevel(sql.LevelReadCommitted)] = true
	SupportedIsolationLevels[driver.IsolationLevel(sql.LevelReadUncommitted)] = true
	SupportedIsolationLevels[driver.IsolationLevel(sql.LevelSerializable)] = true
}
