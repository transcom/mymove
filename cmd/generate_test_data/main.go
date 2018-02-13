package main

import (
	"fmt"
	"github.com/markbates/pop"
	"github.com/namsral/flag"
	"github.com/transcom/mymove/pkg/models"
	"log"
	"time"
)

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

	// Make three TDL records
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

	_, err = dbConnection.ValidateAndSave(&tdl1)
	if err != nil {
		log.Panic(err)
	}

	_, err = dbConnection.ValidateAndSave(&tdl2)
	if err != nil {
		log.Panic(err)
	}

	_, err = dbConnection.ValidateAndSave(&tdl3)
	if err != nil {
		log.Panic(err)
	}

	// Make three TSP table records
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

	// Make three shipment records
	// Grab three UUIDs for individual TDLs
	tdlList := []models.TrafficDistributionList{}
	err = dbConnection.All(&tdlList)
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

	// Make an awarded shipment record
	// Get a shipment ID
	shipmentList := []models.Shipment{}
	err = dbConnection.All(&shipmentList)
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
	award1 := models.ShipmentAward{
		ShipmentID:                      shipmentList[0].ID,
		TransportationServiceProviderID: tspList[0].ID,
		AdministrativeShipment:          false,
	}

	_, err = dbConnection.ValidateAndSave(&award1)
	if err != nil {
		log.Panic(err)
	}

	// Add three best value scores
	bestValueScore1 := models.BestValueScore{
		TransportationServiceProviderID: tspList[0].ID,
		Score: 11,
	}

	_, err = dbConnection.ValidateAndSave(&bestValueScore1)
	if err != nil {
		log.Panic(err)
	}

	bestValueScore2 := models.BestValueScore{
		TransportationServiceProviderID: tspList[1].ID,
		Score: 2,
	}

	_, err = dbConnection.ValidateAndSave(&bestValueScore2)
	if err != nil {
		log.Panic(err)
	}

	bestValueScore3 := models.BestValueScore{
		TransportationServiceProviderID: tspList[2].ID,
		Score: 8,
	}

	_, err = dbConnection.ValidateAndSave(&bestValueScore3)
	if err != nil {
		log.Panic(err)
	}

}
