package rds

import (
	"context"
	"database/sql/driver"
	"github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"
	"time"
)

// NewConnector from the provided configuration fields
func NewConnector(d driver.Driver, api rdsdataserviceiface.RDSDataServiceAPI, resourceARN string, secretARN string, database string) *Connector {
	return &Connector{
		driver:      d,
		rds:         api,
		resourceARN: resourceARN,
		secretARN:   secretARN,
		database:    database,
	}
}

// Connector spits out connections to our database.
type Connector struct {
	driver               driver.Driver
	rds                  rdsdataserviceiface.RDSDataServiceAPI
	resourceARN          string
	secretARN            string
	database             string
	lastSuccessfulWakeup time.Time
}

// Connect returns a connection to the database.
func (r *Connector) Connect(ctx context.Context) (*Connection, error) {
	if r.lastSuccessfulWakeup.Add(time.Minute * 5).Before(time.Now()) {
		err := Wakeup(r.rds, r.resourceARN, r.secretARN, r.database)
		if err != nil {
			return nil, err
		}
		r.lastSuccessfulWakeup = time.Now()
	}

	return NewConnection(ctx, r.rds, r.resourceARN, r.secretARN, r.database), nil
}

// Driver returns the underlying Driver of the Connector, mainly to maintain compatibility with the Driver method on sql.DB.
func (r *Connector) Driver() driver.Driver {
	return r.driver
}
