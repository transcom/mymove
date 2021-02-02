package serviceparamvaluelookups

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestZipSITOriginHHGOriginalAddressLookup() {
	key := models.ServiceItemParamNameZipSITOriginHHGOriginalAddress

	originZip := "30901"
	actualOriginZipSameZip3 := "30907"

	reService := testdatagen.FetchOrMakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: "DOFSIT",
			},
		},
	)

	originAddress := testdatagen.MakeAddress(suite.DB(),
		testdatagen.Assertions{
			Address: models.Address{
				PostalCode: originZip,
			},
		})

	actualOriginSameZip3Address := testdatagen.MakeAddress(suite.DB(),
		testdatagen.Assertions{
			Address: models.Address{
				PostalCode: actualOriginZipSameZip3,
			},
		})

	move := testdatagen.MakeDefaultMove(suite.DB())

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			Move: move,
		})

	mtoServiceItemWithSITOriginZips := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			ReService: reService,
			Move:      move,
			MTOServiceItem: models.MTOServiceItem{
				SITOriginHHGOriginalAddressID: &originAddress.ID,
				SITOriginHHGOriginalAddress:   &originAddress,
				SITOriginHHGActualAddressID:   &actualOriginSameZip3Address.ID,
				SITOriginHHGActualAddress:     &actualOriginSameZip3Address,
			},
		},
	)

	mtoServiceItemNoSITOriginZips := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			ReService:      reService,
			Move:           move,
			MTOServiceItem: models.MTOServiceItem{},
		},
	)

	suite.T().Run("success SIT origin original zip lookup", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItemWithSITOriginZips.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		sitOriginZipOriginal, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		expected := mtoServiceItemWithSITOriginZips.SITOriginHHGOriginalAddress.PostalCode
		suite.Equal(expected, sitOriginZipOriginal)
	})

	suite.T().Run("fail to find SIT origin original zip lookup", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItemNoSITOriginZips.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		sitOriginZipOriginal, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.Equal("", sitOriginZipOriginal)
	})

}
