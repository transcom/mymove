package models_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func TestDateIsWorkday_IsValid(t *testing.T) {
	calendar := dates.NewUSCalendar()
	t.Run("Valid date", func(t *testing.T) {
		date := time.Date(2019, time.January, 24, 0, 0, 0, 0, time.UTC)
		validator := models.DateIsWorkday{Name: "test_date", Field: date, Calendar: calendar}
		errs := validate.NewErrors()

		validator.IsValid(errs)

		if errs.Count() != 0 {
			t.Fatal("There should be no errors")
		}
	})

	t.Run("Weekend date", func(t *testing.T) {
		date := time.Date(2019, time.January, 26, 0, 0, 0, 0, time.UTC)
		validator := models.DateIsWorkday{Name: "test_date", Field: date, Calendar: calendar}
		errs := validate.NewErrors()

		validator.IsValid(errs)

		testErrors := errs.Get("test_date")
		if len(testErrors) != 1 {
			t.Fatal("There should be an error")
		}
		if testErrors[0] != "cannot be on a weekend or holiday, is 2019-01-26 00:00:00 +0000 UTC" {
			t.Fatal("Did not fail with weekend or holiday error")
		}
	})

	t.Run("Holiday date", func(t *testing.T) {
		date := time.Date(2019, time.January, 21, 0, 0, 0, 0, time.UTC)
		validator := models.DateIsWorkday{Name: "test_date", Field: date, Calendar: calendar}
		errs := validate.NewErrors()

		validator.IsValid(errs)

		testErrors := errs.Get("test_date")
		if len(testErrors) != 1 {
			t.Fatal("There should be an error")
		}
		if testErrors[0] != "cannot be on a weekend or holiday, is 2019-01-21 00:00:00 +0000 UTC" {
			t.Fatal("Did not fail with weekend or holiday error")
		}
	})
}

func TestOptionalDateIsWorkday_IsValid(t *testing.T) {
	calendar := dates.NewUSCalendar()
	t.Run("Valid date", func(t *testing.T) {
		date := dates.NextWorkday(*calendar, time.Date(testdatagen.TestYear, time.January, 24, 0, 0, 0, 0, time.UTC))
		validator := models.OptionalDateIsWorkday{Name: "test_date", Field: &date, Calendar: calendar}
		errs := validate.NewErrors()

		validator.IsValid(errs)

		if errs.Count() != 0 {
			t.Fatal("There should be no errors")
		}
	})

	t.Run("Weekend date", func(t *testing.T) {
		date := dates.NextNonWorkday(*calendar, time.Date(testdatagen.TestYear, time.January, 24, 0, 0, 0, 0, time.UTC))
		validator := models.OptionalDateIsWorkday{Name: "test_date", Field: &date, Calendar: calendar}
		errs := validate.NewErrors()

		validator.IsValid(errs)

		testErrors := errs.Get("test_date")
		if len(testErrors) != 1 {
			t.Fatal("There should be an error")
		}
		stringDate := date.Format("2006-01-02 15:04:05 -0700 UTC")
		if testErrors[0] != fmt.Sprintf("cannot be on a weekend or holiday, is %s", stringDate) {
			t.Fatal("Did not fail with weekend or holiday error")
		}
	})

	t.Run("Holiday date", func(t *testing.T) {
		date := time.Date(testdatagen.TestYear, time.January, 1, 0, 0, 0, 0, time.UTC)
		validator := models.OptionalDateIsWorkday{Name: "test_date", Field: &date, Calendar: calendar}
		errs := validate.NewErrors()

		validator.IsValid(errs)

		testErrors := errs.Get("test_date")
		if len(testErrors) != 1 {
			t.Fatal("There should be an error")
		}
		stringDate := date.Format("2006-01-02 15:04:05 -0700 UTC")
		if testErrors[0] != fmt.Sprintf("cannot be on a weekend or holiday, is %v", stringDate) {
			t.Fatal("Did not fail with weekend or holiday error")
		}
	})

	t.Run("No date", func(t *testing.T) {
		validator := models.OptionalDateIsWorkday{Name: "test_date", Calendar: calendar}
		errs := validate.NewErrors()

		validator.IsValid(errs)

		if errs.Count() != 0 {
			t.Fatal("There should be no errors")
		}
	})
}
