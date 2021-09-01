package rds

import (
	"encoding/json"
)

// Config struct used to provide AWS Configuration Credentials
type Config struct {
	ResourceArn string `json:"resource_arn"`
	SecretArn   string `json:"secret_arn"`
	Database    string `json:"database"`
	AWSRegion   string `json:"aws_region"`
	ParseTime   bool   `json:"parse_time"`
}

// ToDSN converts the config to a DSN string
func (o *Config) ToDSN() (string, error) {
	str, err := json.Marshal(o)
	return string(str), err
}

// NewConfigFromDSN assumes that the DSN is a JSON-encoded string
func NewConfigFromDSN(dsn string) (conf *Config, err error) {
	conf = &Config{}
	err = json.Unmarshal([]byte(dsn), conf)
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
