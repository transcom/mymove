package models

import (
	"github.com/gobuffalo/validate"
	"github.com/transcom/mymove/pkg/di"
	"time"
)

// StringPointer allows you to take the address of a string literal.
// It is useful for initializing string pointer fields in model construction
func StringPointer(s string) *string {
	return &s
}

// IntPointer allows you to take the address of a int literal.
// It is useful for initializing int pointer fields in model construction
func IntPointer(i int) *int {
	return &i
}

// Int64Pointer allows you to take the address of a int64 literal.
// It is useful for initializing int64 pointer fields in model construction
func Int64Pointer(i int64) *int64 {
	return &i
}

// Float64Pointer allows you to take the address of a float64 literal.
// It is useful for initializing float64 pointer fields in model construction
func Float64Pointer(i float64) *float64 {
	return &i
}

// TimePointer allows you to take the address of a time.Time literal.
// It is useful for initializing time.Time pointer fields in model construction
func TimePointer(t time.Time) *time.Time {
	return &t
}

// BoolPointer allows you to take the address of a bool literal.
// It is useful for initializing bool pointer fields in model construction
func BoolPointer(b bool) *bool {
	return &b
}

// AddProviders adds all the DI providers from the models package
func AddProviders(c *di.Container) {
	c.MustProvide(NewServiceMemberDB)
	c.MustProvide(NewDocumentDB)
}

// ValidationErrors is a type alias to limit the leakage of POP types from the DB layer
type ValidationErrors *validate.Errors
