package ghcimport

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) importRERateArea(dbTx *pop.Connection) error {
	err := gre.importDomesticRateAreas(dbTx)
	if err != nil {
		return errors.Wrap(err, "importRERateArea failed to import")
	}
	err = gre.importInternationalRateAreas(dbTx)
	if err != nil {
		return errors.Wrap(err, "importRERateArea failed to import")
	}
	return nil
}

func (gre *GHCRateEngineImporter) importDomesticRateAreas(db *pop.Connection) error {

	rateAreaExistMap := make(map[string]bool)

	// have to read international tables to get the domestic rate areas

	// models.StageConusToOconusPrice
	var conusToOconus []models.StageConusToOconusPrice
	err := db.All(&conusToOconus)

	if err != nil {
		return errors.Wrap(err, "Failed to query all StageConusToOconusPrice")
	}
	for _, ra := range conusToOconus {
		if _, ok := rateAreaExistMap[ra.OriginDomesticPriceAreaCode]; !ok {
			// does the rate area already exist in the rate engine
			var rateArea *models.ReRateArea
			rateArea, err = models.FetchReRateAreaItem(db, ra.OriginDomesticPriceAreaCode)
			if err != nil {
				if err.Error() != "sql: no rows in result set" {
					return errors.Wrapf(err, "Failed importing re_rate_area from StageConusToOconusPrice with code <%s>", ra.OriginDomesticPriceAreaCode)
				}
			}

			// if it does exist, compare and update information if different
			if rateArea != nil {
				update := false

				if rateArea.Name != ra.OriginDomesticPriceArea {
					rateArea.Name = ra.OriginDomesticPriceArea
					update = true
				}

				// these are domestic rates
				if rateArea.IsOconus != false {
					rateArea.IsOconus = false
					update = true
				}

				if update {
					var verrs *validate.Errors
					verrs, err = db.ValidateAndSave(rateArea)
					if err != nil || verrs.HasAny() {
						var dbError string
						if err != nil {
							dbError = err.Error()
						}
						if verrs.HasAny() {
							dbError = dbError + verrs.Error()
						}
						return errors.Wrapf(errors.New(dbError), "error saving ReRateArea from StageConusToOconusPrice with rate are ID: %s"+ra.OriginDomesticPriceAreaCode)
					}
				}

				// if it does not exist, insert into ReRateArea
			} else if rateArea == nil {
				// insert into re_rate_area
				newRateArea := models.ReRateArea{
					IsOconus: false,
					Code:     ra.OriginDomesticPriceAreaCode,
					Name:     ra.OriginDomesticPriceArea,
				}
				var verrs *validate.Errors
				verrs, err = db.ValidateAndCreate(&newRateArea)
				if err != nil || verrs.HasAny() {
					var dbError string
					if err != nil {
						dbError = err.Error()
					}
					if verrs.HasAny() {
						dbError = dbError + verrs.Error()
					}
					return errors.Wrapf(errors.New(dbError), "error creating ReRateArea from StageConusToOconusPrice with rate are ID: %s"+ra.OriginDomesticPriceAreaCode)
				}
			}

			// add to map
			rateAreaExistMap[ra.OriginDomesticPriceAreaCode] = true
		}
	}

	// models.StageOconusToConusPrice
	var oconusToConsus []models.StageOconusToConusPrice
	err = db.All(&oconusToConsus)
	if err != nil {
		return errors.Wrap(err, "Failed to query all StageOconusToConusPrice")
	}
	for _, ra := range oconusToConsus {
		if _, ok := rateAreaExistMap[ra.DestinationDomesticPriceAreaCode]; !ok {
			// does the rate area already exist in the rate engine
			rateArea, err := models.FetchReRateAreaItem(db, ra.DestinationDomesticPriceAreaCode)
			if err != nil {
				if err.Error() != "sql: no rows in result set" {
					return errors.Wrapf(err, "Failed importing re_rate_area from StageOconusToConusPrice with code <%s>", ra.DestinationDomesticPriceAreaCode)
				}
			}

			// if it does exist, compare and update information if different
			if rateArea != nil {
				update := false

				if rateArea.Name != ra.DestinationDomesticPriceArea {
					rateArea.Name = ra.DestinationDomesticPriceArea
					update = true
				}

				// these are domestic rates
				if rateArea.IsOconus != false {
					rateArea.IsOconus = false
					update = true
				}

				if update {
					verrs, err := db.ValidateAndSave(rateArea)
					if err != nil || verrs.HasAny() {
						var dbError string
						if err != nil {
							dbError = err.Error()
						}
						if verrs.HasAny() {
							dbError = dbError + verrs.Error()
						}
						return errors.Wrapf(errors.New(dbError), "error saving ReRateArea from StageOconusToConusPrice with rate are ID: %s"+ra.DestinationDomesticPriceAreaCode)
					}
				}

				// if it does not exist, insert into ReRateArea
			} else if rateArea == nil {
				// insert into re_rate_area
				newRateArea := models.ReRateArea{
					IsOconus: false,
					Code:     ra.DestinationDomesticPriceAreaCode,
					Name:     ra.DestinationDomesticPriceArea,
				}
				verrs, err := db.ValidateAndCreate(&newRateArea)
				if err != nil || verrs.HasAny() {
					var dbError string
					if err != nil {
						dbError = err.Error()
					}
					if verrs.HasAny() {
						dbError = dbError + verrs.Error()
					}
					return errors.Wrapf(errors.New(dbError), "error creating ReRateArea from StageOconusToConusPrice with rate are ID: %s"+ra.DestinationDomesticPriceAreaCode)
				}
			}

			// add to map
			rateAreaExistMap[ra.DestinationDomesticPriceAreaCode] = true
		}
	}

	return nil
}

func (gre *GHCRateEngineImporter) importInternationalRateAreas(db *pop.Connection) error {
	var serviceAreas []models.StageInternationalServiceArea

	err := db.All(&serviceAreas)
	if err != nil {
		return errors.Wrap(err, "Failed to query all StageInternationalServiceArea")
	}

	rateAreaExistMap := make(map[string]bool)
	for _, sa := range serviceAreas {
		if _, ok := rateAreaExistMap[sa.RateAreaID]; !ok {
			// query for ReRateArea
			rateArea, err := models.FetchReRateAreaItem(db, sa.RateAreaID)
			if err != nil {
				if err.Error() != "sql: no rows in result set" {
					return errors.Wrapf(err, "Failed importing re_rate_area from StageInternationalServiceArea with code <%s>", sa.RateAreaID)
				}
			}
			// if it does exist, compare and update information if different
			if rateArea != nil {
				update := false

				if rateArea.IsOconus != true {
					rateArea.IsOconus = true
					update = true
				}
				if rateArea.Name != sa.RateArea {
					rateArea.Name = sa.RateArea
					update = true
				}
				if update {
					verrs, err := db.ValidateAndSave(rateArea)
					if err != nil || verrs.HasAny() {
						var dbError string
						if err != nil {
							dbError = err.Error()
						}
						if verrs.HasAny() {
							dbError = dbError + verrs.Error()
						}
						return errors.Wrapf(errors.New(dbError), "error saving ReRateArea from StageInternationalServiceArea with rate are ID: %s"+sa.RateAreaID)
					}
				}
				// if it does not exist, insert into ReRateArea
			} else if rateArea == nil {
				// insert into re_rate_area
				newRateArea := models.ReRateArea{
					IsOconus: true,
					Code:     sa.RateAreaID,
					Name:     sa.RateArea,
				}
				verrs, err := db.ValidateAndCreate(&newRateArea)
				if err != nil || verrs.HasAny() {
					var dbError string
					if err != nil {
						dbError = err.Error()
					}
					if verrs.HasAny() {
						dbError = dbError + verrs.Error()
					}
					return errors.Wrapf(errors.New(dbError), "error creating ReRateArea from StageInternationalServiceArea with rate are ID: %s"+sa.RateAreaID)
				}
			}

			// add to map
			rateAreaExistMap[sa.RateAreaID] = true
		}
	}

	return nil
}