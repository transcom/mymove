package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/unit"
)

// PPMDiscountFetch attempts to fetch the discount rates first for COS D, then 2
// Most PPMs use COS D, but when there is no COS D rate, the calculation is based on Code 2
func PPMDiscountFetch(db *pop.Connection, logger *zap.Logger, move Move, originZip string, destZip string, moveDate time.Time) (unit.DiscountRate, unit.DiscountRate, error) {
	// Try to fetch with COS D.
	lhDiscount, sitDiscount, err := FetchDiscountRates(db,
		originZip,
		destZip,
		"D",
		moveDate)

	if err == nil {
		logger.Info("Found Discount for TDL with COS D.",
			zap.String("moveLocator", move.Locator),
			zap.String("origin_zip", originZip),
			zap.String("destination_zip", destZip),
			zap.Time("move_date", moveDate),
			zap.Float64("lh_discount", lhDiscount.Float64()),
			zap.Float64("sit_discount", sitDiscount.Float64()),
		)
		return lhDiscount, sitDiscount, err
	}

	if err != ErrFetchNotFound {
		return 0, 0, err
	}
	// When COS D not found, COS 2 may have rates.
	lhDiscount, sitDiscount, err = FetchDiscountRates(db,
		originZip,
		destZip,
		"2",
		moveDate)

	if err == nil {
		logger.Info("Found Discount for TDL with COS 2.",
			zap.String("moveLocator", move.Locator),
			zap.String("origin_zip", originZip),
			zap.String("destination_zip", destZip),
			zap.Time("move_date", moveDate),
			zap.Float64("lh_discount", lhDiscount.Float64()),
			zap.Float64("sit_discount", sitDiscount.Float64()),
		)
		return lhDiscount, sitDiscount, err
	}

	logger.Error("Couldn't find Discount for COS D or 2.",
		zap.String("moveLocator", move.Locator),
		zap.String("origin_zip", originZip),
		zap.String("destination_zip", destZip),
		zap.Time("move_date", moveDate),
		zap.Error(err),
	)
	return 0, 0, fmt.Errorf("could not find discount: %w", err)
}
