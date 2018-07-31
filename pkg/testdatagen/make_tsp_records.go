package testdatagen

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
)

const alphanumericBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomSCAC generates a random 4 figure string from allowed alphanumeric bytes to represent the SCAC.
func RandomSCAC() string {
	b := make([]byte, 4)
	for i := range b {
		b[i] = alphanumericBytes[rand.Intn(len(alphanumericBytes))]
	}
	return string(b)
}

// MakeTSP makes a single transportation service provider record.
func MakeTSP(db *pop.Connection, SCAC string) (models.TransportationServiceProvider, error) {
	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: SCAC,
	}

	verrs, err := db.ValidateAndSave(&tsp)
	if verrs.HasAny() {
		err = fmt.Errorf("TSP validation errors: %v", verrs)
	}
	if err != nil {
		log.Panic(err)
	}

	return tsp, err
}

// MakeTSPData creates three TSP records
func MakeTSPData(db *pop.Connection) {
	MakeTSP(db, RandomSCAC())
	MakeTSP(db, RandomSCAC())
	MakeTSP(db, RandomSCAC())
}

// MakeTSPs creates numTSP number of TSP records
// numTSP specifies how many TSPs to create
func MakeTSPs(db *pop.Connection, numTSP int) {
	for i := 0; i < numTSP; i++ {
		MakeTSP(db, RandomSCAC())
	}
}
