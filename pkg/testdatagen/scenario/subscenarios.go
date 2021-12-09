package scenario

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func subScenarioShipmentHHGCancelled(appCtx appcontext.AppContext, allDutyStations []models.DutyStation, originDutyStationsInGBLOC []models.DutyStation) func() {
	db := appCtx.DB()
	return func() {
		createTXO(appCtx)
		createTXOUSMC(appCtx)

		validStatuses := []models.MoveStatus{models.MoveStatusAPPROVED}
		// shipment cancelled was approved before
		approvedDate := time.Now()
		cancelledShipment := models.MTOShipment{Status: models.MTOShipmentStatusCanceled, ApprovedDate: &approvedDate}
		affiliationAirForce := models.AffiliationAIRFORCE
		ordersNumber := "Order1234"
		ordersTypeDetail := internalmessages.OrdersTypeDetailHHGPERMITTED
		tac := "1234"
		// make sure to create moves that does not go to US marines affiliation
		move := createRandomMove(appCtx, validStatuses, allDutyStations, originDutyStationsInGBLOC, true, testdatagen.Assertions{
			Order: models.Order{
				DepartmentIndicator: (*string)(&affiliationAirForce),
				OrdersNumber:        &ordersNumber,
				OrdersTypeDetail:    &ordersTypeDetail,
				TAC:                 &tac,
			},
			Move: models.Move{
				Locator: "HHGCAN",
			},
			ServiceMember: models.ServiceMember{Affiliation: &affiliationAirForce},
			MTOShipment:   cancelledShipment,
		})
		moveManagementUUID := "1130e612-94eb-49a7-973d-72f33685e551"
		testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
			ReService: models.ReService{ID: uuid.FromStringOrNil(moveManagementUUID)},
			MTOServiceItem: models.MTOServiceItem{
				MoveTaskOrderID: move.ID,
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedDate,
			},
		})
	}
}

func subScenarioPPMOfficeQueue(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) func() {
	return func() {
		createPPMOfficeUser(appCtx)

		// PPM Office Queue
		createPPMWithAdvance(appCtx, userUploader, moveRouter)
		createPPMWithNoAdvance(appCtx, userUploader, moveRouter)
		createPPMWithPaymentRequest(appCtx, userUploader, moveRouter)
		createCanceledPPM(appCtx, userUploader, moveRouter)
		createPPMReadyToRequestPayment(appCtx, userUploader, moveRouter)
	}
}

func subScenarioAdditionalPPMUsers(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) func() {
	return func() {
		// Create additional PPM users for mymove tests
		createPPMUsers(appCtx, userUploader)
	}
}

func subScenarioHHGOnboarding(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) func() {
	return func() {
		createTXO(appCtx)
		createTXOUSMC(appCtx)

		// Onboarding
		createUnsubmittedHHGMove(appCtx)
		createUnsubmittedHHGMoveMultiplePickup(appCtx)
		createUnsubmittedHHGMoveMultipleDestinations(appCtx)
		createServiceMemberWithOrdersButNoMoveType(appCtx)
		createServiceMemberWithNoUploadedOrders(appCtx)
		createSubmittedHHGMoveMultiplePickupAmendedOrders(appCtx, userUploader)
	}
}

func subScenarioHHGServicesCounseling(appCtx appcontext.AppContext, userUploader *uploader.UserUploader,
	allDutyStations []models.DutyStation, originDutyStationsInGBLOC []models.DutyStation) func() {
	return func() {
		createTXOServicesCounselor(appCtx)
		createTXOServicesUSMCCounselor(appCtx)

		// Services Counseling
		createHHGNeedsServicesCounseling(appCtx)
		createHHGNeedsServicesCounselingUSMC(appCtx, userUploader)
		createHHGNeedsServicesCounselingUSMC2(appCtx, userUploader)
		createHHGServicesCounselingCompleted(appCtx)
		createHHGNoShipments(appCtx)

		for i := 0; i < 12; i++ {
			validStatuses := []models.MoveStatus{models.MoveStatusNeedsServiceCounseling, models.MoveStatusServiceCounselingCompleted}
			createRandomMove(appCtx, validStatuses, allDutyStations, originDutyStationsInGBLOC, false, testdatagen.Assertions{
				UserUploader: userUploader,
			})
		}
	}
}

func subScenarioTXOQueues(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, logger *zap.Logger) func() {
	return func() {
		createTOO(appCtx)
		createTIO(appCtx)
		createTXO(appCtx)
		createTXOUSMC(appCtx)
		createServicesCounselor(appCtx)
		createTXOServicesCounselor(appCtx)
		createTXOServicesUSMCCounselor(appCtx)

		// TXO Queues
		createNTSMove(appCtx)
		createNTSRMove(appCtx)

		// This allows testing the pagination feature in the TXO queues.
		// Feel free to comment out the loop if you don't need this many moves.
		for i := 1; i < 12; i++ {
			createDefaultHHGMoveWithPaymentRequest(appCtx, userUploader, models.AffiliationAIRFORCE)
		}
		createDefaultHHGMoveWithPaymentRequest(appCtx, userUploader, models.AffiliationMARINES)
	}
}

func subScenarioPaymentRequestCalculations(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader,
	moveRouter services.MoveRouter) func() {
	return func() {
		createTXO(appCtx)
		createTXOUSMC(appCtx)

		// For displaying the Domestic Line Haul calculations displayed on the Payment Requests and Service Item review page
		createHHGMoveWithPaymentRequest(appCtx, userUploader, models.AffiliationAIRFORCE, testdatagen.Assertions{
			Move: models.Move{
				Locator: "SidDLH",
			},
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
			ReService: models.ReService{
				// DLH - Domestic line haul
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"),
			},
		})
		// Locator PARAMS
		createHHGWithPaymentServiceItems(appCtx, primeUploader, moveRouter)
	}
}

func subScenarioPPMAndHHG(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) func() {
	return func() {
		createTXO(appCtx)
		createTXOUSMC(appCtx)

		createMoveWithPPMAndHHG(appCtx, userUploader, moveRouter)
	}
}

func subScenarioDivertedShipments(appCtx appcontext.AppContext, userUploader *uploader.UserUploader,
	allDutyStations []models.DutyStation, originDutyStationsInGBLOC []models.DutyStation) func() {
	return func() {
		createTXO(appCtx)
		createTXOUSMC(appCtx)

		// Create diverted shipments that need TOO approval
		createMoveWithDivertedShipments(appCtx, userUploader)

		// Create diverted shipments that are approved and appear on the Move Task Order page
		createRandomMove(appCtx, nil, allDutyStations, originDutyStationsInGBLOC, true, testdatagen.Assertions{
			UserUploader: userUploader,
			Move: models.Move{
				Status:             models.MoveStatusAPPROVED,
				Locator:            "APRDVS",
				AvailableToPrimeAt: swag.Time(time.Now()),
			},
			MTOShipment: models.MTOShipment{
				Diversion:           true,
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        swag.Time(time.Now()),
				ScheduledPickupDate: swag.Time(time.Now().AddDate(0, 3, 0)),
			},
		})
	}
}

func subScenarioReweighs(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) func() {
	return func() {
		createHHGMoveWithReweigh(appCtx, userUploader)
		createHHGMoveWithBillableWeights(appCtx, userUploader, primeUploader)
		createReweighWithMultipleShipments(appCtx, userUploader, primeUploader, moveRouter)
		createReweighWithShipmentMissingReweigh(appCtx, userUploader, primeUploader, moveRouter)
		createReweighWithShipmentMaxBillableWeightExceeded(appCtx, userUploader, primeUploader, moveRouter)
		createReweighWithShipmentNoEstimatedWeight(appCtx, userUploader, primeUploader, moveRouter)
		createReweighWithShipmentDeprecatedPaymentRequest(appCtx, userUploader, primeUploader, moveRouter)
		createReweighWithMixedShipmentStatuses(appCtx, userUploader)
	}
}

func subScenarioSITExtensions(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) func() {
	return func() {
		createTOO(appCtx)
		createMoveWithSITExtensionHistory(appCtx, userUploader)
		createMoveWithAllPendingTOOActions(appCtx, userUploader, primeUploader)
	}
}

func subScenarioMisc(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader,
	moveRouter services.MoveRouter) func() {
	return func() {
		createTXOServicesCounselor(appCtx)
		createTXOServicesUSMCCounselor(appCtx)

		// A move with missing required order fields
		createMoveWithHHGMissingOrdersInfo(appCtx, moveRouter, userUploader)

		createHHGMoveWith10ServiceItems(appCtx, userUploader)
		createHHGMoveWith2PaymentRequests(appCtx, userUploader)
		createHHGMoveWith2PaymentRequestsReviewedAllRejectedServiceItems(appCtx, userUploader)
		createHHGMoveWithTaskOrderServices(appCtx, userUploader)

		// This one doesn't have submitted shipments. Can we get rid of it?
		// createRecentlyUpdatedHHGMove(appCtx, userUploader)
		createMoveWithHHGAndNTSRPaymentRequest(appCtx, userUploader)
		// This move will still have shipments with some unapproved service items
		// without payment service items
		createMoveWith2ShipmentsAndPaymentRequest(appCtx, userUploader)
		createMoveWith2MinimalShipments(appCtx, userUploader)
		createApprovedMoveWithMinimalShipment(appCtx, userUploader)

		// Prime API
		createWebhookSubscriptionForPaymentRequestUpdate(appCtx)
		// This move below is a PPM move in DRAFT status. It should probably
		// be changed to an HHG move in SUBMITTED status to reflect reality.
		createMoveWithServiceItems(appCtx, userUploader)
		createMoveWithBasicServiceItems(appCtx, userUploader)
		// Sets up a move with a non-default destination duty station address
		// (to more easily spot issues with addresses being overwritten).
		createMoveWithUniqueDestinationAddress(appCtx)
		// Creates a move that has multiple orders uploaded
		createHHGMoveWithMultipleOrdersFiles(appCtx, userUploader, primeUploader)
		createHHGMoveWithAmendedOrders(appCtx, userUploader, primeUploader)
		createHHGMoveWithRiskOfExcess(appCtx, userUploader, primeUploader)

		createMoveWithOriginAndDestinationSIT(appCtx, userUploader)
		createPaymentRequestsWithPartialSITInvoice(appCtx, primeUploader)
	}
}

func subScenarioNTSShipments(
	appCtx appcontext.AppContext,
	userUploader *uploader.UserUploader,
	moveRouter services.MoveRouter,
) func() {
	return func() {
		createTXO(appCtx)
		createTXOServicesCounselor(appCtx)

		// Create some unsubmitted Moves for Customer users
		createMoveWithNTSAndNTSR(
			appCtx,
			sceneOptionsNTS{
				ntsType:     "NTS",
				ntsMoveCode: "P8NTSU",
				moveStatus:  models.MoveStatusDRAFT,
			},
		)
		createMoveWithNTSAndNTSR(
			appCtx,
			sceneOptionsNTS{
				ntsType:     "NTSR",
				ntsMoveCode: "P8NTSU",
				moveStatus:  models.MoveStatusDRAFT,
			},
		)

		// Create some submitted Moves for TXO users
		createMoveWithNTSAndNTSR(
			appCtx,
			sceneOptionsNTS{
				ntsType:     "NTS",
				ntsMoveCode: "P8NTSS",
				moveStatus:  models.MoveStatusSUBMITTED,
			},
		)
		createMoveWithNTSAndNTSR(
			appCtx,
			sceneOptionsNTS{
				ntsType:     "NTSR",
				ntsMoveCode: "P8NTSS",
				moveStatus:  models.MoveStatusSUBMITTED,
			},
		)
	}
}
