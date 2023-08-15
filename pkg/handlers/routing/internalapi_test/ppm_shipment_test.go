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
		// Not sure if I like or dislike the use of internalmessages here. On the one hand, it makes it easier for
		//  people that are used to the way we set up request bodies in the existing integration tests. On the other,
		//  it ties us to the go-swagger types and imports. Alternatively, we could do something like:
		//
		//	body := map[string]string{
		//		"certification_text": "I accept all the liability!",
		//		"signature":          "Best Customer",
		//		"date":               "2023-08-08",
		//	}
		//
		//  and use json.Marshal(body) the same way. I think either option is better ergonomically than the way we did
		//  the first ones with byte slices like:
		//
		//  jsonBody := []byte(`{"certification_text": "I accept all the liability!", "signature": "Best Customer", "date": "2023-08-08"}`)
		//
		//  but at the time I hadn't known we could do it either of these other two ways.
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
