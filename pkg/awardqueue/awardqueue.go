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

func selectTSPToAwardShipment(shipment models.ShipmentWithAwardedTSP) error {
	fmt.Printf("Attempting to award shipment: %v\n", shipment.ID)

	// Query the shipment's TDL
	tdl := models.TrafficDistributionList{}
	err := dbConnection.Find(&tdl, shipment.TrafficDistributionListID)

	// Find TSPs in that TDL sorted by awarded_shipments[asc] and bvs[desc]
	tsps, err := models.FetchTransportationServiceProvidersInTDL(dbConnection, tdl.ID)

	for _, consideredTSP := range tsps {
		fmt.Printf("\tConsidering TSP: %v\n", consideredTSP)

		tsp := models.TransportationServiceProvider{}
		err := dbConnection.Find(&tsp, consideredTSP.TransportationServiceProviderID)
		if err == nil {
			// We found a valid TSP to award to!
			err := models.AwardShipment(dbConnection, shipment.ID, tsp.ID, false)
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

	return err
}

// Run will execute the Award Queue algorithm.
func Run(db *pop.Connection) {
	dbConnection = db

	fmt.Println("TSP Award Queue running.")

	shipments, err := findAllUnawardedShipments()
	if err == nil {
		count := -1
		for i, shipment := range shipments {
			err = selectTSPToAwardShipment(shipment)
			if err != nil {
				fmt.Printf("Failed to award shipment: %s\n", err)
			}
			count = i
		}
		fmt.Printf("Awarded %d shipments.", count+1)
	} else {
		fmt.Printf("Failed to query for shipments: %s", err)
	}
}
