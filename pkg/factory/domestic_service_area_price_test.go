package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *FactorySuite) TestFetchOrMakeDomesticServiceAreaPrice() {
	suite.Run("Successful fetch of domestic service area price", func() {

		id, err := uuid.FromString("070f7c82-fad0-4ae8-9a83-5de87a56472e")
		suite.NoError(err)
		reService := FetchReServiceByCode(suite.DB(), models.ReServiceCodeDLH)
		serviceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea:      "004",
				ServicesSchedule: 2,
			},
			ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
		})

		domesticServiceAreaPrice := FetchOrMakeDomesticServiceAreaPrice(suite.DB(), []Customization{
			{
				Model: models.ReDomesticServiceAreaPrice{
					ContractID:            id,
					ServiceID:             reService.ID,
					DomesticServiceAreaID: serviceArea.ID,
					IsPeakPeriod:          true,
					PriceCents:            unit.Cents(945),
				},
			},
		}, nil)
		suite.NotNil(domesticServiceAreaPrice)
	})
}
