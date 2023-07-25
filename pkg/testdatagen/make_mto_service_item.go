package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// makeServiceItem creates a single service item and associated set relationships
func makeServiceItem(db *pop.Connection, assertions Assertions, isBasicServiceItem bool) models.MTOServiceItem {
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = makeMove(db, assertions)
	}

	var mtoShipmentID *uuid.UUID
	var mtoShipment models.MTOShipment
	if !isBasicServiceItem {
		if isZeroUUID(assertions.MTOShipment.ID) {
			mtoShipment = makeMTOShipment(db, assertions)
			mtoShipmentID = &mtoShipment.ID
		} else {
			mtoShipment = assertions.MTOShipment
			mtoShipmentID = &assertions.MTOShipment.ID
		}
	}

	reService := assertions.ReService
	if isZeroUUID(reService.ID) {
		reService = FetchOrMakeReService(db, assertions)
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

	return MTOServiceItem
}

// MakeMTOServiceItem creates a single MTOServiceItem and associated set relationships
func MakeMTOServiceItem(db *pop.Connection, assertions Assertions) models.MTOServiceItem {
	return makeServiceItem(db, assertions, false)
}

// MakeDefaultMTOServiceItem returns a MTOServiceItem with default values
func MakeDefaultMTOServiceItem(db *pop.Connection) models.MTOServiceItem {
	return MakeMTOServiceItem(db, Assertions{})
}

// MakeMTOServiceItemDomesticCrating makes a domestic crating service item and its associated item and crate
func MakeMTOServiceItemDomesticCrating(db *pop.Connection, assertions Assertions) models.MTOServiceItem {
	mtoServiceItem := MakeMTOServiceItem(db, assertions)

	// Create item
	dimensionItem := MakeMTOServiceItemDimension(db, Assertions{
		MTOServiceItemDimension: assertions.MTOServiceItemDimension,
		MTOServiceItem:          mtoServiceItem,
	})

	// Create crate
	assertions.MTOServiceItemDimensionCrate.Type = models.DimensionTypeCrate
	crateItem := MakeMTOServiceItemDimension(db, Assertions{
		MTOServiceItemDimension: assertions.MTOServiceItemDimensionCrate,
		MTOServiceItem:          mtoServiceItem,
	})

	mtoServiceItem.Dimensions = append(mtoServiceItem.Dimensions, dimensionItem, crateItem)

	return mtoServiceItem
}
