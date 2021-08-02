package rds

import (
	"context"
	"database/sql/driver"
	"github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface"
)

// NewConnector from the provided configuration fields
func NewConnector(d driver.Driver, api rdsdataserviceiface.RDSDataServiceAPI, resourceARN string, secretARN string, database string) *Connector {
	return &Connector{
		driver:      d,
		api:         api,
		resourceARN: resourceARN,
		secretARN:   secretARN,
		database:    database,
	}
}

// Connector spits out connections to our database.
type Connector struct {
	driver      driver.Driver
	api         rdsdataserviceiface.RDSDataServiceAPI
	resourceARN string
	secretARN   string
	database    string
}

// Connect returns a connection to the database.
func (r *Connector) Connect(ctx context.Context) (*Connection, error) {
	return NewConnection(ctx, r.api, r.resourceARN, r.secretARN, r.database), nil
}

// Driver returns the underlying Driver of the Connector, mainly to maintain compatibility with the Driver method on sql.DB.
func (r *Connector) Driver() driver.Driver {
	return r.driver
}
