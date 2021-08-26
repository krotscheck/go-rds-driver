package rds

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"
	"log"
	"strings"
	"time"
)

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

		if len(out.Records) < 1 {
			return fmt.Errorf("invalid response to version request")
		}
		row := out.Records[0]
		if len(row) < 1 {
			return fmt.Errorf("invalid response to version request")
		}
		field := row[0]
		version := aws.StringValue(field.StringValue)
		if strings.Contains(strings.ToLower(version), "postgres") {
			dialect = &DialectPostgres{}
		} else {
			dialect = &DialectMySQL{}
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
