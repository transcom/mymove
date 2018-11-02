package dates

import (
	"log"
	"testing"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/logging/hnyzap"
)

type testCase struct {
	name                       string
	dates                      []time.Time
	includeWeekendsAndHolidays bool
}

func (suite *DatesSuite) TestCreateFutureMoveDates() {
	moveDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	numDays := 7
	usCalendar := NewUSCalendar()

	var cases = []testCase{
		{
			name: "future dates no weekends or holidays",
			dates: []time.Time{
				time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 18, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 19, 0, 0, 0, 0, time.UTC),
			},
			includeWeekendsAndHolidays: false,
		},
		{
			name: "future dates with weekends or holidays",
			dates: []time.Time{
				time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 12, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 14, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 16, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 17, 0, 0, 0, 0, time.UTC),
			},
			includeWeekendsAndHolidays: true,
		},
	}
	for _, testCase := range cases {
		suite.T().Run(testCase.name, func(t *testing.T) {
			dates := CreateFutureMoveDates(moveDate, numDays, testCase.includeWeekendsAndHolidays, usCalendar)
			suite.Equal(testCase.dates, dates, "%v: Future dates did not match, expected %v, got %v", testCase.name, testCase.dates, dates)
		})
	}
}

func (suite *DatesSuite) TestCreatePastMoveDates() {
	moveDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	numDays := 7
	usCalendar := NewUSCalendar()

	var cases = []testCase{
		{
			name: "past dates no weekends or holidays",
			dates: []time.Time{
				time.Date(2018, 12, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 4, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 5, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC),
			},
			includeWeekendsAndHolidays: false,
		},
		{
			name: "past dates with weekends or holidays",
			dates: []time.Time{
				time.Date(2018, 12, 5, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 9, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC),
			},
			includeWeekendsAndHolidays: true,
		},
	}
	for _, testCase := range cases {
		suite.T().Run(testCase.name, func(t *testing.T) {
			dates := CreatePastMoveDates(moveDate, numDays, testCase.includeWeekendsAndHolidays, usCalendar)
			suite.Equal(testCase.dates, dates, "%v: Past dates did not match, expected %v, got %v", testCase.name, testCase.dates, dates)
		})
	}
}

func (suite *DatesSuite) TestCreateValidDatesBetweenTwoDatesEndDateMustBeLater() {
	startDate := time.Date(2018, 12, 11, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2018, 12, 7, 0, 0, 0, 0, time.UTC)
	usCalendar := NewUSCalendar()
	_, err := CreateValidDatesBetweenTwoDates(startDate, endDate, true, false, usCalendar)
	suite.Error(err)
}

type DatesSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *hnyzap.Logger
}

func (suite *DatesSuite) SetupTest() {
	suite.db.TruncateAll()
}

func TestDatesSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &DatesSuite{
		db:     db,
		logger: &hnyzap.Logger{Logger: logger},
	}
	suite.Run(t, hs)
}
