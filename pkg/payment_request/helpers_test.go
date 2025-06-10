package paymentrequest

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

// This test relies on non-truncated, default DB data from re services
func (suite *PaymentRequestHelperSuite) TestResolveReServiceForLookup() {
	suite.Run("MTOServiceItem with a non-empty code (and no swap) returns same ReService code", func() {
		dlhService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDLH)

		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    dlhService,
				LinkOnly: true,
			},
		}, nil)

		reService, err := resolveReServiceForLookup(suite.AppContextForTest(), mtoServiceItem)
		suite.NoError(err)
		suite.Equal(models.ReServiceCodeDLH, reService.Code)
	})

	suite.Run("MTOServiceItem with empty code fails", func() {
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, nil)
		mtoServiceItem.ReService.Code = ""

		_, err := resolveReServiceForLookup(suite.AppContextForTest(), mtoServiceItem)
		suite.Error(err)
	})

	suite.Run("MTOServiceItem with INPK code returns IHPK", func() {
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{Code: models.ReServiceCodeINPK},
			},
		}, nil)

		reService, err := resolveReServiceForLookup(suite.AppContextForTest(), mtoServiceItem)
		suite.NoError(err)
		suite.Equal(models.ReServiceCodeIHPK, reService.Code)
	})
}
