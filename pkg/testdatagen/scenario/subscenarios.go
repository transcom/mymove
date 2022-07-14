package scenario

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func subScenarioShipmentHHGCancelled(appCtx appcontext.AppContext, allDutyLocations []models.DutyLocation, originDutyLocationsInGBLOC []models.DutyLocation) func() {
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
		move := createRandomMove(appCtx, validStatuses, allDutyLocations, originDutyLocationsInGBLOC, true, testdatagen.Assertions{
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

func subScenarioPPMCustomerFlow(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) func() {
	return func() {
		createTXO(appCtx)
		createTXOServicesCounselor(appCtx)
		createTXOUSMC(appCtx)

		// Onboarding
		createUnSubmittedMoveWithMinimumPPMShipment(appCtx, userUploader)
		createUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights(appCtx, userUploader)
		createUnSubmittedMoveWithPPMShipmentThroughAdvanceRequested(appCtx, userUploader)
		createUnsubmittedMoveWithMultipleFullPPMShipmentComplete1(appCtx, userUploader)
		createUnsubmittedMoveWithMultipleFullPPMShipmentComplete2(appCtx, userUploader)
		createSubmittedMoveWithFullPPMShipmentComplete(appCtx, userUploader)
		createUnSubmittedMoveWithFullPPMShipment1(appCtx, userUploader)
		createUnSubmittedMoveWithFullPPMShipment2(appCtx, userUploader)
		createUnSubmittedMoveWithFullPPMShipment3(appCtx, userUploader)
		createSubmittedMoveWithPPMShipment(appCtx, userUploader, moveRouter)
		createMoveWithPPM(appCtx, userUploader, moveRouter)
		createNeedsServicesCounselingWithoutCompletedOrders(appCtx, internalmessages.OrdersTypePERMANENTCHANGEOFSTATION, models.MTOShipmentTypePPM, nil, "SCPPM1")
		createSubmittedMoveWithPPMShipmentForSC(appCtx, userUploader, moveRouter, "PPMSC1")
		// Post-onboarding
		createApprovedMoveWithPPM(appCtx, userUploader)
		createApprovedMoveWithPPMWithActualDateZipsAndAdvanceInfo(appCtx, userUploader)
		createApprovedMoveWithPPMEmptyAboutPage(appCtx, userUploader)
	}
}

func subScenarioHHGServicesCounseling(appCtx appcontext.AppContext, userUploader *uploader.UserUploader,
	allDutyLocations []models.DutyLocation, originDutyLocationsInGBLOC []models.DutyLocation) func() {
	return func() {
		createTXOServicesCounselor(appCtx)
		createTXOServicesUSMCCounselor(appCtx)

		// Services Counseling
		//Order Types -- PCoS, Retr, Sep
		pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
		retirement := internalmessages.OrdersTypeRETIREMENT
		separation := internalmessages.OrdersTypeSEPARATION

		//Shipment Types -- HHG, NTS, NTSR
		hhg := models.MTOShipmentTypeHHG
		nts := models.MTOShipmentTypeHHGIntoNTSDom
		ntsR := models.MTOShipmentTypeHHGOutOfNTSDom

		//Destination Types -- PLEAD, HOR, HOS, OTHER
		plead := models.DestinationTypePlaceEnteredActiveDuty
		hor := models.DestinationTypeHomeOfRecord
		hos := models.DestinationTypeHomeOfSelection
		other := models.DestinationTypeOtherThanAuthorized

		//PCOS - one with nil dest type, 2 others with PLEAD status
		createNeedsServicesCounseling(appCtx, pcos, hhg, nil, "NODEST")
		createNeedsServicesCounseling(appCtx, pcos, nts, &plead, "PLEAD1")
		createNeedsServicesCounseling(appCtx, pcos, nts, &plead, "PLEAD2")

		//Retirees
		createNeedsServicesCounseling(appCtx, retirement, hhg, &hor, "RETIR3")
		createNeedsServicesCounseling(appCtx, retirement, nts, &hos, "RETIR4")
		createNeedsServicesCounseling(appCtx, retirement, ntsR, &other, "RETIR5")
		createNeedsServicesCounseling(appCtx, retirement, hhg, &plead, "RETIR6")

		//Separatees
		createNeedsServicesCounseling(appCtx, separation, hhg, &hor, "SEPAR3")
		createNeedsServicesCounseling(appCtx, separation, nts, &hos, "SEPAR4")
		createNeedsServicesCounseling(appCtx, separation, ntsR, &other, "SEPAR5")
		createNeedsServicesCounseling(appCtx, separation, ntsR, &plead, "SEPAR6")

		//USMC
		createHHGNeedsServicesCounselingUSMC(appCtx, userUploader)
		createHHGNeedsServicesCounselingUSMC2(appCtx, userUploader)
		createHHGServicesCounselingCompleted(appCtx)
		createHHGNoShipments(appCtx)

		for i := 0; i < 12; i++ {
			validStatuses := []models.MoveStatus{models.MoveStatusNeedsServiceCounseling, models.MoveStatusServiceCounselingCompleted}
			createRandomMove(appCtx, validStatuses, allDutyLocations, originDutyLocationsInGBLOC, false, testdatagen.Assertions{
				UserUploader: userUploader,
			})
		}
	}
}

func subScenarioCustomerSupportRemarks(appCtx appcontext.AppContext) func() {
	return func() {
		// Move with a couple of customer support remarks
		remarkMove := testdatagen.MakeMove(appCtx.DB(),
			testdatagen.Assertions{
				Move: models.Move{
					Locator: "SPTRMK",
					Status:  models.MoveStatusSUBMITTED,
				},
			},
		)
		_ = testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{MTOShipment: models.MTOShipment{
			MoveTaskOrderID: remarkMove.ID,
			Status:          models.MTOShipmentStatusSubmitted,
		}})

		officeUser := testdatagen.MakeDefaultOfficeUser(appCtx.DB())
		testdatagen.MakeCustomerSupportRemark(appCtx.DB(), testdatagen.Assertions{
			CustomerSupportRemark: models.CustomerSupportRemark{
				Content: "This is a customer support remark. It can have text content like this." +
					"This comment has some length to it because sometimes people type a lot of thoughts." +
					"For example during this move the customer perhaps called and explained a unique situation" +
					"that they have to me, leading me to leave this note. Hopefully that could turn into " +
					"some sort of helpful action that leads to a resolution that makes things swell for them." +
					"Here's some more text just to make sure I've gotten all my thoughts out, though I do realize" +
					"how meta this whole thing sounds." +
					"Also Grace Griffin told me to write this.",
				OfficeUserID: officeUser.ID,
				MoveID:       remarkMove.ID,
			},
		})
		officeUser2 := testdatagen.MakeDefaultOfficeUser(appCtx.DB())
		testdatagen.MakeCustomerSupportRemark(appCtx.DB(), testdatagen.Assertions{
			CustomerSupportRemark: models.CustomerSupportRemark{
				Content:      "The customer mentioned that there was some damage done to their grandfather clock.",
				OfficeUserID: officeUser2.ID,
				MoveID:       remarkMove.ID,
			},
		})
	}
}

func subScenarioEvaluationReport(appCtx appcontext.AppContext) func() {
	return func() {
		// Move with a few evaluation reports
		move := testdatagen.MakeMove(appCtx.DB(),
			testdatagen.Assertions{
				Move: models.Move{
					Locator: "EVLRPT",
					Status:  models.MoveStatusSUBMITTED,
				},
			},
		)
		shipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{MTOShipment: models.MTOShipment{
			MoveTaskOrderID: move.ID,
			Status:          models.MTOShipmentStatusSubmitted,
		}})
		testdatagen.MakePPMShipment(appCtx.DB(), testdatagen.Assertions{Move: move})

		storageFacility := testdatagen.MakeStorageFacility(appCtx.DB(), testdatagen.Assertions{
			StorageFacility: models.StorageFacility{
				FacilityName: "Storage R Us",
			},
		})

		testdatagen.MakeNTSShipment(appCtx.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				StorageFacility: &storageFacility,
			},
		})
		testdatagen.MakeNTSRShipment(appCtx.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				StorageFacility: &storageFacility,
			},
		})

		officeUser := testdatagen.MakeDefaultOfficeUser(appCtx.DB())
		submittedTime := time.Now()
		remark := "this is a submitted counseling report"
		location := models.EvaluationReportLocationTypeOrigin
		testdatagen.MakeEvaluationReport(appCtx.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				SubmittedAt:        &submittedTime,
				Location:           &location,
				ViolationsObserved: swag.Bool(false),
				Remarks:            &remark,
			},
			Move:       move,
			OfficeUser: officeUser,
		})
		remark1 := "this is a draft counseling report"
		testdatagen.MakeEvaluationReport(appCtx.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				Remarks: &remark1,
			},
			Move:       move,
			OfficeUser: officeUser,
		})
		location = models.EvaluationReportLocationTypeDestination
		remark2 := "this is a submitted shipment report"
		testdatagen.MakeEvaluationReport(appCtx.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				SubmittedAt:        &submittedTime,
				Location:           &location,
				ViolationsObserved: swag.Bool(true),
				Remarks:            &remark2,
			},
			Move:        move,
			OfficeUser:  officeUser,
			MTOShipment: shipment,
		})
		remark3 := "this is a draft shipment report"
		testdatagen.MakeEvaluationReport(appCtx.DB(), testdatagen.Assertions{
			EvaluationReport: models.EvaluationReport{
				Remarks: &remark3,
			},
			Move:        move,
			OfficeUser:  officeUser,
			MTOShipment: shipment,
		})
	}
}

func subScenarioTXOQueues(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) func() {
	return func() {
		createTOO(appCtx)
		createTIO(appCtx)
		createTXO(appCtx)
		createTXOUSMC(appCtx)
		createServicesCounselor(appCtx)
		createTXOServicesCounselor(appCtx)
		createTXOServicesUSMCCounselor(appCtx)
		createQaeCsr(appCtx)

		// TXO Queues
		createNTSMove(appCtx)
		createNTSRMove(appCtx)

		// This allows testing the pagination feature in the TXO queues.
		// Feel free to comment out the loop if you don't need this many moves.
		for i := 1; i < 12; i++ {
			createDefaultHHGMoveWithPaymentRequest(appCtx, userUploader, models.AffiliationAIRFORCE)
		}

		// Marines
		createDefaultHHGMoveWithPaymentRequest(appCtx, userUploader, models.AffiliationMARINES)

		//destination type
		hos := models.DestinationTypeHomeOfSelection
		hor := models.DestinationTypeHomeOfRecord

		//shipment type
		hhg := models.MTOShipmentTypeHHG
		nts := models.MTOShipmentTypeHHGIntoNTSDom
		ntsR := models.MTOShipmentTypeHHGOutOfNTSDom

		//orders type
		retirement := internalmessages.OrdersTypeRETIREMENT
		separation := internalmessages.OrdersTypeSEPARATION

		//Retiree, HOR, HHG
		createMoveWithOptions(appCtx, testdatagen.Assertions{
			Order: models.Order{
				OrdersType: retirement,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType:    hhg,
				DestinationType: &hor,
			},
			Move: models.Move{
				Locator: "R3T1R3",
				Status:  models.MoveStatusSUBMITTED,
			},
			DutyLocation: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		})

		//Retiree, HOS, NTS
		ntsMoveType := models.SelectedMoveTypeNTS
		createMoveWithOptions(appCtx, testdatagen.Assertions{
			Order: models.Order{
				OrdersType: retirement,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType:       nts,
				DestinationType:    &hor,
				UsesExternalVendor: false,
			},
			Move: models.Move{
				Locator:          "R3TNTS",
				Status:           models.MoveStatusSUBMITTED,
				SelectedMoveType: &ntsMoveType,
			},
			DutyLocation: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		})

		//Retiree, HOS, NTSR
		ntsrMoveType := models.SelectedMoveTypeNTSR
		createMoveWithOptions(appCtx, testdatagen.Assertions{
			Order: models.Order{
				OrdersType: retirement,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType:       ntsR,
				DestinationType:    &hos,
				UsesExternalVendor: false,
			},
			Move: models.Move{
				Locator:          "R3TNTR",
				Status:           models.MoveStatusSUBMITTED,
				SelectedMoveType: &ntsrMoveType,
			},
			DutyLocation: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		})

		//Separatee, HOS, hhg
		createMoveWithOptions(appCtx, testdatagen.Assertions{
			Order: models.Order{
				OrdersType: separation,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType:    hhg,
				DestinationType: &hos,
			},
			Move: models.Move{
				Locator: "S3P4R3",
				Status:  models.MoveStatusSUBMITTED,
			},
			DutyLocation: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		})
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
	allDutyLocations []models.DutyLocation, originDutyLocationsInGBLOC []models.DutyLocation) func() {
	return func() {
		createTXO(appCtx)
		createTXOUSMC(appCtx)

		// Create diverted shipments that need TOO approval
		createMoveWithDivertedShipments(appCtx)

		// Create diverted shipments that are approved and appear on the Move Task Order page
		createRandomMove(appCtx, nil, allDutyLocations, originDutyLocationsInGBLOC, true, testdatagen.Assertions{
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

func subScenarioNTSandNTSR(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) func() {
	return func() {
		pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION

		createTXO(appCtx)
		createTXOServicesCounselor(appCtx)

		createNeedsServicesCounselingSingleHHG(appCtx, pcos, "NTSHHG")
		createNeedsServicesCounselingSingleHHG(appCtx, pcos, "NTSRHG")
		createNeedsServicesCounselingMinimalNTSR(appCtx, pcos, "NTSRMN")

		// Create a move with an HHG and NTS prime-handled shipment
		createMoveWithHHGAndNTSShipments(appCtx, "PRINTS", false)

		// Create a move with an HHG and NTS external vendor-handled shipment
		createMoveWithHHGAndNTSShipments(appCtx, "PRXNTS", true)

		// Create a move with only NTS external vendor-handled shipment
		createMoveWithNTSShipment(appCtx, "EXTNTS", true)

		// Create a move with only an NTS external vendor-handled shipment
		createMoveWithNTSShipment(appCtx, "NTSNTS", true)

		// Create a move with an HHG and NTS-release prime-handled shipment
		createMoveWithHHGAndNTSRShipments(appCtx, "PRINTR", false)

		// Create a move with an HHG and NTS-release external vendor-handled shipment
		createMoveWithHHGAndNTSRShipments(appCtx, "PRXNTR", true)

		// Create a move with only an NTS-release external vendor-handled shipment
		createMoveWithNTSRShipment(appCtx, "EXTNTR", true)

		// Create some submitted Moves for TXO users
		createMoveWithHHGAndNTSRMissingInfo(appCtx, moveRouter)
		createMoveWithHHGAndNTSMissingInfo(appCtx, moveRouter)
		createMoveWithNTSAndNTSR(
			appCtx,
			userUploader,
			moveRouter,
			sceneOptionsNTS{
				shipmentMoveCode: "NTSSUB",
				moveStatus:       models.MoveStatusSUBMITTED,
			},
		)

		// uses external vendor
		createMoveWithNTSAndNTSR(
			appCtx,
			userUploader,
			moveRouter,
			sceneOptionsNTS{
				shipmentMoveCode:   "NTSEVR",
				moveStatus:         models.MoveStatusSUBMITTED,
				usesExternalVendor: true,
			},
		)

		// Create some unsubmitted Moves for Customer users
		// uses external vendor
		createMoveWithNTSAndNTSR(
			appCtx,
			userUploader,
			moveRouter,
			sceneOptionsNTS{
				shipmentMoveCode:   "NTSSUN",
				moveStatus:         models.MoveStatusDRAFT,
				usesExternalVendor: true,
			},
		)

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
		// Sets up a move with a non-default destination duty location address
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

func subScenarioPrimeUserAndClientCert(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) func() {
	return func() {
		primeUser := createPrimeUser(appCtx)
		createDevClientCertForUser(appCtx, primeUser)
	}
}
