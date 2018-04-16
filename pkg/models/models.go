package models

// StringPointer allows you to take the address of a string literal.
// It is useful for initializing string pointer fields in model construction
func StringPointer(s string) *string {
	return &s
}
