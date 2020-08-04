package serviceparamvaluelookups

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

const defaultDistance = 1234

type ServiceParamValueLookupsSuite struct {
	testingsuite.PopTestSuite
	logger  Logger
	planner route.Planner
}

func (suite *ServiceParamValueLookupsSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestServiceParamValueLookupsSuite(t *testing.T) {
	planner := &mocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.Anything,
		mock.Anything,
	).Return(defaultDistance, nil)
	planner.On("Zip3TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(defaultDistance, nil)
	planner.On("Zip5TransitDistance",
		"90210",
		"94535",
	).Return(defaultDistance, nil)

	ts := &ServiceParamValueLookupsSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
		planner:      planner,
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *ServiceParamValueLookupsSuite) TestServiceParamValueLookup() {
	suite.T().Run("contract passed in", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{})
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()))

		suite.FatalNoError(err)
		suite.Equal(ghcrateengine.DefaultContractCode, paramLookup.ContractCode)
	})

	suite.T().Run("MTOServiceItem passed in", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{})

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()))

		suite.FatalNoError(err)
		suite.Equal(mtoServiceItem.ID, paramLookup.MTOServiceItemID)
		suite.NotNil(paramLookup.MTOServiceItem)
		suite.Equal(mtoServiceItem.MoveTaskOrderID, paramLookup.MTOServiceItem.MoveTaskOrderID)
	})

	suite.T().Run("nil MTOServiceItemID", func(t *testing.T) {
		badMTOServiceItemID := uuid.Must(uuid.NewV4())
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, badMTOServiceItemID, uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()))

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), fmt.Sprintf("id: %s not found looking for MTOServiceItemID", badMTOServiceItemID))
		var expected *ServiceItemParamKeyData = nil
		suite.Equal(expected, paramLookup)
	})
}
