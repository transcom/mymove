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
type TestAppContext struct { // ADD Want Statement here will be what the test looks for
	DB         Connection // Look for a field whose type is Connection
	testString string
}

// Next Steps:
// Test statments for funcs that take pop.Connection as a parameter or return it as an argument
