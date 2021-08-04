package rds

import "fmt"

// ErrNoPositional is thrown when ordinal parameter keys are used
var ErrNoPositional = fmt.Errorf("only named parameters supported")

// ErrClosed indicates that the connection is closed
var ErrClosed = fmt.Errorf("this connection is closed")
