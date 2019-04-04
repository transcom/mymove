package publicapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	sitop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/storage_in_transits"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func setupStorageInTransitHandlerTest(suite *HandlerSuite) (shipment models.Shipment, sit models.StorageInTransit, user models.OfficeUser) {

	shipment = testdatagen.MakeDefaultShipment(suite.DB())
	user = testdatagen.MakeDefaultOfficeUser(suite.DB())

	assertions := testdatagen.Assertions{
		StorageInTransit: models.StorageInTransit{
			Location:           models.StorageInTransitLocationORIGIN,
			ShipmentID:         shipment.ID,
			EstimatedStartDate: testdatagen.DateInsidePeakRateCycle,
		},
	}
	testdatagen.MakeStorageInTransit(suite.DB(), assertions)
	sit = testdatagen.MakeStorageInTransit(suite.DB(), assertions)

	return shipment, sit, user
}

func storageInTransitPayloadCompare(suite *HandlerSuite, expected *apimessages.StorageInTransit, actual *apimessages.StorageInTransit) {
	suite.Equal(*expected.WarehouseEmail, *actual.WarehouseEmail)
	suite.Equal(*expected.Notes, *actual.Notes)
	suite.Equal(*expected.WarehouseID, *actual.WarehouseID)
	suite.Equal(*expected.Location, *actual.Location)
	suite.Equal(*expected.WarehouseName, *actual.WarehouseName)
	suite.Equal(*expected.WarehousePhone, *actual.WarehousePhone)
	suite.Equal(expected.EstimatedStartDate.String(), actual.EstimatedStartDate.String())
	suite.Equal(expected.Status, actual.Status)
}

func (suite *HandlerSuite) TestIndexStorageInTransitsHandler() {

	shipment, _, user := setupStorageInTransitHandlerTest(suite)

	path := fmt.Sprintf("/shipments/%s/storage_in_transits", shipment.ID.String())
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateOfficeRequest(req, user)
	params := sitop.IndexStorageInTransitsParams{
		HTTPRequest: req,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
	}

	handler := IndexStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.IndexStorageInTransitsOK{}, response)
	okResponse := response.(*sitop.IndexStorageInTransitsOK)

	suite.Equal(2, len(okResponse.Payload))

	// Let's make sure it doesn't authorize with a servicemember user
	serviceMember := testdatagen.MakeDefaultServiceMember(suite.DB())
	serviceMemberReq := httptest.NewRequest("GET", path, nil)
	serviceMemberReq = suite.AuthenticateRequest(serviceMemberReq, serviceMember)
	params = sitop.IndexStorageInTransitsParams{
		HTTPRequest: serviceMemberReq,
		ShipmentID:  strfmt.UUID(shipment.ID.String()),
	}

	handler = IndexStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	_ = handler.Handle(params)

	suite.Assertions.Error(models.ErrFetchForbidden)

}

func (suite *HandlerSuite) TestGetStorageInTransitHandler() {
	shipment, sit, user := setupStorageInTransitHandlerTest(suite)
	sitPayload := payloadForStorageInTransitModel(&sit)

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s", shipment.ID.String(), sit.ID.String())
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateOfficeRequest(req, user)
	params := sitop.GetStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipment.ID.String()),
		StorageInTransitID: strfmt.UUID(sit.ID.String()),
	}

	handler := GetStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)
	suite.Assertions.IsType(&sitop.GetStorageInTransitOK{}, response)

	responsePayload := response.(*sitop.GetStorageInTransitOK).Payload

	storageInTransitPayloadCompare(suite, sitPayload, responsePayload)

	// Let's make sure it fails when a TSP who doesn't own the shipment tries to do a GET on this
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	req = httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	params.HTTPRequest = req
	handler.Handle(params)
	suite.Error(models.ErrFetchForbidden)
}

func (suite *HandlerSuite) TestCreateStorageInTransitHandler() {

	shipment, sit, _ := setupStorageInTransitHandlerTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())

	sit.WarehouseID = "12345"

	sitPayload := payloadForStorageInTransitModel(&sit)

	path := fmt.Sprintf("/shipments/%s/storage_in_transits", shipment.ID.String())
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)

	params := sitop.CreateStorageInTransitParams{
		HTTPRequest:      req,
		ShipmentID:       strfmt.UUID(shipment.ID.String()),
		StorageInTransit: sitPayload,
	}

	handler := CreateStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	_ = handler.Handle(params)

	// we expect this to fail with a forbidden message. The generated TSP does not have rights to the shipment.
	suite.Error(models.ErrFetchForbidden)

	// Now let's do a working one
	assertions := testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			ShipmentID:                      shipment.ID,
		},
	}
	// Create a shipment offer that uses our generated TSP ID and shipment ID so that our TSP has rights to
	// Use these to create a SIT for them.
	testdatagen.MakeShipmentOffer(suite.DB(), assertions)

	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.CreateStorageInTransitCreated{}, response)
	responsePayload := response.(*sitop.CreateStorageInTransitCreated).Payload
	storageInTransitPayloadCompare(suite, sitPayload, responsePayload)
}

func (suite *HandlerSuite) TestApproveStorageInTransitHandler() {
	shipment, sit, user := setupStorageInTransitHandlerTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())

	approvePayload := apimessages.StorageInTransitApprovalPayload{
		AuthorizedStartDate: handlers.FmtDate(testdatagen.DateInsidePeakRateCycle),
		AuthorizationNotes:  *handlers.FmtString("looks good to me"),
	}

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s/approve", shipment.ID.String(), sit.ID.String())
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateOfficeRequest(req, user)
	params := sitop.ApproveStorageInTransitParams{
		HTTPRequest:                     req,
		ShipmentID:                      strfmt.UUID(shipment.ID.String()),
		StorageInTransitID:              strfmt.UUID(sit.ID.String()),
		StorageInTransitApprovalPayload: &approvePayload,
	}

	handler := ApproveStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.ApproveStorageInTransitOK{}, response)
	responsePayload := response.(*sitop.ApproveStorageInTransitOK).Payload

	suite.Equal(string(models.StorageInTransitStatusAPPROVED), responsePayload.Status)

	// Let's make sure it denies a TSP user.
	req = suite.AuthenticateTspRequest(req, tspUser)
	params = sitop.ApproveStorageInTransitParams{
		HTTPRequest:                     req,
		ShipmentID:                      strfmt.UUID(shipment.ID.String()),
		StorageInTransitID:              strfmt.UUID(sit.ID.String()),
		StorageInTransitApprovalPayload: &approvePayload,
	}

	handler = ApproveStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)
	suite.Assertions.IsType(response.(*sitop.ApproveStorageInTransitForbidden), response)

	// Let's make sure it doesn't work if the status is delivered
	sit.Status = models.StorageInTransitStatusDELIVERED
	_, _ = suite.DB().ValidateAndSave(&sit)

	req = suite.AuthenticateOfficeRequest(req, user)
	params = sitop.ApproveStorageInTransitParams{
		HTTPRequest:                     req,
		ShipmentID:                      strfmt.UUID(shipment.ID.String()),
		StorageInTransitID:              strfmt.UUID(sit.ID.String()),
		StorageInTransitApprovalPayload: &approvePayload,
	}

	handler = ApproveStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)
	suite.Assertions.IsType(response.(*sitop.ApproveStorageInTransitConflict), response)

}

func (suite *HandlerSuite) TestDenyStorageInTransitHandler() {
	shipment, sit, user := setupStorageInTransitHandlerTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())

	denyPayload := apimessages.StorageInTransitApprovalPayload{
		AuthorizationNotes: *handlers.FmtString("looks bad to me"),
	}

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s/deny", shipment.ID.String(), sit.ID.String())
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateOfficeRequest(req, user)
	params := sitop.DenyStorageInTransitParams{
		HTTPRequest:                     req,
		ShipmentID:                      strfmt.UUID(shipment.ID.String()),
		StorageInTransitID:              strfmt.UUID(sit.ID.String()),
		StorageInTransitApprovalPayload: &denyPayload,
	}

	handler := DenyStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.DenyStorageInTransitOK{}, response)
	responsePayload := response.(*sitop.DenyStorageInTransitOK).Payload

	suite.Equal(string(models.StorageInTransitStatusDENIED), responsePayload.Status)

	// Let's make sure it denies a TSP user.
	req = suite.AuthenticateTspRequest(req, tspUser)
	params = sitop.DenyStorageInTransitParams{
		HTTPRequest:                     req,
		ShipmentID:                      strfmt.UUID(shipment.ID.String()),
		StorageInTransitID:              strfmt.UUID(sit.ID.String()),
		StorageInTransitApprovalPayload: &denyPayload,
	}

	handler = DenyStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)
	suite.Assertions.IsType(response.(*sitop.DenyStorageInTransitForbidden), response)

	// Let's make sure it doesn't work if the status is delivered
	sit.Status = models.StorageInTransitStatusDELIVERED
	_, _ = suite.DB().ValidateAndSave(&sit)

	req = suite.AuthenticateOfficeRequest(req, user)
	params = sitop.DenyStorageInTransitParams{
		HTTPRequest:                     req,
		ShipmentID:                      strfmt.UUID(shipment.ID.String()),
		StorageInTransitID:              strfmt.UUID(sit.ID.String()),
		StorageInTransitApprovalPayload: &denyPayload,
	}

	handler = DenyStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)
	suite.Assertions.IsType(response.(*sitop.DenyStorageInTransitConflict), response)

}

func (suite *HandlerSuite) TestInSitStorageInTransitHandler() {
	shipment, sit, user := setupStorageInTransitHandlerTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())

	assertions := testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			ShipmentID:                      shipment.ID,
		},
	}
	// Create a shipment offer that uses our generated TSP ID and shipment ID so that our TSP has rights to
	// change the status to in_sit.
	testdatagen.MakeShipmentOffer(suite.DB(), assertions)

	sit.Status = models.StorageInTransitStatusAPPROVED
	_, _ = suite.DB().ValidateAndSave(&sit)

	inSitPayload := apimessages.StorageInTransitInSitPayload{
		ActualStartDate: *handlers.FmtDate(testdatagen.DateInsidePerformancePeriod),
	}

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s/in_sit", shipment.ID.String(), sit.ID.String())
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	params := sitop.InSitStorageInTransitParams{
		HTTPRequest:                  req,
		ShipmentID:                   strfmt.UUID(shipment.ID.String()),
		StorageInTransitID:           strfmt.UUID(sit.ID.String()),
		StorageInTransitInSitPayload: &inSitPayload,
	}

	handler := InSitStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.InSitStorageInTransitOK{}, response)
	responsePayload := response.(*sitop.InSitStorageInTransitOK).Payload

	suite.Equal(string(models.StorageInTransitStatusINSIT), responsePayload.Status)
	suite.Equal(inSitPayload.ActualStartDate, *responsePayload.ActualStartDate)

	// Let's make sure it denies an office user
	req = suite.AuthenticateOfficeRequest(req, user)
	params = sitop.InSitStorageInTransitParams{
		HTTPRequest:                  req,
		ShipmentID:                   strfmt.UUID(shipment.ID.String()),
		StorageInTransitID:           strfmt.UUID(sit.ID.String()),
		StorageInTransitInSitPayload: &inSitPayload,
	}

	handler = InSitStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)

	suite.Assertions.IsType(&sitop.InSitStorageInTransitForbidden{}, response)

	// Let's make sure it won't let us do this if the status is not approved
	req = suite.AuthenticateTspRequest(req, tspUser)
	params = sitop.InSitStorageInTransitParams{
		HTTPRequest:                  req,
		ShipmentID:                   strfmt.UUID(shipment.ID.String()),
		StorageInTransitID:           strfmt.UUID(sit.ID.String()),
		StorageInTransitInSitPayload: &inSitPayload,
	}

	sit.Status = models.StorageInTransitStatusREQUESTED
	_, _ = suite.DB().ValidateAndSave(&sit)

	handler = InSitStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)

	suite.Assertions.IsType(&sitop.InSitStorageInTransitConflict{}, response)

}

func (suite *HandlerSuite) TestDeliverStorageInTransitHandler() {
	shipment, sit, user := setupStorageInTransitHandlerTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())

	assertions := testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			ShipmentID:                      shipment.ID,
		},
	}
	// Create a shipment offer that uses our generated TSP ID and shipment ID so that our TSP has rights to
	// change the status to in_sit.
	testdatagen.MakeShipmentOffer(suite.DB(), assertions)

	sit.Status = models.StorageInTransitStatusINSIT
	_, _ = suite.DB().ValidateAndSave(&sit)

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s/deliver", shipment.ID.String(), sit.ID.String())
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	params := sitop.DeliverStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipment.ID.String()),
		StorageInTransitID: strfmt.UUID(sit.ID.String()),
	}

	handler := DeliverStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.DeliverStorageInTransitOK{}, response)
	responsePayload := response.(*sitop.DeliverStorageInTransitOK).Payload

	suite.Equal(string(models.StorageInTransitStatusDELIVERED), responsePayload.Status)

	// Let's make sure it also works with the released status
	sit.Status = models.StorageInTransitStatusRELEASED
	_, _ = suite.DB().ValidateAndSave(&sit)

	req = suite.AuthenticateTspRequest(req, tspUser)
	params = sitop.DeliverStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipment.ID.String()),
		StorageInTransitID: strfmt.UUID(sit.ID.String()),
	}

	handler = DeliverStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)

	suite.Assertions.IsType(&sitop.DeliverStorageInTransitOK{}, response)
	responsePayload = response.(*sitop.DeliverStorageInTransitOK).Payload

	suite.Equal(string(models.StorageInTransitStatusDELIVERED), responsePayload.Status)

	// Let's make sure it doesn't let us do this if the status isn't in sit or released
	sit.Status = models.StorageInTransitStatusREQUESTED
	_, _ = suite.DB().ValidateAndSave(&sit)

	req = suite.AuthenticateTspRequest(req, tspUser)
	params = sitop.DeliverStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipment.ID.String()),
		StorageInTransitID: strfmt.UUID(sit.ID.String()),
	}

	handler = DeliverStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)

	suite.Assertions.IsType(&sitop.DeliverStorageInTransitConflict{}, response)

	// Let's make sure this fails with an office user
	sit.Status = models.StorageInTransitStatusINSIT
	_, _ = suite.DB().ValidateAndSave(&sit)

	req = httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateOfficeRequest(req, user)
	params = sitop.DeliverStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipment.ID.String()),
		StorageInTransitID: strfmt.UUID(sit.ID.String()),
	}

	handler = DeliverStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)

	suite.Assertions.IsType(&sitop.DeliverStorageInTransitForbidden{}, response)

}

func (suite *HandlerSuite) TestReleaseStorageInTransitHandler() {
	shipment, sit, user := setupStorageInTransitHandlerTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())

	assertions := testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			ShipmentID:                      shipment.ID,
		},
	}
	// Create a shipment offer that uses our generated TSP ID and shipment ID so that our TSP has rights to
	// change the status to in_sit.
	testdatagen.MakeShipmentOffer(suite.DB(), assertions)

	sit.Status = models.StorageInTransitStatusINSIT
	_, _ = suite.DB().ValidateAndSave(&sit)

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s/release", shipment.ID.String(), sit.ID.String())
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	params := sitop.ReleaseStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipment.ID.String()),
		StorageInTransitID: strfmt.UUID(sit.ID.String()),
	}

	handler := ReleaseStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.ReleaseStorageInTransitOK{}, response)
	responsePayload := response.(*sitop.ReleaseStorageInTransitOK).Payload

	suite.Equal(string(models.StorageInTransitStatusRELEASED), responsePayload.Status)

	// Let's make sure this doesn't work if the status isn't 'in sit'
	sit.Status = models.StorageInTransitStatusREQUESTED
	_, _ = suite.DB().ValidateAndSave(&sit)

	req = suite.AuthenticateTspRequest(req, tspUser)
	params = sitop.ReleaseStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipment.ID.String()),
		StorageInTransitID: strfmt.UUID(sit.ID.String()),
	}

	handler = ReleaseStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)

	suite.Assertions.IsType(&sitop.ReleaseStorageInTransitConflict{}, response)

	// Let's make sure it fails if an office user tries to do it
	sit.Status = models.StorageInTransitStatusINSIT
	_, _ = suite.DB().ValidateAndSave(&sit)

	req = suite.AuthenticateOfficeRequest(req, user)
	params = sitop.ReleaseStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipment.ID.String()),
		StorageInTransitID: strfmt.UUID(sit.ID.String()),
	}

	handler = ReleaseStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)

	suite.Assertions.IsType(&sitop.ReleaseStorageInTransitForbidden{}, response)

}

func (suite *HandlerSuite) TestPatchStorageInTransitHandler() {
	shipment, sit, user := setupStorageInTransitHandlerTest(suite)

	sit.WarehouseID = "123456"
	sit.Notes = swag.String("Updated Note")
	sit.WarehouseEmail = swag.String("updated@email.com")

	sitPayload := payloadForStorageInTransitModel(&sit)

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s", shipment.ID.String(), sit.ID.String())
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateOfficeRequest(req, user)

	params := sitop.PatchStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipment.ID.String()),
		StorageInTransitID: strfmt.UUID(sit.ID.String()),
		StorageInTransit:   sitPayload,
	}

	handler := PatchStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.PatchStorageInTransitOK{}, response)

	responsePayload := response.(*sitop.PatchStorageInTransitOK).Payload

	storageInTransitPayloadCompare(suite, sitPayload, responsePayload)

	// Alright now let's make sure we fail out with a bad user
	serviceMemberUser := testdatagen.MakeDefaultServiceMember(suite.DB())
	req = httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateRequest(req, serviceMemberUser)
	handler = PatchStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	handler.Handle(params)
	suite.Error(models.ErrFetchForbidden)

	// Let's also make sure it fails for a TSP user that doesn't have permissions
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	req = httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	params.HTTPRequest = req
	handler = PatchStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	handler.Handle(params)
	suite.Error(models.ErrFetchForbidden)

	// Lastly let's make sure it succeeds when the tsp does have permissions
	// Now let's do a working one
	assertions := testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			ShipmentID:                      shipment.ID,
		},
	}
	// Create a shipment offer that uses our generated TSP ID and shipment ID so that our TSP has rights to
	// Use these to update a SIT for them.
	testdatagen.MakeShipmentOffer(suite.DB(), assertions)
	req = httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	params.HTTPRequest = req
	handler = PatchStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)

	suite.Assertions.IsType(&sitop.PatchStorageInTransitOK{}, response)
	responsePayload = response.(*sitop.PatchStorageInTransitOK).Payload
	storageInTransitPayloadCompare(suite, sitPayload, responsePayload)

}

func (suite *HandlerSuite) TestDeleteStorageInTransitHandler() {

	// Let's make sure that it fails if an Office user tries to do it.
	shipment, sit, officeUser := setupStorageInTransitHandlerTest(suite)

	failPath := fmt.Sprintf("/shipments/%s/storage_in_transits/%s", shipment.ID.String(), sit.ID.String())

	failReq := httptest.NewRequest("DELETE", failPath, nil)
	failReq = suite.AuthenticateOfficeRequest(failReq, officeUser)
	failParams := sitop.DeleteStorageInTransitParams{
		HTTPRequest:        failReq,
		ShipmentID:         strfmt.UUID(shipment.ID.String()),
		StorageInTransitID: strfmt.UUID(sit.ID.String()),
	}
	failHandler := DeleteStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	failHandler.Handle(failParams)
	suite.Error(models.ErrFetchForbidden)

	// Let's have one, you know, be successful
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())

	// Lastly let's make sure it succeeds when the tsp does have permissions
	// Now let's do a working one
	assertions := testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			ShipmentID:                      shipment.ID,
		},
	}
	// Create a shipment offer that uses our generated TSP ID and shipment ID so that our TSP has rights to
	// Use these to create a SIT for them.
	testdatagen.MakeShipmentOffer(suite.DB(), assertions)

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s", shipment.ID.String(), sit.ID.String())

	req := httptest.NewRequest("DELETE", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	params := sitop.DeleteStorageInTransitParams{
		HTTPRequest:        req,
		ShipmentID:         strfmt.UUID(shipment.ID.String()),
		StorageInTransitID: strfmt.UUID(sit.ID.String()),
	}
	handler := DeleteStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)
	suite.Assertions.IsType(&sitop.DeleteStorageInTransitOK{}, response)

}
