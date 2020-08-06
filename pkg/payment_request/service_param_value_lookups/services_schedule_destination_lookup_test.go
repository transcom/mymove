package serviceparamvaluelookups

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestServicesScheduleDestination() {
	key := models.ServiceItemParamNameServicesScheduleDest.String()

	suite.T().Run("lookup ServicesScheduleDestination", func(t *testing.T) {
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
				ServicesSchedule: 3,
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
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("3", valueStr)
	})

	suite.T().Run("lookup ServicesScheduleDestination not found", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				PostalCode: "45007",
			},
		})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				MoveTaskOrder: mtoServiceItem.MoveTaskOrder,
			})

		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Equal("", valueStr)
		suite.Error(err)
		expected := fmt.Sprintf(" with error unable to find domestic service area for 450 under contract code %s", ghcrateengine.DefaultContractCode)
		suite.Contains(err.Error(), expected)
	})

	suite.T().Run("nil DestinationAddress ID", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				MoveTaskOrder: mtoServiceItem.MoveTaskOrder,
			})

		mtoServiceItem.MTOShipment.DestinationAddress = nil
		mtoServiceItem.MTOShipment.DestinationAddressID = nil
		suite.MustSave(&mtoServiceItem.MTOShipment)

		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find destination address for MTOShipment [%s]", mtoServiceItem.MTOShipment.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)
	})

	suite.T().Run("nil MTOShipment ID", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				PostalCode: "45007",
			},
		})
		mtoServiceItem.MTOShipmentID = nil
		suite.MustSave(&mtoServiceItem)

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				MoveTaskOrder: mtoServiceItem.MoveTaskOrder,
			})

		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		suite.Equal("", valueStr)
	})

	suite.T().Run("nil MTOServiceItem ID", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				PostalCode: "45007",
			},
		})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				MoveTaskOrder: mtoServiceItem.MoveTaskOrder,
			})

		// Pass in a non-existent MTOServiceItemID
		invalidMTOServiceItemID := uuid.Must(uuid.NewV4())
		badParamLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, invalidMTOServiceItemID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

		valueStr, err := badParamLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		suite.Equal("", valueStr)
	})
}
