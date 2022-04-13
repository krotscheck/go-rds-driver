package rds_test

import (
	"github.com/krotscheck/go-rds-driver"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Config(t *testing.T) {
	Convey("Config", t, func() {
		conf := rds.NewConfig("resourceARN", "secretARN", "database", "region")
		dsn := conf.ToDSN()
		conf1, err := rds.NewConfigFromDSN(dsn)
		So(err, ShouldBeNil)
		So(conf, ShouldResemble, conf1)
	})

	Convey("Parse", t, func() {
		dsn := "rds://?resource_arn=resourceARN&secret_arn=secretARN&database=database&aws_region=region"
		conf, err := rds.NewConfigFromDSN(dsn)
		So(err, ShouldBeNil)
		So(conf.ResourceArn, ShouldEqual, "resourceARN")
		So(conf.SecretArn, ShouldEqual, "secretARN")
		So(conf.Database, ShouldEqual, "database")
		So(conf.AWSRegion, ShouldEqual, "region")
	})

	Convey("Invalid Scheme", t, func() {
		dsn := "postgres://?resource_arn=resourceARN&secret_arn=secretARN&database=database&aws_region=region"
		_, err := rds.NewConfigFromDSN(dsn)
		So(err, ShouldEqual, rds.ErrInvalidDSNScheme)
	})
}
