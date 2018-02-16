package testdatagen

import (
	"fmt"
	"github.com/markbates/pop"
	"github.com/transcom/mymove/pkg/models"
	"log"
	"time"
)

// MakeShipmentData creates three shipment records
func MakeShipmentData(dbConnection *pop.Connection) {
	// Grab three UUIDs for individual TDLs
	// TODO: should this query be made in main, between creation functions,
	// and then sourced from one central place?
	tdlList := []models.TrafficDistributionList{}
	err := dbConnection.All(&tdlList)
	if err != nil {
		fmt.Println("TDL ID import failed.")
	}

	// Add three shipment table records using UUIDs from TDLs
	time := time.Now()

	shipment1 := models.Shipment{
		TrafficDistributionListID: tdlList[0].ID,
		PickupDate:                time,
		DeliveryDate:              time,
	}

	shipment2 := models.Shipment{
		TrafficDistributionListID: tdlList[1].ID,
		PickupDate:                time,
		DeliveryDate:              time,
	}

	shipment3 := models.Shipment{
		TrafficDistributionListID: tdlList[2].ID,
		PickupDate:                time,
		DeliveryDate:              time,
	}

	_, err = dbConnection.ValidateAndSave(&shipment3)
	if err != nil {
		log.Panic(err)
	}

	_, err = dbConnection.ValidateAndSave(&shipment1)
	if err != nil {
		log.Panic(err)
	}

	_, err = dbConnection.ValidateAndSave(&shipment2)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("make_shipment_records ran")

}
