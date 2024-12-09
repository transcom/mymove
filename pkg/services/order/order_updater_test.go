package order

import (
	"fmt"
	"reflect"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/move"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/swagger/nullable"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *OrderServiceSuite) TestUpdateOrderAsTOO() {
	suite.Run("Returns an error when order is not found", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.UpdateOrderPayload{}
		eTag := ""

		_, _, err = orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), nonexistentUUID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when origin duty location is not found", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildMove(suite.DB(), nil, nil).Orders
		newDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.UpdateOrderPayload{
			NewDutyLocationID:    handlers.FmtUUID(newDutyLocation.ID),
			OriginDutyLocationID: handlers.FmtUUID(nonexistentUUID),
		}
		eTag := etag.GenerateEtag(order.UpdatedAt)

		_, _, err = orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when new duty location is not found", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildMove(suite.DB(), nil, nil).Orders
		originDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.UpdateOrderPayload{
			NewDutyLocationID:    handlers.FmtUUID(nonexistentUUID),
			OriginDutyLocationID: handlers.FmtUUID(originDutyLocation.ID),
		}
		eTag := etag.GenerateEtag(order.UpdatedAt)

		_, _, err = orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when the etag does not match", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildMove(suite.DB(), nil, nil).Orders

		payload := ghcmessages.UpdateOrderPayload{}
		eTag := ""

		_, _, err = orderUpdater.UpdateOrderAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Updates the order when all fields are valid", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		move := factory.BuildServiceCounselingCompletedMove(suite.DB(), nil, nil)
		order := move.Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		updatedOriginDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "77777",
				},
			},
		}, nil)
		updatedGbloc := factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), updatedOriginDutyLocation.Address.PostalCode, "UUUU")
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
		suite.EqualValues(updatedGbloc.GBLOC, *updatedOrder.OriginDutyLocationGBLOC)

		var moveInDB models.Move
		err = suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Equal(move.Status, moveInDB.Status)
	})

	suite.Run("Rolls back transaction if Order is invalid", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildServiceCounselingCompletedMove(suite.DB(), nil, nil).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		updatedOriginDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
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

	suite.Run("Allow Order update to have a missing HHG SAC", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildServiceCounselingCompletedMove(suite.DB(), nil, nil).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		updatedOriginDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
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

	suite.Run("Allow Order update to have a missing NTS SAC", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildServiceCounselingCompletedMove(suite.DB(), nil, nil).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		updatedOriginDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
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

	suite.Run("Allow Order update to have a missing NTS TAC", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildServiceCounselingCompletedMove(suite.DB(), nil, nil).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		updatedOriginDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
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

	suite.Run("Rolls back transaction if Order is invalid", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildServiceCounselingCompletedMove(suite.DB(), nil, nil).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		updatedOriginDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
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

	suite.Run("Rolls back transaction if Order is missing required fields", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		orderWithoutDefaults := factory.BuildOrderWithoutDefaults(suite.DB(), nil, nil)

		factory.BuildServiceCounselingCompletedMove(suite.DB(), []factory.Customization{
			{
				Model:    orderWithoutDefaults,
				LinkOnly: true,
			},
		}, nil)

		eTag := etag.GenerateEtag(orderWithoutDefaults.UpdatedAt)

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		updatedOriginDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
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
	suite.Run("Returns an error when order is not found", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.CounselingUpdateOrderPayload{}
		eTag := ""

		_, _, err = orderUpdater.UpdateOrderAsCounselor(suite.AppContextForTest(), nonexistentUUID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when the etag does not match", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildMove(suite.DB(), nil, nil).Orders

		payload := ghcmessages.CounselingUpdateOrderPayload{}
		eTag := ""

		_, _, err = orderUpdater.UpdateOrderAsCounselor(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Updates the order when it is found", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildNeedsServiceCounselingMove(suite.DB(), nil, nil).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		updatedDestinationDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		updatedOriginDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		ordersType := ghcmessages.OrdersTypeSEPARATION
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
		ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
		grade := ghcmessages.GradeO5

		body := ghcmessages.CounselingUpdateOrderPayload{
			IssueDate:            &dateIssued,
			NewDutyLocationID:    handlers.FmtUUID(updatedDestinationDutyLocation.ID),
			OriginDutyLocationID: handlers.FmtUUID(updatedOriginDutyLocation.ID),
			OrdersType:           ghcmessages.NewOrdersType(ordersType),
			ReportByDate:         &reportByDate,
			DepartmentIndicator:  &deptIndicator,
			OrdersNumber:         handlers.FmtString("ORDER100"),
			OrdersTypeDetail:     &ordersTypeDetail,
			Tac:                  handlers.FmtString("E19A"),
			Sac:                  nullable.NewString("987654321"),
			Grade:                &grade,
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
		suite.EqualValues(*body.OrdersTypeDetail, *updatedOrder.OrdersTypeDetail)
		suite.EqualValues(body.OrdersNumber, updatedOrder.OrdersNumber)
		suite.EqualValues(body.DepartmentIndicator, updatedOrder.DepartmentIndicator)
		suite.EqualValues(body.Tac, updatedOrder.TAC)
		suite.EqualValues(body.Sac.Value, updatedOrder.SAC)
		suite.Equal(*updatedOrder.Entitlement.DBAuthorizedWeight, 16000)
	})

	suite.Run("Updates the PPM actual expense reimbursement when pay grade is civilian", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)

		ppmShipment := factory.BuildPPMShipmentThatNeedsCloseout(suite.DB(), nil, nil)
		move := ppmShipment.Shipment.MoveTaskOrder

		order := move.Orders
		grade := ghcmessages.GradeCIVILIANEMPLOYEE
		body := ghcmessages.CounselingUpdateOrderPayload{
			Grade: &grade,
		}
		eTag := etag.GenerateEtag(order.UpdatedAt)

		var moved models.Move
		err = suite.DB().Find(&moved, move.ID)
		suite.NoError(err)

		_, _, errs := orderUpdater.UpdateOrderAsCounselor(suite.AppContextForTest(), order.ID, body, eTag)
		suite.NoError(errs)

		var updatedPPMShipment models.PPMShipment
		err = suite.DB().Find(&updatedPPMShipment, ppmShipment.ID)

		suite.NoError(err)
		suite.EqualValues(true, *updatedPPMShipment.IsActualExpenseReimbursement)
	})

	suite.Run("Rolls back transaction if Order is invalid", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildOrderWithoutDefaults(suite.DB(), nil, nil)

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		updatedOriginDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		ordersType := ghcmessages.OrdersTypePERMANENTCHANGEOFSTATION
		deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
		ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.CounselingUpdateOrderPayload{
			DepartmentIndicator:  &deptIndicator,
			IssueDate:            &dateIssued,
			NewDutyLocationID:    handlers.FmtUUID(updatedDestinationDutyLocation.ID),
			OriginDutyLocationID: handlers.FmtUUID(updatedOriginDutyLocation.ID),
			OrdersNumber:         handlers.FmtString("ORDER100"),
			OrdersType:           ghcmessages.NewOrdersType(ordersType),
			OrdersTypeDetail:     &ordersTypeDetail,
			ReportByDate:         &reportByDate,
		}

		updatedOrder, _, err := orderUpdater.UpdateOrderAsCounselor(suite.AppContextForTest(), order.ID, payload, eTag)

		// check that we get back a validation error
		suite.EqualError(err, fmt.Sprintf("Invalid input for ID: %s. TransportationAccountingCode cannot be blank.", order.ID))
		suite.Nil(updatedOrder)
		suite.IsType(apperror.InvalidInputError{}, err)
	})
}

func (suite *OrderServiceSuite) TestUpdateAllowanceAsTOO() {
	suite.Run("Returns an error when order is not found", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.UpdateAllowancePayload{}
		eTag := ""

		_, _, err = orderUpdater.UpdateAllowanceAsTOO(suite.AppContextForTest(), nonexistentUUID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when the etag does not match", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildMove(suite.DB(), nil, nil).Orders

		payload := ghcmessages.UpdateAllowancePayload{}
		eTag := ""

		_, _, err = orderUpdater.UpdateAllowanceAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Updates the allowance when all fields are valid and no dependents", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildServiceCounselingCompletedMove(suite.DB(), nil, nil).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := models.Int64Pointer(100)
		proGearWeightSpouse := models.Int64Pointer(10)
		rmeWeight := models.Int64Pointer(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.UpdateAllowancePayload{
			Agency:               &affiliation,
			DependentsAuthorized: models.BoolPointer(true),
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
		suite.Equal(payload.DependentsAuthorized, updatedOrder.Entitlement.DependentsAuthorized)
		suite.Equal(*payload.ProGearWeight, int64(updatedOrder.Entitlement.ProGearWeight))
		suite.Equal(*payload.ProGearWeightSpouse, int64(updatedOrder.Entitlement.ProGearWeightSpouse))
		suite.Equal(*payload.RequiredMedicalEquipmentWeight, int64(updatedOrder.Entitlement.RequiredMedicalEquipmentWeight))
		suite.Equal(*payload.OrganizationalClothingAndIndividualEquipment, updatedOrder.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(*updatedOrder.Entitlement.DBAuthorizedWeight, 16000)
	})

	suite.Run("Updates the allowance when all OCONUS fields are valid with dependents", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildServiceCounselingCompletedMove(suite.DB(), nil, nil).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := models.Int64Pointer(100)
		proGearWeightSpouse := models.Int64Pointer(10)
		rmeWeight := models.Int64Pointer(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.UpdateAllowancePayload{
			Agency:               &affiliation,
			DependentsAuthorized: models.BoolPointer(true),
			Grade:                &grade,
			OrganizationalClothingAndIndividualEquipment: &ocie,
			ProGearWeight:                  proGearWeight,
			ProGearWeightSpouse:            proGearWeightSpouse,
			RequiredMedicalEquipmentWeight: rmeWeight,
			AccompaniedTour:                models.BoolPointer(true),
			DependentsTwelveAndOver:        models.Int64Pointer(2),
			DependentsUnderTwelve:          models.Int64Pointer(4),
		}

		updatedOrder, _, err := orderUpdater.UpdateAllowanceAsTOO(suite.AppContextForTest(), order.ID, payload, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)

		suite.NoError(err)
		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(payload.DependentsAuthorized, updatedOrder.Entitlement.DependentsAuthorized)
		suite.Equal(*payload.ProGearWeight, int64(updatedOrder.Entitlement.ProGearWeight))
		suite.Equal(*payload.ProGearWeightSpouse, int64(updatedOrder.Entitlement.ProGearWeightSpouse))
		suite.Equal(*payload.RequiredMedicalEquipmentWeight, int64(updatedOrder.Entitlement.RequiredMedicalEquipmentWeight))
		suite.Equal(*payload.OrganizationalClothingAndIndividualEquipment, updatedOrder.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(*updatedOrder.Entitlement.DBAuthorizedWeight, 16000)
		suite.Equal(*payload.DependentsTwelveAndOver, int64(*updatedOrder.Entitlement.DependentsTwelveAndOver))
		suite.Equal(*payload.AccompaniedTour, *updatedOrder.Entitlement.AccompaniedTour)
		suite.Equal(*payload.DependentsUnderTwelve, int64(*updatedOrder.Entitlement.DependentsUnderTwelve))
	})

	suite.Run("Updates the allowance when all fields are valid with dependents", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		// Build with dependents trait
		order := factory.BuildServiceCounselingCompletedMove(suite.DB(), nil, []factory.Trait{
			factory.GetTraitHasDependents,
		}).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := models.Int64Pointer(100)
		proGearWeightSpouse := models.Int64Pointer(10)
		rmeWeight := models.Int64Pointer(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.UpdateAllowancePayload{
			Agency:               &affiliation,
			DependentsAuthorized: models.BoolPointer(true),
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
		suite.Equal(payload.DependentsAuthorized, updatedOrder.Entitlement.DependentsAuthorized)
		suite.Equal(*payload.ProGearWeight, int64(updatedOrder.Entitlement.ProGearWeight))
		suite.Equal(*payload.ProGearWeightSpouse, int64(updatedOrder.Entitlement.ProGearWeightSpouse))
		suite.Equal(*payload.RequiredMedicalEquipmentWeight, int64(updatedOrder.Entitlement.RequiredMedicalEquipmentWeight))
		suite.Equal(*payload.OrganizationalClothingAndIndividualEquipment, updatedOrder.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(*updatedOrder.Entitlement.DBAuthorizedWeight, 17500)
	})
}

func (suite *OrderServiceSuite) TestUpdateAllowanceAsCounselor() {
	suite.Run("Returns an error when order is not found", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.CounselingUpdateAllowancePayload{}
		eTag := ""

		_, _, err = orderUpdater.UpdateAllowanceAsCounselor(suite.AppContextForTest(), nonexistentUUID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when the etag does not match", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildMove(suite.DB(), nil, nil).Orders

		payload := ghcmessages.CounselingUpdateAllowancePayload{}
		eTag := ""

		_, _, err = orderUpdater.UpdateAllowanceAsCounselor(suite.AppContextForTest(), order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Updates the entitlement of OCONUS fields", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildNeedsServiceCounselingMove(suite.DB(), nil, nil).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := models.Int64Pointer(100)
		proGearWeightSpouse := models.Int64Pointer(10)
		rmeWeight := models.Int64Pointer(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               &affiliation,
			DependentsAuthorized: models.BoolPointer(true),
			Grade:                &grade,
			OrganizationalClothingAndIndividualEquipment: &ocie,
			ProGearWeight:                  proGearWeight,
			ProGearWeightSpouse:            proGearWeightSpouse,
			RequiredMedicalEquipmentWeight: rmeWeight,
			AccompaniedTour:                models.BoolPointer(true),
			DependentsTwelveAndOver:        models.Int64Pointer(1),
			DependentsUnderTwelve:          models.Int64Pointer(2),
		}

		updatedOrder, _, err := orderUpdater.UpdateAllowanceAsCounselor(suite.AppContextForTest(), order.ID, payload, eTag)
		suite.NoError(err)

		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(payload.DependentsAuthorized, updatedOrder.Entitlement.DependentsAuthorized)
		suite.Equal(*payload.ProGearWeight, int64(updatedOrder.Entitlement.ProGearWeight))
		suite.Equal(*payload.ProGearWeightSpouse, int64(updatedOrder.Entitlement.ProGearWeightSpouse))
		suite.Equal(*payload.RequiredMedicalEquipmentWeight, int64(updatedOrder.Entitlement.RequiredMedicalEquipmentWeight))
		suite.Equal(*payload.OrganizationalClothingAndIndividualEquipment, updatedOrder.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(*payload.AccompaniedTour, *updatedOrder.Entitlement.AccompaniedTour)
		suite.Equal(*payload.DependentsUnderTwelve, int64(*updatedOrder.Entitlement.DependentsUnderTwelve))
		suite.Equal(*payload.DependentsTwelveAndOver, int64(*updatedOrder.Entitlement.DependentsTwelveAndOver))
	})

	suite.Run("Updates the allowance when all fields are valid with dependents authorized but not present", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildNeedsServiceCounselingMove(suite.DB(), nil, nil).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := models.Int64Pointer(100)
		proGearWeightSpouse := models.Int64Pointer(10)
		rmeWeight := models.Int64Pointer(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               &affiliation,
			DependentsAuthorized: models.BoolPointer(true),
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
		suite.EqualValues(payload.Agency, fetchedSM.Affiliation)
		suite.Equal(*updatedOrder.Entitlement.DBAuthorizedWeight, 16000)
	})

	suite.Run("Updates the allowance when all fields are valid with dependents present and authorized", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildNeedsServiceCounselingMove(suite.DB(), nil, []factory.Trait{
			factory.GetTraitHasDependents,
		}).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := models.Int64Pointer(100)
		proGearWeightSpouse := models.Int64Pointer(10)
		rmeWeight := models.Int64Pointer(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               &affiliation,
			DependentsAuthorized: models.BoolPointer(true),
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
		suite.EqualValues(payload.Agency, fetchedSM.Affiliation)
		suite.Equal(*updatedOrder.Entitlement.DBAuthorizedWeight, 17500)
	})

	suite.Run("Updates the allowance when move needs service counseling and order fields are missing", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		orderWithoutDefaults := factory.BuildOrderWithoutDefaults(suite.DB(), nil, nil)
		move := factory.BuildNeedsServiceCounselingMove(suite.DB(), []factory.Customization{
			{
				Model:    orderWithoutDefaults,
				LinkOnly: true,
			},
		}, nil)

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := models.Int64Pointer(100)
		proGearWeightSpouse := models.Int64Pointer(10)
		rmeWeight := models.Int64Pointer(10000)
		eTag := etag.GenerateEtag(orderWithoutDefaults.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               &affiliation,
			DependentsAuthorized: models.BoolPointer(true),
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

	suite.Run("Entire update is aborted when ProGearWeight is over max amount", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildNeedsServiceCounselingMove(suite.DB(), nil, nil).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := models.Int64Pointer(2001)
		proGearWeightSpouse := models.Int64Pointer(10)
		rmeWeight := models.Int64Pointer(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               &affiliation,
			DependentsAuthorized: models.BoolPointer(true),
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

	suite.Run("Entire update is aborted when ProGearWeightSpouse is over max amount", func() {
		moveRouter, err := move.NewMoveRouter()
		suite.FatalNoError(err)
		orderUpdater := NewOrderUpdater(moveRouter)
		order := factory.BuildNeedsServiceCounselingMove(suite.DB(), nil, nil).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.AffiliationAIRFORCE
		ocie := false
		proGearWeight := models.Int64Pointer(100)
		proGearWeightSpouse := models.Int64Pointer(501)
		rmeWeight := models.Int64Pointer(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               &affiliation,
			DependentsAuthorized: models.BoolPointer(true),
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
	moveRouter, err := move.NewMoveRouter()
	suite.FatalNoError(err)
	orderUpdater := NewOrderUpdater(moveRouter)

	setUpOrders := func(setUpPreExistingAmendedOrders bool) *models.Order {
		dutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2}),
				LinkOnly: true,
			},
		}, nil)
		var moves models.Moves

		customs := []factory.Customization{
			{
				Model:    dutyLocation,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
		}

		if setUpPreExistingAmendedOrders {
			customs = append(
				customs,
				factory.Customization{
					Model: models.Document{},
					Type:  &factory.Documents.UploadedAmendedOrders,
				},
			)
		}

		mto := factory.BuildServiceCounselingCompletedMove(suite.DB(), customs, nil)

		order := mto.Orders
		order.Moves = append(moves, mto)

		return &order
	}

	setUpFileToUpload := func() (*runtime.File, func()) {
		file := testdatagen.FixtureRuntimeFile("filled-out-orders.pdf")

		cleanUpFunc := func() {
			fileCloseErr := file.Close()
			suite.NoError(fileCloseErr)
		}

		return file, cleanUpFunc
	}

	suite.Run("Returns a NotFoundErr if the orders are not associated with the service member attempting to upload amended orders", func() {
		order := setUpOrders(false)
		otherServiceMember := factory.BuildExtendedServiceMember(suite.DB(), nil, nil)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: otherServiceMember.ID,
		})

		file, cleanUpFunc := setUpFileToUpload()
		defer cleanUpFunc()

		fakeS3 := storageTest.NewFakeS3Storage(true)

		upload, url, verrs, err := orderUpdater.UploadAmendedOrdersAsCustomer(
			appCtx,
			order.ServiceMember.UserID,
			order.ID,
			file.Data,
			file.Header.Filename,
			fakeS3,
		)

		if suite.Error(err) {
			suite.True(reflect.DeepEqual(models.Upload{}, upload), "Upload should be empty")
			suite.Equal("", url, "URL should be empty")
			suite.NoVerrs(verrs)

			suite.IsType(apperror.NotFoundError{}, err)
			suite.Contains(err.Error(), "not found while looking for order")
		}
	})

	suite.Run("Creates and saves new amendedOrder doc when the order.UploadedAmendedOrders is nil", func() {
		order := setUpOrders(false)
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: order.ServiceMemberID,
		})

		file, cleanUpFunc := setUpFileToUpload()
		defer cleanUpFunc()

		fakeS3 := storageTest.NewFakeS3Storage(true)

		suite.NotEqual(uuid.Nil, order.ServiceMemberID, "ServiceMember has ID that is not 0/empty")
		suite.NotEqual(uuid.Nil, order.ServiceMember.UserID, "ServiceMember.UserID has ID that is not 0/empty")

		upload, url, verrs, err := orderUpdater.UploadAmendedOrdersAsCustomer(
			appCtx,
			order.ServiceMember.UserID,
			order.ID,
			file.Data,
			file.Header.Filename,
			fakeS3)
		suite.NoError(err)
		suite.NoVerrs(verrs)

		expectedChecksum := "+XM59C3+hSg3Qrs0dPRuUhng5IQTWdYZtmcXhEH0SYU="
		if upload.Checksum != expectedChecksum {
			suite.Fail("Did not calculate the correct MD5: expected %s, got %s", expectedChecksum, upload.Checksum)
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
			suite.Fail("Couldn't find expected upload.")
		}
		suite.Equal(upload.ID.String(), findUpload.ID.String(), "found upload in db")
		suite.NotEmpty(url, "URL is populated")
	})

	suite.Run("Returns an error when order is not found", func() {
		order := setUpOrders(false)

		nonexistentOrdersUUID := uuid.Must(uuid.NewV4())

		// No need for a service member in this case because it'll fail on to
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: order.ServiceMemberID,
		})

		file, cleanUpFunc := setUpFileToUpload()
		defer cleanUpFunc()

		fakeS3 := storageTest.NewFakeS3Storage(true)

		_, _, verrs, err := orderUpdater.UploadAmendedOrdersAsCustomer(
			appCtx,
			order.ServiceMember.UserID,
			nonexistentOrdersUUID,
			file.Data,
			file.Header.Filename,
			fakeS3)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), "not found while looking for order")
		suite.NoVerrs(verrs)
	})

	suite.Run("Saves userUpload payload to order.UploadedAmendedOrders if the document already exists", func() {
		order := setUpOrders(true)

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.MilApp,
			ServiceMemberID: order.ServiceMemberID,
		})

		file, cleanUpFunc := setUpFileToUpload()
		defer cleanUpFunc()

		fakeS3 := storageTest.NewFakeS3Storage(true)

		suite.NotEqual(uuid.Nil, order.ServiceMemberID, "ServiceMember has ID that is not 0/empty")
		suite.NotEqual(uuid.Nil, order.ServiceMember.UserID, "ServiceMember.UserID has ID that is not 0/empty")

		_, _, verrs, err := orderUpdater.UploadAmendedOrdersAsCustomer(
			appCtx,
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
		suite.Equal(order.UploadedAmendedOrdersID, orderInDB.UploadedAmendedOrdersID)
		suite.NotNil(order.UploadedAmendedOrders)
		suite.NotNil(orderInDB.UploadedAmendedOrders)
	})
}
