package testdatagen

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"
)

// MakeTSPPerformance makes a single best_value_score record
func MakeTSPPerformance(db *pop.Connection, tsp models.TransportationServiceProvider,
	tdl models.TrafficDistributionList, qualityBand *int, score int, offerCount int,
	peakRateCycle bool) (models.TransportationServiceProviderPerformance, error) {

	tspPerformance := models.TransportationServiceProviderPerformance{
		PerformancePeriodStart:          PerformancePeriodStart,
		PerformancePeriodEnd:            PerformancePeriodEnd,
		PeakRateCycle:                   peakRateCycle,
		TransportationServiceProviderID: tsp.ID,
		TrafficDistributionListID:       tdl.ID,
		QualityBand:                     qualityBand,
		BestValueScore:                  score,
		OfferCount:                      offerCount,
	}

	_, err := db.ValidateAndSave(&tspPerformance)
	if err != nil {
		log.Panic(err)
	}

	return tspPerformance, err
}

// MakeTSPPerformanceData creates three best value score records
// Variable rounds describes how many rounds should have already been offered
// `none` indicates no rounds have been offered, `half` indicates half a round, and `full` a full round.
func MakeTSPPerformanceData(db *pop.Connection, rounds string) {
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
		// For quality band 1, generate a random number between 75 - 100,
		// for quality band 2 between 50-75, etc.
		var offers int
		minBvs := (qualityBand - 1) * 25
		bvs := 100 - (rand.Intn(25) + minBvs)
		// Set rounds according to the flag passed in

		if rounds == "half" {
			if qualityBand == 1 || qualityBand == 2 {
				offers = models.OffersPerQualityBand[qualityBand]
			} else {
				offers = 0
			}
		} else if rounds == "full" {
			offers = models.OffersPerQualityBand[qualityBand]
		} else {
			// default case, no offers
			offers = 0
		}

		MakeTSPPerformance(
			db,
			tspList[rand.Intn(len(tspList))],
			tdlList[rand.Intn(len(tdlList))],
			&qualityBand,
			bvs,
			offers,
			true,
		)
	}
}
