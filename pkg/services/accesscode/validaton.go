package accesscode

import (
	"database/sql"
	"strings"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

type validator interface {
	Validate(appCtx appcontext.AppContext, ac *models.AccessCode) error
}

type validatorFunc func(appcontext.AppContext, *models.AccessCode) error

func (fn validatorFunc) Validate(appCtx appcontext.AppContext, ac *models.AccessCode) error {
	return fn(appCtx, ac)
}

func validateAccessCode(appCtx appcontext.AppContext, ac *models.AccessCode, moveType models.SelectedMoveType, checks ...validator) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, ac); err != nil {
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
		result = apperror.NewInvalidInputError(ac.ID, nil, verrs, "Invalid access code.")
	}
	return result
}

// ValidateAccessCode validates an access code based upon the code and move type. A valid access
// code is assumed to have no `service_member_id`
func checkAccessCode() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, ac *models.AccessCode) error {

		verrs := validate.NewErrors()

		if appCtx.Session() == nil {
			verrs.Add("Unauthorized User", "Active Session")
		}

		splitParams := strings.Split(ac.Code, "-")
		moveType, stringCode := splitParams[0], splitParams[1]

		var code models.AccessCode

		err := appCtx.DB().
			Where("code = ?", stringCode).
			Where("move_type = ?", moveType).
			First(&code)

		if err != nil {
			switch err {
			case sql.ErrNoRows:
				verrs.Add("AccessCode", ac.ID.String()+" Not Found")
			default:
				verrs.Add("AccessCode", err.Error())
			}
		}

		return verrs
	})
}
