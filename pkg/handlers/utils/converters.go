package utils

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// These functions facilitate converting from the go types the db uses
// into the strfmt types that go-swagger uses for payloads.

// FmtUUID converts pop type to go-swagger type
func FmtUUID(u uuid.UUID) *strfmt.UUID {
	fmtUUID := strfmt.UUID(u.String())
	return &fmtUUID
}

// FmtUUIDPtr converts pop type to go-swagger type
func FmtUUIDPtr(u *uuid.UUID) *strfmt.UUID {
	if u == nil {
		return nil
	}
	return FmtUUID(*u)
}

// FmtDateTime converts pop type to go-swagger type
func FmtDateTime(dateTime time.Time) *strfmt.DateTime {
	fmtDateTime := strfmt.DateTime(dateTime)
	return &fmtDateTime
}

// FmtDateTimePtr converts pop type to go-swagger type
func FmtDateTimePtr(dateTime *time.Time) *strfmt.DateTime {
	if dateTime == nil {
		return nil
	}
	return (*strfmt.DateTime)(dateTime)
}

// FmtDate converts pop type to go-swagger type
func FmtDate(date time.Time) *strfmt.Date {
	fmtDate := strfmt.Date(date)
	return &fmtDate
}

// FmtDatePtr converts pop type to go-swagger type
func FmtDatePtr(date *time.Time) *strfmt.Date {
	if date == nil {
		return nil
	}
	return (*strfmt.Date)(date)
}

// FmtPoundPtr converts pop type to go-swagger type
func FmtPoundPtr(weight *unit.Pound) *int64 {
	if weight == nil {
		return nil
	}
	value := weight.Int64()
	return &value
}

// PoundPtrFromInt64Ptr converts pop type to go-swagger type
func PoundPtrFromInt64Ptr(num *int64) *unit.Pound {
	if num == nil {
		return nil
	}
	value := int(*num)
	pound := unit.Pound(value)
	return &pound
}

// FmtURI converts pop type to go-swagger type
func FmtURI(uri string) *strfmt.URI {
	fmtURI := strfmt.URI(uri)
	return &fmtURI
}

// FmtInt64 converts pop type to go-swagger type
func FmtInt64(i int64) *int64 {
	return &i
}

// FmtBool converts pop type to go-swagger type
func FmtBool(b bool) *bool {
	return &b
}

// FmtEmail converts pop type to go-swagger type
func FmtEmail(email string) *strfmt.Email {
	fmtEmail := strfmt.Email(email)
	return &fmtEmail
}

// FmtEmailPtr converts pop type to go-swagger type
func FmtEmailPtr(email *string) *strfmt.Email {
	if email == nil {
		return nil
	}
	return FmtEmail(*email)
}

// StringFromEmail converts pop type to go-swagger type
func StringFromEmail(email *strfmt.Email) *string {
	if email == nil {
		return nil
	}
	emailString := email.String()
	return &emailString
}

// FmtString converts pop type to go-swagger type
func FmtString(s string) *string {
	return &s
}

// FmtSSN converts pop type to go-swagger type
func FmtSSN(s string) *strfmt.SSN {
	ssn := strfmt.SSN(s)
	return &ssn
}

// StringFromSSN converts pop type to go-swagger type
func StringFromSSN(ssn *strfmt.SSN) *string {
	var stringPointer *string
	if ssn != nil {
		plainString := ssn.String()
		stringPointer = &plainString
	}
	return stringPointer
}
