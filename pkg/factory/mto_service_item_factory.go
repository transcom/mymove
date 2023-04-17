package factory

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type mtoServiceItemBuildType byte

const (
	mtoServiceItemBuildBasic mtoServiceItemBuildType = iota
	mtoServiceItemBuildExtended
)

// buildMTOServiceItemWithBuildType creates a single MTOServiceItem.
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func buildMTOServiceItemWithBuildType(db *pop.Connection, customs []Customization, traits []Trait, buildType mtoServiceItemBuildType) models.MTOServiceItem {
	customs = setupCustomizations(customs, traits)

	// Find address customization and extract the custom address
	var cMTOServiceItem models.MTOServiceItem
	if result := findValidCustomization(customs, MTOServiceItem); result != nil {
		cMTOServiceItem = result.Model.(models.MTOServiceItem)
		if result.LinkOnly {
			return cMTOServiceItem
		}
	}

	var mtoShipmentID *uuid.UUID
	var mtoShipment models.MTOShipment
	var move models.Move
	if buildType == mtoServiceItemBuildExtended {
		// BuildMTOShipment creates a move as necessary
		mtoShipment = BuildMTOShipment(db, customs, traits)
		mtoShipmentID = &mtoShipment.ID
		move = mtoShipment.MoveTaskOrder
	} else {
		move = BuildMove(db, customs, traits)
	}

	var reService models.ReService
	if result := findValidCustomization(customs, ReService); result != nil {
		cReService := result.Model.(models.ReService)
		reService = FetchOrBuildReService(db, cReService)
	} else {
		reService = FetchOrBuildReServiceByCode(db, models.ReServiceCode("STEST"))
	}

	// Create default MTOServiceItem
	mtoServiceItem := models.MTOServiceItem{
		MoveTaskOrder:   move,
		MoveTaskOrderID: move.ID,
		MTOShipment:     mtoShipment,
		MTOShipmentID:   mtoShipmentID,
		ReService:       reService,
		ReServiceID:     reService.ID,
		Status:          models.MTOServiceItemStatusSubmitted,
	}

	// only set SITOriginHHGOriginalAddress if a customization is provided
	if result := findValidCustomization(customs, Addresses.SITOriginHHGOriginalAddress); result != nil {
		addressCustoms := convertCustomizationInList(customs, Addresses.SITOriginHHGOriginalAddress, Address)
		address := BuildAddress(db, addressCustoms, traits)
		mtoServiceItem.SITOriginHHGOriginalAddress = &address
		mtoServiceItem.SITOriginHHGOriginalAddressID = &address.ID
	}

	// only set SITOriginHHGActualAddress if a customization is provided
	if result := findValidCustomization(customs, Addresses.SITOriginHHGActualAddress); result != nil {
		addressCustoms := convertCustomizationInList(customs, Addresses.SITOriginHHGActualAddress, Address)
		address := BuildAddress(db, addressCustoms, traits)
		mtoServiceItem.SITOriginHHGActualAddress = &address
		mtoServiceItem.SITOriginHHGActualAddressID = &address.ID
	}

	// only set SITDestinationFinalAddress if a customization is provided
	if result := findValidCustomization(customs, Addresses.SITDestinationFinalAddress); result != nil {
		addressCustoms := convertCustomizationInList(customs, Addresses.SITDestinationFinalAddress, Address)
		address := BuildAddress(db, addressCustoms, traits)
		mtoServiceItem.SITDestinationFinalAddress = &address
		mtoServiceItem.SITDestinationFinalAddressID = &address.ID
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&mtoServiceItem, cMTOServiceItem)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &mtoServiceItem)
	}

	return mtoServiceItem
}

// BuildMTOServiceItem creates a single extended MTOServiceItem
func BuildMTOServiceItem(db *pop.Connection, customs []Customization, traits []Trait) models.MTOServiceItem {
	return buildMTOServiceItemWithBuildType(db, customs, traits, mtoServiceItemBuildExtended)
}

// BuildMTOServiceItemBasic creates a single basic MTOServiceItem
func BuildMTOServiceItemBasic(db *pop.Connection, customs []Customization, traits []Trait) models.MTOServiceItem {
	return buildMTOServiceItemWithBuildType(db, customs, traits, mtoServiceItemBuildBasic)
}
