package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/unit"
)

// PPMDiscountFetch attempts to fetch the discount rates first for COS D, then 2
// Most PPMs use COS D, but when there is no COS D rate, the calculation is based on Code 2
func PPMDiscountFetch(db *pop.Connection, logger *zap.Logger, originZip string, destZip string, moveDate time.Time) (unit.DiscountRate, unit.DiscountRate, error) {
	// Try to fetch with COS D.
	lhDiscount, sitDiscount, err := FetchDiscountRates(db,
		originZip,
		destZip,
		"D",
		moveDate)

	if err == nil {
		logger.Info("Found Discount for TDL with COS D.",
			zap.String("origin_zip", originZip),
			zap.String("destination_zip", destZip),
			zap.Time("move_date", moveDate),
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
			zap.String("origin_zip", originZip),
			zap.String("destination_zip", destZip),
			zap.Time("move_date", moveDate),
		)
		return lhDiscount, sitDiscount, err
	}

	logger.Error("Couldn't find Discount for COS D or 2.",
		zap.String("origin_zip", originZip),
		zap.String("destination_zip", destZip),
		zap.Time("move_date", moveDate),
		zap.Error(err),
	)
	return 0, 0, err
}
