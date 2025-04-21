package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// makeServiceItem creates a single service item and associated set relationships
func makeServiceItem(db *pop.Connection, assertions Assertions, isBasicServiceItem bool) (models.MTOServiceItem, error) {
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		var err error
		moveTaskOrder, err = makeMove(db, assertions)
		if err != nil {
			return models.MTOServiceItem{}, err
		}
	}

	var mtoShipmentID *uuid.UUID
	var mtoShipment models.MTOShipment
	if !isBasicServiceItem {
		if isZeroUUID(assertions.MTOShipment.ID) {
			var err error
			mtoShipment, err = makeMTOShipment(db, assertions)
			if err != nil {
				return models.MTOServiceItem{}, err
			}
			mtoShipmentID = &mtoShipment.ID
		} else {
			mtoShipment = assertions.MTOShipment
			mtoShipmentID = &assertions.MTOShipment.ID
		}
	}

	reService := assertions.ReService
	if isZeroUUID(reService.ID) {
		reService = FetchReService(db, assertions)
	}

	status := assertions.MTOServiceItem.Status
	if status == "" {
		status = models.MTOServiceItemStatusSubmitted
	}

	MTOServiceItem := models.MTOServiceItem{
		MoveTaskOrder:   moveTaskOrder,
		MoveTaskOrderID: moveTaskOrder.ID,
		MTOShipment:     mtoShipment,
		MTOShipmentID:   mtoShipmentID,
		ReService:       reService,
		ReServiceID:     reService.ID,
		Status:          status,
	}

	// Overwrite values with those from assertions
	mergeModels(&MTOServiceItem, assertions.MTOServiceItem)

	mustCreate(db, &MTOServiceItem, assertions.Stub)

	return MTOServiceItem, nil
}

// MakeMTOServiceItem creates a single MTOServiceItem and associated set relationships
func MakeMTOServiceItem(db *pop.Connection, assertions Assertions) (models.MTOServiceItem, error) {
	return makeServiceItem(db, assertions, false)
}

// MakeDefaultMTOServiceItem returns a MTOServiceItem with default values
func MakeDefaultMTOServiceItem(db *pop.Connection) (models.MTOServiceItem, error) {
	return MakeMTOServiceItem(db, Assertions{})
}

// MakeMTOServiceItem creates a single MTOServiceItem and associated set relationships
func MakeStubbedMTOServiceItem(db *pop.Connection) (models.MTOServiceItem, error) {
	return makeServiceItem(db, Assertions{
		Stub: true,
	}, false)
}

// MakeMTOServiceItemBasic creates a single MTOServiceItem that is a basic type, meaning no shipment id associated.
func MakeMTOServiceItemBasic(db *pop.Connection, assertions Assertions) (models.MTOServiceItem, error) {
	return makeServiceItem(db, assertions, true)
}

// MakeMTOServiceItems makes an array of MTOServiceItems
func MakeMTOServiceItems(db *pop.Connection) (models.MTOServiceItems, error) {
	var serviceItemList models.MTOServiceItems
	mtoServiceItem, err := MakeDefaultMTOServiceItem(db)
	if err != nil {
		return models.MTOServiceItems{}, err
	}
	serviceItemList = append(serviceItemList, mtoServiceItem)
	return serviceItemList, nil
}

// MakeMTOServiceItemDomesticCrating makes a domestic crating service item and its associated item and crate
func MakeMTOServiceItemDomesticCrating(db *pop.Connection, assertions Assertions) (models.MTOServiceItem, error) {
	mtoServiceItem, err := MakeMTOServiceItem(db, assertions)
	if err != nil {
		return models.MTOServiceItem{}, err
	}

	// Create item
	dimensionItem, err := MakeMTOServiceItemDimension(db, Assertions{
		MTOServiceItemDimension: assertions.MTOServiceItemDimension,
		MTOServiceItem:          mtoServiceItem,
	})
	if err != nil {
		return models.MTOServiceItem{}, err
	}

	// Create crate
	assertions.MTOServiceItemDimensionCrate.Type = models.DimensionTypeCrate
	crateItem, err := MakeMTOServiceItemDimension(db, Assertions{
		MTOServiceItemDimension: assertions.MTOServiceItemDimensionCrate,
		MTOServiceItem:          mtoServiceItem,
	})
	if err != nil {
		return models.MTOServiceItem{}, err
	}

	mtoServiceItem.Dimensions = append(mtoServiceItem.Dimensions, dimensionItem, crateItem)

	return mtoServiceItem, nil
}
