package internalapi_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *InternalAPISuite) TestUpdateMTOShipment() {
	setUpRequestBody := func(shipment models.MTOShipment) *bytes.Buffer {
		jsonBody := []byte(fmt.Sprintf(`{"customerRemarks": "hello, server!", "shipmentType": "%s"}`, shipment.ShipmentType))
		bodyBuffer := bytes.NewBuffer(jsonBody)
		return bodyBuffer
	}

	suite.Run("Unauthorized mto shipment update by another service member", func() {
		shipment := factory.BuildMTOShipment(suite.DB(), nil, nil)

		maliciousUser := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

		endpointPath := fmt.Sprintf("/internal/mto-shipments/%s", shipment.ID.String())

		body := setUpRequestBody(shipment)

		req := suite.NewAuthenticatedMilRequest("PATCH", endpointPath, body, maliciousUser)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("If-Match", etag.GenerateEtag(shipment.UpdatedAt))

		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)
		resBody, err := io.ReadAll(rr.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
		}

		fmt.Printf("resbody: %s\n", resBody)

		suite.Equal(http.StatusNotFound, rr.Code)
	})

	suite.Run("Unauthorized upload to /mto-shipments/{mtoShipmentId} by user that isn't logged in", func() {
		shipment := factory.BuildMTOShipment(suite.DB(), nil, nil)

		endpointPath := fmt.Sprintf("/internal/mto-shipments/%s", shipment.ID.String())

		body := setUpRequestBody(shipment)

		req := suite.NewMilRequest("PATCH", endpointPath, body)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("If-Match", etag.GenerateEtag(shipment.UpdatedAt))

		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)

		// Happens because we don't have a CSRF token, since they aren't logged in.
		suite.Equal(http.StatusForbidden, rr.Code)
	})

	suite.Run("Authorized patch request to /mto-shipments/{mtoShipmentId}", func() {
		move := factory.BuildSubmittedMove(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			}}, nil)

		endpointPath := fmt.Sprintf("/internal/mto-shipments/%s", shipment.ID.String())

		body := setUpRequestBody(shipment)

		req := suite.NewAuthenticatedMilRequest("PATCH", endpointPath, body, move.Orders.ServiceMember)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("If-Match", etag.GenerateEtag(shipment.UpdatedAt))

		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
	})
}
