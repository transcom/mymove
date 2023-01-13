package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildTransportationOffice creates a single TransportationOffice.
// Also creates, if not provided
// - Address
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildTransportationOffice(db *pop.Connection, customs []Customization, traits []Trait) models.TransportationOffice {
	customs = setupCustomizations(customs, traits)

	// Find TransportationOffice assertion and convert to models.TransportationOffice
	var cTransportationOffice models.TransportationOffice
	if result := findValidCustomization(customs, TransportationOffice); result != nil {
		cTransportationOffice = result.Model.(models.TransportationOffice)
		if result.LinkOnly {
			return cTransportationOffice
		}
	}

	// Find/create the address model
	var address models.Address
	result := findValidCustomization(customs, Address)
	if result != nil {
		address = result.Model.(models.Address)
	}
	address = BuildAddress(db, customs, traits)
	// At this point, address exists. It's either the provided or created address

	// Create transportationOffice
	transportationOffice := models.TransportationOffice{
		Name:      "JPPSO Testy McTest",
		AddressID: address.ID,
		Address:   address,
		Gbloc:     "KKFA",
		Latitude:  1.23445,
		Longitude: -23.34455,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&transportationOffice, cTransportationOffice)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &transportationOffice)
	}
	return transportationOffice
}

// BuildTransportationOfficeWithPhoneLine creates a single TransportationOffice.
// Also creates, if not provided
// - Address
// - OfficePhoneLine
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildTransportationOfficeWithPhoneLine(db *pop.Connection, customs []Customization, traits []Trait) models.TransportationOffice {

	// If an office with a phoneLine is required, best to just build the phone line,
	// which will then create the transportation office.
	phoneLine := BuildOfficePhoneLine(db, customs, traits)
	office := phoneLine.TransportationOffice
	office.PhoneLines = append(office.PhoneLines, phoneLine)

	return office
}

// BuildDefaultTransportationOffice creates one with a phoneline hooked up.
func BuildDefaultTransportationOffice(db *pop.Connection) models.TransportationOffice {
	return BuildTransportationOfficeWithPhoneLine(db, nil, nil)
}
