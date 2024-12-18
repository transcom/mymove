package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *FactorySuite) TestFetchOrMakeDomesticOtherPrice() {
	suite.Run("Successful fetch of domestic other price", func() {

		id, err := uuid.FromString("51393fa4-b31c-40fe-bedf-b692703c46eb")
		suite.NoError(err)
		reService := FetchReServiceByCode(suite.DB(), models.ReServiceCodeDLH)

		domesticOtherPrice := FetchOrMakeDomesticOtherPrice(suite.DB(), []Customization{
			{
				Model: models.ReDomesticOtherPrice{
					ContractID:   id,
					ServiceID:    reService.ID,
					IsPeakPeriod: true,
					Schedule:     1,
					PriceCents:   unit.Cents(945),
				},
			},
		}, nil)
		suite.NotNil(domesticOtherPrice)
	})
}
