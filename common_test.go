package rds_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/krotscheck/go-rds-driver"
	"os"
	"strings"
)

var TestMysqlConfig *rds.Config
var TestPostgresConfig *rds.Config

type TestConfig struct {
	MysqlDBName    string
	MysqlARN       string
	PostgresDBName string
	PostgresARN    string
	SecretARN      string
	AWSRegion      string
}

func init() {
	conf := &TestConfig{}
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		switch pair[0] {
		case "RDS_MYSQL_DB_NAME":
			conf.MysqlDBName = pair[1]
		case "RDS_MYSQL_ARN":
			conf.MysqlARN = pair[1]
		case "RDS_POSTGRES_DB_NAME":
			conf.PostgresDBName = pair[1]
		case "RDS_POSTGRES_ARN":
			conf.PostgresARN = pair[1]
		case "RDS_SECRET_ARN":
			conf.SecretARN = pair[1]
		case "AWS_REGION":
			conf.AWSRegion = pair[1]
		}
	}

	TestMysqlConfig = rds.NewConfig(conf.MysqlARN, conf.SecretARN, conf.MysqlDBName, conf.AWSRegion)
	TestPostgresConfig = rds.NewConfig(conf.PostgresARN, conf.SecretARN, conf.PostgresDBName, conf.AWSRegion)
}

// ExpectWakeup can be used whenever we're mocking out a new connection
func ExpectWakeup(mockRDS *MockRDSDataServiceAPI, conf *rds.Config) {
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
