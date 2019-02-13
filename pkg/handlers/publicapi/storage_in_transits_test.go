package publicapi

import (
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
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
			Location:   models.StorageInTransitLocationORIGIN,
			ShipmentID: shipment.ID,
		},
	}
	testdatagen.MakeStorageInTransit(suite.DB(), assertions)
	assertions.StorageInTransit.Location = models.StorageInTransitLocationDESTINATION
	sit = testdatagen.MakeStorageInTransit(suite.DB(), assertions)

	return shipment, sit, user
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

	//TODO: Add some value assertions
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

	//TODO: Add some value assertions

}
