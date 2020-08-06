package serviceparamvaluelookups

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

func (suite *ServiceParamValueLookupsSuite) TestContractCodeLookup() {
	key := models.ServiceItemParamNameContractCode.String()

	// No database setup needed while this is stubbed out.
	paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), nil)

	suite.T().Run("golden path", func(t *testing.T) {
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal(ghcrateengine.DefaultContractCode, valueStr)
	})
}
