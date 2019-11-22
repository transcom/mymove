package ghcimport

import (
	"fmt"

	"github.com/gobuffalo/uuid"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) importRERateArea(dbTx *pop.Connection) error {
	var err error
	//maps the domestic rate areas to a UUID
	gre.domesticRateAreaToIDMap, err = gre.importDomesticRateAreas(dbTx)
	if err != nil {
		return fmt.Errorf("importRERateArea failed to import: %w", err)
	}
	//maps the international rate areas to a UUID
	gre.internationalRateAreaToIDMap, err = gre.importInternationalRateAreas(dbTx)
	if err != nil {
		return fmt.Errorf("importRERateArea failed to import: %w", err)
	}
	return nil
}

func (gre *GHCRateEngineImporter) importDomesticRateAreas(db *pop.Connection) (map[string]uuid.UUID, error) {
	rateAreaToIDMap := make(map[string]uuid.UUID)

	// have to read international tables to get the domestic rate areas

	// models.StageConusToOconusPrice
	var conusToOconus []models.StageConusToOconusPrice
	err := db.All(&conusToOconus)

	if err != nil {
		return nil, fmt.Errorf("failed to query all StageConusToOconusPrice: %w", err)
	}
	for _, ra := range conusToOconus {
		// does the rate area already exist in the rate engine
		var rateArea *models.ReRateArea
		rateArea, err = models.FetchReRateAreaItem(db, ra.OriginDomesticPriceAreaCode)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return nil, fmt.Errorf("failed importing re_rate_area from StageConusToOconusPrice with code: <%s> error: %w", ra.OriginDomesticPriceAreaCode, err)
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
					return nil, fmt.Errorf("error saving ReRateArea from StageConusToOconusPrice with rate are ID: %s error: %w", ra.OriginDomesticPriceAreaCode, errors.New(dbError))
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
				return nil, fmt.Errorf("error creating ReRateArea from StageConusToOconusPrice with rate are ID: %s error: %w", ra.OriginDomesticPriceAreaCode, errors.New(dbError))
			}
			rateAreaToIDMap[ra.OriginDomesticPriceAreaCode] = newRateArea.ID
		}

		// add to map
		//rateAreaToIDMap[ra.OriginDomesticPriceAreaCode] = newRateArea.ID
	}

	// models.StageOconusToConusPrice
	var oconusToConsus []models.StageOconusToConusPrice
	err = db.All(&oconusToConsus)
	if err != nil {
		return nil, fmt.Errorf("failed to query all StageOconusToConusPrice error: %w", err)
	}
	for _, ra := range oconusToConsus {
		// does the rate area already exist in the rate engine
		rateArea, err := models.FetchReRateAreaItem(db, ra.DestinationDomesticPriceAreaCode)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return nil, fmt.Errorf("Failed importing re_rate_area from StageOconusToConusPrice with code <%s> error: %w", ra.DestinationDomesticPriceAreaCode, err)
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
					return nil, fmt.Errorf("error saving ReRateArea from StageOconusToConusPrice with rate are ID: %s error: %w", ra.DestinationDomesticPriceAreaCode, errors.New(dbError))
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
				return nil, fmt.Errorf("error creating ReRateArea from StageOconusToConusPrice with rate are ID: %s error: %w", ra.DestinationDomesticPriceAreaCode, errors.New(dbError))
			}
			rateAreaToIDMap[ra.DestinationDomesticPriceAreaCode] = newRateArea.ID
		}

		// add to map
		//rateAreaToIDMap[ra.DestinationDomesticPriceAreaCode] = newRateArea.ID
	}

	return rateAreaToIDMap, nil
}

func (gre *GHCRateEngineImporter) importInternationalRateAreas(db *pop.Connection) (map[string]uuid.UUID, error) {
	var serviceAreas []models.StageInternationalServiceArea

	err := db.All(&serviceAreas)
	if err != nil {
		return nil, fmt.Errorf("failed to query all StageInternationalServiceArea: %w", err)
	}

	rateAreaToIDMap := make(map[string]uuid.UUID)
	for _, sa := range serviceAreas {
		// query for ReRateArea
		rateArea, err := models.FetchReRateAreaItem(db, sa.RateAreaID)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				return nil, fmt.Errorf("failed importing re_rate_area from StageInternationalServiceArea with code <%s> error: %w", sa.RateAreaID, err)
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
					return nil, fmt.Errorf("error saving ReRateArea from StageInternationalServiceArea with rate are ID: %s error: %w", sa.RateAreaID, errors.New(dbError))
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
				return nil, fmt.Errorf("error creating ReRateArea from StageInternationalServiceArea with rate are ID: %s error: %w", sa.RateAreaID, errors.New(dbError))
			}
			rateAreaToIDMap[sa.RateAreaID] = newRateArea.ID
		}

		// add to map
		//rateAreaToIDMap[sa.RateAreaID] = newRateArea.ID
	}

	return rateAreaToIDMap, nil
}