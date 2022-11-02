package progearweightticket

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

type progearWeightTicketUpdater struct {
	checks []progearWeightTicketValidator
}

func NewProgearWeightTicketUpdater() services.ProgearWeightTicketUpdater {
	return &progearWeightTicketUpdater{
		checks: updateChecks(),
	}
}

func (f *progearWeightTicketUpdater) UpdateProgearWeightTicket(appCtx appcontext.AppContext, progearWeightTicket models.ProgearWeightTicket, eTag string) (*models.ProgearWeightTicket, error) {
	oldProgearWeightTicket, err := FetchProgearID(appCtx, progearWeightTicket.ID)

	if err != nil {
		return nil, err
	}

	if etag.GenerateEtag(oldProgearWeightTicket.UpdatedAt) != eTag {
		return nil, apperror.NewPreconditionFailedError(oldProgearWeightTicket.ID, nil)
	}

	mergedProgearWeightTicket := mergeProgearWeightTicket(progearWeightTicket, *oldProgearWeightTicket)

	err = validateProgearWeightTicket(appCtx, &mergedProgearWeightTicket, oldProgearWeightTicket, f.checks...)

	if err != nil {
		return nil, err
	}

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := txnCtx.DB().Eager().ValidateAndUpdate(&mergedProgearWeightTicket)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(oldProgearWeightTicket.ID, err, verrs, "")
		} else if err != nil {
			return apperror.NewQueryError("Progear Weight Ticket", err, "")
		}

		return nil
	})

	if txnErr != nil {
		return nil, txnErr
	}

	return &mergedProgearWeightTicket, nil
}

func FetchProgearID(appContext appcontext.AppContext, progearID uuid.UUID) (*models.ProgearWeightTicket, error) {
	var progear models.ProgearWeightTicket

	err := appContext.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload("FullDocument.UserUploads.Upload").Find(&progear, progearID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(progearID, "while looking for Progear")
		default:
			return nil, apperror.NewQueryError("Progear fetch original", err, "")
		}
	}
	// set the updated Document
	progear.Document.UserUploads = progear.Document.UserUploads.FilterDeleted()

	return &progear, nil
}

func mergeProgearWeightTicket(updatedProgearWeightTicket models.ProgearWeightTicket, oldProgearWeightTicket models.ProgearWeightTicket) models.ProgearWeightTicket {
	mergedProgearWeightTicket := oldProgearWeightTicket

	progearWeightTicketStatus := services.SetOptionalStringField((*string)(updatedProgearWeightTicket.Status), (*string)(mergedProgearWeightTicket.Status))
	if progearWeightTicketStatus != nil {
		ppmDocumentStatus := models.PPMDocumentStatus(*progearWeightTicketStatus)
		mergedProgearWeightTicket.Status = &ppmDocumentStatus

		if ppmDocumentStatus == models.PPMDocumentStatusExcluded || ppmDocumentStatus == models.PPMDocumentStatusRejected {
			mergedProgearWeightTicket.Reason = services.SetOptionalStringField(updatedProgearWeightTicket.Reason, mergedProgearWeightTicket.Reason)
		} else {
			// if that status is changed back to approved then we should clear the reason value
			mergedProgearWeightTicket.Reason = nil
		}
	} else {
		mergedProgearWeightTicket.Status = nil
	}

	mergedProgearWeightTicket.BelongsToSelf = services.SetNoNilOptionalBoolField(updatedProgearWeightTicket.BelongsToSelf, mergedProgearWeightTicket.BelongsToSelf)
	mergedProgearWeightTicket.Description = services.SetOptionalStringField(updatedProgearWeightTicket.Description, mergedProgearWeightTicket.Description)
	mergedProgearWeightTicket.HasWeightTickets = services.SetNoNilOptionalBoolField(updatedProgearWeightTicket.HasWeightTickets, mergedProgearWeightTicket.HasWeightTickets)
	mergedProgearWeightTicket.Weight = services.SetOptionalPoundField(updatedProgearWeightTicket.Weight, mergedProgearWeightTicket.Weight)
	mergedProgearWeightTicket.DeletedAt = services.SetOptionalDateTimeField(updatedProgearWeightTicket.DeletedAt, mergedProgearWeightTicket.DeletedAt)

	return mergedProgearWeightTicket
}
