package mtoserviceitem

import (
	"fmt"
	"strconv"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services/query"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type createMTOServiceItemQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	CreateOne(model interface{}) (*validate.Errors, error)
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
	Transaction(fn func(tx *pop.Connection) error) error
}

type mtoServiceItemCreator struct {
	builder          createMTOServiceItemQueryBuilder
	createNewBuilder func(db *pop.Connection) createMTOServiceItemQueryBuilder
}

// CreateMTOServiceItem creates a MTO Service Item
func (o *mtoServiceItemCreator) CreateMTOServiceItem(serviceItem *models.MTOServiceItem) (*models.MTOServiceItems, *validate.Errors, error) {
	var verrs *validate.Errors
	var err error
	var requestedServiceItems models.MTOServiceItems // used in case additional service items need to be auto-created
	var createdServiceItems models.MTOServiceItems

	var move models.Move
	moveID := serviceItem.MoveTaskOrderID
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", moveID),
	}
	// check if Move exists
	err = o.builder.FetchOne(&move, queryFilters)
	if err != nil {
		return nil, nil, services.NewNotFoundError(moveID, "in Moves")
	}

	// find the re service code id
	var reService models.ReService
	reServiceCode := serviceItem.ReService.Code
	queryFilters = []services.QueryFilter{
		query.NewQueryFilter("code", "=", reServiceCode),
	}
	err = o.builder.FetchOne(&reService, queryFilters)
	if err != nil {
		return nil, nil, services.NewNotFoundError(uuid.Nil, fmt.Sprintf("for service item with code: %s", reServiceCode))
	}
	// set re service for service item
	serviceItem.ReServiceID = reService.ID
	serviceItem.Status = models.MTOServiceItemStatusSubmitted

	// We can have two service items that come in from a MTO approval that do not have an MTOShipmentID
	// they are MTO level service items. This should capture that and create them accordingly, they are thankfully
	// also rather basic.
	if serviceItem.MTOShipmentID == nil {
		if serviceItem.ReService.Code == models.ReServiceCodeMS || serviceItem.ReService.Code == models.ReServiceCodeCS {
			serviceItem.Status = "APPROVED"
		}
		verrs, err = o.builder.CreateOne(serviceItem)
		if verrs != nil {
			return nil, verrs, nil
		}
		if err != nil {
			return nil, nil, err
		}

		createdServiceItems = append(createdServiceItems, *serviceItem)

		return &createdServiceItems, nil, nil
	}

	// TODO: Once customer onboarding is built, we can revisit to figure out which service items goes under each type of shipment
	// check if shipment exists linked by MoveTaskOrderID
	var mtoShipment models.MTOShipment
	var mtoShipmentID uuid.UUID

	mtoShipmentID = *serviceItem.MTOShipmentID
	queryFilters = []services.QueryFilter{
		query.NewQueryFilter("id", "=", mtoShipmentID),
		query.NewQueryFilter("move_id", "=", moveID),
	}
	err = o.builder.FetchOne(&mtoShipment, queryFilters)
	if err != nil {
		return nil, nil, services.NewNotFoundError(mtoShipmentID,
			fmt.Sprintf("for mtoShipment with moveID: %s", moveID.String()))
	}

	if serviceItem.ReService.Code == models.ReServiceCodeDOSHUT || serviceItem.ReService.Code == models.ReServiceCodeDDSHUT {
		if mtoShipment.PrimeEstimatedWeight == nil {
			return nil, verrs, services.NewConflictError(reService.ID, "for creating a service item. MTOShipment associated with this service item must have a valid PrimeEstimatedWeight.")
		}
	}

	for index := range serviceItem.CustomerContacts {
		createCustContacts := &serviceItem.CustomerContacts[index]
		err = validateTimeMilitaryField(createCustContacts.TimeMilitary)
		if err != nil {
			return nil, nil, services.NewInvalidInputError(serviceItem.ID, err, nil, err.Error())
		}
	}

	if serviceItem.ReService.Code == models.ReServiceCodeDDFSIT {
		// check if there's another DDFSIT item for this shipment
		err = o.checkDuplicateServiceCodes(serviceItem)
		if err != nil {
			return nil, nil, err
		}

		// create extra DDASIT and DDDSIT service items
		serviceItemDDASIT, errSIT := o.makeExtraSITServiceItem(serviceItem, models.ReServiceCodeDDASIT)
		if errSIT != nil {
			return nil, nil, errSIT
		}
		if serviceItemDDASIT != nil {
			requestedServiceItems = append(requestedServiceItems, *serviceItemDDASIT)
		}

		serviceItemDDDSIT, errSIT := o.makeExtraSITServiceItem(serviceItem, models.ReServiceCodeDDDSIT)
		if errSIT != nil {
			return nil, nil, errSIT
		}
		if serviceItemDDDSIT != nil {
			requestedServiceItems = append(requestedServiceItems, *serviceItemDDDSIT)
		}
	}

	requestedServiceItems = append(requestedServiceItems, *serviceItem)

	// create new items in a transaction in case of failure
	o.builder.Transaction(func(tx *pop.Connection) error {
		// create new builder to use tx
		txBuilder := o.createNewBuilder(tx)

		for serviceItemIndex := range requestedServiceItems {
			requestedServiceItem := &requestedServiceItems[serviceItemIndex]
			verrs, err = txBuilder.CreateOne(requestedServiceItem)
			if verrs != nil || err != nil {
				return fmt.Errorf("%#v %e", verrs, err)
			}

			createdServiceItems = append(createdServiceItems, *requestedServiceItem)

			// create dimensions if any
			for index := range requestedServiceItem.Dimensions {
				createDimension := &requestedServiceItem.Dimensions[index]
				createDimension.MTOServiceItemID = requestedServiceItem.ID
				verrs, err = txBuilder.CreateOne(createDimension)
				if verrs != nil || err != nil {
					return fmt.Errorf("%#v %e", verrs, err)
				}
			}

			// create customer contacts if any
			for index := range requestedServiceItem.CustomerContacts {
				createCustContacts := &requestedServiceItem.CustomerContacts[index]
				createCustContacts.MTOServiceItemID = requestedServiceItem.ID
				verrs, err = txBuilder.CreateOne(createCustContacts)
				if verrs != nil || err != nil {
					return fmt.Errorf("%#v %e", verrs, err)
				}
			}
		}

		return nil
	})

	if verrs != nil && verrs.HasAny() {
		return nil, verrs, nil
	} else if err != nil {
		return nil, verrs, services.NewQueryError("unknown", err, "")
	}

	if move.Status != models.MoveStatusAPPROVALSREQUESTED {
		err := move.SetApprovalsRequested()
		if err != nil {
			return nil, nil, err
		}
		verrs, err := o.builder.UpdateOne(&move, nil)
		if verrs != nil || err != nil {
			return nil, verrs, err
		}
	}

	return &createdServiceItems, nil, nil
}

// checkDuplicateServiceCodes checks if the move or shipment has a duplicate service item with the same code as the one
// requested.
func (o *mtoServiceItemCreator) checkDuplicateServiceCodes(serviceItem *models.MTOServiceItem) error {
	var duplicateServiceItem models.MTOServiceItem

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("move_id", "=", serviceItem.MoveTaskOrderID),
		query.NewQueryFilter("re_service_id", "=", serviceItem.ReServiceID),
	}
	if serviceItem.MTOShipmentID != nil {
		queryFilters = append(queryFilters, query.NewQueryFilter("mto_shipment_id", "=", serviceItem.MTOShipmentID))
	}

	err := o.builder.FetchOne(&duplicateServiceItem, queryFilters)
	if err == nil && duplicateServiceItem.ID != uuid.Nil {
		return services.NewConflictError(duplicateServiceItem.ID,
			fmt.Sprintf("for creating a service item. A service item with reServiceCode %s already exists for this move and/or shipment.", serviceItem.ReService.Code))
	}

	return nil
}

// makeExtraSITServiceItem sets up extra SIT service items if a first-day SIT service item is being created.
func (o *mtoServiceItemCreator) makeExtraSITServiceItem(firstSIT *models.MTOServiceItem, reServiceCode models.ReServiceCode) (*models.MTOServiceItem, error) {
	var reService models.ReService

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("code", "=", reServiceCode),
	}
	err := o.builder.FetchOne(&reService, queryFilters)
	if err != nil {
		return nil, services.NewNotFoundError(uuid.Nil, fmt.Sprintf("for service code: %s", reServiceCode))
	}

	extraServiceItem := models.MTOServiceItem{
		MTOShipmentID:   firstSIT.MTOShipmentID,
		MoveTaskOrderID: firstSIT.MoveTaskOrderID,
		ReServiceID:     reService.ID,
		ReService:       reService,
		SITEntryDate:    firstSIT.SITEntryDate,
		SITPostalCode:   firstSIT.SITPostalCode,
		Reason:          firstSIT.Reason,
		Status:          models.MTOServiceItemStatusSubmitted,
	}

	return &extraServiceItem, nil
}

// NewMTOServiceItemCreator returns a new MTO service item creator
func NewMTOServiceItemCreator(builder createMTOServiceItemQueryBuilder) services.MTOServiceItemCreator {
	// used inside a transaction and mocking
	createNewBuilder := func(db *pop.Connection) createMTOServiceItemQueryBuilder {
		return query.NewQueryBuilder(db)
	}

	return &mtoServiceItemCreator{builder: builder, createNewBuilder: createNewBuilder}
}

func validateTimeMilitaryField(timeMilitary string) error {
	if len(timeMilitary) == 0 {
		return nil
	} else if len(timeMilitary) != 5 {
		return fmt.Errorf("timeMilitary must be in format HHMMZ")
	}

	hours := timeMilitary[:2]
	minutes := timeMilitary[2:4]
	suffix := timeMilitary[len(timeMilitary)-1:]

	hoursInt, err := strconv.Atoi(hours)
	if err != nil {
		return fmt.Errorf("timeMilitary must have a valid number for hours")
	}

	minutesInt, err := strconv.Atoi(minutes)
	if err != nil {
		return fmt.Errorf("timeMilitary must have a valid number for minutes")
	}

	if !(0 <= hoursInt) || !(hoursInt < 24) {
		return fmt.Errorf("timeMilitary hours must be between 00 and 23")
	}
	if !(0 <= minutesInt) || !(minutesInt < 60) {
		return fmt.Errorf("timeMilitary minutes must be between 00 and 59")
	}

	if suffix != "Z" {
		return fmt.Errorf("timeMilitary must end with 'Z'")
	}

	return nil
}
