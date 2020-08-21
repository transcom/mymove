package serviceparamvaluelookups

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestServicesScheduleOrigin() {
	originKey := models.ServiceItemParamNameServicesScheduleOrigin.String()
	destKey := models.ServiceItemParamNameServicesScheduleDest.String()

	originAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			PostalCode: "35007",
		},
	})
	destAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			PostalCode: "45007",
		},
	})

	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PickupAddressID:      &originAddress.ID,
			PickupAddress:        &originAddress,
			DestinationAddressID: &destAddress.ID,
			DestinationAddress:   &destAddress,
		},
	})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			Move: mtoServiceItem.MoveTaskOrder,
		})

	originDomesticServiceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceArea: models.ReDomesticServiceArea{
			ServiceArea:      "004",
			ServicesSchedule: 2,
		},
	})

	testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
		ReZip3: models.ReZip3{
			Contract:            originDomesticServiceArea.Contract,
			DomesticServiceArea: originDomesticServiceArea,
			Zip3:                "350",
		},
	})

	destDomesticServiceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceArea: models.ReDomesticServiceArea{
			Contract:         originDomesticServiceArea.Contract,
			ServiceArea:      "005",
			ServicesSchedule: 3,
		},
	})

	testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
		ReZip3: models.ReZip3{
			Contract:            destDomesticServiceArea.Contract,
			DomesticServiceArea: destDomesticServiceArea,
			Zip3:                "450",
		},
	})

	suite.T().Run("lookup origin ServicesSchedule", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(originKey)
		suite.FatalNoError(err)
		suite.Equal(strconv.Itoa(originDomesticServiceArea.ServicesSchedule), valueStr)
	})

	suite.T().Run("lookup dest ServicesSchedule", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(destKey)
		suite.FatalNoError(err)
		suite.Equal(strconv.Itoa(destDomesticServiceArea.ServicesSchedule), valueStr)
	})

	suite.T().Run("lookup origin ServicesSchedule not found", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				Move: mtoServiceItem.MoveTaskOrder,
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(originKey)
		suite.Equal("", valueStr)
		suite.Error(err)
		expected := fmt.Sprintf(" with error unable to find domestic service area for 902 under contract code %s", ghcrateengine.DefaultContractCode)
		suite.Contains(err.Error(), expected)
	})

	suite.T().Run("lookup dest ServicesSchedule not found", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				Move: mtoServiceItem.MoveTaskOrder,
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(destKey)
		suite.Equal("", valueStr)
		suite.Error(err)
		expected := fmt.Sprintf(" with error unable to find domestic service area for 945 under contract code %s", ghcrateengine.DefaultContractCode)
		suite.Contains(err.Error(), expected)
	})
}
