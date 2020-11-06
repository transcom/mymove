// Package awardqueue implements the Award Queue mechanism as defined in the
// following document:
// https://docs.google.com/document/d/1WEQZya_yVvW6xbPS7j0-7DP8XSoz9DOntLzIv0FAUHM/edit#
//
// Note on terminology: while the system is referred to as the "award queue"
// it is technically awarding "offers" to TSPs, who then need to accept the
// offer.
package awardqueue

import (
	"context"
	"math"

	"github.com/gobuffalo/pop/v5"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

const numQualBands = 4

// Minimum Performance Score (MPS) is the lowest BVS a TSP can have and still be assigned shipments.
// TODO: This will eventually need to be configurable; implement as something other than a constant.
//       Setting to zero for now to make sure that no TSPs are accidentally excluded.
const mps = 0

// AwardQueue encapsulates the TSP award queue process
type AwardQueue struct {
	db     *pop.Connection
	logger Logger
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
func (aq *AwardQueue) assignPerformanceBands(ctx context.Context) error {
	perfGroups, err := models.FetchUnbandedTSPPerformanceGroups(aq.db)
	if err != nil {
		return err
	}

	for _, perfGroup := range perfGroups {
		if err := aq.assignPerformanceBandsForTSPPerformanceGroup(ctx, perfGroup); err != nil {
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
func (aq *AwardQueue) assignPerformanceBandsForTSPPerformanceGroup(ctx context.Context, perfGroup models.TSPPerformanceGroup) error {
	aq.logger.Info("Assigning performance bands",
		zap.String("traffic_distribution_list_id", perfGroup.TrafficDistributionListID.String()),
		zap.String("performance_period_start", perfGroup.PerformancePeriodStart.String()),
		zap.String("performance_period_end", perfGroup.PerformancePeriodEnd.String()),
		zap.String("rate_cycle_start", perfGroup.RateCycleStart.String()),
		zap.String("rate_cycle_end", perfGroup.RateCycleEnd.String()),
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
			aq.logger.Info("Assigning tspPerformance to band", zap.String("tsp_performance_id", performance.ID.String()), zap.Int("band", band+1))
			err := models.AssignQualityBandToTSPPerformance(ctx, aq.db, band+1, performance.ID)
			if err != nil {
				return err
			}
			perfsIndex++
		}
	}
	return nil
}

// waitForLock MUST be called within a transaction!
func waitForLock(ctx context.Context, db *pop.Connection, id int) error {

	// obtain transaction-level advisory-lock
	return db.RawQuery("SELECT pg_advisory_xact_lock($1)", id).Exec()
}

// NewAwardQueue creates a new AwardQueue
func NewAwardQueue(db *pop.Connection, logger Logger) *AwardQueue {
	return &AwardQueue{
		db:     db,
		logger: logger,
	}
}
