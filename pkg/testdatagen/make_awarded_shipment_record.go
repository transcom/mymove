package testdatagen

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/go-openapi/swag"
	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeShipmentAward a single AwardedShipment record
func MakeShipmentAward(db *pop.Connection, shipment models.Shipment,
	tsp models.TransportationServiceProvider, admin bool, accepted *bool,
	rejectionReason *string) (models.ShipmentAward, error) {

	// Add one awarded shipment record using existing shipment and TSP IDs
	shipmentAward := models.ShipmentAward{
		ShipmentID:                      shipment.ID,
		TransportationServiceProviderID: tsp.ID,
		AdministrativeShipment:          admin,
		Accepted:                        accepted,
		RejectionReason:                 rejectionReason,
	}

	_, err := db.ValidateAndSave(&shipmentAward)
	if err != nil {
		log.Panic(err)
	}

	return shipmentAward, err
}

// MakeShipmentAwardData creates one awarded shipment record
func MakeShipmentAwardData(db *pop.Connection) {
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

	// Add one awarded shipment record using existing, random shipment and TSP IDs
	MakeShipmentAward(db,
		shipmentList[rand.Intn(len(shipmentList))],
		tspList[rand.Intn(len(tspList))],
		false,
		swag.Bool(true),
		nil,
	)
}
