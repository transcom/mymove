package publicapi

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	accessorialop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accessorials"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetInvoiceHandler() {
	// When: There is a shipment, tsp/office user and an associated invoice
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []models.ShipmentStatus{models.ShipmentStatusDELIVERED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.DB(), numTspUsers, numShipments, numShipmentOfferSplit, status, models.SelectedMoveTypeHHG)
	suite.NoError(err)

	shipment := shipments[0]

	authedTspUser := tspUsers[0]
	unverifiedTspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Invoice tied to a shipment of authed tsp user
	invoice := testdatagen.MakeInvoice(suite.DB(), testdatagen.Assertions{
		Invoice: models.Invoice{
			ShipmentID: shipment.ID,
		},
	})
	// And: the context contains an unverified tsp user
	req := httptest.NewRequest("GET", "/accessorials", nil)
	req = suite.AuthenticateTspRequest(req, unverifiedTspUser)

	params := accessorialop.GetInvoiceParams{
		HTTPRequest: req,
		InvoiceID:   strfmt.UUID(invoice.ID.String()),
	}

	// And: get invoice is hit
	handler := GetInvoiceHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect 403 forbidden
	suite.Assertions.IsType(&handlers.ErrResponse{}, response)
	errResponse := response.(*handlers.ErrResponse)
	suite.Assertions.Equal(http.StatusForbidden, errResponse.Code)

	// When: verified TSP user is authenticated
	req = suite.AuthenticateTspRequest(req, authedTspUser)
	params.HTTPRequest = req

	// And: get invoice is hit
	handler = GetInvoiceHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)

	// Then: Invoice is returned
	suite.Nil(err)
	suite.Assertions.IsType(&accessorialop.GetInvoiceOK{}, response)
	okResponse := response.(*accessorialop.GetInvoiceOK)
	suite.Equal(strfmt.UUID(invoice.ID.String()), okResponse.Payload.ID)

	// When: Office user is authenticated
	req = suite.AuthenticateOfficeRequest(req, officeUser)
	params.HTTPRequest = req

	// And: get invoice is hit
	handler = GetInvoiceHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response = handler.Handle(params)

	// Then: Invoice is returned
	suite.Nil(err)
	suite.Assertions.IsType(&accessorialop.GetInvoiceOK{}, response)
	okResponse = response.(*accessorialop.GetInvoiceOK)
	suite.Equal(strfmt.UUID(invoice.ID.String()), okResponse.Payload.ID)
}
