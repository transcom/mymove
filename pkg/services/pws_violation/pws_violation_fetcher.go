package pwsviolation

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type pwsViolationsFetcher struct {
}

func NewPWSViolationsFetcher() services.PWSViolationsFetcher {
	return &pwsViolationsFetcher{}
}

func (o pwsViolationsFetcher) GetPWSViolations(appCtx appcontext.AppContext) (*models.PWSViolations, error) {

	pwsViolations := models.PWSViolations{}
	err := appCtx.DB().Order("display_order asc").All(&pwsViolations)
	if err != nil {
		return nil, apperror.NewQueryError("PWSViolations", err, "")
	}
	return &pwsViolations, nil
}
