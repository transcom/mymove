package testdatagen

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeShipmentOffer creates a single shipment offer record
func MakeShipmentOffer(db *pop.Connection, assertions Assertions) models.ShipmentOffer {

	// Test for ShipmentID first before creating a new Shipment
	shipmentID := assertions.ShipmentOffer.ShipmentID
	if isZeroUUID(assertions.ShipmentOffer.ShipmentID) {
		// TODO: Make Shipment and get ID
	}

	// Test for TSP ID first before creating a new TSP
	tspID := assertions.ShipmentOffer.TransportationServiceProviderID
	if isZeroUUID(assertions.ShipmentOffer.TransportationServiceProviderID) {
		// TODO: Make TSP and get ID
	}
	shipmentOffer := models.ShipmentOffer{
		ShipmentID:                      shipmentID,
		TransportationServiceProviderID: tspID,
		AdministrativeShipment:          false,
		Accepted:                        swag.Bool(true),
		RejectionReason:                 nil,
	}

	mergeModels(&shipmentOffer, assertions.ShipmentOffer)

	mustCreate(db, &shipmentOffer)

	return shipmentOffer
}

// MakeDefaultShipmentOffer makes a ShipmentOffer with default values
func MakeDefaultShipmentOffer(db *pop.Connection) models.ShipmentOffer {
	return MakeShipmentOffer(db, Assertions{})
}

// MakeShipmentOfferData creates one offered shipment record
func MakeShipmentOfferData(db *pop.Connection) {
	// Get a shipment ID
	shipmentList := []models.Shipment{}
	err := db.All(&shipmentList)
	if err != nil {
		fmt.Println("Shipment ID import failed.")
	}

	// Get a TSP ID
	tspList := []models.TransportationServiceProvider{}
	err = db.All(&tspList)
	if err != nil {
		fmt.Println("TSP ID import failed.")
	}

	// Add one offered shipment record for each shipment and a random TSP IDs
	for _, shipment := range shipmentList {
		shipmentOfferAssertions := Assertions{
			ShipmentOffer: models.ShipmentOffer{
				ShipmentID:                      shipment.ID,
				TransportationServiceProviderID: tspList[rand.Intn(len(tspList))].ID,
				AdministrativeShipment:          false,
				Accepted:                        swag.Bool(true),
				RejectionReason:                 nil,
			},
		}
		MakeShipmentOffer(db, shipmentOfferAssertions)
	}
}

// CreateShipmentOfferData creates a TSP User, A Shipment, and then them to a Shipment Offer
func CreateShipmentOfferData(db *pop.Connection) (tspUser models.TspUser, shipment models.Shipment, shipmentOffer models.ShipmentOffer) {

	// Given: a TSP User
	newTspUser := MakeDefaultTspUser(db)

	// Shipment is created and saved
	now := time.Now()
	tdl, _ := MakeTDL(
		db,
		DefaultSrcRateArea,
		DefaultDstRegion,
		DefaultCOS)
	market := "dHHG"
	sourceGBLOC := "OHAI"
	newShipment, _ := MakeShipment(db, now, now, now.AddDate(0, 0, 1), tdl, sourceGBLOC, &market)

	// Shipment Offer is created and synced to TSP ID and Shipment ID
	shipmentOfferAssertions := Assertions{
		ShipmentOffer: models.ShipmentOffer{
			ShipmentID:                      newShipment.ID,
			TransportationServiceProviderID: newTspUser.TransportationServiceProviderID,
		},
	}
	newShipmentOffer := MakeShipmentOffer(db, shipmentOfferAssertions)

	return newTspUser, newShipment, newShipmentOffer
}
