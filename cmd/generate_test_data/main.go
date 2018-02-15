package main

import (
	// "fmt"
	"github.com/markbates/pop"
	"github.com/namsral/flag"
	"github.com/transcom/mymove/pkg/testdatagen"
	"log"
	// "time"
	// "reflect"
)

// Hey, refactoring self: you can pull the UUIDs from the objects rather than
// querying the db for them again.
func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presently.")
	flag.Parse()

	//DB connection
	pop.AddLookupPaths(*config)
	dbConnection, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	// fmt.Println(reflect.TypeOf(dbConnection))

	// testdatagen.MakeTDLData(dbConnection)
	// testdatagen.MakeTSPData(dbConnection)
	testdatagen.MakeShipmentData(dbConnection)

	// 	// Make three shipment records
	// 	// Grab three UUIDs for individual TDLs
	// 	tdlList := []models.TrafficDistributionList{}
	// 	err = dbConnection.All(&tdlList)
	// 	if err != nil {
	// 		fmt.Println("TDL ID import failed.")
	// 	}

	// 	// Add three shipment table records using UUIDs from TDLs
	// 	time := time.Now()

	// 	shipment1 := models.Shipment{
	// 		TrafficDistributionListID: tdlList[0].ID,
	// 		PickupDate:                time,
	// 		DeliveryDate:              time,
	// 	}

	// 	shipment2 := models.Shipment{
	// 		TrafficDistributionListID: tdlList[1].ID,
	// 		PickupDate:                time,
	// 		DeliveryDate:              time,
	// 	}

	// 	shipment3 := models.Shipment{
	// 		TrafficDistributionListID: tdlList[2].ID,
	// 		PickupDate:                time,
	// 		DeliveryDate:              time,
	// 	}

	// 	_, err = dbConnection.ValidateAndSave(&shipment3)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}

	// 	_, err = dbConnection.ValidateAndSave(&shipment1)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}

	// 	_, err = dbConnection.ValidateAndSave(&shipment2)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}

	// 	// Make an awarded shipment record
	// 	// Get a shipment ID
	// 	shipmentList := []models.Shipment{}
	// 	err = dbConnection.All(&shipmentList)
	// 	if err != nil {
	// 		fmt.Println("Shipment ID import failed.")
	// 	}

	// 	// Get a TSP ID
	// 	tspList := []models.TransportationServiceProvider{}
	// 	err = dbConnection.All(&tspList)
	// 	if err != nil {
	// 		fmt.Println("TSP ID import failed.")
	// 	}

	// 	// Add one awarded shipment record using existing shipment and TSP IDs
	// 	award1 := models.ShipmentAward{
	// 		ShipmentID:                      shipmentList[0].ID,
	// 		TransportationServiceProviderID: tspList[0].ID,
	// 		AdministrativeShipment:          false,
	// 	}

	// 	_, err = dbConnection.ValidateAndSave(&award1)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}

	// 	// Add three best value scores
	// 	bestValueScore1 := models.BestValueScore{
	// 		TransportationServiceProviderID: tspList[0].ID,
	// 		Score: 11,
	// 		TrafficDistributionListID: tdlList[0].ID,
	// 	}

	// 	_, err = dbConnection.ValidateAndSave(&bestValueScore1)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}

	// 	bestValueScore2 := models.BestValueScore{
	// 		TransportationServiceProviderID: tspList[1].ID,
	// 		Score: 2,
	// 		TrafficDistributionListID: tdlList[1].ID,
	// 	}

	// 	_, err = dbConnection.ValidateAndSave(&bestValueScore2)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}

	// 	bestValueScore3 := models.BestValueScore{
	// 		TransportationServiceProviderID: tspList[2].ID,
	// 		Score: 8,
	// 		TrafficDistributionListID: tdlList[1].ID,
	// 	}

	// 	_, err = dbConnection.ValidateAndSave(&bestValueScore3)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}

}
