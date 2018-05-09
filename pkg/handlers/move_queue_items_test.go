package handlers

import (
	"fmt"

	"net/http/httptest"

	// "github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	queueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/queues"
	// "github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestShowQueueHandler() {
	t := suite.T()
	t.Skip("don't test stubbed out endpoint")

	// Given: An office user
	officeUser := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&officeUser)

	//  A set of orders and a move belonging to those orders
	order, _ := testdatagen.MakeOrder(suite.db)

	newMove := models.Move{
		OrdersID: order.ID,
	}
	suite.mustSave(&newMove)

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/queues/some_queue", nil)
	req = suite.authenticateRequest(req, officeUser)

	params := queueop.ShowQueueParams{
		HTTPRequest: req,
		QueueType:   "new",
	}
	// And: show Queue is queried
	showHandler := ShowQueueHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	okResponse := showResponse.(*queueop.ShowQueueOK)
	moveQueueItem := okResponse.Payload[0]

	// And: Returned query to include our added move
	expectedCustomerName := fmt.Sprintf("%v %v", *order.ServiceMember.FirstName, *order.ServiceMember.LastName)
	if *moveQueueItem.CustomerName != expectedCustomerName {
		t.Errorf("Expected move queue item to have service member name '%v', instead has '%v'", expectedCustomerName, *moveQueueItem.CustomerName)
	}
}
