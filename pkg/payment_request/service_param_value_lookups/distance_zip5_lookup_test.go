package serviceparamvaluelookups

import (
	"strconv"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestDistanceZip5Lookup() {
	key := models.ServiceItemParamNameDistanceZip5.String()

	suite.T().Run("golden path", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				MoveTaskOrder: mtoServiceItem.MoveTaskOrder,
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultDistance)
		suite.Equal(expected, distanceStr)
	})
}
