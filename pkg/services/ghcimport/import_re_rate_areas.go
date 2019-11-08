package ghcimport

import (
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
)

type RERateAreasImporter struct {
	db     *pop.Connection
	logger Logger
}

func (re RERateAreasImporter) Import() error {

	err := re.importDomesticRateAreas()
	if err != nil {
		return err
	}

	err = re.importInternationalRateAreas()
	if err != nil {
		return err
	}

	return nil
}

func (re RERateAreasImporter) Description() string {
	return "re_rate_area importer"
}

func (re RERateAreasImporter) importDomesticRateAreas() error {
	return nil
}

func (re RERateAreasImporter) importInternationalRateAreas() error {

	// models.StageInternationalServiceArea

	var serviceAreas []models.StageInternationalServiceArea
	/*
			query :=
				`select price_millicents, escalation_compounded
		         from re_domestic_linehaul_prices dlp
		         inner join re_contracts c on dlp.contract_id = c.id
		         inner join re_contract_years cy on c.id = cy.contract_id
		         inner join re_domestic_service_areas dsa on dlp.domestic_service_area_id = dsa.id
		         where c.code = $1
		         and $2 between cy.start_date and cy.end_date
		         and dlp.is_peak_period = $3
		         and $4 between dlp.weight_lower and dlp.weight_upper
		         and $5 between dlp.miles_lower and dlp.miles_upper
		         and dsa.service_area = $6;`
			err := p.db.RawQuery(
				query,
				p.contractCode,
				moveDate,
				isPeakPeriod,
				effectiveWeight,
				distance,
				serviceArea).First(&pe)
	*/

	query := `SELECT * FROM stage_international_service_areas`
	err := re.db.RawQuery(query).All(&serviceAreas)
	if err != nil {
		return errors.Wrap(err, "")
	}

	var rateAreaExistMap map[string]bool
	for _, sa := range serviceAreas {
		if _, ok := rateAreaExistMap[sa.RateAreaID]; !ok {
			// query for ReRateArea
			rateArea, err := models.FetchReRateAreaItem(re.db, sa.RateAreaID)
			if err != nil {
				return errors.Wrapf(err, "Failed importing re_rate_area with code <%s>", sa.RateAreaID)
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
					verrs, err := re.db.ValidateAndSave(rateArea)
					if err != nil || verrs.HasAny() {
						var dbError string
						if err != nil {
							dbError = err.Error()
						}
						if verrs.HasAny() {
							dbError = dbError + verrs.Error()
						}
						return errors.Wrapf(errors.New(dbError), "error saving ReRateArea with rate are ID: %s"+sa.RateAreaID)
					}
				}
			}
			// if it does not exist, insert into ReRateArea
			if rateArea == nil {
				// insert into re_rate_area
				newRateArea := models.ReRateArea{
					IsOconus: true,
					Code:     sa.RateAreaID,
					Name:     sa.RateArea,
				}
				verrs, err := re.db.ValidateAndCreate(&newRateArea)
				if err != nil || verrs.HasAny() {
					var dbError string
					if err != nil {
						dbError = err.Error()
					}
					if verrs.HasAny() {
						dbError = dbError + verrs.Error()
					}
					return errors.Wrapf(errors.New(dbError), "error creating ReRateArea with rate are ID: %s"+sa.RateAreaID)
				}

				// add to map
				rateAreaExistMap[rateArea.Code] = true
			}
		}
	}

	return nil
}
