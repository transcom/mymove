package serviceparamvaluelookups

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestServiceAreaLookup() {
	originKey := models.ServiceItemParamNameServiceAreaOrigin
	destKey := models.ServiceItemParamNameServiceAreaDest

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
			ServiceArea: "004",
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
			Contract:    originDomesticServiceArea.Contract,
			ServiceArea: "042",
		},
	})

	testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
		ReZip3: models.ReZip3{
			Contract:            destDomesticServiceArea.Contract,
			DomesticServiceArea: destDomesticServiceArea,
			Zip3:                "450",
		},
	})

	suite.T().Run("origin golden path", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(originKey)
		suite.FatalNoError(err)
		suite.Equal(originDomesticServiceArea.ServiceArea, valueStr)
	})

	suite.T().Run("destination golden path", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(destKey)
		suite.FatalNoError(err)
		suite.Equal(destDomesticServiceArea.ServiceArea, valueStr)
	})
}
