package scenario

import (
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	fakedata "github.com/transcom/mymove/pkg/fakedata_approved"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/random"

	"github.com/pkg/errors"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/testdatagen"
)

// NamedScenario is a data generation scenario that has a name
type NamedScenario struct {
	Name         string
	SubScenarios []string
}

// May15TestYear is a May 15 of TestYear
var May15TestYear = time.Date(testdatagen.TestYear, time.May, 15, 0, 0, 0, 0, time.UTC)

// Oct1TestYear is October 1 of TestYear
var Oct1TestYear = time.Date(testdatagen.TestYear, time.October, 1, 0, 0, 0, 0, time.UTC)

// Dec31TestYear is December 31 of TestYear
var Dec31TestYear = time.Date(testdatagen.TestYear, time.December, 31, 0, 0, 0, 0, time.UTC)

// May14FollowingYear is May 14 of the year AFTER TestYear
var May14FollowingYear = time.Date(testdatagen.TestYear+1, time.May, 14, 0, 0, 0, 0, time.UTC)

func save(db *pop.Connection, model interface{}) error {
	verrs, err := db.ValidateAndSave(model)
	if err != nil {
		return errors.Wrap(err, "Errors encountered saving model")
	}
	if verrs.HasAny() {
		return errors.Errorf("Validation errors encountered saving model: %v", verrs)
	}
	return nil
}

// createRandomMove creates a random move with fake data that has been approved for usage
func createRandomMove(db *pop.Connection, possibleStatuses []models.MoveStatus, allDutyStations []models.DutyStation,
	dutyStationsInGBLOC []models.DutyStation, assertions testdatagen.Assertions) {
	randDays, err := random.GetRandomInt(366)
	if err != nil {
		log.Panic(fmt.Errorf("Unable to generate random integer for submitted move date"), zap.Error(err))
	}
	submittedAt := time.Now().AddDate(0, 0, randDays*-1)

	if assertions.ServiceMember.Affiliation == nil {
		randomAffiliation, err := random.GetRandomInt(5)
		if err != nil {
			log.Panic(fmt.Errorf("Unable to generate random integer for affiliation"), zap.Error(err))
		}
		assertions.ServiceMember.Affiliation = &[]models.ServiceMemberAffiliation{
			models.AffiliationARMY,
			models.AffiliationAIRFORCE,
			models.AffiliationNAVY,
			models.AffiliationCOASTGUARD,
			models.AffiliationMARINES}[randomAffiliation]
	}

	dutyStationCount := len(allDutyStations)
	if assertions.Order.OriginDutyStationID == nil {
		// We can pick any origin duty station not only one in the office user's GBLOC
		if *assertions.ServiceMember.Affiliation == models.AffiliationMARINES {
			randDutyStaionIndex, err := random.GetRandomInt(dutyStationCount)
			if err != nil {
				log.Panic(fmt.Errorf("Unable to generate random integer for duty station"), zap.Error(err))
			}
			assertions.Order.OriginDutyStation = &allDutyStations[randDutyStaionIndex]
			assertions.Order.OriginDutyStationID = &assertions.Order.OriginDutyStation.ID
		} else {
			randDutyStaionIndex, err := random.GetRandomInt(len(dutyStationsInGBLOC))
			if err != nil {
				log.Panic(fmt.Errorf("Unable to generate random integer for duty station"), zap.Error(err))
			}
			assertions.Order.OriginDutyStation = &dutyStationsInGBLOC[randDutyStaionIndex]
			assertions.Order.OriginDutyStationID = &assertions.Order.OriginDutyStation.ID
		}
	}

	if assertions.Order.NewDutyStationID == uuid.Nil {
		randDutyStaionIndex, err := random.GetRandomInt(dutyStationCount)
		if err != nil {
			log.Panic(fmt.Errorf("Unable to generate random integer for duty station"), zap.Error(err))
		}
		assertions.Order.NewDutyStation = allDutyStations[randDutyStaionIndex]
		assertions.Order.NewDutyStationID = assertions.Order.NewDutyStation.ID
	}

	randomFirst, randomLast := fakedata.RandomName()
	assertions.ServiceMember.FirstName = &randomFirst
	assertions.ServiceMember.LastName = &randomLast

	orders := testdatagen.MakeOrderWithoutDefaults(db, assertions)

	if assertions.Move.SubmittedAt == nil {
		assertions.Move.SubmittedAt = &submittedAt
	}

	if assertions.Move.Status == "" {
		randStatusIndex, err := random.GetRandomInt(len(possibleStatuses))
		if err != nil {
			log.Panic(fmt.Errorf("Unable to generate random integer for move status"), zap.Error(err))
		}
		assertions.Move.Status = possibleStatuses[randStatusIndex]

		if assertions.Move.Status == models.MoveStatusServiceCounselingCompleted {
			counseledAt := submittedAt.Add(3 * 24 * time.Hour)
			assertions.Move.ServiceCounselingCompletedAt = &counseledAt
		}
	}
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move:  assertions.Move,
		Order: orders,
	})

	shipmentStatus := models.MTOShipmentStatusSubmitted
	if assertions.MTOShipment.Status != "" {
		shipmentStatus = assertions.MTOShipment.Status
	}

	laterRequestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	laterRequestedDeliveryDate := laterRequestedPickupDate.Add(7 * 24 * time.Hour)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType:          models.MTOShipmentTypeHHG,
			Status:                shipmentStatus,
			RequestedPickupDate:   &laterRequestedPickupDate,
			RequestedDeliveryDate: &laterRequestedDeliveryDate,
			ApprovedDate:          assertions.MTOShipment.ApprovedDate,
			Diversion:             assertions.MTOShipment.Diversion,
		},
	})

	earlierRequestedPickupDate := submittedAt.Add(30 * 24 * time.Hour)
	earlierRequestedDeliveryDate := earlierRequestedPickupDate.Add(7 * 24 * time.Hour)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType:          models.MTOShipmentTypeHHG,
			Status:                shipmentStatus,
			RequestedPickupDate:   &earlierRequestedPickupDate,
			RequestedDeliveryDate: &earlierRequestedDeliveryDate,
			ApprovedDate:          assertions.MTOShipment.ApprovedDate,
			Diversion:             assertions.MTOShipment.Diversion,
		},
	})
}
