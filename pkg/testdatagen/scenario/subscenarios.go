package scenario

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func subScenarioShipmentHHGCancelled(appCfg appconfig.AppConfig, allDutyStations []models.DutyStation, originDutyStationsInGBLOC []models.DutyStation) func() {
	db := appCfg.DB()
	return func() {
		createTXO(appCfg)
		createTXOUSMC(appCfg)

		validStatuses := []models.MoveStatus{models.MoveStatusAPPROVED}
		// shipment cancelled was approved before
		approvedDate := time.Now()
		cancelledShipment := models.MTOShipment{Status: models.MTOShipmentStatusCanceled, ApprovedDate: &approvedDate}
		affiliationAirForce := models.AffiliationAIRFORCE
		ordersNumber := "Order1234"
		ordersTypeDetail := internalmessages.OrdersTypeDetailHHGPERMITTED
		tac := "1234"
		// make sure to create moves that does not go to US marines affiliation
		move := createRandomMove(appCfg, validStatuses, allDutyStations, originDutyStationsInGBLOC, true, testdatagen.Assertions{
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

func subScenarioPPMOfficeQueue(appCfg appconfig.AppConfig, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) func() {
	return func() {
		createPPMOfficeUser(appCfg)

		// PPM Office Queue
		createPPMWithAdvance(appCfg, userUploader, moveRouter)
		createPPMWithNoAdvance(appCfg, userUploader, moveRouter)
		createPPMWithPaymentRequest(appCfg, userUploader, moveRouter)
		createCanceledPPM(appCfg, userUploader, moveRouter)
		createPPMReadyToRequestPayment(appCfg, userUploader, moveRouter)
	}
}

func subScenarioAdditionalPPMUsers(appCfg appconfig.AppConfig, userUploader *uploader.UserUploader) func() {
	return func() {
		// Create additional PPM users for mymove tests
		createPPMUsers(appCfg, userUploader)
	}
}

func subScenarioHHGOnboarding(appCfg appconfig.AppConfig, userUploader *uploader.UserUploader) func() {
	return func() {
		createTXO(appCfg)
		createTXOUSMC(appCfg)

		// Onboarding
		createUnsubmittedHHGMove(appCfg)
		createUnsubmittedMoveWithNTSAndNTSR(appCfg, 1)
		createUnsubmittedMoveWithNTSAndNTSR(appCfg, 2)
		createUnsubmittedHHGMoveMultiplePickup(appCfg)
		createUnsubmittedHHGMoveMultipleDestinations(appCfg)
		createServiceMemberWithOrdersButNoMoveType(appCfg)
		createServiceMemberWithNoUploadedOrders(appCfg)
		createSubmittedHHGMoveMultiplePickupAmendedOrders(appCfg, userUploader)
	}
}

func subScenarioHHGServicesCounseling(appCfg appconfig.AppConfig, userUploader *uploader.UserUploader,
	allDutyStations []models.DutyStation, originDutyStationsInGBLOC []models.DutyStation) func() {
	return func() {
		createTXOServicesCounselor(appCfg)
		createTXOServicesUSMCCounselor(appCfg)

		// Services Counseling
		createHHGNeedsServicesCounseling(appCfg)
		createHHGNeedsServicesCounselingUSMC(appCfg, userUploader)
		createHHGNeedsServicesCounselingUSMC2(appCfg, userUploader)
		createHHGServicesCounselingCompleted(appCfg)
		createHHGNoShipments(appCfg)

		for i := 0; i < 12; i++ {
			validStatuses := []models.MoveStatus{models.MoveStatusNeedsServiceCounseling, models.MoveStatusServiceCounselingCompleted}
			createRandomMove(appCfg, validStatuses, allDutyStations, originDutyStationsInGBLOC, false, testdatagen.Assertions{
				UserUploader: userUploader,
			})
		}
	}
}

func subScenarioTXOQueues(appCfg appconfig.AppConfig, userUploader *uploader.UserUploader, logger *zap.Logger) func() {
	return func() {
		createTOO(appCfg)
		createTIO(appCfg)
		createTXO(appCfg)
		createTXOUSMC(appCfg)
		createServicesCounselor(appCfg)
		createTXOServicesCounselor(appCfg)
		createTXOServicesUSMCCounselor(appCfg)

		// TXO Queues
		createNTSMove(appCfg)
		createNTSRMove(appCfg)

		// This allows testing the pagination feature in the TXO queues.
		// Feel free to comment out the loop if you don't need this many moves.
		for i := 1; i < 12; i++ {
			createDefaultHHGMoveWithPaymentRequest(appCfg, userUploader, models.AffiliationAIRFORCE)
		}
		createDefaultHHGMoveWithPaymentRequest(appCfg, userUploader, models.AffiliationMARINES)
	}
}

func subScenarioPaymentRequestCalculations(appCfg appconfig.AppConfig, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader,
	moveRouter services.MoveRouter) func() {
	return func() {
		createTXO(appCfg)
		createTXOUSMC(appCfg)

		// For displaying the Domestic Line Haul calculations displayed on the Payment Requests and Service Item review page
		createHHGMoveWithPaymentRequest(appCfg, userUploader, models.AffiliationAIRFORCE, testdatagen.Assertions{
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
		createHHGWithPaymentServiceItems(appCfg, primeUploader, moveRouter)
	}
}

func subScenarioPPMAndHHG(appCfg appconfig.AppConfig, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) func() {
	return func() {
		createTXO(appCfg)
		createTXOUSMC(appCfg)

		createMoveWithPPMAndHHG(appCfg, userUploader, moveRouter)
	}
}

func subScenarioDivertedShipments(appCfg appconfig.AppConfig, userUploader *uploader.UserUploader,
	allDutyStations []models.DutyStation, originDutyStationsInGBLOC []models.DutyStation) func() {
	return func() {
		createTXO(appCfg)
		createTXOUSMC(appCfg)

		// Create diverted shipments that need TOO approval
		createMoveWithDivertedShipments(appCfg, userUploader)

		// Create diverted shipments that are approved and appear on the Move Task Order page
		createRandomMove(appCfg, nil, allDutyStations, originDutyStationsInGBLOC, true, testdatagen.Assertions{
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

func subScenarioReweighs(appCfg appconfig.AppConfig, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) func() {
	return func() {
		createHHGMoveWithReweigh(appCfg, userUploader)
		createHHGMoveWithBillableWeights(appCfg, userUploader, primeUploader)
		createReweighWithMultipleShipments(appCfg, userUploader, primeUploader, moveRouter)
		createReweighWithShipmentMissingReweigh(appCfg, userUploader, primeUploader, moveRouter)
		createReweighWithShipmentMaxBillableWeightExceeded(appCfg, userUploader, primeUploader, moveRouter)
		createReweighWithShipmentNoEstimatedWeight(appCfg, userUploader, primeUploader, moveRouter)
	}
}

func subScenarioMisc(appCfg appconfig.AppConfig, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader,
	moveRouter services.MoveRouter) func() {
	return func() {
		createTXOServicesCounselor(appCfg)
		createTXOServicesUSMCCounselor(appCfg)

		// A move with missing required order fields
		createMoveWithHHGMissingOrdersInfo(appCfg, moveRouter)

		createHHGMoveWith10ServiceItems(appCfg, userUploader)
		createHHGMoveWith2PaymentRequests(appCfg, userUploader)
		createHHGMoveWith2PaymentRequestsReviewedAllRejectedServiceItems(appCfg, userUploader)
		createHHGMoveWithTaskOrderServices(appCfg, userUploader)

		// This one doesn't have submitted shipments. Can we get rid of it?
		// createRecentlyUpdatedHHGMove(appCfg, userUploader)
		createMoveWithHHGAndNTSRPaymentRequest(appCfg, userUploader)
		// This move will still have shipments with some unapproved service items
		// without payment service items
		createMoveWith2ShipmentsAndPaymentRequest(appCfg, userUploader)
		createMoveWith2MinimalShipments(appCfg, userUploader)

		// Prime API
		createWebhookSubscriptionForPaymentRequestUpdate(appCfg)
		// This move below is a PPM move in DRAFT status. It should probably
		// be changed to an HHG move in SUBMITTED status to reflect reality.
		createMoveWithServiceItems(appCfg, userUploader)
		createMoveWithBasicServiceItems(appCfg, userUploader)
		// Sets up a move with a non-default destination duty station address
		// (to more easily spot issues with addresses being overwritten).
		createMoveWithUniqueDestinationAddress(appCfg)
		// Creates a move that has multiple orders uploaded
		createHHGMoveWithMultipleOrdersFiles(appCfg, userUploader, primeUploader)
		createHHGMoveWithAmendedOrders(appCfg, userUploader, primeUploader)
		createHHGMoveWithRiskOfExcess(appCfg, userUploader, primeUploader)
	}
}
