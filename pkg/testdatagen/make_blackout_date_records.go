package testdatagen

import (
	"log"
	"math/rand"
	"time"

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
		StartBlackoutDate:               time.Now(),
		EndBlackoutDate:                 time.Now(),
		TrafficDistributionListID:       tdlID,
		SourceGBLOC:                     fmtString("PORK"),
		Market:                          fmtString("dHHG"),
	}

	mergeModels(&blackoutDate, assertions.BlackoutDate)

	mustCreate(db, &blackoutDate)

	return blackoutDate
}

// MakeDefaultBlackoutDate returns a BlackoutDate with default vales
func MakeDefaultBlackoutDate(db *pop.Connection) models.BlackoutDate {
	return MakeBlackoutDate(db, Assertions{})
}

// MakeBlackoutDateData creates three blackoutDate objects and commits them to the blackout_dates table.
func MakeBlackoutDateData(db *pop.Connection) {
	// Make a blackout date with market.
	date1 := MakeDefaultBlackoutDate(db)
	date1.SourceGBLOC = nil
	mustSave(db, &date1)

	// Make a blackout date with a channel.
	date2 := MakeDefaultBlackoutDate(db)
	date2.SourceGBLOC = nil
	date2.Market = nil
	mustSave(db, &date2)

	// Make a blackout date with market and source gbloc.
	MakeDefaultBlackoutDate(db)
}
