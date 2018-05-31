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
		officeUser, _ := testdatagen.MakeOfficeUser(suite.db)

		//  A set of orders and a move belonging to those orders
		order, _ := testdatagen.MakeOrder(suite.db)

		newMove := models.Move{
			OrdersID: order.ID,
			Status:   models.MoveStatus(status),
		}
		suite.mustSave(&newMove)

		_, verrs, locErr := order.CreateNewMove(suite.db, nil)
		suite.False(verrs.HasAny(), "failed to create new move")
		suite.Nil(locErr)

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
		moveQueueItem := okResponse.Payload[0]

		// And: Returned query to include our added move
		// The moveQueueItems are produced by joining Moves, Orders and ServiceMember to each other, so we check the
		// furthest link in that chain
		expectedCustomerName := fmt.Sprintf("%v, %v", *order.ServiceMember.LastName, *order.ServiceMember.FirstName)
		suite.Equal(expectedCustomerName, *moveQueueItem.CustomerName)
	}
}
