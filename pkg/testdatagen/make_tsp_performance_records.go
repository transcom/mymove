package testdatagen

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/markbates/pop"
	"github.com/transcom/mymove/pkg/models"
)

// MakeTspPerformance makes a single best_value_score record
func MakeTspPerformance(db *pop.Connection, tsp models.TransportationServiceProvider,
	tdl models.TrafficDistributionList, qualityBand *int, score int, awardCount int) (models.TransportationServiceProviderPerformance, error) {

	tspPerformance := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          time.Now(),
		PerformancePeriodEnd:            time.Now(),
		TransportationServiceProviderID: tsp.ID,
		TrafficDistributionListID:       tdl.ID,
		QualityBand:                     qualityBand,
		BestValueScore:                  score,
		AwardCount:                      awardCount,
	}

	_, err := db.ValidateAndSave(&tspPerformance)
	if err != nil {
		log.Panic(err)
	}

	return tspPerformance, err
}

// MakeTspPerformanceData creates three best value score records
func MakeTspPerformanceData(db *pop.Connection) {
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

	// Make 4 TspPerformances with random TSPs, random TDLs, different quality bands, and random scores
	for qualityBand := 1; qualityBand < 5; qualityBand++ {
		// For quality band 1, generate a random number between 0 - 25,
		// for quality band 2 between 25-50, etc.
		minBvs := (qualityBand - 1) * 25
		bvs := rand.Intn(25) + minBvs
		MakeTspPerformance(
			db,
			tspList[rand.Intn(len(tspList))],
			tdlList[rand.Intn(len(tdlList))],
			&qualityBand,
			bvs,
			0,
		)
	}
}
