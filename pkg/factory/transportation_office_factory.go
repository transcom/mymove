package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildTransportationOffice creates a single TransportationOffice.
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
		if !result.LinkOnly {
			// TODO replace this with build
			address = testdatagen.MakeDefaultAddress(db)
		}
	}

	// TODO replace this with build
	address = testdatagen.MakeDefaultAddress(db)
	// At this point, address exists. It's either the provided or created address

	// create transportationOffice
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

	//// create ShippingOffice if needed
	result = findValidCustomization(customs, ShippingOffice)
	if result != nil {
		var shippingOffice models.TransportationOffice
		shippingOffice = result.Model.(models.TransportationOffice)
		if !result.LinkOnly {
			// TODO replace this with build
			shippingOffice = BuildTransportationOffice(nil, customs, nil)
		}
		transportationOffice.ShippingOffice = &shippingOffice
	}

	//// If db is false, it's a stub. No need to create in database
	//if db != nil {
	//	mustCreate(db, &transportationOffice)
	//}

	//// create phoneLines model
	//// Don't currently have a testdatagen/factory function to create OfficePhoneLines
	//var phoneLines []models.OfficePhoneLine
	//phoneLine := models.OfficePhoneLine{
	//	TransportationOfficeID: transportationOffice.ID,
	//	TransportationOffice:   transportationOffice,
	//	Number:                 "(510) 555-5555",
	//	IsDsnNumber:            false,
	//	Type:                   "voice",
	//}
	//phoneLines = append(phoneLines, phoneLine)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &transportationOffice)
		//mustCreate(db, &phoneLine)
	}
	//transportationOffice.PhoneLines = phoneLines

	return transportationOffice
}

// BuildDefaultTransportationOffice builds a default TransportationOffice
func BuildDefaultTransportationOffice(db *pop.Connection, customs []Customization, traits []Trait) models.TransportationOffice {
	return BuildTransportationOffice(db, nil, []Trait{GetTraitAdminUserEmail})
}
