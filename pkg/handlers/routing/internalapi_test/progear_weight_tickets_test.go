package internalapi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

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
