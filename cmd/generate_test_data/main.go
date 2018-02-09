package main

import (
	"fmt"
	"github.com/markbates/pop"
	"github.com/namsral/flag"
	"github.com/transcom/mymove/pkg/models"
	"log"
	"time"
)

var config, env, dbConnection String

func init() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presently.")
	flag.Parse()

	//DB connection
	pop.AddLookupPaths(*config)
	dbConnection, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Hey, it's init")
}

func makeTDLData() {
	// Add three TDL records
	tdl1 := models.TrafficDistributionList{
		SourceRateArea:    "california",
		DestinationRegion: "90210",
		CodeOfService:     "2"}

	tdl2 := models.TrafficDistributionList{
		SourceRateArea:    "north carolina",
		DestinationRegion: "27007",
		CodeOfService:     "4"}

	tdl3 := models.TrafficDistributionList{
		SourceRateArea:    "washington",
		DestinationRegion: "98310",
		CodeOfService:     "1"}

	_, err1 := dbConnection.ValidateAndSave(&tdl1)
	if err1 != nil {
		log.Panic(err)
	}

	_, err2 := dbConnection.ValidateAndSave(&tdl2)
	if err2 != nil {
		log.Panic(err)
	}

	_, err3 := dbConnection.ValidateAndSave(&tdl3)
	if err3 != nil {
		log.Panic(err)
	}
}

func makeTSPData() {
	// Add three TSP table records
	tsp1 := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "ABCD",
		Name: "Very Good TSP"}

	tsp2 := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "EFGH",
		Name: "Pretty Alright TSP"}

	tsp3 := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "IJKL",
		Name: "Serviceable and Adequate TSP"}

	_, err = dbConnection.ValidateAndSave(&tsp1)
	if err != nil {
		log.Panic(err)
	}

	_, err = dbConnection.ValidateAndSave(&tsp2)
	if err != nil {
		log.Panic(err)
	}

	_, err = dbConnection.ValidateAndSave(&tsp3)
	if err != nil {
		log.Panic(err)
	}
}

func makeShipmentData() {
	// Grab three UUIDs for individual TDLs
	tdlList := []models.TrafficDistributionList{}
	err1 := dbConnection.RawQuery("SELECT * FROM traffic_distribution_lists").All(&tdlList)
	if err1 != nil {
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
}

func makeAwardedShipmentData() {
	// Get a shipment ID
	shipmentList := []models.Shipment{}
	err1 := dbConnection.RawQuery("SELECT * FROM shipments").One(&shipmentList)
	if err1 != nil {
		fmt.Println("Shipment ID import failed.")
	} else {
		fmt.Print(shipmentList[0].ID)
	}

	// Get a TSP ID
	// tspList := []models.TransportationServiceProvider{}
	// err1 := dbConnection.RawQuery("SELECT * FROM traffic_distribution_lists").All(&tspList)
	// if err1 != nil {
	// 	fmt.Println("TSP ID import failed.")
	// }

	// Add one awarded shipment record using existing
}

// Query for newly made records and print IDs in terminal
// tsps := []models.TransportationServiceProvider{}
// err = dbConnection.All(&tsps)
// if err != nil {
// 	fmt.Print("Error!\n")
// 	fmt.Printf("%v\n", err)
// } else {
// 	for _, v := range tsps {
// 		fmt.Print(v.ID)
// 	}
// }

func main() {
	init()
	makeTDLData()
	makeTSPData()
	makeShipmentData()
	fmt.Println("Pay attention to this now")
	makeAwardedShipmentData()
}
