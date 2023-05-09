package factory

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func BuildMTOServiceItemDimension(db *pop.Connection, customs []Customization, traits []Trait) models.MTOServiceItemDimension {
	customs = setupCustomizations(customs, traits)

	// Find MTOServiceItemDimension customization and extract the custom model
	var cMTOServiceItemDimension models.MTOServiceItemDimension
	if result := findValidCustomization(customs, MTOServiceItemDimension); result != nil {
		cMTOServiceItemDimension = result.Model.(models.MTOServiceItemDimension)
		if result.LinkOnly {
			return cMTOServiceItemDimension
		}
	}

	mtoServiceItem := BuildMTOServiceItem(db, customs, traits)

	// create default MTOServiceItemDimension
	mTOServiceItemDimension := models.MTOServiceItemDimension{
		MTOServiceItem:   mtoServiceItem,
		MTOServiceItemID: mtoServiceItem.ID,
		Type:             models.DimensionTypeItem,
		Length:           12000,
		Height:           12000,
		Width:            12000,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&mTOServiceItemDimension, cMTOServiceItemDimension)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &mTOServiceItemDimension)
	}

	return mTOServiceItemDimension
}
