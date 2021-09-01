package rds

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

// DRIVERNAME is used when configuring your dialector
const DRIVERNAME = "rds"

// NewDriver creates a new driver instance for RDS
func NewDriver() *Driver {
	return &Driver{}
}

// Driver implements the driver.Driver interface for RDS
type Driver struct {
}

// Open returns a new connection to the database.
func (r *Driver) Open(name string) (driver.Conn, error) {
	connector, err := r.OpenConnector(name)
	if err != nil {
		return nil, err
	}
	return connector.Connect(context.Background())
}

// OpenConnector must parse the name in the same format that Driver.Open parses the name parameter.
func (r *Driver) OpenConnector(dsn string) (*Connector, error) {
	conf, err := NewConfigFromDSN(dsn)
	if err != nil {
		return nil, err
	}

	awsConfig := aws.NewConfig().WithRegion(conf.AWSRegion)
	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	rdsAPI := rdsdataservice.New(awsSession)

	return NewConnector(r, rdsAPI, conf), nil
}

func init() {
	sql.Register(DRIVERNAME, NewDriver())
}
