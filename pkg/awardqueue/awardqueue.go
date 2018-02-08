package awardqueue

import (
	"fmt"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

var dbConnection *pop.Connection

func makeSomeStuffUp() {
	tdl1 := models.TrafficDistributionList{
		SourceRateArea:    "here",
		DestinationRegion: "there",
		CodeOfService:     "place",
	}
	dbConnection.ValidateAndSave(&tdl1)
}

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

func findAllTrafficDistributionLists() {
	tdls := []models.TrafficDistributionList{}
	err := dbConnection.All(&tdls)
	if err != nil {
		fmt.Printf("Oh snap! %s", err)
	}

	for i, tdl := range tdls {
		fmt.Printf("TDL %d:\n%v\n", i, tdl)
	}
}
func Run(db *pop.Connection) {
	dbConnection = db

	makeSomeStuffUp()

	fmt.Println("Hello, TSP Award Queue!")

	findAllShipments()
	findAllTrafficDistributionLists()
}
