package appctx_linter_tests

import pop "github.com/gobuffalo/pop/v5"

//type pop struct {
//	Connection string
//}

type AppContext struct {
	ID          string
	Elapsed     int64
	eager       bool
	eagerFields []string
}

// Test pop connection in struct
type TestAppContext struct { // want "Please remove pop.Connection from the struct if not in models"
	DB         *pop.Connection // Look for a field whose type is Connection
	testString string
}

// No want statement because the linter isn't flagged here
func TestAppContextFalse(db *pop.Connection) {}

// NOTE: We don't need a want statement here because we are testing tat the code passes
func TestAppCtxTrueFunc(appCtx AppContext) {}
