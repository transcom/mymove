package services

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
)

// DateSelectionChecker is the service object interface to determine business rule for given date
//
//go:generate mockery --name DateSelectionChecker
type DateSelectionChecker interface {
	IsDateWeekendHoliday(appCtx appcontext.AppContext, countryCode string, date time.Time) (*IsDateWeekendHolidayInfo, error)
}

type IsDateWeekendHolidayInfo struct {
	CountryCode string
	CountryName string
	Date        time.Time
	IsWeekend   bool
	IsHoliday   bool
}
