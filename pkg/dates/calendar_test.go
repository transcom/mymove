package dates

import (
	"fmt"
	"testing"
	"time"
)

func TestNextWorkday(t *testing.T) {
	var dateTests = []struct {
		name string
		in   time.Time
		out  time.Time
	}{
		{
			"No weekend or holiday",
			time.Date(2019, 1, 24, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 1, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			"Holiday",
			time.Date(2019, 12, 25, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 12, 26, 0, 0, 0, 0, time.UTC),
		},
		{
			"Weekend",
			time.Date(2019, 1, 25, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 1, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			"Weekend and holiday",
			time.Date(2019, 1, 18, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 1, 22, 0, 0, 0, 0, time.UTC),
		},
	}
	cal := NewUSCalendar()
	for _, dt := range dateTests {
		t.Run(dt.name, func(t *testing.T) {
			nextDate := NextWorkday(*cal, dt.in)
			if nextDate != dt.out {
				t.Fatal(fmt.Sprintf("Actual date: %v is not equal to expected date: %v", nextDate, dt.out))
			}
		})
	}
}

func TestNextNonWorkday(t *testing.T) {
	var dateTests = []struct {
		name string
		in   time.Time
		out  time.Time
	}{
		{
			"Saturday after weekday",
			time.Date(2019, 1, 24, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 1, 26, 0, 0, 0, 0, time.UTC),
		},
		{
			"Saturday after Sunday",
			time.Date(2019, 1, 27, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 2, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			"Holiday after weekday",
			time.Date(2019, 12, 23, 0, 0, 0, 0, time.UTC),
			time.Date(2019, 12, 25, 0, 0, 0, 0, time.UTC),
		},
	}
	cal := NewUSCalendar()
	for _, dt := range dateTests {
		t.Run(dt.name, func(t *testing.T) {
			nextDate := NextNonWorkday(*cal, dt.in)
			if nextDate != dt.out {
				t.Fatal(fmt.Sprintf("Actual date: %v is not equal to expected date: %v", nextDate, dt.out))
			}
		})
	}
}
