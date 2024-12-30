// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package supportapi

import (
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	portlocation "github.com/transcom/mymove/pkg/services/port_location"
	"github.com/transcom/mymove/pkg/services/query"
)

// Create a service item on a Move with Approvals Requested status
func (suite *HandlerSuite) createServiceItem() models.MTOServiceItem {
	move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
	serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	return serviceItem
}

func (suite *HandlerSuite) TestUpdateMTOServiceItemStatusHandlerApproveSuccess() {

	// TESTCASE SCENARIO
	// Under test: UpdateMTOServiceItemStatusHandler
	// Mocked:     None
	// Set up:     We create an MTO service item in the DB, then try to approve it.
	// Expected outcome:
	//             Success, MTO service item is approved

	// SETUP
	// Create a service item on a move
	mtoServiceItem := suite.createServiceItem()
	// Update the service item so that it has an existing RejectionReason
	// because we want to test that it becomes nil when the service item is
	// approved.
	reason := "should not update reason"
	mtoServiceItem.RejectionReason = &reason
	suite.MustSave(&mtoServiceItem)

	// Create a request to the endpoint
	request := httptest.NewRequest("PATCH", "/mto-service-items/{mtoServiceItemID}/status", nil)

	requestPayload := &supportmessages.UpdateMTOServiceItemStatus{
		Status: supportmessages.MTOServiceItemStatusAPPROVED,
	}
	params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
		HTTPRequest:      request,
		MtoServiceItemID: mtoServiceItem.ID.String(),
		Body:             requestPayload,
		IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
	}

	handlerConfig := suite.HandlerConfig()
	queryBuilder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
	addressCreator := address.NewAddressCreator()
	portLocationFetcher := portlocation.NewPortLocationFetcher()
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	handler := UpdateMTOServiceItemStatusHandler{handlerConfig,
		mtoserviceitem.NewMTOServiceItemUpdater(planner, queryBuilder, moveRouter, shipmentFetcher, addressCreator, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer(), ghcrateengine.NewDomesticDestinationSITDeliveryPricer(), ghcrateengine.NewDomesticOriginSITFuelSurchargePricer(), portLocationFetcher),
	}

	// CALL FUNCTION UNDER TEST
	suite.Nil(params.Body.Validate(strfmt.Default))
	response := handler.Handle(params)

	// CHECK RESULTS
	suite.IsNotErrResponse(response)
	mtoServiceItemResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK)
	mtoServiceItemPayload := mtoServiceItemResponse.Payload
	suite.Assertions.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)

	// Check the status is APPROVED
	suite.Equal(supportmessages.MTOServiceItemStatusAPPROVED, mtoServiceItemPayload.Status())
	// Check that RejectionReason was set to nil
	suite.Nil(mtoServiceItemPayload.RejectionReason())
}

func (suite *HandlerSuite) TestUpdateMTOServiceItemStatusHandlerRejectSuccess() {

	// TESTCASE SCENARIO
	// Under test: UpdateMTOServiceItemStatusHandler
	// Mocked:     None
	// Set up:     We create an MTO service item in the DB, then try to reject it.
	// Expected outcome:
	//             Success, MTO service item is rejected, rejectionReason is populated

	// SETUP
	// Create a service item on a move
	mtoServiceItem := suite.createServiceItem()

	request := httptest.NewRequest("PATCH", "/mto-service-items/{mtoServiceItemID}/status", nil)
	requestPayload := &supportmessages.UpdateMTOServiceItemStatus{
		Status:          supportmessages.MTOServiceItemStatusREJECTED,
		RejectionReason: models.StringPointer("Should definitely update the reason"),
	}
	params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
		HTTPRequest:      request,
		MtoServiceItemID: mtoServiceItem.ID.String(),
		Body:             requestPayload,
		IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
	}

	handlerConfig := suite.HandlerConfig()
	queryBuilder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
	addressCreator := address.NewAddressCreator()
	portLocationFetcher := portlocation.NewPortLocationFetcher()
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	handler := UpdateMTOServiceItemStatusHandler{handlerConfig,
		mtoserviceitem.NewMTOServiceItemUpdater(planner, queryBuilder, moveRouter, shipmentFetcher, addressCreator, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer(), ghcrateengine.NewDomesticDestinationSITDeliveryPricer(), ghcrateengine.NewDomesticOriginSITFuelSurchargePricer(), portLocationFetcher),
	}

	// CALL FUNCTION UNDER TEST
	suite.Nil(params.Body.Validate(strfmt.Default))
	response := handler.Handle(params)

	// CHECK RESULTS
	suite.IsNotErrResponse(response)
	mtoServiceItemResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK)
	mtoServiceItemPayload := mtoServiceItemResponse.Payload

	suite.Assertions.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)
	suite.Equal(supportmessages.MTOServiceItemStatusREJECTED, mtoServiceItemPayload.Status())
	suite.Equal(requestPayload.RejectionReason, mtoServiceItemPayload.RejectionReason())
}

func (suite *HandlerSuite) TestUpdateMTOServiceItemStatusHandlerRejectionFailedNoReason() {

	// TESTCASE SCENARIO
	// Under test: UpdateMTOServiceItemStatusHandler
	// Mocked:     None
	// Set up:     We create an MTO service item in the DB, then try to reject it, but fail
	//             to send a rejectionReason
	// Expected outcome:
	//             Fail, RejectionReason must be provided to service item is not updated.

	// SETUP
	// Create a service item on a move
	mtoServiceItem := suite.createServiceItem()

	request := httptest.NewRequest("PATCH", "/mto-service-items/{mtoServiceItemID}/status", nil)
	requestPayload := &supportmessages.UpdateMTOServiceItemStatus{
		Status:          supportmessages.MTOServiceItemStatusREJECTED,
		RejectionReason: nil,
	}

	params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
		HTTPRequest:      request,
		MtoServiceItemID: mtoServiceItem.ID.String(),
		Body:             requestPayload,
		IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
	}

	handlerConfig := suite.HandlerConfig()
	queryBuilder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
	addressCreator := address.NewAddressCreator()
	portLocationFetcher := portlocation.NewPortLocationFetcher()
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	handler := UpdateMTOServiceItemStatusHandler{handlerConfig,
		mtoserviceitem.NewMTOServiceItemUpdater(planner, queryBuilder, moveRouter, shipmentFetcher, addressCreator, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer(), ghcrateengine.NewDomesticDestinationSITDeliveryPricer(), ghcrateengine.NewDomesticOriginSITFuelSurchargePricer(), portLocationFetcher),
	}

	// CALL FUNCTION UNDER TEST
	suite.Nil(params.Body.Validate(strfmt.Default))
	response := handler.Handle(params)

	// CHECK RESULTS
	suite.Assertions.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusUnprocessableEntity{}, response)

	// Check that the status in DB is still SUBMITTED, not APPROVED or REJECTED
	serviceItemInDB := models.MTOServiceItem{}
	suite.DB().Find(&serviceItemInDB, mtoServiceItem.ID)
	suite.Equal(models.MTOServiceItemStatusSubmitted, serviceItemInDB.Status)

}
