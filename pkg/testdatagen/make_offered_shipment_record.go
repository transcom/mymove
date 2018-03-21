package testdatagen

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/go-openapi/swag"
	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeShipmentOffer a single ShipmentOffer record
func MakeShipmentOffer(db *pop.Connection, shipment models.Shipment,
	tsp models.TransportationServiceProvider, admin bool, accepted *bool,
	rejectionReason *string) (models.ShipmentOffer, error) {

	// Add one offered shipment record using existing shipment and TSP IDs
	shipmentOffer := models.ShipmentOffer{
		ShipmentID:                      shipment.ID,
		TransportationServiceProviderID: tsp.ID,
		AdministrativeShipment:          admin,
		Accepted:                        accepted,
		RejectionReason:                 rejectionReason,
	}

	_, err := db.ValidateAndSave(&shipmentOffer)
	if err != nil {
		log.Panic(err)
	}

	return shipmentOffer, err
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

	// Add one offered shipment record using existing, random shipment and TSP IDs
	MakeShipmentOffer(db,
		shipmentList[rand.Intn(len(shipmentList))],
		tspList[rand.Intn(len(tspList))],
		false,
		swag.Bool(true),
		nil,
	)
}
