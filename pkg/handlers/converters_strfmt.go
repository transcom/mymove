package handlers

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/unit"
)

// These functions facilitate converting from the go types the db uses
// into the strfmt types that go-swagger uses for payloads.

// FmtDatePtrToPopPtr converts go-swagger type to pop type
func FmtDatePtrToPopPtr(date *strfmt.Date) *time.Time {
	if date == nil {
		return nil
	}

	fmtDate := time.Time(*date)
	return &fmtDate
}

// FmtDateTimePtrToPopPtr converts go-swagger type to pop type
func FmtDateTimePtrToPopPtr(date *strfmt.DateTime) *time.Time {
	if date == nil {
		return nil
	}

	fmtDate := time.Time(*date)
	return &fmtDate
}

// FmtDateTimePtrToPop converts go-swagger type time to pop time
func FmtDateTimePtrToPop(date *strfmt.DateTime) time.Time {
	if date == nil {
		return time.Time{} // Empty time literal
	}
	fmtTime := time.Time(*date)
	return fmtTime
}

// FmtInt64PtrToPopPtr converts go-swagger type to pop type
func FmtInt64PtrToPopPtr(c *int64) *unit.Cents {
	if c == nil {
		return nil
	}

	fmtCents := unit.Cents(*c)
	return &fmtCents
}

// FmtUUIDPtrToPopPtr converts go-swagger uuid type to pop type
func FmtUUIDPtrToPopPtr(u *strfmt.UUID) *uuid.UUID {
	if u == nil {
		return nil
	}
	fmtUUID := uuid.FromStringOrNil(u.String())
	return &fmtUUID
}

// FmtUUIDToPop converts go-swagger uuid type to pop type
func FmtUUIDToPop(u strfmt.UUID) uuid.UUID {
	fmtUUID := uuid.FromStringOrNil(u.String())
	return fmtUUID
}
