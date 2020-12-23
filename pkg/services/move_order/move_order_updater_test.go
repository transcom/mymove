package moveorder

import (
	"testing"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveOrderServiceSuite) TestMoveOrderUpdater() {
	expectedMoveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	expectedMoveOrder := expectedMoveTaskOrder.Orders

	queryBuilder := query.NewQueryBuilder(suite.DB())
	moveOrderUpdater := NewMoveOrderUpdater(suite.DB(), queryBuilder)

	suite.T().Run("NotFoundError when order id doesn't exit", func(t *testing.T) {
		_, err := moveOrderUpdater.UpdateMoveOrder(uuid.Nil, "", models.Order{})
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("PreconditionsError when etag is stale", func(t *testing.T) {
		staleEtag := etag.GenerateEtag(expectedMoveOrder.UpdatedAt.Add(-1 * time.Minute))
		_, err := moveOrderUpdater.UpdateMoveOrder(expectedMoveOrder.ID, staleEtag, models.Order{ID: expectedMoveOrder.ID})
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("Orders fields are updated without entitlement", func(t *testing.T) {
		defaultMoveOrder := testdatagen.MakeDefaultMove(suite.DB()).Orders

		newDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
		issueDate := time.Now().Add(-48 * time.Hour)
		reportByDate := time.Now().Add(72 * time.Hour)
		ordersTypeDetail := internalmessages.OrdersTypeDetailINSTRUCTION20WEEKS
		updatedMoveOrder := models.Order{
			ID:                  defaultMoveOrder.ID,
			OriginDutyStationID: &newDutyStation.ID,
			NewDutyStationID:    newDutyStation.ID,
			IssueDate:           issueDate,
			ReportByDate:        reportByDate,
			DepartmentIndicator: swag.String("COAST_GUARD"),
			OrdersType:          internalmessages.OrdersTypeSEPARATION,
			OrdersTypeDetail:    &ordersTypeDetail,
			Grade:               swag.String(string(models.ServiceMemberRankO10)),
			OrdersNumber:        swag.String("1122334455"),
			TAC:                 swag.String("8843"),
			SAC:                 swag.String("7766"),
		}

		expectedETag := etag.GenerateEtag(defaultMoveOrder.UpdatedAt)
		actualOrder, err := moveOrderUpdater.UpdateMoveOrder(defaultMoveOrder.ID, expectedETag, updatedMoveOrder)

		suite.NoError(err)
		suite.Equal(updatedMoveOrder.ID, actualOrder.ID)
		suite.Equal(updatedMoveOrder.NewDutyStationID, actualOrder.NewDutyStation.ID)
		suite.Equal(updatedMoveOrder.OriginDutyStationID.String(), actualOrder.OriginDutyStation.ID.String())
		suite.Equal(updatedMoveOrder.IssueDate, actualOrder.IssueDate)
		suite.Equal(updatedMoveOrder.ReportByDate, actualOrder.ReportByDate)
		suite.Equal(updatedMoveOrder.OrdersType, actualOrder.OrdersType)
		suite.Equal(updatedMoveOrder.OrdersTypeDetail, actualOrder.OrdersTypeDetail)
		suite.Equal(updatedMoveOrder.OrdersNumber, actualOrder.OrdersNumber)
		suite.Equal(updatedMoveOrder.DepartmentIndicator, actualOrder.DepartmentIndicator)
		suite.Equal(updatedMoveOrder.TAC, actualOrder.TAC)
		suite.Equal(updatedMoveOrder.SAC, actualOrder.SAC)
		suite.Equal(updatedMoveOrder.Grade, actualOrder.Grade)
	})

	suite.T().Run("Entitlement is updated with authorizedWeight or dependentsAuthorized", func(t *testing.T) {
		defaultMoveOrder := testdatagen.MakeDefaultMove(suite.DB()).Orders
		updatedMoveOrder := models.Order{
			ID:                  defaultMoveOrder.ID,
			OriginDutyStationID: defaultMoveOrder.OriginDutyStationID,
			NewDutyStationID:    defaultMoveOrder.NewDutyStationID,
			IssueDate:           defaultMoveOrder.IssueDate,
			ReportByDate:        defaultMoveOrder.ReportByDate,
			OrdersType:          defaultMoveOrder.OrdersType,
			Entitlement: &models.Entitlement{
				DBAuthorizedWeight:   swag.Int(20000),
				DependentsAuthorized: swag.Bool(true),
			},
		}

		expectedETag := etag.GenerateEtag(defaultMoveOrder.UpdatedAt)
		actualOrder, err := moveOrderUpdater.UpdateMoveOrder(defaultMoveOrder.ID, expectedETag, updatedMoveOrder)

		suite.NoError(err)
		suite.Equal(swag.Int(20000), actualOrder.Entitlement.DBAuthorizedWeight)
		suite.Equal(swag.Bool(true), actualOrder.Entitlement.DependentsAuthorized)
	})
}