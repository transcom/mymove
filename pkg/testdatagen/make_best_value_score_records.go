package testdatagen

import (
	"fmt"
	"github.com/markbates/pop"
	"github.com/transcom/mymove/pkg/models"
	"log"
)

// MakeBestValueScoreRecords creates three best value score records
func MakeBestValueScoreRecords(dbConnection *pop.Connection) {
	// These two queries duplicate ones in other testdatagen files; not optimal
	tspList := []models.TransportationServiceProvider{}
	err := dbConnection.All(&tspList)
	if err != nil {
		fmt.Println("TSP ID import failed.")
	}

	tdlList := []models.TrafficDistributionList{}
	err = dbConnection.All(&tdlList)
	if err != nil {
		fmt.Println("TDL ID import failed.")
	}

	bestValueScore1 := models.BestValueScore{
		TransportationServiceProviderID: tspList[0].ID,
		Score: 11,
		TrafficDistributionListID: tdlList[0].ID,
	}

	bestValueScore2 := models.BestValueScore{
		TransportationServiceProviderID: tspList[1].ID,
		Score: 2,
		TrafficDistributionListID: tdlList[1].ID,
	}

	bestValueScore3 := models.BestValueScore{
		TransportationServiceProviderID: tspList[2].ID,
		Score: 8,
		TrafficDistributionListID: tdlList[1].ID,
	}

	_, err = dbConnection.ValidateAndSave(&bestValueScore1)
	if err != nil {
		log.Panic(err)
	}

	_, err = dbConnection.ValidateAndSave(&bestValueScore2)
	if err != nil {
		log.Panic(err)
	}

	_, err = dbConnection.ValidateAndSave(&bestValueScore3)
	if err != nil {
		log.Panic(err)
	}
}
