package rds_test

import (
	"database/sql/driver"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata/types"
	"github.com/krotscheck/go-rds-driver"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Dialect(t *testing.T) {

	Convey("ConvertNamedValue", t, func() {

		Convey("Null Values", func() {
			var UInt8 []uint8
			var UInt8Ptr *[]uint8

			values := []driver.Value{
				UInt8,
				UInt8Ptr,
			}

			for _, v := range values {
				namedValue := driver.NamedValue{
					Name:  "name",
					Value: v,
				}
				result, err := rds.ConvertNamedValue(namedValue)
				So(err, ShouldBeNil)
				So(result, ShouldResemble, &types.SqlParameter{
					Name: aws.String("name"),
					Value: &types.Field{
						IsNull: aws.Bool(true),
					},
				})
			}
		})
	})
}
