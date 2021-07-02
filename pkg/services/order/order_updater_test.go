package order

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"

	storageTest "github.com/transcom/mymove/pkg/storage/test"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/handlers"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *OrderServiceSuite) TestUpdateOrderAsTOO() {
	suite.T().Run("Returns an error when order is not found", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.UpdateOrderPayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateOrderAsTOO(nonexistentUUID, payload, eTag)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when origin duty station is not found", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders
		newDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.UpdateOrderPayload{
			NewDutyStationID:    handlers.FmtUUID(newDutyStation.ID),
			OriginDutyStationID: handlers.FmtUUID(nonexistentUUID),
		}
		eTag := etag.GenerateEtag(order.UpdatedAt)

		_, _, err := orderUpdater.UpdateOrderAsTOO(order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when new duty station is not found", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders
		originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.UpdateOrderPayload{
			NewDutyStationID:    handlers.FmtUUID(nonexistentUUID),
			OriginDutyStationID: handlers.FmtUUID(originDutyStation.ID),
		}
		eTag := etag.GenerateEtag(order.UpdatedAt)

		_, _, err := orderUpdater.UpdateOrderAsTOO(order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when the etag does not match", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders

		payload := ghcmessages.UpdateOrderPayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateOrderAsTOO(order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("Updates the order when all fields are valid", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		order := testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{}).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
		updatedOriginDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
		ordersType := ghcmessages.OrdersTypeSEPARATION
		deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
		ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.UpdateOrderPayload{
			DepartmentIndicator: &deptIndicator,
			IssueDate:           &dateIssued,
			NewDutyStationID:    handlers.FmtUUID(updatedDestinationDutyStation.ID),
			OriginDutyStationID: handlers.FmtUUID(updatedOriginDutyStation.ID),
			OrdersNumber:        handlers.FmtString("ORDER100"),
			OrdersType:          ordersType,
			OrdersTypeDetail:    &ordersTypeDetail,
			ReportByDate:        &reportByDate,
			Tac:                 handlers.FmtString("E19A"),
			Sac:                 handlers.FmtString("987654321"),
		}

		updatedOrder, _, err := orderUpdater.UpdateOrderAsTOO(order.ID, payload, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)

		fetchedSM := models.ServiceMember{}
		_ = suite.DB().EagerPreload("DutyStation").Find(&fetchedSM, order.ServiceMember.ID)

		suite.NoError(err)
		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(payload.NewDutyStationID.String(), updatedOrder.NewDutyStation.ID.String())
		suite.Equal(payload.OriginDutyStationID.String(), updatedOrder.OriginDutyStation.ID.String())
		suite.Equal(time.Time(*payload.IssueDate), updatedOrder.IssueDate)
		suite.Equal(time.Time(*payload.ReportByDate), updatedOrder.ReportByDate)
		suite.EqualValues(payload.OrdersType, updatedOrder.OrdersType)
		suite.EqualValues(payload.OrdersTypeDetail, updatedOrder.OrdersTypeDetail)
		suite.Equal(payload.OrdersNumber, updatedOrder.OrdersNumber)
		suite.EqualValues(payload.DepartmentIndicator, updatedOrder.DepartmentIndicator)
		suite.Equal(payload.Tac, updatedOrder.TAC)
		suite.Equal(payload.Sac, updatedOrder.SAC)
		suite.EqualValues(&updatedOriginDutyStation.ID, fetchedSM.DutyStationID)
		suite.EqualValues(updatedOriginDutyStation.ID, fetchedSM.DutyStation.ID)
		suite.EqualValues(updatedOriginDutyStation.Name, fetchedSM.DutyStation.Name)
	})

	suite.T().Run("Rolls back transaction if Order is invalid", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		order := testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{}).Orders

		emptyStrSAC := ""
		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))
		updatedDestinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
		updatedOriginDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
		ordersType := ghcmessages.OrdersTypeSEPARATION
		deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
		ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.UpdateOrderPayload{
			DepartmentIndicator: &deptIndicator,
			IssueDate:           &dateIssued,
			NewDutyStationID:    handlers.FmtUUID(updatedDestinationDutyStation.ID),
			OriginDutyStationID: handlers.FmtUUID(updatedOriginDutyStation.ID),
			OrdersNumber:        handlers.FmtString("ORDER100"),
			OrdersType:          ordersType,
			OrdersTypeDetail:    &ordersTypeDetail,
			ReportByDate:        &reportByDate,
			Tac:                 handlers.FmtString("E19A"),
			Sac:                 &emptyStrSAC, // this will trigger a validation error on Order model
		}

		updatedOrder, _, err := orderUpdater.UpdateOrderAsTOO(order.ID, payload, eTag)

		// check that we get back a validation error
		suite.EqualError(err, fmt.Sprintf("Invalid input for id: %s. SAC can not be blank.", order.ID))
		suite.Nil(updatedOrder)
		suite.IsType(services.InvalidInputError{}, err)
	})
}

func (suite *OrderServiceSuite) TestUpdateOrderAsCounselor() {
	suite.T().Run("Returns an error when order is not found", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.CounselingUpdateOrderPayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateOrderAsCounselor(nonexistentUUID, payload, eTag)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when the etag does not match", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders

		payload := ghcmessages.CounselingUpdateOrderPayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateOrderAsCounselor(order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("Updates the order when it is found", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		order := testdatagen.MakeNeedsServiceCounselingMove(suite.DB()).Orders

		dateIssued := strfmt.Date(time.Now().Add(-48 * time.Hour))
		updatedDestinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
		updatedOriginDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
		ordersType := ghcmessages.OrdersTypeSEPARATION
		reportByDate := strfmt.Date(time.Now().Add(72 * time.Hour))

		body := ghcmessages.CounselingUpdateOrderPayload{
			IssueDate:           &dateIssued,
			NewDutyStationID:    handlers.FmtUUID(updatedDestinationDutyStation.ID),
			OriginDutyStationID: handlers.FmtUUID(updatedOriginDutyStation.ID),
			OrdersType:          ordersType,
			ReportByDate:        &reportByDate,
		}

		eTag := etag.GenerateEtag(order.UpdatedAt)

		updatedOrder, _, err := orderUpdater.UpdateOrderAsCounselor(order.ID, body, eTag)
		suite.NoError(err)

		var orderInDB models.Order
		err = suite.DB().Find(&orderInDB, order.ID)

		suite.NoError(err)
		suite.Equal(order.ID.String(), updatedOrder.ID.String())
		suite.Equal(body.NewDutyStationID.String(), updatedOrder.NewDutyStation.ID.String())
		suite.Equal(body.OriginDutyStationID.String(), updatedOrder.OriginDutyStation.ID.String())
		suite.Equal(time.Time(*body.IssueDate), updatedOrder.IssueDate)
		suite.Equal(time.Time(*body.ReportByDate), updatedOrder.ReportByDate)
		suite.EqualValues(body.OrdersType, updatedOrder.OrdersType)
	})
}

func (suite *OrderServiceSuite) TestUpdateAllowanceAsTOO() {
	suite.T().Run("Returns an error when order is not found", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.UpdateAllowancePayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateAllowanceAsTOO(nonexistentUUID, payload, eTag)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when the etag does not match", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders

		payload := ghcmessages.UpdateAllowancePayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateAllowanceAsTOO(order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("Updates the allowance when all fields are valid", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		order := testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{}).Orders

		newAuthorizedWeight := int64(10000)
		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.BranchAIRFORCE
		ocie := false
		proGearWeight := swag.Int64(100)
		proGearWeightSpouse := swag.Int64(10)
		rmeWeight := swag.Int64(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.UpdateAllowancePayload{
			Agency:               affiliation,
			AuthorizedWeight:     &newAuthorizedWeight,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
			OrganizationalClothingAndIndividualEquipment: &ocie,
			ProGearWeight:                  proGearWeight,
			ProGearWeightSpouse:            proGearWeightSpouse,
			RequiredMedicalEquipmentWeight: rmeWeight,
		}

		updatedOrder, _, err := orderUpdater.UpdateAllowanceAsTOO(order.ID, payload, eTag)
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
		orderUpdater := NewOrderUpdater(suite.DB())
		nonexistentUUID := uuid.Must(uuid.NewV4())

		payload := ghcmessages.CounselingUpdateAllowancePayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateAllowanceAsCounselor(nonexistentUUID, payload, eTag)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when the etag does not match", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		order := testdatagen.MakeDefaultMove(suite.DB()).Orders

		payload := ghcmessages.CounselingUpdateAllowancePayload{}
		eTag := ""

		_, _, err := orderUpdater.UpdateAllowanceAsCounselor(order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("Updates the allowance when all fields are valid", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		order := testdatagen.MakeNeedsServiceCounselingMove(suite.DB()).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.BranchAIRFORCE
		ocie := false
		proGearWeight := swag.Int64(100)
		proGearWeightSpouse := swag.Int64(10)
		rmeWeight := swag.Int64(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               affiliation,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
			OrganizationalClothingAndIndividualEquipment: &ocie,
			ProGearWeight:                  proGearWeight,
			ProGearWeightSpouse:            proGearWeightSpouse,
			RequiredMedicalEquipmentWeight: rmeWeight,
		}

		updatedOrder, _, err := orderUpdater.UpdateAllowanceAsCounselor(order.ID, payload, eTag)
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
		suite.EqualValues(payload.Agency, *fetchedSM.Affiliation)
	})

	suite.T().Run("Updates the allowance when move needs service counseling and order fields are missing", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		orderWithoutDefaults := testdatagen.MakeOrderWithoutDefaults(suite.DB(), testdatagen.Assertions{})
		move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusNeedsServiceCounseling,
			},
			Order: orderWithoutDefaults,
		})

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.BranchAIRFORCE
		ocie := false
		proGearWeight := swag.Int64(100)
		proGearWeightSpouse := swag.Int64(10)
		rmeWeight := swag.Int64(10000)
		eTag := etag.GenerateEtag(orderWithoutDefaults.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               affiliation,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
			OrganizationalClothingAndIndividualEquipment: &ocie,
			ProGearWeight:                  proGearWeight,
			ProGearWeightSpouse:            proGearWeightSpouse,
			RequiredMedicalEquipmentWeight: rmeWeight,
		}

		updatedOrder, _, err := orderUpdater.UpdateAllowanceAsCounselor(orderWithoutDefaults.ID, payload, eTag)
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
		suite.EqualValues(payload.Agency, *fetchedSM.Affiliation)

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
		orderUpdater := NewOrderUpdater(suite.DB())
		order := testdatagen.MakeNeedsServiceCounselingMove(suite.DB()).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.BranchAIRFORCE
		ocie := false
		proGearWeight := swag.Int64(2001)
		proGearWeightSpouse := swag.Int64(10)
		rmeWeight := swag.Int64(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               affiliation,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
			OrganizationalClothingAndIndividualEquipment: &ocie,
			ProGearWeight:                  proGearWeight,
			ProGearWeightSpouse:            proGearWeightSpouse,
			RequiredMedicalEquipmentWeight: rmeWeight,
		}

		updatedOrder, _, err := orderUpdater.UpdateAllowanceAsCounselor(order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)

		var orderInDB models.Order
		err = suite.DB().EagerPreload("Entitlement").Find(&orderInDB, order.ID)

		suite.NoError(err)
		suite.NotEqual(payload.ProGearWeight, orderInDB.Entitlement.ProGearWeight)
		suite.Nil(updatedOrder)
	})

	suite.T().Run("Entire update is aborted when ProGearWeightSpouse is over max amount", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		order := testdatagen.MakeNeedsServiceCounselingMove(suite.DB()).Orders

		grade := ghcmessages.GradeO5
		affiliation := ghcmessages.BranchAIRFORCE
		ocie := false
		proGearWeight := swag.Int64(100)
		proGearWeightSpouse := swag.Int64(501)
		rmeWeight := swag.Int64(10000)
		eTag := etag.GenerateEtag(order.UpdatedAt)

		payload := ghcmessages.CounselingUpdateAllowancePayload{
			Agency:               affiliation,
			DependentsAuthorized: swag.Bool(true),
			Grade:                &grade,
			OrganizationalClothingAndIndividualEquipment: &ocie,
			ProGearWeight:                  proGearWeight,
			ProGearWeightSpouse:            proGearWeightSpouse,
			RequiredMedicalEquipmentWeight: rmeWeight,
		}

		updatedOrder, _, err := orderUpdater.UpdateAllowanceAsCounselor(order.ID, payload, eTag)

		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)

		var orderInDB models.Order
		err = suite.DB().EagerPreload("Entitlement").Find(&orderInDB, order.ID)

		suite.NoError(err)
		suite.NotEqual(payload.ProGearWeightSpouse, orderInDB.Entitlement.ProGearWeightSpouse)
		suite.Nil(updatedOrder)
	})
}

func (suite *OrderServiceSuite) TestUploadAmendedOrdersForCustomer() {

	suite.T().Run("Creates and saves new amendedOrder doc when the order.UploadedAmendedOrders is nil", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		dutyStation := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
			DutyStation: models.DutyStation{
				Address: testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{}),
			},
		})
		var moves models.Moves
		mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})

		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				OriginDutyStation: &dutyStation,
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

		logger, zapErr := zap.NewDevelopment()
		suite.NoError(zapErr)

		suite.NotEqual(uuid.Nil, order.ServiceMemberID, "ServiceMember has ID that is not 0/empty")
		suite.NotEqual(uuid.Nil, order.ServiceMember.UserID, "ServiceMember.UserID has ID that is not 0/empty")

		upload, url, verrs, err := orderUpdater.UploadAmendedOrdersAsCustomer(
			logger,
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
		orderUpdater := NewOrderUpdater(suite.DB())
		nonexistentUUID := uuid.Must(uuid.NewV4())

		file := testdatagen.FixtureRuntimeFile("test.pdf")
		defer func() {
			fileCloseErr := file.Close()
			suite.NoError(fileCloseErr)
		}()

		fakeS3 := storageTest.NewFakeS3Storage(true)

		logger, zapErr := zap.NewDevelopment()
		suite.NoError(zapErr)

		_, _, verrs, err := orderUpdater.UploadAmendedOrdersAsCustomer(
			logger,
			nonexistentUUID,
			nonexistentUUID,
			file.Data,
			file.Header.Filename,
			fakeS3)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
		suite.Contains(err.Error(), "not found while looking for order")
		suite.NoVerrs(verrs)
	})

	suite.T().Run("Saves userUpload payload to order.UploadedAmendedOrders if the document already exists", func(t *testing.T) {
		orderUpdater := NewOrderUpdater(suite.DB())
		dutyStation := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
			DutyStation: models.DutyStation{
				Address: testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{}),
			},
		})
		var moves models.Moves
		mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{})

		document := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{})
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				OriginDutyStation:       &dutyStation,
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

		logger, zapErr := zap.NewDevelopment()
		suite.NoError(zapErr)

		suite.NotEqual(uuid.Nil, order.ServiceMemberID, "ServiceMember has ID that is not 0/empty")
		suite.NotEqual(uuid.Nil, order.ServiceMember.UserID, "ServiceMember.UserID has ID that is not 0/empty")

		_, _, verrs, err := orderUpdater.UploadAmendedOrdersAsCustomer(
			logger,
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
