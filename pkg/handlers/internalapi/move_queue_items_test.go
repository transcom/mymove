package internalapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	queueop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/queues"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

var statusToQueueMap = map[string]string{
	"SUBMITTED": "new",
	"APPROVED":  "ppm",
}

func (suite *HandlerSuite) TestShowQueueHandler() {
	for status, queueType := range statusToQueueMap {

		suite.DB().TruncateAll()

		// Given: An office user
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

		//  A set of orders and a move belonging to those orders
		order := testdatagen.MakeDefaultOrder(suite.DB())

		moveShow := true
		newMove := models.Move{
			OrdersID: order.ID,
			Status:   models.MoveStatus(status),
			Show:     &moveShow,
		}
		suite.MustSave(&newMove)

		// Make a PPM
		testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
			PersonallyProcuredMove: models.PersonallyProcuredMove{
				Move:   newMove,
				MoveID: newMove.ID,
			},
		})

		// And: the context contains the auth values
		path := "/queues/" + queueType
		req := httptest.NewRequest("GET", path, nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		params := queueop.ShowQueueParams{
			HTTPRequest: req,
			QueueType:   queueType,
		}
		// And: show Queue is queried
		showHandler := ShowQueueHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
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
		user := testdatagen.MakeDefaultServiceMember(suite.DB())

		// And: the context contains the auth values
		path := "/queues/" + queueType
		req := httptest.NewRequest("GET", path, nil)
		req = suite.AuthenticateRequest(req, user)

		params := queueop.ShowQueueParams{
			HTTPRequest: req,
			QueueType:   queueType,
		}
		// And: show Queue is queried
		showHandler := ShowQueueHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
		showResponse := showHandler.Handle(params)

		// Then: Expect a 403 status code
		suite.Assertions.IsType(&queueop.ShowQueueForbidden{}, showResponse)
	}
}

func (suite *HandlerSuite) TestGetMoveQueueItemsComboMoveDate() {
	suite.SetupTest()

	// Given: An office user
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Make a PPM
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{},
	})

	move := &ppm.Move
	move.Status = "SUBMITTED"
	suite.DB().Save(move)

	pickupDate := testdatagen.NextValidMoveDate

	// Make a shipment
	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &pickupDate,
			ActualPickupDate:    &pickupDate,
			Status:              models.ShipmentStatusSUBMITTED,
			Move:                ppm.Move,
			MoveID:              ppm.Move.ID,
		},
	})

	// And: the context contains the auth values
	path := "/queues/new"
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)
	params := queueop.ShowQueueParams{
		HTTPRequest: req,
		QueueType:   "new",
	}

	// And: show Queue is queried
	showHandler := ShowQueueHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	okResponse := showResponse.(*queueop.ShowQueueOK)

	suite.NotEmpty(okResponse.Payload)

	moveQueueItem := okResponse.Payload[0]

	// And: expect the moveQueueItem's move date to be the actual pickup date
	expectedMoveDate := fmt.Sprintf("%v", handlers.FmtDate(*shipment.ActualPickupDate))
	actualMoveDate := fmt.Sprintf("%v", *moveQueueItem.MoveDate)
	suite.Equal(expectedMoveDate, actualMoveDate)
}

func (suite *HandlerSuite) TestGetMoveQueueItemsComboSubmittedDate() {
	suite.SetupTest()

	// Given: An office user
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	// Make a PPM
	ppmSubmitDate := time.Now()
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			SubmitDate: &ppmSubmitDate,
		},
	})

	move := &ppm.Move
	move.Status = "SUBMITTED"
	suite.DB().Save(move)

	pickupDate := testdatagen.NextValidMoveDate

	// Make a shipment
	hhgSubmitDate := time.Now()
	testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &pickupDate,
			ActualPickupDate:    &pickupDate,
			Status:              models.ShipmentStatusSUBMITTED,
			Move:                ppm.Move,
			MoveID:              ppm.Move.ID,
			SubmitDate:          &hhgSubmitDate,
		},
	})

	testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			RequestedPickupDate: &pickupDate,
			ActualPickupDate:    &pickupDate,
			Status:              models.ShipmentStatusSUBMITTED,
			Move:                ppm.Move,
			MoveID:              ppm.Move.ID,
		},
	})

	// And: the context contains the auth values
	path := "/queues/new"
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateOfficeRequest(req, officeUser)
	params := queueop.ShowQueueParams{
		HTTPRequest: req,
		QueueType:   "new",
	}

	showHandler := ShowQueueHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	okResponse := showResponse.(*queueop.ShowQueueOK)

	suite.Equal(2, len(okResponse.Payload))

	moveQueueItem1 := okResponse.Payload[0]
	moveQueueItem2 := okResponse.Payload[1]

	resultPPMDate := *handlers.FmtDatePtrToPopPtr(moveQueueItem1.SubmittedDate)
	resultHHGDate := *handlers.FmtDatePtrToPopPtr(moveQueueItem2.SubmittedDate)

	suite.Equal(hhgSubmitDate.Format(time.UnixDate), resultHHGDate.Format(time.UnixDate))
	suite.Equal(ppmSubmitDate.Format(time.UnixDate), resultPPMDate.Format(time.UnixDate))
}
