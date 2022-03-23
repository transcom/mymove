package ghcapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/gen/ghcmessages"

	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/queues"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/mocks"
	order "github.com/transcom/mymove/pkg/services/order"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetMoveQueuesHandler() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	hhgMoveType := models.SelectedMoveTypeHHG
	// Default Origin Duty Location GBLOC is LKNQ
	hhgMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusSUBMITTED,
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: hhgMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	// Create a move with an origin duty location outside of office user GBLOC
	excludedMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
		},
		TransportationOffice: models.TransportationOffice{
			Gbloc: "AGFM",
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: excludedMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "90210", officeUser.TransportationOffice.Gbloc)
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "50309", officeUser.TransportationOffice.Gbloc)

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetMovesQueueHandler{
		context,
		order.NewOrderFetcher(),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueOK{}, response)
	payload := response.(*queues.GetMovesQueueOK).Payload

	order := hhgMove.Orders
	result := payload.QueueMoves[0]
	deptIndicator := *result.DepartmentIndicator
	suite.Len(payload.QueueMoves, 1)
	suite.Equal(hhgMove.ID.String(), result.ID.String())
	suite.Equal(order.ServiceMember.ID.String(), result.Customer.ID.String())
	suite.Equal(*order.DepartmentIndicator, string(deptIndicator))
	suite.Equal(order.OriginDutyLocation.TransportationOffice.Gbloc, string(result.OriginGBLOC))
	suite.Equal(order.OriginDutyLocation.ID.String(), result.OriginDutyLocation.ID.String())
	suite.Equal(hhgMove.Locator, result.Locator)
	suite.Equal(int64(1), result.ShipmentsCount)
}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerMoveInfo() {
	suite.Run("displays move attributes for all move types returned by ListOrders", func() {
		stub := testdatagen.Assertions{Stub: true}

		// Stub HHG move
		hhgMove := testdatagen.MakeHHGMoveWithShipment(suite.DB(), stub)
		hhgMove.ShipmentGBLOC = append(hhgMove.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: "LKNQ"})

		// Stub HHG_PPM move
		hhgPPMMove := testdatagen.MakeHHGPPMMoveWithShipment(suite.DB(), stub)
		hhgPPMMove.ShipmentGBLOC = append(hhgPPMMove.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: "LKNQ"})

		// Stub NTS move
		ntsMove := testdatagen.MakeNTSMoveWithShipment(suite.DB(), stub)
		ntsMove.ShipmentGBLOC = append(ntsMove.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: "LKNQ"})

		// Stub NTSR move
		ntsrMove := testdatagen.MakeNTSRMoveWithShipment(suite.DB(), stub)
		ntsrMove.ShipmentGBLOC = append(ntsrMove.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: "LKNQ"})

		var expectedMoves []models.Move
		expectedMoves = append(expectedMoves, hhgMove, hhgPPMMove, ntsMove, ntsrMove)

		officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), stub)

		orderFetcher := mocks.OrderFetcher{}
		orderFetcher.On("ListOrders", mock.AnythingOfType("*appcontext.appContext"),
			officeUser.ID, mock.Anything).Return(expectedMoves, 4, nil)

		request := httptest.NewRequest("GET", "/queues/moves", nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
		}
		context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
		handler := GetMovesQueueHandler{
			context,
			&orderFetcher,
		}
		response := handler.Handle(params)
		payload := response.(*queues.GetMovesQueueOK).Payload
		moves := payload.QueueMoves

		suite.Equal(4, len(moves))
		for i := range moves {
			suite.Equal(moves[i].Locator, expectedMoves[i].Locator)
			suite.Equal(string(moves[i].Status), string(expectedMoves[i].Status))
			suite.Equal(moves[i].ShipmentsCount, int64(len(expectedMoves[i].MTOShipments)))
		}
	})
}

func (suite *HandlerSuite) TestGetMoveQueuesBranchFilter() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "90210", officeUser.TransportationOffice.Gbloc)
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "50309", officeUser.TransportationOffice.Gbloc)

	hhgMoveType := models.SelectedMoveTypeHHG

	move := models.Move{
		SelectedMoveType: &hhgMoveType,
		Status:           models.MoveStatusSUBMITTED,
	}

	shipment := models.MTOShipment{
		Status: models.MTOShipmentStatusSubmitted,
	}

	// Create an order where the service member has an ARMY affiliation (default)
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move:        move,
		MTOShipment: shipment,
	})

	// Create an order where the service member has an AIR_FORCE affiliation
	airForce := models.AffiliationAIRFORCE
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MTOShipment: shipment,
		Move:        move,
		ServiceMember: models.ServiceMember{
			Affiliation: &airForce,
		},
	})

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
		Branch:      models.StringPointer("AIR_FORCE"),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetMovesQueueHandler{
		context,
		order.NewOrderFetcher(),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueOK{}, response)
	payload := response.(*queues.GetMovesQueueOK).Payload

	result := payload.QueueMoves[0]

	suite.Equal(1, len(payload.QueueMoves))
	suite.Equal("AIR_FORCE", result.Customer.Agency)
}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerStatuses() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "90210", officeUser.TransportationOffice.Gbloc)
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "50309", officeUser.TransportationOffice.Gbloc)

	hhgMoveType := models.SelectedMoveTypeHHG
	// Default Origin Duty Location GBLOC is LKNQ
	hhgMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusSUBMITTED,
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: hhgMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	// Create a shipment on hhgMove that has Rejected status
	rejectionReason := "unnecessary"
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: hhgMove,
		MTOShipment: models.MTOShipment{
			Status:          models.MTOShipmentStatusRejected,
			RejectionReason: &rejectionReason,
		},
	})
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "06001", "AGFM")

	// Create an order with an origin duty location outside of office user GBLOC
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Name:  "Fort Punxsutawney",
			Gbloc: "AGFM",
		},
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
			PickupAddress: &models.Address{
				PostalCode: "06001",
			},
		},
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
		},
	})

	request := httptest.NewRequest("GET", "/move-task-orders/{moveTaskOrderID}", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetMovesQueueHandler{
		context,
		order.NewOrderFetcher(),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	payload := response.(*queues.GetMovesQueueOK).Payload
	suite.NoError(payload.Validate(strfmt.Default))
	result := payload.QueueMoves[0]

	suite.Equal(ghcmessages.MoveStatus("SUBMITTED"), result.Status)

	// let's test for the Move approved status
	hhgMove.Status = models.MoveStatusAPPROVED
	_, _ = suite.DB().ValidateAndSave(&hhgMove)

	response = handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueOK{}, response)
	payload = response.(*queues.GetMovesQueueOK).Payload

	result = payload.QueueMoves[0]

	suite.Equal(ghcmessages.MoveStatus("APPROVED"), result.Status)

	// Now let's test Approvals requested
	hhgMove.Status = models.MoveStatusAPPROVALSREQUESTED
	_, _ = suite.DB().ValidateAndSave(&hhgMove)

	response = handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueOK{}, response)
	payload = response.(*queues.GetMovesQueueOK).Payload

	result = payload.QueueMoves[0]

	suite.Equal(ghcmessages.MoveStatus("APPROVALS REQUESTED"), result.Status)

}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerFilters() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "90210", officeUser.TransportationOffice.Gbloc)
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "50309", officeUser.TransportationOffice.Gbloc)
	hhgMoveType := models.SelectedMoveTypeHHG
	submittedMove := models.Move{
		SelectedMoveType: &hhgMoveType,
		Status:           models.MoveStatusSUBMITTED,
	}
	submittedShipment := models.MTOShipment{
		Status: models.MTOShipmentStatusSubmitted,
	}
	airForce := models.AffiliationAIRFORCE

	// New move with AIR_FORCE service member affiliation to test branch filter
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move:        submittedMove,
		MTOShipment: submittedShipment,
		ServiceMember: models.ServiceMember{
			Affiliation: &airForce,
		},
	})

	// Approvals requested
	approvedMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusAPPROVALSREQUESTED,
		},
	})
	testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		Move: approvedMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusApproved,
		},
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusSubmitted,
		},
	})

	// Move approved
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusApproved,
		},
	})

	// Move DRAFT and CANCELLED should not be included
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusDRAFT,
		},
		MTOShipment: submittedShipment,
	})
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusCANCELED,
		},
		MTOShipment: submittedShipment,
	})

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)

	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetMovesQueueHandler{
		context,
		order.NewOrderFetcher(),
	}

	suite.Run("loads results with all STATUSes selected", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
			Status: []string{
				string(models.MoveStatusSUBMITTED),
				string(models.MoveStatusAPPROVED),
				string(models.MoveStatusAPPROVALSREQUESTED),
			},
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetMovesQueueOK).Payload
		suite.NoError(payload.Validate(strfmt.Default))
		suite.EqualValues(3, payload.TotalCount)
		suite.Len(payload.QueueMoves, 3)
		// test that the moves are sorted by status descending
		suite.Equal(ghcmessages.MoveStatus("SUBMITTED"), payload.QueueMoves[0].Status)
	})

	suite.Run("loads results with all STATUSes and 1 page selected", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
			Status: []string{
				string(models.MoveStatusSUBMITTED),
				string(models.MoveStatusAPPROVED),
				string(models.MoveStatusAPPROVALSREQUESTED),
			},
			PerPage: swag.Int64(1),
			Page:    swag.Int64(1),
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetMovesQueueOK).Payload
		suite.EqualValues(3, payload.TotalCount)
		suite.Len(payload.QueueMoves, 1)
	})

	suite.Run("loads results with one STATUS selected", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
			Status: []string{
				string(models.MoveStatusSUBMITTED),
			},
			Page:    swag.Int64(1),
			PerPage: swag.Int64(1),
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetMovesQueueOK).Payload
		suite.EqualValues(1, payload.TotalCount)
		suite.Len(payload.QueueMoves, 1)
		suite.EqualValues(string(models.MoveStatusSUBMITTED), payload.QueueMoves[0].Status)
	})

	suite.Run("Excludes draft and canceled moves when STATUS params is empty", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetMovesQueueOK).Payload
		moves := payload.QueueMoves
		var actualStatuses []string
		for _, move := range moves {
			actualStatuses = append(actualStatuses, string(move.Status))
		}
		expectedStatuses := [3]string{"SUBMITTED", "APPROVED", "APPROVALS REQUESTED"}

		suite.EqualValues(3, payload.TotalCount)
		suite.Len(payload.QueueMoves, 3)
		suite.ElementsMatch(expectedStatuses, actualStatuses)
	})

	suite.Run("1 result with status New Move and branch AIR_FORCE", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
			Status: []string{
				string(models.MoveStatusSUBMITTED),
			},
			Branch: models.StringPointer("AIR_FORCE"),
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetMovesQueueOK).Payload
		suite.EqualValues(1, payload.TotalCount)
		suite.Len(payload.QueueMoves, 1)
		suite.EqualValues(string(models.MoveStatusSUBMITTED), payload.QueueMoves[0].Status)
		suite.Equal("AIR_FORCE", payload.QueueMoves[0].Customer.Agency)
	})

	suite.Run("No results with status New Move and branch ARMY", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
			Status: []string{
				string(models.MoveStatusSUBMITTED),
			},
			Branch: models.StringPointer("ARMY"),
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetMovesQueueOK).Payload
		suite.EqualValues(0, payload.TotalCount)
		suite.Len(payload.QueueMoves, 0)
	})
}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerCustomerInfoFilters() {
	dutyLocation1 := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
		DutyLocation: models.DutyLocation{
			Name: "This Other Station",
		},
	})

	dutyLocation2 := testdatagen.MakeDefaultDutyLocation(suite.DB())

	officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{})

	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	hhgMoveType := models.SelectedMoveTypeHHG
	// Default Origin Duty Location GBLOC is LKNQ

	serviceMember1 := testdatagen.MakeServiceMember(suite.DB(), testdatagen.Assertions{
		Stub: true,
		ServiceMember: models.ServiceMember{
			FirstName: models.StringPointer("Zoya"),
			LastName:  models.StringPointer("Darvish"),
			Edipi:     models.StringPointer("11111"),
		},
	})

	serviceMember2 := testdatagen.MakeServiceMember(suite.DB(), testdatagen.Assertions{
		Stub: true,
		ServiceMember: models.ServiceMember{
			FirstName: models.StringPointer("Owen"),
			LastName:  models.StringPointer("Nance"),
			Edipi:     models.StringPointer("22222"),
		},
	})

	move1 := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusSUBMITTED,
		},
		Order: models.Order{
			OriginDutyLocation:   &dutyLocation1,
			OriginDutyLocationID: &dutyLocation1.ID,
			NewDutyLocation:      dutyLocation1,
			NewDutyLocationID:    dutyLocation1.ID,
		},
		ServiceMember: serviceMember1,
	})

	move2 := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusSUBMITTED,
		},
		Order: models.Order{
			OriginDutyLocation:   &dutyLocation2,
			OriginDutyLocationID: &dutyLocation2.ID,
			NewDutyLocation:      dutyLocation2,
			NewDutyLocationID:    dutyLocation2.ID,
		},
		ServiceMember: serviceMember2,
	})

	shipment := models.MTOShipment{
		Status: models.MTOShipmentStatusSubmitted,
	}

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move:        move1,
		MTOShipment: shipment,
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move:        move2,
		MTOShipment: shipment,
	})

	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "90210", officeUser.TransportationOffice.Gbloc)
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "50309", officeUser.TransportationOffice.Gbloc)

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)

	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetMovesQueueHandler{
		context,
		order.NewOrderFetcher(),
	}

	suite.Run("returns unfiltered results", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetMovesQueueOK).Payload

		suite.Len(payload.QueueMoves, 2)
	})

	suite.Run("returns results matching last name search term", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
			LastName:    models.StringPointer("Nan"),
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetMovesQueueOK).Payload
		result := payload.QueueMoves[0]

		suite.Len(payload.QueueMoves, 1)
		suite.Equal("Nance", result.Customer.LastName)
	})

	suite.Run("returns results matching Dod ID search term", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
			DodID:       serviceMember1.Edipi,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetMovesQueueOK).Payload
		result := payload.QueueMoves[0]

		suite.Len(payload.QueueMoves, 1)
		suite.Equal("11111", result.Customer.DodID)
	})

	suite.Run("returns results matching Move ID search term", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
			Locator:     &move1.Locator,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetMovesQueueOK).Payload
		result := payload.QueueMoves[0]

		suite.Len(payload.QueueMoves, 1)
		suite.Equal(move1.Locator, result.Locator)

	})

	suite.Run("returns results matching OriginDutyLocation name search term", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest:        request,
			OriginDutyLocation: &dutyLocation1.Name,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetMovesQueueOK).Payload
		result := payload.QueueMoves[0]

		suite.Len(payload.QueueMoves, 1)
		suite.Equal("This Other Station", result.OriginDutyLocation.Name)
	})

	suite.Run("returns results with multiple filters applied", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest:        request,
			LastName:           models.StringPointer("Dar"),
			DodID:              serviceMember1.Edipi,
			Locator:            &move1.Locator,
			OriginDutyLocation: &dutyLocation1.Name,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetMovesQueueOK).Payload

		suite.Len(payload.QueueMoves, 1)
	})

}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerUnauthorizedRole() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTIO,
	})

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetMovesQueueHandler{
		context,
		order.NewOrderFetcher(),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueForbidden{}, response)
}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerUnauthorizedUser() {
	serviceUser := testdatagen.MakeDefaultServiceMember(suite.DB())
	serviceUser.User.Roles = append(serviceUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeCustomer,
	})

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateRequest(request, serviceUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetMovesQueueHandler{
		context,
		order.NewOrderFetcher(),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueForbidden{}, response)
}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerEmptyResults() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	// Create an order with an origin duty location outside of office user GBLOC
	hhgMoveType := models.SelectedMoveTypeHHG
	excludedMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
		},
		TransportationOffice: models.TransportationOffice{
			Gbloc: "AGFM",
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: excludedMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetMovesQueueHandler{
		context,
		order.NewOrderFetcher(),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetMovesQueueOK{}, response)
	payload := response.(*queues.GetMovesQueueOK).Payload

	suite.Len(payload.QueueMoves, 0)
}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueHandler() {
	officeUser := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{})

	// Default Origin Duty Location GBLOC is LKNQ
	hhgMove := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

	// we need a mapping for the pickup address postal code to our user's gbloc
	testdatagen.MakePostalCodeToGBLOC(suite.DB(),
		hhgMove.MTOShipments[0].PickupAddress.PostalCode,
		officeUser.TransportationOffice.Gbloc)

	// Fake this as a day and a half in the past so floating point age values can be tested
	prevCreatedAt := time.Now().Add(time.Duration(time.Hour * -36))

	actualPaymentRequest := testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		Move: hhgMove,
		PaymentRequest: models.PaymentRequest{
			CreatedAt: prevCreatedAt,
		},
	})

	testdatagen.MakeDefaultPaymentRequest(suite.DB())

	request := httptest.NewRequest("GET", "/queues/payment-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetPaymentRequestsQueueParams{
		HTTPRequest: request,
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetPaymentRequestsQueueHandler{
		context,
		paymentrequest.NewPaymentRequestListFetcher(),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	suite.Assertions.IsType(&queues.GetPaymentRequestsQueueOK{}, response)
	payload := response.(*queues.GetPaymentRequestsQueueOK).Payload

	// unfortunately, what we return and what our swagger definition
	// says are pretty far apart
	// we don't return the associated addresses for the locations
	// and the status returned is from the query string not the
	// defined PaymentRequestStatus enum as indicated in the swagger
	// definition
	//
	// suite.NoError(payload.Validate(strfmt.Default))

	suite.Len(payload.QueuePaymentRequests, 1)

	paymentRequest := *payload.QueuePaymentRequests[0]

	suite.Equal(actualPaymentRequest.ID.String(), paymentRequest.ID.String())
	suite.Equal(actualPaymentRequest.MoveTaskOrderID.String(), paymentRequest.MoveID.String())
	suite.Equal(hhgMove.Orders.ServiceMemberID.String(), paymentRequest.Customer.ID.String())
	suite.Equal(string(paymentRequest.Status), "Payment requested")
	suite.Equal("KKFA", string(paymentRequest.OriginGBLOC))

	//createdAt := actualPaymentRequest.CreatedAt
	age := float64(2)
	deptIndicator := *paymentRequest.DepartmentIndicator

	suite.Equal(age, paymentRequest.Age)
	// TODO: Standardize time format
	//suite.Equal(createdAt.Format("2006-01-02T15:04:05.000Z07:00"), paymentRequest.SubmittedAt.String()) // swagger formats to milliseconds
	suite.Equal(hhgMove.Locator, paymentRequest.Locator)

	suite.Equal(*hhgMove.Orders.DepartmentIndicator, string(deptIndicator))
}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueSubmittedAtFilter() {
	officeUser := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{})

	outOfRangeDate, _ := time.Parse("2006-01-02", "2020-10-10")

	hhgMove1 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
	hhgMove2 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

	// we need a mapping for the pickup address postal code to our user's gbloc
	testdatagen.MakePostalCodeToGBLOC(suite.DB(),
		hhgMove1.MTOShipments[0].PickupAddress.PostalCode,
		officeUser.TransportationOffice.Gbloc)

	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			CreatedAt:       outOfRangeDate,
			MoveTaskOrderID: hhgMove1.ID,
			MoveTaskOrder:   hhgMove1,
		},
	})

	createdAtTime := time.Date(2020, 10, 29, 0, 0, 0, 0, time.UTC)
	testdatagen.MakePaymentRequest(suite.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			CreatedAt:       createdAtTime,
			MoveTaskOrderID: hhgMove2.ID,
			MoveTaskOrder:   hhgMove2,
		},
	})

	request := httptest.NewRequest("GET", "/queues/payment-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)

	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetPaymentRequestsQueueHandler{
		context,
		paymentrequest.NewPaymentRequestListFetcher(),
	}
	suite.Run("returns unfiltered results", func() {
		params := queues.GetPaymentRequestsQueueParams{
			HTTPRequest: request,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		suite.Assertions.IsType(&queues.GetPaymentRequestsQueueOK{}, response)
		payload := response.(*queues.GetPaymentRequestsQueueOK).Payload

		suite.Len(payload.QueuePaymentRequests, 2)
	})

	suite.Run("returns unfiltered paginated results", func() {
		params := queues.GetPaymentRequestsQueueParams{
			HTTPRequest: request,
			Page:        swag.Int64(1),
			PerPage:     swag.Int64(1),
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		suite.Assertions.IsType(&queues.GetPaymentRequestsQueueOK{}, response)
		payload := response.(*queues.GetPaymentRequestsQueueOK).Payload

		suite.Len(payload.QueuePaymentRequests, 1)
		// Total count is more than the perPage
		suite.Equal(int64(2), payload.TotalCount)
	})

	suite.Run("returns results matching SubmittedAt date", func() {
		submittedAtDate := strfmt.DateTime(time.Date(2020, 10, 29, 0, 0, 0, 0, time.UTC))

		params := queues.GetPaymentRequestsQueueParams{
			HTTPRequest: request,
			SubmittedAt: &submittedAtDate,
		}

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetPaymentRequestsQueueOK).Payload

		suite.Len(payload.QueuePaymentRequests, 1)
	})

}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueHandlerUnauthorizedRole() {
	officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

	request := httptest.NewRequest("GET", "/queues/payment-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetPaymentRequestsQueueParams{
		HTTPRequest: request,
		Page:        swag.Int64(1),
		PerPage:     swag.Int64(1),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetPaymentRequestsQueueHandler{
		context,
		paymentrequest.NewPaymentRequestListFetcher(),
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&queues.GetPaymentRequestsQueueForbidden{}, response)
}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueHandlerServerError() {
	officeUser := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

	paymentRequestListFetcher := mocks.PaymentRequestListFetcher{}

	paymentRequestListFetcher.On("FetchPaymentRequestList", mock.AnythingOfType("*appcontext.appContext"),
		officeUser.ID,
		mock.Anything,
		mock.Anything).Return(nil, 0, errors.New("database query error"))

	request := httptest.NewRequest("GET", "/queues/payment-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetPaymentRequestsQueueParams{
		HTTPRequest: request,
		Page:        swag.Int64(1),
		PerPage:     swag.Int64(1),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetPaymentRequestsQueueHandler{
		context,
		&paymentRequestListFetcher,
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&queues.GetPaymentRequestsQueueInternalServerError{}, response)
}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueHandlerEmptyResults() {
	officeUser := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

	paymentRequestListFetcher := mocks.PaymentRequestListFetcher{}

	paymentRequestListFetcher.On("FetchPaymentRequestList", mock.AnythingOfType("*appcontext.appContext"),
		officeUser.ID,
		mock.Anything,
		mock.Anything).Return(&models.PaymentRequests{}, 0, nil)

	request := httptest.NewRequest("GET", "/queues/payment-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetPaymentRequestsQueueParams{
		HTTPRequest: request,
		Page:        swag.Int64(1),
		PerPage:     swag.Int64(1),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	handler := GetPaymentRequestsQueueHandler{
		context,
		&paymentRequestListFetcher,
	}

	response := handler.Handle(params)

	suite.Assertions.IsType(&queues.GetPaymentRequestsQueueOK{}, response)
	payload := response.(*queues.GetPaymentRequestsQueueOK).Payload

	suite.Len(payload.QueuePaymentRequests, 0)
	suite.Equal(int64(0), payload.TotalCount)
}

type servicesCounselingSubtestData struct {
	needsCounselingMove             models.Move
	counselingCompletedMove         models.Move
	marineCorpsMove                 models.Move
	officeUser                      models.OfficeUser
	needsCounselingEarliestShipment models.MTOShipment
	counselingCompletedShipment     models.MTOShipment
	handler                         GetServicesCounselingQueueHandler
	request                         *http.Request
}

func (suite *HandlerSuite) makeServicesCounselingSubtestData() (subtestData *servicesCounselingSubtestData) {
	subtestData = &servicesCounselingSubtestData{}
	subtestData.officeUser = testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{})

	hhgMoveType := models.SelectedMoveTypeHHG
	submittedAt := time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC)
	// Default Origin Duty Location GBLOC is LKNQ
	subtestData.needsCounselingMove = testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusNeedsServiceCounseling,
			SubmittedAt:      &submittedAt,
		},
	})

	requestedPickupDate := time.Date(2021, 04, 01, 0, 0, 0, 0, time.UTC)
	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: subtestData.needsCounselingMove,
		MTOShipment: models.MTOShipment{
			RequestedPickupDate: &requestedPickupDate,
			Status:              models.MTOShipmentStatusSubmitted,
		},
	})

	// May have to create postalcodetogbolc for office user
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "50309", subtestData.officeUser.TransportationOffice.Gbloc)

	earlierRequestedPickup := requestedPickupDate.Add(-7 * 24 * time.Hour)
	subtestData.needsCounselingEarliestShipment = testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: subtestData.needsCounselingMove,
		MTOShipment: models.MTOShipment{
			RequestedPickupDate: &earlierRequestedPickup,
			Status:              models.MTOShipmentStatusSubmitted,
		},
	})

	earlierSubmittedAt := submittedAt.Add(-1 * 24 * time.Hour)
	subtestData.counselingCompletedMove = testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusServiceCounselingCompleted,
			SubmittedAt:      &earlierSubmittedAt,
		},
	})

	subtestData.counselingCompletedShipment = testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: subtestData.counselingCompletedMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})

	// Create a move with an origin location outside of office user GBLOC	testdatagen.MakePostalCodeToGBLOC(suite.DB(), "30813", "AGFM")
	dutyLocationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "Fort Gordon",
			City:           "Augusta",
			State:          "GA",
			PostalCode:     "30813",
			Country:        swag.String("United States"),
		},
	})

	originDutyLocation := testdatagen.MakeDutyLocation(suite.DB(), testdatagen.Assertions{
		DutyLocation: models.DutyLocation{
			Name:      "Fort Sam Houston",
			AddressID: dutyLocationAddress.ID,
			Address:   dutyLocationAddress,
		},
	})

	excludedGBLOCMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusNeedsServiceCounseling,
		},
		OriginDutyLocation: originDutyLocation,
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: excludedGBLOCMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
		PickupAddress: models.Address{
			PostalCode: "06001",
		},
	})

	// Create a move with an origin duty location outside of office user GBLOC
	excludedStatusMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusSUBMITTED,
		},
		Order: models.Order{
			OriginDutyLocation: &originDutyLocation,
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: excludedStatusMove,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
			PickupAddress: &models.Address{
				PostalCode: "06001",
			},
		},
	})

	marineCorpsAffiliation := models.AffiliationMARINES
	subtestData.marineCorpsMove = testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusNeedsServiceCounseling,
			SubmittedAt:      &submittedAt,
		},
		ServiceMember: models.ServiceMember{
			Affiliation: &marineCorpsAffiliation,
		},
	})

	testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: subtestData.marineCorpsMove,
		MTOShipment: models.MTOShipment{
			RequestedPickupDate: &requestedPickupDate,
			Status:              models.MTOShipmentStatusSubmitted,
		},
	})

	request := httptest.NewRequest("GET", "/queues/counseling", nil)
	subtestData.request = suite.AuthenticateOfficeRequest(request, subtestData.officeUser)
	context := handlers.NewHandlerContext(suite.DB(), suite.Logger())
	subtestData.handler = GetServicesCounselingQueueHandler{
		context,
		order.NewOrderFetcher(),
	}

	return subtestData
}

func (suite *HandlerSuite) TestGetServicesCounselingQueueHandler() {
	suite.Run("returns moves in the needs counseling status by default", func() {
		subtestData := suite.makeServicesCounselingSubtestData()

		params := queues.GetServicesCounselingQueueParams{
			HTTPRequest: subtestData.request,
			Sort:        swag.String("branch"),
			Order:       swag.String("asc"),
		}
		response := subtestData.handler.Handle(params)
		suite.IsNotErrResponse(response)

		suite.Assertions.IsType(&queues.GetServicesCounselingQueueOK{}, response)
		payload := response.(*queues.GetServicesCounselingQueueOK).Payload

		order := subtestData.needsCounselingMove.Orders
		result1 := payload.QueueMoves[0]
		result2 := payload.QueueMoves[1]

		suite.Len(payload.QueueMoves, 2)
		suite.Equal(order.ServiceMember.ID.String(), result1.Customer.ID.String())
		suite.Equal(*order.ServiceMember.Edipi, result1.Customer.DodID)
		suite.Equal(subtestData.needsCounselingMove.Locator, result1.Locator)
		suite.EqualValues(subtestData.needsCounselingMove.Status, result1.Status)
		suite.Equal(subtestData.needsCounselingEarliestShipment.RequestedPickupDate.Format(time.RFC3339Nano), (time.Time)(*result1.RequestedMoveDate).Format(time.RFC3339Nano))
		suite.Equal(subtestData.needsCounselingMove.SubmittedAt.Format(time.RFC3339Nano), (time.Time)(*result1.SubmittedAt).Format(time.RFC3339Nano))
		suite.Equal(order.ServiceMember.Affiliation.String(), result1.Customer.Agency)
		suite.Equal(order.NewDutyLocation.ID.String(), result1.DestinationDutyLocation.ID.String())

		suite.EqualValues(subtestData.needsCounselingMove.Status, result2.Status)
		suite.Equal("MARINES", result2.Customer.Agency)
	})

	suite.Run("returns moves in the needs counseling and services counseling complete statuses when both filters are selected", func() {
		subtestData := suite.makeServicesCounselingSubtestData()
		params := queues.GetServicesCounselingQueueParams{
			HTTPRequest: subtestData.request,
			Status:      []string{string(models.MoveStatusNeedsServiceCounseling), string(models.MoveStatusServiceCounselingCompleted)},
		}
		response := subtestData.handler.Handle(params)
		suite.IsNotErrResponse(response)

		suite.Assertions.IsType(&queues.GetServicesCounselingQueueOK{}, response)
		payload := response.(*queues.GetServicesCounselingQueueOK).Payload

		suite.Len(payload.QueueMoves, 3)

		for _, move := range payload.QueueMoves {
			// Test that only moves with postal code in the officer user gbloc are returned
			suite.Equal("50309", *move.OriginDutyLocation.Address.PostalCode)

			// Fail if a move has a status other than the two target ones
			if models.MoveStatus(move.Status) != models.MoveStatusNeedsServiceCounseling && models.MoveStatus(move.Status) != models.MoveStatusServiceCounselingCompleted {
				suite.Fail("Test does not return moves with the correct statuses.")
			}
		}
	})

	suite.Run("responds with forbidden error when user is not an office user", func() {
		subtestData := suite.makeServicesCounselingSubtestData()
		ppmOfficeUser := testdatagen.MakePPMOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

		request := httptest.NewRequest("GET", "/queues/counseling", nil)
		request = suite.AuthenticateOfficeRequest(request, ppmOfficeUser)

		params := queues.GetServicesCounselingQueueParams{
			HTTPRequest: request,
		}
		response := subtestData.handler.Handle(params)
		suite.IsNotErrResponse(response)

		suite.Assertions.IsType(&queues.GetServicesCounselingQueueForbidden{}, response)
	})
}
