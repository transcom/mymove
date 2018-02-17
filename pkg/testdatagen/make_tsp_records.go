package testdatagen

import (
	"github.com/markbates/pop"
	"github.com/transcom/mymove/pkg/models"
	"log"
)

// MakeTSP makes a single transportation service provider record.
func MakeTSP(db *pop.Connection, name string, SCAC string) error {
	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: SCAC,
		Name: name}

	_, err := db.ValidateAndSave(&tsp)
	if err != nil {
		log.Panic(err)
	}

	return err
}

// MakeTSPData creates three TSP records
func MakeTSPData(db *pop.Connection) {
	MakeTSP(db, "Very Good TSP", "ABCD")
	MakeTSP(db, "Pretty Alright TSP", "EFGH")
	MakeTSP(db, "Serviceable and Adequate TSP", "IJKL")
}
