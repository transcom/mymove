package handlers

import (
	"github.com/go-openapi/strfmt"
	"github.com/markbates/pop/nulls"
)

func pointerFromString(s string) *string                      { return &s }
func pointerFromSUUID(u strfmt.UUID) *strfmt.UUID             { return &u }
func pointerFromSDateTime(d strfmt.DateTime) *strfmt.DateTime { return &d }

func pointerFromNullString(s nulls.String) *string {
	var p *string
	if s.Valid {
		p = &s.String
	}
	return p
}

func nullFromPointerString(s *string) nulls.String {
	n := nulls.String{}
	if s != nil {
		n = nulls.NewString(*s)
	}
	return n
}
