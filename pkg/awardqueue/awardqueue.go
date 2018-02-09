package awardqueue

import (
	"fmt"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

var dbConnection *pop.Connection

func findAllUnawardedShipments() ([]models.ShipmentWithAwardedTSP, error) {
	shipments, err := models.FetchAwardedShipments(dbConnection)
	return shipments, err
}

func awardShipment(shipment models.ShipmentWithAwardedTSP) error {
	fmt.Printf("Attempting to award shipment: %v\n", shipment.ID)

	// Query shipment's TDL
	tdl := models.TrafficDistributionList{}
	err := dbConnection.Find(&tdl, shipment.TrafficDistributionListId)

	// Query TSPs in that TDL sorted by awarded_shipments[asc] and bvs[desc]
	tsps, err := models.FetchTransportationServiceProvidersInTDL(dbConnection, tdl)

	for _, tsp := range tsps {
		fmt.Printf("Considering TSP: %v\n", tsp)
	}

	return err
}

/*Run will execute the Award Queue algorithm described below.
- Given all unawarded shipments...
- Query TSPs in the TDL, sorted by awarded_shipments[asc] and bvs[desc]
- Create awarded_shipment for the shipment<->tsp
*/
func Run(db *pop.Connection) {
	dbConnection = db

	fmt.Println("TSP Award Queue running.")

	shipments, err := findAllUnawardedShipments()
	if err == nil {
		for _, shipment := range shipments {
			awardShipment(shipment)
		}
	} else {
		fmt.Printf("Failed to query for shipments: %s", err)
	}
}
