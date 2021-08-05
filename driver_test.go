package rds_test

import (
	"database/sql"
	"github.com/krotscheck/go-rds-driver"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

//go:generate mockgen -package rds_test -destination rdsdataservice_mocks_test.go github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface RDSDataServiceAPI

func Test_Driver(t *testing.T) {
	Convey("Driver", t, func() {
		So(sql.Drivers(), ShouldContain, rds.DRIVERNAME)

		testDSN, err := TestConfig.ToDSN()
		So(err, ShouldBeNil)

		driver := rds.NewDriver()

		Convey("Open", func() {
			Convey("With Invalid DSN", func() {
				connection, err := driver.Open("invalid_dsn")
				So(err, ShouldNotBeNil)
				So(connection, ShouldBeNil)
			})

			Convey("With Valid DSN", func() {
				connection, err := driver.Open(testDSN)
				So(err, ShouldBeNil)
				So(connection, ShouldNotBeNil)
				err = connection.Close()
				So(err, ShouldBeNil)
			})
		})

		Convey("OpenConnector", func() {

			Convey("With Invalid DSN", func() {
				connector, err := driver.OpenConnector("invalid_dsn")
				So(err, ShouldNotBeNil)
				So(connector, ShouldBeNil)
			})

			Convey("With Valid DSN", func() {
				connector, err := driver.OpenConnector(testDSN)
				So(err, ShouldBeNil)
				So(connector, ShouldNotBeNil)
			})
		})
	})
}
