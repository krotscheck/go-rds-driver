package rds

import "fmt"

// ErrNoMixedParams is thrown if parameters are mixed
var ErrNoMixedParams = fmt.Errorf("please do not mix ordinal and named parameters")

// ErrClosed indicates that the connection is closed
var ErrClosed = fmt.Errorf("this connection is closed")
