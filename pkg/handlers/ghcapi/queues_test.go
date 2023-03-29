package ghcapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/queues"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/mocks"
	order "github.com/transcom/mymove/pkg/services/order"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetMoveQueuesHandler() {
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	hhgMoveType := models.SelectedMoveTypeHHG
	// Default Origin Duty Location GBLOC is KKFA
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

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	handlerConfig := suite.HandlerConfig()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&queues.GetMovesQueueOK{}, response)
	payload := response.(*queues.GetMovesQueueOK).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))

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
		gbloc := "LKNQ"
		stub := testdatagen.Assertions{Stub: true}

		// Stub HHG move
		hhgMove := testdatagen.MakeHHGMoveWithShipment(suite.DB(), stub)
		hhgMove.ShipmentGBLOC = append(hhgMove.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: &gbloc})

		// Stub HHG_PPM move
		hhgPPMMove := testdatagen.MakeHHGPPMMoveWithShipment(suite.DB(), stub)
		hhgPPMMove.ShipmentGBLOC = append(hhgPPMMove.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: &gbloc})

		// Stub NTS move
		ntsMove := testdatagen.MakeNTSMoveWithShipment(suite.DB(), stub)
		ntsMove.ShipmentGBLOC = append(ntsMove.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: &gbloc})

		// Stub NTSR move
		ntsrMove := testdatagen.MakeNTSRMoveWithShipment(suite.DB(), stub)
		ntsrMove.ShipmentGBLOC = append(ntsrMove.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: &gbloc})

		var expectedMoves []models.Move
		expectedMoves = append(expectedMoves, hhgMove, hhgPPMMove, ntsMove, ntsrMove)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		orderFetcher := mocks.OrderFetcher{}
		orderFetcher.On("ListOrders", mock.AnythingOfType("*appcontext.appContext"),
			officeUser.ID, mock.Anything).Return(expectedMoves, 4, nil)

		request := httptest.NewRequest("GET", "/queues/moves", nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
		}
		handlerConfig := suite.HandlerConfig()
		handler := GetMovesQueueHandler{
			handlerConfig,
			&orderFetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

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
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

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
	handlerConfig := suite.HandlerConfig()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&queues.GetMovesQueueOK{}, response)
	payload := response.(*queues.GetMovesQueueOK).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))

	result := payload.QueueMoves[0]

	suite.Equal(1, len(payload.QueueMoves))
	suite.Equal("AIR_FORCE", result.Customer.Agency)
}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerStatuses() {
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	hhgMoveType := models.SelectedMoveTypeHHG
	// Default Origin Duty Location GBLOC is KKFA
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
	handlerConfig := suite.HandlerConfig()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	payload := response.(*queues.GetMovesQueueOK).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))

	result := payload.QueueMoves[0]

	suite.Equal(ghcmessages.MoveStatus("SUBMITTED"), result.Status)

	// let's test for the Move approved status
	hhgMove.Status = models.MoveStatusAPPROVED
	_, _ = suite.DB().ValidateAndSave(&hhgMove)

	// Validate incoming payload: no body to validate

	response = handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&queues.GetMovesQueueOK{}, response)
	payload = response.(*queues.GetMovesQueueOK).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))

	result = payload.QueueMoves[0]

	suite.Equal(ghcmessages.MoveStatus("APPROVED"), result.Status)

	// Now let's test Approvals requested
	hhgMove.Status = models.MoveStatusAPPROVALSREQUESTED
	_, _ = suite.DB().ValidateAndSave(&hhgMove)

	// Validate incoming payload: no body to validate

	response = handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&queues.GetMovesQueueOK{}, response)
	payload = response.(*queues.GetMovesQueueOK).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))

	result = payload.QueueMoves[0]

	suite.Equal(ghcmessages.MoveStatus("APPROVALS REQUESTED"), result.Status)

}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerFilters() {
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

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

	handlerConfig := suite.HandlerConfig()
	handler := GetMovesQueueHandler{
		handlerConfig,
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

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
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

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)

		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

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

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.EqualValues(1, payload.TotalCount)
		suite.Len(payload.QueueMoves, 1)
		suite.EqualValues(string(models.MoveStatusSUBMITTED), payload.QueueMoves[0].Status)
	})

	suite.Run("Excludes draft and canceled moves when STATUS params is empty", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

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

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

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

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.EqualValues(0, payload.TotalCount)
		suite.Len(payload.QueueMoves, 0)
	})
}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerCustomerInfoFilters() {
	dutyLocation1 := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: models.DutyLocation{
				Name: "This Other Station",
			},
		},
	}, nil)

	dutyLocation2 := factory.BuildDutyLocation(suite.DB(), nil, nil)

	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	hhgMoveType := models.SelectedMoveTypeHHG
	// Default Origin Duty Location GBLOC is KKFA

	serviceMember1 := factory.BuildServiceMember(nil, []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName: models.StringPointer("Zoya"),
				LastName:  models.StringPointer("Darvish"),
				Edipi:     models.StringPointer("11111"),
			},
		},
	}, nil)

	serviceMember2 := factory.BuildServiceMember(nil, []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName: models.StringPointer("Owen"),
				LastName:  models.StringPointer("Nance"),
				Edipi:     models.StringPointer("22222"),
			},
		},
	}, nil)

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

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)

	handlerConfig := suite.HandlerConfig()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
	}

	suite.Run("returns unfiltered results", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Len(payload.QueueMoves, 2)
	})

	suite.Run("returns results matching last name search term", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
			LastName:    models.StringPointer("Nan"),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		result := payload.QueueMoves[0]

		suite.Len(payload.QueueMoves, 1)
		suite.Equal("Nance", result.Customer.LastName)
	})

	suite.Run("returns results matching Dod ID search term", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
			DodID:       serviceMember1.Edipi,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		result := payload.QueueMoves[0]

		suite.Len(payload.QueueMoves, 1)
		suite.Equal("11111", result.Customer.DodID)
	})

	suite.Run("returns results matching Move ID search term", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
			Locator:     &move1.Locator,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		result := payload.QueueMoves[0]

		suite.Len(payload.QueueMoves, 1)
		suite.Equal(move1.Locator, result.Locator)

	})

	suite.Run("returns results matching OriginDutyLocation name search term", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest:        request,
			OriginDutyLocation: &dutyLocation1.Name,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

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

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Len(payload.QueueMoves, 1)
	})

}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerUnauthorizedRole() {
	officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTIO})

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	handlerConfig := suite.HandlerConfig()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&queues.GetMovesQueueForbidden{}, response)
	payload := response.(*queues.GetMovesQueueForbidden).Payload

	// Validate outgoing payload: nil payload
	suite.Nil(payload)
}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerUnauthorizedUser() {
	serviceUser := factory.BuildServiceMember(suite.DB(), nil, nil)
	serviceUser.User.Roles = append(serviceUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeCustomer,
	})

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateRequest(request, serviceUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	handlerConfig := suite.HandlerConfig()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&queues.GetMovesQueueForbidden{}, response)
	payload := response.(*queues.GetMovesQueueForbidden).Payload

	// Validate outgoing payload: nil payload
	suite.Nil(payload)
}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerEmptyResults() {
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
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
	handlerConfig := suite.HandlerConfig()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&queues.GetMovesQueueOK{}, response)
	payload := response.(*queues.GetMovesQueueOK).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))

	suite.Len(payload.QueueMoves, 0)
}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueHandler() {
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTIO})

	// Default Origin Duty Location GBLOC is KKFA
	hhgMove := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

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
	handlerConfig := suite.HandlerConfig()
	handler := GetPaymentRequestsQueueHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestListFetcher(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&queues.GetPaymentRequestsQueueOK{}, response)
	payload := response.(*queues.GetPaymentRequestsQueueOK).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))

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
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTIO})

	outOfRangeDate, _ := time.Parse("2006-01-02", "2020-10-10")

	hhgMove1 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})
	hhgMove2 := testdatagen.MakeHHGMoveWithShipment(suite.DB(), testdatagen.Assertions{})

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

	handlerConfig := suite.HandlerConfig()
	handler := GetPaymentRequestsQueueHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestListFetcher(),
	}
	suite.Run("returns unfiltered results", func() {
		params := queues.GetPaymentRequestsQueueParams{
			HTTPRequest: request,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&queues.GetPaymentRequestsQueueOK{}, response)
		payload := response.(*queues.GetPaymentRequestsQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Len(payload.QueuePaymentRequests, 2)
	})

	suite.Run("returns unfiltered paginated results", func() {
		params := queues.GetPaymentRequestsQueueParams{
			HTTPRequest: request,
			Page:        swag.Int64(1),
			PerPage:     swag.Int64(1),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&queues.GetPaymentRequestsQueueOK{}, response)
		payload := response.(*queues.GetPaymentRequestsQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

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

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		payload := response.(*queues.GetPaymentRequestsQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Len(payload.QueuePaymentRequests, 1)
	})

}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueHandlerUnauthorizedRole() {
	officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})

	request := httptest.NewRequest("GET", "/queues/payment-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetPaymentRequestsQueueParams{
		HTTPRequest: request,
		Page:        swag.Int64(1),
		PerPage:     swag.Int64(1),
	}
	handlerConfig := suite.HandlerConfig()
	handler := GetPaymentRequestsQueueHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestListFetcher(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsType(&queues.GetPaymentRequestsQueueForbidden{}, response)
	payload := response.(*queues.GetPaymentRequestsQueueForbidden).Payload

	// Validate outgoing payload: nil payload
	suite.Nil(payload)
}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueHandlerServerError() {
	officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTIO})

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
	handlerConfig := suite.HandlerConfig()
	handler := GetPaymentRequestsQueueHandler{
		handlerConfig,
		&paymentRequestListFetcher,
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)

	suite.IsType(&queues.GetPaymentRequestsQueueInternalServerError{}, response)
	payload := response.(*queues.GetPaymentRequestsQueueInternalServerError).Payload

	// Validate outgoing payload: nil payload
	suite.Nil(payload)
}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueHandlerEmptyResults() {
	officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTIO})

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
	handlerConfig := suite.HandlerConfig()
	handler := GetPaymentRequestsQueueHandler{
		handlerConfig,
		&paymentRequestListFetcher,
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsType(&queues.GetPaymentRequestsQueueOK{}, response)
	payload := response.(*queues.GetPaymentRequestsQueueOK).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))

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
	subtestData.officeUser = factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})

	hhgMoveType := models.SelectedMoveTypeHHG
	submittedAt := time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC)
	// Default Origin Duty Location GBLOC is KKFA
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
			RequestedPickupDate:   &requestedPickupDate,
			RequestedDeliveryDate: &requestedPickupDate,
			Status:                models.MTOShipmentStatusSubmitted,
		},
	})

	earlierRequestedPickup := requestedPickupDate.Add(-7 * 24 * time.Hour)
	subtestData.needsCounselingEarliestShipment = testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: subtestData.needsCounselingMove,
		MTOShipment: models.MTOShipment{
			RequestedPickupDate:   &earlierRequestedPickup,
			RequestedDeliveryDate: &requestedPickupDate,
			Status:                models.MTOShipmentStatusSubmitted,
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

	// Create a move with an origin duty location outside of office user GBLOC
	dutyLocationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "Fort Gordon",
				City:           "Augusta",
				State:          "GA",
				PostalCode:     "77777",
				Country:        models.StringPointer("United States"),
			},
		},
	}, nil)

	// Create a custom postal code to GBLOC
	testdatagen.MakePostalCodeToGBLOC(suite.DB(), dutyLocationAddress.PostalCode, "UUUU")
	originDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: models.DutyLocation{
				Name: "Fort Sam Houston",
			},
		},
		{
			Model:    dutyLocationAddress,
			LinkOnly: true,
		},
	}, nil)

	// Create a move with an origin duty location outside of office user GBLOC
	excludedGBLOCMove := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusNeedsServiceCounseling,
		},
		Order: models.Order{
			OriginDutyLocation: &originDutyLocation,
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
	handlerConfig := suite.HandlerConfig()
	subtestData.handler = GetServicesCounselingQueueHandler{
		handlerConfig,
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

		// Validate incoming payload: no body to validate

		response := subtestData.handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&queues.GetServicesCounselingQueueOK{}, response)
		payload := response.(*queues.GetServicesCounselingQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

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

		// Validate incoming payload: no body to validate

		response := subtestData.handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&queues.GetServicesCounselingQueueOK{}, response)
		payload := response.(*queues.GetServicesCounselingQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

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
		user := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTIO})

		request := httptest.NewRequest("GET", "/queues/counseling", nil)
		request = suite.AuthenticateOfficeRequest(request, user)

		params := queues.GetServicesCounselingQueueParams{
			HTTPRequest: request,
		}

		// Validate incoming payload: no body to validate

		response := subtestData.handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&queues.GetServicesCounselingQueueForbidden{}, response)
		payload := response.(*queues.GetServicesCounselingQueueForbidden).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}
