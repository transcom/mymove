package testdatagen

import (
	"log"
	"math/rand"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// No test includes a zip3 or a volume_move value, as we're not
// currently fully implementing those.

// MakeBlackoutDate creates a test blackoutDate object to add to the database.
func MakeBlackoutDate(db *pop.Connection, assertions Assertions) models.BlackoutDate {
	tspID := assertions.BlackoutDate.TransportationServiceProviderID
	if isZeroUUID(tspID) {
		// Fetches random TSP
		tspList := []models.TransportationServiceProvider{}
		err := db.All(&tspList)
		if err != nil {
			log.Panic(err)
		}
		tspID = tspList[rand.Intn(len(tspList))].ID
	}

	tdlID := assertions.BlackoutDate.TrafficDistributionListID
	if tdlID == nil {
		// Fetches random TDL
		tdlList := []models.TrafficDistributionList{}
		err := db.All(&tdlList)
		if err != nil {
			log.Panic(err)
		}
		tdlID = &tdlList[rand.Intn(len(tdlList))].ID
	}

	blackoutDate := models.BlackoutDate{
		TransportationServiceProviderID: tspID,
		StartBlackoutDate:               NextValidMoveDate,
		EndBlackoutDate:                 NextValidMoveDate,
		TrafficDistributionListID:       tdlID,
		SourceGBLOC:                     stringPointer("PORK"),
		Market:                          stringPointer("dHHG"),
	}

	mergeModels(&blackoutDate, assertions.BlackoutDate)

	mustCreate(db, &blackoutDate)

	return blackoutDate
}

// MakeDefaultBlackoutDate returns a BlackoutDate with default vales
func MakeDefaultBlackoutDate(db *pop.Connection) models.BlackoutDate {
	return MakeBlackoutDate(db, Assertions{})
}
