package models

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// EagerAssociations are a collection of named associations
type EagerAssociations []string

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

// Int32Pointer allows you to take the address of a int32 literal.
// It is useful for initializing int32 pointer fields in model construction
func Int32Pointer(i int32) *int32 {
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

// PoundPointer allows you to get the pointer to a unit.Pound literal.
// It is useful for initializing unit.Pound pointer fields in model construction
func PoundPointer(p unit.Pound) *unit.Pound {
	return &p
}

// CentPointer allows you to get the pointer to a unit.Cent literal.
// It is useful for initializing unit.Cent pointer fields in model construction
func CentPointer(c unit.Cents) *unit.Cents {
	return &c
}

// UUIDPointer allows you to get the pointer to a uuid.UUID literal.
// It is useful for initializing uuid.UUID pointer fields in model construction
func UUIDPointer(u uuid.UUID) *uuid.UUID {
	return &u
}
