package rds_test

import (
	"context"
	"database/sql/driver"
	"github.com/golang/mock/gomock"
	"github.com/krotscheck/go-rds-driver"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Statement(t *testing.T) {
	ctx := context.Background()
	conf := rds.NewConfig("resourceARN", "secretARN", "database", "region")

	Convey("Statement", t, func() {
		contrl := gomock.NewController(t)
		d := rds.NewDriver()
		mockRDS := NewMockAWSClientInterface(contrl)
		ExpectWakeup(mockRDS, conf)

		connector := rds.NewConnector(d, mockRDS, conf)
		connection, err := connector.Connect(ctx)
		So(err, ShouldBeNil)

		dStmnt, err := connection.Prepare("select 1")
		So(err, ShouldBeNil)
		So(dStmnt, ShouldNotBeNil)

		stmnt, ok := dStmnt.(*rds.Statement)
		So(ok, ShouldBeTrue)

		Convey("ConvertOrdinal", func() {
			vars := []driver.Value{"one", "two", "three", "four"}

			Convey("Valid", func() {
				vars := stmnt.ConvertOrdinal(vars)
				So(err, ShouldBeNil)

				So(vars[0], ShouldResemble, driver.NamedValue{Ordinal: 1, Value: "one"})
				So(vars[1], ShouldResemble, driver.NamedValue{Ordinal: 2, Value: "two"})
				So(vars[2], ShouldResemble, driver.NamedValue{Ordinal: 3, Value: "three"})
				So(vars[3], ShouldResemble, driver.NamedValue{Ordinal: 4, Value: "four"})
			})
		})
	})
}
