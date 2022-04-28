package appctx_linter_tests

import pop "github.com/gobuffalo/pop/v6"

type AppContext struct {
	ID          string
	Elapsed     int64
	eager       bool
	eagerFields []string
}

// TestAppContext Test pop connection in struct
type TestAppContext struct { // want "Please remove pop.Connection from the struct if not in allowed places. See pkg/appcontext-linter/appctx.go for valid placements."
	DB         *pop.Connection
	testString string
}

// TestHandler Test pop connection in another struct. Want to make sure both get flagged
type TestHandler struct { // want "Please remove pop.Connection from the struct if not in allowed places. See pkg/appcontext-linter/appctx.go for valid placements."
	DB         *pop.Connection
	BackupDB   *pop.Connection
	testString string
}

// TestFuncWithPopConnection func that takes in *pop.Connection as a param.
func TestFuncWithPopConnection(db *pop.Connection) {} // want "Please use appcontext instead of pop.Connection."

// TestFuncWithAppContext NOTE: We don't need a want statement here because we are testing tat the code passes
func TestFuncWithAppContext(appCtx AppContext) {}
