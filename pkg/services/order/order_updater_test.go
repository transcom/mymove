package order

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

func (suite *OrderServiceSuite) TestOrderUpdater() {
	expectedMove := testdatagen.MakeDefaultMove(suite.DB())
	expectedOrder := expectedMove.Orders

	orderUpdater := NewOrderUpdater(suite.DB())

	suite.T().Run("NotFoundError when order id doesn't exit", func(t *testing.T) {
		_, err := orderUpdater.UpdateOrder("", models.Order{})
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("PreconditionsError when etag is stale", func(t *testing.T) {
		staleEtag := etag.GenerateEtag(expectedOrder.UpdatedAt.Add(-1 * time.Minute))
		_, err := orderUpdater.UpdateOrder(staleEtag, models.Order{ID: expectedOrder.ID})
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("Orders fields are updated without entitlement", func(t *testing.T) {
		defaultOrder := testdatagen.MakeDefaultMove(suite.DB()).Orders

		newDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
		issueDate := time.Now().Add(-48 * time.Hour)
		reportByDate := time.Now().Add(72 * time.Hour)
		ordersTypeDetail := internalmessages.OrdersTypeDetailINSTRUCTION20WEEKS
		updatedOrder := models.Order{
			ID:                  defaultOrder.ID,
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

		expectedETag := etag.GenerateEtag(defaultOrder.UpdatedAt)
		actualOrder, err := orderUpdater.UpdateOrder(expectedETag, updatedOrder)

		suite.NoError(err)
		suite.Equal(updatedOrder.ID, actualOrder.ID)
		suite.Equal(updatedOrder.NewDutyStationID, actualOrder.NewDutyStation.ID)
		suite.Equal(updatedOrder.OriginDutyStationID.String(), actualOrder.OriginDutyStation.ID.String())
		suite.Equal(updatedOrder.IssueDate, actualOrder.IssueDate)
		suite.Equal(updatedOrder.ReportByDate, actualOrder.ReportByDate)
		suite.Equal(updatedOrder.OrdersType, actualOrder.OrdersType)
		suite.Equal(updatedOrder.OrdersTypeDetail, actualOrder.OrdersTypeDetail)
		suite.Equal(updatedOrder.OrdersNumber, actualOrder.OrdersNumber)
		suite.Equal(updatedOrder.DepartmentIndicator, actualOrder.DepartmentIndicator)
		suite.Equal(updatedOrder.TAC, actualOrder.TAC)
		suite.Equal(updatedOrder.SAC, actualOrder.SAC)
		suite.Equal(updatedOrder.Grade, actualOrder.Grade)
	})

	suite.T().Run("Service member affiliation updated if order affiliation updated", func(t *testing.T) {
		defaultOrder := testdatagen.MakeDefaultMove(suite.DB()).Orders
		serviceMember := defaultOrder.ServiceMember

		newAffiliation := models.AffiliationNAVY

		suite.NotEqual(serviceMember.Affiliation, newAffiliation)

		var serviceMemberPatch models.ServiceMember

		serviceMemberPatch.Affiliation = &newAffiliation

		updatedOrder := models.Order{
			ID:                  defaultOrder.ID,
			OriginDutyStationID: defaultOrder.OriginDutyStationID,
			NewDutyStationID:    defaultOrder.NewDutyStationID,
			IssueDate:           defaultOrder.IssueDate,
			ReportByDate:        defaultOrder.ReportByDate,
			OrdersType:          defaultOrder.OrdersType,
			ServiceMember:       serviceMemberPatch,
		}

		expectedETag := etag.GenerateEtag(defaultOrder.UpdatedAt)
		_, err := orderUpdater.UpdateOrder(expectedETag, updatedOrder)

		suite.NoError(err)

		fetchedSM := models.ServiceMember{}
		_ = suite.DB().Find(&fetchedSM, serviceMember.ID)

		suite.EqualValues(newAffiliation, *fetchedSM.Affiliation)
	})

	suite.T().Run("Service member rank updated if order grade updated", func(t *testing.T) {
		defaultOrder := testdatagen.MakeDefaultMove(suite.DB()).Orders
		serviceMember := defaultOrder.ServiceMember

		newRank := models.ServiceMemberRankE2

		suite.NotEqual(serviceMember.Rank, newRank)

		updatedOrder := models.Order{
			ID:                  defaultOrder.ID,
			OriginDutyStationID: defaultOrder.OriginDutyStationID,
			NewDutyStationID:    defaultOrder.NewDutyStationID,
			IssueDate:           defaultOrder.IssueDate,
			ReportByDate:        defaultOrder.ReportByDate,
			OrdersType:          defaultOrder.OrdersType,
			Grade:               (*string)(&newRank),
		}

		expectedETag := etag.GenerateEtag(defaultOrder.UpdatedAt)
		actualOrder, err := orderUpdater.UpdateOrder(expectedETag, updatedOrder)

		suite.NoError(err)
		suite.Equal(newRank, models.ServiceMemberRank(*actualOrder.Grade))

		fetchedSM := models.ServiceMember{}
		_ = suite.DB().Find(&fetchedSM, serviceMember.ID)

		suite.EqualValues(newRank, *fetchedSM.Rank)
	})

	suite.T().Run("Service member current duty station updated if order origin duty station updated", func(t *testing.T) {
		defaultOrder := testdatagen.MakeDefaultMove(suite.DB()).Orders
		serviceMember := defaultOrder.ServiceMember

		newDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())

		suite.NotEqual(defaultOrder.OriginDutyStationID, newDutyStation.ID)

		updatedOrder := models.Order{
			ID:                  defaultOrder.ID,
			OriginDutyStationID: &newDutyStation.ID,
			NewDutyStationID:    defaultOrder.NewDutyStationID,
			IssueDate:           defaultOrder.IssueDate,
			ReportByDate:        defaultOrder.ReportByDate,
			OrdersType:          defaultOrder.OrdersType,
		}

		expectedETag := etag.GenerateEtag(defaultOrder.UpdatedAt)
		actualOrder, err := orderUpdater.UpdateOrder(expectedETag, updatedOrder)

		suite.NoError(err)
		suite.Equal(updatedOrder.ID, actualOrder.ID)
		suite.Equal(updatedOrder.OriginDutyStationID.String(), actualOrder.OriginDutyStation.ID.String())

		fetchedSM := models.ServiceMember{}
		_ = suite.DB().EagerPreload("DutyStation").Find(&fetchedSM, serviceMember.ID)

		suite.EqualValues(&newDutyStation.ID, fetchedSM.DutyStationID)
		suite.EqualValues(newDutyStation.ID, fetchedSM.DutyStation.ID)
		suite.EqualValues(newDutyStation.Name, fetchedSM.DutyStation.Name)
	})

	suite.T().Run("Entitlement is updated with authorizedWeight or dependentsAuthorized", func(t *testing.T) {
		defaultOrder := testdatagen.MakeDefaultMove(suite.DB()).Orders
		updatedOrder := models.Order{
			ID:                  defaultOrder.ID,
			OriginDutyStationID: defaultOrder.OriginDutyStationID,
			NewDutyStationID:    defaultOrder.NewDutyStationID,
			IssueDate:           defaultOrder.IssueDate,
			ReportByDate:        defaultOrder.ReportByDate,
			OrdersType:          defaultOrder.OrdersType,
			Entitlement: &models.Entitlement{
				DBAuthorizedWeight:   swag.Int(20000),
				DependentsAuthorized: swag.Bool(true),
			},
		}

		expectedETag := etag.GenerateEtag(defaultOrder.UpdatedAt)
		actualOrder, err := orderUpdater.UpdateOrder(expectedETag, updatedOrder)

		suite.NoError(err)
		suite.Equal(swag.Int(20000), actualOrder.Entitlement.DBAuthorizedWeight)
		suite.Equal(swag.Bool(true), actualOrder.Entitlement.DependentsAuthorized)
	})

	suite.T().Run("Transaction rolled back after Order model validation error", func(t *testing.T) {
		defaultOrder := testdatagen.MakeDefaultMove(suite.DB()).Orders
		serviceMember := defaultOrder.ServiceMember

		// update service member to compare after a failed transaction
		updateAffiliation := models.AffiliationCOASTGUARD
		serviceMember.Affiliation = &updateAffiliation

		emptyStrSAC := ""
		updatedOrder := models.Order{
			ID:                  defaultOrder.ID,
			OriginDutyStationID: defaultOrder.OriginDutyStationID,
			NewDutyStationID:    defaultOrder.NewDutyStationID,
			IssueDate:           defaultOrder.IssueDate,
			ReportByDate:        defaultOrder.ReportByDate,
			OrdersType:          defaultOrder.OrdersType,
			Entitlement: &models.Entitlement{ // try to update entitlement and see that it's not updated after failed transaction
				DBAuthorizedWeight:   swag.Int(20000),
				DependentsAuthorized: swag.Bool(false),
			},
			ServiceMember: serviceMember, // this is to make sure we're updating other models so we can check after a failed transaction
			SAC:           &emptyStrSAC,  // this will trigger validation error on Order model
		}

		expectedETag := etag.GenerateEtag(defaultOrder.UpdatedAt)
		actualOrder, err := orderUpdater.UpdateOrder(expectedETag, updatedOrder)

		// check that we get back a validation error
		suite.EqualError(err, fmt.Sprintf("Invalid input for id: %s. SAC can not be blank.", defaultOrder.ID))
		suite.Nil(actualOrder)

		// make sure that service member is not updated as well
		// we expect the affiliation to not have been updated, which is expected to be ARMY
		fetchedSM := models.ServiceMember{}
		_ = suite.DB().Find(&fetchedSM, serviceMember.ID)
		suite.EqualValues(models.AffiliationARMY, *fetchedSM.Affiliation)

		// check that entitlement is not updated as well
		fetchedEntitlement := models.Entitlement{}
		_ = suite.DB().Find(&fetchedEntitlement, defaultOrder.Entitlement.ID)
		suite.NotEqual(20000, *fetchedEntitlement.DBAuthorizedWeight)
		suite.EqualValues(true, *fetchedEntitlement.DependentsAuthorized)
	})
}
