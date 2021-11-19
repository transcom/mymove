package testdatagen

import (
	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

// MakeStorageFacility creates a single SIT Extension and associated set relationships
func MakeStorageFacility(db *pop.Connection, assertions Assertions) models.StorageFacility {
	lotNumber := "1234"
	phone := "5555555555"
	email := "storage@email.com"
	address := assertions.StorageFacility.Address

	if address.StreetAddress1 == "" {
		address = MakeAddress(db, assertions)
	}

	storageFacility := models.StorageFacility{
		FacilityName: "Storage R Us",
		LotNumber:    &lotNumber,
		Address:      address,
		AddressID:    address.ID,
		Phone:        &phone,
		Email:        &email,
	}

	storageFacility.Address = address
	storageFacility.AddressID = address.ID

	mergeModels(&storageFacility, assertions.StorageFacility)

	mustCreate(db, &storageFacility, assertions.Stub)

	return storageFacility
}
