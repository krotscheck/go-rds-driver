package rds_test

import (
	"database/sql/driver"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata/types"
	"github.com/krotscheck/go-rds-driver"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_Dialect(t *testing.T) {

	Convey("ConvertNamedValue", t, func() {

		Convey("time.Time is formatted as UTC regardless of input timezone", func() {
			helsinki, _ := time.LoadLocation("Europe/Helsinki") // UTC+2/UTC+3
			t := time.Date(2024, 3, 15, 14, 30, 45, 123456000, helsinki)

			namedValue := driver.NamedValue{Name: "ts", Value: t}
			result, err := rds.ConvertNamedValue(namedValue)

			So(err, ShouldBeNil)
			So(result, ShouldResemble, types.SqlParameter{
				Name:     aws.String("ts"),
				TypeHint: types.TypeHintTimestamp,
				Value:    &types.FieldMemberStringValue{Value: "2024-03-15 12:30:45.123456"},
			})
		})

		Convey("time.Time already in UTC is formatted unchanged", func() {
			t := time.Date(2024, 3, 15, 12, 30, 45, 0, time.UTC)

			namedValue := driver.NamedValue{Name: "ts", Value: t}
			result, err := rds.ConvertNamedValue(namedValue)

			So(err, ShouldBeNil)
			So(result, ShouldResemble, types.SqlParameter{
				Name:     aws.String("ts"),
				TypeHint: types.TypeHintTimestamp,
				Value:    &types.FieldMemberStringValue{Value: "2024-03-15 12:30:45"},
			})
		})

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
				So(result, ShouldResemble, types.SqlParameter{
					Name:  aws.String("name"),
					Value: &types.FieldMemberIsNull{Value: true},
				})
			}
		})
	})
}
