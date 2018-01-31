package handlers

import (
	"github.com/go-openapi/strfmt"
	"github.com/markbates/pop/nulls"
)

func pointerFromString(s string) *string                      { return &s }
func pointerFromSUUID(u strfmt.UUID) *strfmt.UUID             { return &u }
func pointerFromSDateTime(d strfmt.DateTime) *strfmt.DateTime { return &d }
