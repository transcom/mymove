package serviceparamvaluelookups

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestServiceAreaOrigin() {
	key := models.ServiceItemParamNameServiceAreaOrigin.String()

	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			PostalCode: "35007",
		},
	})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			MoveTaskOrder: mtoServiceItem.MoveTaskOrder,
		})

	domesticServiceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceArea: models.ReDomesticServiceArea{
			ServiceArea: "004",
		},
	})

	testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
		ReZip3: models.ReZip3{
			Contract:            domesticServiceArea.Contract,
			DomesticServiceArea: domesticServiceArea,
			Zip3:                "350",
		},
	})

	paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

	suite.T().Run("golden path", func(t *testing.T) {

		valueStr, err := paramLookup.ServiceParamValue(key)

		suite.FatalNoError(err)
		suite.Equal("004", valueStr)
	})
}
