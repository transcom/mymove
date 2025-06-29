package movetaskorder

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/entitlements"
	"github.com/transcom/mymove/pkg/services/order"
	"github.com/transcom/mymove/pkg/services/query"
)

type moveTaskOrderUpdater struct {
	moveTaskOrderFetcher
	builder                    UpdateMoveTaskOrderQueryBuilder
	serviceItemCreator         services.MTOServiceItemCreator
	moveRouter                 services.MoveRouter
	signedCertificationCreator services.SignedCertificationCreator
	signedCertificationUpdater services.SignedCertificationUpdater
	estimator                  services.PPMEstimator
}

// NewMoveTaskOrderUpdater creates a new struct with the service dependencies
func NewMoveTaskOrderUpdater(builder UpdateMoveTaskOrderQueryBuilder, serviceItemCreator services.MTOServiceItemCreator, moveRouter services.MoveRouter, signedCertificationCreator services.SignedCertificationCreator, signedCertificationUpdater services.SignedCertificationUpdater, estimator services.PPMEstimator) services.MoveTaskOrderUpdater {
	// Fetcher dependency
	waf := entitlements.NewWeightAllotmentFetcher()

	return &moveTaskOrderUpdater{moveTaskOrderFetcher{
		waf: waf,
	}, builder, serviceItemCreator, moveRouter, signedCertificationCreator, signedCertificationUpdater, estimator}
}

// UpdateStatusServiceCounselingCompleted updates the status on the move (move task order) to service counseling completed
func (o moveTaskOrderUpdater) UpdateStatusServiceCounselingCompleted(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string) (*models.Move, error) {
	// Fetch the move and associations.
	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	move, fetchErr := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if fetchErr != nil {
		return &models.Move{}, fetchErr
	}

	// Check the If-Match header against existing eTag before updating.
	encodedUpdatedAt := etag.GenerateEtag(move.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return &models.Move{}, apperror.NewPreconditionFailedError(move.ID, nil)
	}

	transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
		// Update move status, verifying that move/shipments are in expected state.
		err := o.moveRouter.CompleteServiceCounseling(appCtx, move)
		if err != nil {
			return err
		}

		//When submiting a move for approval - remove the SC assigned user
		move.SCCounselingAssignedID = nil

		// Save the move.
		var verrs *validate.Errors
		verrs, err = appCtx.DB().ValidateAndSave(move)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(move.ID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		// If this is move has a PPM, then we also need to adjust other statuses:
		//   - set MTO shipment status to APPROVED
		//   - set PPM shipment status to WAITING_ON_CUSTOMER
		// TODO: Perhaps this could be part of the shipment router. PPMs are a separate model/table,
		//   so would need to figure out how they factor in.
		if move.HasPPM() {
			// Note: Avoiding the copy of the element in the range so we can preserve the changes to the
			// statuses when we return the entire move tree.

			for i := range move.MTOShipments { // We should only change for PPM shipments.
				if move.MTOShipments[i].PPMShipment != nil {
					move.MTOShipments[i].Status = models.MTOShipmentStatusApproved

					verrs, err = appCtx.DB().ValidateAndSave(&move.MTOShipments[i])
					if verrs != nil && verrs.HasAny() {
						return apperror.NewInvalidInputError(move.MTOShipments[i].ID, nil, verrs, "")
					}
					if err != nil {
						return err
					}

					var ppm models.PPMShipment
					err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
						EagerPreload(
							"Shipment",
							"WeightTickets",
							"MovingExpenses",
							"ProgearWeightTickets",
							"W2Address.Country",
							"PickupAddress.Country",
							"SecondaryPickupAddress.Country",
							"TertiaryPickupAddress.Country",
							"DestinationAddress.Country",
							"SecondaryDestinationAddress.Country",
							"TertiaryDestinationAddress.Country",
						).
						Where("shipment_id = ?", move.MTOShipments[i].ID).First(&ppm)

					if err != nil {
						switch err {
						case sql.ErrNoRows:
							return apperror.NewNotFoundError(move.MTOShipments[i].ID, "while looking for PPMShipment by MTO ShipmentID")
						default:
							return apperror.NewQueryError("PPMShipment", err, "unable to find PPMShipment")
						}
					}
					// if the customer has their max incentive empty in the db, we need to update this value before proceeding
					if ppm.MaxIncentive == nil {
						estimatedIncentive, estimatedSITCost, err := o.estimator.EstimateIncentiveWithDefaultChecks(appCtx, ppm, &ppm)
						if err != nil {
							return err
						}
						move.MTOShipments[i].PPMShipment.EstimatedIncentive = estimatedIncentive
						move.MTOShipments[i].PPMShipment.SITEstimatedCost = estimatedSITCost

						maxIncentive, err := o.estimator.MaxIncentive(appCtx, ppm, &ppm)
						if err != nil {
							return err
						}
						move.MTOShipments[i].PPMShipment.MaxIncentive = maxIncentive
					}

					move.MTOShipments[i].PPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer
					now := time.Now()
					move.MTOShipments[i].PPMShipment.ApprovedAt = &now

					verrs, err = appCtx.DB().ValidateAndSave(move.MTOShipments[i].PPMShipment)
					if verrs != nil && verrs.HasAny() {
						return apperror.NewInvalidInputError(move.MTOShipments[i].PPMShipment.ID, nil, verrs, "")
					}
					if err != nil {
						return err
					}

					err = o.SignCertificationPPMCounselingCompleted(appCtx, move.ID, move.MTOShipments[i].PPMShipment.ID)
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	if transactionError != nil {
		return &models.Move{}, transactionError
	}

	return move, nil
}

// UpdateReviewedBillableWeightsAt updates the BillableWeightsReviewedAt field on the move (move task order)
func (o moveTaskOrderUpdater) UpdateReviewedBillableWeightsAt(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string) (*models.Move, error) {
	var err error

	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	move, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return &models.Move{}, err
	}

	transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
		// update field for move
		now := time.Now()
		if move.BillableWeightsReviewedAt == nil {
			move.BillableWeightsReviewedAt = &now
		}

		// Check the If-Match header against existing eTag before updating
		encodedUpdatedAt := etag.GenerateEtag(move.UpdatedAt)
		if encodedUpdatedAt != eTag {
			return apperror.NewPreconditionFailedError(move.ID, err)
		}

		err = appCtx.DB().Update(move)
		return err
	})
	if transactionError != nil {
		return &models.Move{}, transactionError
	}

	return move, nil
}

// UpdateTIORemarks updates the TIORemarks field on the move (move task order)
func (o moveTaskOrderUpdater) UpdateTIORemarks(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string, remarks string) (*models.Move, error) {
	var err error

	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	move, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return &models.Move{}, err
	}

	transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
		// update field for move
		move.TIORemarks = &remarks

		// Check the If-Match header against existing eTag before updating
		encodedUpdatedAt := etag.GenerateEtag(move.UpdatedAt)
		if encodedUpdatedAt != eTag {
			return apperror.NewPreconditionFailedError(move.ID, err)
		}

		err = appCtx.DB().Update(move)
		return err
	})
	if transactionError != nil {
		return &models.Move{}, transactionError
	}

	return move, nil
}

// ApproveMoveAndCreateServiceItems approves a Move and
// creates Move-level service items (counseling and move management) if the
// TOO selected them. If the move received service counseling, the counseling
// service item will automatically be created without the TOO having to select it.
func (o *moveTaskOrderUpdater) ApproveMoveAndCreateServiceItems(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string,
	includeServiceCodeMS bool, includeServiceCodeCS bool) (*models.Move, error) {

	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	move, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return &models.Move{}, err
	}

	existingETag := etag.GenerateEtag(move.UpdatedAt)
	if existingETag != eTag {
		return &models.Move{}, apperror.NewPreconditionFailedError(move.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	//When approving a shipment - remove the assigned TOO user
	move.TOOTaskOrderAssignedID = nil

	updateMove := false
	if move.ApprovedAt == nil {
		updateMove = true
		err = o.moveRouter.Approve(appCtx, move)
		if err != nil {
			return &models.Move{}, apperror.NewConflictError(move.ID, err.Error())
		}
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if updateMove {
			err = o.updateMove(txnAppCtx, move, order.CheckRequiredFields())
			if err != nil {
				return err
			}
		}

		// When provided, this will create and approve these Move-level service items.
		if includeServiceCodeMS && !move.IsPPMOnly() {
			err = o.createMoveLevelServiceItem(txnAppCtx, *move, models.ReServiceCodeMS)
		}

		if err != nil {
			return err
		}

		if includeServiceCodeCS {
			err = o.createMoveLevelServiceItem(txnAppCtx, *move, models.ReServiceCodeCS)
		}

		return err
	})

	if transactionError != nil {
		return &models.Move{}, transactionError
	}

	return move, nil
}

// MakeAvailableToPrime makes the move available to prime
func (o *moveTaskOrderUpdater) MakeAvailableToPrime(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID) (*models.Move, bool, error) {
	var move *models.Move
	var wasMadeAvailableToPrime = false

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		searchParams := services.MoveTaskOrderFetcherParams{
			IncludeHidden:   false,
			MoveTaskOrderID: moveTaskOrderID,
		}
		var err error
		move, err = o.FetchMoveTaskOrder(txnAppCtx, &searchParams)
		if err != nil {
			return err
		}

		if move.AvailableToPrimeAt == nil {
			now := time.Now()
			move.AvailableToPrimeAt = &now

			err = o.updateMove(txnAppCtx, move, order.CheckRequiredFields())
			if err != nil {
				return err
			}
			wasMadeAvailableToPrime = true
		}
		return nil
	})

	if transactionError != nil {
		return &models.Move{}, false, transactionError
	}

	return move, wasMadeAvailableToPrime, nil
}

func (o *moveTaskOrderUpdater) updateMove(appCtx appcontext.AppContext, move *models.Move, checks ...order.Validator) error {
	if verr := order.ValidateOrder(&move.Orders, checks...); verr != nil {
		return verr
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(move)

	if verrs != nil && verrs.HasAny() {
		return apperror.NewInvalidInputError(move.ID, nil, verrs, "")
	}

	return err
}

func (o *moveTaskOrderUpdater) createMoveLevelServiceItem(appCtx appcontext.AppContext, move models.Move, code models.ReServiceCode) error {
	now := time.Now()

	siCreator := o.serviceItemCreator

	_, verrs, err := siCreator.CreateMTOServiceItem(appCtx, &models.MTOServiceItem{
		MoveTaskOrderID: move.ID,
		MTOShipmentID:   nil,
		ReService:       models.ReService{Code: code},
		Status:          models.MTOServiceItemStatusApproved,
		ApprovedAt:      &now,
	})

	if err != nil {
		if errors.Is(err, models.ErrInvalidTransition) {
			return apperror.NewConflictError(move.ID, err.Error())
		}
		return err
	}

	if verrs != nil && verrs.HasAny() {
		return apperror.NewInvalidInputError(move.ID, nil, verrs, "")
	}

	return nil
}

// UpdateMoveTaskOrderQueryBuilder is the query builder for updating MTO
type UpdateMoveTaskOrderQueryBuilder interface {
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
}

// UpdatePostCounselingInfo updates the counseling info
func (o *moveTaskOrderUpdater) UpdatePostCounselingInfo(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID, eTag string) (*models.Move, error) {
	// Fetch the move and associations.
	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:            false,
		MoveTaskOrderID:          moveTaskOrderID,
		ExcludeExternalShipments: true,
	}
	moveTaskOrder, fetchErr := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if fetchErr != nil {
		return &models.Move{}, fetchErr
	}

	approvedForPrimeCounseling := false
	for _, serviceItem := range moveTaskOrder.MTOServiceItems {
		if serviceItem.ReService.Code == models.ReServiceCodeCS && serviceItem.Status == models.MTOServiceItemStatusApproved {
			approvedForPrimeCounseling = true
			break
		}
	}
	if !approvedForPrimeCounseling {
		return &models.Move{}, apperror.NewConflictError(moveTaskOrderID, "Counseling is not an approved service item")
	}

	transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
		// Check the If-Match header against existing eTag before updating.
		encodedUpdatedAt := etag.GenerateEtag(moveTaskOrder.UpdatedAt)
		if encodedUpdatedAt != eTag {
			return apperror.NewPreconditionFailedError(moveTaskOrderID, nil)
		}

		now := time.Now()
		moveTaskOrder.PrimeCounselingCompletedAt = &now

		verrs, err := appCtx.DB().ValidateAndSave(moveTaskOrder)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(moveTaskOrderID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		// Note: Avoiding the copy of the element in the range so we can preserve the changes to the
		// statuses when we return the entire move tree.
		for i := range moveTaskOrder.MTOShipments {
			if moveTaskOrder.MTOShipments[i].PPMShipment != nil {
				moveTaskOrder.MTOShipments[i].PPMShipment.Status = models.PPMShipmentStatusWaitingOnCustomer
				moveTaskOrder.MTOShipments[i].PPMShipment.ApprovedAt = &now

				verrs, err = appCtx.DB().ValidateAndSave(moveTaskOrder.MTOShipments[i].PPMShipment)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(moveTaskOrder.MTOShipments[i].PPMShipment.ID, nil, verrs, "")
				}
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if transactionError != nil {
		return &models.Move{}, transactionError
	}

	return moveTaskOrder, nil
}

// ShowHide changes the value in the "Show" field for a Move. This can be either True or False and indicates if the move has been deactivated or not.
func (o *moveTaskOrderUpdater) ShowHide(appCtx appcontext.AppContext, moveID uuid.UUID, show *bool) (*models.Move, error) {
	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   true, // We need to search every move to change its status
		MoveTaskOrderID: moveID,
	}
	move, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return nil, err
	}

	if show == nil {
		return nil, apperror.NewInvalidInputError(moveID, nil, nil, "The 'show' field must be either True or False - it cannot be empty")
	}

	move.Show = show
	verrs, err := appCtx.DB().ValidateAndSave(move)
	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(move.ID, err, verrs, "Invalid input found while updating the Move")
	} else if err != nil {
		return nil, apperror.NewQueryError("Move", err, "")
	}

	// Get the updated Move and return
	updatedMove, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return nil, apperror.NewQueryError("Move", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}

	return updatedMove, nil
}

// UpdatePPMType updates the PPMType field on the move (move task order)
func (o moveTaskOrderUpdater) UpdatePPMType(appCtx appcontext.AppContext, moveTaskOrderID uuid.UUID) (*models.Move, error) {
	searchParams := services.MoveTaskOrderFetcherParams{
		IncludeHidden:   false,
		MoveTaskOrderID: moveTaskOrderID,
	}
	move, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return nil, err
	}

	transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
		if move.IsPPMOnly() { // Only PPM Shipments in the move
			ppmType := models.MovePPMTypeFULL
			move.PPMType = &ppmType
		} else if move.HasPPM() { // At least 1 PPM Shipment in the move
			ppmType := models.MovePPMTypePARTIAL
			move.PPMType = &ppmType
		} else {
			move.PPMType = nil
		}
		// update PPMType Column for move in DB
		err = appCtx.DB().UpdateColumns(move, "ppm_type")
		if err != nil {
			return err
		}

		return err
	})
	if transactionError != nil {
		return move, transactionError
	}

	// Get the updated Move and return
	updatedMove, err := o.FetchMoveTaskOrder(appCtx, &searchParams)
	if err != nil {
		return nil, apperror.NewQueryError("Move", err, fmt.Sprintf("Unexpected error after saving: %v", err))
	}

	return updatedMove, nil
}

func (o moveTaskOrderUpdater) SignCertificationPPMCounselingCompleted(appCtx appcontext.AppContext, moveID uuid.UUID, ppmShipmentID uuid.UUID) error {
	// Retrieve if PPM has certificate
	signedCertifications, err := models.FetchSignedCertificationPPMByType(appCtx.DB(), appCtx.Session(), moveID, ppmShipmentID, models.SignedCertificationTypePreCloseoutReviewedPPMPAYMENT)
	if err != nil {
		return err
	}

	signatureText := fmt.Sprintf("%s %s", appCtx.Session().FirstName, appCtx.Session().LastName)

	if len(signedCertifications) == 0 {
		// Add new certificate
		now := time.Now()
		certificateType := models.SignedCertificationTypePreCloseoutReviewedPPMPAYMENT
		signedCertification := models.SignedCertification{
			SubmittingUserID:  appCtx.Session().UserID,
			MoveID:            moveID,
			PpmID:             models.UUIDPointer(ppmShipmentID),
			CertificationType: &certificateType,
			CertificationText: "Confirmed: Reviewed Waiting On Customer ",
			Signature:         signatureText,
			Date:              now,
		}
		_, err := o.signedCertificationCreator.CreateSignedCertification(appCtx, signedCertification)
		if err != nil {
			return err
		}
	} else {
		// Update existing certificate. Ensure only one
		eTag := etag.GenerateEtag(signedCertifications[0].UpdatedAt)
		// Update with current counselor information
		signedCertifications[0].SubmittingUserID = appCtx.Session().UserID
		signedCertifications[0].Signature = signatureText
		_, err := o.signedCertificationUpdater.UpdateSignedCertification(appCtx, *signedCertifications[0], eTag)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateStatusServiceCounselingSendPPMToCustomer updates the status on a PPM Shipment and creates required certs
func (o moveTaskOrderUpdater) UpdateStatusServiceCounselingSendPPMToCustomer(appCtx appcontext.AppContext, ppmShipment models.PPMShipment, eTag string, move *models.Move) (*models.PPMShipment, error) {
	// Check the If-Match header against existing eTag before updating.
	encodedUpdatedAt := etag.GenerateEtag(ppmShipment.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return &models.PPMShipment{}, apperror.NewPreconditionFailedError(move.ID, nil)
	}

	transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
		if ppmShipment.Shipment.Status != models.MTOShipmentStatusApproved {
			ppmShipment.Shipment.Status = models.MTOShipmentStatusApproved
			verrs, err := appCtx.DB().ValidateAndSave(&ppmShipment.Shipment)
			if verrs != nil && verrs.HasAny() {
				return apperror.NewInvalidInputError(ppmShipment.Shipment.ID, nil, verrs, "")
			}
			if err != nil {
				return err
			}
		}

		// Pull old ppmshipment for estimator update
		var ppm models.PPMShipment
		err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
			EagerPreload(
				"Shipment",
				"WeightTickets",
				"MovingExpenses",
				"ProgearWeightTickets",
				"W2Address.Country",
				"PickupAddress.Country",
				"SecondaryPickupAddress.Country",
				"TertiaryPickupAddress.Country",
				"DestinationAddress.Country",
				"SecondaryDestinationAddress.Country",
				"TertiaryDestinationAddress.Country",
			).
			Where("shipment_id = ?", ppmShipment.ShipmentID).First(&ppm)

		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return apperror.NewNotFoundError(ppmShipment.ID, "while looking for PPMShipment by MTO ShipmentID")
			default:
				return apperror.NewQueryError("PPMShipment", err, "unable to find PPMShipment")
			}
		}
		// if the customer has their max incentive empty in the db, we need to update this value before proceeding
		if ppmShipment.MaxIncentive == nil {
			estimatedIncentive, estimatedSITCost, err := o.estimator.EstimateIncentiveWithDefaultChecks(appCtx, ppm, &ppmShipment)
			if err != nil {
				return err
			}
			ppmShipment.EstimatedIncentive = estimatedIncentive
			ppmShipment.SITEstimatedCost = estimatedSITCost

			maxIncentive, err := o.estimator.MaxIncentive(appCtx, ppm, &ppmShipment)
			if err != nil {
				return err
			}
			ppmShipment.MaxIncentive = maxIncentive
		}

		ppmShipment.Status = models.PPMShipmentStatusWaitingOnCustomer
		now := time.Now()
		ppmShipment.ApprovedAt = &now

		verrs, err := appCtx.DB().ValidateAndSave(&ppmShipment)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(ppmShipment.ID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		err = o.SignCertificationPPMCounselingCompleted(appCtx, move.ID, ppmShipment.ID)
		if err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return &models.PPMShipment{}, transactionError
	}

	return &ppmShipment, nil
}
