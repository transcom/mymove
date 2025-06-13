package gunsafeweightticket

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

type gunSafeWeightTicketUpdater struct {
	checks []gunSafeWeightTicketValidator
}

// NewCustomerGunSafeWeightTicketUpdater creates a new gunSafeWeightTicketUpdater struct with the checks it needs for a customer
func NewCustomerGunSafeWeightTicketUpdater() services.GunSafeWeightTicketUpdater {
	return &gunSafeWeightTicketUpdater{
		checks: basicChecksForCustomer(),
	}
}

func NewOfficeGunSafeWeightTicketUpdater() services.GunSafeWeightTicketUpdater {
	return &gunSafeWeightTicketUpdater{
		checks: basicChecksForOffice(),
	}
}

// UpdateGunSafeWeightTicket updates a gunSafeWeightTicket
func (f *gunSafeWeightTicketUpdater) UpdateGunSafeWeightTicket(appCtx appcontext.AppContext, gunSafeWeightTicket models.GunSafeWeightTicket, eTag string) (*models.GunSafeWeightTicket, error) {
	// get existing GunSafeWeightTicket
	originalGunSafeWeightTicket, err := FetchGunSafeWeightTicketByIDExcludeDeletedUploads(appCtx, gunSafeWeightTicket.ID)
	if err != nil {
		return nil, err
	}

	// verify ETag
	if etag.GenerateEtag(originalGunSafeWeightTicket.UpdatedAt) != eTag {

		return nil, apperror.NewPreconditionFailedError(originalGunSafeWeightTicket.ID, nil)
	}

	mergedGunSafeWeightTicket := mergeGunSafeWeightTicket(gunSafeWeightTicket, *originalGunSafeWeightTicket)

	// validate updated model
	if err := validateGunSafeWeightTicket(appCtx, &mergedGunSafeWeightTicket, originalGunSafeWeightTicket, f.checks...); err != nil {
		return nil, err
	}

	if appCtx.Session().IsMilApp() {
		if mergedGunSafeWeightTicket.Weight != nil {
			mergedGunSafeWeightTicket.SubmittedWeight = mergedGunSafeWeightTicket.Weight
		}
		if mergedGunSafeWeightTicket.HasWeightTickets != nil {
			mergedGunSafeWeightTicket.SubmittedHasWeightTickets = mergedGunSafeWeightTicket.HasWeightTickets
		}
	}

	// update the DB record
	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := txnCtx.DB().ValidateAndUpdate(&mergedGunSafeWeightTicket)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(originalGunSafeWeightTicket.ID, err, verrs, "invalid input found while updating the GunSafeWeightTicket")
		} else if err != nil {
			return apperror.NewQueryError("GunSafeWeightTicket update", err, "")
		}

		return nil
	})

	if txnErr != nil {
		return nil, txnErr
	}

	return &mergedGunSafeWeightTicket, nil
}

func mergeGunSafeWeightTicket(gunSafeWeightTicket models.GunSafeWeightTicket, originalGunSafeWeightTicket models.GunSafeWeightTicket) models.GunSafeWeightTicket {
	mergedGunSafeWeightTicket := originalGunSafeWeightTicket

	mergedGunSafeWeightTicket.Description = services.SetOptionalStringField(gunSafeWeightTicket.Description, mergedGunSafeWeightTicket.Description)
	mergedGunSafeWeightTicket.Weight = services.SetNoNilOptionalPoundField(gunSafeWeightTicket.Weight, mergedGunSafeWeightTicket.Weight)
	mergedGunSafeWeightTicket.HasWeightTickets = services.SetNoNilOptionalBoolField(gunSafeWeightTicket.HasWeightTickets, mergedGunSafeWeightTicket.HasWeightTickets)
	mergedGunSafeWeightTicket.Reason = services.SetOptionalStringField(gunSafeWeightTicket.Reason, mergedGunSafeWeightTicket.Reason)
	status := services.SetOptionalStringField((*string)(gunSafeWeightTicket.Status), (*string)(mergedGunSafeWeightTicket.Status))
	if status != nil {
		ppmDocStatus := models.PPMDocumentStatus(*status)
		mergedGunSafeWeightTicket.Status = &ppmDocStatus
	} else {
		mergedGunSafeWeightTicket.Status = nil
	}

	return mergedGunSafeWeightTicket
}

func FetchGunSafeWeightTicketByIDExcludeDeletedUploads(appContext appcontext.AppContext, gunSafeWeightTicketID uuid.UUID) (*models.GunSafeWeightTicket, error) {
	var gunSafeWeightTicket models.GunSafeWeightTicket
	findGunSafeWeightTicketQuery := appContext.DB().Q().Scope(utilities.ExcludeDeletedScope(models.GunSafeWeightTicket{})).EagerPreload(
		"Document.UserUploads.Upload",
	)
	if appContext.Session().IsMilApp() {
		findGunSafeWeightTicketQuery.
			LeftJoin("documents", "documents.id = gunsafe_weight_tickets.document_id").
			Where("documents.service_member_id = ?", appContext.Session().ServiceMemberID)
	}
	err := findGunSafeWeightTicketQuery.Find(&gunSafeWeightTicket, gunSafeWeightTicketID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(gunSafeWeightTicketID, "while looking for GunSafeWeightTicket")
		default:
			return nil, apperror.NewQueryError("GunSafeWeightTicket fetch original", err, "")
		}
	}

	gunSafeWeightTicket.Document.UserUploads.FilterDeleted()

	return &gunSafeWeightTicket, nil
}
