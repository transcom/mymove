package appctx_linter_tests

//type pop struct {
//	Connection string
//}

type Connection struct {
	ID          string
	Elapsed     int64
	eager       bool
	eagerFields []string
}

// Test pop connection in struct
// Summary: [linter] - [linter type code] - [Linter summary] // want "Please remove pop.Connection from struct if not in models and use appContext" #
type TestAppContextFalse struct { // ADD Want Statement here will be what the test looks for
	DB         Connection // Look for a field whose type is Connection
	testString string
}

// NOTE: We don't need a want statement here because we are testing tat the code passes
type TestAppContextTrue struct {
	appCtx     appCtx.DB // Look for a field whose type is Connection
	testString string
}

// Next Steps:
// Test statements for structs that take pop.Connection as a parameter or return it as an argument
