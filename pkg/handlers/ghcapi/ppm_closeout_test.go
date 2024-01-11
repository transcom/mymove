package ghcapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/factory"
	ppmcloseoutops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	ppmcloseout "github.com/transcom/mymove/pkg/services/ppm_closeout"
)

func (suite *HandlerSuite) TestGetPPMCloseoutHandler() {
	// Success integration test
	suite.Run("Successful fetch (integration) test", func() {
		shipment := factory.BuildPPMShipmentWithAllDocTypesApproved(suite.DB(), nil)
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		handlerConfig := suite.HandlerConfig()
		fetcher := ppmcloseout.NewPPMCloseoutFetcher()
		request := httptest.NewRequest("GET", fmt.Sprintf("/ppm-shipments/%s/closeout", shipment.ID.String()), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := ppmcloseoutops.GetPPMCloseoutParams{
			HTTPRequest:   request,
			PpmShipmentID: strfmt.UUID(shipment.ID.String()),
		}

		handler := GetPPMCloseoutHandler{
			HandlerConfig:      handlerConfig,
			ppmCloseoutFetcher: fetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&ppmcloseoutops.GetPPMCloseoutOK{}, response)
		payload := response.(*ppmcloseoutops.GetPPMCloseoutOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}
