package testdatagen

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// makeServiceItem creates a single service item and associated set relationships
func makeServiceItem(db *pop.Connection, assertions Assertions, isBasicServiceItem bool) models.MTOServiceItem {
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMove(db, assertions)
	}

	var MTOShipmentID *uuid.UUID
	var MTOShipment models.MTOShipment
	if !isBasicServiceItem {
		if isZeroUUID(assertions.MTOShipment.ID) {
			MTOShipment = MakeMTOShipment(db, assertions)
			MTOShipmentID = &MTOShipment.ID
		} else {
			MTOShipment = assertions.MTOShipment
			MTOShipmentID = &assertions.MTOShipment.ID
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
		MTOShipment:     MTOShipment,
		MTOShipmentID:   MTOShipmentID,
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

// MakeMTOServiceItemBasic creates a single MTOServiceItem that is a basic type, meaning no shipment id associated.
func MakeMTOServiceItemBasic(db *pop.Connection, assertions Assertions) models.MTOServiceItem {
	return makeServiceItem(db, assertions, true)
}

// MakeMTOServiceItems makes an array of MTOServiceItems
func MakeMTOServiceItems(db *pop.Connection) models.MTOServiceItems {
	var serviceItemList models.MTOServiceItems
	serviceItemList = append(serviceItemList, MakeDefaultMTOServiceItem(db))
	return serviceItemList
}
