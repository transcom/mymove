package testdatagen

import (
	"fmt"
	"github.com/markbates/pop"
	"github.com/transcom/mymove/pkg/models"
	"log"
)

// MakeTSP makes a single transportation service provider record.
func MakeTSP(db *pop.Connection, name string, SCAC string) (models.TransportationServiceProvider, error) {
	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: SCAC,
		Name: name}

	_, err := db.ValidateAndSave(&tsp)
	if err != nil {
		log.Panic(err)
	}

	return tsp, err
}

// MakeTSPData creates three TSP records
func MakeTSPData(db *pop.Connection) {
	MakeTSP(db, "Very Good TSP", "ABCD")
	MakeTSP(db, "Pretty Alright TSP", "EFGH")
	MakeTSP(db, "Serviceable and Adequate TSP", "IJKL")
}

// MakeMoreTSP creates numTSP number of TSP records
func MakeTSPs(db *pop.Connection, numTSP int) {
	for i := 0; i < numTSP; i++ {
		tspName := fmt.Sprintf("Just another TSP %d", i)
		tspSCAC := fmt.Sprintf("XYZ%d", i)
		MakeTSP(db, tspName, tspSCAC)
	}

}
