package serviceparamvaluelookups

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestContractCodeLookup() {
	key := models.ServiceItemParamNameContractCode.String()

	suite.T().Run("golden path", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{})
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()))
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal(ghcrateengine.DefaultContractCode, valueStr)
	})
}
