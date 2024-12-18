package ghcapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/queues"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	movelocker "github.com/transcom/mymove/pkg/services/lock_move"
	"github.com/transcom/mymove/pkg/services/mocks"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	officeusercreator "github.com/transcom/mymove/pkg/services/office_user"
	order "github.com/transcom/mymove/pkg/services/order"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
)

func (suite *HandlerSuite) TestGetMoveQueuesHandler() {
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(), []roles.RoleType{roles.RoleTypeTOO})
	factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(), []roles.RoleType{roles.RoleTypeTIO})
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	// Default Origin Duty Location GBLOC is KKFA
	hhgMove := factory.BuildSubmittedMove(suite.DB(), nil, nil)

	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    hhgMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	// Create a move with an origin duty location outside of office user GBLOC
	excludedMove := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Gbloc: "AGFM",
			},
			Type: &factory.TransportationOffices.CloseoutOffice,
		},
	}, nil)

	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    excludedMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	handlerConfig := suite.HandlerConfig()
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&queues.GetMovesQueueOK{}, response)
	payload := response.(*queues.GetMovesQueueOK).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))

	suite.Len(payload.QueueMoves[0].AvailableOfficeUsers, 1)
	suite.Equal(payload.QueueMoves[0].AvailableOfficeUsers[0].OfficeUserID.String(), officeUser.ID.String())

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

func (suite *HandlerSuite) TestGetDestinationRequestsQueuesHandler() {
	// default GBLOC is KKFA
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(), []roles.RoleType{roles.RoleTypeTOO})
	officeUser.User.Roles = append(officeUser.User.Roles, roles.Role{
		RoleType: roles.RoleTypeTOO,
	})

	postalCode := "90210"
	postalCode2 := "73064"
	factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), "90210", "KKFA")
	factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), "73064", "JEAT")

	// setting up two moves, one we will see and the other we won't
	move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVALSREQUESTED,
				Show:   models.BoolPointer(true),
			},
		}}, nil)

	destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{PostalCode: postalCode},
		},
	}, nil)
	shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    destinationAddress,
			LinkOnly: true,
		},
	}, nil)

	// destination service item in SUBMITTED status
	factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    shipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusSubmitted,
			},
		},
	}, nil)

	move2 := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVALSREQUESTED,
				Show:   models.BoolPointer(true),
			},
		}}, nil)

	destinationAddress2 := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{PostalCode: postalCode2},
		},
	}, nil)
	shipment2 := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
		{
			Model:    move2,
			LinkOnly: true,
		},
		{
			Model:    destinationAddress2,
			LinkOnly: true,
		},
	}, nil)

	// destination shuttle
	factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDSHUT,
			},
		},
		{
			Model:    move2,
			LinkOnly: true,
		},
		{
			Model:    shipment2,
			LinkOnly: true,
		},
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusSubmitted,
			},
		},
	}, nil)

	request := httptest.NewRequest("GET", "/queues/destination-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetDestinationRequestsQueueParams{
		HTTPRequest: request,
	}
	handlerConfig := suite.HandlerConfig()
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetDestinationRequestsQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&queues.GetDestinationRequestsQueueOK{}, response)
	payload := response.(*queues.GetDestinationRequestsQueueOK).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))

	result := payload.QueueMoves[0]
	suite.Len(payload.QueueMoves, 1)
	suite.Equal(move.ID.String(), result.ID.String())
}

func (suite *HandlerSuite) TestListPrimeMovesHandler() {
	// Default Origin Duty Location GBLOC is KKFA
	move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

	request := httptest.NewRequest("GET", "/queues/listPrimeMoves", nil)
	params := queues.ListPrimeMovesParams{
		HTTPRequest: request,
	}
	handlerConfig := suite.HandlerConfig()
	handler := ListPrimeMovesHandler{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	movesList := response.(*queues.ListPrimeMovesOK).Payload.QueueMoves

	suite.Equal(1, len(movesList))
	suite.Equal(move.ID.String(), movesList[0].ID.String())
}

func (suite *HandlerSuite) TestGetMoveQueuesHandlerMoveInfo() {
	suite.Run("displays move attributes for all move types returned by ListOrders", func() {
		gbloc := "LKNQ"

		// Stub HHG move
		hhgMove := factory.BuildMoveWithShipment(nil, []factory.Customization{
			{
				Model: models.Move{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		hhgMove.ShipmentGBLOC = append(hhgMove.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: &gbloc})

		// Stub HHG_PPM move
		hhgPPMMove := factory.BuildMoveWithShipment(nil, []factory.Customization{
			{
				Model: models.Move{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		hhgPPMMove.ShipmentGBLOC = append(hhgPPMMove.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: &gbloc})

		// Stub NTS move
		ntsMove := factory.BuildMoveWithShipment(nil, []factory.Customization{
			{
				Model: models.Move{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHGIntoNTS,
				},
			},
		}, nil)
		ntsMove.ShipmentGBLOC = append(ntsMove.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: &gbloc})

		// Stub NTSR move
		ntsrMove := factory.BuildMoveWithShipment(nil, []factory.Customization{
			{
				Model: models.Move{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHGOutOfNTSDom,
				},
			},
		}, nil)
		ntsrMove.ShipmentGBLOC = append(ntsrMove.ShipmentGBLOC, models.MoveToGBLOC{GBLOC: &gbloc})

		var expectedMoves []models.Move
		expectedMoves = append(expectedMoves, hhgMove, hhgPPMMove, ntsMove, ntsrMove)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		orderFetcher := mocks.OrderFetcher{}
		orderFetcher.On("ListOrders", mock.AnythingOfType("*appcontext.appContext"),
			officeUser.ID, roles.RoleTypeTOO, mock.Anything).Return(expectedMoves, 4, nil)

		request := httptest.NewRequest("GET", "/queues/moves", nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
		}
		handlerConfig := suite.HandlerConfig()
		mockUnlocker := movelocker.NewMoveUnlocker()
		handler := GetMovesQueueHandler{
			handlerConfig,
			&orderFetcher,
			mockUnlocker,
			officeusercreator.NewOfficeUserFetcherPop(),
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

	move := models.Move{
		Status: models.MoveStatusSUBMITTED,
	}

	shipment := models.MTOShipment{
		Status: models.MTOShipmentStatusSubmitted,
	}

	// Create an order where the service member has an ARMY affiliation (default)
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model: move,
		},
		{
			Model: shipment,
		},
	}, nil)

	// Create an order where the service member has an AIR_FORCE affiliation
	airForce := models.AffiliationAIRFORCE
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model: shipment,
		},
		{
			Model: move,
		},
		{
			Model: models.ServiceMember{
				Affiliation: &airForce,
			},
		},
	}, nil)

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
		Branch:      models.StringPointer("AIR_FORCE"),
	}
	handlerConfig := suite.HandlerConfig()
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
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

	// Default Origin Duty Location GBLOC is KKFA
	hhgMove := factory.BuildSubmittedMove(suite.DB(), nil, nil)

	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    hhgMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	// Create a shipment on hhgMove that has Rejected status
	rejectionReason := "unnecessary"
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    hhgMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:          models.MTOShipmentStatusRejected,
				RejectionReason: &rejectionReason,
			},
		},
	}, nil)
	factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), "06001", "AGFM")

	// Create an order with an origin duty location outside of office user GBLOC
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Name:  "Fort Punxsutawney",
				Gbloc: "AGFM",
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model: models.Address{
				PostalCode: "06001",
			},
			Type: &factory.Addresses.PickupAddress,
		},
	}, nil)

	request := httptest.NewRequest("GET", "/move-task-orders/{moveTaskOrderID}", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	handlerConfig := suite.HandlerConfig()
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)

	payload := response.(*queues.GetMovesQueueOK).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))

	result := payload.QueueMoves[0]

	suite.Equal(ghcmessages.MoveStatus("SUBMITTED"), result.Status)

	// let's test for the ServiceCounselingCompleted status
	hhgMove.Status = models.MoveStatusServiceCounselingCompleted
	_, _ = suite.DB().ValidateAndSave(&hhgMove)

	// Validate incoming payload: no body to validate

	response = handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&queues.GetMovesQueueOK{}, response)
	payload = response.(*queues.GetMovesQueueOK).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))

	result = payload.QueueMoves[0]

	suite.Equal(ghcmessages.MoveStatus("SERVICE COUNSELING COMPLETED"), result.Status)

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

	submittedMove := models.Move{
		Status: models.MoveStatusSUBMITTED,
	}
	submittedShipment := models.MTOShipment{
		Status: models.MTOShipmentStatusSubmitted,
	}
	airForce := models.AffiliationAIRFORCE

	// New move with AIR_FORCE service member affiliation to test branch filter
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model: submittedMove,
		},
		{
			Model: submittedShipment,
		},
		{
			Model: models.ServiceMember{
				Affiliation: &airForce,
			},
		},
	}, nil)

	// Approvals requested
	approvedMove := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)

	factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    approvedMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusSubmitted,
			},
		},
	}, nil)

	// Service Counseling Completed Move
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusServiceCounselingCompleted,
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
	}, nil)

	// Move DRAFT and CANCELLED should not be included
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusDRAFT,
			},
		},
		{
			Model: submittedShipment,
		},
	}, nil)
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusCANCELED,
			},
		},
		{
			Model: submittedShipment,
		},
	}, nil)

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)

	handlerConfig := suite.HandlerConfig()
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
	}

	suite.Run("loads results with all STATUSes selected", func() {
		params := queues.GetMovesQueueParams{
			HTTPRequest: request,
			Status: []string{
				string(models.MoveStatusSUBMITTED),
				string(models.MoveStatusAPPROVALSREQUESTED),
				string(models.MoveStatusServiceCounselingCompleted),
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
				string(models.MoveStatusAPPROVALSREQUESTED),
				string(models.MoveStatusServiceCounselingCompleted),
			},
			PerPage: models.Int64Pointer(1),
			Page:    models.Int64Pointer(1),
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
			Page:    models.Int64Pointer(1),
			PerPage: models.Int64Pointer(1),
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
		expectedStatuses := [3]string{"SUBMITTED", "APPROVALS REQUESTED", "SERVICE COUNSELING COMPLETED"}

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

	// Default Origin Duty Location GBLOC is KKFA

	serviceMember1 := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName: models.StringPointer("Zoya"),
				LastName:  models.StringPointer("Darvish"),
				Edipi:     models.StringPointer("11111"),
			},
		},
	}, nil)

	serviceMember2 := factory.BuildServiceMember(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				FirstName: models.StringPointer("Owen"),
				LastName:  models.StringPointer("Nance"),
				Edipi:     models.StringPointer("22222"),
			},
		},
	}, nil)

	move1 := factory.BuildSubmittedMove(suite.DB(), []factory.Customization{
		{
			Model:    dutyLocation1,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    dutyLocation1,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model:    serviceMember1,
			LinkOnly: true,
		},
	}, nil)

	move2 := factory.BuildSubmittedMove(suite.DB(), []factory.Customization{
		{
			Model:    dutyLocation2,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model:    dutyLocation2,
			LinkOnly: true,
			Type:     &factory.DutyLocations.NewDutyLocation,
		},
		{
			Model:    serviceMember2,
			LinkOnly: true,
		},
	}, nil)

	shipment := models.MTOShipment{
		Status: models.MTOShipmentStatusSubmitted,
	}

	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move1,
			LinkOnly: true,
		},
		{
			Model: shipment,
		},
	}, nil)

	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    move2,
			LinkOnly: true,
		},
		{
			Model: shipment,
		},
	}, nil)

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)

	handlerConfig := suite.HandlerConfig()
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
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
			HTTPRequest:  request,
			CustomerName: models.StringPointer("Nan"),
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
			Edipi:       serviceMember1.Edipi,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsNotErrResponse(response)
		payload := response.(*queues.GetMovesQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		result := payload.QueueMoves[0]

		suite.Len(payload.QueueMoves, 1)
		suite.Equal("11111", result.Customer.Edipi)
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
		var originDutyLocations []string
		originDutyLocations = append(originDutyLocations, dutyLocation1.Name)
		params := queues.GetMovesQueueParams{
			HTTPRequest:        request,
			OriginDutyLocation: originDutyLocations,
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
		var originDutyLocations []string
		originDutyLocations = append(originDutyLocations, dutyLocation1.Name)
		params := queues.GetMovesQueueParams{
			HTTPRequest:        request,
			CustomerName:       models.StringPointer("Dar"),
			Edipi:              serviceMember1.Edipi,
			Locator:            &move1.Locator,
			OriginDutyLocation: originDutyLocations,
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
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
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
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
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
	excludedMove := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model: models.TransportationOffice{
				Gbloc: "AGFM",
			},
			Type: &factory.TransportationOffices.CloseoutOffice,
		},
	}, nil)
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    excludedMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	request := httptest.NewRequest("GET", "/queues/moves", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetMovesQueueParams{
		HTTPRequest: request,
	}
	handlerConfig := suite.HandlerConfig()
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetMovesQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
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
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(), []roles.RoleType{roles.RoleTypeTIO})
	factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(), []roles.RoleType{roles.RoleTypeTOO})

	// Default Origin Duty Location GBLOC is KKFA
	hhgMove := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
	// Fake this as a day and a half in the past so floating point age values can be tested
	prevCreatedAt := time.Now().Add(time.Duration(time.Hour * -36))

	actualPaymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model:    hhgMove,
			LinkOnly: true,
		},
		{
			Model: models.PaymentRequest{
				CreatedAt: prevCreatedAt,
			},
		},
	}, nil)

	factory.BuildPaymentRequest(suite.DB(), nil, nil)

	request := httptest.NewRequest("GET", "/queues/payment-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetPaymentRequestsQueueParams{
		HTTPRequest: request,
	}
	handlerConfig := suite.HandlerConfig()
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetPaymentRequestsQueueHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestListFetcher(),
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
	}

	// Validate incoming payload: no body to validate
	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&queues.GetPaymentRequestsQueueOK{}, response)
	payload := response.(*queues.GetPaymentRequestsQueueOK).Payload

	// Validate outgoing payload
	suite.NoError(payload.Validate(strfmt.Default))

	suite.Len(payload.QueuePaymentRequests, 1)
	suite.Len(payload.QueuePaymentRequests[0].AvailableOfficeUsers, 1)
	suite.Equal(payload.QueuePaymentRequests[0].AvailableOfficeUsers[0].OfficeUserID.String(), officeUser.ID.String())

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

	hhgMove1 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
	hhgMove2 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

	factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				CreatedAt: outOfRangeDate,
			},
		},
		{
			Model:    hhgMove1,
			LinkOnly: true,
		},
	}, nil)

	createdAtTime := time.Date(2020, 10, 29, 0, 0, 0, 0, time.UTC)
	factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				CreatedAt: createdAtTime,
			},
		},
		{
			Model:    hhgMove2,
			LinkOnly: true,
		},
	}, nil)

	request := httptest.NewRequest("GET", "/queues/payment-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)

	handlerConfig := suite.HandlerConfig()
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetPaymentRequestsQueueHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestListFetcher(),
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
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
			Page:        models.Int64Pointer(1),
			PerPage:     models.Int64Pointer(1),
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
		Page:        models.Int64Pointer(1),
		PerPage:     models.Int64Pointer(1),
	}
	handlerConfig := suite.HandlerConfig()
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetPaymentRequestsQueueHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestListFetcher(),
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
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
		Page:        models.Int64Pointer(1),
		PerPage:     models.Int64Pointer(1),
	}
	handlerConfig := suite.HandlerConfig()
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetPaymentRequestsQueueHandler{
		handlerConfig,
		&paymentRequestListFetcher,
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)

	suite.IsType(&queues.GetPaymentRequestsQueueInternalServerError{}, response)
	payload := response.(*queues.GetPaymentRequestsQueueInternalServerError).Payload

	// Validate outgoing payload: nil payload
	suite.Nil(payload)
}

func (suite *HandlerSuite) TestGetPaymentRequestsQueueHandlerEmptyResults() {
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTIO})

	paymentRequestListFetcher := mocks.PaymentRequestListFetcher{}

	paymentRequestListFetcher.On("FetchPaymentRequestList", mock.AnythingOfType("*appcontext.appContext"),
		officeUser.ID,
		mock.Anything,
		mock.Anything).Return(&models.PaymentRequests{}, 0, nil)

	request := httptest.NewRequest("GET", "/queues/payment-requests", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)
	params := queues.GetPaymentRequestsQueueParams{
		HTTPRequest: request,
		Page:        models.Int64Pointer(1),
		PerPage:     models.Int64Pointer(1),
	}
	handlerConfig := suite.HandlerConfig()
	mockUnlocker := movelocker.NewMoveUnlocker()
	handler := GetPaymentRequestsQueueHandler{
		handlerConfig,
		&paymentRequestListFetcher,
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
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
	ppmNeedsCloseoutMove            models.Move
	officeUser                      models.OfficeUser
	needsCounselingEarliestShipment models.MTOShipment
	counselingCompletedShipment     models.MTOShipment
	handler                         GetServicesCounselingQueueHandler
	request                         *http.Request
}

func (suite *HandlerSuite) makeServicesCounselingSubtestData() (subtestData *servicesCounselingSubtestData) {
	subtestData = &servicesCounselingSubtestData{}
	subtestData.officeUser = factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(), []roles.RoleType{roles.RoleTypeServicesCounselor})

	submittedAt := time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC)
	// Default Origin Duty Location GBLOC is KKFA
	subtestData.needsCounselingMove = factory.BuildNeedsServiceCounselingMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				SubmittedAt: &submittedAt,
			},
		},
	}, nil)

	requestedPickupDate := time.Date(2021, 04, 01, 0, 0, 0, 0, time.UTC)
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.needsCounselingMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedPickupDate,
				Status:                models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)
	subtestData.ppmNeedsCloseoutMove = factory.BuildMoveWithPPMShipment(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				SubmittedAt:      &submittedAt,
				Status:           models.MoveStatusServiceCounselingCompleted,
				CloseoutOfficeID: &transportationOffice.ID,
			},
		},
		{
			Model: models.MTOShipment{
				RequestedPickupDate:   &requestedPickupDate,
				RequestedDeliveryDate: &requestedPickupDate,
				Status:                models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model: models.PPMShipment{
				Status: models.PPMShipmentStatusNeedsCloseout,
			},
		},
	}, nil)

	earlierRequestedPickup := requestedPickupDate.Add(-7 * 24 * time.Hour)
	subtestData.needsCounselingEarliestShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.needsCounselingMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				RequestedPickupDate:   &earlierRequestedPickup,
				RequestedDeliveryDate: &requestedPickupDate,
				Status:                models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	earlierSubmittedAt := submittedAt.Add(-1 * 24 * time.Hour)
	subtestData.counselingCompletedMove = factory.BuildServiceCounselingCompletedMove(suite.DB(), []factory.Customization{
		{
			Model: models.Move{
				SubmittedAt: &earlierSubmittedAt,
			},
		},
	}, nil)

	subtestData.counselingCompletedShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.counselingCompletedMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	// Create a move with an origin duty location outside of office user GBLOC
	dutyLocationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				StreetAddress1: "Fort Eisenhower",
				City:           "Fort Eisenhower",
				State:          "GA",
				PostalCode:     "77777",
			},
		},
	}, nil)

	// Create a custom postal code to GBLOC
	factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), dutyLocationAddress.PostalCode, "UUUU")
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
	excludedGBLOCMove := factory.BuildNeedsServiceCounselingMove(suite.DB(), []factory.Customization{
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
	}, nil)
	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    excludedGBLOCMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model: models.Address{
				PostalCode: "06001",
			},
		},
	}, nil)

	excludedStatusMove := factory.BuildSubmittedMove(suite.DB(), []factory.Customization{
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
	}, nil)

	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    excludedStatusMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model: models.Address{
				PostalCode: "06001",
			},
			Type: &factory.Addresses.PickupAddress,
		},
	}, nil)

	marineCorpsAffiliation := models.AffiliationMARINES

	subtestData.marineCorpsMove = factory.BuildNeedsServiceCounselingMove(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &marineCorpsAffiliation,
			},
		},
		{
			Model: models.Move{
				SubmittedAt: &submittedAt,
			},
		},
	}, nil)

	factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    subtestData.marineCorpsMove,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				RequestedPickupDate: &requestedPickupDate,
				Status:              models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	request := httptest.NewRequest("GET", "/queues/counseling", nil)
	subtestData.request = suite.AuthenticateOfficeRequest(request, subtestData.officeUser)
	handlerConfig := suite.HandlerConfig()
	mockUnlocker := movelocker.NewMoveUnlocker()
	subtestData.handler = GetServicesCounselingQueueHandler{
		handlerConfig,
		order.NewOrderFetcher(),
		mockUnlocker,
		officeusercreator.NewOfficeUserFetcherPop(),
	}

	return subtestData
}

func (suite *HandlerSuite) TestGetServicesCounselingQueueHandler() {
	suite.Run("returns moves in the needs counseling status by default", func() {
		subtestData := suite.makeServicesCounselingSubtestData()

		params := queues.GetServicesCounselingQueueParams{
			HTTPRequest: subtestData.request,
			Sort:        models.StringPointer("branch"),
			Order:       models.StringPointer("asc"),
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

		suite.Len(payload.QueueMoves[0].AvailableOfficeUsers, 1)
		suite.Equal(subtestData.officeUser.ID.String(), payload.QueueMoves[0].AvailableOfficeUsers[0].OfficeUserID.String())

		suite.Len(payload.QueueMoves, 2)
		suite.Equal(order.ServiceMember.ID.String(), result1.Customer.ID.String())
		suite.Equal(*order.ServiceMember.Edipi, result1.Customer.Edipi)
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

		suite.Len(payload.QueueMoves, 5)

		for _, move := range payload.QueueMoves {
			// Test that only moves with postal code in the officer user gbloc are returned
			suite.Equal("50309", *move.OriginDutyLocation.Address.PostalCode)

			// Fail if a move has a status other than the two target ones
			if models.MoveStatus(move.Status) != models.MoveStatusNeedsServiceCounseling && models.MoveStatus(move.Status) != models.MoveStatusServiceCounselingCompleted {
				suite.Fail("Test does not return moves with the correct statuses.")
			}
		}
	})

	suite.Run("returns moves in the needs closeout status when NeedsPPMCloseout is true", func() {
		subtestData := suite.makeServicesCounselingSubtestData()

		needsPpmCloseout := true
		params := queues.GetServicesCounselingQueueParams{
			HTTPRequest:      subtestData.request,
			NeedsPPMCloseout: &needsPpmCloseout,
		}

		// Validate incoming payload: no body to validate
		response := subtestData.handler.Handle(params)
		suite.IsNotErrResponse(response)
		suite.IsType(&queues.GetServicesCounselingQueueOK{}, response)
		payload := response.(*queues.GetServicesCounselingQueueOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Len(payload.QueueMoves, 1)

		for _, move := range payload.QueueMoves {
			// Fail if a ppm has a status other than needs closeout
			if models.MoveStatus(move.PpmStatus) != models.MoveStatus(models.PPMShipmentStatusNeedsCloseout) {
				suite.Fail("Test does not return moves with the correct status.")
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
