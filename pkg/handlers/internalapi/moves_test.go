//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
//RA: in which this would be considered a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package internalapi

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"os"
	"time"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/mocks"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	moveop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/moves"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/notifications"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestPatchMoveHandler() {
	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	var newType = internalmessages.SelectedMoveTypeHHGPPM
	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}
	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}
	// And: a move is patched
	handler := PatchMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	// Then: expect a 200 status code
	suite.Assertions.IsType(&moveop.PatchMoveCreated{}, response)
	okResponse := response.(*moveop.PatchMoveCreated)

	// And: Returned query to include our added move
	suite.Assertions.Equal(&newType, okResponse.Payload.SelectedMoveType)
}

func (suite *HandlerSuite) TestPatchMoveHandlerWrongUser() {
	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.DB())
	// And: another logged in user
	anotherUser := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains a different user
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, anotherUser)

	var newType = internalmessages.SelectedMoveTypeHHGPPM
	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.CheckResponseForbidden(response)
}

func (suite *HandlerSuite) TestPatchMoveHandlerNoMove() {
	// Given: a logged in user and no Move
	user := testdatagen.MakeDefaultServiceMember(suite.DB())

	moveUUID := uuid.Must(uuid.NewV4())

	// And: the context contains a logged in user
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, user)

	var newType = internalmessages.SelectedMoveTypeHHGPPM
	patchPayload := internalmessages.PatchMovePayload{
		SelectedMoveType: &newType,
	}

	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(moveUUID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.CheckResponseNotFound(response)
}

func (suite *HandlerSuite) TestPatchMoveHandlerNoType() {
	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("PATCH", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	patchPayload := internalmessages.PatchMovePayload{}
	params := moveop.PatchMoveParams{
		HTTPRequest:      req,
		MoveID:           strfmt.UUID(move.ID.String()),
		PatchMovePayload: &patchPayload,
	}

	handler := PatchMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	response := handler.Handle(params)

	suite.Assertions.IsType(&moveop.PatchMoveCreated{}, response)
	okResponse := response.(*moveop.PatchMoveCreated)

	suite.Assertions.Equal(move.ID.String(), okResponse.Payload.ID.String())
}

func (suite *HandlerSuite) TestShowMoveHandler() {

	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.DB())

	// And: the context contains the auth values
	req := httptest.NewRequest("GET", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	params := moveop.ShowMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: show Move is queried
	showHandler := ShowMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	showResponse := showHandler.Handle(params)

	// Then: Expect a 200 status code
	suite.Assertions.IsType(&moveop.ShowMoveOK{}, showResponse)
	okResponse := showResponse.(*moveop.ShowMoveOK)

	// And: Returned query to include our added move
	suite.Assertions.Equal(move.OrdersID.String(), okResponse.Payload.OrdersID.String())

}

func (suite *HandlerSuite) TestShowMoveWrongUser() {
	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.DB())
	// And: another logged in user
	anotherUser := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains the auth values for not logged-in user
	req := httptest.NewRequest("GET", "/moves/some_id", nil)
	req = suite.AuthenticateRequest(req, anotherUser)

	showMoveParams := moveop.ShowMoveParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
	}
	// And: Show move is queried
	showHandler := ShowMoveHandler{handlers.NewHandlerContext(suite.DB(), suite.TestLogger())}
	showResponse := showHandler.Handle(showMoveParams)
	// Then: expect a forbidden response
	suite.CheckResponseForbidden(showResponse)

}

func (suite *HandlerSuite) TestSubmitMoveForApprovalHandler() {
	suite.Run("Submits ppm success", func() {
		// Given: a set of orders, a move, user and servicemember
		ppm := testdatagen.MakeDefaultPPM(suite.DB())
		move := ppm.Move

		// And: the context contains the auth values
		req := httptest.NewRequest("POST", "/moves/some_id/submit", nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
		certType := internalmessages.SignedCertificationTypeCreateSHIPMENT
		signingDate := strfmt.DateTime(time.Now())
		ppmID := strfmt.UUID(ppm.ID.String())
		certificate := internalmessages.CreateSignedCertificationPayload{
			CertificationText:        swag.String("This is your legal message"),
			CertificationType:        &certType,
			PersonallyProcuredMoveID: &ppmID,
			Date:                     &signingDate,
			Signature:                swag.String("Jane Doe"),
		}
		newSubmitMoveForApprovalPayload := internalmessages.SubmitMoveForApprovalPayload{Certificate: &certificate}

		params := moveop.SubmitMoveForApprovalParams{
			HTTPRequest:                  req,
			MoveID:                       strfmt.UUID(move.ID.String()),
			SubmitMoveForApprovalPayload: &newSubmitMoveForApprovalPayload,
		}
		// When: a move is submitted
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		context.SetNotificationSender(notifications.NewStubNotificationSender("milmovelocal", suite.TestLogger()))
		handler := SubmitMoveHandler{context, moverouter.NewMoveStatusRouter(suite.DB())}
		response := handler.Handle(params)

		// Then: expect a 200 status code
		suite.Assertions.IsType(&moveop.SubmitMoveForApprovalOK{}, response)
		okResponse := response.(*moveop.SubmitMoveForApprovalOK)
		updatedMove, err := models.FetchMoveByMoveID(suite.DB(), move.ID)
		suite.NoError(err)

		// And: Returned query to have a submitted status
		suite.Assertions.Equal(internalmessages.MoveStatusSUBMITTED, okResponse.Payload.Status)
		// And: Expect move's PPM's advance to have "Requested" status
		suite.Assertions.Equal(
			internalmessages.ReimbursementStatusREQUESTED,
			*okResponse.Payload.PersonallyProcuredMoves[0].Advance.Status)
		suite.Assertions.NotNil(okResponse.Payload.SubmittedAt)

		actualSubmittedAt := updatedMove.SubmittedAt
		currentTime := time.Now()
		diff := currentTime.Sub(*actualSubmittedAt)
		diffInSeconds := diff.Seconds()
		var oneSecond float64
		oneSecond = 1.000000

		// Test that the move was submitted within a few seconds of the current time.
		// This is better than asserting that it's not Nil, and avoids trying to mock
		// time.Now() or having to pass in a date to move.Submit just to be able to test it.
		suite.Assertions.LessOrEqual(diffInSeconds, oneSecond)
	})
	suite.Run("Submits hhg shipment success", func() {
		// Given: a set of orders, a move, user and servicemember
		hhg := testdatagen.MakeDefaultMTOShipment(suite.DB())
		move := hhg.MoveTaskOrder

		// And: the context contains the auth values
		req := httptest.NewRequest("POST", "/moves/some_id/submit", nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
		certType := internalmessages.SignedCertificationTypeCreateSHIPMENT
		signingDate := strfmt.DateTime(time.Now())
		certificate := internalmessages.CreateSignedCertificationPayload{
			CertificationText: swag.String("This is your legal message"),
			CertificationType: &certType,
			Date:              &signingDate,
			Signature:         swag.String("Jane Doe"),
		}
		newSubmitMoveForApprovalPayload := internalmessages.SubmitMoveForApprovalPayload{Certificate: &certificate}

		params := moveop.SubmitMoveForApprovalParams{
			HTTPRequest:                  req,
			MoveID:                       strfmt.UUID(move.ID.String()),
			SubmitMoveForApprovalPayload: &newSubmitMoveForApprovalPayload,
		}
		// And: a move is submitted
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		context.SetNotificationSender(notifications.NewStubNotificationSender("milmovelocal", suite.TestLogger()))
		handler := SubmitMoveHandler{context, moverouter.NewMoveStatusRouter(suite.DB())}
		response := handler.Handle(params)

		// Then: expect a 200 status code
		suite.Assertions.IsType(&moveop.SubmitMoveForApprovalOK{}, response)
		okResponse := response.(*moveop.SubmitMoveForApprovalOK)

		// And: Returned query to have a submitted status
		suite.Assertions.Equal(internalmessages.MoveStatusSUBMITTED, okResponse.Payload.Status)
	})

}

func (suite *HandlerSuite) TestSubmitMoveForServiceCounselingHandler() {
	suite.Run("Routes to service counseling when feature flag is true", func() {
		os.Setenv("FEATURE_FLAG_SERVICE_COUNSELING", "true")
		// Given: a set of orders, a move, user and servicemember
		move := testdatagen.MakeDefaultMove(suite.DB())

		// And: the context contains the auth values
		req := httptest.NewRequest("POST", "/moves/some_id/submit", nil)
		req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)
		certType := internalmessages.SignedCertificationTypeCreateSHIPMENT
		signingDate := strfmt.DateTime(time.Now())
		certificate := internalmessages.CreateSignedCertificationPayload{
			CertificationText: swag.String("This is your legal message"),
			CertificationType: &certType,
			Date:              &signingDate,
			Signature:         swag.String("Jane Doe"),
		}
		newSubmitMoveForApprovalPayload := internalmessages.SubmitMoveForApprovalPayload{Certificate: &certificate}

		params := moveop.SubmitMoveForApprovalParams{
			HTTPRequest:                  req,
			MoveID:                       strfmt.UUID(move.ID.String()),
			SubmitMoveForApprovalPayload: &newSubmitMoveForApprovalPayload,
		}
		// When: a move is submitted
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		context.SetNotificationSender(notifications.NewStubNotificationSender("milmovelocal", suite.TestLogger()))
		handler := SubmitMoveHandler{context, moverouter.NewMoveStatusRouter(suite.DB())}
		response := handler.Handle(params)

		// Then: expect a 200 status code
		suite.Assertions.IsType(&moveop.SubmitMoveForApprovalOK{}, response)
		okResponse := response.(*moveop.SubmitMoveForApprovalOK)

		updatedMove, err := models.FetchMoveByMoveID(suite.DB(), move.ID)
		suite.NoError(err)

		actualSubmittedAt := updatedMove.SubmittedAt
		currentTime := time.Now()
		diff := currentTime.Sub(*actualSubmittedAt)
		diffInSeconds := diff.Seconds()
		var oneSecond float64
		oneSecond = 1.000000

		// Test that the move was submitted within a few seconds of the current time.
		// This is better than asserting that it's not Nil, and avoids trying to mock
		// time.Now() or having to pass in a date to move.SendToServiceCounseling just
		// to be able to test it.
		suite.Assertions.LessOrEqual(diffInSeconds, oneSecond)

		suite.Equal(models.MoveStatusNeedsServiceCounseling, updatedMove.Status)
		// And: Returned query to have a needs service counseling status
		suite.Equal(internalmessages.MoveStatusNEEDSSERVICECOUNSELING, okResponse.Payload.Status)
		suite.NotNil(okResponse.Payload.SubmittedAt)
	})
}

func (suite *HandlerSuite) TestShowMoveDatesSummaryHandler() {
	dutyStationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "Fort Gordon",
			City:           "Augusta",
			State:          "GA",
			PostalCode:     "30813",
			Country:        swag.String("United States"),
		},
	})

	dutyStation := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name:        "Fort Sam Houston",
			Affiliation: internalmessages.AffiliationARMY,
			AddressID:   dutyStationAddress.ID,
			Address:     dutyStationAddress,
		},
	})

	rank := models.ServiceMemberRankE4
	serviceMember := testdatagen.MakeServiceMember(suite.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			Rank:          &rank,
			DutyStationID: &dutyStation.ID,
			DutyStation:   dutyStation,
		},
	})

	newDutyStationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "n/a",
			City:           "San Antonio",
			State:          "TX",
			PostalCode:     "78234",
			Country:        swag.String("United States"),
		},
	})

	newDutyStation := testdatagen.MakeDutyStation(suite.DB(), testdatagen.Assertions{
		DutyStation: models.DutyStation{
			Name:        "Fort Gordon",
			Affiliation: internalmessages.AffiliationARMY,
			AddressID:   newDutyStationAddress.ID,
			Address:     newDutyStationAddress,
		},
	})

	move := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID:  serviceMember.ID,
			ServiceMember:    serviceMember,
			ReportByDate:     time.Date(2018, 10, 31, 0, 0, 0, 0, time.UTC),
			NewDutyStationID: newDutyStation.ID,
			NewDutyStation:   newDutyStation,
			HasDependents:    true,
			SpouseHasProGear: true,
		},
	})

	path := fmt.Sprintf("/moves/%s/move_dates", move.ID.String())
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	moveID := strfmt.UUID(move.ID.String())
	moveDate := strfmt.Date(time.Date(2018, 10, 10, 0, 0, 0, 0, time.UTC))
	params := moveop.ShowMoveDatesSummaryParams{
		HTTPRequest: req,
		MoveID:      moveID,
		MoveDate:    moveDate,
	}
	planner := &mocks.Planner{}
	planner.On("TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(1125, nil)

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetPlanner(planner)

	showHandler := ShowMoveDatesSummaryHandler{context}
	response := showHandler.Handle(params)

	suite.IsType(&moveop.ShowMoveDatesSummaryOK{}, response)
	okResponse := response.(*moveop.ShowMoveDatesSummaryOK)

	id := move.ID.String() + ":" + moveDate.String()
	suite.Equal(id, *okResponse.Payload.ID)
	suite.Equal(moveID, *okResponse.Payload.MoveID)
	suite.Equal(moveDate, *okResponse.Payload.MoveDate)

	pack := []strfmt.Date{
		strfmt.Date(time.Date(2018, 10, 5, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 9, 0, 0, 0, 0, time.UTC)),
	}
	suite.Equal(pack, okResponse.Payload.Pack)

	pickup := []strfmt.Date{
		strfmt.Date(time.Date(2018, 10, 10, 0, 0, 0, 0, time.UTC)),
	}
	suite.Equal(pickup, okResponse.Payload.Pickup)

	transit := []strfmt.Date{
		strfmt.Date(time.Date(2018, 10, 11, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 12, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 13, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 14, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 15, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 16, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 17, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 18, 0, 0, 0, 0, time.UTC)),
		strfmt.Date(time.Date(2018, 10, 19, 0, 0, 0, 0, time.UTC)),
	}
	suite.Equal(transit, okResponse.Payload.Transit)

	delivery := []strfmt.Date{
		strfmt.Date(time.Date(2018, 10, 22, 0, 0, 0, 0, time.UTC)),
	}
	suite.Equal(delivery, okResponse.Payload.Delivery)

	report := []strfmt.Date{
		strfmt.Date(move.Orders.ReportByDate),
	}
	suite.Equal(report, okResponse.Payload.Report)
}

func (suite *HandlerSuite) TestShowMoveDatesSummaryForbiddenUser() {
	// Given: a set of orders, a move, user and servicemember
	move := testdatagen.MakeDefaultMove(suite.DB())
	// And: another logged in user
	anotherUser := testdatagen.MakeDefaultServiceMember(suite.DB())

	// And: the context contains the auth values for not logged-in user
	req := httptest.NewRequest("GET", "/moves/some_id/", nil)
	req = suite.AuthenticateRequest(req, anotherUser)

	moveDate := strfmt.Date(time.Date(2018, 10, 10, 0, 0, 0, 0, time.UTC))
	params := moveop.ShowMoveDatesSummaryParams{
		HTTPRequest: req,
		MoveID:      strfmt.UUID(move.ID.String()),
		MoveDate:    moveDate,
	}
	planner := &mocks.Planner{}
	planner.On("TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(1125, nil)

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	context.SetPlanner(planner)

	showHandler := ShowMoveDatesSummaryHandler{context}
	response := showHandler.Handle(params)

	// Then: expect a forbidden response
	suite.CheckResponseForbidden(response)

}

func (suite *HandlerSuite) TestShowShipmentSummaryWorksheet() {
	testdatagen.MakeTariff400ngItemRate(suite.DB(), testdatagen.Assertions{
		Tariff400ngItemRate: models.Tariff400ngItemRate{
			Code:     "210A",
			Schedule: models.IntPointer(1),
		},
	})
	testdatagen.MakeTariff400ngItemRate(suite.DB(), testdatagen.Assertions{
		Tariff400ngItemRate: models.Tariff400ngItemRate{
			Code:     "225A",
			Schedule: models.IntPointer(1),
		},
	})
	testdatagen.MakeDefaultTariff400ngItem(suite.DB())
	testdatagen.MakeTariff400ngServiceArea(suite.DB(), testdatagen.Assertions{
		Tariff400ngServiceArea: models.Tariff400ngServiceArea{
			ServiceArea: "296",
		},
	})
	testdatagen.MakeTariff400ngServiceArea(suite.DB(), testdatagen.Assertions{
		Tariff400ngServiceArea: models.Tariff400ngServiceArea{
			ServiceArea: "208",
		},
	})
	lhr := models.Tariff400ngLinehaulRate{
		DistanceMilesLower: 1,
		DistanceMilesUpper: 10000,
		WeightLbsLower:     1,
		WeightLbsUpper:     10000,
		RateCents:          20000,
		Type:               "ConusLinehaul",
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&lhr)
	fpr := models.Tariff400ngFullPackRate{
		Schedule:           1,
		WeightLbsLower:     1,
		WeightLbsUpper:     10000,
		RateCents:          100,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&fpr)
	fupr := models.Tariff400ngFullUnpackRate{
		Schedule:           1,
		EffectiveDateLower: testdatagen.PeakRateCycleStart,
		EffectiveDateUpper: testdatagen.PeakRateCycleEnd,
	}
	suite.MustSave(&fupr)
	tdl := testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			SourceRateArea:    "US53",
			DestinationRegion: "12",
		},
	})
	testdatagen.MakeTSPPerformance(suite.DB(),
		testdatagen.Assertions{
			TransportationServiceProviderPerformance: models.TransportationServiceProviderPerformance{
				TrafficDistributionListID: tdl.ID,
			},
		})

	move := testdatagen.MakeDefaultMove(suite.DB())
	netWeight := unit.Pound(1000)
	ppm := testdatagen.MakePPM(suite.DB(), testdatagen.Assertions{
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			MoveID:                move.ID,
			ActualMoveDate:        &testdatagen.DateInsidePerformancePeriod,
			NetWeight:             &netWeight,
			PickupPostalCode:      models.StringPointer("50303"),
			DestinationPostalCode: models.StringPointer("30814"),
		},
	})
	certificationType := models.SignedCertificationTypePPMPAYMENT
	testdatagen.MakeSignedCertification(suite.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			SubmittingUserID:         move.Orders.ServiceMember.UserID,
			MoveID:                   move.ID,
			PersonallyProcuredMoveID: &ppm.ID,
			CertificationType:        &certificationType,
		},
	})

	req := httptest.NewRequest("GET", "/moves/some_id/shipment_summary_worksheet", nil)
	req = suite.AuthenticateRequest(req, move.Orders.ServiceMember)

	preparationDate := strfmt.Date(time.Date(2019, time.January, 1, 1, 1, 1, 1, time.UTC))
	params := moveop.ShowShipmentSummaryWorksheetParams{
		HTTPRequest:     req,
		MoveID:          strfmt.UUID(move.ID.String()),
		PreparationDate: preparationDate,
	}

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	planner := &mocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.Anything,
		mock.Anything,
	).Return(1044, nil)
	context.SetPlanner(planner)

	handler := ShowShipmentSummaryWorksheetHandler{context}
	response := handler.Handle(params)

	suite.Assertions.IsType(&moveop.ShowShipmentSummaryWorksheetOK{}, response)
	okResponse := response.(*moveop.ShowShipmentSummaryWorksheetOK)

	// check that the payload wasn't empty
	buf := new(bytes.Buffer)
	bytesRead, err := buf.ReadFrom(okResponse.Payload)
	suite.NoError(err)
	suite.NotZero(bytesRead)
}
