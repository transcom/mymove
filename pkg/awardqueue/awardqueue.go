package awardqueue

import (
	"fmt"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

var dbConnection *pop.Connection

// This function was made just to get my Golang legs on and play with data
func findAllUnawardedShipments() ([]models.Shipment, error) {
	shipments := []models.Shipment{}
	err := dbConnection.All(&shipments)

	if err != nil {
		fmt.Printf("Oh snap! %s", err)
		return nil, err
	}
	fmt.Printf("Shipments:\n%v\n", shipments)

	return shipments, nil
}

func awardShipment(shipment models.Shipment) {
	fmt.Printf("Trying to award shipment:\n%v\n", shipment)
}

/*Run will execute the Award Queue algorithm described below.
- Given all unawarded shipments...
- Query shipment's TDL
- Query TSPs in the TDL, sorted by awarded_shipments[asc] and bvs[desc]
- Create awarded_shipment for the shipment<->tsp
*/
func Run(db *pop.Connection) {
	dbConnection = db

	fmt.Println("Hello, TSP Award Queue!")

	shipments, err := findAllUnawardedShipments()
	if err == nil {
		for _, shipment := range shipments {
			awardShipment(shipment)
		}
	} else {
		fmt.Printf("Failed to query for shipments!")
	}
}
