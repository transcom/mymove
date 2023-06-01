package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	sitaddressupdate "github.com/transcom/mymove/pkg/services/sit_address_update"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/trace"
)

func (suite *HandlerSuite) TestListMTOServiceItemHandler() {
	reServiceID, _ := uuid.NewV4()
	serviceItemID, _ := uuid.NewV4()
	mtoShipmentID, _ := uuid.NewV4()
	var mtoID uuid.UUID

	setupTestData := func() (models.User, models.MTOServiceItems) {
		mto := factory.BuildMove(suite.DB(), nil, nil)
		mtoID = mto.ID
		reService := factory.BuildReService(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					ID:   reServiceID,
					Code: "TEST10000",
				},
			},
		}, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{ID: mtoShipmentID},
			},
		}, nil)
		requestUser := factory.BuildUser(nil, nil, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					ID: serviceItemID,
				},
			},
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)

		year, month, day := time.Now().Date()
		aWeekAgo := time.Date(year, month, day-7, 0, 0, 0, 0, time.UTC)
		departureDate := aWeekAgo.Add(time.Hour * 24 * 30)
		originSit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &aWeekAgo,
					SITDepartureDate: &departureDate,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOPSIT,
				},
			},
		}, nil)

		sitAddressUpdate := factory.BuildSITAddressUpdate(suite.DB(), []factory.Customization{{Model: originSit,
			LinkOnly: true}}, nil)
		originSit.SITAddressUpdates = []models.SITAddressUpdate{sitAddressUpdate}

		customerContact := testdatagen.MakeMTOServiceItemCustomerContact(suite.DB(), testdatagen.Assertions{})
		destinationSit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					CustomerContacts: models.MTOServiceItemCustomerContacts{customerContact},
				},
			},
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDFSIT,
					Name: "Destination 1st Day SIT",
				},
			},
		}, nil)

		serviceItems := models.MTOServiceItems{serviceItem, originSit, destinationSit}

		return requestUser, serviceItems
	}

	suite.Run("Successful list fetch - Integration Test", func() {
		requestUser, serviceItems := setupTestData()
		req := httptest.NewRequest("GET", fmt.Sprintf("/move_task_orders/%s/mto_service_items", mtoID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := mtoserviceitemop.ListMTOServiceItemsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: *handlers.FmtUUID(serviceItems[0].MoveTaskOrderID),
		}

		queryBuilder := query.NewQueryBuilder()
		listFetcher := fetch.NewListFetcher(queryBuilder)
		fetcher := fetch.NewFetcher(queryBuilder)
		handler := ListMTOServiceItemsHandler{
			suite.HandlerConfig(),
			listFetcher,
			fetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.ListMTOServiceItemsOK{}, response)
		okResponse := response.(*mtoserviceitemop.ListMTOServiceItemsOK)

		// Validate outgoing payload
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		suite.Len(okResponse.Payload, 3)
		suite.Equal(serviceItems[0].ID.String(), okResponse.Payload[0].ID.String())
		suite.Equal(serviceItems[1].ID.String(), okResponse.Payload[1].ID.String())

		// Validate that SITAddressUpdates are included in payload
		suite.Len(okResponse.Payload[1].SitAddressUpdates, 1)
		suite.Equal(serviceItems[1].SITAddressUpdates[0].ID.String(), okResponse.Payload[1].SitAddressUpdates[0].ID.String())

		// Validate that the Customer Contacts were included in the payload
		suite.Len(okResponse.Payload[2].CustomerContacts, 1)
		suite.Equal(serviceItems[2].CustomerContacts[0].TimeMilitary, okResponse.Payload[2].CustomerContacts[0].TimeMilitary)
	})

	suite.Run("Failure list fetch - Internal Server Error", func() {
		requestUser, serviceItems := setupTestData()
		req := httptest.NewRequest("GET", fmt.Sprintf("/move_task_orders/%s/mto_service_items", mtoID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := mtoserviceitemop.ListMTOServiceItemsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: *handlers.FmtUUID(serviceItems[0].MoveTaskOrderID),
		}
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOServiceItemsHandler{
			suite.HandlerConfig(),
			&mockListFetcher,
			&mockFetcher,
		}

		internalServerErr := errors.New("ServerError")

		mockFetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil)

		mockListFetcher.On("FetchRecordList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(internalServerErr)

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.ListMTOServiceItemsInternalServerError{}, response)
		payload := response.(*mtoserviceitemop.ListMTOServiceItemsInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Failure list fetch - 404 Not Found - Move Task Order ID", func() {
		requestUser, serviceItems := setupTestData()
		req := httptest.NewRequest("GET", fmt.Sprintf("/move_task_orders/%s/mto_service_items", mtoID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := mtoserviceitemop.ListMTOServiceItemsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: *handlers.FmtUUID(serviceItems[0].MoveTaskOrderID),
		}

		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOServiceItemsHandler{
			suite.HandlerConfig(),
			&mockListFetcher,
			&mockFetcher,
		}

		notfound := errors.New("Not found error")

		mockFetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(notfound)

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.ListMTOServiceItemsNotFound{}, response)
		payload := response.(*mtoserviceitemop.ListMTOServiceItemsNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) createServiceItem() (models.MTOServiceItem, models.Move) {
	move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
	serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	return serviceItem, move
}

func (suite *HandlerSuite) TestUpdateMTOServiceItemStatusHandler() {
	moveTaskOrderID, _ := uuid.NewV4()
	serviceItemID, _ := uuid.NewV4()
	var requestUser models.User

	setupTestData := func() mtoserviceitemop.UpdateMTOServiceItemStatusParams {
		requestUser = factory.BuildUser(nil, nil, nil)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/mto_service_items/%s/status",
			moveTaskOrderID, serviceItemID), nil)

		req = suite.AuthenticateUserRequest(req, requestUser)
		params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
			HTTPRequest:      req,
			IfMatch:          etag.GenerateEtag(time.Now()),
			Body:             &ghcmessages.PatchMTOServiceItemStatusPayload{Status: "APPROVED"},
			MoveTaskOrderID:  moveTaskOrderID.String(),
			MtoServiceItemID: serviceItemID.String(),
		}
		return params
	}

	// With this first set of tests we'll use mocked service object responses so that we can make sure the handler
	// is returning the right HTTP code given a set of circumstances.
	suite.Run("404 - not found response", func() {
		params := setupTestData()
		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		fetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(errors.New("Not found error")).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusNotFound{}, response)
		payload := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("200 - success response", func() {
		params := setupTestData()
		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		fetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()

		serviceItemStatusUpdater.On("ApproveOrRejectServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&models.MTOServiceItem{ID: serviceItemID}, nil).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)
		payload := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("412 - precondition failed response", func() {
		params := setupTestData()

		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		fetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()

		serviceItemStatusUpdater.On("ApproveOrRejectServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, apperror.NewPreconditionFailedError(serviceItemID, errors.New("oh no"))).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusPreconditionFailed{}, response)
		payload := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("500 - internal server error response", func() {
		params := setupTestData()

		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		fetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()

		serviceItemStatusUpdater.On("ApproveOrRejectServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, errors.New("oh no")).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusInternalServerError{}, response)
		payload := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("422 - unprocessable entity response", func() {
		params := setupTestData()

		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		params.MtoServiceItemID = ""
		fetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusUnprocessableEntity{}, response)
		payload := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	// With this we'll do a happy path integration test to ensure that the use of the service object
	// by the handler is working as expected.
	suite.Run("Successful rejected status update - Integration test", func() {

		queryBuilder := query.NewQueryBuilder()
		mtoServiceItem, move := suite.createServiceItem()
		requestUser := factory.BuildUser(nil, nil, nil)

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/mto_service_items/%s/status",
			moveTaskOrderID, serviceItemID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		rejectionReason := "No justification given"
		params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
			HTTPRequest:      req,
			IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
			Body:             &ghcmessages.PatchMTOServiceItemStatusPayload{Status: "REJECTED", RejectionReason: &rejectionReason},
			MoveTaskOrderID:  move.ID.String(),
			MtoServiceItemID: mtoServiceItem.ID.String(),
		}

		fetcher := fetch.NewFetcher(queryBuilder)
		moveRouter := moverouter.NewMoveRouter()
		addressCreator := address.NewAddressCreator()
		mtoServiceItemStatusUpdater := mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder, moveRouter, addressCreator)

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: mtoServiceItemStatusUpdater,
			Fetcher:               fetcher,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)
		okResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK)

		// Validate outgoing payload
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		suite.Equal(string(models.MTOServiceItemStatusRejected), string(okResponse.Payload.Status))
		suite.NotNil(okResponse.Payload.RejectedAt)
		suite.Equal(rejectionReason, *okResponse.Payload.RejectionReason)
	})

	// With this we'll do a happy path integration test to ensure that the use of the service object
	// by the handler is working as expected.
	suite.Run("Successful status update of MTO service item and event trigger", func() {
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moverouter.NewMoveRouter()
		addressCreator := address.NewAddressCreator()
		mtoServiceItem, availableMove := suite.createServiceItem()
		requestUser := factory.BuildUser(nil, nil, nil)
		availableMoveID := availableMove.ID
		mtoServiceItemID := mtoServiceItem.ID

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/mto_service_items/%s/status", availableMoveID, mtoServiceItemID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
			HTTPRequest:      req,
			IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
			Body:             &ghcmessages.PatchMTOServiceItemStatusPayload{Status: "APPROVED"},
			MoveTaskOrderID:  availableMoveID.String(),
			MtoServiceItemID: mtoServiceItemID.String(),
		}

		fetcher := fetch.NewFetcher(queryBuilder)
		mtoServiceItemStatusUpdater := mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder, moveRouter, addressCreator)

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: mtoServiceItemStatusUpdater,
			Fetcher:               fetcher,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)
		okResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK)

		// Validate outgoing payload
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		suite.Equal(string(models.MTOServiceItemStatusApproved), string(okResponse.Payload.Status))
		suite.NotNil(okResponse.Payload.ApprovedAt)
		suite.HasWebhookNotification(mtoServiceItemID, traceID)

		impactedMove := models.Move{}
		_ = suite.DB().Find(&impactedMove, okResponse.Payload.MoveTaskOrderID)
		suite.Equal(models.MoveStatusAPPROVED, impactedMove.Status)
	})
}

func (suite *HandlerSuite) TestCreateSITAddressUpdate() {
	mockPlanner := &routemocks.Planner{}
	mockedDistance := 55
	mockPlanner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(mockedDistance, nil)
	serviceItemUpdater := mtoserviceitem.NewMTOServiceItemUpdater(
		query.NewQueryBuilder(),
		moverouter.NewMoveRouter(),
		address.NewAddressCreator(),
	)
	sitAddressUpdateCreator := sitaddressupdate.NewApprovedOfficeSITAddressUpdateCreator(
		mockPlanner,
		address.NewAddressCreator(),
		serviceItemUpdater,
	)

	suite.Run("Returns 200, creates new SIT extension, and updates SIT days allowance on shipment without an allowance when validations pass", func() {
		handlerConfig := suite.HandlerConfig()
		handler := CreateSITAddressUpdateHandler{
			handlerConfig,
			sitAddressUpdateCreator,
		}

		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationOriginalAddress,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		req := httptest.NewRequest("POST", fmt.Sprintf("/service-items/%s/sit-address-update/", serviceItem.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		officeRemarks := "new office remarks"
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		createParams := mtoserviceitemop.CreateSITAddressUpdateParams{
			HTTPRequest: req,
			Body: &ghcmessages.CreateSITAddressUpdate{
				NewAddress: &ghcmessages.Address{
					City:           &newAddress.City,
					Country:        newAddress.Country,
					PostalCode:     &newAddress.PostalCode,
					State:          &newAddress.State,
					StreetAddress1: &newAddress.StreetAddress1,
					StreetAddress2: newAddress.StreetAddress2,
					StreetAddress3: newAddress.StreetAddress3,
				},
				OfficeRemarks: &officeRemarks,
			},
			MtoServiceItemID: *handlers.FmtUUID(serviceItem.ID),
		}

		// Validate incoming payload
		suite.NoError(createParams.Body.Validate(strfmt.Default))

		response := handler.Handle(createParams)
		suite.IsType(&mtoserviceitemop.CreateSITAddressUpdateOK{}, response)
		payload := response.(*mtoserviceitemop.CreateSITAddressUpdateOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Len(payload.SitAddressUpdates, 1)
		suite.Equal(*createParams.Body.NewAddress.City, *payload.SitAddressUpdates[0].NewAddress.City)
		suite.Equal(*createParams.Body.NewAddress.Country, *payload.SitAddressUpdates[0].NewAddress.Country)
		suite.Equal(*createParams.Body.NewAddress.PostalCode, *payload.SitAddressUpdates[0].NewAddress.PostalCode)
		suite.Equal(*createParams.Body.NewAddress.State, *payload.SitAddressUpdates[0].NewAddress.State)
		suite.Equal(*createParams.Body.NewAddress.StreetAddress1, *payload.SitAddressUpdates[0].NewAddress.StreetAddress1)
		suite.Equal(*createParams.Body.NewAddress.StreetAddress2, *payload.SitAddressUpdates[0].NewAddress.StreetAddress2)
		suite.Equal(*createParams.Body.NewAddress.StreetAddress3, *payload.SitAddressUpdates[0].NewAddress.StreetAddress3)

		suite.Equal(*createParams.Body.NewAddress.City, *payload.SitDestinationFinalAddress.City)
		suite.Equal(*createParams.Body.NewAddress.Country, *payload.SitDestinationFinalAddress.Country)
		suite.Equal(*createParams.Body.NewAddress.PostalCode, *payload.SitDestinationFinalAddress.PostalCode)
		suite.Equal(*createParams.Body.NewAddress.State, *payload.SitDestinationFinalAddress.State)
		suite.Equal(*createParams.Body.NewAddress.StreetAddress1, *payload.SitDestinationFinalAddress.StreetAddress1)
		suite.Equal(*createParams.Body.NewAddress.StreetAddress2, *payload.SitDestinationFinalAddress.StreetAddress2)
		suite.Equal(*createParams.Body.NewAddress.StreetAddress3, *payload.SitDestinationFinalAddress.StreetAddress3)

		suite.Equal(serviceItem.SITDestinationFinalAddress.ID.String(), payload.SitAddressUpdates[0].OldAddress.ID.String())
		suite.Equal(serviceItem.SITDestinationFinalAddress.City, *payload.SitAddressUpdates[0].OldAddress.City)
		suite.Equal(*serviceItem.SITDestinationFinalAddress.Country, *payload.SitAddressUpdates[0].OldAddress.Country)
		suite.Equal(serviceItem.SITDestinationFinalAddress.PostalCode, *payload.SitAddressUpdates[0].OldAddress.PostalCode)
		suite.Equal(serviceItem.SITDestinationFinalAddress.State, *payload.SitAddressUpdates[0].OldAddress.State)
		suite.Equal(serviceItem.SITDestinationFinalAddress.StreetAddress1, *payload.SitAddressUpdates[0].OldAddress.StreetAddress1)
		suite.Equal(*serviceItem.SITDestinationFinalAddress.StreetAddress2, *payload.SitAddressUpdates[0].OldAddress.StreetAddress2)
		suite.Equal(*serviceItem.SITDestinationFinalAddress.StreetAddress3, *payload.SitAddressUpdates[0].OldAddress.StreetAddress3)

		suite.Require().NotNil(*payload.SitAddressUpdates[0].OfficeRemarks)
		suite.Equal(officeRemarks, *payload.SitAddressUpdates[0].OfficeRemarks)
	})

	suite.Run("Returns a 403 when the office user is not a TOO", func() {
		handlerConfig := suite.HandlerConfig()
		handler := CreateSITAddressUpdateHandler{
			handlerConfig,
			sitAddressUpdateCreator,
		}
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationOriginalAddress,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		req := httptest.NewRequest("POST", fmt.Sprintf("/service-items/%s/sit-address-update/", serviceItem.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		officeRemarks := "new office remarks"
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		createParams := mtoserviceitemop.CreateSITAddressUpdateParams{
			HTTPRequest: req,
			Body: &ghcmessages.CreateSITAddressUpdate{
				NewAddress: &ghcmessages.Address{
					City:           &newAddress.City,
					Country:        newAddress.Country,
					PostalCode:     &newAddress.PostalCode,
					State:          &newAddress.State,
					StreetAddress1: &newAddress.StreetAddress1,
					StreetAddress2: newAddress.StreetAddress2,
					StreetAddress3: newAddress.StreetAddress3,
				},
				OfficeRemarks: &officeRemarks,
			},
			MtoServiceItemID: *handlers.FmtUUID(serviceItem.ID),
		}

		// Validate incoming payload
		suite.NoError(createParams.Body.Validate(strfmt.Default))

		response := handler.Handle(createParams)
		suite.IsType(&mtoserviceitemop.CreateSITAddressUpdateForbidden{}, response)
		payload := response.(*mtoserviceitemop.CreateSITAddressUpdateForbidden).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.IsType(&ghcmessages.Error{}, payload)
	})

	suite.Run("Returns 404 when creator returns NotFoundError", func() {
		creator := &mocks.ApprovedSITAddressUpdateRequestCreator{}
		creator.On(
			"CreateApprovedSITAddressUpdate",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.SITAddressUpdate")).
			Return(nil, apperror.NotFoundError{})

		handlerConfig := suite.HandlerConfig()
		handler := CreateSITAddressUpdateHandler{
			handlerConfig,
			creator,
		}

		fakeID := uuid.Must(uuid.NewV4())

		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		req := httptest.NewRequest("POST", fmt.Sprintf("/service-items/%s/sit-address-update/", fakeID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		officeRemarks := "new office remarks"
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		createParams := mtoserviceitemop.CreateSITAddressUpdateParams{
			HTTPRequest: req,
			Body: &ghcmessages.CreateSITAddressUpdate{
				NewAddress: &ghcmessages.Address{
					City:           &newAddress.City,
					Country:        newAddress.Country,
					PostalCode:     &newAddress.PostalCode,
					State:          &newAddress.State,
					StreetAddress1: &newAddress.StreetAddress1,
					StreetAddress2: newAddress.StreetAddress2,
					StreetAddress3: newAddress.StreetAddress3,
				},
				OfficeRemarks: &officeRemarks,
			},
			MtoServiceItemID: *handlers.FmtUUID(fakeID),
		}

		// Validate incoming payload
		suite.NoError(createParams.Body.Validate(strfmt.Default))

		response := handler.Handle(createParams)
		suite.IsType(&mtoserviceitemop.CreateSITAddressUpdateNotFound{}, response)
		payload := response.(*mtoserviceitemop.CreateSITAddressUpdateNotFound).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
		suite.IsType(&ghcmessages.Error{}, payload)
	})

	suite.Run("Returns 422 when creator returns validation errors", func() {
		handlerConfig := suite.HandlerConfig()
		handler := CreateSITAddressUpdateHandler{
			handlerConfig,
			sitAddressUpdateCreator,
		}

		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusRejected,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationOriginalAddress,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationFinalAddress,
			},
		}, nil)

		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		req := httptest.NewRequest("POST", fmt.Sprintf("/service-items/%s/sit-address-update/", serviceItem.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		officeRemarks := "new office remarks"
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		createParams := mtoserviceitemop.CreateSITAddressUpdateParams{
			HTTPRequest: req,
			Body: &ghcmessages.CreateSITAddressUpdate{
				NewAddress: &ghcmessages.Address{
					City:           &newAddress.City,
					Country:        newAddress.Country,
					PostalCode:     &newAddress.PostalCode,
					State:          &newAddress.State,
					StreetAddress1: &newAddress.StreetAddress1,
					StreetAddress2: newAddress.StreetAddress2,
					StreetAddress3: newAddress.StreetAddress3,
				},
				OfficeRemarks: &officeRemarks,
			},
			MtoServiceItemID: *handlers.FmtUUID(serviceItem.ID),
		}

		// Validate incoming payload
		suite.NoError(createParams.Body.Validate(strfmt.Default))

		response := handler.Handle(createParams)
		suite.IsType(&mtoserviceitemop.CreateSITAddressUpdateUnprocessableEntity{}, response)
		payload := response.(*mtoserviceitemop.CreateSITAddressUpdateUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 500 when approver returns unexpected error", func() {
		creator := &mocks.ApprovedSITAddressUpdateRequestCreator{}
		creator.On(
			"CreateApprovedSITAddressUpdate",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.SITAddressUpdate")).
			Return(nil, apperror.InternalServerError{})

		handlerConfig := suite.HandlerConfig()
		handler := CreateSITAddressUpdateHandler{
			handlerConfig,
			creator,
		}

		fakeID := uuid.Must(uuid.NewV4())

		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		req := httptest.NewRequest("POST", fmt.Sprintf("/service-items/%s/sit-address-update/", fakeID), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		officeRemarks := "new office remarks"
		newAddress := factory.BuildAddress(nil, nil, []factory.Trait{factory.GetTraitAddress3})
		createParams := mtoserviceitemop.CreateSITAddressUpdateParams{
			HTTPRequest: req,
			Body: &ghcmessages.CreateSITAddressUpdate{
				NewAddress: &ghcmessages.Address{
					City:           &newAddress.City,
					Country:        newAddress.Country,
					PostalCode:     &newAddress.PostalCode,
					State:          &newAddress.State,
					StreetAddress1: &newAddress.StreetAddress1,
					StreetAddress2: newAddress.StreetAddress2,
					StreetAddress3: newAddress.StreetAddress3,
				},
				OfficeRemarks: &officeRemarks,
			},
			MtoServiceItemID: *handlers.FmtUUID(fakeID),
		}

		// Validate incoming payload
		suite.NoError(createParams.Body.Validate(strfmt.Default))

		response := handler.Handle(createParams)
		suite.IsType(&mtoserviceitemop.CreateSITAddressUpdateInternalServerError{}, response)
		payload := response.(*mtoserviceitemop.CreateSITAddressUpdateInternalServerError).Payload

		// Validate outgoing payload
		suite.Nil(payload)
	})
}
