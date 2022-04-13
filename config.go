package rds

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	keyResourceARN = "resource_arn"
	keySecretARN   = "secret_arn"
	keyDatabase    = "database"
	keyAWSRegion   = "aws_region"
	keyParseTime   = "parse_time"
)

// Config struct used to provide AWS Configuration Credentials
type Config struct {
	ResourceArn string
	SecretArn   string
	Database    string
	AWSRegion   string
	ParseTime   bool
}

// ToDSN converts the config to a DSN string
func (o *Config) ToDSN() string {
	v := url.Values{}
	v.Add(keyResourceARN, o.ResourceArn)
	v.Add(keySecretARN, o.SecretArn)
	v.Add(keyDatabase, o.Database)
	v.Add(keyAWSRegion, o.AWSRegion)
	v.Add(keyParseTime, strconv.FormatBool(o.ParseTime))

	return fmt.Sprintf("%s://?%s", DRIVERNAME, v.Encode())
}

// NewConfigFromDSN assumes that the DSN is a JSON-encoded string
func NewConfigFromDSN(dsn string) (conf *Config, err error) {
	conf = &Config{}

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	// Make sure the scheme is correct.
	if u.Scheme != DRIVERNAME {
		return nil, ErrInvalidDSNScheme
	}

	// All the actual data is in the Query
	values := u.Query()
	conf.ResourceArn = values.Get(keyResourceARN)
	conf.SecretArn = values.Get(keySecretARN)
	conf.Database = values.Get(keyDatabase)
	conf.AWSRegion = values.Get(keyAWSRegion)

	// Swallow the error here because default is fine.
	parseTime, _ := strconv.ParseBool(values.Get(keyParseTime))
	conf.ParseTime = parseTime

	return
}

// NewConfig with specific
func NewConfig(resourceARN string, secretARN string, database string, awsRegion string) (conf *Config) {
	return &Config{
		ResourceArn: resourceARN,
		SecretArn:   secretARN,
		Database:    database,
		AWSRegion:   awsRegion,
	}
}
