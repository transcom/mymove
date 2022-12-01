package cli

import pop "github.com/gobuffalo/pop/v6"

// These should be ignored because the package is on the allow list.

// TestFuncWithPopConnection func that takes in *pop.Connection as a param.
func TestFuncWithPopConnection(db *pop.Connection) {}
