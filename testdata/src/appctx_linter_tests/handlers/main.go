package handlers

import pop "github.com/gobuffalo/pop/v5"

// TestAppContext Test pop connection in struct.
type TestAppContext struct { // want "Please remove pop.Connection from the struct if not in appcontext"
	DB         *pop.Connection // Look for a field whose type is Connection
	testString string
}

// handlerContext should not be flagged because it's a special exception.
type handlerContext struct {
	DB         *pop.Connection
	testString string
}

// TestFuncWithPopConnection func that takes in *pop.Connection as a param.
func TestFuncWithPopConnection(db *pop.Connection) {} // want "Please use appcontext instead of pop.Connection"
