package publicapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gobuffalo/validate"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"

	"github.com/transcom/mymove/mocks"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	sitop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/storage_in_transits"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestIndexStorageInTransitsHandler() {
	shipmentID, err := uuid.NewV4()
	suite.NoError(err)

	storageInTransitID, err := uuid.NewV4()
	suite.NoError(err)

	storageInTransitID2, err := uuid.NewV4()
	suite.NoError(err)

	officeUserID, err := uuid.NewV4()
	suite.NoError(err)

	userID, err := uuid.NewV4()
	suite.NoError(err)

	user := models.OfficeUser{ID: officeUserID, UserID: &userID}

	path := fmt.Sprintf("/shipments/%s/storage_in_transits", shipmentID)
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateOfficeRequest(req, user)
	params := sitop.IndexStorageInTransitsParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipmentID.String()),
	}

	storageInTransitIndexer := &mocks.StorageInTransitsIndexer{}

	handler := IndexStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(),
		suite.TestLogger()),
		storageInTransitIndexer,
	}

	returnSits := []models.StorageInTransit{
		{ID: storageInTransitID},
		{ID: storageInTransitID2},
	}

	// Happy path
	storageInTransitIndexer.On("IndexStorageInTransits",
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
	).Return(returnSits, nil).Once()

	response := handler.Handle(params)
	suite.Assertions.IsType(&sitop.IndexStorageInTransitsOK{}, response)
	responsePayload := response.(*sitop.IndexStorageInTransitsOK).Payload
	suite.Equal(2, len(responsePayload))

	expectedError := models.ErrFetchForbidden
	storageInTransitIndexer.On("IndexStorageInTransits",
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
	).Return(nil, expectedError).Once()

	// Forbidden Scenario
	response = handler.Handle(params)
	expectedResponse := &handlers.ErrResponse{
		Code: http.StatusForbidden,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)

}

func (suite *HandlerSuite) TestGetStorageInTransitHandler() {
	shipmentID, err := uuid.NewV4()
	suite.NoError(err)
	storageInTransitID, err := uuid.NewV4()
	suite.NoError(err)
	sitPayload := apimessages.StorageInTransit{
		ID: strfmt.UUID(storageInTransitID.String()),
	}

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s", shipmentID, storageInTransitID)
	req := httptest.NewRequest("GET", path, nil)
	officeUserID, err := uuid.NewV4()
	suite.NoError(err)
	userID, err := uuid.NewV4()
	suite.NoError(err)
	user := models.OfficeUser{ID: officeUserID, UserID: &userID}
	req = suite.AuthenticateOfficeRequest(req, user)
	params := sitop.GetStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipmentID.String()),
		StorageInTransitID: strfmt.UUID(storageInTransitID.String()),
	}

	storageInTransitByIDFetcher := &mocks.StorageInTransitByIDFetcher{}
	handler := GetStorageInTransitHandler{
		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		storageInTransitByIDFetcher,
	}
	returnSit := models.StorageInTransit{ID: storageInTransitID}
	// Happy path
	storageInTransitByIDFetcher.On(
		"FetchStorageInTransitByID",
		storageInTransitID,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
	).Return(&returnSit, nil).Once()
	response := handler.Handle(params)
	suite.Assertions.IsType(&sitop.GetStorageInTransitOK{}, response)

	responsePayload := response.(*sitop.GetStorageInTransitOK).Payload

	suite.Equal(sitPayload.ID, responsePayload.ID)

	// Forbidden scenario
	expectedError := models.ErrFetchForbidden
	storageInTransitByIDFetcher.On(
		"FetchStorageInTransitByID",
		storageInTransitID,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
	).Return(nil, expectedError).Once()

	response = handler.Handle(params)

	expectedResponse := &handlers.ErrResponse{
		Code: http.StatusForbidden,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)
}

func (suite *HandlerSuite) TestCreateStorageInTransitHandler() {
	shipmentID, err := uuid.NewV4()
	suite.NoError(err)

	storageInTransitID, err := uuid.NewV4()
	suite.NoError(err)
	storageInTransit := models.StorageInTransit{ID: storageInTransitID}
	payload := payloadForStorageInTransitModel(&storageInTransit)

	tspUserID, err := uuid.NewV4()
	suite.NoError(err)

	userID, err := uuid.NewV4()
	suite.NoError(err)

	tspUser := models.TspUser{ID: tspUserID, UserID: &userID}

	path := fmt.Sprintf("/shipments/%s/storage_in_transits", shipmentID)
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := sitop.CreateStorageInTransitParams{
		HTTPRequest:      req,
		ShipmentID:       strfmt.UUID(shipmentID.String()),
		StorageInTransit: payload,
	}

	storageInTransitCreator := &mocks.StorageInTransitCreator{}

	handler := CreateStorageInTransitHandler{
		handlers.NewHandlerContext(suite.DB(),
			suite.TestLogger()),
		storageInTransitCreator,
	}
	// Happy path
	storageInTransitCreator.On("CreateStorageInTransit",
		*payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
	).Return(&storageInTransit, validate.NewErrors(), nil).Once()

	response := handler.Handle(params)
	suite.Assertions.IsType(&sitop.CreateStorageInTransitCreated{}, response)
	responsePayload := response.(*sitop.CreateStorageInTransitCreated).Payload
	suite.Equal(storageInTransit.ID.String(), responsePayload.ID.String())

	// Forbidden scenario
	expectedError := models.ErrFetchForbidden
	storageInTransitCreator.On("CreateStorageInTransit",
		*payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
	).Return(nil, validate.NewErrors(), expectedError).Once()

	response = handler.Handle(params)
	expectedResponse := &handlers.ErrResponse{
		Code: http.StatusForbidden,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)
}

func (suite *HandlerSuite) TestApproveStorageInTransitHandler() {
	shipmentID, err := uuid.NewV4()
	suite.NoError(err)

	storageInTransitID, err := uuid.NewV4()
	suite.NoError(err)
	storageInTransit := models.StorageInTransit{ID: storageInTransitID}

	officeUserID, err := uuid.NewV4()
	suite.NoError(err)

	userID, err := uuid.NewV4()
	suite.NoError(err)

	user := models.OfficeUser{ID: officeUserID, UserID: &userID}

	payload := apimessages.StorageInTransitApprovalPayload{
		AuthorizedStartDate: *handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
		AuthorizationNotes:  handlers.FmtString("looks good to me"),
	}

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s/approve", shipmentID, storageInTransitID)
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateOfficeRequest(req, user)
	params := sitop.ApproveStorageInTransitParams{
		HTTPRequest:                     req,
		ShipmentID:                      strfmt.UUID(shipmentID.String()),
		StorageInTransitID:              strfmt.UUID(storageInTransitID.String()),
		StorageInTransitApprovalPayload: &payload,
	}

	storageInTransitApprover := &mocks.StorageInTransitApprover{}

	handler := ApproveStorageInTransitHandler{
		handlers.NewHandlerContext(suite.DB(),
			suite.TestLogger()),
		storageInTransitApprover,
	}
	// Happy path
	storageInTransitApprover.On("ApproveStorageInTransit",
		payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(&storageInTransit, validate.NewErrors(), nil).Once()

	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.ApproveStorageInTransitOK{}, response)
	responsePayload := response.(*sitop.ApproveStorageInTransitOK).Payload
	suite.Equal(storageInTransitID.String(), responsePayload.ID.String())

	// Forbidden scenario
	expectedError := models.ErrFetchForbidden
	storageInTransitApprover.On("ApproveStorageInTransit",
		payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(nil, validate.NewErrors(), expectedError).Once()

	response = handler.Handle(params)
	expectedResponse := &handlers.ErrResponse{
		Code: http.StatusForbidden,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)

	// Write conflict scenario
	expectedError = models.ErrWriteConflict
	storageInTransitApprover.On("ApproveStorageInTransit",
		payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(nil, validate.NewErrors(), expectedError).Once()

	response = handler.Handle(params)
	expectedResponse = &handlers.ErrResponse{
		Code: http.StatusConflict,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)
}

func (suite *HandlerSuite) TestDenyStorageInTransitHandler() {
	shipmentID, err := uuid.NewV4()
	suite.NoError(err)

	storageInTransitID, err := uuid.NewV4()
	suite.NoError(err)
	storageInTransit := models.StorageInTransit{ID: storageInTransitID}

	officeUserID, err := uuid.NewV4()
	suite.NoError(err)

	userID, err := uuid.NewV4()
	suite.NoError(err)

	user := models.OfficeUser{ID: officeUserID, UserID: &userID}

	payload := apimessages.StorageInTransitDenialPayload{
		AuthorizationNotes: *handlers.FmtString("looks good to me"),
	}

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s/deny", shipmentID, storageInTransitID)
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateOfficeRequest(req, user)
	params := sitop.DenyStorageInTransitParams{
		HTTPRequest:                 req,
		ShipmentID:                  strfmt.UUID(shipmentID.String()),
		StorageInTransitID:          strfmt.UUID(storageInTransitID.String()),
		StorageInTransitDenyPayload: &payload,
	}

	storageInTransitDenier := &mocks.StorageInTransitDenier{}

	handler := DenyStorageInTransitHandler{
		handlers.NewHandlerContext(suite.DB(),
			suite.TestLogger()),
		storageInTransitDenier,
	}
	// Happy path
	storageInTransitDenier.On("DenyStorageInTransit",
		payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(&storageInTransit, validate.NewErrors(), nil).Once()

	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.DenyStorageInTransitOK{}, response)
	responsePayload := response.(*sitop.DenyStorageInTransitOK).Payload
	suite.Equal(storageInTransitID.String(), responsePayload.ID.String())

	// Forbidden scenario
	expectedError := models.ErrFetchForbidden
	storageInTransitDenier.On("DenyStorageInTransit",
		payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(nil, validate.NewErrors(), expectedError).Once()

	response = handler.Handle(params)
	expectedResponse := &handlers.ErrResponse{
		Code: http.StatusForbidden,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)

	// Write conflict scenario
	expectedError = models.ErrWriteConflict
	storageInTransitDenier.On("DenyStorageInTransit",
		payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(nil, validate.NewErrors(), expectedError).Once()

	response = handler.Handle(params)
	expectedResponse = &handlers.ErrResponse{
		Code: http.StatusConflict,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)

}

func (suite *HandlerSuite) TestInSitStorageInTransitHandler() {
	shipmentID, err := uuid.NewV4()
	suite.NoError(err)

	storageInTransitID, err := uuid.NewV4()
	suite.NoError(err)
	storageInTransit := models.StorageInTransit{ID: storageInTransitID}

	tspUserID, err := uuid.NewV4()
	suite.NoError(err)

	userID, err := uuid.NewV4()
	suite.NoError(err)

	user := models.TspUser{ID: tspUserID, UserID: &userID}

	payload := apimessages.StorageInTransitInSitPayload{
		ActualStartDate: *handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
	}

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s/place_into_sit", shipmentID, storageInTransitID)
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, user)
	params := sitop.InSitStorageInTransitParams{
		HTTPRequest:                  req,
		ShipmentID:                   strfmt.UUID(shipmentID.String()),
		StorageInTransitID:           strfmt.UUID(storageInTransitID.String()),
		StorageInTransitInSitPayload: &payload,
	}

	storageInTransitInSITPlacer := &mocks.StorageInTransitInSITPlacer{}

	handler := InSitStorageInTransitHandler{
		handlers.NewHandlerContext(suite.DB(),
			suite.TestLogger()),
		storageInTransitInSITPlacer,
	}
	// Happy path
	storageInTransitInSITPlacer.On("PlaceIntoSITStorageInTransit",
		payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(&storageInTransit, validate.NewErrors(), nil).Once()

	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.InSitStorageInTransitOK{}, response)
	responsePayload := response.(*sitop.InSitStorageInTransitOK).Payload
	suite.Equal(storageInTransitID.String(), responsePayload.ID.String())

	// Forbidden scenario
	expectedError := models.ErrFetchForbidden
	storageInTransitInSITPlacer.On("PlaceIntoSITStorageInTransit",
		payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(nil, validate.NewErrors(), expectedError).Once()

	response = handler.Handle(params)
	expectedResponse := &handlers.ErrResponse{
		Code: http.StatusForbidden,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)

	// Write conflict scenario
	expectedError = models.ErrWriteConflict
	storageInTransitInSITPlacer.On("PlaceIntoSITStorageInTransit",
		payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(nil, validate.NewErrors(), expectedError).Once()

	response = handler.Handle(params)
	expectedResponse = &handlers.ErrResponse{
		Code: http.StatusConflict,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)
}

func (suite *HandlerSuite) TestDeliverStorageInTransitHandler() {
	shipmentID, err := uuid.NewV4()
	suite.NoError(err)

	storageInTransitID, err := uuid.NewV4()
	suite.NoError(err)
	storageInTransit := models.StorageInTransit{ID: storageInTransitID}

	tspUserID, err := uuid.NewV4()
	suite.NoError(err)

	userID, err := uuid.NewV4()
	suite.NoError(err)

	user := models.TspUser{ID: tspUserID, UserID: &userID}

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s/deliver", shipmentID, storageInTransitID)
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, user)
	params := sitop.DeliverStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipmentID.String()),
		StorageInTransitID: strfmt.UUID(storageInTransitID.String()),
	}

	storageInTransitInSITDeliverer := &mocks.StorageInTransitDeliverer{}

	handler := DeliverStorageInTransitHandler{
		handlers.NewHandlerContext(suite.DB(),
			suite.TestLogger()),
		storageInTransitInSITDeliverer,
	}
	// Happy path
	storageInTransitInSITDeliverer.On("DeliverStorageInTransit",
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(&storageInTransit, validate.NewErrors(), nil).Once()

	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.DeliverStorageInTransitOK{}, response)
	responsePayload := response.(*sitop.DeliverStorageInTransitOK).Payload
	suite.Equal(storageInTransitID.String(), responsePayload.ID.String())

	// Forbidden scenario
	expectedError := models.ErrFetchForbidden
	storageInTransitInSITDeliverer.On("DeliverStorageInTransit",
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(nil, validate.NewErrors(), expectedError).Once()

	response = handler.Handle(params)
	expectedResponse := &handlers.ErrResponse{
		Code: http.StatusForbidden,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)

	// Write conflict scenario
	expectedError = models.ErrWriteConflict
	storageInTransitInSITDeliverer.On("DeliverStorageInTransit",
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(nil, validate.NewErrors(), expectedError).Once()

	response = handler.Handle(params)
	expectedResponse = &handlers.ErrResponse{
		Code: http.StatusConflict,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)
}

func (suite *HandlerSuite) TestReleaseStorageInTransitHandler() {
	shipmentID, err := uuid.NewV4()
	suite.NoError(err)

	storageInTransitID, err := uuid.NewV4()
	suite.NoError(err)
	storageInTransit := models.StorageInTransit{ID: storageInTransitID}

	tspUserID, err := uuid.NewV4()
	suite.NoError(err)

	userID, err := uuid.NewV4()
	suite.NoError(err)

	user := models.TspUser{ID: tspUserID, UserID: &userID}
	payload := apimessages.StorageInTransitReleasePayload{
		ReleasedOn: *handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
	}

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s/release", shipmentID, storageInTransitID)
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, user)
	params := sitop.ReleaseStorageInTransitParams{
		HTTPRequest:                      req,
		ShipmentID:                       strfmt.UUID(shipmentID.String()),
		StorageInTransitID:               strfmt.UUID(storageInTransitID.String()),
		StorageInTransitOnReleasePayload: &payload,
	}

	storageInTransitReleaser := &mocks.StorageInTransitReleaser{}

	handler := ReleaseStorageInTransitHandler{
		handlers.NewHandlerContext(suite.DB(),
			suite.TestLogger()),
		storageInTransitReleaser,
	}
	// Happy path
	storageInTransitReleaser.On("ReleaseStorageInTransit",
		payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(&storageInTransit, validate.NewErrors(), nil).Once()

	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.ReleaseStorageInTransitOK{}, response)
	responsePayload := response.(*sitop.ReleaseStorageInTransitOK).Payload
	suite.Equal(storageInTransitID.String(), responsePayload.ID.String())

	// Forbidden scenario
	expectedError := models.ErrFetchForbidden
	storageInTransitReleaser.On("ReleaseStorageInTransit",
		payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(nil, validate.NewErrors(), expectedError).Once()

	response = handler.Handle(params)
	expectedResponse := &handlers.ErrResponse{
		Code: http.StatusForbidden,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)

	// Write conflict scenario
	expectedError = models.ErrWriteConflict
	storageInTransitReleaser.On("ReleaseStorageInTransit",
		payload,
		shipmentID,
		auth.SessionFromRequestContext(params.HTTPRequest),
		storageInTransitID,
	).Return(nil, validate.NewErrors(), expectedError).Once()

	response = handler.Handle(params)
	expectedResponse = &handlers.ErrResponse{
		Code: http.StatusConflict,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)
}

func (suite *HandlerSuite) TestPatchStorageInTransitHandler() {
	shipmentID, err := uuid.NewV4()
	suite.NoError(err)

	storageInTransitID, err := uuid.NewV4()
	suite.NoError(err)

	officeUserID, err := uuid.NewV4()
	suite.NoError(err)

	userID, err := uuid.NewV4()
	suite.NoError(err)

	user := models.OfficeUser{ID: officeUserID, UserID: &userID}
	storageInTransit := models.StorageInTransit{ID: storageInTransitID, ShipmentID: shipmentID}
	payload := apimessages.StorageInTransit{
		ID: *handlers.FmtUUID(storageInTransitID),
	}

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s", shipmentID, storageInTransitID)
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateOfficeRequest(req, user)

	storageInTransitPatcher := &mocks.StorageInTransitPatcher{}

	params := sitop.PatchStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipmentID.String()),
		StorageInTransitID: strfmt.UUID(storageInTransitID.String()),
		StorageInTransit:   &payload,
	}

	// Happy path
	storageInTransitPatcher.On("PatchStorageInTransit",
		payload,
		shipmentID,
		storageInTransitID,
		auth.SessionFromRequestContext(params.HTTPRequest),
	).Return(&storageInTransit, validate.NewErrors(), nil).Once()

	handler := PatchStorageInTransitHandler{
		handlers.NewHandlerContext(suite.DB(),
			suite.TestLogger()),
		storageInTransitPatcher,
	}

	response := handler.Handle(params)
	suite.Assertions.IsType(&sitop.PatchStorageInTransitOK{}, response)
	responsePayload := response.(*sitop.PatchStorageInTransitOK).Payload
	suite.Equal(storageInTransit.ID.String(), responsePayload.ID.String())

	// Forbidden scenario
	expectedError := models.ErrFetchForbidden
	storageInTransitPatcher.On("PatchStorageInTransit",
		payload,
		shipmentID,
		storageInTransitID,
		auth.SessionFromRequestContext(params.HTTPRequest),
	).Return(nil, validate.NewErrors(), expectedError).Once()

	response = handler.Handle(params)
	expectedResponse := &handlers.ErrResponse{
		Code: http.StatusForbidden,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)
}

func (suite *HandlerSuite) TestDeleteStorageInTransitHandler() {

	shipmentID, err := uuid.NewV4()
	suite.NoError(err)

	storageInTransitID, err := uuid.NewV4()
	suite.NoError(err)

	tspUserID, err := uuid.NewV4()
	suite.NoError(err)

	userID, err := uuid.NewV4()
	suite.NoError(err)

	user := models.TspUser{ID: tspUserID, UserID: &userID}

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s", shipmentID, storageInTransitID)
	req := httptest.NewRequest("DELETE", path, nil)
	req = suite.AuthenticateTspRequest(req, user)
	params := sitop.DeleteStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipmentID.String()),
		StorageInTransitID: strfmt.UUID(storageInTransitID.String()),
	}

	storageInTransitDeleter := &mocks.StorageInTransitDeleter{}

	handler := DeleteStorageInTransitHandler{
		handlers.NewHandlerContext(suite.DB(),
			suite.TestLogger()),
		storageInTransitDeleter,
	}
	// Happy path
	storageInTransitDeleter.On("DeleteStorageInTransit",
		shipmentID,
		storageInTransitID,
		auth.SessionFromRequestContext(params.HTTPRequest),
	).Return(nil).Once()

	response := handler.Handle(params)
	suite.Assertions.IsType(&sitop.DeleteStorageInTransitOK{}, response)

	// Forbidden scenario
	expectedError := models.ErrFetchForbidden
	storageInTransitDeleter.On("DeleteStorageInTransit",
		shipmentID,
		storageInTransitID,
		auth.SessionFromRequestContext(params.HTTPRequest),
	).Return(expectedError).Once()

	response = handler.Handle(params)
	expectedResponse := &handlers.ErrResponse{
		Code: http.StatusForbidden,
		Err:  expectedError,
	}
	suite.Equal(expectedResponse, response)
}
