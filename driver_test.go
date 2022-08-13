package rds_test

import (
	"database/sql"
	"testing"

	"github.com/jonbretman/go-rds-driver"
	. "github.com/smartystreets/goconvey/convey"
)

//go:generate mockgen -package rds_test -destination client_mocks_test.go . AWSClientInterface

func Test_Driver(t *testing.T) {
	Convey("Driver", t, func() {
		So(sql.Drivers(), ShouldContain, rds.DRIVERNAME)

		testDSN := TestMysqlConfig.ToDSN()
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
