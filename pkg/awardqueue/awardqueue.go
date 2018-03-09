package awardqueue

import (
	"fmt"
	"math"
	"time"

	"github.com/markbates/pop"
	"github.com/satori/go.uuid"

	"github.com/transcom/mymove/pkg/models"
)

const numQualBands = 4

// Minimum Performance Score (MPS) is the lowest BVS a TSP can have and still be assigned shipments.
// TODO: This will eventually need to be configurable; implement as something other than a constant.
const mps = 10

type qualityBand models.TransportationServiceProviderPerformances
type qualityBands []qualityBand

// AwardQueue encapsulates the TSP award queue process
type AwardQueue struct {
	db *pop.Connection
}

func (aq *AwardQueue) findAllUnawardedShipments() ([]models.PossiblyAwardedShipment, error) {
	shipments, err := models.FetchAwardedShipments(aq.db)
	return shipments, err
}

// AttemptShipmentAward will attempt to take the given Shipment and award it to
// a TSP.
func (aq *AwardQueue) attemptShipmentAward(shipment models.PossiblyAwardedShipment) (*models.ShipmentAward, error) {
	fmt.Printf("Attempting to award shipment: %v\n", shipment.ID)

	// Query the shipment's TDL
	tdl := models.TrafficDistributionList{}
	err := aq.db.Find(&tdl, shipment.TrafficDistributionListID)

	if err != nil {
		return nil, fmt.Errorf("Cannot find TDL in database: %s", err)
	}

	tspPerformances, err := models.FetchTSPPerformanceForAwardQueue(aq.db, tdl.ID, mps)

	if err != nil {
		return nil, fmt.Errorf("Cannot award. Database error: %s", err)
	}

	if len(tspPerformances) == 0 {
		return nil, fmt.Errorf("Cannot award. No TSPs found in TDL (%v)", tdl.ID)
	}

	var shipmentAward *models.ShipmentAward

	for _, tspPerformance := range tspPerformances {
		tsp := models.TransportationServiceProvider{}
		tspBlackoutDatesPresent := aq.checkTSPBlackoutDates(tsp.ID, shipment.PickupDate)

		if err := aq.db.Find(&tsp, tspPerformance.TransportationServiceProviderID); err == nil {
			fmt.Printf("\tAttempting to award to TSP: %v. \n", tsp.Name)
			fmt.Printf("\tQuerying TSP %v for blackout dates for blackout dates.\n", tsp.Name)
			if tspBlackoutDatesPresent == true {
				shipmentAward, err = models.CreateShipmentAward(aq.db, shipment.ID, tsp.ID, true)
				fmt.Printf("\tFailed to award to TSP: %v\n", err)
			} else {
				shipmentAward, err = models.CreateShipmentAward(aq.db, shipment.ID, tsp.ID, false)
				if err == nil {
					fmt.Print("\tShipment awarded to TSP!\n")
					break
				} else {
					fmt.Printf("\tFailed to award to TSP: %v\n", err)
				}
			}
		}
	}
	return shipmentAward, err
}

func (aq *AwardQueue) assignUnawardedShipments() {
	fmt.Println("TSP Award Queue running.")

	shipments, err := aq.findAllUnawardedShipments()
	if err == nil {
		count := 0
		for _, shipment := range shipments {
			_, err = aq.attemptShipmentAward(shipment)
			if err != nil {
				fmt.Printf("\tFailed to award shipment: %s\n", err)
			} else {
				count++
			}
		}
		fmt.Printf("Awarded %d shipments.\n", count)
	} else {
		fmt.Printf("Failed to query for shipments: %s", err)
	}
}

// getTSPsPerBand determines how many TSPs should be assigned to each Quality Band
// If the number of TSPs in the TDL does not divide evenly into 4 bands, the remainder
// is divided from the top band down.
//
// count is the number of TSPs to distribute.
func getTSPsPerBand(count int) []int {
	bands := make([]int, numQualBands)
	base := int(math.Floor(float64(count) / float64(numQualBands)))
	for i := range bands {
		bands[i] = base
	}

	for i := 0; i < count%numQualBands; i++ {
		bands[i]++
	}
	return bands
}

// assignPerformanceBands loops through each TDL and assigns any
// TransportationServiceProviderPerformances without a quality band to a band.
func (aq *AwardQueue) assignPerformanceBands() error {

	// for each TDL with pending performances
	tdls, err := models.FetchTDLsAwaitingBandAssignment(aq.db)
	if err != nil {
		return err
	}

	for _, tdl := range tdls {
		if err := aq.assignPerformanceBandsForTDL(tdl); err != nil {
			return err
		}
	}
	return nil
}

// assignPerformanceBandsForTDL loops through a TDL's TransportationServiceProviderPerformances
// and assigns a QualityBand to each one.
//
// This assumes that all TransportationServiceProviderPerformances have been properly
// created and have a valid BestValueScore.
func (aq *AwardQueue) assignPerformanceBandsForTDL(tdl models.TrafficDistributionList) error {
	fmt.Printf("Assigning performance bands for TDL %s\n", tdl.ID)

	perfs, err := models.FetchTSPPerformanceForQualityBandAssignment(aq.db, tdl.ID, mps)
	if err != nil {
		return err
	}

	perfsIndex := 0
	bands := getTSPsPerBand(len(perfs))

	for band, count := range bands {
		for i := 0; i < count; i++ {
			performance := perfs[perfsIndex]
			fmt.Printf("Assigning tspp %s to band %d\n", performance.ID, band+1)
			err := models.AssignQualityBandToTSPPerformance(aq.db, band+1, performance.ID)
			if err != nil {
				return err
			}
			perfsIndex++
		}
	}
	return nil
}

// NewAwardQueue creates a new AwardQueue
func NewAwardQueue(db *pop.Connection) *AwardQueue {
	return &AwardQueue{db: db}
}

// Run will execute the award queue algorithm.
func Run(db *pop.Connection) error {
	queue := NewAwardQueue(db)

	if err := queue.assignPerformanceBands(); err != nil {
		return err
	}

	// This method should also return an error
	queue.assignUnawardedShipments()
	return nil
}

// checkTSPBlackoutDates searches the blackout_dates table by TSP ID and then compares
// start_blackout_date and end_blackout_date to a submitted pickup date to see if it falls
// within the window created by the blackout date record.
func (aq *AwardQueue) checkTSPBlackoutDates(tspid uuid.UUID, pickupDate time.Time) bool {
	blackoutDates, err := models.FetchTSPBlackoutDates(aq.db, tspid)

	if err != nil {
		fmt.Println("Error retrieving blackout dates.")
	}

	if len(blackoutDates) == 0 {
		return false
	}

	// Checks to see if pickupDate is equal to the start or end dates of the blackout period
	// or if the pickupDate falls between the start and end.
	for _, blackoutDate := range blackoutDates {
		if (pickupDate.After(blackoutDate.StartBlackoutDate) && pickupDate.Before(blackoutDate.EndBlackoutDate)) ||
			pickupDate.Equal(blackoutDate.EndBlackoutDate) ||
			pickupDate.Equal(blackoutDate.StartBlackoutDate) {
			return true
		}
	}
	// Needs a return at the end; feels unsafe to return false (or true!) outside of tests.
	// Is the solution in the vein of "if err != nil"?
	return false // added for debugging; do not keep this in a PR.
}
