package weightticket

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/unit"
)

type weightTicketUpdater struct {
	checks []weightTicketValidator
	services.WeightTicketFetcher
	ppmShipmentUpdater services.PPMShipmentUpdater
}

// NewCustomerWeightTicketUpdater creates a new weightTicketUpdater struct with the checks it needs for a customer
func NewCustomerWeightTicketUpdater(fetcher services.WeightTicketFetcher, ppmUpdater services.PPMShipmentUpdater) services.WeightTicketUpdater {
	return &weightTicketUpdater{
		checks:              basicChecksForCustomer(),
		WeightTicketFetcher: fetcher,
		ppmShipmentUpdater:  ppmUpdater,
	}
}

func NewOfficeWeightTicketUpdater(fetcher services.WeightTicketFetcher, ppmUpdater services.PPMShipmentUpdater) services.WeightTicketUpdater {
	return &weightTicketUpdater{
		checks:              basicChecksForOffice(),
		WeightTicketFetcher: fetcher,
		ppmShipmentUpdater:  ppmUpdater,
	}
}

// UpdateWeightTicket updates a weightTicket
func (f *weightTicketUpdater) UpdateWeightTicket(appCtx appcontext.AppContext, weightTicket models.WeightTicket, eTag string) (*models.WeightTicket, error) {
	// get existing WeightTicket
	originalWeightTicket, err := f.GetWeightTicket(appCtx, weightTicket.ID)
	if err != nil {
		return nil, err
	}

	// verify ETag
	if etag.GenerateEtag(originalWeightTicket.UpdatedAt) != eTag {
		return nil, apperror.NewPreconditionFailedError(originalWeightTicket.ID, nil)
	}

	mergedWeightTicket := mergeWeightTicket(weightTicket, *originalWeightTicket)

	// validate updated model
	if err = validateWeightTicket(appCtx, &mergedWeightTicket, originalWeightTicket, f.checks...); err != nil {
		return nil, err
	}

	hasTotalWeightChanged := hasTotalWeightChanged(*originalWeightTicket, mergedWeightTicket)
	var currentPPMShipment models.PPMShipment
	if hasTotalWeightChanged {
		ppmShipmentFromDB, issue := ppmshipment.FindPPMShipmentAndWeightTickets(appCtx, originalWeightTicket.PPMShipmentID)
		if issue != nil {
			return nil, issue
		}
		currentPPMShipment = *ppmShipmentFromDB
		for i := range currentPPMShipment.WeightTickets {
			if currentPPMShipment.WeightTickets[i].ID == mergedWeightTicket.ID {
				currentPPMShipment.WeightTickets[i] = mergedWeightTicket
			}
		}
	}

	// update the DB record
	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		// if weight changes call update PPMShipment with new weightTicket
		if hasTotalWeightChanged {
			_, err = f.ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(txnCtx, &currentPPMShipment, currentPPMShipment.ShipmentID)
			if err != nil {
				return err
			}
		}

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

func mergeWeightTicket(weightTicket models.WeightTicket, originalWeightTicket models.WeightTicket) models.WeightTicket {
	mergedWeightTicket := originalWeightTicket

	mergedWeightTicket.VehicleDescription = services.SetOptionalStringField(weightTicket.VehicleDescription, mergedWeightTicket.VehicleDescription)
	mergedWeightTicket.EmptyWeight = services.SetNoNilOptionalPoundField(weightTicket.EmptyWeight, mergedWeightTicket.EmptyWeight)
	mergedWeightTicket.MissingEmptyWeightTicket = services.SetNoNilOptionalBoolField(weightTicket.MissingEmptyWeightTicket, mergedWeightTicket.MissingEmptyWeightTicket)
	mergedWeightTicket.FullWeight = services.SetNoNilOptionalPoundField(weightTicket.FullWeight, mergedWeightTicket.FullWeight)
	mergedWeightTicket.MissingFullWeightTicket = services.SetNoNilOptionalBoolField(weightTicket.MissingFullWeightTicket, mergedWeightTicket.MissingFullWeightTicket)
	mergedWeightTicket.OwnsTrailer = services.SetNoNilOptionalBoolField(weightTicket.OwnsTrailer, mergedWeightTicket.OwnsTrailer)
	mergedWeightTicket.TrailerMeetsCriteria = services.SetNoNilOptionalBoolField(weightTicket.TrailerMeetsCriteria, mergedWeightTicket.TrailerMeetsCriteria)
	mergedWeightTicket.AdjustedNetWeight = services.SetNoNilOptionalPoundField(weightTicket.AdjustedNetWeight, mergedWeightTicket.AdjustedNetWeight)
	mergedWeightTicket.NetWeightRemarks = services.SetOptionalStringField(weightTicket.NetWeightRemarks, mergedWeightTicket.NetWeightRemarks)
	mergedWeightTicket.Reason = services.SetOptionalStringField(weightTicket.Reason, mergedWeightTicket.Reason)
	status := services.SetOptionalStringField((*string)(weightTicket.Status), (*string)(mergedWeightTicket.Status))
	if status != nil {
		ppmDocStatus := models.PPMDocumentStatus(*status)
		mergedWeightTicket.Status = &ppmDocStatus
	} else {
		mergedWeightTicket.Status = nil
	}

	return mergedWeightTicket
}

func hasTotalWeightChanged(originalWeightTicket, newWeightTicket models.WeightTicket) bool {
	var newWeight unit.Pound
	var oldWeight unit.Pound

	if newWeightTicket.FullWeight != nil && newWeightTicket.EmptyWeight != nil {
		newWeight = *newWeightTicket.FullWeight - *newWeightTicket.EmptyWeight
	}
	if originalWeightTicket.FullWeight != nil && originalWeightTicket.EmptyWeight != nil {
		oldWeight = *originalWeightTicket.FullWeight - *originalWeightTicket.EmptyWeight
	}

	return newWeight != oldWeight
}
