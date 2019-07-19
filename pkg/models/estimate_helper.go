package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/unit"
)

type ppmDiscountFetchParams struct {
	OriginZip string
	DestZip   string
	Date      time.Time
	Cos       string
	RunFetch  bool
}

func tryPPMDiscountFetch(db *pop.Connection, logger Logger, fetchParams ppmDiscountFetchParams) (unit.DiscountRate, unit.DiscountRate, error) {

	// Try to fetch
	lhDiscount, sitDiscount, err := FetchDiscountRates(db,
		fetchParams.OriginZip,
		fetchParams.DestZip,
		fetchParams.Cos,
		fetchParams.Date)

	if err == nil {
		logger.Info("Found Discount for TDL with",
			zap.String("COS", fetchParams.Cos),
			zap.String("origin_zip", fetchParams.OriginZip),
			zap.String("destination_zip", fetchParams.DestZip),
			zap.Time("date", fetchParams.Date),
			zap.Float64("lh_discount", lhDiscount.Float64()),
			zap.Float64("sit_discount", sitDiscount.Float64()),
		)
		return lhDiscount, sitDiscount, err
	}

	return 0, 0, err
}

// PPMDiscountFetch attempts to fetch the discount rates first for COS D, then 2
// Most PPMs use COS D, but when there is no COS D rate, the calculation is based on Code 2
func PPMDiscountFetch(db *pop.Connection, logger Logger, originZip string, destZip string, moveDate time.Time, bookDate time.Time, allowBookDate bool) (unit.DiscountRate, unit.DiscountRate, error) {

	datesForFetch := []ppmDiscountFetchParams{
		{
			OriginZip: originZip,
			DestZip:   destZip,
			Date:      moveDate,
			Cos:       "D",
			RunFetch:  true,
		},
		{
			OriginZip: originZip,
			DestZip:   destZip,
			Date:      moveDate,
			Cos:       "2",
			RunFetch:  true,
		},
		{
			OriginZip: originZip,
			DestZip:   destZip,
			Date:      bookDate,
			Cos:       "D",
			RunFetch:  allowBookDate,
		},
		{
			OriginZip: originZip,
			DestZip:   destZip,
			Date:      bookDate,
			Cos:       "2",
			RunFetch:  allowBookDate,
		},
	}

	logger.Info("Fetching PPM Discount for TDL with ", zap.Time("move_date", moveDate))
	if allowBookDate {
		logger.Info("PPM Discount for TDL is allowed to use Book Date of move ", zap.Time("book_date", bookDate))
	}

	err := ErrFetchNotFound
	var lhDiscount unit.DiscountRate
	var sitDiscount unit.DiscountRate
	for _, params := range datesForFetch {
		if params.RunFetch {
			if err == ErrFetchNotFound {
				lhDiscount, sitDiscount, err = tryPPMDiscountFetch(db, logger, params)
				if err == nil {
					logger.Info("Found Discount for TDL with",
						zap.String("COS", params.Cos),
						zap.String("origin_zip", originZip),
						zap.String("destination_zip", destZip),
						zap.Time("using date for fetch", params.Date),
						zap.Time("with move_date", moveDate),
						zap.Float64("lh_discount", lhDiscount.Float64()),
						zap.Float64("sit_discount", sitDiscount.Float64()),
					)
					if allowBookDate {
						logger.Info("Found Discount for TDL allowing use of book date", zap.Time("book_date", bookDate))
					}
					return lhDiscount, sitDiscount, err
				}
			}
		}
	}

	logger.Error("Couldn't find Discount for COS D or 2.",
		zap.String("origin_zip", originZip),
		zap.String("destination_zip", destZip),
		zap.Time("move_date", moveDate),
		zap.Error(err),
	)
	if allowBookDate {
		logger.Info("Couldn't find Discount with date including", zap.Time("book_date", bookDate))
	}
	return 0, 0, err
}
