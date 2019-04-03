package handlers

import (
	"time"

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

// FmtInt64PtrToPopPtr converts go-swagger type to pop type
func FmtInt64PtrToPopPtr(c *int64) *unit.Cents {
	if c == nil {
		return nil
	}

	fmtCents := unit.Cents(*c)
	return &fmtCents
}
