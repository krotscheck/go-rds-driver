package rds_test

import (
	"database/sql"
	"github.com/krotscheck/rds"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

//go:generate mockgen -package rds_test -destination rdsdataservice_mocks_test.go github.com/aws/aws-sdk-go/service/rdsdataservice/rdsdataserviceiface RDSDataServiceAPI

func Test_Driver(t *testing.T) {
	Convey("Driver", t, func() {
		So(sql.Drivers(), ShouldContain, rds.DRIVERNAME)
	})
}
