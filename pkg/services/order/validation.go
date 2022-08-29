package order

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// Validator is the interface for the various validations we might want to
// define.
type Validator interface {
	Validate(order *models.Order) error
}

type validatorFunc func(*models.Order) error

func (fn validatorFunc) Validate(order *models.Order) error {
	return fn(order)
}

// ValidateOrder accepts a range of validator functions and runs each one
func ValidateOrder(order *models.Order, checks ...Validator) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(order); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				// accumulate validation errors
				verrs.Append(e)
			default:
				// non-validation errors have priority,
				// and short-circuit doing any further checks
				return err
			}
		}
	}
	if verrs.HasAny() {
		result = apperror.NewInvalidInputError(order.ID, nil, verrs, "")
	}
	return result
}

// StringIsPresentAfterSubmission checks presence of fields after an order has been submitted
// This validation only applies when a TOO is updating orders.
type StringIsPresentAfterSubmission struct {
	Name  string
	Field *string
}

// IsValid adds an error if the field is blank
func (validate *StringIsPresentAfterSubmission) IsValid(errors *validate.Errors) {
	if validate.Field == nil || *validate.Field == "" {
		errors.Add(validators.GenerateKey(validate.Name), fmt.Sprintf("%s cannot be blank.", validate.Name))
	}
}

// OrdersTypeDetailIsPresentAfterSubmission validates that orders type field is present
type OrdersTypeDetailIsPresentAfterSubmission struct {
	Name  string
	Field *internalmessages.OrdersTypeDetail
}

// IsValid adds an error if the string value is blank.
func (validate *OrdersTypeDetailIsPresentAfterSubmission) IsValid(errors *validate.Errors) {
	if validate.Field == nil || string(*validate.Field) == "" {
		errors.Add(validators.GenerateKey(validate.Name), fmt.Sprintf("%s cannot be blank.", validate.Name))
	}
}

// CheckRequiredFields ensures the presence of certain order fields before a TOO
// can approve a move or update an order.
func CheckRequiredFields() Validator {
	return validatorFunc(func(order *models.Order) error {
		verrs := validate.Validate(
			&StringIsPresentAfterSubmission{Name: "TransportationAccountingCode", Field: order.TAC},
			&StringIsPresentAfterSubmission{Name: "DepartmentIndicator", Field: order.DepartmentIndicator},
			&StringIsPresentAfterSubmission{Name: "OrdersNumber", Field: order.OrdersNumber},
			&OrdersTypeDetailIsPresentAfterSubmission{Name: "OrdersTypeDetail", Field: order.OrdersTypeDetail},
		)

		return verrs
	})
}
