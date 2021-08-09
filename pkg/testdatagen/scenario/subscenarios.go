package scenario

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func subScenarioShipmentHHGCancelled(db *pop.Connection, allDutyStations []models.DutyStation, originDutyStationsInGBLOC []models.DutyStation) func() {
	return func() {
		createTXO(db)
		createTXOUSMC(db)

		validStatuses := []models.MoveStatus{models.MoveStatusAPPROVED}
		// shipment cancelled was approved before
		approvedDate := time.Now()
		cancelledShipment := models.MTOShipment{Status: models.MTOShipmentStatusCanceled, ApprovedDate: &approvedDate}
		affiliationAirForce := models.AffiliationAIRFORCE
		ordersNumber := "Order1234"
		ordersTypeDetail := internalmessages.OrdersTypeDetailHHGPERMITTED
		tac := "1234"
		// make sure to create moves that does not go to US marines affiliation
		move := createRandomMove(db, validStatuses, allDutyStations, originDutyStationsInGBLOC, true, testdatagen.Assertions{
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

func subScenarioPPMOfficeQueue(db *pop.Connection, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) func() {
	return func() {
		createPPMOfficeUser(db)

		// PPM Office Queue
		createPPMWithAdvance(db, userUploader, moveRouter)
		createPPMWithNoAdvance(db, userUploader, moveRouter)
		createPPMWithPaymentRequest(db, userUploader, moveRouter)
		createCanceledPPM(db, userUploader, moveRouter)
		createPPMReadyToRequestPayment(db, userUploader, moveRouter)
	}
}

func subScenarioAdditionalPPMUsers(db *pop.Connection, userUploader *uploader.UserUploader) func() {
	return func() {
		// Create additional PPM users for mymove tests
		createPPMUsers(db, userUploader)
	}
}

func subScenarioHHGOnboarding(db *pop.Connection, userUploader *uploader.UserUploader) func() {
	return func() {
		createTXO(db)
		createTXOUSMC(db)

		// Onboarding
		createUnsubmittedHHGMove(db)
		createUnsubmittedMoveWithNTSAndNTSR(db, 1)
		createUnsubmittedMoveWithNTSAndNTSR(db, 2)
		createUnsubmittedHHGMoveMultiplePickup(db)
		createUnsubmittedHHGMoveMultipleDestinations(db)
		createServiceMemberWithOrdersButNoMoveType(db)
		createServiceMemberWithNoUploadedOrders(db)
		createSubmittedHHGMoveMultiplePickupAmendedOrders(db, userUploader)
	}
}

func subScenarioHHGServicesCounseling(db *pop.Connection, userUploader *uploader.UserUploader,
	allDutyStations []models.DutyStation, originDutyStationsInGBLOC []models.DutyStation) func() {
	return func() {
		createTXOServicesCounselor(db)
		createTXOServicesUSMCCounselor(db)

		// Services Counseling
		createHHGNeedsServicesCounseling(db)
		createHHGNeedsServicesCounselingUSMC(db, userUploader)
		createHHGNeedsServicesCounselingUSMC2(db, userUploader)
		createHHGServicesCounselingCompleted(db)
		createHHGNoShipments(db)

		for i := 0; i < 12; i++ {
			validStatuses := []models.MoveStatus{models.MoveStatusNeedsServiceCounseling, models.MoveStatusServiceCounselingCompleted}
			createRandomMove(db, validStatuses, allDutyStations, originDutyStationsInGBLOC, false, testdatagen.Assertions{
				UserUploader: userUploader,
			})
		}
	}
}

func subScenarioTXOQueues(db *pop.Connection, userUploader *uploader.UserUploader, logger *zap.Logger) func() {
	return func() {
		createTOO(db)
		createTIO(db)
		createTXO(db)
		createTXOUSMC(db)
		createServicesCounselor(db)
		createTXOServicesCounselor(db)
		createTXOServicesUSMCCounselor(db)

		// TXO Queues
		createNTSMove(db)
		createNTSRMove(db)

		// This allows testing the pagination feature in the TXO queues.
		// Feel free to comment out the loop if you don't need this many moves.
		for i := 1; i < 12; i++ {
			createDefaultHHGMoveWithPaymentRequest(db, userUploader, logger, models.AffiliationAIRFORCE)
		}
		createDefaultHHGMoveWithPaymentRequest(db, userUploader, logger, models.AffiliationMARINES)
	}
}

func subScenarioPaymentRequestCalculations(db *pop.Connection, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader,
	moveRouter services.MoveRouter, logger *zap.Logger) func() {
	return func() {
		createTXO(db)
		createTXOUSMC(db)

		// For displaying the Domestic Line Haul calculations displayed on the Payment Requests and Service Item review page
		createHHGMoveWithPaymentRequest(db, userUploader, logger, models.AffiliationAIRFORCE, testdatagen.Assertions{
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
		createHHGWithPaymentServiceItems(db, primeUploader, logger, moveRouter)
	}
}

func subScenarioPPMAndHHG(db *pop.Connection, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) func() {
	return func() {
		createTXO(db)
		createTXOUSMC(db)

		createMoveWithPPMAndHHG(db, userUploader, moveRouter)
	}
}

func subScenarioDivertedShipments(db *pop.Connection, userUploader *uploader.UserUploader,
	allDutyStations []models.DutyStation, originDutyStationsInGBLOC []models.DutyStation) func() {
	return func() {
		createTXO(db)
		createTXOUSMC(db)

		// Create diverted shipments that need TOO approval
		createMoveWithDivertedShipments(db, userUploader)

		// Create diverted shipments that are approved and appear on the Move Task Order page
		createRandomMove(db, nil, allDutyStations, originDutyStationsInGBLOC, true, testdatagen.Assertions{
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

func subScenarioReweighs(db *pop.Connection, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) func() {
	return func() {
		createHHGMoveWithReweigh(db, userUploader)
		createHHGMoveWithBillableWeights(db, userUploader, primeUploader)
		createReweighWithMultipleShipments(db, userUploader, primeUploader, moveRouter)
		createReweighWithShipmentMissingReweigh(db, userUploader, primeUploader, moveRouter)
		createReweighWithShipmentMaxBillableWeightExceeded(db, userUploader, primeUploader, moveRouter)
		createReweighWithShipmentNoEstimatedWeight(db, userUploader, primeUploader, moveRouter)
	}
}

func subScenarioMisc(db *pop.Connection, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader,
	moveRouter services.MoveRouter) func() {
	return func() {
		createTXOServicesCounselor(db)
		createTXOServicesUSMCCounselor(db)

		// A move with missing required order fields
		createMoveWithHHGMissingOrdersInfo(db, moveRouter)

		createHHGMoveWith10ServiceItems(db, userUploader)
		createHHGMoveWith2PaymentRequests(db, userUploader)
		createHHGMoveWith2PaymentRequestsReviewedAllRejectedServiceItems(db, userUploader)
		createHHGMoveWithTaskOrderServices(db, userUploader)

		// This one doesn't have submitted shipments. Can we get rid of it?
		// createRecentlyUpdatedHHGMove(db, userUploader)
		createMoveWithHHGAndNTSRPaymentRequest(db, userUploader)
		// This move will still have shipments with some unapproved service items
		// without payment service items
		createMoveWith2ShipmentsAndPaymentRequest(db, userUploader)
		createMoveWith2MinimalShipments(db, userUploader)

		// Prime API
		createWebhookSubscriptionForPaymentRequestUpdate(db)
		// This move below is a PPM move in DRAFT status. It should probably
		// be changed to an HHG move in SUBMITTED status to reflect reality.
		createMoveWithServiceItems(db, userUploader)
		createMoveWithBasicServiceItems(db, userUploader)
		// Sets up a move with a non-default destination duty station address
		// (to more easily spot issues with addresses being overwritten).
		createMoveWithUniqueDestinationAddress(db)
		// Creates a move that has multiple orders uploaded
		createHHGMoveWithMultipleOrdersFiles(db, userUploader, primeUploader)
		createHHGMoveWithAmendedOrders(db, userUploader, primeUploader)
		createHHGMoveWithRiskOfExcess(db, userUploader, primeUploader)
	}
}
