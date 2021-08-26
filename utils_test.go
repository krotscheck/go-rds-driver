package rds_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/krotscheck/go-rds-driver"
)

// ExpectWakeup can be used whenever we're mocking out a new connection
func ExpectWakeup(mockRDS *MockRDSDataServiceAPI, conf *rds.Config) {

	//ExpectedStatement(conf, "/* wakeup */ SELECT VERSION()", nil)

	mockRDS.EXPECT().
		ExecuteStatement(&rdsdataservice.ExecuteStatementInput{
			Database:    aws.String("database"),
			ResourceArn: aws.String("resourceARN"),
			SecretArn:   aws.String("secretARN"),
			Sql:         aws.String("/* wakeup */ SELECT VERSION()"),
			Parameters:  []*rdsdataservice.SqlParameter{},
		}).
		AnyTimes().
		Return(&rdsdataservice.ExecuteStatementOutput{
			Records: [][]*rdsdataservice.Field{
				{
					&rdsdataservice.Field{
						StringValue: aws.String("5.7.0"),
					},
				},
			},
		}, nil)
}
