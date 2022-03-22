package mtoshipment

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/getlantern/deepcopy"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/query"
)

// UpdateMTOShipmentQueryBuilder is the query builder for updating MTO Shipments
type UpdateMTOShipmentQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
	CreateOne(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error)
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
	Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error)
	FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
}

type mtoShipmentUpdater struct {
	builder UpdateMTOShipmentQueryBuilder
	services.Fetcher
	planner      route.Planner
	moveRouter   services.MoveRouter
	moveWeights  services.MoveWeights
	sender       notifications.NotificationSender
	recalculator services.PaymentRequestShipmentRecalculator
}

// NewMTOShipmentUpdater creates a new struct with the service dependencies
func NewMTOShipmentUpdater(builder UpdateMTOShipmentQueryBuilder, fetcher services.Fetcher, planner route.Planner, moveRouter services.MoveRouter, moveWeights services.MoveWeights, sender notifications.NotificationSender, recalculator services.PaymentRequestShipmentRecalculator) services.MTOShipmentUpdater {
	return &mtoShipmentUpdater{
		builder,
		fetch.NewFetcher(builder),
		planner,
		moveRouter,
		moveWeights,
		sender,
		recalculator,
	}
}

// setNewShipmentFields validates the updated shipment
func setNewShipmentFields(appCtx appcontext.AppContext, dbShipment *models.MTOShipment, requestedUpdatedShipment *models.MTOShipment) {
	if requestedUpdatedShipment.RequestedPickupDate != nil {
		dbShipment.RequestedPickupDate = requestedUpdatedShipment.RequestedPickupDate
	}

	dbShipment.Diversion = requestedUpdatedShipment.Diversion
	dbShipment.UsesExternalVendor = requestedUpdatedShipment.UsesExternalVendor

	if requestedUpdatedShipment.RequestedDeliveryDate != nil {
		dbShipment.RequestedDeliveryDate = requestedUpdatedShipment.RequestedDeliveryDate
	}

	if requestedUpdatedShipment.PrimeActualWeight != nil {
		dbShipment.PrimeActualWeight = requestedUpdatedShipment.PrimeActualWeight
	}

	if requestedUpdatedShipment.FirstAvailableDeliveryDate != nil {
		dbShipment.FirstAvailableDeliveryDate = requestedUpdatedShipment.FirstAvailableDeliveryDate
	}

	if requestedUpdatedShipment.ActualPickupDate != nil {
		dbShipment.ActualPickupDate = requestedUpdatedShipment.ActualPickupDate
	}

	if requestedUpdatedShipment.ScheduledPickupDate != nil {
		dbShipment.ScheduledPickupDate = requestedUpdatedShipment.ScheduledPickupDate
	}

	if requestedUpdatedShipment.PrimeEstimatedWeight != nil {
		now := time.Now()
		dbShipment.PrimeEstimatedWeight = requestedUpdatedShipment.PrimeEstimatedWeight
		dbShipment.PrimeEstimatedWeightRecordedDate = &now
	}

	if requestedUpdatedShipment.NTSRecordedWeight != nil {
		dbShipment.NTSRecordedWeight = requestedUpdatedShipment.NTSRecordedWeight
	}

	if requestedUpdatedShipment.PickupAddress != nil {
		dbShipment.PickupAddress = requestedUpdatedShipment.PickupAddress
	}

	if requestedUpdatedShipment.DestinationAddress != nil {
		dbShipment.DestinationAddress = requestedUpdatedShipment.DestinationAddress
	}

	if requestedUpdatedShipment.DestinationType != nil {
		dbShipment.DestinationType = requestedUpdatedShipment.DestinationType
	}

	if requestedUpdatedShipment.SecondaryPickupAddress != nil {
		dbShipment.SecondaryPickupAddress = requestedUpdatedShipment.SecondaryPickupAddress
	}

	if requestedUpdatedShipment.SecondaryDeliveryAddress != nil {
		dbShipment.SecondaryDeliveryAddress = requestedUpdatedShipment.SecondaryDeliveryAddress
	}

	if requestedUpdatedShipment.ShipmentType != "" {
		dbShipment.ShipmentType = requestedUpdatedShipment.ShipmentType
	}

	if requestedUpdatedShipment.Status != "" {
		dbShipment.Status = requestedUpdatedShipment.Status
	}

	if requestedUpdatedShipment.RequiredDeliveryDate != nil {
		dbShipment.RequiredDeliveryDate = requestedUpdatedShipment.RequiredDeliveryDate
	}

	if requestedUpdatedShipment.PrimeEstimatedWeightRecordedDate != nil {
		dbShipment.PrimeEstimatedWeightRecordedDate = requestedUpdatedShipment.PrimeEstimatedWeightRecordedDate
	}

	if requestedUpdatedShipment.CustomerRemarks != nil {
		dbShipment.CustomerRemarks = requestedUpdatedShipment.CustomerRemarks
	}

	if requestedUpdatedShipment.CounselorRemarks != nil {
		dbShipment.CounselorRemarks = requestedUpdatedShipment.CounselorRemarks
	}

	if requestedUpdatedShipment.BillableWeightCap != nil {
		dbShipment.BillableWeightCap = requestedUpdatedShipment.BillableWeightCap
	}

	if requestedUpdatedShipment.BillableWeightJustification != nil {
		dbShipment.BillableWeightJustification = requestedUpdatedShipment.BillableWeightJustification
	}

	if requestedUpdatedShipment.TACType != nil && *requestedUpdatedShipment.TACType == "" {
		dbShipment.TACType = requestedUpdatedShipment.TACType
	} else if requestedUpdatedShipment.TACType != nil {
		dbShipment.TACType = requestedUpdatedShipment.TACType
	}

	if requestedUpdatedShipment.SACType != nil && *requestedUpdatedShipment.SACType == "" {
		dbShipment.SACType = requestedUpdatedShipment.SACType
	} else if requestedUpdatedShipment.SACType != nil {
		dbShipment.SACType = requestedUpdatedShipment.SACType
	}

	if requestedUpdatedShipment.ServiceOrderNumber != nil {
		dbShipment.ServiceOrderNumber = requestedUpdatedShipment.ServiceOrderNumber
	}

	if requestedUpdatedShipment.StorageFacility != nil {
		dbShipment.StorageFacility = requestedUpdatedShipment.StorageFacility
	}

	//// TODO: move mtoagent creation into service: Should not update MTOAgents here because we don't have an eTag
	if len(requestedUpdatedShipment.MTOAgents) > 0 {
		agentsToCreateOrUpdate := []models.MTOAgent{}
		for _, newAgentInfo := range requestedUpdatedShipment.MTOAgents {
			// if no record exists in the db
			if newAgentInfo.ID == uuid.Nil {
				newAgentInfo.MTOShipmentID = requestedUpdatedShipment.ID
				agentsToCreateOrUpdate = append(agentsToCreateOrUpdate, newAgentInfo)
			} else {
				foundAgent := false
				// make sure there is an existing record in the db
				for i, dbAgent := range dbShipment.MTOAgents {
					if foundAgent {
						break
					}
					if dbAgent.ID == newAgentInfo.ID {
						foundAgent = true
						if newAgentInfo.MTOAgentType != "" && newAgentInfo.MTOAgentType != dbAgent.MTOAgentType {
							dbShipment.MTOAgents[i].MTOAgentType = newAgentInfo.MTOAgentType
						}

						if newAgentInfo.FirstName != nil {
							dbShipment.MTOAgents[i].FirstName = newAgentInfo.FirstName
						}

						if newAgentInfo.LastName != nil {
							dbShipment.MTOAgents[i].LastName = newAgentInfo.LastName
						}

						if newAgentInfo.Email != nil {
							dbShipment.MTOAgents[i].Email = newAgentInfo.Email
						}

						if newAgentInfo.Phone != nil {
							dbShipment.MTOAgents[i].Phone = newAgentInfo.Phone
						}
						agentsToCreateOrUpdate = append(agentsToCreateOrUpdate, dbShipment.MTOAgents[i])
					}
				}
			}
		}
		dbShipment.MTOAgents = agentsToCreateOrUpdate // don't return unchanged existing agents
	}
}

// StaleIdentifierError is used when optimistic locking determines that the identifier refers to stale data
type StaleIdentifierError struct {
	StaleIdentifier string
}

func (e StaleIdentifierError) Error() string {
	return fmt.Sprintf("stale identifier: %s", e.StaleIdentifier)
}

//CheckIfMTOShipmentCanBeUpdated checks if a shipment should be updatable
func (f *mtoShipmentUpdater) CheckIfMTOShipmentCanBeUpdated(appCtx appcontext.AppContext, mtoShipment *models.MTOShipment, session *auth.Session) (bool, error) {
	if session.IsOfficeApp() && session.IsOfficeUser() {
		isServiceCounselor := session.Roles.HasRole(roles.RoleTypeServicesCounselor)
		isTOO := session.Roles.HasRole(roles.RoleTypeTOO)
		isTIO := session.Roles.HasRole(roles.RoleTypeTIO)
		switch mtoShipment.Status {
		case models.MTOShipmentStatusSubmitted:
			if isServiceCounselor || isTOO {
				return true, nil
			}
		case models.MTOShipmentStatusApproved:
			if isTIO || isTOO {
				return true, nil
			}
		case models.MTOShipmentStatusCancellationRequested:
			if isTOO {
				return true, nil
			}
		case models.MTOShipmentStatusCanceled:
			if isTOO {
				return true, nil
			}
		case models.MTOShipmentStatusDiversionRequested:
			if isTOO {
				return true, nil
			}
		default:
			return false, nil
		}

		return false, nil
	}

	return true, nil
}

func (f *mtoShipmentUpdater) RetrieveMTOShipment(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (*models.MTOShipment, error) {
	var shipment models.MTOShipment

	err := appCtx.DB().EagerPreload(
		"MoveTaskOrder",
		"PickupAddress",
		"DestinationAddress",
		"SecondaryPickupAddress",
		"SecondaryDeliveryAddress",
		"MTOAgents",
		"SITExtensions",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOServiceItems.CustomerContacts",
		"StorageFacility.Address").Find(&shipment, mtoShipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(mtoShipmentID, "Shipment not found")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	return &shipment, nil
}

// UpdateMTOShipmentOffice updates the mto shipment
// TODO: apply the subset of business logic validations
// that would be appropriate for the OFFICE USER
func (f *mtoShipmentUpdater) UpdateMTOShipmentOffice(appCtx appcontext.AppContext, mtoShipment *models.MTOShipment, eTag string) (*models.MTOShipment, error) {
	return f.updateMTOShipment(
		appCtx,
		mtoShipment,
		eTag,
		checkStatus(),
	)
}

// UpdateMTOShipmentCustomer updates the mto shipment
// TODO: apply the subset of business logic validations
// that would be appropriate for the CUSTOMER
func (f *mtoShipmentUpdater) UpdateMTOShipmentCustomer(appCtx appcontext.AppContext, mtoShipment *models.MTOShipment, eTag string) (*models.MTOShipment, error) {
	return f.updateMTOShipment(
		appCtx,
		mtoShipment,
		eTag,
		checkStatus(),
	)
}

// UpdateMTOShipmentPrime updates the mto shipment
// TODO: apply the subset of business logic validations
// that would be appropriate for the PRIME
func (f *mtoShipmentUpdater) UpdateMTOShipmentPrime(appCtx appcontext.AppContext, mtoShipment *models.MTOShipment, eTag string) (*models.MTOShipment, error) {
	return f.updateMTOShipment(
		appCtx,
		mtoShipment,
		eTag,
		checkStatus(),
		checkAvailToPrime(),
	)
}

//updateMTOShipment updates the mto shipment
func (f *mtoShipmentUpdater) updateMTOShipment(appCtx appcontext.AppContext, mtoShipment *models.MTOShipment, eTag string, checks ...validator) (*models.MTOShipment, error) {
	oldShipment, err := f.RetrieveMTOShipment(appCtx, mtoShipment.ID)
	if err != nil {
		return nil, err
	}

	// run the (read-only) validations
	if verr := validateShipment(appCtx, mtoShipment, oldShipment, checks...); verr != nil {
		return nil, verr
	}

	var dbShipment models.MTOShipment
	err = deepcopy.Copy(&dbShipment, oldShipment) // save the original db version, oldShipment will be modified
	if err != nil {
		return nil, fmt.Errorf("error copying shipment data %w", err)
	}
	setNewShipmentFields(appCtx, oldShipment, mtoShipment)
	newShipment := oldShipment // old shipment has now been updated with requested changes
	// db version is used to check if agents need creating or updating
	err = f.updateShipmentRecord(appCtx, &dbShipment, newShipment, eTag)
	if err != nil {
		switch err.(type) {
		case StaleIdentifierError:
			return nil, apperror.NewPreconditionFailedError(mtoShipment.ID, err)
		default:
			return nil, err
		}
	}

	updatedShipment, err := f.RetrieveMTOShipment(appCtx, mtoShipment.ID)
	if err != nil {
		return nil, err
	}

	return updatedShipment, nil
}

// Takes the validated shipment input and updates the database using a transaction. If any part of the
// update fails, the entire transaction will be rolled back.
func (f *mtoShipmentUpdater) updateShipmentRecord(appCtx appcontext.AppContext, dbShipment *models.MTOShipment, newShipment *models.MTOShipment, eTag string) error {
	var autoReweighShipments models.MTOShipments
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// temp optimistic locking solution til query builder is re-tooled to handle nested updates
		encodedUpdatedAt := etag.GenerateEtag(newShipment.UpdatedAt)

		if encodedUpdatedAt != eTag {
			return StaleIdentifierError{StaleIdentifier: eTag}
		}

		if newShipment.DestinationAddress != nil {
			// If there is an existing DestinationAddressID associated
			// with the shipment, grab it.
			if dbShipment.DestinationAddressID != nil {
				newShipment.DestinationAddress.ID = *dbShipment.DestinationAddressID
			}
			// If there is an existing DestinationAddressID, tx.Save will use it
			// to find and update the existing record. If there isn't, it will create
			// a new record.
			err := txnAppCtx.DB().Save(newShipment.DestinationAddress)
			if err != nil {
				return err
			}
			// Make sure the shipment has the updated DestinationAddressID to store
			// in mto_shipments table
			newShipment.DestinationAddressID = &newShipment.DestinationAddress.ID

		}

		if newShipment.PickupAddress != nil {
			if dbShipment.PickupAddressID != nil {
				newShipment.PickupAddress.ID = *dbShipment.PickupAddressID
			}

			err := txnAppCtx.DB().Save(newShipment.PickupAddress)
			if err != nil {
				return err
			}

			newShipment.PickupAddressID = &newShipment.PickupAddress.ID
		}

		if newShipment.SecondaryPickupAddress != nil {
			if dbShipment.SecondaryPickupAddressID != nil {
				newShipment.SecondaryPickupAddress.ID = *dbShipment.SecondaryPickupAddressID
			}

			err := txnAppCtx.DB().Save(newShipment.SecondaryPickupAddress)
			if err != nil {
				return err
			}

			newShipment.SecondaryPickupAddressID = &newShipment.SecondaryPickupAddress.ID
		}

		if newShipment.SecondaryDeliveryAddress != nil {
			if dbShipment.SecondaryDeliveryAddressID != nil {
				newShipment.SecondaryDeliveryAddress.ID = *dbShipment.SecondaryDeliveryAddressID
			}

			err := txnAppCtx.DB().Save(newShipment.SecondaryDeliveryAddress)
			if err != nil {
				return err
			}

			newShipment.SecondaryDeliveryAddressID = &newShipment.SecondaryDeliveryAddress.ID
		}

		if newShipment.StorageFacility != nil {
			if dbShipment.StorageFacilityID != nil {
				newShipment.StorageFacility.ID = *dbShipment.StorageFacilityID
			}

			if dbShipment.StorageFacility != nil && dbShipment.StorageFacility.AddressID != uuid.Nil {
				newShipment.StorageFacility.Address.ID = dbShipment.StorageFacility.AddressID
				newShipment.StorageFacility.AddressID = dbShipment.StorageFacility.AddressID
			}
			err := txnAppCtx.DB().Save(&newShipment.StorageFacility.Address)
			if err != nil {
				return err
			}

			err = txnAppCtx.DB().Save(newShipment.StorageFacility)
			if err != nil {
				return err
			}

			newShipment.StorageFacilityID = &newShipment.StorageFacility.ID
		}

		if len(newShipment.MTOAgents) != 0 {
			for i := range newShipment.MTOAgents {
				copyOfAgent := newShipment.MTOAgents[i]

				for j := range dbShipment.MTOAgents {
					dbAgent := dbShipment.MTOAgents[j]
					// if the updates already have an agent in the system
					if dbAgent.ID == copyOfAgent.ID {
						if err := txnAppCtx.DB().Update(&copyOfAgent); err != nil {
							return err
						}
					}
				}
				if copyOfAgent.ID == uuid.Nil {
					// create a new agent if it doesn't already exist
					verrs, err := f.builder.CreateOne(txnAppCtx, &copyOfAgent)
					if verrs != nil && verrs.HasAny() {
						return verrs
					}
					if err != nil {
						return err
					}
				}
			}
		}

		// If the estimated weight was updated on an approved shipment then it would mean the move could qualify for
		// excess weight risk depending on the weight allowance and other shipment estimated weights
		if newShipment.PrimeEstimatedWeight != nil {
			if dbShipment.PrimeEstimatedWeight == nil || *newShipment.PrimeEstimatedWeight != *dbShipment.PrimeEstimatedWeight {
				/*
					TODO: If the move was already in risk of excess we need to set the status back to APPROVED if
					the new shipment estimated weight drops it out of the range. Can potentially reuse
					moveRouter.ApproveAmmendedOrders if we also add checks for excess weight there and orders
					acknowledgement
				*/
				move, verrs, err := f.moveWeights.CheckExcessWeight(txnAppCtx, dbShipment.MoveTaskOrderID, *newShipment)
				if verrs != nil && verrs.HasAny() {
					return errors.New(verrs.Error())
				}
				if err != nil {
					return err
				}

				existingMoveStatus := move.Status
				err = f.moveRouter.SendToOfficeUser(txnAppCtx, move)
				if err != nil {
					return err
				}

				if existingMoveStatus != move.Status {
					err = txnAppCtx.DB().Update(move)
					if err != nil {
						return err
					}
				}
			}
		}

		if newShipment.PrimeActualWeight != nil {
			if dbShipment.PrimeActualWeight == nil || *newShipment.PrimeActualWeight != *dbShipment.PrimeActualWeight {
				var err error
				autoReweighShipments, err = f.moveWeights.CheckAutoReweigh(txnAppCtx, dbShipment.MoveTaskOrderID, newShipment)
				if err != nil {
					return err
				}
			}
		}

		// Check that only NTS Release shipment uses that NTSRecordedWeight field
		if newShipment.NTSRecordedWeight != nil && newShipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTSDom {
			errMessage := fmt.Sprintf("field NTSRecordedWeight cannot be set for shipment type %s", string(newShipment.ShipmentType))
			return apperror.NewInvalidInputError(newShipment.ID, nil, nil, errMessage)
		}

		// If the max allowable weight for a shipment has been adjusted set a flag to recalculate payment requests for
		// this shipment
		runShipmentRecalculate := false
		if newShipment.BillableWeightCap != nil {
			// new billable cap has a value and it is not the same as the previous value
			if dbShipment.BillableWeightCap == nil || *newShipment.BillableWeightCap != *dbShipment.BillableWeightCap {
				runShipmentRecalculate = true
			}
		} else if dbShipment.BillableWeightCap != nil {
			// setting the billable cap back to nil (where previously it wasn't)
			runShipmentRecalculate = true
		}

		// A diverted shipment gets set to the SUBMITTED status automatically:
		if !dbShipment.Diversion && newShipment.Diversion {
			newShipment.Status = models.MTOShipmentStatusSubmitted

			// Get the move
			var move models.Move
			err := txnAppCtx.DB().Find(&move, dbShipment.MoveTaskOrderID)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					return apperror.NewNotFoundError(dbShipment.MoveTaskOrderID, "looking for Move")
				default:
					return apperror.NewQueryError("Move", err, "")
				}
			}

			existingMoveStatus := move.Status
			err = f.moveRouter.SendToOfficeUser(txnAppCtx, &move)
			if err != nil {
				return err
			}

			// only update if the move status has actually changed
			if existingMoveStatus != move.Status {
				err = txnAppCtx.DB().Update(&move)
				if err != nil {
					return err
				}
			}
		}

		if newShipment.TACType != nil && *newShipment.TACType == "" {
			newShipment.TACType = nil
		} else if newShipment.TACType == nil {
			newShipment.TACType = dbShipment.TACType
		}

		if newShipment.SACType != nil && *newShipment.SACType == "" {
			newShipment.SACType = nil
		} else if newShipment.SACType == nil {
			newShipment.SACType = dbShipment.SACType
		}

		if err := txnAppCtx.DB().Update(newShipment); err != nil {
			return err
		}

		//
		// Perform shipment recalculate payment request
		//
		if runShipmentRecalculate {
			_, err := f.recalculator.ShipmentRecalculatePaymentRequest(txnAppCtx, dbShipment.ID)
			if err != nil {
				return err
			}
		}

		//
		// Done with updates to shipment
		//
		return nil
	})

	if transactionError != nil {
		// Two possible types of transaction errors to handle
		if t, ok := transactionError.(StaleIdentifierError); ok {
			return apperror.NewPreconditionFailedError(dbShipment.ID, t)
		}
		return apperror.NewQueryError("mtoShipment", transactionError, "")
	}

	if len(autoReweighShipments) > 0 {
		for _, shipment := range autoReweighShipments {
			err := f.sender.SendNotification(appCtx,
				notifications.NewReweighRequested(shipment.MoveTaskOrderID, shipment),
			)
			if err != nil {
				return err
			}
		}
	}

	return nil

}

type mtoShipmentStatusUpdater struct {
	builder   UpdateMTOShipmentQueryBuilder
	siCreator services.MTOServiceItemCreator
	planner   route.Planner
}

// UpdateMTOShipmentStatus updates MTO Shipment Status
func (o *mtoShipmentStatusUpdater) UpdateMTOShipmentStatus(appCtx appcontext.AppContext, shipmentID uuid.UUID, status models.MTOShipmentStatus, rejectionReason *string, eTag string) (*models.MTOShipment, error) {
	shipment, err := fetchShipment(appCtx, shipmentID, o.builder)
	if err != nil {
		return nil, err
	}

	if status != models.MTOShipmentStatusRejected {
		rejectionReason = nil
	}

	// here we determine if the current shipment status is diversion requested before updating
	wasShipmentDiversionRequested := shipment.Status == models.MTOShipmentStatusDiversionRequested
	shipmentRouter := NewShipmentRouter()

	switch status {
	case models.MTOShipmentStatusCancellationRequested:
		err = shipmentRouter.RequestCancellation(appCtx, shipment)
	case models.MTOShipmentStatusApproved:
		err = shipmentRouter.Approve(appCtx, shipment)
	case models.MTOShipmentStatusCanceled:
		err = shipmentRouter.Cancel(appCtx, shipment)
	case models.MTOShipmentStatusDiversionRequested:
		err = shipmentRouter.RequestDiversion(appCtx, shipment)
	case models.MTOShipmentStatusRejected:
		err = shipmentRouter.Reject(appCtx, shipment, rejectionReason)
	default:
		return nil, ConflictStatusError{id: shipment.ID, transitionFromStatus: shipment.Status, transitionToStatus: status}
	}

	if err != nil {
		return nil, err
	}

	// calculate required delivery date to save it with the shipment
	if shipment.Status == models.MTOShipmentStatusApproved {
		err = o.setRequiredDeliveryDate(appCtx, shipment)
		if err != nil {
			return nil, err
		}
	}

	verrs, err := o.builder.UpdateOne(appCtx, shipment, &eTag)

	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return nil, apperror.NewPreconditionFailedError(shipment.ID, err)
		default:
			return nil, err
		}
	}

	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(shipment.ID, nil, verrs, "There was an issue with validating the updates")
	}

	// after updating shipment
	// create shipment level service items if shipment status was NOT diversion requested before it was updated
	// and current status is approved
	createSSI := shipment.Status == models.MTOShipmentStatusApproved && !wasShipmentDiversionRequested
	if createSSI {
		err = o.createShipmentServiceItems(appCtx, shipment)
		if err != nil {
			return nil, err
		}
	}

	return shipment, nil
}

// createShipmentServiceItems creates shipment level service items
func (o *mtoShipmentStatusUpdater) createShipmentServiceItems(appCtx appcontext.AppContext, shipment *models.MTOShipment) error {
	reServiceCodes := reServiceCodesForShipment(*shipment)
	serviceItemsToCreate := constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceCodes)
	for _, serviceItem := range serviceItemsToCreate {
		copyOfServiceItem := serviceItem // Make copy to avoid implicit memory aliasing of items from a range statement.
		_, verrs, err := o.siCreator.CreateMTOServiceItem(appCtx, &copyOfServiceItem)

		if verrs != nil && verrs.HasAny() {
			invalidInputError := apperror.NewInvalidInputError(shipment.ID, nil, verrs, "There was an issue creating service items for the shipment")
			return invalidInputError
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// setRequiredDeliveryDate set the calculated delivery date for the shipment
func (o *mtoShipmentStatusUpdater) setRequiredDeliveryDate(appCtx appcontext.AppContext, shipment *models.MTOShipment) error {
	if shipment.ScheduledPickupDate != nil &&
		shipment.RequiredDeliveryDate == nil &&
		(shipment.PrimeEstimatedWeight != nil || shipment.NTSRecordedWeight != nil) {

		var pickupLocation *models.Address
		var deliveryLocation *models.Address
		weight := shipment.PrimeEstimatedWeight

		switch shipment.ShipmentType {
		case models.MTOShipmentTypeHHGIntoNTSDom:
			if shipment.StorageFacility == nil || shipment.StorageFacility.AddressID == uuid.Nil {
				return errors.Errorf("StorageFacility is required for %s shipments", models.MTOShipmentTypeHHGIntoNTSDom)
			}
			err := appCtx.DB().Load(shipment.StorageFacility, "Address")
			if err != nil {
				return apperror.NewNotFoundError(shipment.StorageFacility.AddressID, "looking for MTOShipment.StorageFacility.Address")
			}

			pickupLocation = shipment.PickupAddress
			deliveryLocation = &shipment.StorageFacility.Address
		case models.MTOShipmentTypeHHGOutOfNTSDom:
			if shipment.StorageFacility == nil || shipment.StorageFacility.AddressID == uuid.Nil {
				return errors.Errorf("StorageFacility is required for %s shipments", models.MTOShipmentTypeHHGOutOfNTSDom)
			}
			err := appCtx.DB().Load(shipment.StorageFacility, "Address")
			if err != nil {
				return apperror.NewNotFoundError(shipment.StorageFacility.AddressID, "looking for MTOShipment.StorageFacility.Address")
			}
			pickupLocation = &shipment.StorageFacility.Address
			deliveryLocation = shipment.DestinationAddress
			weight = shipment.NTSRecordedWeight
		default:
			pickupLocation = shipment.PickupAddress
			deliveryLocation = shipment.DestinationAddress
		}
		requiredDeliveryDate, calcErr := CalculateRequiredDeliveryDate(appCtx, o.planner, *pickupLocation, *deliveryLocation, *shipment.ScheduledPickupDate, weight.Int())
		if calcErr != nil {
			return calcErr
		}

		shipment.RequiredDeliveryDate = requiredDeliveryDate
	}

	return nil
}

func fetchShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID, builder UpdateMTOShipmentQueryBuilder) (*models.MTOShipment, error) {
	var shipment models.MTOShipment

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", shipmentID),
	}
	err := builder.FetchOne(appCtx, &shipment, queryFilters)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(shipmentID, "Shipment not found")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	return &shipment, nil
}

func reServiceCodesForShipment(shipment models.MTOShipment) []models.ReServiceCode {
	// We will detect the type of shipment we're working with and then call a helper with the correct
	// default service items that we want created as a side effect.
	// More info in MB-1140: https://dp3.atlassian.net/browse/MB-1140

	switch shipment.ShipmentType {
	case models.MTOShipmentTypeHHG, models.MTOShipmentTypeHHGLongHaulDom:
		// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Packing, and Dom Unpacking.
		return []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDPK,
			models.ReServiceCodeDUPK,
		}
	case models.MTOShipmentTypeHHGShortHaulDom:
		// Need to create: Dom Shorthaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Packing, Dom Unpacking
		return []models.ReServiceCode{
			models.ReServiceCodeDSH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDPK,
			models.ReServiceCodeDUPK,
		}
	case models.MTOShipmentTypeHHGIntoNTSDom:
		// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom NTS Packing
		return []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDNPK,
		}
	case models.MTOShipmentTypeHHGOutOfNTSDom:
		// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Unpacking
		return []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDUPK,
		}
	case models.MTOShipmentTypeMotorhome:
		// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Mobile Home Factor
		return []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDMHF,
		}
	case models.MTOShipmentTypeBoatHaulAway:
		// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Haul Away Boat Factor
		return []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDBHF,
		}
	case models.MTOShipmentTypeBoatTowAway:
		// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Tow Away Boat Factor
		return []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDBTF,
		}
	}

	return []models.ReServiceCode{}
}

// CalculateRequiredDeliveryDate function is used to get a distance calculation using the pickup and destination addresses. It then uses
// the value returned to make a fetch on the ghc_domestic_transit_times table and returns a required delivery date
// based on the max_days_transit_time.
func CalculateRequiredDeliveryDate(appCtx appcontext.AppContext, planner route.Planner, pickupAddress models.Address, destinationAddress models.Address, pickupDate time.Time, weight int) (*time.Time, error) {
	// Okay, so this is something to get us able to take care of the 20 day condition over in the gdoc linked in this
	// story: https://dp3.atlassian.net/browse/MB-1141
	// We unfortunately didn't get a lot of guidance regarding vicinity. So for now we're taking zip codes that are the
	// explicitly mentioned 20 day cities and those in the same county (that I've manually compiled together here).
	// If a move is in that group it adds 20 days, if it's not in that group, but is in Alaska it adds 10 days.
	// Else it will not do either of those things.
	// The cities for 20 days are: Adak, Kodiak, Juneau, Ketchikan, and Sitka. As well as others in their 'vicinity.'
	twentyDayAKZips := [28]string{"99546", "99547", "99591", "99638", "99660", "99685", "99692", "99550", "99608",
		"99615", "99619", "99624", "99643", "99644", "99697", "99650", "99801", "99802", "99803", "99811", "99812",
		"99950", "99824", "99850", "99901", "99928", "99950", "99835"}

	// Get a distance calculation between pickup and destination addresses.
	distance, err := planner.TransitDistance(appCtx, &pickupAddress, &destinationAddress)
	if err != nil {
		return nil, err
	}
	// Query the ghc_domestic_transit_times table for the max transit time
	var ghcDomesticTransitTime models.GHCDomesticTransitTime
	err = appCtx.DB().Where("distance_miles_lower <= ? "+
		"AND distance_miles_upper >= ? "+
		"AND weight_lbs_lower <= ? "+
		"AND (weight_lbs_upper >= ? OR weight_lbs_upper = 0)",
		distance, distance, weight, weight).First(&ghcDomesticTransitTime)

	if err != nil {
		return nil, errors.Errorf("failed to find transit time for shipment of %d lbs weight and %d mile distance", weight, distance)
	}

	// Add the max transit time to the pickup date to get the new required delivery date
	requiredDeliveryDate := pickupDate.AddDate(0, 0, ghcDomesticTransitTime.MaxDaysTransitTime)

	// Let's add some days if we're dealing with an alaska shipment.
	if destinationAddress.State == "AK" {
		for _, zip := range twentyDayAKZips {
			if destinationAddress.PostalCode == zip {
				// Add an extra 10 days here, so that after we add the 10 for being in AK we wind up with a total of 20
				requiredDeliveryDate = requiredDeliveryDate.AddDate(0, 0, 10)
				break
			}
		}
		// Add an extra 10 days for being in AK
		requiredDeliveryDate = requiredDeliveryDate.AddDate(0, 0, 10)
	}

	// return the value
	return &requiredDeliveryDate, nil
}

// This private function is used to generically construct service items when shipments are approved.
func constructMTOServiceItemModels(shipmentID uuid.UUID, mtoID uuid.UUID, reServiceCodes []models.ReServiceCode) models.MTOServiceItems {
	serviceItems := make(models.MTOServiceItems, len(reServiceCodes))
	currentTime := swag.Time(time.Now())
	for i, reServiceCode := range reServiceCodes {
		serviceItem := models.MTOServiceItem{
			MoveTaskOrderID: mtoID,
			MTOShipmentID:   &shipmentID,
			ReService:       models.ReService{Code: reServiceCode},
			Status:          "APPROVED",
			ApprovedAt:      currentTime,
		}
		serviceItems[i] = serviceItem
	}
	return serviceItems
}

// NewMTOShipmentStatusUpdater creates a new MTO Shipment Status Updater
func NewMTOShipmentStatusUpdater(builder UpdateMTOShipmentQueryBuilder, siCreator services.MTOServiceItemCreator, planner route.Planner) services.MTOShipmentStatusUpdater {
	return &mtoShipmentStatusUpdater{builder, siCreator, planner}
}

// ConflictStatusError returns an error for a conflict in status
type ConflictStatusError struct {
	id                        uuid.UUID
	transitionFromStatus      models.MTOShipmentStatus
	transitionToStatus        models.MTOShipmentStatus
	transitionAllowedStatuses *[]models.MTOShipmentStatus
}

// Error is the string representation of the error
func (e ConflictStatusError) Error() string {
	if e.transitionAllowedStatuses != nil {
		return fmt.Sprintf("Shipment with id '%s' can only transition to status '%s' from %+q, but its current status is '%s'",
			e.id.String(), e.transitionToStatus, *e.transitionAllowedStatuses, e.transitionFromStatus)
	}

	return ""
}

// MTOShipmentsMTOAvailableToPrime checks if a given shipment is available to the Prime
// TODO: trend away from using this method, it represents *business logic* and should
// ideally be done only as an internal check, rather than relying on being invoked
// by the handler layer
func (f mtoShipmentUpdater) MTOShipmentsMTOAvailableToPrime(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (bool, error) {
	err := checkAvailToPrime().Validate(appCtx, &models.MTOShipment{ID: mtoShipmentID}, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}
