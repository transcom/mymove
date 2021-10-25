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
type TestAppContext struct {
	DB         Connection
	testString string
}
