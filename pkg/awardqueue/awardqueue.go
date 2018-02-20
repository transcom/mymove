package awardqueue

import (
	"fmt"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

var db *pop.Connection

func findAllUnawardedShipments() ([]models.PossiblyAwardedShipment, error) {
	shipments, err := models.FetchAwardedShipments(db)
	return shipments, err
}

// AttemptShipmentAward will attempt to take the given Shipment and award it to
// a TSP.
func AttemptShipmentAward(shipment models.PossiblyAwardedShipment) (*models.ShipmentAward, error) {
	fmt.Printf("Attempting to award shipment: %v\n", shipment.ID)

	// Query the shipment's TDL
	tdl := models.TrafficDistributionList{}
	err := db.Find(&tdl, shipment.TrafficDistributionListID)

	// Find TSPs in that TDL sorted by shipment_awards[asc] and bvs[desc]
	tsps, err := models.FetchTransportationServiceProvidersInTDL(db, tdl.ID)

	if len(tsps) == 0 {
		return nil, fmt.Errorf("Cannot award. No TSPs found in TDL (%v)", tdl.ID)
	}

	var shipmentAward *models.ShipmentAward

	for _, consideredTSP := range tsps {
		fmt.Printf("\tConsidering TSP: %v\n", consideredTSP.Name)

		tsp := models.TransportationServiceProvider{}
		err := db.Find(&tsp, consideredTSP.ID)
		if err == nil {
			// We found a valid TSP to award to!
			shipmentAward, err = models.CreateShipmentAward(db, shipment.ID, tsp.ID, false)
			if err == nil {
				fmt.Print("\tShipment awarded to TSP!\n")
				break
			} else {
				fmt.Printf("\tFailed to award to TSP: %v\n", err)
			}
		} else {
			fmt.Printf("\tFailed to award to TSP: %v\n", err)
		}
	}

	return shipmentAward, err
}

// Run will execute the Award Queue algorithm.
func Run(dbConnection *pop.Connection) {
	db = dbConnection

	fmt.Println("TSP Award Queue running.")

	shipments, err := findAllUnawardedShipments()
	if err == nil {
		count := 0
		for _, shipment := range shipments {
			_, err = AttemptShipmentAward(shipment)
			if err != nil {
				fmt.Printf("\tFailed to award shipment: %s\n", err)
			} else {
				count++
			}
		}
		fmt.Printf("Awarded %d shipments.\n", count)
	} else {
		fmt.Printf("Failed to query for shipments: %s", err)
	}
}
