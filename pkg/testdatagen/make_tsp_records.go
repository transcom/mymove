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
func MakeTSP(db *pop.Connection, assertions Assertions) models.TransportationServiceProvider {

	scac := assertions.TransportationServiceProvider.StandardCarrierAlphaCode
	if scac == "" {
		scac = RandomSCAC()
	}
	tsp := models.TransportationServiceProvider{
		StandardCarrierAlphaCode: scac,
	}

	verrs, err := db.ValidateAndSave(&tsp)
	if verrs.HasAny() {
		err = fmt.Errorf("TSP validation errors: %v", verrs)
	}
	if err != nil {
		log.Panic(err)
	}

	return tsp
}

// MakeDefaultTSP makes a TSP with default values
func MakeDefaultTSP(db *pop.Connection) models.TransportationServiceProvider {
	return MakeTSP(db, Assertions{})
}

// MakeTSPs creates numTSP number of TSP records
// numTSP specifies how many TSPs to create
func MakeTSPs(db *pop.Connection, numTSP int) {
	for i := 0; i < numTSP; i++ {
		MakeTSP(db, Assertions{
			TransportationServiceProvider: models.TransportationServiceProvider{
				StandardCarrierAlphaCode: RandomSCAC(),
			},
		})
	}
}
