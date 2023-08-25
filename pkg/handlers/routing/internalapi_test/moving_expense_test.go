package internalapi_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *InternalAPISuite) TestUploadMovingExpense() {
	suite.Run("Authorized post to /ppm-shipments/{ppmShipmentId}/moving-expenses", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		serviceMember := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
		endpointPath := fmt.Sprintf("/internal/ppm-shipments/%s/moving-expenses", ppmShipment.ID.String())

		req := suite.NewAuthenticatedMilRequest("POST", endpointPath, nil, serviceMember)
		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)
		suite.Equal(http.StatusCreated, rr.Code)
	})

	suite.Run("Unauthorized post to /ppm-shipments/{ppmShipmentId}/moving-expenses", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		endpointPath := fmt.Sprintf("/internal/ppm-shipments/%s/moving-expenses", ppmShipment.ID.String())

		req := suite.NewAuthenticatedMilRequest("POST", endpointPath, nil, serviceMember)
		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)
		suite.Equal(http.StatusNotFound, rr.Code)
	})
}

func (suite *InternalAPISuite) TestUpdateMovingExpense() {
	setUpRequestBody := func() *bytes.Buffer {
		jsonBody := []byte(`{"movingExpenseType": "PACKING_MATERIALS", "description": "sample text", "paidWithGTCC": false, "amount": 2000, "missingReceipt": false}`)
		bodyBuffer := bytes.NewBuffer(jsonBody)
		return bodyBuffer
	}

	suite.Run("Authorized patch to /ppm-shipments/{ppmShipmentId}/moving-expenses/{movingExpenseId}", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		movingExpense := factory.BuildMovingExpense(suite.DB(), []factory.Customization{
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		serviceMember := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
		endpointPath := fmt.Sprintf("/internal/ppm-shipments/%s/moving-expenses/%s", ppmShipment.ID.String(), movingExpense.ID.String())

		body := setUpRequestBody()

		req := suite.NewAuthenticatedMilRequest("PATCH", endpointPath, body, serviceMember)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("If-Match", etag.GenerateEtag(movingExpense.UpdatedAt))

		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)
		suite.Equal(http.StatusOK, rr.Code)
	})

	suite.Run("Unauthorized patch to /ppm-shipments/{ppmShipmentId}/moving-expenses/{movingExpenseId}", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		movingExpense := factory.BuildMovingExpense(suite.DB(), []factory.Customization{
			{
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		endpointPath := fmt.Sprintf("/internal/ppm-shipments/%s/moving-expenses/%s", ppmShipment.ID.String(), movingExpense.ID.String())

		body := setUpRequestBody()

		req := suite.NewAuthenticatedMilRequest("PATCH", endpointPath, body, serviceMember)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("If-Match", etag.GenerateEtag(movingExpense.UpdatedAt))

		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)
		suite.Equal(http.StatusNotFound, rr.Code)
	})
}
