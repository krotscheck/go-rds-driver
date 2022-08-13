package rds_test

import (
	"database/sql/driver"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	"github.com/aws/smithy-go/middleware"
	"github.com/jonbretman/go-rds-driver"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Rows(t *testing.T) {
	oneResult := &rdsdata.ExecuteStatementOutput{
		ColumnMetadata:         nil,
		FormattedRecords:       nil,
		GeneratedFields:        nil,
		NumberOfRecordsUpdated: 0,
		Records:                nil,
		ResultMetadata:         middleware.Metadata{},
	}
	twoResult := &rdsdata.ExecuteStatementOutput{
		ColumnMetadata:         nil,
		FormattedRecords:       nil,
		GeneratedFields:        nil,
		NumberOfRecordsUpdated: 0,
		Records:                nil,
		ResultMetadata:         middleware.Metadata{},
	}

	Convey("Row with single result set", t, func() {
		row := rds.NewRows(rds.NewMySQL(TestMysqlConfig), []*rdsdata.ExecuteStatementOutput{
			oneResult,
		})
		rowResult, ok := row.(driver.RowsNextResultSet)
		So(ok, ShouldBeTrue)
		So(rowResult.HasNextResultSet(), ShouldBeFalse)
		So(rowResult.NextResultSet(), ShouldEqual, io.EOF)
	})

	Convey("Row with multiple result sets", t, func() {
		row := rds.NewRows(rds.NewMySQL(TestMysqlConfig), []*rdsdata.ExecuteStatementOutput{
			oneResult,
			twoResult,
		})
		rowResult, ok := row.(driver.RowsNextResultSet)
		So(ok, ShouldBeTrue)
		So(rowResult.HasNextResultSet(), ShouldBeTrue)
		So(rowResult.NextResultSet(), ShouldBeNil)
		So(rowResult.HasNextResultSet(), ShouldBeFalse)
		So(rowResult.NextResultSet(), ShouldEqual, io.EOF)
	})
}
