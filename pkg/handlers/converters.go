package handlers

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/satori/go.uuid"
)

// These functions facilitate converting from the go types the db uses
// into the strfmt types that go-swagger uses for payloads.
func fmtUUID(u uuid.UUID) *strfmt.UUID {
	fmtUUID := strfmt.UUID(u.String())
	return &fmtUUID
}

func fmtUUIDPtr(u *uuid.UUID) *strfmt.UUID {
	var su *strfmt.UUID
	if u == nil {
		su = nil
	} else {
		su = fmtUUID(*u)
	}
	return su
}

func fmtDateTime(dateTime time.Time) *strfmt.DateTime {
	fmtDateTime := strfmt.DateTime(dateTime)
	return &fmtDateTime
}

func fmtDate(date time.Time) *strfmt.Date {
	fmtDate := strfmt.Date(date)
	return &fmtDate
}

func fmtInt64(i int) *int64 {
	fmtInt := int64(i)
	return &fmtInt
}

func stringFromSSN(ssn *strfmt.SSN) *string {
	var stringPointer *string
	if ssn != nil {
		plainString := ssn.String()
		stringPointer = &plainString
	}
	return stringPointer
}
