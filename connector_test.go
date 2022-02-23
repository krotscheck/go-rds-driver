package rds_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/krotscheck/go-rds-driver"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Connector(t *testing.T) {
	conf := rds.NewConfig("resourceARN", "secretARN", "database", "region")

	Convey("Connector", t, func() {
		contrl := gomock.NewController(t)
		d := rds.NewDriver()
		mockRDS := NewMockAWSClientInterface(contrl)
		connector := rds.NewConnector(d, mockRDS, conf)
		ctx := context.Background()

		Convey("Connect", func() {
			ExpectWakeup(mockRDS, conf)

			connection, err := connector.Connect(ctx)
			So(err, ShouldBeNil)
			So(connection, ShouldNotBeNil)
		})

		Convey("Driver", func() {
			So(connector.Driver(), ShouldEqual, d)
		})
	})
}
