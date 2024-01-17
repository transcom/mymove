package ghcapi

import (
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	ppmcloseoutops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	ppmcloseout "github.com/transcom/mymove/pkg/services/ppm_closeout"
)

func (suite *HandlerSuite) TestGetPPMCloseoutHandler() {
	// Success integration test
	suite.Run("Successful fetch (integration) test", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		handlerConfig := suite.HandlerConfig()
		fetcher := ppmcloseout.NewPPMCloseoutFetcher(suite.HandlerConfig().DTODPlanner())
		request := httptest.NewRequest("GET", fmt.Sprintf("/ppm-shipments/%s/closeout", ppmShipment.ID.String()), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := ppmcloseoutops.GetPPMCloseoutParams{
			HTTPRequest:   request,
			PpmShipmentID: strfmt.UUID(ppmShipment.ID.String()),
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

	// 404 response
	suite.Run("404 response when the service returns not found", func() {
		uuidForShipment, _ := uuid.NewV4()
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		handlerConfig := suite.HandlerConfig()
		fetcher := ppmcloseout.NewPPMCloseoutFetcher(suite.HandlerConfig().DTODPlanner())
		request := httptest.NewRequest("GET", fmt.Sprintf("/ppm-shipments/%s/closeout", uuidForShipment.String()), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := ppmcloseoutops.GetPPMCloseoutParams{
			HTTPRequest:   request,
			PpmShipmentID: strfmt.UUID(uuidForShipment.String()),
		}

		handler := GetPPMCloseoutHandler{
			HandlerConfig:      handlerConfig,
			ppmCloseoutFetcher: fetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&ppmcloseoutops.GetPPMCloseoutNotFound{}, response)
		payload := response.(*ppmcloseoutops.GetPPMCloseoutNotFound).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}
