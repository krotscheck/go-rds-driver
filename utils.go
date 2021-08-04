package rds

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"
	"log"
	"time"
)

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
	if name == "" {
		err = ErrNoPositional
		return
	}

	if arg.Value == nil {
		value = &rdsdataservice.SqlParameter{
			Name:  &name,
			Value: &rdsdataservice.Field{IsNull: aws.Bool(true)},
		}
		return
	}
	var f *rdsdataservice.Field
	switch t := arg.Value.(type) {
	case string:
		f = &rdsdataservice.Field{StringValue: aws.String(t)}
	case []byte:
		f = &rdsdataservice.Field{BlobValue: t}
	case bool:
		f = &rdsdataservice.Field{BooleanValue: &t}
	case float64:
		f = &rdsdataservice.Field{DoubleValue: &t}
	case int64:
		f = &rdsdataservice.Field{LongValue: &t}
	default:
		err = fmt.Errorf("%s is unsupported type: %#v", name, value)
		return
	}
	value = &rdsdataservice.SqlParameter{
		Name:  &name,
		Value: f,
	}
	return
}

var SupportedIsolationLevels = map[driver.IsolationLevel]bool{}

// Wakeup the cluster if it's dormant
func Wakeup(r rdsdataserviceiface.RDSDataServiceAPI, resourceARN string, secretARN string, database string) error {
	request := &rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(resourceARN),
		Database:    aws.String(database),
		SecretArn:   aws.String(secretARN),
		Sql:         aws.String("/* wakeup */ SELECT 1"), // This works for all databases, I think.
		Parameters:  []*rdsdataservice.SqlParameter{},
	}

	return retry(10, time.Second, func() (err error) {
		_, err = r.ExecuteStatement(request)
		return
	})
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
