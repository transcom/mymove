package models

import pop "github.com/gobuffalo/pop/v6"

// These should be ignored because the package is on the allow list.

// TestAppContext Test pop connection in struct
type TestAppContext struct {
	DB         *pop.Connection // Look for a field whose type is Connection
	testString string
}

// TestFuncWithPopConnection func that takes in *pop.Connection as a param.
func TestFuncWithPopConnection(db *pop.Connection) {}
