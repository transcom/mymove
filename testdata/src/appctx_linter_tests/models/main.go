package models

import pop "github.com/gobuffalo/pop/v5"

// TestAppContext Test pop connection in struct. This one should be ignored because it's in a package on the allow list.
type TestAppContext struct {
	DB         *pop.Connection // Look for a field whose type is Connection
	testString string
}
