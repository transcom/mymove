package factory

import (
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type boatBuildType byte

const (
	boatBuildStandard boatBuildType = iota
	boatBuildStandardTowAway
	boatBuildStandardHaulAway
)

// buildBoatShipmentWithBuildType does the actual work
// It builds
//   - MTOShipment and associated set relationships
//
// These will be created if and only if a customization is provided
//   - W2Address
func buildBoatShipmentWithBuildType(db *pop.Connection, customs []Customization, traits []Trait, buildType boatBuildType) models.BoatShipment {
	customs = setupCustomizations(customs, traits)

	// Find boatShipment assertion and convert to models.BoatShipment
	var cBoatShipment models.BoatShipment
	if result := findValidCustomization(customs, BoatShipment); result != nil {
		cBoatShipment = result.Model.(models.BoatShipment)
		if result.LinkOnly {
			return cBoatShipment
		}
	}

	traits = append(traits, GetTraitBoatShipment)
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

	if buildType == boatBuildStandardTowAway {
		shipment.ShipmentType = models.MTOShipmentTypeBoatTowAway
	}

	if buildType == boatBuildStandardHaulAway {
		shipment.ShipmentType = models.MTOShipmentTypeBoatHaulAway
	}

	boatShipment := models.BoatShipment{
		ShipmentID:     shipment.ID,
		Shipment:       shipment,
		Type:           models.BoatShipmentTypeHaulAway,
		Year:           models.IntPointer(2000),
		Make:           models.StringPointer("Boat Make"),
		Model:          models.StringPointer("Boat Model"),
		LengthInInches: models.IntPointer(300),
		WidthInInches:  models.IntPointer(108),
		HeightInInches: models.IntPointer(72),
		HasTrailer:     models.BoolPointer(true),
		IsRoadworthy:   models.BoolPointer(false),
	}

	if buildType == boatBuildStandardTowAway {
		boatShipment.Type = models.BoatShipmentTypeTowAway
		boatShipment.HasTrailer = models.BoolPointer(true)
		boatShipment.IsRoadworthy = models.BoolPointer(true)
	}

	if buildType == boatBuildStandardHaulAway {
		boatShipment.Type = models.BoatShipmentTypeHaulAway
		boatShipment.HasTrailer = models.BoolPointer(false)
		boatShipment.IsRoadworthy = models.BoolPointer(false)
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&boatShipment, cBoatShipment)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &boatShipment)
	}

	boatShipment.Shipment.BoatShipment = &boatShipment

	return boatShipment
}

func BuildBoatShipment(db *pop.Connection, customs []Customization, traits []Trait) models.BoatShipment {
	return buildBoatShipmentWithBuildType(db, customs, traits, boatBuildStandard)
}
func BuildBoatShipmentTowAway(db *pop.Connection, customs []Customization, traits []Trait) models.BoatShipment {
	return buildBoatShipmentWithBuildType(db, customs, traits, boatBuildStandardTowAway)
}
func BuildBoatShipmentHaulAway(db *pop.Connection, customs []Customization, traits []Trait) models.BoatShipment {
	return buildBoatShipmentWithBuildType(db, customs, traits, boatBuildStandardHaulAway)
}

// ------------------------
//        TRAITS
// ------------------------

func GetTraitBoatShipment() []Customization {
	return []Customization{
		{
			Model: models.MTOShipment{
				Status:       models.MTOShipmentStatusSubmitted,
				ShipmentType: models.MTOShipmentTypeBoatHaulAway,
			},
		},
	}
}
