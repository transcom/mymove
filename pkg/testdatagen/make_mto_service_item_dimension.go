package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMTOServiceItemDimension creates a single MTOServiceItemDimension and associated set relationships
func MakeMTOServiceItemDimension(db *pop.Connection, assertions Assertions) models.MTOServiceItemDimension {
	MTOServiceItem := assertions.MTOServiceItem
	if isZeroUUID(MTOServiceItem.ID) {
		MTOServiceItem = MakeMTOServiceItem(db, assertions)
	}

	MTOServiceItemDimension := models.MTOServiceItemDimension{
		MTOServiceItemID: MTOServiceItem.ID,
		MTOServiceItem:   MTOServiceItem,
		Type:             models.DimensionTypeItem,
		Length:           12000,
		Height:           12000,
		Width:            12000,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	// Overwrite values with those from assertions
	mergeModels(&MTOServiceItemDimension, assertions.MTOServiceItemDimension)

	mustCreate(db, &MTOServiceItemDimension, assertions.Stub)

	return MTOServiceItemDimension
}

// MakeDefaultMTOServiceItemDimension returns a MTOServiceItemDimension with default values
func MakeDefaultMTOServiceItemDimension(db *pop.Connection) models.MTOServiceItemDimension {
	return MakeMTOServiceItemDimension(db, Assertions{})
}
