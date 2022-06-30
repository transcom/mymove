package weightticket

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type weightTicketUpdater struct {
	checks []weightTicketValidator
}

// NewCustomerWeightTicketUpdater creates a new weightTicketUpdater struct with the checks it needs for a customer
func NewCustomerWeightTicketUpdater() services.WeightTicketUpdater {
	return &weightTicketUpdater{
		checks: basicChecks(),
	}
}

// UpdateWeightTicket updates a weightTicket
func (f *weightTicketUpdater) UpdateWeightTicket(appCtx appcontext.AppContext, weightTicket models.WeightTicket, eTag string) (*models.WeightTicket, error) {
	// get existing WeightTicket
	originalWeightTicket, err := models.FetchWeightTicketByIDExcludeDeletedUploads(appCtx.DB(), weightTicket.ID)
	if err != nil {
		return nil, err
	}

	// verify ETag
	if etag.GenerateEtag(originalWeightTicket.UpdatedAt) != eTag {
		return nil, apperror.NewPreconditionFailedError(originalWeightTicket.ID, nil)
	}

	// merge
	mergedWeightTicket := *originalWeightTicket
	mergedWeightTicket.VehicleDescription = services.SetOptionalStringField(weightTicket.VehicleDescription, mergedWeightTicket.VehicleDescription)
	mergedWeightTicket.EmptyWeight = services.SetNoNilOptionalPoundField(weightTicket.EmptyWeight, mergedWeightTicket.EmptyWeight)
	mergedWeightTicket.MissingEmptyWeightTicket = services.SetNoNilOptionalBoolField(weightTicket.MissingEmptyWeightTicket, mergedWeightTicket.MissingEmptyWeightTicket)
	mergedWeightTicket.FullWeight = services.SetNoNilOptionalPoundField(weightTicket.FullWeight, mergedWeightTicket.FullWeight)
	mergedWeightTicket.MissingFullWeightTicket = services.SetNoNilOptionalBoolField(weightTicket.MissingFullWeightTicket, mergedWeightTicket.MissingFullWeightTicket)
	mergedWeightTicket.OwnsTrailer = services.SetNoNilOptionalBoolField(weightTicket.OwnsTrailer, mergedWeightTicket.OwnsTrailer)
	mergedWeightTicket.TrailerMeetsCriteria = services.SetNoNilOptionalBoolField(weightTicket.TrailerMeetsCriteria, mergedWeightTicket.TrailerMeetsCriteria)

	// validate updated model
	if err := validateWeightTicket(appCtx, &mergedWeightTicket, originalWeightTicket, f.checks...); err != nil {
		return nil, err
	}

	// update the DB record
	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := txnCtx.DB().ValidateAndUpdate(&mergedWeightTicket)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(originalWeightTicket.ID, err, verrs, "invalid input found while updating the WeightTicket")
		} else if err != nil {
			return apperror.NewQueryError("WeightTicket update", err, "")
		}

		return nil
	})

	if txnErr != nil {
		return nil, txnErr
	}

	return &mergedWeightTicket, nil
}
