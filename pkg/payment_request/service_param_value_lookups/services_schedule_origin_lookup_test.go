package serviceparamvaluelookups

import (
	"fmt"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestServicesScheduleOrigin() {
	key := models.ServiceItemParamNameServicesScheduleOrigin.String()

	suite.T().Run("lookup ServicesScheduleOrigin", func(t *testing.T) {
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
				ServicesSchedule: 2,
			},
		})

		testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            domesticServiceArea.Contract,
				DomesticServiceArea: domesticServiceArea,
				Zip3:                "350",
			},
		})

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("2", valueStr)
	})

	suite.T().Run("lookup ServicesScheduleOrigin not found", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				PostalCode: "45007",
			},
		})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				MoveTaskOrder: mtoServiceItem.MoveTaskOrder,
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Equal("", valueStr)
		suite.Error(err)
		expected := fmt.Sprintf(" with error unable to find domestic service area for 450 under contract code %s", ghcrateengine.DefaultContractCode)
		suite.Contains(err.Error(), expected)
	})

	suite.T().Run("nil PickupAddress ID", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				MoveTaskOrder: mtoServiceItem.MoveTaskOrder,
			})

		mtoServiceItem.MTOShipment.PickupAddress = nil
		mtoServiceItem.MTOShipment.PickupAddressID = nil
		suite.MustSave(&mtoServiceItem.MTOShipment)

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.Contains(err.Error(), "looking for PickupAddressID")
		suite.Equal("", valueStr)
	})
}
