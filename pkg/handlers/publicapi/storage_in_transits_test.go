package publicapi

import (
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	sitop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/storage_in_transits"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"net/http/httptest"
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
}

func (suite *HandlerSuite) TestCreateStorageInTransitHandler() {

	shipment, sit, user := setupStorageInTransitHandlerTest(suite)

	sit.WarehouseID = "12345"

	sitPayload := payloadForStorageInTransitModel(&sit)

	path := fmt.Sprintf("/shipments/%s/storage_in_transits", shipment.ID.String())
	fmt.Println(path)
	req := httptest.NewRequest("POST", path, nil)
	req = suite.AuthenticateOfficeRequest(req, user)

	params := sitop.CreateStorageInTransitParams{
		HTTPRequest:      req,
		ShipmentID:       strfmt.UUID(shipment.ID.String()),
		StorageInTransit: sitPayload,
	}

	handler := CreateStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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

	sitPayload := payloadForStorageInTransitModel(&sit)

	path := fmt.Sprintf("/shipments/%s/storage_in_transits/%s", shipment.ID.String(), sit.ID.String())
	fmt.Println(path)
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
}

func (suite *HandlerSuite) TestDeleteStorageInTransitHandler() {
	shipment, sit, user := setupStorageInTransitHandlerTest(suite)

	path := fmt.Sprintf("/shipments/%s/storage_in_transits", shipment.ID.String())
	fmt.Println(path)
	req := httptest.NewRequest("DELETE", path, nil)
	req = suite.AuthenticateOfficeRequest(req, user)

	params := sitop.DeleteStorageInTransitParams{
		HTTPRequest:        req,
		StorageInTransitID: strfmt.UUID(sit.ID.String()),
	}

	handler := DeleteStorageInTransitHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&sitop.DeleteStorageInTransitOK{}, response)
}
