package handlers

import (
	"fmt"
	"net/http/httptest"

	queueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/queues"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

var statusToQueueMap = map[string]string{
	"SUBMITTED": "new",
	"APPROVED":  "ppm",
}

func (suite *HandlerSuite) TestShowQueueHandler() {
	for status, queueType := range statusToQueueMap {

		suite.db.TruncateAll()

		// Given: An office user
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)

		//  A set of orders and a move belonging to those orders
		order := testdatagen.MakeDefaultOrder(suite.db)

		newMove := models.Move{
			OrdersID: order.ID,
			Status:   models.MoveStatus(status),
		}
		suite.mustSave(&newMove)

		// Make a PPM
		newMove.CreatePPM(suite.db,
			nil,
			models.Int64Pointer(8000),
			models.TimePointer(testdatagen.DateInsidePeakRateCycle),
			models.StringPointer("72017"),
			models.BoolPointer(false),
			nil,
			models.StringPointer("60605"),
			models.BoolPointer(false),
			nil,
			models.StringPointer("estimate sit"),
			true,
			nil,
		)

		// And: the context contains the auth values
		path := "/queues/" + queueType
		req := httptest.NewRequest("GET", path, nil)
		req = suite.authenticateOfficeRequest(req, officeUser)

		params := queueop.ShowQueueParams{
			HTTPRequest: req,
			QueueType:   queueType,
		}
		// And: show Queue is queried
		showHandler := ShowQueueHandler(NewHandlerContext(suite.db, suite.logger))
		showResponse := showHandler.Handle(params)

		// Then: Expect a 200 status code
		okResponse := showResponse.(*queueop.ShowQueueOK)
		fmt.Printf("status: %v res: %v", status, okResponse)
		moveQueueItem := okResponse.Payload[0]

		// And: Returned query to include our added move
		// The moveQueueItems are produced by joining Moves, Orders and ServiceMember to each other, so we check the
		// furthest link in that chain
		expectedCustomerName := fmt.Sprintf("%v, %v", *order.ServiceMember.LastName, *order.ServiceMember.FirstName)
		suite.Equal(expectedCustomerName, *moveQueueItem.CustomerName)
	}
}

func (suite *HandlerSuite) TestShowQueueHandlerForbidden() {
	for _, queueType := range statusToQueueMap {

		// Given: A non-office user
		user := testdatagen.MakeDefaultServiceMember(suite.db)

		// And: the context contains the auth values
		path := "/queues/" + queueType
		req := httptest.NewRequest("GET", path, nil)
		req = suite.authenticateRequest(req, user)

		params := queueop.ShowQueueParams{
			HTTPRequest: req,
			QueueType:   queueType,
		}
		// And: show Queue is queried
		showHandler := ShowQueueHandler(NewHandlerContext(suite.db, suite.logger))
		showResponse := showHandler.Handle(params)

		// Then: Expect a 403 status code
		suite.Assertions.IsType(&queueop.ShowQueueForbidden{}, showResponse)
	}
}
