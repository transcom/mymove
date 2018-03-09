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
// TODO: refactor this method to ensure the transaction is wrapping what it needs to
func (aq *AwardQueue) attemptShipmentAward(shipment models.PossiblyAwardedShipment) (*models.ShipmentAward, error) {
	fmt.Printf("Attempting to award shipment: %v\n", shipment.ID)

	// Query the shipment's TDL
	tdl := models.TrafficDistributionList{}
	err := aq.db.Find(&tdl, shipment.TrafficDistributionListID)

	if err != nil {
		return nil, fmt.Errorf("Cannot find TDL in database: %s", err)
	}

	var shipmentAward *models.ShipmentAward

	// TODO: We need to loop here, because if a TSP has a blackout date we need to try again.
	// we _also_ want to watch out for inifite loops, because if all the TSPs in the selection
	// have blackout dates (imagine a 1-TSP-TDL, with a blackout date) we will keep awarding
	// administrative shipments forever.
	foundAvailableTSP := false
	loopCount := 0
	blackoutRetries := 1000

	for !foundAvailableTSP && loopCount < blackoutRetries {
		loopCount++

		tspPerformance, err := models.NextEligibleTSPPerformance(aq.db, tdl.ID)

		if err != nil {
			return nil, fmt.Errorf("Cannot award. Error: %s", err)
		}

		err = aq.db.Transaction(func(tx *pop.Connection) error {
			tsp := models.TransportationServiceProvider{}

			if err := aq.db.Find(&tsp, tspPerformance.TransportationServiceProviderID); err == nil {
				fmt.Printf("\tAttempting to award to TSP: %s\n", tsp.Name)

				tspBlackoutDatesPresent, err := aq.CheckTSPBlackoutDates(tsp.ID, shipment.PickupDate)
				if err == nil {
					shipmentAward, err = models.CreateShipmentAward(aq.db, shipment.ID, tsp.ID, tspBlackoutDatesPresent)
				} else {
					return err
				}

				if err == nil {
					if err = models.IncrementTSPPerformanceAwardCount(aq.db, tspPerformance.ID); err == nil {
						if tspBlackoutDatesPresent == true {
							fmt.Printf("\tShipment pickup date is during a blackout period. Awarding Administrative Shipment to TSP.\n")
						} else {
							// TODO: AwardCount is off by 1
							fmt.Printf("\tShipment awarded to TSP! TSP now has %d shipment awards.\n", tspPerformance.AwardCount+1)
							foundAvailableTSP = true
						}
						return nil
					}
				} else {
					fmt.Printf("\tFailed to award to TSP: %v\n", err)
				}
			}

			fmt.Printf("\tFailed to award to TSP: %v\n", err)
			return err
		})
		if !foundAvailableTSP {
			fmt.Printf("\tChecking for another TSP. Tries left: %d\n", blackoutRetries-loopCount)
		}
	}

	if loopCount == blackoutRetries {
		return nil, fmt.Errorf("Could not find a TSP without blackout dates in %d tries", blackoutRetries)
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

// CheckTSPBlackoutDates searches the blackout_dates table by TSP ID and then compares start_blackout_date and end_blackout_date to a submitted pickup date to see if it falls within the window created by the blackout date record.
func (aq *AwardQueue) CheckTSPBlackoutDates(tspID uuid.UUID, pickupDate time.Time) (bool, error) {
	blackoutDates, err := models.FetchTSPBlackoutDates(aq.db, tspID)

	if err != nil {
		return false, fmt.Errorf("Error retrieving blackout dates from database: %s", err)
	}

	if len(blackoutDates) == 0 {
		return false, nil
	}

	// Checks to see if pickupDate is equal to the start or end dates of the blackout period
	// or if the pickupDate falls between the start and end.
	for _, blackoutDate := range blackoutDates {
		fmt.Printf("Comparing blackout date: %s < %s < %s\n", blackoutDate.StartBlackoutDate, pickupDate, blackoutDate.EndBlackoutDate)
		if (pickupDate.After(blackoutDate.StartBlackoutDate) && pickupDate.Before(blackoutDate.EndBlackoutDate)) ||
			pickupDate.Equal(blackoutDate.EndBlackoutDate) ||
			pickupDate.Equal(blackoutDate.StartBlackoutDate) {
			return true, nil
		}
	}

	return false, nil
}
