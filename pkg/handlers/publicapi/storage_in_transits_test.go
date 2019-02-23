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

func (suite *HandlerSuite) TestPatchStorageInTransitHandler() {
	shipment, sit, user := setupStorageInTransitHandlerTest(suite)

	sit.WarehouseID = "123456"
	sit.Notes = swag.String("Updated Note")
	sit.WarehouseEmail = swag.String("updated@email.com")
	sit.Status = models.StorageInTransitStatusAPPROVED

	sitPayload := payloadForStorageInTransitModel(&sit)

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s", shipment.ID.String(), sit.ID.String())
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateOfficeRequest(req, user)

	params := sitop.PatchStorageInTransitParams{
		HTTPRequest:      req,
		ShipmentID:       strfmt.UUID(shipment.ID.String()),
		StorageInTransit: sitPayload,
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
	// Use these to create a SIT for them.
	testdatagen.MakeShipmentOffer(suite.DB(), assertions)
	req = httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateTspRequest(req, tspUser)
	params.HTTPRequest = req
	handler = PatchStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)

	suite.Assertions.IsType(&sitop.PatchStorageInTransitOK{}, response)
	responsePayload = response.(*sitop.PatchStorageInTransitOK).Payload
	storageInTransitPayloadCompare(suite, sitPayload, responsePayload)

	// Let's also double check to make sure that we didn't change the status. This shouldn't be something we're
	// doing from the patch handler.
	savedStorageInTransit, _ := models.FetchStorageInTransitByID(suite.DB(), sit.ID)
	suite.Equal(models.StorageInTransitStatusREQUESTED, savedStorageInTransit.Status)

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
