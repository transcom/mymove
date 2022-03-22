package services

import (
	"time"

	"github.com/transcom/mymove/pkg/unit"
)

// ValidationFunc is a type representing the signature for a function that validates a service/model
type ValidationFunc func() error

// CheckValidationData runs through a list of ValidationFuncs to check for errors
func CheckValidationData(checks []ValidationFunc) error {
	var err error
	for _, check := range checks {
		err = check()
		if err != nil {
			return err
		}
	}
	return nil
}

// SetOptionalDateTimeField sets the correct new value for the updated date field. Can be nil.
func SetOptionalDateTimeField(newDate *time.Time, oldDate *time.Time) *time.Time {
	// check if the user wanted to keep this field the same:
	if newDate == nil {
		return oldDate
	}

	// check if the user wanted to nullify the value in this field:
	if newDate.IsZero() {
		return nil
	}

	return newDate // return the new intended value
}

// SetOptionalStringField sets the correct new value for the updated string field. Can be nil.
func SetOptionalStringField(newString *string, oldString *string) *string {
	// check if the user wanted to keep this field the same:
	if newString == nil {
		return oldString
	}

	// check if the user wanted to nullify the value in this field:
	if *newString == "" {
		return nil
	}

	return newString // return the new intended value
}

// SetOptionalPoundField sets the correct new value for the updated weight field. Can be nil.
func SetOptionalPoundField(newWeight *unit.Pound, oldWeight *unit.Pound) *unit.Pound {
	// check if the user wanted to keep this field the same:
	if newWeight == nil {
		return oldWeight
	}

	// check if the user wanted to nullify the value in this field:
	if *newWeight == 0 {
		return nil
	}

	return newWeight // return the new intended value
}

// SetNoNilOptionalCentField sets the correct new value for the updated cent field. Can be nil.
func SetNoNilOptionalCentField(newCent *unit.Cents, oldCent *unit.Cents) *unit.Cents {
	// check if the user wanted to keep this field the same:
	if newCent == nil {
		return oldCent
	}

	return newCent // return the new intended value
}

// SetNoNilOptionalPoundField sets the correct new value for the updated weight field.
func SetNoNilOptionalPoundField(newWeight *unit.Pound, oldWeight *unit.Pound) *unit.Pound {
	// check if the user wanted to keep this field the same:
	if newWeight == nil {
		return oldWeight
	}

	return newWeight // return the new intended value
}

// SetNoNNilOptionalInt32Field sets the correct new value for the updated int32 field.
func SetNoNNilOptionalInt32Field(newInt *int32, oldInt *int32) *int32 {
	// check if the user wanted to keep this field the same:
	if newInt == nil {
		return oldInt
	}

	return newInt // return the new intended value
}

// SetNoNilOptionalBoolField sets the correct new value for the updated bool field.
func SetNoNilOptionalBoolField(newBool *bool, oldBool *bool) *bool {
	// check if the user wanted to keep this field the same:
	if newBool == nil {
		return oldBool
	}

	return newBool // return the new intended value
}
