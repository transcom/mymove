package handlers

import pop "github.com/gobuffalo/pop/v6"

// TestAppContext Test pop connection in struct.
type TestAppContext struct { // want "Please remove pop.Connection from the struct if not in allowed places. See pkg/appcontext-linter/appctx.go for valid placements."
	DB         *pop.Connection // Look for a field whose type is Connection
	testString string
}

// Config should not be flagged because it's a special exception.
type Config struct {
	DB         *pop.Connection
	testString string
}

// NewHandlerConfig should not be flagged because it's a special exception.
func NewHandlerConfig(db *pop.Connection) {}
