package rds

import "fmt"

// ErrNoPositional is thrown when ordinal parameter keys are used
var ErrNoPositional = fmt.Errorf("only named parameters supported")
