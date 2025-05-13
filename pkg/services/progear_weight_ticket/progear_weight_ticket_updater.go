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

// NewCustomerProgearWeightTicketUpdater creates a new progearWeightTicketUpdater struct with the checks it needs for a customer
func NewCustomerProgearWeightTicketUpdater() services.ProgearWeightTicketUpdater {
	return &progearWeightTicketUpdater{
		checks: basicChecksForCustomer(),
	}
}

func NewOfficeProgearWeightTicketUpdater() services.ProgearWeightTicketUpdater {
	return &progearWeightTicketUpdater{
		checks: basicChecksForOffice(),
	}
}

// UpdateProgearWeightTicket updates a progearWeightTicket
func (f *progearWeightTicketUpdater) UpdateProgearWeightTicket(appCtx appcontext.AppContext, progearWeightTicket models.ProgearWeightTicket, eTag string) (*models.ProgearWeightTicket, error) {
	// get existing ProgearWeightTicket
	originalProgearWeightTicket, err := FetchProgearWeightTicketByIDExcludeDeletedUploads(appCtx, progearWeightTicket.ID)
	if err != nil {
		return nil, err
	}

	// verify ETag
	if etag.GenerateEtag(originalProgearWeightTicket.UpdatedAt) != eTag {

		return nil, apperror.NewPreconditionFailedError(originalProgearWeightTicket.ID, nil)
	}

	mergedProgearWeightTicket := mergeProgearWeightTicket(progearWeightTicket, *originalProgearWeightTicket)

	// validate updated model
	if err := validateProgearWeightTicket(appCtx, &mergedProgearWeightTicket, originalProgearWeightTicket, f.checks...); err != nil {
		return nil, err
	}

	if appCtx.Session().IsMilApp() {
		if mergedProgearWeightTicket.Weight != nil {
			mergedProgearWeightTicket.SubmittedWeight = mergedProgearWeightTicket.Weight
		}
		if mergedProgearWeightTicket.BelongsToSelf != nil {
			mergedProgearWeightTicket.SubmittedBelongsToSelf = mergedProgearWeightTicket.BelongsToSelf
		}
		if mergedProgearWeightTicket.HasWeightTickets != nil {
			mergedProgearWeightTicket.SubmittedHasWeightTickets = mergedProgearWeightTicket.HasWeightTickets
		}
	}

	// update the DB record
	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := txnCtx.DB().ValidateAndUpdate(&mergedProgearWeightTicket)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(originalProgearWeightTicket.ID, err, verrs, "invalid input found while updating the ProgearWeightTicket")
		} else if err != nil {
			return apperror.NewQueryError("ProgearWeightTicket update", err, "")
		}

		if err := txnCtx.DB().
			RawQuery("SELECT update_actual_progear_weight_totals($1)", mergedProgearWeightTicket.PPMShipmentID).
			Exec(); err != nil {
			return apperror.NewQueryError("update_actual_progear_weight_totals", err, "")
		}

		return nil
	})

	if txnErr != nil {
		return nil, txnErr
	}

	return &mergedProgearWeightTicket, nil
}

func mergeProgearWeightTicket(progearWeightTicket models.ProgearWeightTicket, originalProgearWeightTicket models.ProgearWeightTicket) models.ProgearWeightTicket {
	mergedProgearWeightTicket := originalProgearWeightTicket

	mergedProgearWeightTicket.Description = services.SetOptionalStringField(progearWeightTicket.Description, mergedProgearWeightTicket.Description)
	mergedProgearWeightTicket.Weight = services.SetNoNilOptionalPoundField(progearWeightTicket.Weight, mergedProgearWeightTicket.Weight)
	mergedProgearWeightTicket.HasWeightTickets = services.SetNoNilOptionalBoolField(progearWeightTicket.HasWeightTickets, mergedProgearWeightTicket.HasWeightTickets)
	mergedProgearWeightTicket.BelongsToSelf = services.SetNoNilOptionalBoolField(progearWeightTicket.BelongsToSelf, mergedProgearWeightTicket.BelongsToSelf)
	mergedProgearWeightTicket.Reason = services.SetOptionalStringField(progearWeightTicket.Reason, mergedProgearWeightTicket.Reason)
	status := services.SetOptionalStringField((*string)(progearWeightTicket.Status), (*string)(mergedProgearWeightTicket.Status))
	if status != nil {
		ppmDocStatus := models.PPMDocumentStatus(*status)
		mergedProgearWeightTicket.Status = &ppmDocStatus
	} else {
		mergedProgearWeightTicket.Status = nil
	}

	return mergedProgearWeightTicket
}

func FetchProgearWeightTicketByIDExcludeDeletedUploads(appContext appcontext.AppContext, progearWeightTicketID uuid.UUID) (*models.ProgearWeightTicket, error) {
	var progearWeightTicket models.ProgearWeightTicket
	findProgearWeightTicketQuery := appContext.DB().Q().Scope(utilities.ExcludeDeletedScope(models.ProgearWeightTicket{})).EagerPreload(
		"Document.UserUploads.Upload",
	)
	if appContext.Session().IsMilApp() {
		findProgearWeightTicketQuery.
			LeftJoin("documents", "documents.id = progear_weight_tickets.document_id").
			Where("documents.service_member_id = ?", appContext.Session().ServiceMemberID)
	}
	err := findProgearWeightTicketQuery.Find(&progearWeightTicket, progearWeightTicketID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(progearWeightTicketID, "while looking for ProgearWeightTicket")
		default:
			return nil, apperror.NewQueryError("ProgearWeightTicket fetch original", err, "")
		}
	}

	progearWeightTicket.Document.UserUploads.FilterDeleted()

	return &progearWeightTicket, nil
}
