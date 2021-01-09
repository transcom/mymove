package testdatagen

import (
	"log"
	"time"

	"github.com/gobuffalo/pop/v5"
	//"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// SetupServiceAreaRateArea sets up contract, service area, rate area, zip3
// returns contractYear, serviceArea, rateArea, reZip3
func SetupServiceAreaRateArea(db *pop.Connection, serviceAreaStr string, address models.Address) (models.ReContractYear, models.ReDomesticServiceArea, models.ReRateArea, models.ReZip3) {
	contractYear := MakeReContractYear(db,
		Assertions{
			ReContractYear: models.ReContractYear{
				Escalation:           1.0197,
				EscalationCompounded: 1.04071,
				StartDate:            time.Date(GHCTestYear, time.January, 1, 0, 0, 0, 0, time.UTC),
				EndDate:              time.Date(GHCTestYear, time.December, 31, 0, 0, 0, 0, time.UTC),
			},
		})

	serviceArea := MakeReDomesticServiceArea(db,
		Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				Contract:    contractYear.Contract,
				ServiceArea: serviceAreaStr,
			},
		})

	rateArea := models.ReRateArea{
		ContractID: contractYear.Contract.ID,
		IsOconus:   false,
		Code:       "US47",
		Name:       address.State,
		Contract:   contractYear.Contract,
	}
	mustSave(db, &rateArea)
	err := db.Q().Where("code = ?", "US47").First(&rateArea)
	if err != nil {
		log.Panic(err)
	}

	reZip3 := MakeReZip3(db, Assertions{
		ReZip3: models.ReZip3{
			ContractID:            contractYear.Contract.ID,
			Zip3:                  address.PostalCode[0:3],
			BasePointCity:         address.City,
			State:                 address.State,
			DomesticServiceAreaID: serviceArea.ID,
			RateAreaID:            &rateArea.ID,
			HasMultipleRateAreas:  false,
			Contract:              contractYear.Contract,
			DomesticServiceArea:   serviceArea,
			RateArea:              &rateArea,
		},
	})

	return contractYear, serviceArea, rateArea, reZip3
}
