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
	if u == nil {
		return nil
	}
	return fmtUUID(*u)
}

func fmtDateTime(dateTime time.Time) *strfmt.DateTime {
	fmtDateTime := strfmt.DateTime(dateTime)
	return &fmtDateTime
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

func fmtInt64(i int64) *int64 {
	return &i
}

func fmtBool(b bool) *bool {
	return &b
}

func fmtEmail(email string) *strfmt.Email {
	fmtEmail := strfmt.Email(email)
	return &fmtEmail
}

func fmtEmailPtr(email *string) *strfmt.Email {
	if email == nil {
		return nil
	}
	return fmtEmail(*email)
}

func stringFromEmail(email *strfmt.Email) *string {
	if email == nil {
		return nil
	}
	emailString := email.String()
	return &emailString
}

func fmtString(s string) *string {
	return &s
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
