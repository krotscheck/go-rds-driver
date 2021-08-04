# go-rds-driver
A golang sql Driver for the Amazon Aurora Serverless data api.

> **Note:** The serverless data api only supports _named_ query parameters, not ordinal ones. Please write your queries accordingly.

## Getting Started

The `dsn` used in this driver is a JSON encoded string
of the required aws-sdk parameters. The string may be generated
by using the provided `Config` type and its `ToDSN` method.
```go
conf := &rds.Config{
    ResourceArn: "...",
    SecretArn:   "...",
    Database:    "...",
    AWSRegion:   "...",
}
dsn, err := conf.ToDSN()

db.ConnPool, err = sql.Open(rds.DRIVERNAME, dsn)
```

## Usage with Gorm

The above caveat with the Serverless Data API makes usage of gorm tricky. While you can easily restrict yourself
to only using named parameters in your application code, the current implementation of the Gorm
Migrators (as of Aug 1, 2021) exclusively uses ordinal parameters. Please be careful when using this driver.

```go
// RDS using the MySQL Dialector
conf := mysql.Config{
    DriverName: rds.DRIVERNAME,
    DSN:        dsn,
}
dialector := mysql.New(conf)


// RDS using the Postgres Dialector
conf := postgres.Config{
    DriverName: rds.DRIVERNAME,
    DSN:        dsn,
}
dialector := postgres.New(conf)
```

## Running the tests
You will need to provide your own RDS cluster to run tests for this package.
Once provisioned, set the following environment parameters:
```
export RDS_TEST_RESOURCE_ARN=""
export RDS_TEST_SECRET_ARN=""
export RDS_TEST_DATABASE=""
export AWS_REGION=""
```

Executing tests and generating reports can be done via the provided makefile.
```shell
make clean checks
```

## Acknowledgments
This implementation heavily inspired by [what came before](https://github.com/graveyard/rds/tree/birthday).

## TODO List
* Increase test coverage.
* Figure out CI/CD, if necessary.
* Detect database engine underlying it all.
