package internalapi_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *InternalAPISuite) TestUploadProGearWeightTicket() {
	suite.Run("Authorized post to /ppm-shipments/{ppmShipmentId}/pro-gear-weight-tickets", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		serviceMember := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
		endpointPath := fmt.Sprintf("/internal/ppm-shipments/%s/pro-gear-weight-tickets", ppmShipment.ID.String())

		req := suite.NewAuthenticatedMilRequest("POST", endpointPath, nil, serviceMember)
		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)
		suite.Equal(http.StatusCreated, rr.Code)
	})

	suite.Run("Unauthorized post to /ppm-shipments/{ppmShipmentId}/pro-gear-weight-tickets", func() {
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, nil)
		serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		endpointPath := fmt.Sprintf("/internal/ppm-shipments/%s/pro-gear-weight-tickets", ppmShipment.ID.String())

		req := suite.NewAuthenticatedMilRequest("POST", endpointPath, nil, serviceMember)
		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)
		suite.Equal(http.StatusNotFound, rr.Code)
	})
}

func (suite *InternalAPISuite) TestUpdateProgearWeightTicket() {
	setUpRequestBody := func() *bytes.Buffer {
		jsonBody := []byte(`{"description": "true"}`)
		bodyBuffer := bytes.NewBuffer(jsonBody)
		return bodyBuffer
	}

	suite.Run("Unauthorized progear weight ticket update by another service member", func() {
		progearWeightTicket := factory.BuildProgearWeightTicket(suite.DB(), nil, nil)

		ppmShipment := factory.BuildPPMShipment(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		maliciousUser := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		endpointPath := fmt.Sprintf("/internal/ppm-shipments/%s/pro-gear-weight-tickets/%s", ppmShipment.ID.String(), progearWeightTicket.ID.String())

		body := setUpRequestBody()

		req := suite.NewAuthenticatedMilRequest("PATCH", endpointPath, body, maliciousUser)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("If-Match", etag.GenerateEtag(progearWeightTicket.UpdatedAt))

		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusNotFound, rr.Code)
	})

	suite.Run("Unauthorized progear weight ticket update by user that isn't logged in", func() {
		progearWeightTicket := factory.BuildProgearWeightTicket(suite.DB(), nil, nil)

		ppmShipment := factory.BuildPPMShipment(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)

		endpointPath := fmt.Sprintf("/internal/ppm-shipments/%s/pro-gear-weight-tickets/%s", ppmShipment.ID.String(), progearWeightTicket.ID.String())

		body := setUpRequestBody()

		req := suite.NewMilRequest("PATCH", endpointPath, body)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("If-Match", etag.GenerateEtag(progearWeightTicket.UpdatedAt))

		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)

		// Happens because we don't have a CSRF token, since they aren't logged in.
		suite.Equal(http.StatusForbidden, rr.Code)
	})

	suite.Run("Authorized patch request to update progear weight ticket", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		document := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)
		move := factory.BuildSubmittedMove(suite.DB(), []factory.Customization{
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
		}, nil)

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			}, {
				Model:    serviceMember,
				LinkOnly: true,
			}, {
				Model:    shipment,
				LinkOnly: true,
			},
		}, nil)

		progearWeightTicket := factory.BuildProgearWeightTicket(suite.DB(), []factory.Customization{
			{
				Model:    document,
				LinkOnly: true,
			}, {
				Model:    ppmShipment,
				LinkOnly: true,
			},
		}, nil)

		endpointPath := fmt.Sprintf("/internal/ppm-shipments/%s/pro-gear-weight-tickets/%s", progearWeightTicket.PPMShipmentID.String(), progearWeightTicket.ID.String())

		body := setUpRequestBody()
		req := suite.NewAuthenticatedMilRequest("PATCH", endpointPath, body, progearWeightTicket.Document.ServiceMember)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("If-Match", etag.GenerateEtag(progearWeightTicket.UpdatedAt))

		rr := httptest.NewRecorder()

		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
	})
}
