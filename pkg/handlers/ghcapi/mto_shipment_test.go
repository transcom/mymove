package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	shipmentorchestrator "github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/swagger/nullable"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/trace"

	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/go-openapi/strfmt"

	"github.com/gofrs/uuid"

	routemocks "github.com/transcom/mymove/pkg/route/mocks"

	"github.com/gobuffalo/validate/v3"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	shipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/shipment"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"

	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type listMTOShipmentsSubtestData struct {
	mtoAgent       models.MTOAgent
	mtoServiceItem models.MTOServiceItem
	shipments      models.MTOShipments
	params         mtoshipmentops.ListMTOShipmentsParams
	sitExtension   models.SITExtension
	sit            models.MTOServiceItem
}

func (suite *HandlerSuite) makeListMTOShipmentsSubtestData() (subtestData *listMTOShipmentsSubtestData) {
	subtestData = &listMTOShipmentsSubtestData{}

	mto := testdatagen.MakeDefaultMove(suite.DB())

	storageFacility := testdatagen.MakeDefaultStorageFacility(suite.DB())

	sitAllowance := int(90)
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status:           models.MTOShipmentStatusApproved,
			CounselorRemarks: handlers.FmtString("counselor remark"),
			SITDaysAllowance: &sitAllowance,
			StorageFacility:  &storageFacility,
		},
	})

	secondShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	// third shipment with destination address and type
	destinationAddress := testdatagen.MakeDefaultAddress(suite.DB())
	destinationType := models.DestinationTypeHomeOfRecord
	thirdShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusSubmitted,
			DestinationAddressID: &destinationAddress.ID,
			DestinationType:      &destinationType,
		},
	})

	subtestData.mtoAgent = testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipmentID: mtoShipment.ID,
		},
	})
	subtestData.mtoServiceItem = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			MTOShipmentID: &mtoShipment.ID,
		},
	})

	ppm := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
	})

	// testdatagen.MakeDOFSITReService(suite.DB(), testdatagen.Assertions{})

	year, month, day := time.Now().Date()
	lastMonthEntry := time.Date(year, month, day-37, 0, 0, 0, 0, time.UTC)
	lastMonthDeparture := time.Date(year, month, day-30, 0, 0, 0, 0, time.UTC)
	testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate:     &lastMonthEntry,
			SITDepartureDate: &lastMonthDeparture,
			Status:           models.MTOServiceItemStatusApproved,
		},
		Move:        mto,
		MTOShipment: mtoShipment,
		ReService: models.ReService{
			Code: models.ReServiceCodeDOPSIT,
		},
	})

	aWeekAgo := time.Date(year, month, day-7, 0, 0, 0, 0, time.UTC)
	departureDate := aWeekAgo.Add(time.Hour * 24 * 30)
	subtestData.sit = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate:     &aWeekAgo,
			SITDepartureDate: &departureDate,
			Status:           models.MTOServiceItemStatusApproved,
		},
		Move:        mto,
		MTOShipment: mtoShipment,
		ReService: models.ReService{
			Code: models.ReServiceCodeDOPSIT,
		},
	})

	subtestData.sitExtension = testdatagen.MakeSITExtension(suite.DB(), testdatagen.Assertions{
		SITExtension: models.SITExtension{
			MTOShipmentID: mtoShipment.ID,
		},
	})

	subtestData.shipments = models.MTOShipments{mtoShipment, secondShipment, thirdShipment, ppm.Shipment}
	requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{})

	req := httptest.NewRequest("GET", fmt.Sprintf("/move_task_orders/%s/mto_shipments", mto.ID.String()), nil)
	req = suite.AuthenticateOfficeRequest(req, requestUser)

	subtestData.params = mtoshipmentops.ListMTOShipmentsParams{
		HTTPRequest:     req,
		MoveTaskOrderID: *handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
	}

	return subtestData
}

func (suite *HandlerSuite) TestListMTOShipmentsHandler() {
	suite.Run("Successful list fetch - Integration Test", func() {
		subtestData := suite.makeListMTOShipmentsSubtestData()
		params := subtestData.params
		shipments := subtestData.shipments
		mtoAgent := subtestData.mtoAgent
		mtoServiceItem := subtestData.mtoServiceItem
		sitExtension := subtestData.sitExtension

		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			mtoshipment.NewMTOShipmentFetcher(),
			mtoshipment.NewShipmentSITStatus(),
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsOK{}, response)

		okResponse := response.(*mtoshipmentops.ListMTOShipmentsOK)
		suite.Len(okResponse.Payload, 4)

		payloadShipment := okResponse.Payload[0]
		suite.Equal(shipments[0].ID.String(), payloadShipment.ID.String())
		suite.Equal(*shipments[0].CounselorRemarks, *payloadShipment.CounselorRemarks)
		suite.Equal(mtoAgent.ID.String(), payloadShipment.MtoAgents[0].ID.String())
		suite.Equal(mtoServiceItem.ID.String(), payloadShipment.MtoServiceItems[0].ID.String())
		suite.Equal(sitExtension.ID.String(), payloadShipment.SitExtensions[0].ID.String())
		suite.Equal(shipments[0].StorageFacility.ID.String(), payloadShipment.StorageFacility.ID.String())
		suite.Equal(shipments[0].StorageFacility.Address.ID.String(), payloadShipment.StorageFacility.Address.ID.String())

		payloadShipment2 := okResponse.Payload[1]
		suite.Equal(shipments[1].ID.String(), payloadShipment2.ID.String())
		suite.Nil(payloadShipment2.StorageFacility)

		suite.Equal(int64(190), *payloadShipment.SitDaysAllowance)
		suite.Equal(mtoshipment.OriginSITLocation, payloadShipment.SitStatus.Location)
		suite.Equal(int64(7), *payloadShipment.SitStatus.DaysInSIT)
		suite.Equal(int64(176), *payloadShipment.SitStatus.TotalDaysRemaining)
		suite.Equal(int64(14), *payloadShipment.SitStatus.TotalSITDaysUsed) // 7 from the previous SIT and 7 from the current
		suite.Equal(subtestData.sit.SITEntryDate.Format(strfmt.MarshalFormat), payloadShipment.SitStatus.SitEntryDate.String())
		suite.Equal(subtestData.sit.SITDepartureDate.Format(strfmt.MarshalFormat), payloadShipment.SitStatus.SitDepartureDate.String())

		suite.Len(payloadShipment.SitStatus.PastSITServiceItems, 1)
		year, month, day := time.Now().Date()
		lastMonthEntry := time.Date(year, month, day-37, 0, 0, 0, 0, time.UTC)
		suite.Equal(lastMonthEntry.Format(strfmt.MarshalFormat), payloadShipment.SitStatus.PastSITServiceItems[0].SitEntryDate.String())

		// This one has a destination shipment type
		payloadShipment3 := okResponse.Payload[2]
		suite.Equal(string(models.DestinationTypeHomeOfRecord), string(*payloadShipment3.DestinationType))

		payloadShipment4 := okResponse.Payload[3]
		suite.NotNil(payloadShipment4.PpmShipment)
		suite.Equal(shipments[3].ID.String(), payloadShipment4.PpmShipment.ShipmentID.String())
		suite.Equal(shipments[3].PPMShipment.ID.String(), payloadShipment4.PpmShipment.ID.String())
	})

	suite.Run("Failure list fetch - Internal Server Error", func() {
		subtestData := suite.makeListMTOShipmentsSubtestData()
		params := subtestData.params
		mockMTOShipmentFetcher := &mocks.MTOShipmentFetcher{}

		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			mockMTOShipmentFetcher,
			mtoshipment.NewShipmentSITStatus(),
		}

		mockMTOShipmentFetcher.On("ListMTOShipments", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID")).Return(nil, apperror.NewQueryError("MTOShipment", errors.New("query error"), ""))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsInternalServerError{}, response)
	})

	suite.Run("Failure list fetch - 404 Not Found - Move Task Order ID", func() {
		subtestData := suite.makeListMTOShipmentsSubtestData()
		params := subtestData.params

		mockMTOShipmentFetcher := &mocks.MTOShipmentFetcher{}

		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			mockMTOShipmentFetcher,
			mtoshipment.NewShipmentSITStatus(),
		}

		mockMTOShipmentFetcher.On("ListMTOShipments", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID")).Return(nil, apperror.NewNotFoundError(uuid.FromStringOrNil(params.MoveTaskOrderID.String()), "move not found"))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsNotFound{}, response)
	})
}

func (suite *HandlerSuite) TestDeleteShipmentHandler() {
	suite.Run("Returns a 403 when the office user is not a service counselor", func() {
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		uuid := uuid.Must(uuid.NewV4())
		deleter := &mocks.ShipmentDeleter{}

		deleter.AssertNumberOfCalls(suite.T(), "DeleteShipment", 0)

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/shipments/%s", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := shipmentops.DeleteShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
		}

		response := handler.Handle(deletionParams)
		suite.IsType(&shipmentops.DeleteShipmentForbidden{}, response)
	})

	suite.Run("Returns 204 when all validations pass", func() {
		shipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())
		officeUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(shipment.MoveTaskOrderID, nil)

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := shipmentops.DeleteShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(deletionParams)

		suite.IsType(&shipmentops.DeleteShipmentNoContent{}, response)
	})

	suite.Run("Returns 404 when deleter returns NotFoundError", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		officeUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.NotFoundError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := shipmentops.DeleteShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(deletionParams)
		suite.IsType(&shipmentops.DeleteShipmentNotFound{}, response)
	})

	suite.Run("Returns 403 when deleter returns ForbiddenError", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		officeUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.ForbiddenError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := shipmentops.DeleteShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(deletionParams)
		suite.IsType(&shipmentops.DeleteShipmentForbidden{}, response)
	})

	suite.Run("Returns 422 - Unprocessable Enitity error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		officeUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.UnprocessableEntityError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := shipmentops.DeleteShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(deletionParams)
		suite.IsType(&shipmentops.DeleteShipmentUnprocessableEntity{}, response)
	})

	suite.Run("Returns 409 - Conflict error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		officeUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.ConflictError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := shipmentops.DeleteShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(deletionParams)
		suite.IsType(&shipmentops.DeleteShipmentConflict{}, response)
	})

}

func (suite *HandlerSuite) TestApproveShipmentHandler() {
	suite.Run("Returns 200 when all validations pass", func() {
		move := testdatagen.MakeAvailableMove(suite.DB())
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		})
		// Populate the reServices table with codes needed by the
		// HHG_LONGHAUL_DOMESTIC shipment type
		reServiceCodes := []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDPK,
			models.ReServiceCodeDUPK,
		}
		for _, serviceCode := range reServiceCodes {
			testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
				ReService: models.ReService{
					Code:      serviceCode,
					Name:      "test",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			})
		}

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		builder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		approver := mtoshipment.NewShipmentApprover(
			mtoshipment.NewShipmentRouter(),
			mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter),
			&routemocks.Planner{},
		)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}

		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentOK{}, response)
		suite.HasWebhookNotification(shipment.ID, traceID)
	})

	suite.Run("Returns a 403 when the office user is not a TOO", func() {
		officeUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		uuid := uuid.Must(uuid.NewV4())
		approver := &mocks.ShipmentApprover{}

		approver.AssertNumberOfCalls(suite.T(), "ApproveShipment", 0)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
			IfMatch:     etag.GenerateEtag(time.Now()),
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentForbidden{}, response)
	})

	suite.Run("Returns 404 when approver returns NotFoundError", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		approver := &mocks.ShipmentApprover{}

		approver.On("ApproveShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.NotFoundError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentNotFound{}, response)
	})

	suite.Run("Returns 409 when approver returns Conflict Error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		approver := &mocks.ShipmentApprover{}

		approver.On("ApproveShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.ConflictError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentConflict{}, response)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(time.Now())
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		approver := &mocks.ShipmentApprover{}

		approver.On("ApproveShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.PreconditionFailedError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentPreconditionFailed{}, response)
	})

	suite.Run("Returns 422 when approver returns validation errors", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		approver := &mocks.ShipmentApprover{}

		approver.On("ApproveShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.InvalidInputError{ValidationErrors: &validate.Errors{}})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentUnprocessableEntity{}, response)
	})

	suite.Run("Returns 500 when approver returns unexpected error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		approver := &mocks.ShipmentApprover{}

		approver.On("ApproveShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, errors.New("UnexpectedError"))

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestRequestShipmentDiversionHandler() {
	suite.Run("Returns 200 when all validations pass", func() {
		move := testdatagen.MakeAvailableMove(suite.DB())
		shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
			Move: move,
		})

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		requester := mtoshipment.NewShipmentDiversionRequester(
			mtoshipment.NewShipmentRouter(),
		)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			mtoshipment.NewShipmentSITStatus(),
		}

		approveParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionOK{}, response)
		suite.HasWebhookNotification(shipment.ID, traceID)
	})

	suite.Run("Returns a 403 when the office user is not a TOO", func() {
		officeUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		uuid := uuid.Must(uuid.NewV4())
		requester := &mocks.ShipmentDiversionRequester{}

		requester.AssertNumberOfCalls(suite.T(), "RequestShipmentDiversion", 0)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
			IfMatch:     etag.GenerateEtag(time.Now()),
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionForbidden{}, response)
	})

	suite.Run("Returns 404 when requester returns NotFoundError", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		requester := &mocks.ShipmentDiversionRequester{}

		requester.On("RequestShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.NotFoundError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionNotFound{}, response)
	})

	suite.Run("Returns 409 when requester returns Conflict Error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		requester := &mocks.ShipmentDiversionRequester{}

		requester.On("RequestShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, mtoshipment.ConflictStatusError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionConflict{}, response)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(time.Now())
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		requester := &mocks.ShipmentDiversionRequester{}

		requester.On("RequestShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.PreconditionFailedError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionPreconditionFailed{}, response)
	})

	suite.Run("Returns 422 when requester returns validation errors", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		requester := &mocks.ShipmentDiversionRequester{}

		requester.On("RequestShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.InvalidInputError{ValidationErrors: &validate.Errors{}})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionUnprocessableEntity{}, response)
	})

	suite.Run("Returns 500 when requester returns unexpected error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		requester := &mocks.ShipmentDiversionRequester{}

		requester.On("RequestShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, errors.New("UnexpectedError"))

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestApproveShipmentDiversionHandler() {
	suite.Run("Returns 200 when all validations pass", func() {
		move := testdatagen.MakeAvailableMove(suite.DB())
		shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:    models.MTOShipmentStatusSubmitted,
				Diversion: true,
			},
			Move: move,
		})

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		approver := mtoshipment.NewShipmentDiversionApprover(
			mtoshipment.NewShipmentRouter(),
		)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}

		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionOK{}, response)
		suite.HasWebhookNotification(shipment.ID, traceID)
	})

	suite.Run("Returns a 403 when the office user is not a TOO", func() {
		officeUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		uuid := uuid.Must(uuid.NewV4())
		approver := &mocks.ShipmentDiversionApprover{}

		approver.AssertNumberOfCalls(suite.T(), "ApproveShipmentDiversion", 0)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
			IfMatch:     etag.GenerateEtag(time.Now()),
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionForbidden{}, response)
	})

	suite.Run("Returns 404 when approver returns NotFoundError", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		approver := &mocks.ShipmentDiversionApprover{}

		approver.On("ApproveShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.NotFoundError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionNotFound{}, response)
	})

	suite.Run("Returns 409 when approver returns Conflict Error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		approver := &mocks.ShipmentDiversionApprover{}

		approver.On("ApproveShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, mtoshipment.ConflictStatusError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionConflict{}, response)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(time.Now())
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		approver := &mocks.ShipmentDiversionApprover{}

		approver.On("ApproveShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.PreconditionFailedError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionPreconditionFailed{}, response)
	})

	suite.Run("Returns 422 when approver returns validation errors", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		approver := &mocks.ShipmentDiversionApprover{}

		approver.On("ApproveShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.InvalidInputError{ValidationErrors: &validate.Errors{}})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionUnprocessableEntity{}, response)
	})

	suite.Run("Returns 500 when approver returns unexpected error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		approver := &mocks.ShipmentDiversionApprover{}

		approver.On("ApproveShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, errors.New("UnexpectedError"))

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestRejectShipmentHandler() {
	reason := "reason"

	suite.Run("Returns 200 when all validations pass", func() {
		move := testdatagen.MakeAvailableMove(suite.DB())
		shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: move,
		})

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		rejecter := mtoshipment.NewShipmentRejecter(
			mtoshipment.NewShipmentRouter(),
		)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}

		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentOK{}, response)
		suite.HasWebhookNotification(shipment.ID, traceID)
	})

	suite.Run("Returns a 403 when the office user is not a TOO", func() {
		officeUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		uuid := uuid.Must(uuid.NewV4())
		rejecter := &mocks.ShipmentRejecter{}

		rejecter.AssertNumberOfCalls(suite.T(), "RejectShipment", 0)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
			IfMatch:     etag.GenerateEtag(time.Now()),
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentForbidden{}, response)
	})

	suite.Run("Returns 404 when rejecter returns NotFoundError", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		rejecter := &mocks.ShipmentRejecter{}

		rejecter.On("RejectShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &reason).Return(nil, apperror.NotFoundError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentNotFound{}, response)
	})

	suite.Run("Returns 409 when rejecter returns Conflict Error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		rejecter := &mocks.ShipmentRejecter{}

		rejecter.On("RejectShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &reason).Return(nil, mtoshipment.ConflictStatusError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentConflict{}, response)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(time.Now())
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		rejecter := &mocks.ShipmentRejecter{}

		rejecter.On("RejectShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &reason).Return(nil, apperror.PreconditionFailedError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentPreconditionFailed{}, response)
	})

	suite.Run("Returns 422 when rejecter returns validation errors", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		rejecter := &mocks.ShipmentRejecter{}

		rejecter.On("RejectShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &reason).Return(nil, apperror.InvalidInputError{ValidationErrors: &validate.Errors{}})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentUnprocessableEntity{}, response)
	})

	suite.Run("Returns 500 when rejecter returns unexpected error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		rejecter := &mocks.ShipmentRejecter{}

		rejecter.On("RejectShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &reason).Return(nil, errors.New("UnexpectedError"))

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentInternalServerError{}, response)
	})

	suite.Run("Requires rejection reason in Body of request", func() {
		move := testdatagen.MakeAvailableMove(suite.DB())
		shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: move,
		})
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		rejecter := mtoshipment.NewShipmentRejecter(
			mtoshipment.NewShipmentRouter(),
		)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body:        &ghcmessages.RejectShipment{},
		}

		suite.Error(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentUnprocessableEntity{}, response)
	})
}

func (suite *HandlerSuite) TestRequestShipmentCancellationHandler() {
	suite.Run("Returns 200 when all validations pass", func() {
		move := testdatagen.MakeAvailableMove(suite.DB())
		shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
			Move: move,
		})

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		canceler := mtoshipment.NewShipmentCancellationRequester(
			mtoshipment.NewShipmentRouter(),
		)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			mtoshipment.NewShipmentSITStatus(),
		}

		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationOK{}, response)
		suite.HasWebhookNotification(shipment.ID, traceID)
	})

	suite.Run("Returns a 403 when the office user is not a TOO", func() {
		officeUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		uuid := uuid.Must(uuid.NewV4())
		canceler := &mocks.ShipmentCancellationRequester{}

		canceler.AssertNumberOfCalls(suite.T(), "RequestShipmentCancellation", 0)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
			IfMatch:     etag.GenerateEtag(time.Now()),
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationForbidden{}, response)
	})

	suite.Run("Returns 404 when canceler returns NotFoundError", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		canceler := &mocks.ShipmentCancellationRequester{}

		canceler.On("RequestShipmentCancellation", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.NotFoundError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationNotFound{}, response)
	})

	suite.Run("Returns 409 when canceler returns Conflict Error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		canceler := &mocks.ShipmentCancellationRequester{}

		canceler.On("RequestShipmentCancellation", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, mtoshipment.ConflictStatusError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationConflict{}, response)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(time.Now())
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		canceler := &mocks.ShipmentCancellationRequester{}

		canceler.On("RequestShipmentCancellation", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.PreconditionFailedError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationPreconditionFailed{}, response)
	})

	suite.Run("Returns 422 when canceler returns validation errors", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		canceler := &mocks.ShipmentCancellationRequester{}

		canceler.On("RequestShipmentCancellation", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.InvalidInputError{ValidationErrors: &validate.Errors{}})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationUnprocessableEntity{}, response)
	})

	suite.Run("Returns 500 when canceler returns unexpected error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		canceler := &mocks.ShipmentCancellationRequester{}

		canceler.On("RequestShipmentCancellation", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, errors.New("UnexpectedError"))

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			mtoshipment.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestRequestShipmentReweighHandler() {
	suite.Run("Returns 200 when all validations pass", func() {
		move := testdatagen.MakeAvailableMove(suite.DB())
		shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
			Move: move,
		})

		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		reweighRequester := mtoshipment.NewShipmentReweighRequester()

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-reweigh", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		handlerConfig.SetNotificationSender(suite.TestNotificationSender())
		planner := &routemocks.Planner{}
		planner.On("TransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveRouter := moverouter.NewMoveRouter()
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)

		handler := RequestShipmentReweighHandler{
			handlerConfig,
			reweighRequester,
			mtoshipment.NewShipmentSITStatus(),
			updater,
		}

		approveParams := shipmentops.RequestShipmentReweighParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		response := handler.Handle(approveParams)
		okResponse := response.(*shipmentops.RequestShipmentReweighOK)
		payload := okResponse.Payload
		suite.IsType(&shipmentops.RequestShipmentReweighOK{}, response)
		suite.Equal(strfmt.UUID(shipment.ID.String()), payload.ShipmentID)
		suite.EqualValues(models.ReweighRequesterTOO, payload.RequestedBy)
		suite.WithinDuration(time.Now(), (time.Time)(payload.RequestedAt), 2*time.Second)
		suite.HasWebhookNotification(shipment.ID, traceID)
	})

	suite.Run("Returns a 403 when the office user is not a TOO", func() {
		officeUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		uuid := uuid.Must(uuid.NewV4())
		reweighRequester := &mocks.ShipmentReweighRequester{}

		reweighRequester.AssertNumberOfCalls(suite.T(), "RequestShipmentReweigh", 0)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-reweigh", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		planner := &routemocks.Planner{}
		planner.On("TransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveRouter := moverouter.NewMoveRouter()
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)

		handler := RequestShipmentReweighHandler{
			handlerConfig,
			reweighRequester,
			mtoshipment.NewShipmentSITStatus(),
			updater,
		}
		approveParams := shipmentops.RequestShipmentReweighParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
		}

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentReweighForbidden{}, response)
	})

	suite.Run("Returns 404 when reweighRequester returns NotFoundError", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		reweighRequester := &mocks.ShipmentReweighRequester{}

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-reweigh", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		planner := &routemocks.Planner{}
		planner.On("TransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveRouter := moverouter.NewMoveRouter()
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)

		handler := RequestShipmentReweighHandler{
			handlerConfig,
			reweighRequester,
			mtoshipment.NewShipmentSITStatus(),
			updater,
		}
		params := shipmentops.RequestShipmentReweighParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}
		reweighRequester.On("RequestShipmentReweigh", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, models.ReweighRequesterTOO).Return(nil, apperror.NotFoundError{})

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RequestShipmentReweighNotFound{}, response)
	})

	suite.Run("Returns 409 when reweighRequester returns Conflict Error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		reweighRequester := &mocks.ShipmentReweighRequester{}

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-reweigh", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		planner := &routemocks.Planner{}
		planner.On("TransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveRouter := moverouter.NewMoveRouter()
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)

		handler := RequestShipmentReweighHandler{
			handlerConfig,
			reweighRequester,
			mtoshipment.NewShipmentSITStatus(),
			updater,
		}
		params := shipmentops.RequestShipmentReweighParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		reweighRequester.On("RequestShipmentReweigh", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, models.ReweighRequesterTOO).Return(nil, apperror.ConflictError{})

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RequestShipmentReweighConflict{}, response)
	})

	suite.Run("Returns 422 when reweighRequester returns validation errors", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		reweighRequester := &mocks.ShipmentReweighRequester{}

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-reweigh", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		planner := &routemocks.Planner{}
		planner.On("TransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveRouter := moverouter.NewMoveRouter()
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)

		handler := RequestShipmentReweighHandler{
			handlerConfig,
			reweighRequester,
			mtoshipment.NewShipmentSITStatus(),
			updater,
		}
		params := shipmentops.RequestShipmentReweighParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}
		reweighRequester.On("RequestShipmentReweigh", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, models.ReweighRequesterTOO).Return(nil, apperror.InvalidInputError{ValidationErrors: &validate.Errors{}})

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RequestShipmentReweighUnprocessableEntity{}, response)
	})

	suite.Run("Returns 500 when reweighRequester returns unexpected error", func() {
		shipment := testdatagen.MakeStubbedShipment(suite.DB())
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		reweighRequester := &mocks.ShipmentReweighRequester{}

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-reweigh", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		planner := &routemocks.Planner{}
		planner.On("TransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveRouter := moverouter.NewMoveRouter()
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)

		handler := RequestShipmentReweighHandler{
			handlerConfig,
			reweighRequester,
			mtoshipment.NewShipmentSITStatus(),
			updater,
		}
		params := shipmentops.RequestShipmentReweighParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		reweighRequester.On("RequestShipmentReweigh", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, models.ReweighRequesterTOO).Return(nil, errors.New("UnexpectedError"))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RequestShipmentReweighInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestApproveSITExtensionHandler() {
	suite.Run("Returns 200 and updates SIT days allowance when validations pass", func() {
		sitDaysAllowance := 20
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Entitlement: models.Entitlement{
				StorageInTransit: &sitDaysAllowance,
			},
		})
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				SITDaysAllowance: &sitDaysAllowance,
			},
			Move: move,
		})

		year, month, day := time.Now().Date()
		lastMonthEntry := time.Date(year, month, day-37, 0, 0, 0, 0, time.UTC)
		lastMonthDeparture := time.Date(year, month, day-30, 0, 0, 0, 0, time.UTC)
		testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOServiceItem: models.MTOServiceItem{
				SITEntryDate:     &lastMonthEntry,
				SITDepartureDate: &lastMonthDeparture,
				Status:           models.MTOServiceItemStatusApproved,
			},
			Move:        move,
			MTOShipment: mtoShipment,
			ReService: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		})
		sitExtension := testdatagen.MakePendingSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		moveRouter := moverouter.NewMoveRouter()
		sitExtensionApprover := mtoshipment.NewSITExtensionApprover(moveRouter)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/shipments/%s/sit-extension/%s/approve", mtoShipment.ID.String(), sitExtension.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := ApproveSITExtensionHandler{
			handlerConfig,
			sitExtensionApprover,
			mtoshipment.NewShipmentSITStatus(),
		}
		approvedDays := int64(10)
		officeRemarks := "new office remarks"
		approveParams := shipmentops.ApproveSITExtensionParams{
			HTTPRequest: req,
			IfMatch:     eTag,
			Body: &ghcmessages.ApproveSITExtension{
				ApprovedDays:  &approvedDays,
				OfficeRemarks: &officeRemarks,
			},
			ShipmentID:     *handlers.FmtUUID(mtoShipment.ID),
			SitExtensionID: *handlers.FmtUUID(sitExtension.ID),
		}
		response := handler.Handle(approveParams)
		okResponse := response.(*shipmentops.ApproveSITExtensionOK)
		payload := okResponse.Payload
		suite.IsType(&shipmentops.ApproveSITExtensionOK{}, response)
		suite.Equal(int64(30), *payload.SitDaysAllowance)
		suite.Equal("APPROVED", payload.SitExtensions[0].Status)
		suite.Require().NotNil(payload.SitExtensions[0].OfficeRemarks)
		suite.Equal(officeRemarks, *payload.SitExtensions[0].OfficeRemarks)
	})
}

func (suite *HandlerSuite) TestDenySITExtensionHandler() {
	suite.Run("Returns 200 when validations pass", func() {
		sitDaysAllowance := 20
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				SITDaysAllowance: &sitDaysAllowance,
			},
			Move: move,
		})
		sitExtension := testdatagen.MakePendingSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		moveRouter := moverouter.NewMoveRouter()
		sitExtensionDenier := mtoshipment.NewSITExtensionDenier(moveRouter)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/shipments/%s/sit-extension/%s/deny", mtoShipment.ID.String(), sitExtension.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := DenySITExtensionHandler{
			handlerConfig,
			sitExtensionDenier,
			mtoshipment.NewShipmentSITStatus(),
		}
		officeRemarks := "new office remarks on denial of extension"
		denyParams := shipmentops.DenySITExtensionParams{
			HTTPRequest: req,
			IfMatch:     eTag,
			Body: &ghcmessages.DenySITExtension{
				OfficeRemarks: &officeRemarks,
			},
			ShipmentID:     *handlers.FmtUUID(mtoShipment.ID),
			SitExtensionID: *handlers.FmtUUID(sitExtension.ID),
		}
		response := handler.Handle(denyParams)
		okResponse := response.(*shipmentops.DenySITExtensionOK)
		payload := okResponse.Payload
		suite.IsType(&shipmentops.DenySITExtensionOK{}, response)
		suite.Equal("DENIED", payload.SitExtensions[0].Status)
	})
}

func (suite *HandlerSuite) CreateSITExtensionAsTOO() {
	suite.Run("Returns 200, creates new SIT extension, and updates SIT days allowance on shipment without an allowance when validations pass", func() {
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{},
		})

		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		sitExtensionCreatorAsTOO := mtoshipment.NewCreateSITExtensionAsTOO()
		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/sit-extension/", mtoShipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := CreateSITExtensionAsTOOHandler{
			handlerConfig,
			sitExtensionCreatorAsTOO,
			mtoshipment.NewShipmentSITStatus(),
		}
		approvedDays := int64(10)
		officeRemarks := "new office remarks"
		requestReason := "OTHER"
		createParams := shipmentops.CreateSITExtensionAsTOOParams{
			HTTPRequest: req,
			IfMatch:     eTag,
			Body: &ghcmessages.CreateSITExtensionAsTOO{
				ApprovedDays:  &approvedDays,
				OfficeRemarks: &officeRemarks,
				RequestReason: &requestReason,
			},
			ShipmentID: *handlers.FmtUUID(mtoShipment.ID),
		}
		suite.NoError(createParams.Body.Validate(strfmt.Default))

		response := handler.Handle(createParams)
		okResponse := response.(*shipmentops.CreateSITExtensionAsTOOOK)
		payload := okResponse.Payload
		suite.IsType(&shipmentops.CreateSITExtensionAsTOOOK{}, response)
		suite.Equal(int64(10), *payload.SitDaysAllowance)
		suite.Equal("APPROVED", payload.SitExtensions[0].Status)
		suite.Require().NotNil(payload.SitExtensions[0].OfficeRemarks)
		suite.Equal(officeRemarks, *payload.SitExtensions[0].OfficeRemarks)
	})

	suite.Run("Returns 200, creates new SIT extension, and updates SIT days allowance on shipment that already has an allowance when validations pass", func() {
		sitDaysAllowance := 20
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				SITDaysAllowance: &sitDaysAllowance,
			},
		})

		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		sitExtensionCreatorAsTOO := mtoshipment.NewCreateSITExtensionAsTOO()
		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/sit-extension/", mtoShipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		handler := CreateSITExtensionAsTOOHandler{
			handlerConfig,
			sitExtensionCreatorAsTOO,
			mtoshipment.NewShipmentSITStatus(),
		}
		approvedDays := int64(10)
		officeRemarks := "new office remarks"
		requestReason := "OTHER"
		createParams := shipmentops.CreateSITExtensionAsTOOParams{
			HTTPRequest: req,
			IfMatch:     eTag,
			Body: &ghcmessages.CreateSITExtensionAsTOO{
				ApprovedDays:  &approvedDays,
				OfficeRemarks: &officeRemarks,
				RequestReason: &requestReason,
			},
			ShipmentID: *handlers.FmtUUID(mtoShipment.ID),
		}
		suite.NoError(createParams.Body.Validate(strfmt.Default))

		response := handler.Handle(createParams)
		okResponse := response.(*shipmentops.CreateSITExtensionAsTOOOK)
		payload := okResponse.Payload
		suite.IsType(&shipmentops.CreateSITExtensionAsTOOOK{}, response)
		suite.Equal(int64(30), *payload.SitDaysAllowance)
		suite.Equal("APPROVED", payload.SitExtensions[0].Status)
		suite.Require().NotNil(payload.SitExtensions[0].OfficeRemarks)
		suite.Equal(officeRemarks, *payload.SitExtensions[0].OfficeRemarks)
	})
}

type createMTOShipmentSubtestData struct {
	builder *query.Builder
	params  mtoshipmentops.CreateMTOShipmentParams
	traceID uuid.UUID
}

func (suite *HandlerSuite) makeCreateMTOShipmentSubtestData() (subtestData *createMTOShipmentSubtestData) {
	subtestData = &createMTOShipmentSubtestData{}

	mto := testdatagen.MakeAvailableMove(suite.DB())
	pickupAddress := testdatagen.MakeDefaultAddress(suite.DB())
	destinationAddress := testdatagen.MakeDefaultAddress(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move:        mto,
		MTOShipment: models.MTOShipment{},
	})

	mtoShipment.MoveTaskOrderID = mto.ID

	subtestData.builder = query.NewQueryBuilder()

	req := httptest.NewRequest("POST", "/mto-shipments", nil)

	// Set the traceID so we can use it to find the webhook notification
	traceID, err := uuid.NewV4()
	suite.FatalNoError(err, "Error creating a new trace ID.")
	req = req.WithContext(trace.NewContext(req.Context(), traceID))

	subtestData.traceID = traceID
	shipmentType := ghcmessages.MTOShipmentType(mtoShipment.ShipmentType)
	subtestData.params = mtoshipmentops.CreateMTOShipmentParams{
		HTTPRequest: req,
		Body: &ghcmessages.CreateMTOShipment{
			MoveTaskOrderID:     handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
			Agents:              nil,
			CustomerRemarks:     handlers.FmtString("customer remark"),
			CounselorRemarks:    handlers.FmtString("counselor remark"),
			RequestedPickupDate: handlers.FmtDatePtr(mtoShipment.RequestedPickupDate),
			ShipmentType:        &shipmentType,
		},
	}
	subtestData.params.Body.DestinationAddress.Address = ghcmessages.Address{
		City:           &destinationAddress.City,
		Country:        destinationAddress.Country,
		PostalCode:     &destinationAddress.PostalCode,
		State:          &destinationAddress.State,
		StreetAddress1: &destinationAddress.StreetAddress1,
		StreetAddress2: destinationAddress.StreetAddress2,
		StreetAddress3: destinationAddress.StreetAddress3,
	}
	subtestData.params.Body.PickupAddress.Address = ghcmessages.Address{
		City:           &pickupAddress.City,
		Country:        pickupAddress.Country,
		PostalCode:     &pickupAddress.PostalCode,
		State:          &pickupAddress.State,
		StreetAddress1: &pickupAddress.StreetAddress1,
		StreetAddress2: pickupAddress.StreetAddress2,
		StreetAddress3: pickupAddress.StreetAddress3,
	}

	return subtestData
}

func (suite *HandlerSuite) TestCreateMTOShipmentHandler() {
	moveRouter := moverouter.NewMoveRouter()

	suite.Run("Successful POST - Integration Test", func() {
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		subtestData := suite.makeCreateMTOShipmentSubtestData()
		builder := subtestData.builder
		params := subtestData.params

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
		ppmEstimator := mocks.PPMEstimator{}
		ppmCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmCreator)
		sitStatus := mtoshipment.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitStatus,
		}
		response := handler.Handle(params)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		createMTOShipmentPayload := okResponse.Payload
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)

		suite.Require().Equal(ghcmessages.MTOShipmentStatusSUBMITTED, createMTOShipmentPayload.Status, "MTO Shipment should have been submitted")
		suite.Require().Equal(createMTOShipmentPayload.ShipmentType, ghcmessages.MTOShipmentTypeHHG, "MTO Shipment should be an HHG")
		suite.Equal(int64(models.DefaultServiceMemberSITDaysAllowance), *createMTOShipmentPayload.SitDaysAllowance)
		suite.Equal(string("customer remark"), *createMTOShipmentPayload.CustomerRemarks)
		suite.Equal(string("counselor remark"), *createMTOShipmentPayload.CounselorRemarks)
	})

	suite.Run("POST failure - 500", func() {
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		subtestData := suite.makeCreateMTOShipmentSubtestData()
		params := subtestData.params

		shipmentCreator := mocks.ShipmentCreator{}
		sitStatus := mtoshipment.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			&shipmentCreator,
			sitStatus,
		}

		err := errors.New("ServerError")

		shipmentCreator.On("CreateShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.MTOShipment"),
		).Return(nil, err)

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentInternalServerError{}, response)
	})

	suite.Run("POST failure - 422 -- Bad agent IDs set on shipment", func() {
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		subtestData := suite.makeCreateMTOShipmentSubtestData()
		builder := subtestData.builder
		params := subtestData.params

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
		ppmEstimator := mocks.PPMEstimator{}
		ppmCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmCreator)
		sitStatus := mtoshipment.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitStatus,
		}

		badID := params.Body.MoveTaskOrderID
		agent := &ghcmessages.MTOAgent{
			ID:            *badID,
			MtoShipmentID: *badID,
			FirstName:     handlers.FmtString("Mary"),
		}

		paramsBadIDs := params
		paramsBadIDs.Body.Agents = ghcmessages.MTOAgents{agent}

		response := handler.Handle(paramsBadIDs)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)
		suite.NotEmpty(typedResponse.Payload.InvalidFields)
	})

	suite.Run("POST failure - 422 - invalid input, missing pickup address", func() {
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		subtestData := suite.makeCreateMTOShipmentSubtestData()
		builder := subtestData.builder
		params := subtestData.params

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
		ppmEstimator := mocks.PPMEstimator{}
		ppmCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmCreator)
		sitStatus := mtoshipment.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitStatus,
		}

		badParams := params
		badParams.Body.PickupAddress.Address.StreetAddress1 = nil

		suite.NoError(badParams.Body.Validate(strfmt.Default))

		response := handler.Handle(badParams)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)
		// CreateMTOShipment is returning services.NewInvalidInputError without any validation errors
		// so InvalidFields won't be added to the payload.
		suite.Empty(typedResponse.Payload.InvalidFields)
	})

	suite.Run("POST failure - 404 -- not found", func() {
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		subtestData := suite.makeCreateMTOShipmentSubtestData()
		builder := subtestData.builder
		params := subtestData.params

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
		ppmEstimator := mocks.PPMEstimator{}
		ppmCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmCreator)
		sitStatus := mtoshipment.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitStatus,
		}

		uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
		badParams := params
		badParams.Body.MoveTaskOrderID = handlers.FmtUUID(uuid.FromStringOrNil(uuidString))

		response := handler.Handle(badParams)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
	})

	suite.Run("POST failure - 400 -- nil body", func() {
		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())

		subtestData := suite.makeCreateMTOShipmentSubtestData()
		builder := subtestData.builder

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
		ppmEstimator := mocks.PPMEstimator{}
		ppmCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmCreator)
		sitStatus := mtoshipment.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitStatus,
		}

		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		paramsNilBody := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
		}
		response := handler.Handle(paramsNilBody)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentBadRequest{}, response)
	})
}

func (suite *HandlerSuite) TestCreateMTOShipmentHandlerUsingPPM() {
	suite.Run("Successful POST - Integration Test (PPM, all fields)", func() {
		// Make a move along with an attached minimal shipment. Shouldn't matter what's in them.
		move := testdatagen.MakeDefaultMove(suite.DB())
		hhgShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: move,
		})

		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moverouter.NewMoveRouter())
		ppmEstimator := mocks.PPMEstimator{}
		ppmCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmCreator)
		sitStatus := mtoshipment.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitStatus,
		}

		shipmentType := ghcmessages.MTOShipmentTypePPM
		expectedDepartureDate := hhgShipment.RequestedPickupDate
		pickupPostalCode := "30907"
		secondaryPickupPostalCode := "30809"
		destinationPostalCode := "36106"
		secondaryDestinationPostalCode := "36101"
		sitExpected := true
		sitLocation := ghcmessages.SITLocationTypeDESTINATION
		sitEstimatedWeight := unit.Pound(1700)
		sitEstimatedEntryDate := expectedDepartureDate.AddDate(0, 0, 5)
		sitEstimatedDepartureDate := sitEstimatedEntryDate.AddDate(0, 0, 20)
		estimatedWeight := unit.Pound(3000)
		hasProGear := true
		proGearWeight := unit.Pound(300)
		spouseProGearWeight := unit.Pound(200)
		estimatedIncentive := 654321
		req := httptest.NewRequest("POST", "/mto-shipments", nil)
		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &ghcmessages.CreateMTOShipment{
				MoveTaskOrderID: handlers.FmtUUID(move.ID),
				ShipmentType:    &shipmentType,
				PpmShipment: &ghcmessages.CreatePPMShipment{
					ExpectedDepartureDate:          handlers.FmtDatePtr(expectedDepartureDate),
					PickupPostalCode:               &pickupPostalCode,
					SecondaryPickupPostalCode:      &secondaryPickupPostalCode,
					DestinationPostalCode:          &destinationPostalCode,
					SecondaryDestinationPostalCode: &secondaryDestinationPostalCode,
					SitExpected:                    &sitExpected,
					SitLocation:                    &sitLocation,
					SitEstimatedWeight:             handlers.FmtPoundPtr(&sitEstimatedWeight),
					SitEstimatedEntryDate:          handlers.FmtDate(sitEstimatedEntryDate),
					SitEstimatedDepartureDate:      handlers.FmtDate(sitEstimatedDepartureDate),
					EstimatedWeight:                handlers.FmtPoundPtr(&estimatedWeight),
					HasProGear:                     &hasProGear,
					ProGearWeight:                  handlers.FmtPoundPtr(&proGearWeight),
					SpouseProGearWeight:            handlers.FmtPoundPtr(&spouseProGearWeight),
				},
			},
		}

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), nil).Once()

		response := handler.Handle(params)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)

		// Check MTOShipment fields.
		payload := okResponse.Payload
		suite.NotZero(payload.ID)
		suite.NotEqual(uuid.Nil.String(), payload.ID.String())
		suite.Equal(move.ID.String(), payload.MoveTaskOrderID.String())
		suite.Equal(ghcmessages.MTOShipmentTypePPM, payload.ShipmentType)
		suite.Equal(ghcmessages.MTOShipmentStatusSUBMITTED, payload.Status)
		suite.NotZero(payload.CreatedAt)
		suite.NotZero(payload.UpdatedAt)

		// Check PPMShipment fields.
		ppmPayload := payload.PpmShipment
		if suite.NotNil(ppmPayload) {
			suite.NotZero(ppmPayload.ID)
			suite.NotEqual(uuid.Nil.String(), ppmPayload.ID.String())
			suite.EqualDatePtr(expectedDepartureDate, ppmPayload.ExpectedDepartureDate)
			suite.Equal(pickupPostalCode, *ppmPayload.PickupPostalCode)
			suite.Equal(&secondaryPickupPostalCode, ppmPayload.SecondaryPickupPostalCode)
			suite.Equal(destinationPostalCode, *ppmPayload.DestinationPostalCode)
			suite.Equal(&secondaryDestinationPostalCode, ppmPayload.SecondaryDestinationPostalCode)
			suite.Equal(sitExpected, *ppmPayload.SitExpected)
			suite.Equal(&sitLocation, ppmPayload.SitLocation)
			suite.Equal(handlers.FmtPoundPtr(&sitEstimatedWeight), ppmPayload.SitEstimatedWeight)
			suite.Equal(handlers.FmtDate(sitEstimatedEntryDate), ppmPayload.SitEstimatedEntryDate)
			suite.Equal(handlers.FmtDate(sitEstimatedDepartureDate), ppmPayload.SitEstimatedDepartureDate)
			suite.Equal(handlers.FmtPoundPtr(&estimatedWeight), ppmPayload.EstimatedWeight)
			suite.Equal(&hasProGear, ppmPayload.HasProGear)
			suite.Equal(handlers.FmtPoundPtr(&proGearWeight), ppmPayload.ProGearWeight)
			suite.Equal(handlers.FmtPoundPtr(&spouseProGearWeight), ppmPayload.SpouseProGearWeight)
			suite.Equal(ghcmessages.PPMShipmentStatusSUBMITTED, ppmPayload.Status)
			suite.Equal(int64(estimatedIncentive), *ppmPayload.EstimatedIncentive)
			suite.NotZero(ppmPayload.CreatedAt)
			suite.NotZero(ppmPayload.UpdatedAt)
		}
	})

	suite.Run("Successful POST - Integration Test (PPM, minimal fields)", func() {
		// Make a move along with an attached minimal shipment. Shouldn't matter what's in them.
		move := testdatagen.MakeDefaultMove(suite.DB())
		hhgShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: move,
		})

		handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moverouter.NewMoveRouter())
		ppmEstimator := mocks.PPMEstimator{}
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmshipment.NewPPMShipmentCreator(&ppmEstimator))
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			mtoshipment.NewShipmentSITStatus(),
		}

		shipmentType := ghcmessages.MTOShipmentTypePPM
		expectedDepartureDate := hhgShipment.RequestedPickupDate
		pickupPostalCode := "29212"
		destinationPostalCode := "78234"
		sitExpected := false
		estimatedWeight := unit.Pound(2450)
		hasProGear := false
		estimatedIncentive := 123456
		req := httptest.NewRequest("POST", "/mto-shipments", nil)
		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &ghcmessages.CreateMTOShipment{
				MoveTaskOrderID: handlers.FmtUUID(move.ID),
				ShipmentType:    &shipmentType,
				PpmShipment: &ghcmessages.CreatePPMShipment{
					ExpectedDepartureDate: handlers.FmtDatePtr(expectedDepartureDate),
					PickupPostalCode:      &pickupPostalCode,
					DestinationPostalCode: &destinationPostalCode,
					SitExpected:           &sitExpected,
					EstimatedWeight:       handlers.FmtPoundPtr(&estimatedWeight),
					HasProGear:            &hasProGear,
				},
			},
		}

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), nil).Once()

		response := handler.Handle(params)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)

		// Check MTOShipment fields.
		payload := okResponse.Payload
		suite.NotZero(payload.ID)
		suite.NotEqual(uuid.Nil.String(), payload.ID.String())
		suite.Equal(move.ID.String(), payload.MoveTaskOrderID.String())
		suite.Equal(ghcmessages.MTOShipmentTypePPM, payload.ShipmentType)
		suite.Equal(ghcmessages.MTOShipmentStatusSUBMITTED, payload.Status)
		suite.NotZero(payload.CreatedAt)
		suite.NotZero(payload.UpdatedAt)

		// Check PPMShipment fields.
		ppmPayload := payload.PpmShipment
		if suite.NotNil(ppmPayload) {
			suite.NotZero(ppmPayload.ID)
			suite.NotEqual(uuid.Nil.String(), ppmPayload.ID.String())
			suite.EqualDatePtr(expectedDepartureDate, ppmPayload.ExpectedDepartureDate)
			suite.Equal(pickupPostalCode, *ppmPayload.PickupPostalCode)
			suite.Equal(destinationPostalCode, *ppmPayload.DestinationPostalCode)
			suite.Equal(sitExpected, *ppmPayload.SitExpected)
			suite.Equal(handlers.FmtPoundPtr(&estimatedWeight), ppmPayload.EstimatedWeight)
			suite.Equal(&hasProGear, ppmPayload.HasProGear)
			suite.Equal(ghcmessages.PPMShipmentStatusSUBMITTED, ppmPayload.Status)
			suite.Equal(int64(estimatedIncentive), *ppmPayload.EstimatedIncentive)
			suite.NotZero(ppmPayload.CreatedAt)
			suite.NotZero(ppmPayload.UpdatedAt)
		}
	})
}

func (suite *HandlerSuite) getUpdateShipmentParams(originalShipment models.MTOShipment) mtoshipmentops.UpdateMTOShipmentParams {
	servicesCounselor := testdatagen.MakeDefaultOfficeUser(suite.DB())
	servicesCounselor.User.Roles = append(servicesCounselor.User.Roles, roles.Role{
		RoleType: roles.RoleTypeServicesCounselor,
	})
	pickupAddress := testdatagen.MakeDefaultAddress(suite.DB())
	pickupAddress.StreetAddress1 = "123 Fake Test St NW"
	destinationAddress := testdatagen.MakeDefaultAddress(suite.DB())
	destinationAddress.StreetAddress1 = "54321 Test Fake Rd SE"
	customerRemarks := "help"
	counselorRemarks := "counselor approved"
	billableWeightCap := int64(8000)
	billableWeightJustification := "Unable to perform reweigh because shipment was already unloaded."
	mtoAgent := testdatagen.MakeDefaultMTOAgent(suite.DB())
	agents := ghcmessages.MTOAgents{&ghcmessages.MTOAgent{
		FirstName: mtoAgent.FirstName,
		LastName:  mtoAgent.LastName,
		Email:     mtoAgent.Email,
		Phone:     mtoAgent.Phone,
		AgentType: string(mtoAgent.MTOAgentType),
	}}

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/mto_shipments/%s", originalShipment.MoveTaskOrderID.String(), originalShipment.ID.String()), nil)
	req = suite.AuthenticateOfficeRequest(req, servicesCounselor)

	eTag := etag.GenerateEtag(originalShipment.UpdatedAt)

	now := strfmt.Date(time.Now())
	payload := ghcmessages.UpdateShipment{
		BillableWeightJustification: &billableWeightJustification,
		BillableWeightCap:           &billableWeightCap,
		RequestedPickupDate:         &now,
		RequestedDeliveryDate:       &now,
		ShipmentType:                ghcmessages.MTOShipmentTypeHHG,
		CustomerRemarks:             &customerRemarks,
		CounselorRemarks:            &counselorRemarks,
		Agents:                      agents,
		TacType:                     nullable.NewString("NTS"),
		SacType:                     nullable.NewString(""),
	}
	payload.DestinationAddress.Address = ghcmessages.Address{
		City:           &destinationAddress.City,
		Country:        destinationAddress.Country,
		PostalCode:     &destinationAddress.PostalCode,
		State:          &destinationAddress.State,
		StreetAddress1: &destinationAddress.StreetAddress1,
		StreetAddress2: destinationAddress.StreetAddress2,
		StreetAddress3: destinationAddress.StreetAddress3,
	}
	payload.PickupAddress.Address = ghcmessages.Address{
		City:           &pickupAddress.City,
		Country:        pickupAddress.Country,
		PostalCode:     &pickupAddress.PostalCode,
		State:          &pickupAddress.State,
		StreetAddress1: &pickupAddress.StreetAddress1,
		StreetAddress2: pickupAddress.StreetAddress2,
		StreetAddress3: pickupAddress.StreetAddress3,
	}

	params := mtoshipmentops.UpdateMTOShipmentParams{
		HTTPRequest: req,
		ShipmentID:  *handlers.FmtUUID(originalShipment.ID),
		Body:        &payload,
		IfMatch:     eTag,
	}

	return params
}

func (suite *HandlerSuite) TestUpdateShipmentHandler() {
	planner := &routemocks.Planner{}
	planner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	moveRouter := moverouter.NewMoveRouter()
	moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

	// Get shipment payment request recalculator service
	creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

	suite.Run("Successful PATCH - Integration Test", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)
		handler := UpdateShipmentHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			fetcher,
			updater,
			mtoshipment.NewShipmentSITStatus(),
		}

		hhgLOAType := models.LOATypeHHG
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:             models.MTOShipmentStatusSubmitted,
				UsesExternalVendor: true,
				TACType:            &hhgLOAType,
			},
		})
		params := suite.getUpdateShipmentParams(oldShipment)

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		updatedShipment := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload
		suite.Equal(oldShipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(params.Body.BillableWeightCap, updatedShipment.BillableWeightCap)
		suite.Equal(params.Body.BillableWeightJustification, updatedShipment.BillableWeightJustification)
		suite.Equal(params.Body.CustomerRemarks, updatedShipment.CustomerRemarks)
		suite.Equal(params.Body.CounselorRemarks, updatedShipment.CounselorRemarks)
		suite.Equal(params.Body.PickupAddress.StreetAddress1, updatedShipment.PickupAddress.StreetAddress1)
		suite.Equal(params.Body.DestinationAddress.StreetAddress1, updatedShipment.DestinationAddress.StreetAddress1)
		suite.Equal(params.Body.RequestedPickupDate.String(), updatedShipment.RequestedPickupDate.String())
		suite.Equal(params.Body.Agents[0].FirstName, updatedShipment.MtoAgents[0].FirstName)
		suite.Equal(params.Body.Agents[0].LastName, updatedShipment.MtoAgents[0].LastName)
		suite.Equal(params.Body.Agents[0].Email, updatedShipment.MtoAgents[0].Email)
		suite.Equal(params.Body.Agents[0].Phone, updatedShipment.MtoAgents[0].Phone)
		suite.Equal(params.Body.Agents[0].AgentType, updatedShipment.MtoAgents[0].AgentType)
		suite.Equal(oldShipment.ID.String(), string(updatedShipment.MtoAgents[0].MtoShipmentID))
		suite.NotEmpty(updatedShipment.MtoAgents[0].ID)
		suite.Equal(params.Body.RequestedDeliveryDate.String(), updatedShipment.RequestedDeliveryDate.String())
		suite.Equal(oldShipment.UsesExternalVendor, updatedShipment.UsesExternalVendor)
		suite.Equal(*params.Body.TacType.Value, string(*updatedShipment.TacType))
		suite.Nil(updatedShipment.SacType)
	})

	suite.Run("PATCH failure - 400 -- nil body", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)
		handler := UpdateShipmentHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			fetcher,
			updater,
			mtoshipment.NewShipmentSITStatus(),
		}

		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		})
		params := suite.getUpdateShipmentParams(oldShipment)
		params.Body = nil

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
	})

	suite.Run("PATCH failure - 404 -- not found", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)
		handler := UpdateShipmentHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			fetcher,
			updater,
			mtoshipment.NewShipmentSITStatus(),
		}

		uuidString := handlers.FmtUUID(uuid.FromStringOrNil("d874d002-5582-4a91-97d3-786e8f66c763"))
		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		})
		params := suite.getUpdateShipmentParams(oldShipment)
		params.ShipmentID = *uuidString

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)
		handler := UpdateShipmentHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			fetcher,
			updater,
			mtoshipment.NewShipmentSITStatus(),
		}

		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		})
		params := suite.getUpdateShipmentParams(oldShipment)
		params.IfMatch = "intentionally-bad-if-match-header-value"

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentPreconditionFailed{}, response)
	})

	suite.Run("PATCH failure - 412 -- shipment shouldn't be updatable", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator)
		handler := UpdateShipmentHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			fetcher,
			updater,
			mtoshipment.NewShipmentSITStatus(),
		}

		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusDraft,
			},
		})

		params := suite.getUpdateShipmentParams(oldShipment)

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentPreconditionFailed{}, response)
	})

	suite.Run("PATCH failure - 500", func() {
		builder := query.NewQueryBuilder()
		mockUpdater := mocks.MTOShipmentUpdater{}
		fetcher := fetch.NewFetcher(builder)
		handler := UpdateShipmentHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			fetcher,
			&mockUpdater,
			mtoshipment.NewShipmentSITStatus(),
		}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateMTOShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, err)
		mockUpdater.On("RetrieveMTOShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, err)
		mockUpdater.On("CheckIfMTOShipmentCanBeUpdated",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil, err)

		oldShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		})
		params := suite.getUpdateShipmentParams(oldShipment)

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentInternalServerError{}, response)
	})

}
