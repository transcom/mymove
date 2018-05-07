package handlers

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

// PPMDiscountFetch attempts to fetch the discount rates first for COS D, then 2
// Most PPMs use COS D, but when there is no COS D rate, the calculation is based on Code 2
func PPMDiscountFetch(db *pop.Connection, logger *zap.Logger, originZip string, destZip string, moveDate time.Time) (float64, float64, error) {
	lhDiscount, sitDiscount, err := models.FetchDiscountRates(db,
		originZip,
		destZip,
		"D",
		moveDate)

	if err != nil {
		if err != models.ErrFetchNotFound {
			fmt.Println(err)
			return 0, 0, err
		}
		lhDiscount, sitDiscount, err = models.FetchDiscountRates(db,
			originZip,
			destZip,
			"2",
			moveDate)

		if err != nil {
			logger.Info("Couldn't find Discount for COS D or 2.",
				zap.String("origin_zip", originZip),
				zap.String("destination_zip", destZip),
				zap.Time("move_date", moveDate),
				zap.Error(err),
			)
			return 0, 0, err
		}
		logger.Info("Found Discount for TDL with COS 2.",
			zap.String("origin_zip", originZip),
			zap.String("destination_zip", destZip),
			zap.Time("move_date", moveDate),
		)
	} else {
		logger.Info("Found Discount for TDL with COS D.",
			zap.String("origin_zip", originZip),
			zap.String("destination_zip", destZip),
			zap.Time("move_date", moveDate),
		)
	}
	return lhDiscount, sitDiscount, err
}
