package order

import (
	"fmt"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/swagger/nullable"

	"github.com/go-openapi/strfmt"

	storageTest "github.com/transcom/mymove/pkg/storage/test"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/handlers"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *OrderServiceSuite) TestUpdateOrderAsTOO() {
	suite.T().Run("Returns an error when order is not found", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.UpdateOrderPayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), nonexistentUUID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when origin duty location is not found", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders
		newDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.UpdateOrderPayload{
			NewDutyLocationID:    handlers.FmtUUID(newDutyLocation.ID),
			OriginDutyLocationID: handlers.FmtUUID(nonexistentUUID),
		}
		eTag := etag.GenerateEtag(order.UpdatedAt)

		_, _, err := orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when new duty location is not found", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders
		originDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.UpdateOrderPayload{
			NewDutyLocationID:    handlers.FmtUUID(nonexistentUUID),
			OriginDutyLocationID: handlers.FmtUUID(originDutyLocation.ID),
		}
		eTag := etag.GenerateEtag(order.UpdatedAt)

		_, _, err := orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when the etag does not match", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders

		payload := ghcmessages.UpdateOrderPayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.T().Run("Updates the order when all fields are valid", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		move := testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{})
		order := move.Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		updatedOriginDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		ordersType := ghcmessages.OrdersTypeSEPARATION
		deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
		ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.UpdateOrderPayload{
			DepartmentIndicator:  &deptIndicator,
			IssueDate:            &dateIssued,
			NewDutyLocationID:    handlers.FmtUUID(updatedDestinationDutyLocation.ID),
			OriginDutyLocationID: handlers.FmtUUID(updatedOriginDutyLocation.ID),
			OrdersNumber:         handlers.FmtString("ORDER100"),
			OrdersType:           ghcmessages.NewOrdersType(ordersType),
			OrdersTypeDetail:     &ordersTypeDetail,
			ReportByDate:         &reportByDate,
			Tac:                  handlers.FmtString("E19A"),
			Sac:                  nullable.NewString("987654321"),
		}

		updatedOrder, _, err := orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)

		fetchedSM := models.ServiceMember{}
		_ = suite.DB().EagerPreload("DutyLocation").Find(&fetchedSM, order.ServiceMember.ID)

		suite.NoError(err)
		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(payload.NewDutyLocationID.String(), updatedOrder.NewDutyLocation.ID.String())
		suite.Equal(payload.OriginDutyLocationID.String(), updatedOrder.OriginDutyLocation.ID.String())
		suite.Equal(time.Time(*payload.IssueDate), updatedOrder.IssueDate)
		suite.Equal(time.Time(*payload.ReportByDate), updatedOrder.ReportByDate)
		suite.EqualValues(*payload.OrdersType, updatedOrder.OrdersType)
		suite.EqualValues(payload.OrdersTypeDetail, updatedOrder.OrdersTypeDetail)
		suite.Equal(payload.OrdersNumber, updatedOrder.OrdersNumber)
		suite.EqualValues(payload.DepartmentIndicator, updatedOrder.DepartmentIndicator)
		suite.Equal(payload.Tac, updatedOrder.TAC)
		suite.Equal(payload.Sac.Value, updatedOrder.SAC)
		suite.EqualValues(&updatedOriginDutyLocation.ID, fetchedSM.DutyLocationID)
		suite.EqualValues(updatedOriginDutyLocation.ID, fetchedSM.DutyLocation.ID)
		suite.EqualValues(updatedOriginDutyLocation.Name, fetchedSM.DutyLocation.Name)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Equal(move.Status, moveInDB.Status)
	})

	suite.T().Run("Rolls back transaction if Order is invalid", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{}).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		updatedOriginDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		ordersType := ghcmessages.OrdersTypeSEPARATION
		deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
		ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.UpdateOrderPayload{
			DepartmentIndicator:  &deptIndicator,
			IssueDate:            &dateIssued,
			NewDutyLocationID:    handlers.FmtUUID(updatedDestinationDutyLocation.ID),
			OriginDutyLocationID: handlers.FmtUUID(updatedOriginDutyLocation.ID),
			OrdersNumber:         handlers.FmtString("ORDER100"),
			OrdersType:           ghcmessages.NewOrdersType(ordersType),
			OrdersTypeDetail:     &ordersTypeDetail,
			ReportByDate:         &reportByDate,
			Tac:                  handlers.FmtString(""), // this will trigger a validation error on Order model
		}

		updatedOrder, _, err := orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)

		// check that we get back a validation error
		suite.EqualError(err, fmt.Sprintf("Invalid input for ID: %s. TransportationAccountingCode cannot be blank.", order.ID))
		suite.Nil(updatedOrder)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.T().Run("Allow Order update to have a missing HHG SAC", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{}).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		updatedOriginDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		ordersType := ghcmessages.OrdersTypeSEPARATION
		deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
		ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.UpdateOrderPayload{
			DepartmentIndicator:  &deptIndicator,
			IssueDate:            &dateIssued,
			NewDutyLocationID:    handlers.FmtUUID(updatedDestinationDutyLocation.ID),
			OriginDutyLocationID: handlers.FmtUUID(updatedOriginDutyLocation.ID),
			OrdersNumber:         handlers.FmtString("ORDER100"),
			OrdersType:           ghcmessages.NewOrdersType(ordersType),
			OrdersTypeDetail:     &ordersTypeDetail,
			ReportByDate:         &reportByDate,
			Tac:                  handlers.FmtString("E19A"),
			Sac:                  nullable.NewNullString(),
		}

		updatedOrder, _, err := orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)

		// check that we get back a validation error
		suite.NoError(err)
		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(payload.NewDutyLocationID.String(), updatedOrder.NewDutyLocation.ID.String())
		suite.Equal(payload.OriginDutyLocationID.String(), updatedOrder.OriginDutyLocation.ID.String())
		suite.Equal(time.Time(*payload.IssueDate), updatedOrder.IssueDate)
		suite.Equal(time.Time(*payload.ReportByDate), updatedOrder.ReportByDate)
		suite.EqualValues(*payload.OrdersType, updatedOrder.OrdersType)
		suite.EqualValues(payload.OrdersTypeDetail, updatedOrder.OrdersTypeDetail)
		suite.Equal(payload.OrdersNumber, updatedOrder.OrdersNumber)
		suite.EqualValues(payload.DepartmentIndicator, updatedOrder.DepartmentIndicator)
		suite.Equal(payload.Tac, updatedOrder.TAC)
		suite.Nil(updatedOrder.SAC)
	})

	suite.T().Run("Allow Order update to have a missing NTS SAC", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{}).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		updatedOriginDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		ordersType := ghcmessages.OrdersTypeSEPARATION
		deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
		ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.UpdateOrderPayload{
			DepartmentIndicator:  &deptIndicator,
			IssueDate:            &dateIssued,
			NewDutyLocationID:    handlers.FmtUUID(updatedDestinationDutyLocation.ID),
			OriginDutyLocationID: handlers.FmtUUID(updatedOriginDutyLocation.ID),
			OrdersNumber:         handlers.FmtString("ORDER100"),
			OrdersType:           ghcmessages.NewOrdersType(ordersType),
			OrdersTypeDetail:     &ordersTypeDetail,
			ReportByDate:         &reportByDate,
			Tac:                  handlers.FmtString("E19A"),
			NtsSac:               nullable.NewNullString(),
		}

		updatedOrder, _, err := orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)

		// check that we get back a validation error
		suite.NoError(err)
		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(payload.NewDutyLocationID.String(), updatedOrder.NewDutyLocation.ID.String())
		suite.Equal(payload.OriginDutyLocationID.String(), updatedOrder.OriginDutyLocation.ID.String())
		suite.Equal(time.Time(*payload.IssueDate), updatedOrder.IssueDate)
		suite.Equal(time.Time(*payload.ReportByDate), updatedOrder.ReportByDate)
		suite.EqualValues(*payload.OrdersType, updatedOrder.OrdersType)
		suite.EqualValues(payload.OrdersTypeDetail, updatedOrder.OrdersTypeDetail)
		suite.Equal(payload.OrdersNumber, updatedOrder.OrdersNumber)
		suite.EqualValues(payload.DepartmentIndicator, updatedOrder.DepartmentIndicator)
		suite.Equal(payload.Tac, updatedOrder.TAC)
		suite.Nil(updatedOrder.NtsSAC)
	})

	suite.T().Run("Allow Order update to have a missing NTS TAC", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{}).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		updatedOriginDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		ordersType := ghcmessages.OrdersTypeSEPARATION
		deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
		ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.UpdateOrderPayload{
			DepartmentIndicator:  &deptIndicator,
			IssueDate:            &dateIssued,
			NewDutyLocationID:    handlers.FmtUUID(updatedDestinationDutyLocation.ID),
			OriginDutyLocationID: handlers.FmtUUID(updatedOriginDutyLocation.ID),
			OrdersNumber:         handlers.FmtString("ORDER100"),
			OrdersType:           ghcmessages.NewOrdersType(ordersType),
			OrdersTypeDetail:     &ordersTypeDetail,
			ReportByDate:         &reportByDate,
			Tac:                  handlers.FmtString("E19A"),
			NtsTac:               nullable.NewNullString(),
		}

		updatedOrder, _, err := orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)

		// check that we get back a validation error
		suite.NoError(err)
		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(payload.NewDutyLocationID.String(), updatedOrder.NewDutyLocation.ID.String())
		suite.Equal(payload.OriginDutyLocationID.String(), updatedOrder.OriginDutyLocation.ID.String())
		suite.Equal(time.Time(*payload.IssueDate), updatedOrder.IssueDate)
		suite.Equal(time.Time(*payload.ReportByDate), updatedOrder.ReportByDate)
		suite.EqualValues(*payload.OrdersType, updatedOrder.OrdersType)
		suite.EqualValues(payload.OrdersTypeDetail, updatedOrder.OrdersTypeDetail)
		suite.Equal(payload.OrdersNumber, updatedOrder.OrdersNumber)
		suite.EqualValues(payload.DepartmentIndicator, updatedOrder.DepartmentIndicator)
		suite.Equal(payload.Tac, updatedOrder.TAC)
		suite.Nil(updatedOrder.NtsTAC)
	})

	suite.T().Run("Rolls back transaction if Order is invalid", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{}).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		updatedOriginDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		ordersType := ghcmessages.OrdersTypeSEPARATION
		deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
		ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.UpdateOrderPayload{
			DepartmentIndicator:  &deptIndicator,
			IssueDate:            &dateIssued,
			NewDutyLocationID:    handlers.FmtUUID(updatedDestinationDutyLocation.ID),
			OriginDutyLocationID: handlers.FmtUUID(updatedOriginDutyLocation.ID),
			OrdersNumber:         handlers.FmtString("ORDER100"),
			OrdersType:           ghcmessages.NewOrdersType(ordersType),
			OrdersTypeDetail:     &ordersTypeDetail,
			ReportByDate:         &reportByDate,
			Tac:                  handlers.FmtString(""), // this will trigger a validation error on Order model
			Sac:                  nullable.NewString(""),
		}

		updatedOrder, _, err := orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)

		// check that we get back a validation error
		suite.EqualError(err, fmt.Sprintf("Invalid input for ID: %s. TransportationAccountingCode cannot be blank.", order.ID))
		suite.Nil(updatedOrder)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.T().Run("Rolls back transaction if Order is missing required fields", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		orderWithoutDefaults := testdatagen.MakeOrderWithoutDefaults(suite.DB(), testdatagen.Assertions{})
		testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusServiceCounselingCompleted,
			},
			Order: orderWithoutDefaults,
		})

		eTag := etag.GenerateEtag(orderWithoutDefaults.UpdatedAt)

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		updatedOriginDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		ordersType := ghcmessages.OrdersTypeSEPARATION

		payload := ghcmessages.UpdateOrderPayload{
			IssueDate:            &dateIssued,
			NewDutyLocationID:    handlers.FmtUUID(updatedDestinationDutyLocation.ID),
			OriginDutyLocationID: handlers.FmtUUID(updatedOriginDutyLocation.ID),
			OrdersType:           ghcmessages.NewOrdersType(ordersType),
			ReportByDate:         &reportByDate,
		}

		suite.NoError(payload.Validate(strfmt.Default))

		updatedOrder, _, err := orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), orderWithoutDefaults.ID, payload, eTag)

		suite.Contains(err.Error(), fmt.Sprintf("Invalid input for ID: %s.", orderWithoutDefaults.ID))
		suite.Contains(err.Error(), "DepartmentIndicator cannot be blank.")
		suite.Contains(err.Error(), "OrdersTypeDetail cannot be blank.")
		suite.Contains(err.Error(), "TransportationAccountingCode cannot be blank.")
		suite.Contains(err.Error(), "OrdersNumber cannot be blank.")
		suite.Nil(updatedOrder)
		suite.IsType(apperror.InvalidInputError{}, err)
	})
}

func (suite *OrderServiceSuite) TestUpdateOrderAsCounselor() {
	suite.T().Run("Returns an error when order is not found", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.CounselingUpdateOrderPayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateOrderAsCounselor(suite.AppContextForTest(), nonexistentUUID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when the etag does not match", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders

		payload := ghcmessages.CounselingUpdateOrderPayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateOrderAsCounselor(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.T().Run("Updates the order when it is found", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeNeedsServiceCounselingMove(suite.DB()).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		updatedDestinationDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		updatedOriginDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
		ordersType := ghcmessages.OrdersTypeSEPARATION
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))

		body := ghcmessages.CounselingUpdateOrderPayload{
			IssueDate:            &dateIssued,
			NewDutyLocationID:    handlers.FmtUUID(updatedDestinationDutyLocation.ID),
			OriginDutyLocationID: handlers.FmtUUID(updatedOriginDutyLocation.ID),
			OrdersType:           ghcmessages.NewOrdersType(ordersType),
			ReportByDate:         &reportByDate,
		}

		eTag := etag.GenerateEtag(order.UpdatedAt)

		updatedOrder, _, err := orderUpdater.UpdateOrderAsCounselor(suite.AppContextForTest(), order.ID, body, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)

		suite.NoError(err)
		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(body.NewDutyLocationID.String(), updatedOrder.NewDutyLocation.ID.String())
		suite.Equal(body.OriginDutyLocationID.String(), updatedOrder.OriginDutyLocation.ID.String())
		suite.Equal(time.Time(*body.IssueDate), updatedOrder.IssueDate)
		suite.Equal(time.Time(*body.ReportByDate), updatedOrder.ReportByDate)
		suite.EqualValues(*body.OrdersType, updatedOrder.OrdersType)
	})
}

func (suite *OrderServiceSuite) TestUpdateAllowanceAsTOO() {
	suite.T().Run("Returns an error when order is not found", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.UpdateAllowancePayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateAllowanceAsTOO(suite.AppContextForTest(), nonexistentUUID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when the etag does not match", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders

		payload := ghcmessages.UpdateAllowancePayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateAllowanceAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.T().Run("Updates the allowance when all fields are valid", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{}).Orders

		newAuthorizedWeight := int64(10000)
		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := swag.Int64(100)
		proGearWeightSpouse := swag.Int64(10)
		rmeWeight := swag.Int64(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.UpdateAllowancePayload{
			Agency:               &affiliation,
			AuthorizedWeight:     &newAuthorizedWeight,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
			OrganizationalClothingAndIndividualEquipment: &ocie,
			ProGearWeight:                  proGearWeight,
			ProGearWeightSpouse:            proGearWeightSpouse,
			RequiredMedicalEquipmentWeight: rmeWeight,
		}

		updatedOrder, _, err := orderUpdater.UpdateAllowanceAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)

		suite.NoError(err)
		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(*payload.AuthorizedWeight, int64(*updatedOrder.Entitlement.DBAuthorizedWeight))
		suite.Equal(payload.DependentsAuthorized, updatedOrder.Entitlement.DependentsAuthorized)
		suite.Equal(*payload.ProGearWeight, int64(updatedOrder.Entitlement.ProGearWeight))
		suite.Equal(*payload.ProGearWeightSpouse, int64(updatedOrder.Entitlement.ProGearWeightSpouse))
		suite.Equal(*payload.RequiredMedicalEquipmentWeight, int64(updatedOrder.Entitlement.RequiredMedicalEquipmentWeight))
		suite.Equal(*payload.OrganizationalClothingAndIndividualEquipment, updatedOrder.Entitlement.OrganizationalClothingAndIndividualEquipment)
	})
}

func (suite *OrderServiceSuite) TestUpdateAllowanceAsCounselor() {
	suite.T().Run("Returns an error when order is not found", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.CounselingUpdateAllowancePayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateAllowanceAsCounselor(suite.AppContextForTest(), nonexistentUUID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when the etag does not match", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders

		payload := ghcmessages.CounselingUpdateAllowancePayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateAllowanceAsCounselor(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.T().Run("Updates the allowance when all fields are valid", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeNeedsServiceCounselingMove(suite.DB()).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := swag.Int64(100)
		proGearWeightSpouse := swag.Int64(10)
		rmeWeight := swag.Int64(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               &affiliation,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
			OrganizationalClothingAndIndividualEquipment: &ocie,
			ProGearWeight:                  proGearWeight,
			ProGearWeightSpouse:            proGearWeightSpouse,
			RequiredMedicalEquipmentWeight: rmeWeight,
		}

		updatedOrder, _, err := orderUpdater.UpdateAllowanceAsCounselor(suite.AppContextForTest(), order.ID, payload, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)

		fetchedSM := models.ServiceMember{}
		_ = suite.DB().Find(&fetchedSM, order.ServiceMember.ID)

		suite.NoError(err)
		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(payload.DependentsAuthorized, updatedOrder.Entitlement.DependentsAuthorized)
		suite.Equal(*payload.ProGearWeight, int64(updatedOrder.Entitlement.ProGearWeight))
		suite.Equal(*payload.ProGearWeightSpouse, int64(updatedOrder.Entitlement.ProGearWeightSpouse))
		suite.Equal(*payload.RequiredMedicalEquipmentWeight, int64(updatedOrder.Entitlement.RequiredMedicalEquipmentWeight))
		suite.Equal(*payload.OrganizationalClothingAndIndividualEquipment, updatedOrder.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.EqualValues(*payload.Grade, *fetchedSM.Rank)
		suite.EqualValues(payload.Agency, fetchedSM.Affiliation)
	})

	suite.T().Run("Updates the allowance when move needs service counseling and order fields are missing", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		orderWithoutDefaults := testdatagen.MakeOrderWithoutDefaults(suite.DB(), testdatagen.Assertions{})
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusNeedsServiceCounseling,
			},
			Order: orderWithoutDefaults,
		})

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := swag.Int64(100)
		proGearWeightSpouse := swag.Int64(10)
		rmeWeight := swag.Int64(10000)
		eTag := etag.GenerateEtag(orderWithoutDefaults.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               &affiliation,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
			OrganizationalClothingAndIndividualEquipment: &ocie,
			ProGearWeight:                  proGearWeight,
			ProGearWeightSpouse:            proGearWeightSpouse,
			RequiredMedicalEquipmentWeight: rmeWeight,
		}

		updatedOrder, _, err := orderUpdater.UpdateAllowanceAsCounselor(suite.AppContextForTest(), orderWithoutDefaults.ID, payload, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, orderWithoutDefaults.ID)

		fetchedSM := models.ServiceMember{}
		_ = suite.DB().Find(&fetchedSM, orderWithoutDefaults.ServiceMember.ID)

		suite.NoError(err)
		suite.Equal(orderWithoutDefaults.ID.String(), updatedOrder.ID.String())
		suite.Equal(payload.DependentsAuthorized, updatedOrder.Entitlement.DependentsAuthorized)
		suite.Equal(*payload.ProGearWeight, int64(updatedOrder.Entitlement.ProGearWeight))
		suite.Equal(*payload.ProGearWeightSpouse, int64(updatedOrder.Entitlement.ProGearWeightSpouse))
		suite.Equal(*payload.RequiredMedicalEquipmentWeight, int64(updatedOrder.Entitlement.RequiredMedicalEquipmentWeight))
		suite.Equal(*payload.OrganizationalClothingAndIndividualEquipment, updatedOrder.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.EqualValues(*payload.Grade, *fetchedSM.Rank)
		suite.EqualValues(payload.Agency, fetchedSM.Affiliation)

		// make sure that there are missing submission fields and move is in correct status
		fetchedMove := models.Move{}
		_ = suite.DB().Find(&fetchedMove, move.ID)
		suite.Equal(models.MoveStatusNeedsServiceCounseling, fetchedMove.Status)
		suite.Nil(updatedOrder.TAC)
		suite.Nil(updatedOrder.SAC)
		suite.Nil(updatedOrder.DepartmentIndicator)
		suite.Nil(updatedOrder.OrdersTypeDetail)
	})

	suite.T().Run("Entire update is aborted when ProGearWeight is over max amount", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeNeedsServiceCounselingMove(suite.DB()).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := swag.Int64(2001)
		proGearWeightSpouse := swag.Int64(10)
		rmeWeight := swag.Int64(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               &affiliation,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
			OrganizationalClothingAndIndividualEquipment: &ocie,
			ProGearWeight:                  proGearWeight,
			ProGearWeightSpouse:            proGearWeightSpouse,
			RequiredMedicalEquipmentWeight: rmeWeight,
		}

		updatedOrder, _, err := orderUpdater.UpdateAllowanceAsCounselor(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		var orderInDB models.Order
		err = suite.DB().EagerPreload("Entitlement").Find(&orderInDB, order.ID)

		suite.NoError(err)
		suite.NotEqual(payload.ProGearWeight, orderInDB.Entitlement.ProGearWeight)
		suite.Nil(updatedOrder)
	})

	suite.T().Run("Entire update is aborted when ProGearWeightSpouse is over max amount", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		order := testdatagen.MakeNeedsServiceCounselingMove(suite.DB()).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := swag.Int64(100)
		proGearWeightSpouse := swag.Int64(501)
		rmeWeight := swag.Int64(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               &affiliation,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
			OrganizationalClothingAndIndividualEquipment: &ocie,
			ProGearWeight:                  proGearWeight,
			ProGearWeightSpouse:            proGearWeightSpouse,
			RequiredMedicalEquipmentWeight: rmeWeight,
		}

		updatedOrder, _, err := orderUpdater.UpdateAllowanceAsCounselor(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)

		var orderInDB models.Order
		err = suite.DB().EagerPreload("Entitlement").Find(&orderInDB, order.ID)

		suite.NoError(err)
		suite.NotEqual(payload.ProGearWeightSpouse, orderInDB.Entitlement.ProGearWeightSpouse)
		suite.Nil(updatedOrder)
	})
}

func (suite *OrderServiceSuite) TestUploadAmendedOrdersForCustomer() {

	suite.T().Run("Creates and saves new amendedOrder doc when the order.UploadedAmendedOrders is nil", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		dutyLocation := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
			DutyLocation: models.DutyLocation{
				Address: testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{}),
			},
		})
		var moves models.Moves
		mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})

		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				OriginDutyLocation: &dutyLocation,
			},
			Move: mto,
		})
		order.Moves = append(moves, mto)

		file := testdatagen.FixtureRuntimeFile("test.pdf")
		defer func() {
			fileCloseErr := file.Close()
			suite.NoError(fileCloseErr)
		}()

		fakeS3 := storageTest.NewFakeS3Storage(true)

		suite.NotEqual(uuid.Nil, order.ServiceMemberID, "ServiceMember has ID that is not 0/empty")
		suite.NotEqual(uuid.Nil, order.ServiceMember.UserID, "ServiceMember.UserID has ID that is not 0/empty")

		upload, url, verrs, err := orderUpdater.UploadAmendedOrdersAsCustomer(
			suite.AppContextForTest(),
			order.ServiceMember.UserID,
			order.ID,
			file.Data,
			file.Header.Filename,
			fakeS3)
		suite.NoError(err)
		suite.NoVerrs(verrs)

		expectedChecksum := "nOE6HwzyE4VEDXn67ULeeA=="
		if upload.Checksum != expectedChecksum {
			t.Errorf("Did not calculate the correct MD5: expected %s, got %s", expectedChecksum, upload.Checksum)
		}

		var orderInDB models.Order
		err = suite.DB().
			EagerPreload("UploadedAmendedOrders").
			Find(&orderInDB, order.ID)

		suite.NoError(err)
		suite.Equal(orderInDB.ID.String(), order.ID.String())
		suite.NotNil(orderInDB.UploadedAmendedOrders)

		findUpload := models.Upload{}
		err = suite.DB().Find(&findUpload, upload.ID)
		if err != nil {
			t.Fatalf("Couldn't find expected upload.")
		}
		suite.Equal(upload.ID.String(), findUpload.ID.String(), "found upload in db")
		suite.NotEmpty(url, "URL is populated")
	})

	suite.T().Run("Returns an error when order is not found", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())

		file := testdatagen.FixtureRuntimeFile("test.pdf")
		defer func() {
			fileCloseErr := file.Close()
			suite.NoError(fileCloseErr)
		}()

		fakeS3 := storageTest.NewFakeS3Storage(true)

		_, _, verrs, err := orderUpdater.UploadAmendedOrdersAsCustomer(
			suite.AppContextForTest(),
			nonexistentUUID,
			nonexistentUUID,
			file.Data,
			file.Header.Filename,
			fakeS3)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "not found while looking for order")
		suite.NoVerrs(verrs)
	})

	suite.T().Run("Saves userUpload payload to order.UploadedAmendedOrders if the document already exists", func(t *testing.T) {
		moveRouter := move.NewMoveRouter()
		orderUpdater := NewOrderUpdater(moveRouter)
		dutyLocation := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
			DutyLocation: models.DutyLocation{
				Address: testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{}),
			},
		})
		var moves models.Moves
		mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})

		document := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{})
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				OriginDutyLocation:      &dutyLocation,
				UploadedAmendedOrders:   &document,
				UploadedAmendedOrdersID: &document.ID,
			},
			Move: mto,
		})
		order.Moves = append(moves, mto)

		file := testdatagen.FixtureRuntimeFile("test.pdf")
		defer func() {
			fileCloseErr := file.Close()
			suite.NoError(fileCloseErr)
		}()

		fakeS3 := storageTest.NewFakeS3Storage(true)

		suite.NotEqual(uuid.Nil, order.ServiceMemberID, "ServiceMember has ID that is not 0/empty")
		suite.NotEqual(uuid.Nil, order.ServiceMember.UserID, "ServiceMember.UserID has ID that is not 0/empty")

		_, _, verrs, err := orderUpdater.UploadAmendedOrdersAsCustomer(
			suite.AppContextForTest(),
			order.ServiceMember.UserID,
			order.ID,
			file.Data,
			file.Header.Filename,
			fakeS3)
		suite.NoError(err)
		suite.NoVerrs(verrs)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().
			EagerPreload("UploadedAmendedOrders").
			Find(&orderInDB, order.ID)

		suite.NoError(err)
		suite.NotNil(orderInDB.ID)
		suite.NotNil(orderInDB.UploadedAmendedOrders)
		suite.Equal(document.ID, *orderInDB.UploadedAmendedOrdersID)
		suite.NotNil(order.UploadedAmendedOrders)
		suite.NotNil(orderInDB.UploadedAmendedOrders)
	})
}
