package handlers

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"
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

func timeFromDateTime(dateTime strfmt.DateTime) time.Time {
	fmtTime := time.Time(dateTime)
	return fmtTime
}

func fmtDate(date time.Time) *strfmt.Date {
	fmtDate := strfmt.Date(date)
	return &fmtDate
}

func fmtDatePtr(date *time.Time) *strfmt.Date {
	if date == nil {
		return nil
	}
	return (*strfmt.Date)(date)
}

func fmtURI(uri string) *strfmt.URI {
	fmtURI := strfmt.URI(uri)
	return &fmtURI
}

func fmtInt64(i int) *int64 {
	fmtInt := int64(i)
	return &fmtInt
}

func fmtBool(b bool) *bool {
	return &b
}

func fmtEmail(email string) *strfmt.Email {
	fmtEmail := strfmt.Email(email)
	return &fmtEmail
}

func stringFromEmail(email *strfmt.Email) *string {
	if email == nil {
		return nil
	}
	emailString := email.String()
	return &emailString
}

func fmtSSN(s string) *strfmt.SSN {
	ssn := strfmt.SSN(s)
	return &ssn
}

func stringFromSSN(ssn *strfmt.SSN) *string {
	var stringPointer *string
	if ssn != nil {
		plainString := ssn.String()
		stringPointer = &plainString
	}
	return stringPointer
}
