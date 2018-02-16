package testdatagen

import (
	"fmt"
	"github.com/markbates/pop"
	"github.com/transcom/mymove/pkg/models"
	"log"
)

// MakeAwardedShipmentData creates one awarded shipment record
func MakeAwardedShipmentData(dbConnection *pop.Connection) {
	// Make an awarded shipment record
	// Get a shipment ID
	shipmentList := []models.Shipment{}
	err := dbConnection.All(&shipmentList)
	if err != nil {
		fmt.Println("Shipment ID import failed.")
	}

	// Get a TSP ID
	tspList := []models.TransportationServiceProvider{}
	err = dbConnection.All(&tspList)
	if err != nil {
		fmt.Println("TSP ID import failed.")
	}

	// Add one awarded shipment record using existing shipment and TSP IDs
	awardedShipment := models.ShipmentAward{
		ShipmentID:                      shipmentList[0].ID,
		TransportationServiceProviderID: tspList[0].ID,
		AdministrativeShipment:          false,
	}

	_, err = dbConnection.ValidateAndSave(&awardedShipment)
	if err != nil {
		log.Panic(err)
	}
}
