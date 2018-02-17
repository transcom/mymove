package testdatagen

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/markbates/pop"
	"github.com/transcom/mymove/pkg/models"
)

// MakeBestValueScore makes a single best_value_score record
func MakeBestValueScore(db *pop.Connection, tsp models.TransportationServiceProvider,
	tdl models.TrafficDistributionList, score int) error {

	bestValueScore := models.BestValueScore{
		TransportationServiceProviderID: tsp.ID,
		Score: score,
		TrafficDistributionListID: tdl.ID,
	}

	_, err := db.ValidateAndSave(&bestValueScore)
	if err != nil {
		log.Panic(err)
	}

	return err
}

// MakeBestValueScoreData creates three best value score records
func MakeBestValueScoreData(db *pop.Connection) {
	// These two queries duplicate ones in other testdatagen files; not optimal
	tspList := []models.TransportationServiceProvider{}
	err := db.All(&tspList)
	if err != nil {
		fmt.Println("TSP ID import failed.")
	}

	tdlList := []models.TrafficDistributionList{}
	err = db.All(&tdlList)
	if err != nil {
		fmt.Println("TDL ID import failed.")
	}

	// Make 3 BestValueScores with random TSPs, random TDLs, and random scores
	for i := 0; i < 3; i++ {
		MakeBestValueScore(
			db,
			tspList[rand.Intn(len(tspList))],
			tdlList[rand.Intn(len(tdlList))],
			rand.Intn(99),
		)
	}
}
