package rds

import (
	"context"
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"
	"log"
	"strings"
	"time"
)

// NewConnector from the provided configuration fields
func NewConnector(d driver.Driver, api rdsdataserviceiface.RDSDataServiceAPI, conf *Config) *Connector {
	return &Connector{
		driver: d,
		rds:    api,
		conf:   conf,
	}
}

// Connector spits out connections to our database.
type Connector struct {
	driver               driver.Driver
	rds                  rdsdataserviceiface.RDSDataServiceAPI
	conf                 *Config
	lastSuccessfulWakeup time.Time
	dialect              Dialect
}

// Connect returns a connection to the database.
func (r *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	if r.lastSuccessfulWakeup.Add(time.Minute * 5).Before(time.Now()) {
		dialect, err := r.Wakeup()
		if err != nil {
			return nil, err
		}
		r.dialect = dialect
		r.lastSuccessfulWakeup = time.Now()
	}

	return NewConnection(ctx, r.rds, r.conf.ResourceArn, r.conf.SecretArn, r.conf.Database, r.dialect), nil
}

// Driver returns the underlying Driver of the Connector, mainly to maintain compatibility with the Driver method on sql.DB.
func (r *Connector) Driver() driver.Driver {
	return r.driver
}

// Wakeup the cluster if it's dormant
func (r *Connector) Wakeup() (dialect Dialect, err error) {
	request := &rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(r.conf.ResourceArn),
		Database:    aws.String(r.conf.Database),
		SecretArn:   aws.String(r.conf.SecretArn),
		Sql:         aws.String("/* wakeup */ SELECT VERSION()"), // This works for all databases, I think.
		Parameters:  []*rdsdataservice.SqlParameter{},
	}

	err = r.retry(10, time.Second, func() error {
		out, err := r.rds.ExecuteStatement(request)
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
			dialect = NewPostgres(r.conf)
		} else {
			dialect = NewMySQL(r.conf)
		}
		return err
	})
	return
}

func (r *Connector) retry(attempts int, sleep time.Duration, callback func() error) (err error) {
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
