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

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

const awardQueueLockID = 1
const numQualBands = 4

// Minimum Performance Score (MPS) is the lowest BVS a TSP can have and still be assigned shipments.
// TODO: This will eventually need to be configurable; implement as something other than a constant.
//       Setting to zero for now to make sure that no TSPs are accidentally excluded.
const mps = 0

// AwardQueue encapsulates the TSP award queue process
type AwardQueue struct {
	db     *pop.Connection
	logger *zap.Logger
}

func (aq *AwardQueue) findAllUnassignedShipments() (models.Shipments, error) {
	shipments, err := models.FetchUnofferedShipments(aq.db)
	return shipments, err
}

// attemptShipmentOffer will attempt to take the given Shipment and award it to
// a TSP.
// TODO: refactor this method to ensure the transaction is wrapping what it needs to
func (aq *AwardQueue) attemptShipmentOffer(shipment models.Shipment) (*models.ShipmentOffer, error) {
	aq.logger.Info("Attempting to offer shipment",
		zap.Any("shipment_id", shipment.ID),
		zap.Any("traffic_distribution_list_id", shipment.TrafficDistributionListID))

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
	firstEligibleTSPPerformance, err := models.NextEligibleTSPPerformance(aq.db, tdl.ID, *shipment.BookDate,
		*shipment.RequestedPickupDate)
	if err != nil {
		return nil, err
	}

	firstTSPid := firstEligibleTSPPerformance.ID
	tspPerformance := firstEligibleTSPPerformance
	foundAvailableTSP := false
	loopCount := 0

	for !foundAvailableTSP {

		if loopCount != 0 && tspPerformance.ID == firstTSPid {
			return nil, fmt.Errorf("could not find a TSP without blackout dates in %d tries", loopCount)
		}
		loopCount++

		tsp := models.TransportationServiceProvider{}
		if err := aq.db.Find(&tsp, tspPerformance.TransportationServiceProviderID); err == nil {
			aq.logger.Info("Attempting to offer to TSP", zap.Object("tsp", tsp))

			isAdministrativeShipment, err := aq.ShipmentWithinBlackoutDates(tsp.ID, shipment)
			if err != nil {
				aq.logger.Error("Failed to determine if shipment is within TSP blackout dates", zap.Error(err))
				return nil, err
			}

			shipmentOffer, err = models.CreateShipmentOffer(aq.db, shipment.ID, tsp.ID, tspPerformance.ID, isAdministrativeShipment)
			if err == nil {
				if tspPerformance, err = models.IncrementTSPPerformanceOfferCount(aq.db, tspPerformance.ID); err == nil {
					if isAdministrativeShipment == true {
						aq.logger.Info("Shipment pickup date is during a blackout period. Awarding Administrative Shipment to TSP.")
					} else {
						qb := -1
						if tspPerformance.QualityBand != nil {
							qb = *tspPerformance.QualityBand
						}

						aq.logger.Info("Shipment offered to TSP!",
							zap.Int("quality_band", qb),
							zap.Int("offer_count", tspPerformance.OfferCount))
						foundAvailableTSP = true

						// Award the shipment
						if err := models.AwardShipment(aq.db, shipment.ID); err != nil {
							aq.logger.Error("Failed to set shipment as awarded",
								zap.Stringer("shipment ID", shipment.ID), zap.Error(err))
							return nil, err
						}
					}
				} else {
					aq.logger.Error("Failed to increment offer count", zap.Error(err))
				}
			} else {
				aq.logger.Error("Failed to offer to TSP", zap.Error(err))
			}
		}

		if !foundAvailableTSP {
			aq.logger.Info("Selected TSP has blackouts. Checking for another TSP.")

			tspPerformance, err = models.NextEligibleTSPPerformance(aq.db, tdl.ID, *shipment.BookDate,
				*shipment.RequestedPickupDate)
			if err != nil {
				return nil, err
			}
		}
	}

	if shipmentOffer == nil {
		err = fmt.Errorf("shipment not awarded; no TSPs found")
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
				aq.logger.Error("Failed to offer shipment", zap.Error(err))
			} else {
				count++
			}
		}
		aq.logger.Info("Awarded some shipments.", zap.Int("total_count", count))
	} else {
		aq.logger.Error("Failed to query for shipments", zap.Error(err))
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

// assignPerformanceBands loops through each unique TransportationServiceProviderPerformances group
// and assigns any unbanded TransportationServiceProviderPerformances to a band.
func (aq *AwardQueue) assignPerformanceBands() error {

	perfGroups, err := models.FetchUnbandedTSPPerformanceGroups(aq.db)
	if err != nil {
		return err
	}

	for _, perfGroup := range perfGroups {
		if err := aq.assignPerformanceBandsForTSPPerformanceGroup(perfGroup); err != nil {
			return err
		}
	}

	return nil
}

// assignPerformanceBandsForTSPPerformanceGroup loops through the TSPPs for a given TSPP grouping
// and assigns a QualityBand to each one.
//
// This assumes that all TransportationServiceProviderPerformances have been properly created and
// have a valid BestValueScore.
func (aq *AwardQueue) assignPerformanceBandsForTSPPerformanceGroup(perfGroup models.TSPPerformanceGroup) error {
	aq.logger.Info("Assigning performance bands",
		zap.Any("traffic_distribution_list_id", perfGroup.TrafficDistributionListID),
		zap.Any("performance_period_start", perfGroup.PerformancePeriodStart),
		zap.Any("performance_period_end", perfGroup.PerformancePeriodEnd),
		zap.Any("rate_cycle_start", perfGroup.RateCycleStart),
		zap.Any("rate_cycle_end", perfGroup.RateCycleEnd),
	)

	perfs, err := models.FetchTSPPerformancesForQualityBandAssignment(aq.db, perfGroup, mps)
	if err != nil {
		return err
	}

	perfsIndex := 0
	bands := getTSPsPerBand(len(perfs))

	for band, count := range bands {
		for i := 0; i < count; i++ {
			performance := perfs[perfsIndex]
			aq.logger.Info("Assigning tspPerformance to band", zap.Any("tsp_performance_id", performance.ID), zap.Int("band", band+1))
			err := models.AssignQualityBandToTSPPerformance(aq.db, band+1, performance.ID)
			if err != nil {
				return err
			}
			perfsIndex++
		}
	}
	return nil
}

// ShipmentWithinBlackoutDates searches the blackout_dates table by TSP ID and shipment details
// to see if it falls within the window created by the blackout date record and if it matches on
// optional fields COS, channel, GBLOC, and market.
func (aq *AwardQueue) ShipmentWithinBlackoutDates(tspID uuid.UUID, shipment models.Shipment) (bool, error) {
	blackoutDates, err := models.FetchTSPBlackoutDates(aq.db, tspID, shipment)

	if err != nil {
		return false, errors.Wrap(err, "Error retrieving blackout dates from database")
	}

	return len(blackoutDates) != 0, nil
}

// Run will execute the award queue algorithm.
func (aq *AwardQueue) Run() error {
	originalDB := aq.db
	defer func() { aq.db = originalDB }()
	return aq.db.Transaction(func(tx *pop.Connection) error {
		// ensure that all parts of the AQ run inside the transaction
		aq.db = tx

		aq.logger.Info("Waiting to acquire advisory lock...")
		if err := waitForLock(tx, awardQueueLockID); err != nil {
			return err
		}
		aq.logger.Info("Acquired pg_advisory_xact_lock")

		// TODO: for testing purposes, will be removed shortly
		// time.Sleep(time.Second * 10)

		if err := aq.assignPerformanceBands(); err != nil {
			return err
		}

		// This method should also return an error
		aq.assignShipments()

		return nil
	})
}

// waitForLock MUST be called within a transaction!
func waitForLock(db *pop.Connection, id int) error {
	// obtain transaction-level advisory-lock
	return db.RawQuery("SELECT pg_advisory_xact_lock($1)", id).Exec()
}

// NewAwardQueue creates a new AwardQueue
func NewAwardQueue(db *pop.Connection, logger *zap.Logger) *AwardQueue {
	return &AwardQueue{db: db, logger: logger}
}
