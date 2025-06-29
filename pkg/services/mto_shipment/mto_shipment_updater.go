package mtoshipment

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
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
	addressUpdater services.AddressUpdater
	addressCreator services.AddressCreator
	planner        route.Planner
	moveRouter     services.MoveRouter
	moveWeights    services.MoveWeights
	recalculator   services.PaymentRequestShipmentRecalculator
	checks         []validator
}

// NewMTOShipmentUpdater creates a new struct with the service dependencies
func NewMTOShipmentUpdater(builder UpdateMTOShipmentQueryBuilder, _ services.Fetcher, planner route.Planner, moveRouter services.MoveRouter, moveWeights services.MoveWeights, sender notifications.NotificationSender, recalculator services.PaymentRequestShipmentRecalculator, addressUpdater services.AddressUpdater, addressCreator services.AddressCreator) services.MTOShipmentUpdater {
	return &mtoShipmentUpdater{
		builder,
		fetch.NewFetcher(builder),
		addressUpdater,
		addressCreator,
		planner,
		moveRouter,
		moveWeights,
		recalculator,
		[]validator{},
	}
}

// TODO: apply the subset of business logic validations
// that would be appropriate for the CUSTOMER
func NewCustomerMTOShipmentUpdater(builder UpdateMTOShipmentQueryBuilder, _ services.Fetcher, planner route.Planner, moveRouter services.MoveRouter, moveWeights services.MoveWeights, sender notifications.NotificationSender, recalculator services.PaymentRequestShipmentRecalculator, addressUpdater services.AddressUpdater, addressCreator services.AddressCreator) services.MTOShipmentUpdater {
	return &mtoShipmentUpdater{
		builder,
		fetch.NewFetcher(builder),
		addressUpdater,
		addressCreator,
		planner,
		moveRouter,
		moveWeights,
		recalculator,
		[]validator{checkStatus(), MTOShipmentHasTertiaryAddressWithNoSecondaryAddressUpdate(), checkUBShipmentOCONUSRequirement(), MTOShipmentHasValidRequestedPickupDate()},
	}
}

func NewOfficeMTOShipmentUpdater(builder UpdateMTOShipmentQueryBuilder, _ services.Fetcher, planner route.Planner, moveRouter services.MoveRouter, moveWeights services.MoveWeights, sender notifications.NotificationSender, recalculator services.PaymentRequestShipmentRecalculator, addressUpdater services.AddressUpdater, addressCreator services.AddressCreator) services.MTOShipmentUpdater {
	return &mtoShipmentUpdater{
		builder,
		fetch.NewFetcher(builder),
		addressUpdater,
		addressCreator,
		planner,
		moveRouter,
		moveWeights,
		recalculator,
		[]validator{checkStatus(), checkUpdateAllowed(), MTOShipmentHasTertiaryAddressWithNoSecondaryAddressUpdate(), checkUBShipmentOCONUSRequirement(), MTOShipmentHasValidRequestedPickupDate()},
	}
}

// TODO: apply the subset of business logic validations
// that would be appropriate for the PRIME
func NewPrimeMTOShipmentUpdater(builder UpdateMTOShipmentQueryBuilder, _ services.Fetcher, planner route.Planner, moveRouter services.MoveRouter, moveWeights services.MoveWeights, sender notifications.NotificationSender, recalculator services.PaymentRequestShipmentRecalculator, addressUpdater services.AddressUpdater, addressCreator services.AddressCreator) services.MTOShipmentUpdater {
	return &mtoShipmentUpdater{
		builder,
		fetch.NewFetcher(builder),
		addressUpdater,
		addressCreator,
		planner,
		moveRouter,
		moveWeights,
		recalculator,
		[]validator{checkStatus(), checkAvailToPrime(), checkPrimeValidationsOnModel(planner), MTOShipmentHasTertiaryAddressWithNoSecondaryAddressUpdate(), checkUBShipmentOCONUSRequirement()},
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

	if requestedUpdatedShipment.ActualDeliveryDate != nil {
		dbShipment.ActualDeliveryDate = requestedUpdatedShipment.ActualDeliveryDate
	}

	if requestedUpdatedShipment.ScheduledDeliveryDate != nil {
		dbShipment.ScheduledDeliveryDate = requestedUpdatedShipment.ScheduledDeliveryDate
	}

	if requestedUpdatedShipment.PrimeEstimatedWeight != nil {
		now := time.Now()
		dbShipment.PrimeEstimatedWeight = requestedUpdatedShipment.PrimeEstimatedWeight
		dbShipment.PrimeEstimatedWeightRecordedDate = &now
	}

	if requestedUpdatedShipment.NTSRecordedWeight != nil {
		dbShipment.NTSRecordedWeight = requestedUpdatedShipment.NTSRecordedWeight
	}

	if requestedUpdatedShipment.PickupAddress != nil && dbShipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTS {
		dbShipment.PickupAddress = requestedUpdatedShipment.PickupAddress
	}

	if requestedUpdatedShipment.DestinationAddress != nil && dbShipment.ShipmentType != models.MTOShipmentTypeHHGIntoNTS {
		dbShipment.DestinationAddress = requestedUpdatedShipment.DestinationAddress
	}

	if requestedUpdatedShipment.DestinationType != nil {
		dbShipment.DestinationType = requestedUpdatedShipment.DestinationType
	}

	// If HasSecondaryPickupAddress is false, we want to remove the secondary address as well as the tertiary address
	// Otherwise, if a non-nil address is in the payload, we should save it
	if requestedUpdatedShipment.HasSecondaryPickupAddress != nil && !*requestedUpdatedShipment.HasSecondaryPickupAddress {
		dbShipment.HasSecondaryPickupAddress = requestedUpdatedShipment.HasSecondaryPickupAddress
		dbShipment.SecondaryPickupAddress = nil
		dbShipment.SecondaryPickupAddressID = nil
		requestedUpdatedShipment.HasTertiaryPickupAddress = models.BoolPointer(false)
	} else if requestedUpdatedShipment.SecondaryPickupAddress != nil {
		dbShipment.SecondaryPickupAddress = requestedUpdatedShipment.SecondaryPickupAddress
		dbShipment.HasSecondaryPickupAddress = models.BoolPointer(true)
	}

	// If HasSecondaryDeliveryAddress is false, we want to remove the secondary address as well as the tertiary address
	// Otherwise, if a non-nil address is in the payload, we should save it
	if requestedUpdatedShipment.HasSecondaryDeliveryAddress != nil && !*requestedUpdatedShipment.HasSecondaryDeliveryAddress {
		dbShipment.HasSecondaryDeliveryAddress = requestedUpdatedShipment.HasSecondaryDeliveryAddress
		dbShipment.SecondaryDeliveryAddress = nil
		dbShipment.SecondaryDeliveryAddressID = nil
		requestedUpdatedShipment.HasTertiaryDeliveryAddress = models.BoolPointer(false)
	} else if requestedUpdatedShipment.SecondaryDeliveryAddress != nil {
		dbShipment.SecondaryDeliveryAddress = requestedUpdatedShipment.SecondaryDeliveryAddress
		dbShipment.HasSecondaryDeliveryAddress = models.BoolPointer(true)
	}

	// If HasTertiaryPickupAddress is false, we want to remove the address
	// Otherwise, if a non-nil address is in the payload, we should save it
	if requestedUpdatedShipment.HasTertiaryPickupAddress != nil && !*requestedUpdatedShipment.HasTertiaryPickupAddress {
		dbShipment.HasTertiaryPickupAddress = requestedUpdatedShipment.HasTertiaryPickupAddress
		dbShipment.TertiaryPickupAddress = nil
		dbShipment.TertiaryPickupAddressID = nil
	} else if requestedUpdatedShipment.TertiaryPickupAddress != nil {
		dbShipment.TertiaryPickupAddress = requestedUpdatedShipment.TertiaryPickupAddress
		dbShipment.HasTertiaryPickupAddress = models.BoolPointer(true)
	}

	// If HasTertiaryDeliveryAddress is false, we want to remove the address
	// Otherwise, if a non-nil address is in the payload, we should save it
	if requestedUpdatedShipment.HasTertiaryDeliveryAddress != nil && !*requestedUpdatedShipment.HasTertiaryDeliveryAddress {
		dbShipment.HasTertiaryDeliveryAddress = requestedUpdatedShipment.HasTertiaryDeliveryAddress
		dbShipment.TertiaryDeliveryAddress = nil
		dbShipment.TertiaryDeliveryAddressID = nil
	} else if requestedUpdatedShipment.TertiaryDeliveryAddress != nil {
		dbShipment.TertiaryDeliveryAddress = requestedUpdatedShipment.TertiaryDeliveryAddress
		dbShipment.HasTertiaryDeliveryAddress = models.BoolPointer(true)
	}

	if requestedUpdatedShipment.ShipmentType != "" {
		dbShipment.ShipmentType = requestedUpdatedShipment.ShipmentType
	}

	if requestedUpdatedShipment.Status != "" {
		dbShipment.Status = requestedUpdatedShipment.Status
	}

	if requestedUpdatedShipment.ApprovedDate != nil {
		dbShipment.ApprovedDate = requestedUpdatedShipment.ApprovedDate
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

	if requestedUpdatedShipment.TACType != nil {
		dbShipment.TACType = requestedUpdatedShipment.TACType
	}

	if requestedUpdatedShipment.SACType != nil {
		dbShipment.SACType = requestedUpdatedShipment.SACType
	}

	if requestedUpdatedShipment.ServiceOrderNumber != nil {
		dbShipment.ServiceOrderNumber = requestedUpdatedShipment.ServiceOrderNumber
	}

	if requestedUpdatedShipment.StorageFacility != nil {
		dbShipment.StorageFacility = requestedUpdatedShipment.StorageFacility
	}

	if requestedUpdatedShipment.ActualProGearWeight != nil {
		dbShipment.ActualProGearWeight = requestedUpdatedShipment.ActualProGearWeight
	}

	if requestedUpdatedShipment.ActualSpouseProGearWeight != nil {
		dbShipment.ActualSpouseProGearWeight = requestedUpdatedShipment.ActualSpouseProGearWeight
	}

	if requestedUpdatedShipment.OriginSITAuthEndDate != nil {
		dbShipment.OriginSITAuthEndDate = requestedUpdatedShipment.OriginSITAuthEndDate
	}

	if requestedUpdatedShipment.DestinationSITAuthEndDate != nil {
		dbShipment.DestinationSITAuthEndDate = requestedUpdatedShipment.DestinationSITAuthEndDate
	}

	//// TODO: move mtoagent creation into service: Should not update MTOAgents here because we don't have an eTag
	if len(requestedUpdatedShipment.MTOAgents) > 0 {
		var agentsToCreateOrUpdate []models.MTOAgent
		for _, newAgentInfo := range requestedUpdatedShipment.MTOAgents {
			// if no record exists in the db
			if newAgentInfo.ID == uuid.Nil {
				newAgentInfo.MTOShipmentID = requestedUpdatedShipment.ID
				if newAgentInfo.FirstName != nil && *newAgentInfo.FirstName == "" {
					newAgentInfo.FirstName = nil
				}
				if newAgentInfo.LastName != nil && *newAgentInfo.LastName == "" {
					newAgentInfo.LastName = nil
				}
				if newAgentInfo.Email != nil && *newAgentInfo.Email == "" {
					newAgentInfo.Email = nil
				}
				if newAgentInfo.Phone != nil && *newAgentInfo.Phone == "" {
					newAgentInfo.Phone = nil
				}
				// If no fields are set, then we do not want to create the MTO agent
				if newAgentInfo.FirstName == nil && newAgentInfo.LastName == nil && newAgentInfo.Email == nil && newAgentInfo.Phone == nil {
					continue
				}
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
						dbShipment.MTOAgents[i].FirstName = services.SetOptionalStringField(newAgentInfo.FirstName, dbShipment.MTOAgents[i].FirstName)
						dbShipment.MTOAgents[i].LastName = services.SetOptionalStringField(newAgentInfo.LastName, dbShipment.MTOAgents[i].LastName)
						dbShipment.MTOAgents[i].Email = services.SetOptionalStringField(newAgentInfo.Email, dbShipment.MTOAgents[i].Email)
						dbShipment.MTOAgents[i].Phone = services.SetOptionalStringField(newAgentInfo.Phone, dbShipment.MTOAgents[i].Phone)
						// If no fields are set, then we will soft-delete the MTO agent
						if dbShipment.MTOAgents[i].FirstName == nil && dbShipment.MTOAgents[i].LastName == nil && dbShipment.MTOAgents[i].Email == nil && dbShipment.MTOAgents[i].Phone == nil {
							err := utilities.SoftDestroy(appCtx.DB(), &dbShipment.MTOAgents[i])
							if err != nil {
								appCtx.Logger().Error("Error soft destroying MTO Agent.")
								continue
							}
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

// UpdateMTOShipment updates the mto shipment
func (f *mtoShipmentUpdater) UpdateMTOShipment(appCtx appcontext.AppContext, mtoShipment *models.MTOShipment, eTag string, api string) (*models.MTOShipment, error) {
	eagerAssociations := []string{"MoveTaskOrder",
		"PickupAddress",
		"DestinationAddress",
		"SecondaryPickupAddress",
		"SecondaryDeliveryAddress",
		"TertiaryPickupAddress",
		"TertiaryDeliveryAddress",
		"SITDurationUpdates",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOServiceItems.CustomerContacts",
		"StorageFacility.Address",
		"Reweigh",
		"ShipmentLocator",
	}

	oldShipment, err := FindShipment(appCtx, mtoShipment.ID, eagerAssociations...)
	if err != nil {
		return nil, err
	}

	var agents []models.MTOAgent
	err = appCtx.DB().Scope(utilities.ExcludeDeletedScope()).Where("mto_shipment_id = ?", mtoShipment.ID).All(&agents)
	if err != nil {
		return nil, err
	}
	oldShipment.MTOAgents = agents

	// run the (read-only) validations
	if verr := validateShipment(appCtx, mtoShipment, oldShipment, f.checks...); verr != nil {
		return nil, verr
	}

	// save the original db version, oldShipment will be modified
	var dbShipment models.MTOShipment
	err = copier.CopyWithOption(&dbShipment, oldShipment, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, fmt.Errorf("error copying shipment data %w", err)
	}
	setNewShipmentFields(appCtx, oldShipment, mtoShipment)
	newShipment := oldShipment // old shipment has now been updated with requested changes
	// db version is used to check if agents need creating or updating
	err = f.updateShipmentRecord(appCtx, &dbShipment, newShipment, eTag)
	if err != nil {
		switch typedErr := err.(type) {
		case StaleIdentifierError:
			return nil, apperror.NewPreconditionFailedError(mtoShipment.ID, typedErr)
		case apperror.InvalidInputError:
			return nil, apperror.NewInvalidInputError(mtoShipment.ID, typedErr, typedErr.ValidationErrors, "Invalid input found while updating the shipment")
		default:
			return nil, err
		}
	}

	updatedShipment, err := FindShipment(appCtx, mtoShipment.ID, eagerAssociations...)
	if err != nil {
		return nil, err
	}

	var updatedAgents []models.MTOAgent
	err = appCtx.DB().Scope(utilities.ExcludeDeletedScope()).Where("mto_shipment_id = ?", mtoShipment.ID).All(&updatedAgents)
	if err != nil {
		return nil, err
	}
	updatedShipment.MTOAgents = updatedAgents

	// As the API is passed through, we want to apply a breaking change without duplicating lots of code.
	// 'prime' is the V1 version of this endpoint. All endpoints besides the prime should be utilizing new logic
	// of this function where it no longer calls UpdateDestinationSITServiceItemsAddress. UpdateDestinationSITServiceItemsAddress
	// has been deprecated out of this function per E-04819
	if api == "prime" {
		err = UpdateDestinationSITServiceItemsAddress(appCtx, updatedShipment)
		if err != nil {
			return nil, err
		}
	}

	return updatedShipment, nil
}

// Takes the validated shipment input and updates the database using a transaction. If any part of the
// update fails, the entire transaction will be rolled back.
func (f *mtoShipmentUpdater) updateShipmentRecord(appCtx appcontext.AppContext, dbShipment *models.MTOShipment, newShipment *models.MTOShipment, eTag string) error {
	var verrs *validate.Errors
	var move *models.Move
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// temp optimistic locking solution til query builder is re-tooled to handle nested updates
		updatedAt, err := etag.DecodeEtag(eTag)
		if err != nil {
			return StaleIdentifierError{StaleIdentifier: eTag}
		}

		if !updatedAt.Equal(dbShipment.UpdatedAt) {
			return StaleIdentifierError{StaleIdentifier: eTag}
		}

		// TODO: We currently can't distinguish between a nil DestinationAddress meaning to "clear field"
		//   vs "don't touch" the field, so we can't safely reset a nil DestinationAddress to the duty
		//   location address for an HHG like we do in the MTOShipmentCreator now.  See MB-15718.

		if newShipment.DestinationAddress != nil && newShipment.ShipmentType != models.MTOShipmentTypeHHGIntoNTS {
			// If there is an existing DestinationAddressID associated
			// with the shipment, grab it.
			if dbShipment.DestinationAddressID != nil {
				newShipment.DestinationAddress.ID = *dbShipment.DestinationAddressID
			}

			// Only call the address updater service if there is an original destination address to be updated at all
			if dbShipment.DestinationAddress != nil {
				newDestinationAddress, destAddErr := f.addressUpdater.UpdateAddress(txnAppCtx, newShipment.DestinationAddress, etag.GenerateEtag(dbShipment.DestinationAddress.UpdatedAt))
				if destAddErr != nil {
					return destAddErr
				}
				// Make sure the shipment has the updated DestinationAddressID to store
				// in mto_shipments table
				newShipment.DestinationAddressID = &newDestinationAddress.ID
				newShipment.DestinationAddress = newDestinationAddress
			} else if newShipment.DestinationAddressID == nil {
				// There is no original address to update
				if newShipment.DestinationAddress.ID == uuid.Nil {
					// And this new address does not have an ID.
					// We need to create a new one.
					newDestinationAddress, newDestAddErr := f.addressCreator.CreateAddress(appCtx, newShipment.DestinationAddress)
					if newDestAddErr != nil {
						return newDestAddErr
					}
					newShipment.DestinationAddressID = &newDestinationAddress.ID
					newShipment.DestinationAddress = newDestinationAddress
				} else {
					// Otherwise, there is no original address to update and this new address already has an ID
					newShipment.DestinationAddressID = &newShipment.DestinationAddress.ID
				}

			}

		}

		if newShipment.PickupAddress != nil && newShipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTS {
			if dbShipment.PickupAddressID != nil {
				newShipment.PickupAddress.ID = *dbShipment.PickupAddressID
			}

			// If there is an existing, original address then we need to update it
			if dbShipment.PickupAddress != nil {
				newPickupAddress, newPickupErr := f.addressUpdater.UpdateAddress(txnAppCtx, newShipment.PickupAddress, etag.GenerateEtag(dbShipment.PickupAddress.UpdatedAt))
				if newPickupErr != nil {
					return newPickupErr
				}

				newShipment.PickupAddressID = &newPickupAddress.ID
				newShipment.PickupAddress = newPickupAddress
			} else if newShipment.PickupAddressID == nil {
				// There is no original address to update
				if newShipment.PickupAddress.ID == uuid.Nil {
					// And this new address does not have an ID.
					// We need to create a new one.
					newPickupAddress, newPickupAddCreateErr := f.addressCreator.CreateAddress(appCtx, newShipment.PickupAddress)
					if newPickupAddCreateErr != nil {
						return newPickupAddCreateErr
					}
					newShipment.PickupAddressID = &newPickupAddress.ID
					newShipment.PickupAddress = newPickupAddress
				} else {
					// Otherwise, there is no original address to update and this new address already has an ID
					newShipment.PickupAddressID = &newShipment.PickupAddress.ID
				}

			}
		}
		if newShipment.HasSecondaryPickupAddress != nil {
			if !*newShipment.HasSecondaryPickupAddress {
				newShipment.SecondaryDeliveryAddressID = nil
			}
		}
		if newShipment.HasSecondaryDeliveryAddress != nil {
			if !*newShipment.HasSecondaryDeliveryAddress {
				newShipment.SecondaryDeliveryAddressID = nil
			}
		}

		if newShipment.HasSecondaryPickupAddress != nil && *newShipment.HasSecondaryPickupAddress && newShipment.SecondaryPickupAddress != nil {
			if dbShipment.SecondaryPickupAddressID != nil {
				newShipment.SecondaryPickupAddress.ID = *dbShipment.SecondaryPickupAddressID
			}

			if dbShipment.SecondaryPickupAddress != nil {
				// Secondary pickup address exists, meaning it should be updated
				newSecondaryPickupAddress, newSecondaryPickupUpdateErr := f.addressUpdater.UpdateAddress(txnAppCtx, newShipment.SecondaryPickupAddress, etag.GenerateEtag(dbShipment.SecondaryPickupAddress.UpdatedAt))
				if newSecondaryPickupUpdateErr != nil {
					return newSecondaryPickupUpdateErr
				}
				newShipment.SecondaryPickupAddressID = &newSecondaryPickupAddress.ID
			} else if newShipment.SecondaryPickupAddressID == nil {
				// Secondary pickup address appears to not exist yet, meaning it should be created
				if newShipment.SecondaryPickupAddress.ID == uuid.Nil {
					newSecondaryPickupAddress, newSecondaryPickupCreateErr := f.addressCreator.CreateAddress(txnAppCtx, newShipment.SecondaryPickupAddress)
					if newSecondaryPickupCreateErr != nil {
						return newSecondaryPickupCreateErr
					}
					newShipment.SecondaryPickupAddressID = &newSecondaryPickupAddress.ID
				} else {
					// No original address to update, and the new address already has an ID so we should just assign it to the shipment
					newShipment.SecondaryPickupAddressID = &newShipment.SecondaryPickupAddress.ID
				}
			}
		}

		if newShipment.SecondaryDeliveryAddress != nil {
			if dbShipment.SecondaryDeliveryAddressID != nil {
				newShipment.SecondaryDeliveryAddress.ID = *dbShipment.SecondaryDeliveryAddressID
			}

			if dbShipment.SecondaryDeliveryAddress != nil {
				// Secondary delivery address exists, meaning it should be updated
				newSecondaryDeliveryAddress, secondaryDeliveryUpdateErr := f.addressUpdater.UpdateAddress(txnAppCtx, newShipment.SecondaryDeliveryAddress, etag.GenerateEtag(dbShipment.SecondaryDeliveryAddress.UpdatedAt))
				if secondaryDeliveryUpdateErr != nil {
					return secondaryDeliveryUpdateErr
				}
				newShipment.SecondaryDeliveryAddressID = &newSecondaryDeliveryAddress.ID
			} else if newShipment.SecondaryDeliveryAddressID == nil {
				// Secondary delivery address appears to not exist yet, meaning it should be created
				if newShipment.SecondaryDeliveryAddress.ID == uuid.Nil {
					newSecondaryDeliveryAddress, secondaryDeliveryCreateErr := f.addressCreator.CreateAddress(txnAppCtx, newShipment.SecondaryDeliveryAddress)
					if secondaryDeliveryCreateErr != nil {
						return secondaryDeliveryCreateErr
					}
					newShipment.SecondaryDeliveryAddressID = &newSecondaryDeliveryAddress.ID
				} else {
					// No original address to update, and the new address already has an ID so we should just assign it to the shipment
					newShipment.SecondaryDeliveryAddressID = &newShipment.SecondaryDeliveryAddress.ID
				}
			}
		}

		if newShipment.HasTertiaryPickupAddress != nil {
			if !*newShipment.HasTertiaryPickupAddress {
				newShipment.TertiaryPickupAddressID = nil
			}
		}
		if newShipment.HasTertiaryDeliveryAddress != nil {
			if !*newShipment.HasTertiaryDeliveryAddress {
				newShipment.TertiaryDeliveryAddressID = nil
			}
		}

		if newShipment.HasTertiaryPickupAddress != nil && *newShipment.HasTertiaryPickupAddress && newShipment.TertiaryPickupAddress != nil {
			if dbShipment.TertiaryPickupAddress != nil {
				newShipment.TertiaryPickupAddress.ID = *dbShipment.TertiaryPickupAddressID
			}

			if dbShipment.TertiaryPickupAddress != nil {
				// Tertiary pickup address exists, meaning it should be updated
				newTertiaryPickupAddress, newTertiaryPickupUpdateErr := f.addressUpdater.UpdateAddress(txnAppCtx, newShipment.TertiaryPickupAddress, etag.GenerateEtag(dbShipment.TertiaryPickupAddress.UpdatedAt))
				if newTertiaryPickupUpdateErr != nil {
					return newTertiaryPickupUpdateErr
				}
				newShipment.TertiaryPickupAddressID = &newTertiaryPickupAddress.ID
			} else if newShipment.TertiaryPickupAddressID == nil {
				// Tertiary pickup address appears to not exist yet, meaning it should be created
				if newShipment.TertiaryPickupAddress.ID == uuid.Nil {
					newTertiaryPickupAddress, newTertiaryPickupCreateErr := f.addressCreator.CreateAddress(txnAppCtx, newShipment.TertiaryPickupAddress)
					if newTertiaryPickupCreateErr != nil {
						return newTertiaryPickupCreateErr
					}
					newShipment.TertiaryPickupAddressID = &newTertiaryPickupAddress.ID
				} else {
					// No original address to update, and the new address already has an ID so we should just assign it to the shipment
					newShipment.TertiaryPickupAddressID = &newShipment.TertiaryPickupAddress.ID
				}
			}
		}

		if newShipment.TertiaryDeliveryAddress != nil {
			if dbShipment.TertiaryDeliveryAddressID != nil {
				newShipment.TertiaryDeliveryAddress.ID = *dbShipment.TertiaryDeliveryAddressID
			}

			if dbShipment.TertiaryDeliveryAddress != nil {
				// Tertiary delivery address exists, meaning it should be updated
				newTertiaryDeliveryAddress, tertiaryDeliveryUpdateErr := f.addressUpdater.UpdateAddress(txnAppCtx, newShipment.TertiaryDeliveryAddress, etag.GenerateEtag(dbShipment.TertiaryDeliveryAddress.UpdatedAt))
				if tertiaryDeliveryUpdateErr != nil {
					return tertiaryDeliveryUpdateErr
				}
				newShipment.TertiaryDeliveryAddressID = &newTertiaryDeliveryAddress.ID
			} else if newShipment.TertiaryDeliveryAddressID == nil {
				// Tertiary delivery address appears to not exist yet, meaning it should be created
				if newShipment.TertiaryDeliveryAddress.ID == uuid.Nil {
					newTertiaryDeliveryAddress, tertiaryDeliveryCreateErr := f.addressCreator.CreateAddress(txnAppCtx, newShipment.TertiaryDeliveryAddress)
					if tertiaryDeliveryCreateErr != nil {
						return tertiaryDeliveryCreateErr
					}
					newShipment.TertiaryDeliveryAddressID = &newTertiaryDeliveryAddress.ID
				} else {
					// No original address to update, and the new address already has an ID so we should just assign it to the shipment
					newShipment.TertiaryDeliveryAddressID = &newShipment.TertiaryDeliveryAddress.ID
				}
			}
		}

		if newShipment.StorageFacility != nil {
			if dbShipment.StorageFacilityID != nil {
				newShipment.StorageFacility.ID = *dbShipment.StorageFacilityID
			}

			if dbShipment.StorageFacility != nil && dbShipment.StorageFacility.AddressID != uuid.Nil {
				newShipment.StorageFacility.Address.ID = dbShipment.StorageFacility.AddressID
				newShipment.StorageFacility.AddressID = dbShipment.StorageFacility.AddressID
			}
			if dbShipment.StorageFacility != nil {
				// Storage facility address exists, meaning we should update
				newStorageFacilityAddress, storageFacilityUpdateErr := f.addressUpdater.UpdateAddress(txnAppCtx, &newShipment.StorageFacility.Address, etag.GenerateEtag(dbShipment.StorageFacility.Address.UpdatedAt))
				if storageFacilityUpdateErr != nil {
					return storageFacilityUpdateErr
				}
				// Assign updated storage facility address to the updated shipment
				newShipment.StorageFacility.AddressID = newStorageFacilityAddress.ID
				newShipment.StorageFacility.Address = *newStorageFacilityAddress
			} else {
				// Make sure that the new storage facility address doesn't already have an ID.
				// If it does, we just assign it. Otherwise, we need to create it.
				if newShipment.StorageFacility.Address.ID != uuid.Nil && newShipment.StorageFacility.AddressID == uuid.Nil {
					// Assign
					newShipment.StorageFacility.AddressID = newShipment.StorageFacility.ID
				} else if newShipment.StorageFacility.Address.ID == uuid.Nil {
					// Create
					newStorageFacilityAddress, storageFacilityCreateErr := f.addressCreator.CreateAddress(txnAppCtx, &newShipment.StorageFacility.Address)
					if storageFacilityCreateErr != nil {
						return storageFacilityCreateErr
					}
					// Assign newly created storage facility address to the updated shipment
					newShipment.StorageFacility.AddressID = newStorageFacilityAddress.ID
				}
			}

			err = txnAppCtx.DB().Save(newShipment.StorageFacility)
			if err != nil {
				return err
			}

			newShipment.StorageFacilityID = &newShipment.StorageFacility.ID

			// For NTS-Release set the pick up address to the storage facility
			if newShipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTS {
				newShipment.PickupAddressID = &newShipment.StorageFacility.AddressID
				newShipment.PickupAddress = &newShipment.StorageFacility.Address

			}
			// For NTS set the destination address to the storage facility
			if newShipment.ShipmentType == models.MTOShipmentTypeHHGIntoNTS {
				newShipment.DestinationAddressID = &newShipment.StorageFacility.AddressID
				newShipment.DestinationAddress = &newShipment.StorageFacility.Address
			}
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
					verrs, err = f.builder.CreateOne(txnAppCtx, &copyOfAgent)
					if verrs != nil && verrs.HasAny() {
						return verrs
					}
					if err != nil {
						return err
					}
				}
			}
		}

		// var move *models.Move
		// If the estimated weight was updated on an approved shipment then it would mean the move could qualify for
		// excess weight risk depending on the weight allowance and other shipment estimated weights
		if newShipment.PrimeEstimatedWeight != nil || (newShipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTS && newShipment.NTSRecordedWeight != nil) {
			// checking if the total of shipment weight & new prime estimated weight is 90% or more of allowed weight
			move, verrs, err = f.moveWeights.CheckExcessWeight(txnAppCtx, dbShipment.MoveTaskOrderID, *newShipment)
			if verrs != nil && verrs.HasAny() {
				return errors.New(verrs.Error())
			}
			if err != nil {
				return err
			}

			// we only want to update the authorized weight if the shipment is approved and the previous weight is nil
			// otherwise, shipment_updater will handle updating authorized weight when a shipment is approved
			if (dbShipment.PrimeEstimatedWeight == nil || (newShipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTS && newShipment.NTSRecordedWeight == nil)) &&
				(newShipment.Status == models.MTOShipmentStatusApproved || newShipment.Status == models.MTOShipmentStatusApprovalsRequested) {
				// updates to prime estimated weight should change the authorized weight of the entitlement
				// which can be manually adjusted by an office user if needed
				err = updateAuthorizedWeight(appCtx, newShipment, move)
				if err != nil {
					return err
				}
			}

			if dbShipment.PrimeEstimatedWeight == nil || *newShipment.PrimeEstimatedWeight != *dbShipment.PrimeEstimatedWeight {
				existingMoveStatus := move.Status
				// if the move is in excess weight risk and the TOO has not acknowledge that, need to change move status to "Approvals Requested"
				// this will trigger the TOO to acknowledged the excess right, which populates ExcessWeightAcknowledgedAt
				if move.ExcessWeightQualifiedAt != nil && move.ExcessWeightAcknowledgedAt == nil || move.ExcessUnaccompaniedBaggageWeightQualifiedAt != nil && move.ExcessUnaccompaniedBaggageWeightAcknowledgedAt == nil {
					err = f.moveRouter.SendToOfficeUser(txnAppCtx, move)
					if err != nil {
						return err
					}
				}

				if existingMoveStatus != move.Status {
					err = txnAppCtx.DB().Update(move)
					if err != nil {
						return err
					}
				}
			}
		}

		if (dbShipment.Status == models.MTOShipmentStatusApproved ||
			dbShipment.Status == models.MTOShipmentStatusApprovalsRequested) &&
			(dbShipment.PrimeEstimatedWeight == nil ||
				*newShipment.PrimeEstimatedWeight != *dbShipment.PrimeEstimatedWeight ||
				(newShipment.PrimeActualWeight != nil && dbShipment.PrimeActualWeight == nil) ||
				(newShipment.PrimeActualWeight != nil && dbShipment.PrimeActualWeight != nil && *newShipment.PrimeActualWeight != *dbShipment.PrimeActualWeight)) {
			err := f.moveWeights.CheckAutoReweigh(txnAppCtx, dbShipment.MoveTaskOrderID, newShipment)
			if err != nil {
				return err
			}
		}

		// Check that only NTS Release shipment uses that NTSRecordedWeight field
		if newShipment.NTSRecordedWeight != nil && newShipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTS {
			errorMsg := fmt.Sprintf("field NTSRecordedWeight cannot be set for shipment type %s", string(newShipment.ShipmentType))
			verrs := validate.NewErrors()
			verrs.Add("NTSRecordedWeight error", errorMsg)
			return apperror.NewInvalidInputError(newShipment.ID, nil, verrs, errorMsg)
		}

		weightsCalculator := NewShipmentBillableWeightCalculator()
		calculatedBillableWeight := weightsCalculator.CalculateShipmentBillableWeight(dbShipment).CalculatedBillableWeight

		// If the max allowable weight for a shipment has been adjusted set a flag to recalculate payment requests for
		// this shipment
		runShipmentRecalculate := false
		if newShipment.BillableWeightCap != nil {
			// new billable cap has a value and it is not the same as the previous value
			if *newShipment.BillableWeightCap != *calculatedBillableWeight {
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

		// when populating the market_code column, it is considered domestic if both pickup & dest are CONUS addresses
		if newShipment.ShipmentType != models.MTOShipmentTypePPM {
			newShipment = models.DetermineShipmentMarketCode(newShipment)
		}

		if newShipment.ShipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
			isShipmentOCONUS := models.IsShipmentOCONUS(*newShipment)
			if isShipmentOCONUS != nil && !*isShipmentOCONUS {
				errorMsg := "At least one address for a UB shipment must be OCONUS"
				ubVerrs := validate.NewErrors()
				ubVerrs.Add("UB shipment error", errorMsg)
				return apperror.NewInvalidInputError(uuid.Nil, nil, ubVerrs, errorMsg)
			}
		}

		// RDD for UB shipments only need the pick up date, shipment origin address and destination address to determine required delivery date
		// Other shipment types have thier RDD set in the checkPrimeValidationsOnModel validator for UpdateMTOShipment which calls this function
		if newShipment.ScheduledPickupDate != nil && !newShipment.ScheduledPickupDate.IsZero() && newShipment.ShipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
			calculatedRDD, err := CalculateRequiredDeliveryDate(appCtx, f.planner, *newShipment.PickupAddress, *newShipment.DestinationAddress, *newShipment.ScheduledPickupDate, nil, newShipment.MoveTaskOrderID, newShipment.ShipmentType)
			if err != nil {
				return err
			}

			newShipment.RequiredDeliveryDate = calculatedRDD
		}

		if err := txnAppCtx.DB().Update(newShipment); err != nil {
			return err
		}

		// if the shipment has an estimated weight, we need to update the service item pricing
		// we only need to do this if the estimated weight, primary addresses, and pickup date are being updated since those all impact pricing
		// we will compare data here to see if we even need to update the pricing
		if newShipment.MarketCode == models.MarketCodeInternational &&
			(newShipment.PrimeEstimatedWeight != nil ||
				newShipment.PickupAddress != nil && dbShipment.PickupAddress != nil && newShipment.PickupAddress.PostalCode != dbShipment.PickupAddress.PostalCode ||
				newShipment.DestinationAddress != nil && dbShipment.DestinationAddress != nil && newShipment.DestinationAddress.PostalCode != dbShipment.DestinationAddress.PostalCode ||
				newShipment.RequestedPickupDate != nil && newShipment.RequestedPickupDate.Format("2006-01-02") != dbShipment.RequestedPickupDate.Format("2006-01-02")) {

			// Recalculate SIT service items using latest mileage.
			// This is to ensure when UpdateEstimatedPricingForShipmentBasicServiceItems
			// is executed it is using most up to date mileage if address changed using service_item.sit_delivery_miles.
			err = UpdateSITServiceItemsSITIfPostalCodeChanged(f.planner, appCtx, f.addressCreator, newShipment)
			if err != nil {
				return err
			}

			portZip, portType, err := models.GetPortLocationInfoForShipment(appCtx.DB(), newShipment.ID)
			if err != nil {
				return err
			}
			// if we don't have the port data, then we won't worry about pricing POEFSC/PODFSC because we need the distance from/to the ports
			if portZip != nil && portType != nil {
				var pickupZip string
				var destZip string
				// if the port type is POEFSC this means the shipment is CONUS -> OCONUS (pickup -> port)
				// if the port type is PODFSC this means the shipment is OCONUS -> CONUS (port -> destination)
				if *portType == models.ReServiceCodePOEFSC.String() {
					pickupZip = newShipment.PickupAddress.PostalCode
					destZip = *portZip
				} else if *portType == models.ReServiceCodePODFSC.String() {
					pickupZip = *portZip
					destZip = newShipment.DestinationAddress.PostalCode
				}
				// we need to get the mileage from DTOD first, the db proc will consume that
				mileage, err := f.planner.ZipTransitDistance(appCtx, pickupZip, destZip)
				if err != nil {
					return err
				}

				// update the service item pricing if relevant fields have changed
				err = models.UpdateEstimatedPricingForShipmentBasicServiceItems(appCtx.DB(), newShipment, &mileage)
				if err != nil {
					return err
				}
			} else {
				// if we don't have the port data, that's okay - we can update the other service items except for PODFSC/POEFSC
				err = models.UpdateEstimatedPricingForShipmentBasicServiceItems(appCtx.DB(), newShipment, nil)
				if err != nil {
					return err
				}
			}
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
		if t, ok := transactionError.(StaleIdentifierError); ok {
			return apperror.NewPreconditionFailedError(dbShipment.ID, t)
		}
		if t, ok := transactionError.(apperror.InvalidInputError); ok {
			return apperror.NewInvalidInputError(dbShipment.ID, t, t.ValidationErrors, "There was an issue with validating the shipment update")
		}
		return apperror.NewQueryError("mtoShipment", transactionError, "")
	}

	return nil

}

type mtoShipmentStatusUpdater struct {
	builder   UpdateMTOShipmentQueryBuilder
	siCreator services.MTOServiceItemCreator
	planner   route.Planner
}

// UpdateMTOShipmentStatus updates MTO Shipment Status
func (o *mtoShipmentStatusUpdater) UpdateMTOShipmentStatus(appCtx appcontext.AppContext, shipmentID uuid.UUID, status models.MTOShipmentStatus, rejectionReason *string, diversionReason *string, eTag string) (*models.MTOShipment, error) {
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
		err = shipmentRouter.RequestDiversion(appCtx, shipment, diversionReason)
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
		var weight *int
		if shipment.PrimeEstimatedWeight != nil {
			weight = models.IntPointer(shipment.PrimeEstimatedWeight.Int())
		}

		switch shipment.ShipmentType {
		case models.MTOShipmentTypeHHGIntoNTS:
			if shipment.StorageFacility == nil || shipment.StorageFacility.AddressID == uuid.Nil {
				return errors.Errorf("StorageFacility is required for %s shipments", models.MTOShipmentTypeHHGIntoNTS)
			}
			err := appCtx.DB().Load(shipment.StorageFacility, "Address", "Address.Country")
			if err != nil {
				return apperror.NewNotFoundError(shipment.StorageFacility.AddressID, "looking for MTOShipment.StorageFacility.Address")
			}

			pickupLocation = shipment.PickupAddress
			deliveryLocation = &shipment.StorageFacility.Address
		case models.MTOShipmentTypeHHGOutOfNTS:
			if shipment.StorageFacility == nil || shipment.StorageFacility.AddressID == uuid.Nil {
				return errors.Errorf("StorageFacility is required for %s shipments", models.MTOShipmentTypeHHGOutOfNTS)
			}
			err := appCtx.DB().Load(shipment.StorageFacility, "Address", "Address.Country")
			if err != nil {
				return apperror.NewNotFoundError(shipment.StorageFacility.AddressID, "looking for MTOShipment.StorageFacility.Address")
			}
			pickupLocation = &shipment.StorageFacility.Address
			deliveryLocation = shipment.DestinationAddress
			if shipment.NTSRecordedWeight != nil {
				weight = models.IntPointer(shipment.NTSRecordedWeight.Int())
			}
		default:
			pickupLocation = shipment.PickupAddress
			deliveryLocation = shipment.DestinationAddress
		}
		requiredDeliveryDate, calcErr := CalculateRequiredDeliveryDate(appCtx, o.planner, *pickupLocation, *deliveryLocation, *shipment.ScheduledPickupDate, weight, shipment.MoveTaskOrderID, shipment.ShipmentType)
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

	// international shipment service items are created in the shipment_approver
	if shipment.MarketCode != models.MarketCodeInternational {
		switch shipment.ShipmentType {
		case models.MTOShipmentTypeHHG:
			originZIP3 := shipment.PickupAddress.PostalCode[0:3]
			destinationZIP3 := shipment.DestinationAddress.PostalCode[0:3]

			if originZIP3 == destinationZIP3 {
				return []models.ReServiceCode{
					models.ReServiceCodeDSH,
					models.ReServiceCodeFSC,
					models.ReServiceCodeDOP,
					models.ReServiceCodeDDP,
					models.ReServiceCodeDPK,
					models.ReServiceCodeDUPK,
				}
			}

			// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Packing, and Dom Unpacking.
			return []models.ReServiceCode{
				models.ReServiceCodeDLH,
				models.ReServiceCodeFSC,
				models.ReServiceCodeDOP,
				models.ReServiceCodeDDP,
				models.ReServiceCodeDPK,
				models.ReServiceCodeDUPK,
			}
		case models.MTOShipmentTypeHHGIntoNTS:
			// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom NTS Packing
			return []models.ReServiceCode{
				models.ReServiceCodeDLH,
				models.ReServiceCodeFSC,
				models.ReServiceCodeDOP,
				models.ReServiceCodeDDP,
				models.ReServiceCodeDNPK,
			}
		case models.MTOShipmentTypeHHGOutOfNTS:
			// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Unpacking
			return []models.ReServiceCode{
				models.ReServiceCodeDLH,
				models.ReServiceCodeFSC,
				models.ReServiceCodeDOP,
				models.ReServiceCodeDDP,
				models.ReServiceCodeDUPK,
			}
		case models.MTOShipmentTypeMobileHome:
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
	}

	return []models.ReServiceCode{}
}

// CalculateRequiredDeliveryDate function is used to get the requied delivery date for a shipment. For shipments that have at least one domestic address
// a distance calculation using the pickup and destination addresses is made. It then uses the returned value and the shipment's weight to make a fetch
// on the ghc_domestic_transit_times table and adds the max_days_transit_time to the pickup date required delivery date. For shipments with with an
// OCONUS address additional days from re_intl_transit_times are added. For UB and shipments with two OCONUS addresses the value from re_intl_transit_times
// is added directly to the pick up date.
func CalculateRequiredDeliveryDate(appCtx appcontext.AppContext, planner route.Planner, pickupAddress models.Address, destinationAddress models.Address, pickupDate time.Time, weight *int, moveID uuid.UUID, shipmentType models.MTOShipmentType) (*time.Time, error) {
	var requiredDeliveryDate time.Time

	destinationIsAlaska, err := destinationAddress.IsAddressAlaska()
	if err != nil {
		return nil, fmt.Errorf("destination address is nil for move ID: %s", moveID)
	}
	pickupIsAlaska, err := pickupAddress.IsAddressAlaska()
	if err != nil {
		return nil, fmt.Errorf("pickup address is nil for move ID: %s", moveID)
	}

	// If either address is CONUS get the domestic transit time based on weight and distance
	if (!destinationIsAlaska || !pickupIsAlaska) && shipmentType != models.MTOShipmentTypeUnaccompaniedBaggage {
		if weight == nil {
			return nil, fmt.Errorf("unable to calculate domestic transit time due to missing weight for move ID: %s", moveID)
		}

		distance, err := planner.ZipTransitDistance(appCtx, pickupAddress.PostalCode, destinationAddress.PostalCode)
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
			return nil, errors.Errorf("failed to find transit time for shipment of %d lbs weight and %d mile distance", *weight, distance)
		}

		// Add the max transit time to the pickup date to get the new required delivery date
		requiredDeliveryDate = pickupDate.AddDate(0, 0, ghcDomesticTransitTime.MaxDaysTransitTime)
	}

	if destinationIsAlaska || pickupIsAlaska {
		var intlTransTime models.InternationalTransitTime

		contract, err := models.FetchContractForMove(appCtx, moveID)
		if err != nil {
			return nil, fmt.Errorf("error fetching contract for move ID: %s", moveID)
		}

		pickupAddressRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), pickupAddress.ID, &uuid.Nil, contract.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching pickup rate area id for address ID: %s", pickupAddress.ID)
		}

		destinationAddressRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), destinationAddress.ID, &uuid.Nil, contract.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching destination rate area id for address ID: %s", destinationAddress.ID)
		}

		intlTransTime, err = models.FetchInternationalTransitTime(appCtx.DB(), pickupAddressRateAreaID, destinationAddressRateAreaID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return nil, fmt.Errorf("no international transit time found for pickup rate area ID: %s and destination rate area ID: %s", pickupAddressRateAreaID, destinationAddressRateAreaID)
			default:
				return nil, err
			}
		}

		if shipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
			if intlTransTime.UbTransitTime != nil {
				// For UBs we add the transit time directly to the pick date
				requiredDeliveryDate = pickupDate.AddDate(0, 0, *intlTransTime.UbTransitTime)
			}
		} else if destinationIsAlaska && pickupIsAlaska {
			if intlTransTime.HhgTransitTime != nil {
				// For AK/AK HHGs we add the transit time directly to the pick date
				requiredDeliveryDate = pickupDate.AddDate(0, 0, *intlTransTime.HhgTransitTime)
			}
		} else {
			if intlTransTime.HhgTransitTime != nil {
				// For CONUS/AK HHGs we add the DB AK transit time to the weight and distance based domestic transit time from above
				requiredDeliveryDate = requiredDeliveryDate.AddDate(0, 0, *intlTransTime.HhgTransitTime)
			}
		}
	}

	return &requiredDeliveryDate, nil
}

// This private function is used to generically construct service items when shipments are approved.
func constructMTOServiceItemModels(shipmentID uuid.UUID, mtoID uuid.UUID, reServiceCodes []models.ReServiceCode) models.MTOServiceItems {
	serviceItems := make(models.MTOServiceItems, len(reServiceCodes))
	currentTime := models.TimePointer(time.Now())
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

// UpdateDestinationSITServiceItemsAddress updates destination SIT service items attached to a shipment
// this updates the final_destination_address to be the same as the shipment's destination_address
func UpdateDestinationSITServiceItemsAddress(appCtx appcontext.AppContext, shipment *models.MTOShipment) error {
	// getting the shipment and service items with code in case they weren't passed in
	eagerAssociations := []string{"MTOServiceItems.ReService.Code"}
	mtoShipment, err := FindShipment(appCtx, shipment.ID, eagerAssociations...)
	if err != nil {
		return err
	}

	mtoServiceItems := mtoShipment.MTOServiceItems

	// Only update these serviceItems address ID
	serviceItemsToUpdate := []models.ReServiceCode{models.ReServiceCodeDDDSIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDASIT, models.ReServiceCodeDDSFSC,
		models.ReServiceCodeIDDSIT, models.ReServiceCodeIDFSIT, models.ReServiceCodeIDASIT, models.ReServiceCodeIDSFSC}

	for _, serviceItem := range mtoServiceItems {

		// Only update the address ID if it is not up to date with the shipment destination address ID
		if slices.Contains(serviceItemsToUpdate, serviceItem.ReService.Code) && serviceItem.SITDestinationFinalAddressID != shipment.DestinationAddressID {

			newServiceItem := serviceItem
			newServiceItem.SITDestinationFinalAddressID = shipment.DestinationAddressID

			transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
				// update service item final destination address ID to match shipment address ID
				verrs, err := txnCtx.DB().ValidateAndUpdate(&newServiceItem)
				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(shipment.ID, err, verrs, "invalid input found while updating final destination address of service item")
				} else if err != nil {
					return apperror.NewQueryError("Service item", err, "")
				}

				return nil
			})

			if transactionError != nil {
				return transactionError
			}
		}
	}

	return nil
}

func UpdateDestinationSITServiceItemsSITDeliveryMiles(planner route.Planner, appCtx appcontext.AppContext, shipment *models.MTOShipment, newAddress *models.Address, TOOApprovalRequired bool) error {
	eagerAssociations := []string{"MTOServiceItems.ReService.Code", "MTOServiceItems.SITDestinationOriginalAddress"}
	mtoShipment, err := FindShipment(appCtx, shipment.ID, eagerAssociations...)
	if err != nil {
		return err
	}

	mtoServiceItems := mtoShipment.MTOServiceItems
	for _, s := range mtoServiceItems {
		serviceItem := s
		reServiceCode := serviceItem.ReService.Code
		if reServiceCode == models.ReServiceCodeDDDSIT ||
			reServiceCode == models.ReServiceCodeDDSFSC ||
			reServiceCode == models.ReServiceCodeIDDSIT ||
			reServiceCode == models.ReServiceCodeIDSFSC {

			var milesCalculated int

			if TOOApprovalRequired {
				if serviceItem.SITDestinationOriginalAddress != nil {
					// if TOO approval was required, shipment destination address has been updated at this point
					milesCalculated, err = planner.ZipTransitDistance(appCtx, shipment.DestinationAddress.PostalCode, serviceItem.SITDestinationOriginalAddress.PostalCode)
				}
			} else {
				// if TOO approval was not required, use the newAddress
				milesCalculated, err = planner.ZipTransitDistance(appCtx, newAddress.PostalCode, serviceItem.SITDestinationOriginalAddress.PostalCode)
			}
			if err != nil {
				return err
			}

			serviceItem.SITDeliveryMiles = &milesCalculated

			mtoServiceItems = append(mtoServiceItems, serviceItem)
		}
	}
	transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		// update service item final SITDeliveryMiles
		verrs, err := txnCtx.DB().ValidateAndUpdate(&mtoServiceItems)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(shipment.ID, err, verrs, "invalid input found while updating final destination address of service item")
		} else if err != nil {
			return apperror.NewQueryError("Service item", err, "")
		}

		return nil
	})

	if transactionError != nil {
		return transactionError
	}

	return nil
}

func UpdateSITServiceItemsSITIfPostalCodeChanged(planner route.Planner, appCtx appcontext.AppContext, addressCreator services.AddressCreator, newShipment *models.MTOShipment) error {
	var expectedSITs = []models.ReServiceCode{models.ReServiceCodeDOPSIT, models.ReServiceCodeDOSFSC,
		models.ReServiceCodeIOPSIT, models.ReServiceCodeIOSFSC, models.ReServiceCodeDDDSIT,
		models.ReServiceCodeDDSFSC, models.ReServiceCodeIDDSIT, models.ReServiceCodeIDSFSC}

	containsSIT := false
	for _, serviceItem := range newShipment.MTOServiceItems {
		if slices.Contains(expectedSITs, serviceItem.ReService.Code) {
			containsSIT = true
		}
	}

	if !containsSIT {
		return nil
	}

	eagerAssociations := []string{"DestinationAddress", "MTOServiceItems.ReService.Code", "MTOServiceItems.SITOriginHHGActualAddress", "MTOServiceItems.SITDestinationFinalAddress", "MTOServiceItems.SITDestinationOriginalAddress"}
	mtoShipment, err := FindShipment(appCtx, newShipment.ID, eagerAssociations...)
	if err != nil {
		return err
	}

	mtoServiceItems := mtoShipment.MTOServiceItems
	for _, s := range mtoServiceItems {
		serviceItem := s
		reServiceCode := serviceItem.ReService.Code
		var milesCalculated int

		if reServiceCode == models.ReServiceCodeDOPSIT ||
			reServiceCode == models.ReServiceCodeDOSFSC ||
			reServiceCode == models.ReServiceCodeIOPSIT ||
			reServiceCode == models.ReServiceCodeIOSFSC {

			milesCalculated, err = planner.ZipTransitDistance(appCtx, serviceItem.SITOriginHHGActualAddress.PostalCode, newShipment.PickupAddress.PostalCode)
			if err != nil {
				return err
			}
			serviceItem.SITDeliveryMiles = &milesCalculated
			mtoServiceItems = append(mtoServiceItems, serviceItem)
		}

		if reServiceCode == models.ReServiceCodeDDDSIT ||
			reServiceCode == models.ReServiceCodeDDSFSC ||
			reServiceCode == models.ReServiceCodeIDDSIT ||
			reServiceCode == models.ReServiceCodeIDSFSC {

			// init using shipment destination if SITDestinationOriginalAddress is not set during pre-approval
			originalDestination := mtoShipment.DestinationAddress.PostalCode
			if serviceItem.SITDestinationOriginalAddress != nil {
				originalDestination = serviceItem.SITDestinationOriginalAddress.PostalCode
			}
			milesCalculated, err = planner.ZipTransitDistance(appCtx, originalDestination, serviceItem.SITDestinationFinalAddress.PostalCode)
			if err != nil {
				return err
			}
			serviceItem.SITDeliveryMiles = &milesCalculated
			mtoServiceItems = append(mtoServiceItems, serviceItem)

		}
	}
	transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := txnCtx.DB().ValidateAndUpdate(&mtoServiceItems)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(newShipment.ID, err, verrs, "invalid input found while updating final destination address of service item")
		} else if err != nil {
			return apperror.NewQueryError("Service item", err, "")
		}

		return nil
	})

	if transactionError != nil {
		return transactionError
	}

	return nil
}

func updateAuthorizedWeight(appCtx appcontext.AppContext, shipment *models.MTOShipment, move *models.Move) error {
	var dBAuthorizedWeight int
	if shipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTS {
		dBAuthorizedWeight = int(*shipment.PrimeEstimatedWeight)
	} else {
		dBAuthorizedWeight = int(*shipment.NTSRecordedWeight)
	}
	if len(move.MTOShipments) != 0 {
		for _, mtoShipment := range move.MTOShipments {
			if (mtoShipment.Status == models.MTOShipmentStatusApproved || mtoShipment.Status == models.MTOShipmentStatusApprovalsRequested) && mtoShipment.ID != shipment.ID {
				if mtoShipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTS {
					//uses PrimeEstimatedWeight for HHG and NTS shipments
					if mtoShipment.PrimeEstimatedWeight != nil {
						dBAuthorizedWeight += int(*mtoShipment.PrimeEstimatedWeight)
					}
				} else {
					//used NTSRecordedWeight for NTSRShipments
					if mtoShipment.NTSRecordedWeight != nil {
						dBAuthorizedWeight += int(*mtoShipment.NTSRecordedWeight)
					}
				}
			}
		}
	}
	dBAuthorizedWeight = int(math.Round(float64(dBAuthorizedWeight) * 1.10))
	entitlement := move.Orders.Entitlement
	entitlement.DBAuthorizedWeight = &dBAuthorizedWeight
	verrs, err := appCtx.DB().ValidateAndUpdate(entitlement)

	if verrs != nil && verrs.HasAny() {
		invalidInputError := apperror.NewInvalidInputError(shipment.ID, nil, verrs, "There was an issue with validating the updates")
		return invalidInputError
	}
	if err != nil {
		return err
	}

	return nil
}
