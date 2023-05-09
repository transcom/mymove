package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeStorageFacility creates a single SIT Extension and associated set relationships
func MakeStorageFacility(db *pop.Connection, assertions Assertions) models.StorageFacility {
	address := assertions.StorageFacility.Address
	if isZeroUUID(address.ID) {
		address = MakeAddress(db, assertions)
	}

	storageFacility := models.StorageFacility{
		FacilityName: "Storage R Us",
		Address:      address,
		AddressID:    address.ID,
		LotNumber:    models.StringPointer("1234"),
		Phone:        models.StringPointer("555-555-5555"),
		Email:        models.StringPointer("storage@email.com"),
	}

	// Overwrite values with those from assertions
	mergeModels(&storageFacility, assertions.StorageFacility)

	mustCreate(db, &storageFacility, assertions.Stub)

	return storageFacility
}

// MakeDefaultStorageFacility makes a single StorageFacility with default values
func MakeDefaultStorageFacility(db *pop.Connection) models.StorageFacility {
	return MakeStorageFacility(db, Assertions{})
}
