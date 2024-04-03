package mtoshipment

import (
	"database/sql"
	"fmt"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type createMTOShipmentQueryBuilder interface {
	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
	CreateOne(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error)
	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
}

type mtoShipmentCreator struct {
	builder createMTOShipmentQueryBuilder
	services.Fetcher
	moveRouter services.MoveRouter
	addressCreator services.AddressCreator
	checks     []validator
}

// NewMTOShipmentCreatorV1 creates a new struct with the service dependencies
// This is utilized in Prime API V1
func NewMTOShipmentCreatorV1(builder createMTOShipmentQueryBuilder, fetcher services.Fetcher, moveRouter services.MoveRouter, addressCreator services.AddressCreator) services.MTOShipmentCreator {
	return &mtoShipmentCreator{
		builder,
		fetcher,
		moveRouter,
		addressCreator,
		[]validator{protectV1Diversion()},
	}
}

// NewMTOShipmentCreator creates a new struct with the service dependencies
// This is utilized in Prime API V2
func NewMTOShipmentCreatorV2(builder createMTOShipmentQueryBuilder, fetcher services.Fetcher, moveRouter services.MoveRouter, addressCreator services.AddressCreator) services.MTOShipmentCreator {
	return &mtoShipmentCreator{
		builder,
		fetcher,
		moveRouter,
		addressCreator,
		[]validator{checkDiversionValid(), childDiversionPrimeWeightRule()},
	}
}

// CreateMTOShipment creates the mto shipment
func (f mtoShipmentCreator) CreateMTOShipment(appCtx appcontext.AppContext, shipment *models.MTOShipment) (*models.MTOShipment, error) {
	var verrs *validate.Errors
	var err error

	for _, check := range f.checks {
		if err = check.Validate(appCtx, shipment, nil); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				// Accumulate all validation errors
				verrs.Append(e)
			default:
				// Non-validation errors have priority and short-circuit doing any further checks
				return nil, err
			}
		}
	}

	serviceItems := shipment.MTOServiceItems
	err = checkShipmentIDFields(shipment, serviceItems)

	if err != nil {
		return nil, err
	}

	// Check shipment fields that should be there or not based on shipment type.
	if shipment.ShipmentType == models.MTOShipmentTypeHHGOutOfNTSDom {
		if shipment.RequestedPickupDate != nil {
			return nil, apperror.NewInvalidInputError(uuid.Nil, nil, verrs,
				fmt.Sprintf("RequestedPickupDate should not be set when creating a %s shipment", shipment.ShipmentType))
		}
	} else if shipment.ShipmentType != models.MTOShipmentTypePPM {
		// No need for a PPM to have a RequestedPickupDate
		if shipment.RequestedPickupDate == nil || shipment.RequestedPickupDate.IsZero() {
			return nil, apperror.NewInvalidInputError(uuid.Nil, nil, verrs,
				fmt.Sprintf("RequestedPickupDate is required to create a %s shipment", shipment.ShipmentType))
		}
		if shipment.NTSRecordedWeight != nil {
			return nil, apperror.NewInvalidInputError(uuid.Nil, nil, verrs,
				fmt.Sprintf("NTSRecordedWeight should not be set when creating a %s shipment", shipment.ShipmentType))
		}
	}

	var move models.Move
	moveID := shipment.MoveTaskOrderID

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", moveID),
	}

	// check if Move exists
	err = f.builder.FetchOne(appCtx, &move, queryFilters)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(moveID, "for move")
		default:
			return nil, apperror.NewQueryError("Move", err, "")
		}
	}

	if appCtx.Session() != nil && appCtx.Session().IsMilApp() {
		if move.Orders.ServiceMemberID != appCtx.Session().ServiceMemberID {
			return nil, apperror.NewNotFoundError(appCtx.Session().ServiceMemberID, "for service member")
		}
	}

	if serviceItems != nil {
		serviceItemsList := make(models.MTOServiceItems, 0, len(serviceItems))

		for _, serviceItem := range serviceItems {
			// find the re service code id
			var reService models.ReService
			reServiceCode := serviceItem.ReService.Code
			queryFilters = []services.QueryFilter{
				query.NewQueryFilter("code", "=", reServiceCode),
			}
			err = f.builder.FetchOne(appCtx, &reService, queryFilters)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					return nil, apperror.NewNotFoundError(uuid.Nil, fmt.Sprintf("for service item with code: %s", reServiceCode))
				default:
					return nil, apperror.NewQueryError("ReService", err, "")
				}
			}
			// set re service for service item
			serviceItem.ReServiceID = reService.ID
			serviceItem.Status = models.MTOServiceItemStatusSubmitted

			serviceItemsList = append(serviceItemsList, serviceItem)
		}
		shipment.MTOServiceItems = serviceItemsList
	}

	// Populate the destination address fields with the new duty location's address when
	// we have an HHG with no destination address, but don't copy over any street fields.
	if shipment.ShipmentType == models.MTOShipmentTypeHHG && shipment.DestinationAddress == nil {
		err = appCtx.DB().Load(&move, "Orders.NewDutyLocation.Address")
		if err != nil {
			return nil, apperror.NewQueryError("Orders", err, "")
		}
		newDutyLocationAddress := move.Orders.NewDutyLocation.Address

		county, errCounty := models.FindCountyByZipCode(appCtx.DB(), newDutyLocationAddress.PostalCode)
		if errCounty != nil {
			return nil, errCounty
		}

		shipment.DestinationAddress = &models.Address{
			StreetAddress1: "N/A", // can't use an empty string given the model validations
			City:           newDutyLocationAddress.City,
			State:          newDutyLocationAddress.State,
			PostalCode:     newDutyLocationAddress.PostalCode,
			Country:        newDutyLocationAddress.Country,
			County:         &county,
		}
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		// create pickup and destination addresses
		if shipment.PickupAddress != nil {
			pickupAddress, err := f.addressCreator.CreateAddress(txnAppCtx, shipment.PickupAddress)
			if err != nil {
				return fmt.Errorf("failed to create pickup address %#v %e", verrs, err)
			}
			shipment.PickupAddress = pickupAddress
			shipment.PickupAddressID = &shipment.PickupAddress.ID

		} else if shipment.ShipmentType != models.MTOShipmentTypeHHGOutOfNTSDom && shipment.ShipmentType != models.MTOShipmentTypePPM {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "PickupAddress is required to create an HHG or NTS type MTO shipment")
		}

		if shipment.SecondaryPickupAddress != nil {
			secondaryPickupAddress, err := f.addressCreator.CreateAddress(txnAppCtx, shipment.SecondaryPickupAddress)
			if err != nil {
				return fmt.Errorf("failed to create secondary pickup address %#v %e", verrs, err)
			}
			shipment.SecondaryPickupAddress = secondaryPickupAddress
			shipment.SecondaryPickupAddressID = &shipment.SecondaryPickupAddress.ID
		}

		if shipment.DestinationAddress != nil {
			destinationAddress, err := f.addressCreator.CreateAddress(txnAppCtx, shipment.DestinationAddress)
			if err != nil {
				return fmt.Errorf("failed to create destination address %#v %e", verrs, err)
			}
			shipment.DestinationAddress = destinationAddress
			shipment.DestinationAddressID = &shipment.DestinationAddress.ID
		}

		if shipment.SecondaryDeliveryAddress != nil {
			secondaryDeliveryAddress, err := f.addressCreator.CreateAddress(txnAppCtx, shipment.SecondaryDeliveryAddress)
			if err != nil {
				return fmt.Errorf("failed to create secondary delivery address %#v %e", verrs, err)
			}
			shipment.SecondaryDeliveryAddress = secondaryDeliveryAddress
			shipment.SecondaryDeliveryAddressID = &shipment.SecondaryDeliveryAddress.ID
		}

		if shipment.StorageFacility != nil {
			storageFacility, err := f.addressCreator.CreateAddress(txnAppCtx, &shipment.StorageFacility.Address)
			if err != nil {
				return fmt.Errorf("failed to create storage facility address %#v %e", verrs, err)
			}
			shipment.StorageFacility.Address = *storageFacility
			shipment.StorageFacility.AddressID = shipment.StorageFacility.Address.ID

			verrs, err = f.builder.CreateOne(txnAppCtx, shipment.StorageFacility)
			if verrs != nil || err != nil {
				return fmt.Errorf("failed to create storage facility %#v %e", verrs, err)
			}
			shipment.StorageFacilityID = &shipment.StorageFacility.ID
		}

		//assign status to shipment draft by default
		if shipment.Status != models.MTOShipmentStatusSubmitted {
			shipment.Status = models.MTOShipmentStatusDraft
		}

		// Assign default SITDaysAllowance based on customer type...but we only have service members right now.
		// Once we introduce more, this logic will have to change.
		defaultSITDays := int(models.DefaultServiceMemberSITDaysAllowance)
		shipment.SITDaysAllowance = &defaultSITDays

		// create a shipment
		verrs, err = f.builder.CreateOne(txnAppCtx, shipment)

		if verrs != nil || err != nil {
			return fmt.Errorf("failed to create shipment %s %e", verrs.Error(), err)
		}

		// create MTOAgents List
		if shipment.MTOAgents != nil {
			agentsList := make(models.MTOAgents, 0, len(shipment.MTOAgents))

			for _, agent := range shipment.MTOAgents {
				copyOfAgent := agent
				copyOfAgent.MTOShipmentID = shipment.ID
				if copyOfAgent.FirstName != nil && *copyOfAgent.FirstName == "" {
					copyOfAgent.FirstName = nil
				}
				if copyOfAgent.LastName != nil && *copyOfAgent.LastName == "" {
					copyOfAgent.LastName = nil
				}
				if copyOfAgent.Email != nil && *copyOfAgent.Email == "" {
					copyOfAgent.Email = nil
				}
				if copyOfAgent.Phone != nil && *copyOfAgent.Phone == "" {
					copyOfAgent.Phone = nil
				}
				// If no fields are set, then we do not want to create the MTO agent
				if copyOfAgent.FirstName == nil && copyOfAgent.LastName == nil && copyOfAgent.Email == nil && copyOfAgent.Phone == nil {
					continue
				}
				verrs, err = f.builder.CreateOne(txnAppCtx, &copyOfAgent)
				if verrs != nil && verrs.HasAny() {
					return verrs
				}
				if err != nil {
					return err
				}

				for _, agentInList := range agentsList {
					if agentInList.MTOAgentType == copyOfAgent.MTOAgentType {
						return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "MTOAgents can only contain one agent of each type")
					}
				}

				agentsList = append(agentsList, copyOfAgent)
			}
			shipment.MTOAgents = agentsList
		}

		// create MTOServiceItems List
		if shipment.MTOServiceItems != nil {
			serviceItemsList := make(models.MTOServiceItems, 0, len(shipment.MTOServiceItems))

			for _, serviceItem := range shipment.MTOServiceItems {
				copyOfServiceItem := serviceItem // Make copy to avoid implicit memory aliasing of items from a range statement.
				copyOfServiceItem.MTOShipmentID = &shipment.ID
				copyOfServiceItem.MoveTaskOrderID = shipment.MoveTaskOrderID

				verrs, err = f.builder.CreateOne(txnAppCtx, &copyOfServiceItem)
				if verrs != nil && verrs.HasAny() {
					return verrs
				}
				if err != nil {
					return err
				}
				serviceItemsList = append(serviceItemsList, copyOfServiceItem)
			}
			shipment.MTOServiceItems = serviceItemsList
		}

		// transition the move to "Approvals Requested" if a shipment was created with the "Submitted" status:
		if shipment.Status == models.MTOShipmentStatusSubmitted && move.Status == models.MoveStatusAPPROVED {
			err = f.moveRouter.SendToOfficeUser(txnAppCtx, &move)
			if err != nil {
				return err
			}
			verrs, err = f.builder.UpdateOne(txnAppCtx, &move, nil)
			if err != nil {
				return err
			}
			if verrs != nil && verrs.HasAny() {
				return verrs
			}
		}

		return nil
	})

	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Unable to create shipment")
	} else if err != nil {
		return nil, apperror.NewQueryError("unknown", err, "")
	}

	return shipment, transactionError
}

// checkShipmentIDFields checks that the client hasn't attempted to set ID fields that will be generated/auto-set
func checkShipmentIDFields(shipment *models.MTOShipment, serviceItems models.MTOServiceItems) error {
	verrs := validate.NewErrors()

	if shipment.MTOAgents != nil && len(shipment.MTOAgents) > 0 {
		for _, agent := range shipment.MTOAgents {
			if agent.ID != uuid.Nil {
				verrs.Add("agents:id", "cannot be set for new agents")
			}
			if agent.MTOShipmentID != uuid.Nil {
				verrs.Add("agents:mtoShipmentID", "cannot be set for agents created with a shipment")
			}
		}
	}

	if len(serviceItems) > 0 {
		for _, item := range serviceItems {
			if item.ID != uuid.Nil {
				verrs.Add("mtoServiceItems:id", "cannot be set for new service items")
			}
			if item.MTOShipmentID != nil && *item.MTOShipmentID != uuid.Nil {
				verrs.Add("mtoServiceItems:mtoShipmentID", "cannot be set for service items created with a shipment")
			}
		}
	}

	addressMsg := "cannot be set for new addresses"
	if shipment.PickupAddress != nil && shipment.PickupAddress.ID != uuid.Nil {
		verrs.Add("pickupAddress:id", addressMsg)
	}
	if shipment.DestinationAddress != nil && shipment.DestinationAddress.ID != uuid.Nil {
		verrs.Add("destinationAddress:id", addressMsg)
	}
	if shipment.SecondaryPickupAddress != nil && shipment.SecondaryPickupAddress.ID != uuid.Nil {
		verrs.Add("secondaryPickupAddress:id", addressMsg)
	}
	if shipment.SecondaryDeliveryAddress != nil && shipment.SecondaryDeliveryAddress.ID != uuid.Nil {
		verrs.Add("SecondaryDeliveryAddress:id", addressMsg)
	}

	if verrs.HasAny() {
		return apperror.NewInvalidInputError(uuid.Nil, nil, verrs, "Fields that cannot be set found while creating new shipment.")
	}

	return nil
}
