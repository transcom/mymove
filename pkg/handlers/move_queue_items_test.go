package handlers

import (
	"fmt"

	"net/http/httptest"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/auth/context"
	queueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/queues"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *HandlerSuite) TestShowQueueHandler() {
	t := suite.T()

	// Given: An office user
	officeUser := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&officeUser)

	//  A service member and a move belonging to that service member
	smUser := models.User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "servicememeber@example.com",
	}
	suite.mustSave(&smUser)
	newServiceMember := models.ServiceMember{
		UserID:    smUser.ID,
		FirstName: swag.String("Nino"),
		LastName:  swag.String("Panino"),
	}
	suite.mustSave(&newServiceMember)

	newMove := models.Move{
		UserID: smUser.ID,
	}
	suite.mustSave(&newMove)

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/queues/some_queue", nil)
	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, officeUser.ID, "fake token")
	ctx = context.PopulateUserModel(ctx, officeUser)
	req = req.WithContext(ctx)

	params := queueop.ShowQueueParams{
		HTTPRequest: req,
		QueueType:   "new_moves",
	}
	// And: show Queue is queried
	showHandler := ShowQueueHandler(NewHandlerContext(suite.db, suite.logger))
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	okResponse := showResponse.(*queueop.ShowQueueOK)
	moveQueueItem := okResponse.Payload[0]

	// And: Returned query to include our added move
	expectedCustomerName := fmt.Sprintf("%v, %v", newServiceMember.LastName, newServiceMember.FirstName)
	if *moveQueueItem.CustomerName != expectedCustomerName {
		t.Errorf("Expected move queue item to have service member name '%v', instead has '%v'", expectedCustomerName, moveQueueItem.CustomerName)
	}

}

// func (suite *HandlerSuite) TestShowMoveWrongUser() {
// 	t := suite.T()

// 	// Given: A move with a not-logged-in user and a separate logged-in user
// 	notLoggedInUser := models.User{
// 		LoginGovUUID:  uuid.Must(uuid.NewV4()),
// 		LoginGovEmail: "email@example.com",
// 	}
// 	suite.mustSave(&notLoggedInUser)

// 	loggedInUser := models.User{
// 		LoginGovUUID:  uuid.Must(uuid.NewV4()),
// 		LoginGovEmail: "email2@example.com",
// 	}
// 	suite.mustSave(&loggedInUser)

// 	// When: A move is created for not-logged-in-user
// 	var selectedType = internalmessages.SelectedMoveTypeCOMBO
// 	newMove := models.Move{
// 		UserID:           notLoggedInUser.ID,
// 		SelectedMoveType: &selectedType,
// 	}
// 	suite.mustSave(&newMove)

// 	// And: the context contains the auth values for logged-in user
// 	req := httptest.NewRequest("GET", "/moves/some_id", nil)
// 	ctx := req.Context()
// 	ctx = context.PopulateAuthContext(ctx, loggedInUser.ID, "fake token")
// 	ctx = context.PopulateUserModel(ctx, loggedInUser)
// 	req = req.WithContext(ctx)
// 	showMoveParams := moveop.ShowMoveParams{
// 		HTTPRequest: req,
// 		MoveID:      strfmt.UUID(newMove.ID.String()),
// 	}
// 	// And: Show move is queried
// 	showHandler := ShowMoveHandler(NewHandlerContext(suite.db, suite.logger))
// 	showResponse := showHandler.Handle(showMoveParams)

// 	_, ok := showResponse.(*moveop.ShowMoveForbidden)
// 	if !ok {
// 		t.Fatalf("Request failed: %#v", showResponse)
// 	}
// }
