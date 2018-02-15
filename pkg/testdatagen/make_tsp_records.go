package testdatagen

import (
	"fmt"
	"github.com/markbates/pop"
	"github.com/transcom/mymove/pkg/models"
	"log"
)

// MakeTSPData creates three TSP records
func MakeTSPData(dbConnection *pop.Connection) {
	tsp1 := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "ABCD",
		Name: "Very Good TSP"}

	tsp2 := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "EFGH",
		Name: "Pretty Alright TSP"}

	tsp3 := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: "IJKL",
		Name: "Serviceable and Adequate TSP"}

	_, err := dbConnection.ValidateAndSave(&tsp1)
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

	fmt.Println("make_tsp_records ran")
}
