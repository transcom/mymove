package factory

import (
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type mobileHomeBuildType byte

const (
	mobileHomeBuildStandard mobileHomeBuildType = iota
)

// buildMobileHomeShipmentWithBuildType does the actual work
// It builds
//   - MTOShipment and associated set relationships
//
// These will be created if and only if a customization is provided
//   - W2Address
func buildMobileHomeShipmentWithBuildType(db *pop.Connection, customs []Customization, traits []Trait, buildType mobileHomeBuildType) models.MobileHome {
	customs = setupCustomizations(customs, traits)

	// Find mobileHomeShipment assertion and convert to models.MobileHome
	var cMobileHomeShipment models.MobileHome
	if result := findValidCustomization(customs, MobileHome); result != nil {
		cMobileHomeShipment = result.Model.(models.MobileHome)
		if result.LinkOnly {
			return cMobileHomeShipment
		}
	}

	traits = append(traits, GetTraitMobileHomeShipment)
	shipment := BuildMTOShipment(db, customs, traits)

	serviceMember := shipment.MoveTaskOrder.Orders.ServiceMember
	if serviceMember.ResidentialAddressID == nil {
		log.Panic("Created shipment has service member without ResidentialAddressID")
	}
	if serviceMember.ResidentialAddress == nil {
		var address models.Address
		err := db.Find(&address, serviceMember.ResidentialAddressID)
		if err != nil {
			log.Panicf("Cannot find address with ID %s: %s",
				serviceMember.ResidentialAddressID, err)
		}
		serviceMember.ResidentialAddress = &address
	}

	if buildType == mobileHomeBuildStandard {
		shipment.ShipmentType = models.MTOShipmentTypeMobileHome
	}

	date := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	mobileHomeShipment := models.MobileHome{
	ShipmentID:     shipment.ID,
	Shipment:       shipment,
	Year:           models.IntPointer(2000),
	Make:           models.StringPointer("Mobile Home Make"),
	Model:          models.StringPointer("Mobile Home Model"),
	LengthInInches: models.IntPointer(300),
	WidthInInches:  models.IntPointer(108),
	HeightInInches: models.IntPointer(72),
	CreatedAt: date,
}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&mobileHomeShipment, cMobileHomeShipment)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &mobileHomeShipment)
	}

	mobileHomeShipment.Shipment.MobileHome = &mobileHomeShipment

	return mobileHomeShipment
}

func BuildMobileHomeShipment(db *pop.Connection, customs []Customization, traits []Trait) models.MobileHome {
	return buildMobileHomeShipmentWithBuildType(db, customs, traits, mobileHomeBuildStandard)
}

// ------------------------
//        TRAITS
// ------------------------

func GetTraitMobileHomeShipment() []Customization {
	return []Customization{
		{
			Model: models.MTOShipment{
				Status:       models.MTOShipmentStatusSubmitted,
				ShipmentType: models.MTOShipmentTypeMobileHome,
			},
		},
	}
}