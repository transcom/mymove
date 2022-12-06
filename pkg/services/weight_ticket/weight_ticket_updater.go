package weightticket

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type weightTicketUpdater struct {
	checks             []weightTicketValidator
	ppmShipmentUpdater services.PPMShipmentUpdater
}

// NewCustomerWeightTicketUpdater creates a new weightTicketUpdater struct with the checks it needs for a customer
func NewCustomerWeightTicketUpdater(ppmUpdater services.PPMShipmentUpdater) services.WeightTicketUpdater {
	return &weightTicketUpdater{
		checks:             basicChecksForCustomer(),
		ppmShipmentUpdater: ppmUpdater,
	}
}

func NewOfficeWeightTicketUpdater(ppmUpdater services.PPMShipmentUpdater) services.WeightTicketUpdater {
	return &weightTicketUpdater{
		checks:             basicChecksForOffice(),
		ppmShipmentUpdater: ppmUpdater,
	}
}

// UpdateWeightTicket updates a weightTicket
func (f *weightTicketUpdater) UpdateWeightTicket(appCtx appcontext.AppContext, weightTicket models.WeightTicket, eTag string) (*models.WeightTicket, error) {
	// get existing WeightTicket
	// tk we have the id
	originalWeightTicket, err := FetchWeightTicketByIDExcludeDeletedUploads(appCtx, weightTicket.ID)
	if err != nil {
		return nil, err
	}

	// verify ETag
	if etag.GenerateEtag(originalWeightTicket.UpdatedAt) != eTag {
		return nil, apperror.NewPreconditionFailedError(originalWeightTicket.ID, nil)
	}

	mergedWeightTicket := mergeWeightTicket(weightTicket, *originalWeightTicket)

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

func mergeWeightTicket(weightTicket models.WeightTicket, originalWeightTicket models.WeightTicket) models.WeightTicket {
	mergedWeightTicket := originalWeightTicket

	mergedWeightTicket.VehicleDescription = services.SetOptionalStringField(weightTicket.VehicleDescription, mergedWeightTicket.VehicleDescription)
	mergedWeightTicket.EmptyWeight = services.SetNoNilOptionalPoundField(weightTicket.EmptyWeight, mergedWeightTicket.EmptyWeight)
	mergedWeightTicket.MissingEmptyWeightTicket = services.SetNoNilOptionalBoolField(weightTicket.MissingEmptyWeightTicket, mergedWeightTicket.MissingEmptyWeightTicket)
	mergedWeightTicket.FullWeight = services.SetNoNilOptionalPoundField(weightTicket.FullWeight, mergedWeightTicket.FullWeight)
	mergedWeightTicket.MissingFullWeightTicket = services.SetNoNilOptionalBoolField(weightTicket.MissingFullWeightTicket, mergedWeightTicket.MissingFullWeightTicket)
	mergedWeightTicket.OwnsTrailer = services.SetNoNilOptionalBoolField(weightTicket.OwnsTrailer, mergedWeightTicket.OwnsTrailer)
	mergedWeightTicket.TrailerMeetsCriteria = services.SetNoNilOptionalBoolField(weightTicket.TrailerMeetsCriteria, mergedWeightTicket.TrailerMeetsCriteria)
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

func FetchWeightTicketByIDExcludeDeletedUploads(appContext appcontext.AppContext, weightTicketID uuid.UUID) (*models.WeightTicket, error) {
	var weightTicket models.WeightTicket

	err := appContext.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"EmptyDocument.UserUploads.Upload",
			"FullDocument.UserUploads.Upload",
			"ProofOfTrailerOwnershipDocument.UserUploads.Upload",
		).
		Find(&weightTicket, weightTicketID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(weightTicketID, "while looking for WeightTicket")
		default:
			return nil, apperror.NewQueryError("WeightTicket fetch original", err, "")
		}
	}

	weightTicket.EmptyDocument.UserUploads = weightTicket.EmptyDocument.UserUploads.FilterDeleted()
	weightTicket.FullDocument.UserUploads = weightTicket.FullDocument.UserUploads.FilterDeleted()
	weightTicket.ProofOfTrailerOwnershipDocument.UserUploads = weightTicket.ProofOfTrailerOwnershipDocument.UserUploads.FilterDeleted()

	return &weightTicket, nil
}

func FetchWeightTicketAndPPMShipment(appContext appcontext.AppContext, weightTicketID uuid.UUID) (*models.WeightTicket, error) {
	var weightTicket models.WeightTicket

	err := appContext.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"EmptyDocument.UserUploads.Upload",
			"FullDocument.UserUploads.Upload",
			"ProofOfTrailerOwnershipDocument.UserUploads.Upload",
			"PPMShipment.Shipment",
		).
		Find(&weightTicket, weightTicketID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(weightTicketID, "while looking for WeightTicket")
		default:
			return nil, apperror.NewQueryError("WeightTicket fetch original", err, "")
		}
	}

	weightTicket.EmptyDocument.UserUploads = weightTicket.EmptyDocument.UserUploads.FilterDeleted()
	weightTicket.FullDocument.UserUploads = weightTicket.FullDocument.UserUploads.FilterDeleted()
	weightTicket.ProofOfTrailerOwnershipDocument.UserUploads = weightTicket.ProofOfTrailerOwnershipDocument.UserUploads.FilterDeleted()

	return &weightTicket, nil
}
