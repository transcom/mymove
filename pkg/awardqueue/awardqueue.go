package awardqueue

import (
	"fmt"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

var dbConnection *pop.Connection

// This function was made just to get my Golang legs on and play with data
func makeSomeStuffUp() {
	tdl1 := models.TrafficDistributionList{
		SourceRateArea:    "here",
		DestinationRegion: "there",
		CodeOfService:     "place",
	}
	dbConnection.ValidateAndSave(&tdl1)
}

// This function was made just to get my Golang legs on and play with data
func findAllShipments() error {
	shipments := []models.Shipment{}
	err := dbConnection.All(&shipments)
	if err != nil {
		fmt.Printf("Oh snap! %s", err)
		return err
	}
	fmt.Printf("Shipments:\n%v\n", shipments)

	return nil
}

// This function was made just to get my Golang legs on and play with data
func findAllTrafficDistributionLists() error {
	tdls := []models.TrafficDistributionList{}
	err := dbConnection.All(&tdls)
	if err != nil {
		fmt.Printf("Oh snap! %s", err)
		return err
	}

	for i, tdl := range tdls {
		fmt.Printf("TDL %d:\n%v\n", i, tdl)
	}

	return nil
}

/*Run will execute the Award Queue algorithm described below.

- Query for TDLs, so we know what markets exist
- Query for shipments that do not have a matching awarded_shipment
- Group shipments by the TDL they belong to
- For each TDL:
  - Sort Shipments by BVS
  - Query TSPs in this TDL, joined on & sorted by BVS
  - In order of descending BVS, create AwardedShipments matching
    shipments to TSPs.
*/
func Run(db *pop.Connection) {
	dbConnection = db

	makeSomeStuffUp()

	fmt.Println("Hello, TSP Award Queue!")

	findAllShipments()
	findAllTrafficDistributionLists()

}
