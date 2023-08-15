package internalapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
)

func (suite *InternalAPISuite) TestSubmitPPMShipmentDocumentation() {
	// setUpRequestBody sets up the request body for the ppm document submission request.
	setUpRequestBody := func() *bytes.Buffer {
		body := &internalmessages.SavePPMShipmentSignedCertification{
			CertificationText: handlers.FmtString("I accept all the liability!"),
			Signature:         handlers.FmtString("Best Customer"),
			Date:              handlers.FmtDate(time.Now()),
		}

		jsonBody, err := json.Marshal(body)

		suite.FatalNoError(err)

		bodyBuffer := bytes.NewBuffer(jsonBody)

		return bodyBuffer
	}

	suite.Run("Unauthorized call to /ppm-shipments/{ppmShipmentId}/submit-ppm-shipment-documentation by another service member", func() {
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, factory.GetTraitActiveServiceMemberUser())

		maliciousUser := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

		endpointPath := fmt.Sprintf("/internal/ppm-shipments/%s/submit-ppm-shipment-documentation", ppmShipment.ID.String())

		body := setUpRequestBody()

		req := suite.NewAuthenticatedMilRequest(http.MethodPost, endpointPath, body, maliciousUser)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusNotFound, rr.Code)
	})

	suite.Run("Unauthorized call to /ppm-shipments/{ppmShipmentId}/submit-ppm-shipment-documentation by user that isn't logged in", func() {
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, factory.GetTraitActiveServiceMemberUser())

		endpointPath := fmt.Sprintf("/internal/ppm-shipments/%s/submit-ppm-shipment-documentation", ppmShipment.ID.String())

		body := setUpRequestBody()

		req := suite.NewMilRequest(http.MethodPost, endpointPath, body)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)

		// Happens because we don't have a CSRF token, since they aren't logged in.
		suite.Equal(http.StatusForbidden, rr.Code)
	})

	suite.Run("Authorized call to /ppm-shipments/{ppmShipmentId}/submit-ppm-shipment-documentation", func() {
		ppmShipment := factory.BuildPPMShipmentReadyForFinalCustomerCloseOut(suite.DB(), nil, factory.GetTraitActiveServiceMemberUser())

		endpointPath := fmt.Sprintf("/internal/ppm-shipments/%s/submit-ppm-shipment-documentation", ppmShipment.ID.String())

		body := setUpRequestBody()

		req := suite.NewAuthenticatedMilRequest(http.MethodPost, endpointPath, body, ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember)

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
	})
}
