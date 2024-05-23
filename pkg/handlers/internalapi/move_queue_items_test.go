package internalapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	queueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/queues"
	"github.com/transcom/mymove/pkg/models/roles"
)

var statusToQueueMap = map[string]string{
	"SUBMITTED":         "new",
	"APPROVED":          "ppm_approved",
	"PAYMENT_REQUESTED": "ppm_payment_requested",
	"COMPLETED":         "ppm_completed",
}

func (suite *HandlerSuite) TestShowQueueHandlerForbidden() {
	for _, queueType := range statusToQueueMap {

		// Given: A non-office user
		user := factory.BuildServiceMember(suite.DB(), nil, nil)

		// And: the context contains the auth values
		path := "/queues/" + queueType
		req := httptest.NewRequest("GET", path, nil)
		req = suite.AuthenticateRequest(req, user)

		params := queueop.ShowQueueParams{
			HTTPRequest: req,
			QueueType:   queueType,
		}

		// And: show Queue is queried
		showHandler := ShowQueueHandler{suite.HandlerConfig()}
		showResponse := showHandler.Handle(params)

		// Then: Expect a 403 status code
		suite.Assertions.IsType(&queueop.ShowQueueForbidden{}, showResponse)
	}
}

func (suite *HandlerSuite) TestShowQueueHandlerNotFound() {

	// Given: An office user
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

	// And: the context contains the auth values
	queueType := "queue_not_found"
	path := "/queues/" + queueType
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)

	params := queueop.ShowQueueParams{
		HTTPRequest: req,
		QueueType:   queueType,
	}
	// And: show Queue is queried
	showHandler := ShowQueueHandler{suite.HandlerConfig()}
	showResponse := showHandler.Handle(params)

	// Then: Expect a 404 status code
	suite.CheckResponseNotFound(showResponse)
}
