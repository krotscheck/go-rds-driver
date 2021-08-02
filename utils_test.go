package rds_test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
	"github.com/krotscheck/rds"
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

	_, err = rdsAPI.ExecuteStatement(&rdsdataservice.ExecuteStatementInput{
		ResourceArn: aws.String(TestConfig.ResourceArn),
		SecretArn:   aws.String(TestConfig.SecretArn),
		Sql:         aws.String(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", TestConfig.Database)),
	})
	if err != nil {
		log.Fatal(err)
	}
}
