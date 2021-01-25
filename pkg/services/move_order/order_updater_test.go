package moveorder

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MoveOrderServiceSuite) TestMoveOrderUpdater() {
	expectedMoveTaskOrder := testdatagen.MakeDefaultMove(suite.DB())
	expectedMoveOrder := expectedMoveTaskOrder.Orders

	moveOrderUpdater := NewOrderUpdater(suite.DB())

	suite.T().Run("NotFoundError when order id doesn't exit", func(t *testing.T) {
		_, err := moveOrderUpdater.UpdateOrder("", models.Order{})
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("PreconditionsError when etag is stale", func(t *testing.T) {
		staleEtag := etag.GenerateEtag(expectedMoveOrder.UpdatedAt.Add(-1 * time.Minute))
		_, err := moveOrderUpdater.UpdateOrder(staleEtag, models.Order{ID: expectedMoveOrder.ID})
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
		actualOrder, err := moveOrderUpdater.UpdateOrder(expectedETag, updatedMoveOrder)

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
		actualOrder, err := moveOrderUpdater.UpdateOrder(expectedETag, updatedMoveOrder)

		suite.NoError(err)
		suite.Equal(swag.Int(20000), actualOrder.Entitlement.DBAuthorizedWeight)
		suite.Equal(swag.Bool(true), actualOrder.Entitlement.DependentsAuthorized)
	})

	suite.T().Run("Transaction rolled back after Order model validation error", func(t *testing.T) {
		defaultMoveOrder := testdatagen.MakeDefaultMove(suite.DB()).Orders
		serviceMember := defaultMoveOrder.ServiceMember

		// update service member to compare after a failed transaction
		updateAffiliation := models.AffiliationCOASTGUARD
		serviceMember.Affiliation = &updateAffiliation

		emptyStrSAC := ""
		updatedMoveOrder := models.Order{
			ID:                  defaultMoveOrder.ID,
			OriginDutyStationID: defaultMoveOrder.OriginDutyStationID,
			NewDutyStationID:    defaultMoveOrder.NewDutyStationID,
			IssueDate:           defaultMoveOrder.IssueDate,
			ReportByDate:        defaultMoveOrder.ReportByDate,
			OrdersType:          defaultMoveOrder.OrdersType,
			Entitlement: &models.Entitlement{ // try to update entitlement and see that it's not updated after failed transaction
				DBAuthorizedWeight:   swag.Int(20000),
				DependentsAuthorized: swag.Bool(false),
			},
			ServiceMember: serviceMember, // this is to make sure we're updating other models so we can check after a failed transaction
			SAC:           &emptyStrSAC,  // this will trigger validation error on Order model
		}

		expectedETag := etag.GenerateEtag(defaultMoveOrder.UpdatedAt)
		actualOrder, err := moveOrderUpdater.UpdateOrder(expectedETag, updatedMoveOrder)

		// check that we get back a validation error
		suite.EqualError(err, fmt.Sprintf("Invalid input for id: %s. SAC can not be blank.", defaultMoveOrder.ID))
		suite.Nil(actualOrder)

		// make sure that service member is not updated as well
		// we expect the affiliation to not have been updated, which is expected to be ARMY
		fetchedSM := models.ServiceMember{}
		_ = suite.DB().Find(&fetchedSM, serviceMember.ID)
		suite.EqualValues(models.AffiliationARMY, *fetchedSM.Affiliation)

		// check that entitlement is not updated as well
		fetchedEntitlement := models.Entitlement{}
		_ = suite.DB().Find(&fetchedEntitlement, defaultMoveOrder.Entitlement.ID)
		suite.NotEqual(20000, *fetchedEntitlement.DBAuthorizedWeight)
		suite.EqualValues(true, *fetchedEntitlement.DependentsAuthorized)
	})
}
