package serviceparamvaluelookups

import (
	"strconv"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ServiceParamValueLookupsSuite) TestExternalCrateLookup() {
	suite.Run("ExternalCrate is true", func() {
		externalCrate := true
		mtoServiceItem := models.MTOServiceItem{
			ExternalCrate: &externalCrate,
		}

		paramLookup := ExternalCrateLookup{ServiceItem: mtoServiceItem}
		valueStr, err := paramLookup.lookup(suite.AppContextForTest(), nil)

		suite.FatalNoError(err)
		expected := strconv.FormatBool(externalCrate)
		suite.Equal(expected, valueStr)
	})

	suite.Run("ExternalCrate is false", func() {
		externalCrate := false
		mtoServiceItem := models.MTOServiceItem{
			ExternalCrate: &externalCrate,
		}

		paramLookup := ExternalCrateLookup{ServiceItem: mtoServiceItem}
		valueStr, err := paramLookup.lookup(suite.AppContextForTest(), nil)

		suite.FatalNoError(err)
		expected := strconv.FormatBool(externalCrate)
		suite.Equal(expected, valueStr)
	})

	suite.Run("ExternalCrate is nil", func() {
		mtoServiceItem := models.MTOServiceItem{
			ExternalCrate: nil,
		}

		paramLookup := ExternalCrateLookup{ServiceItem: mtoServiceItem}
		valueStr, err := paramLookup.lookup(suite.AppContextForTest(), nil)

		suite.FatalNoError(err)
		expected := "false"
		suite.Equal(expected, valueStr)
	})
}
