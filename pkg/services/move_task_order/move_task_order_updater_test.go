package movetaskorder

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/etag"
	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderUpdater_UpdatePostCounselingInfo() {
	expectedOrder := testdatagen.MakeDefaultOrder(suite.DB())
	expectedMTO := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: expectedOrder,
	})

	queryBuilder := query.NewQueryBuilder(suite.DB())
	mtoUpdater := NewMoveTaskOrderUpdater(suite.DB(), queryBuilder, mtoserviceitem.NewMTOServiceItemCreator(queryBuilder))
	body := movetaskorderops.UpdateMTOPostCounselingInformationBody{
		PpmType:            "FULL",
		PpmEstimatedWeight: 3000,
		PointOfContact:     "user@prime.com",
	}

	suite.T().Run("MTO post counseling information is updated succesfully", func(t *testing.T) {
		eTag := base64.StdEncoding.EncodeToString([]byte(expectedMTO.UpdatedAt.Format(time.RFC3339Nano)))

		actualMTO, err := mtoUpdater.UpdatePostCounselingInfo(expectedMTO.ID, body, eTag)

		suite.NoError(err)

		suite.NotZero(expectedMTO.ID, actualMTO.ID)
		suite.Equal(expectedMTO.Orders.ID, actualMTO.Orders.ID)
		suite.NotZero(actualMTO.Orders)
		suite.NotNil(expectedMTO.ReferenceID)
		suite.NotNil(expectedMTO.Locator)
		suite.Nil(expectedMTO.AvailableToPrimeAt)
		suite.NotEqual(expectedMTO.Status, models.MoveStatusCANCELED)

		suite.NotNil(expectedMTO.Orders.ServiceMember.FirstName)
		suite.NotNil(expectedMTO.Orders.ServiceMember.LastName)
		suite.NotNil(expectedMTO.Orders.NewDutyStation.Address.City)
		suite.NotNil(expectedMTO.Orders.NewDutyStation.Address.State)
	})

	suite.T().Run("Etag is stale", func(t *testing.T) {
		eTag := etag.GenerateEtag(time.Now())
		_, err := mtoUpdater.UpdatePostCounselingInfo(expectedMTO.ID, body, eTag)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})
}
