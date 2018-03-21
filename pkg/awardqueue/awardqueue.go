// Package awardqueue implements the Award Queue mechanism as defined in the
// following document:
// https://docs.google.com/document/d/1WEQZya_yVvW6xbPS7j0-7DP8XSoz9DOntLzIv0FAUHM/edit#
//
// Note on terminology: while the system is referred to as the "award queue"
// it is technically awarding "offers" to TSPs, who then need to accept the
// offer.
package awardqueue

import (
	"fmt"
	"math"
	"time"

	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

const numQualBands = 4

// Minimum Performance Score (MPS) is the lowest BVS a TSP can have and still be assigned shipments.
// TODO: This will eventually need to be configurable; implement as something other than a constant.
const mps = 10

// AwardQueue encapsulates the TSP award queue process
type AwardQueue struct {
	db     *pop.Connection
	logger *zap.SugaredLogger
}

func (aq *AwardQueue) findAllUnassignedShipments() ([]models.ShipmentWithOffer, error) {
	shipments, err := models.FetchShipments(aq.db, true)
	return shipments, err
}

// attemptShipmentOffer will attempt to take the given Shipment and award it to
// a TSP.
// TODO: refactor this method to ensure the transaction is wrapping what it needs to
func (aq *AwardQueue) attemptShipmentOffer(shipment models.ShipmentWithOffer) (*models.ShipmentOffer, error) {
	aq.logger.Infof("Attempting to offer shipment: %s", shipment.ID)

	// Query the shipment's TDL
	tdl := models.TrafficDistributionList{}
	err := aq.db.Find(&tdl, shipment.TrafficDistributionListID)

	if err != nil {
		return nil, errors.Wrap(err, "Cannot find TDL in database")
	}

	var shipmentOffer *models.ShipmentOffer

	// We need to loop here, because if a TSP has a blackout date we need to try again.
	// We _also_ want to watch out for infinite loops, because if all the TSPs in the selection
	// have blackout dates (imagine a 1-TSP-TDL, with a blackout date) we will keep awarding
	// administrative shipments forever.
	firstEligibleTSPPerformance, err := models.NextEligibleTSPPerformance(aq.db, tdl.ID, shipment.BookDate)
	firstTSPid := firstEligibleTSPPerformance.ID
	foundAvailableTSP := false
	loopCount := 0

	for !foundAvailableTSP {

		tspPerformance, err := models.NextEligibleTSPPerformance(aq.db, tdl.ID, shipment.BookDate)

		if loopCount != 0 && tspPerformance.ID == firstTSPid {
			return nil, fmt.Errorf("Could not find a TSP without blackout dates in %d tries", loopCount)
		}
		loopCount++
		if err != nil {
			return nil, err
		}

		err = aq.db.Transaction(func(tx *pop.Connection) error {
			tsp := models.TransportationServiceProvider{}
			if err := aq.db.Find(&tsp, tspPerformance.TransportationServiceProviderID); err == nil {
				aq.logger.Infof("Attempting to offer to TSP: %s", tsp.Name)

				isAdministrativeShipment, err := aq.ShipmentWithinBlackoutDates(tsp.ID, shipment.PickupDate)
				if err != nil {
					return err
				}

				shipmentOffer, err = models.CreateShipmentOffer(aq.db, shipment.ID, tsp.ID, isAdministrativeShipment)
				if err == nil {
					if err = models.IncrementTSPPerformanceOfferCount(aq.db, tspPerformance.ID); err == nil {
						if isAdministrativeShipment == true {
							aq.logger.Info("Shipment pickup date is during a blackout period. Awarding Administrative Shipment to TSP.")
						} else {
							// TODO: OfferCount is off by 1
							aq.logger.Infof("Shipment offered to TSP! TSP now has %d shipment offers.", tspPerformance.OfferCount+1)
							foundAvailableTSP = true
						}
						return nil
					}
				} else {
					aq.logger.Errorf("Failed to offer to TSP: %s", err)
				}
			}

			aq.logger.Errorf("Failed to offer to TSP: %s", err)
			return err
		})

		if !foundAvailableTSP {
			aq.logger.Info("Checking for another TSP.")
		}
	}

	return shipmentOffer, err
}

// assignShipments searches for all shipments that haven't been offered
// yet to a TSP, and attempts to generate offers for each of them.
func (aq *AwardQueue) assignShipments() {
	aq.logger.Info("TSP Award Queue running.")

	shipments, err := aq.findAllUnassignedShipments()
	if err == nil {
		count := 0
		for _, shipment := range shipments {
			_, err = aq.attemptShipmentOffer(shipment)
			if err != nil {
				aq.logger.Errorf("Failed to offer shipment: %s", err)
			} else {
				count++
			}
		}
		aq.logger.Infof("Awarded %d shipments.", count)
	} else {
		aq.logger.Errorf("Failed to query for shipments %s", err)
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
	aq.logger.Infof("Assigning performance bands for TDL %s", tdl.ID)

	perfs, err := models.FetchTSPPerformanceForQualityBandAssignment(aq.db, tdl.ID, mps)
	if err != nil {
		return err
	}

	perfsIndex := 0
	bands := getTSPsPerBand(len(perfs))

	for band, count := range bands {
		for i := 0; i < count; i++ {
			performance := perfs[perfsIndex]
			aq.logger.Infof("Assigning tspp %s to band %d", performance.ID, band+1)
			err := models.AssignQualityBandToTSPPerformance(aq.db, band+1, performance.ID)
			if err != nil {
				return err
			}
			perfsIndex++
		}
	}
	return nil
}

// Run will execute the award queue algorithm.
func (aq *AwardQueue) Run() error {
	if err := aq.assignPerformanceBands(); err != nil {
		return err
	}

	// This method should also return an error
	aq.assignShipments()
	return nil
}

// ShipmentWithinBlackoutDates searches the blackout_dates table by TSP ID and then compares start_blackout_date and end_blackout_date
// to a submitted pickup date to see if it falls within the window created by the blackout date record.
func (aq *AwardQueue) ShipmentWithinBlackoutDates(tspID uuid.UUID, pickupDate time.Time) (bool, error) {
	blackoutDates, err := models.FetchTSPBlackoutDates(aq.db, tspID)

	if err != nil {
		return false, errors.Wrap(err, "Error retrieving blackout dates from database")
	}

	if len(blackoutDates) == 0 {
		return false, nil
	}

	// Checks to see if pickupDate is equal to the start or end dates of the blackout period
	// or if the pickupDate falls between the start and end.
	for _, blackoutDate := range blackoutDates {
		aq.logger.Debugf("Evaluating whether pickup date is between blackout dates (%s <= %s <= %s)", blackoutDate.StartBlackoutDate, pickupDate, blackoutDate.EndBlackoutDate)

		if (pickupDate.After(blackoutDate.StartBlackoutDate) && pickupDate.Before(blackoutDate.EndBlackoutDate)) ||
			pickupDate.Equal(blackoutDate.EndBlackoutDate) ||
			pickupDate.Equal(blackoutDate.StartBlackoutDate) {
			return true, nil
		}
	}

	return false, nil
}

// NewAwardQueue creates a new AwardQueue
func NewAwardQueue(db *pop.Connection, logger *zap.Logger) *AwardQueue {
	return &AwardQueue{db: db, logger: logger.Sugar()}
}
