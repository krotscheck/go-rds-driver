package rds_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/krotscheck/rds"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"os"
)

// TestConfig to use when making integration test calls
var TestConfig *rds.Config

func init() {
	// Testing requires a few environment parameters
	resourceARN := os.Getenv("RDS_TEST_RESOURCE_ARN")
	if resourceARN == "" {
		log.Fatal("Missing test environment parameter: RDS_TEST_RESOURCE_ARN")
	}
	secretARN := os.Getenv("RDS_TEST_SECRET_ARN")
	if secretARN == "" {
		log.Fatal("Missing test environment parameter: RDS_TEST_SECRET_ARN")
	}
	database := os.Getenv("RDS_TEST_DATABASE")
	if database == "" {
		log.Fatal("Missing test environment parameter: RDS_TEST_DATABASE")
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		log.Fatal("Missing test environment parameter: AWS_REGION")
	}

	TestConfig = rds.NewConfig(resourceARN, secretARN, database, region)

	// Make sure the database exists...
	awsConfig := aws.NewConfig().WithRegion(TestConfig.AWSRegion)
	awsSession, err := session.NewSession(awsConfig)
	if err != nil {
		log.Fatal(err)
	}
	rdsAPI := rdsdataservice.New(awsSession)

	// Wakeup the cluster
	err = rds.Wakeup(rdsAPI, resourceARN, secretARN, database)
	if err != nil {
		log.Fatal(err)
	}

	_, err = rdsAPI.ExecuteStatement(&rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(TestConfig.ResourceArn),
		SecretArn:   aws.String(TestConfig.SecretArn),
		Sql:         aws.String(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", TestConfig.Database)),
	})
	if err != nil {
		log.Fatal(err)
	}
}

// ExpectWakeup can be used whenever we're mocking out a new connection
func ExpectWakeup(mockRDS *MockRDSDataServiceAPI, conf *rds.Config) {
	mockRDS.EXPECT().ExecuteStatement(ExpectedStatement(conf, "/* wakeup */ SELECT 1", nil)).AnyTimes().Return(nil, nil)
}

// ExpectTransaction to be started
func ExpectTransaction(ctx context.Context, mockRDS *MockRDSDataServiceAPI, conf *rds.Config, transactionID string, readonly bool, isolation sql.IsolationLevel) {
	mockRDS.EXPECT().
		BeginTransactionWithContext(ctx, &rdsdataservice.BeginTransactionInput{
			Database:    aws.String(conf.Database),
			ResourceArn: aws.String(conf.ResourceArn),
			SecretArn:   aws.String(conf.SecretArn),
		}).
		Times(1).
		Return(&rdsdataservice.BeginTransactionOutput{TransactionId: aws.String(transactionID)}, nil)

	mockRDS.EXPECT().
		ExecuteStatementWithContext(ctx, &rdsdataservice.ExecuteStatementInput{
			Database:    aws.String(conf.Database),
			ResourceArn: aws.String(conf.ResourceArn),
			SecretArn:   aws.String(conf.SecretArn),
			Sql:         aws.String("SET TRANSACTION ISOLATION LEVEL :isolation, :readonly"),
			Parameters: []*rdsdataservice.SqlParameter{
				{
					Name:  aws.String("isolation"),
					Value: &rdsdataservice.Field{StringValue: aws.String(isolation.String())},
				},
				{
					Name:  aws.String("readonly"),
					Value: &rdsdataservice.Field{StringValue: aws.String("READ WRITE")},
				},
			},
		}).
		Times(1).
		Return(&rdsdataservice.BeginTransactionOutput{TransactionId: aws.String(transactionID)}, nil)
}

// ExpectQuery in a test
func ExpectedStatement(conf *rds.Config, query string, args []driver.NamedValue) *rdsdataservice.ExecuteStatementInput {
	params, err := rds.ConvertNamedValues(args)
	So(err, ShouldBeNil)
	return &rdsdataservice.ExecuteStatementInput{
		Database:    aws.String(conf.Database),
		ResourceArn: aws.String(conf.ResourceArn),
		SecretArn:   aws.String(conf.SecretArn),
		Sql:         aws.String(query),
		Parameters:  params,
	}
}
