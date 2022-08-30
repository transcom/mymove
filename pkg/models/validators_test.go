package models_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type stringList []string

func (sl stringList) Contains(needle string) bool {
	for _, s := range sl {
		if s == needle {
			return true
		}
	}
	return false
}

func (sl stringList) Contents() []string {
	return sl
}

func TestStringInList_IsValid(t *testing.T) {
	validTypes := stringList{"image/png", "image/jpeg", "application/pdf"}
	for _, validType := range validTypes {
		t.Run(validType, func(t *testing.T) {
			validator := models.NewStringInList(validType, "fieldName", validTypes)

			errs := validate.NewErrors()
			validator.IsValid(errs)

			if errs.Count() != 0 {
				t.Errorf("wrong number of errors; expected %d, got %d", 0, errs.Count())
			}
		})
	}

	invalidTypes := stringList{"image/gif", "video/mp4", "application/octet-stream"}
	for _, invalidType := range invalidTypes {
		t.Run(invalidType, func(t *testing.T) {
			validator := models.NewStringInList(invalidType, "fieldName", validTypes)

			errs := validate.NewErrors()
			validator.IsValid(errs)

			if errs.Count() != 1 {
				t.Fatal("There should be one error")
			}

			expected := fmt.Sprintf("'%s' is not in the list [%s].", invalidType, strings.Join(validTypes, ", "))
			fieldErrors := errs.Get("field_name")

			if len(fieldErrors) != 1 {
				t.Errorf("wrong number of errors; expected %d, got %d", 1, len(fieldErrors))
			}
			if fieldErrors[0] != expected {
				t.Errorf("wrong validation message; expected %s, got %s", expected, fieldErrors[0])
			}
		})
	}
}

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

func TestOptionalStringInclusion_IsValid(t *testing.T) {
	targetStrings := []string{"aaa", "bbb", "ccc"}
	fieldName := "test_string"

	t.Run("String in list", func(t *testing.T) {
		testString := "bbb"
		validator := models.OptionalStringInclusion{Name: fieldName, List: targetStrings, Field: &testString}
		errs := validate.NewErrors()

		validator.IsValid(errs)

		if errs.Count() != 0 {
			t.Fatal("There should be no errors")
		}
	})

	t.Run("String not in list", func(t *testing.T) {
		testString := "zzz"
		validator := models.OptionalStringInclusion{Name: fieldName, List: targetStrings, Field: &testString}
		errs := validate.NewErrors()

		validator.IsValid(errs)

		testErrors := errs.Get(fieldName)
		if len(testErrors) != 1 {
			t.Fatal("There should be an error")
		}

		expected := fmt.Sprintf("%s is not in the list [%s].", fieldName, strings.Join(targetStrings, ", "))
		if testErrors[0] != expected {
			t.Errorf("wrong validation message; expected %s, got %s", expected, testErrors[0])
		}
	})

	t.Run("String is nil", func(t *testing.T) {
		validator := models.OptionalStringInclusion{Name: "test_string", List: targetStrings, Field: nil}
		errs := validate.NewErrors()

		validator.IsValid(errs)

		if errs.Count() != 0 {
			t.Fatal("There should be no errors")
		}
	})
}

func TestFloat64IsPresent_IsValid(t *testing.T) {
	fieldName := "number"

	t.Run("Float64 is non-zero", func(t *testing.T) {
		validator := models.Float64IsPresent{Name: fieldName, Field: 3.14}
		errs := validate.NewErrors()
		validator.IsValid(errs)

		if errs.Count() != 0 {
			t.Fatal("There should be no errors")
		}
	})

	t.Run("Float64 is not provided", func(t *testing.T) {
		validator := models.Float64IsPresent{Name: fieldName}
		errs := validate.NewErrors()
		validator.IsValid(errs)

		if errs.Count() != 1 {
			t.Fatal("There should be one error")
		}

		testErrors := errs.Get(fieldName)
		expected := fmt.Sprintf("%s can not be blank.", validator.Name)
		if testErrors[0] != expected {
			t.Errorf("wrong validation message; expected %s, got %s", expected, testErrors[0])
		}
	})

	t.Run("Float64 is set to 0; has custom error message", func(t *testing.T) {
		customMessage := "Validation failed"
		validator := models.Float64IsPresent{Name: fieldName, Field: 0, Message: customMessage}
		errs := validate.NewErrors()
		validator.IsValid(errs)

		if errs.Count() != 1 {
			t.Fatal("There should be one error")
		}

		testErrors := errs.Get(fieldName)
		if testErrors[0] != customMessage {
			t.Errorf("wrong validation message; expected %s, got %s", customMessage, testErrors[0])
		}
	})
}

func TestFloat64IsGreaterThan_IsValid(t *testing.T) {
	fieldName := "number"

	t.Run("Float64 is greater than compared", func(t *testing.T) {
		validator := models.Float64IsGreaterThan{Name: fieldName, Field: 2, Compared: 1}
		errs := validate.NewErrors()
		validator.IsValid(errs)

		if errs.Count() != 0 {
			t.Fatal("There should be no errors")
		}
	})

	t.Run("Float64 is less than compared", func(t *testing.T) {
		validator := models.Float64IsGreaterThan{Name: fieldName, Field: 1, Compared: 2}
		errs := validate.NewErrors()
		validator.IsValid(errs)

		if errs.Count() != 1 {
			t.Fatal("There should be one error")
		}

		testErrors := errs.Get(fieldName)
		expected := fmt.Sprintf("%f is not greater than %f.", validator.Field, validator.Compared)
		if testErrors[0] != expected {
			t.Errorf("wrong validation message; expected %s, got %s", expected, testErrors[0])
		}
	})

	t.Run("Float64 is equal to compared; has custom error message", func(t *testing.T) {
		customMessage := "Validation failed"
		validator := models.Float64IsGreaterThan{Name: fieldName, Field: 0, Compared: 0, Message: customMessage}
		errs := validate.NewErrors()
		validator.IsValid(errs)

		if errs.Count() != 1 {
			t.Fatal("There should be one error")
		}

		testErrors := errs.Get(fieldName)
		if testErrors[0] != customMessage {
			t.Errorf("wrong validation message; expected %s, got %s", customMessage, testErrors[0])
		}
	})
}

func Test_OptionalUUIDIsPresent(t *testing.T) {
	id, err := uuid.NewV4()

	if err != nil {
		t.Fatal(err)
	}

	// positive tests

	// test with filled id
	v := models.OptionalUUIDIsPresent{Name: "Name", Field: &id}
	errors := validate.NewErrors()
	v.IsValid(errors)
	if errors.Count() > 0 {
		t.Errorf("got errors when should be valid: %v", errors)
	}
	// test with nil pointer
	v = models.OptionalUUIDIsPresent{Name: "Name", Field: nil}
	errors = validate.NewErrors()
	v.IsValid(errors)
	if errors.Count() > 0 {
		t.Errorf("got errors when should be valid: %v", errors)
	}

	// negative test

	// test with empty id, this is equivalent to uuid.Nil
	emptyUUID := uuid.UUID{}
	v = models.OptionalUUIDIsPresent{Name: "Name", Field: &emptyUUID}
	errors = validate.NewErrors()
	v.IsValid(errors)
	if errors.Count() != 1 {
		t.Errorf("got wrong number of errors: %v", errors)
	}
	if errors.Get("name")[0] != "Name can not be blank." {
		t.Errorf("wrong error; expected %s, got %s", "Name can not be blank.", errors.Get("name")[0])
	}
}

func Test_OptionalPoundIsNonNegative_isValid(t *testing.T) {
	name := "pound"
	positiveLb := unit.Pound(10)
	negativeLb := unit.Pound(-20)
	zeroLb := unit.Pound(0)

	t.Run("field with positive value succeeds", func(t *testing.T) {
		v := models.OptionalPoundIsNonNegative{
			Name:  name,
			Field: &positiveLb,
		}
		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() != 0 {
			t.Errorf("got errors when should be valid: %v", errs)
		}
	})

	t.Run("field nil value succeeds", func(t *testing.T) {
		// test with nil pointer
		v := models.OptionalPoundIsNonNegative{Name: "Name", Field: nil}
		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() > 0 {
			t.Errorf("got errors when should be valid: %v", errs)
		}
	})

	t.Run("field with zero value succeeds", func(t *testing.T) {
		v := models.OptionalPoundIsNonNegative{
			Name:  name,
			Field: &zeroLb,
		}
		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() > 0 {
			t.Errorf("got errors when should be valid: %v", errs)
		}
	})

	t.Run("field with negative value fails", func(t *testing.T) {
		v := models.OptionalPoundIsNonNegative{
			Name:  name,
			Field: &negativeLb,
		}
		errs := validate.NewErrors()
		v.IsValid(errs)

		if errs.Count() != 1 {
			t.Fatal("There should be one error")
		}

		testErrors := errs.Get(name)
		expected := fmt.Sprintf("%d is less than zero.", *v.Field)
		if testErrors[0] != expected {
			t.Errorf("wrong validation message; expected %s, got %s", expected, testErrors[0])
		}
	})
}

func Test_OptionalPoundIsPositive_isValid(t *testing.T) {
	name := "pound"
	positiveLb := unit.Pound(10)
	negativeLb := unit.Pound(-20)
	zeroLb := unit.Pound(0)

	t.Run("field with positive value succeeds", func(t *testing.T) {
		v := models.OptionalPoundIsPositive{
			Name:  name,
			Field: &positiveLb,
		}
		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() != 0 {
			t.Errorf("got errors when should be valid: %v", errs)
		}
	})

	t.Run("field nil value succeeds", func(t *testing.T) {
		// test with nil pointer
		v := models.OptionalPoundIsPositive{Name: "Name", Field: nil}
		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() > 0 {
			t.Errorf("got errors when should be valid: %v", errs)
		}
	})

	t.Run("field with negative value fails", func(t *testing.T) {
		v := models.OptionalPoundIsPositive{
			Name:  name,
			Field: &negativeLb,
		}
		errs := validate.NewErrors()
		v.IsValid(errs)

		if errs.Count() != 1 {
			t.Fatal("There should be one error")
		}

		testErrors := errs.Get(name)
		expected := fmt.Sprintf("%d is less than or equal to zero", *v.Field)
		if testErrors[0] != expected {
			t.Errorf("wrong validation message; expected %s, got %s", expected, testErrors[0])
		}
	})

	t.Run("field with zero value fails", func(t *testing.T) {
		v := models.OptionalPoundIsPositive{
			Name:  name,
			Field: &zeroLb,
		}
		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() != 1 {
			t.Fatal("There should be one error")
		}

		testErrors := errs.Get(name)
		expected := fmt.Sprintf("%d is less than or equal to zero", *v.Field)
		if testErrors[0] != expected {
			t.Errorf("wrong validation message; expected %s, got %s", expected, testErrors[0])
		}
	})
}

func Test_MustBeBothNilOrBothNotNil_IsValid(t *testing.T) {
	make := "Honda"
	model := "Civic"

	t.Run("fields both have value succeeds", func(t *testing.T) {

		v := models.MustBeBothNilOrBothHaveValue{
			FieldName1:  "VehicleMake",
			FieldValue1: &make,
			FieldName2:  "VehicleModel",
			FieldValue2: &model,
		}
		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() != 0 {
			t.Errorf("got errors when should be valid: %v", errs)
		}
	})

	t.Run("fields are both nil succeeds", func(t *testing.T) {
		v := models.MustBeBothNilOrBothHaveValue{
			FieldName1:  "VehicleMake",
			FieldValue1: nil,
			FieldName2:  "VehicleModel",
			FieldValue2: nil,
		}
		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() != 0 {
			t.Errorf("got errors when should be valid: %v", errs)
		}
	})

	t.Run("first field is nil and the second field has value fails", func(t *testing.T) {
		v := models.MustBeBothNilOrBothHaveValue{
			FieldName1:  "VehicleMake",
			FieldValue1: nil,
			FieldName2:  "VehicleModel",
			FieldValue2: &model,
		}

		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() == 0 {
			t.Errorf("should throw an error if %v is empty and %v is filled: %v", v.FieldName1, v.FieldName2, errs)
		}
	})

	t.Run("first field has value and the second field is nil fails", func(t *testing.T) {
		v := models.MustBeBothNilOrBothHaveValue{
			FieldName1:  "VehicleMake",
			FieldValue1: &make,
			FieldName2:  "VehicleModel",
			FieldValue2: nil,
		}

		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() == 0 {
			t.Errorf("should throw an error if %v is filled and %v is empty: %v", v.FieldName1, v.FieldName2, errs)
		}
	})
}

func Test_ItemCanFitInsideCrate_IsValid(t *testing.T) {
	makeInt32 := func(i int) *int32 {
		val := int32(i)
		return &val
	}

	item := primemessages.MTOServiceItemDimension{
		Height: makeInt32(10000),
		Length: makeInt32(10000),
		Width:  makeInt32(10000),
	}
	crate := primemessages.MTOServiceItemDimension{
		Height: makeInt32(20000),
		Length: makeInt32(20000),
		Width:  makeInt32(20000),
	}

	t.Run("item and crate both have values succeeds", func(t *testing.T) {
		v := models.ItemCanFitInsideCrate{
			Name:         "Item",
			Item:         &item,
			NameCompared: "Crate",
			Crate:        &crate,
		}
		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() != 0 {
			t.Errorf("got errors when should be valid: %v", errs)
		}
	})

	t.Run("item and crate both nil fails", func(t *testing.T) {
		v := models.ItemCanFitInsideCrate{
			Name:         "Item",
			Item:         nil,
			NameCompared: "Crate",
			Crate:        nil,
		}
		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() == 0 {
			t.Errorf("got no errors when should be invalid")
		}
	})

	t.Run("item and crate dimension nil values fails", func(t *testing.T) {
		v := models.ItemCanFitInsideCrate{
			Name:         "Item",
			Item:         &primemessages.MTOServiceItemDimension{},
			NameCompared: "Crate",
			Crate:        &primemessages.MTOServiceItemDimension{},
		}
		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() == 0 {
			t.Errorf("got no errors when should be invalid")
		}
	})

	t.Run("item dimensions greater than or equal to crate dimensions fails", func(t *testing.T) {
		v := models.ItemCanFitInsideCrate{
			Name: "Item",
			Item: &primemessages.MTOServiceItemDimension{
				Height: makeInt32(0),
				Length: makeInt32(0),
				Width:  makeInt32(0),
			},
			NameCompared: "Crate",
			Crate: &primemessages.MTOServiceItemDimension{
				Height: makeInt32(0),
				Length: makeInt32(0),
				Width:  makeInt32(0),
			},
		}
		errs := validate.NewErrors()
		v.IsValid(errs)
		if errs.Count() == 0 {
			t.Errorf("got no errors when should be invalid")
		}
	})
}
