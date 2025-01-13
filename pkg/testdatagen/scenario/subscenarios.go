package scenario

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
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
		move := createRandomMove(appCtx, validStatuses, allDutyLocations, originDutyLocationsInGBLOC, true,
			nil,
			models.Move{
				Locator: "HHGCAN",
			},
			cancelledShipment,
			models.Order{
				DepartmentIndicator: (*string)(&affiliationAirForce),
				OrdersNumber:        &ordersNumber,
				OrdersTypeDetail:    &ordersTypeDetail,
				TAC:                 &tac,
			},
			models.ServiceMember{Affiliation: &affiliationAirForce},
		)
		moveManagementUUID := "1130e612-94eb-49a7-973d-72f33685e551"
		factory.BuildMTOServiceItemBasic(db, []factory.Customization{
			{
				Model: models.ReService{ID: uuid.FromStringOrNil(moveManagementUUID)},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status:     models.MTOServiceItemStatusApproved,
					ApprovedAt: &approvedDate,
				},
			},
		}, nil)
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

func subScenarioPPMCloseOut(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) func() {
	return func() {
		createServicesCounselor(appCtx)
		createServicesCounselorForCloseoutWithGbloc(appCtx, uuid.Must(uuid.FromString("8bed04b0-9c64-4fb2-bc1a-65223319f109")), "ppm.processing.navy@office.mil", "NAVY")
		createServicesCounselorForCloseoutWithGbloc(appCtx, uuid.Must(uuid.FromString("96b45e64-a501-49ce-9f8d-975bcb7b417c")), "ppm.processing.coastguard@office.mil", "USCG")
		createServicesCounselorForCloseoutWithGbloc(appCtx, uuid.Must(uuid.FromString("f38cb4ed-fa1f-4f92-a52d-9695ba8cc85c")), "ppm.processing.marinecorps@office.mil", "TVCB")

		// PPM Closeout
		createMovesForEachBranch(appCtx, userUploader)
		CreateMoveWithCloseoutOffice(appCtx, MoveCreatorInfo{
			UserID:      uuid.Must(uuid.NewV4()),
			Email:       "closeoutoffice@ppm.closeout",
			SmID:        uuid.Must(uuid.NewV4()),
			FirstName:   "CLOSEOUT",
			LastName:    "OFFICE",
			MoveID:      uuid.Must(uuid.NewV4()),
			MoveLocator: "CLSOFF",
		}, userUploader)
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

		CreateSubmittedMoveWithPPMShipmentForSC(appCtx, userUploader, moveRouter,
			MoveCreatorInfo{
				UserID:      uuid.Must(uuid.NewV4()),
				Email:       "complete@ppm.submitted",
				FirstName:   "PPMSC",
				LastName:    "Submitted",
				SmID:        uuid.Must(uuid.NewV4()),
				MoveLocator: "PPMSC1",
				MoveID:      uuid.Must(uuid.NewV4()),
			},
		)
		CreateSubmittedMoveWithPPMShipmentForSC(appCtx, userUploader, moveRouter,
			MoveCreatorInfo{
				UserID:      uuid.Must(uuid.NewV4()),
				Email:       "complete@ppm.submitted",
				FirstName:   "PPMSC",
				LastName:    "Submitted",
				SmID:        uuid.Must(uuid.NewV4()),
				MoveLocator: "PPMADD",
				MoveID:      uuid.Must(uuid.NewV4()),
			},
		)
		CreateSubmittedMoveWithPPMShipmentForSC(appCtx, userUploader, moveRouter,
			MoveCreatorInfo{
				UserID:      uuid.Must(uuid.NewV4()),
				Email:       "complete@ppm.submitted",
				FirstName:   "PPMSC",
				LastName:    "Submitted",
				SmID:        uuid.Must(uuid.NewV4()),
				MoveLocator: "PPMSCF",
				MoveID:      uuid.Must(uuid.NewV4()),
			},
		)
		createSubmittedMoveWithPPMShipmentForSCWithSIT(appCtx, userUploader, moveRouter, "PPMSIT")
		// Post-onboarding
		createApprovedMoveWithPPM(appCtx, userUploader)
		createApprovedMoveWithPPMWithAboutFormComplete(appCtx, userUploader)
		createApprovedMoveWithPPMWithAboutFormComplete2(appCtx, userUploader)
		createApprovedMoveWithPPMWithAboutFormComplete3(appCtx, userUploader)
		createApprovedMoveWithPPMWithAboutFormComplete4(appCtx, userUploader)
		createApprovedMoveWithPPMWithAboutFormComplete5(appCtx, userUploader)
		createApprovedMoveWithPPMWithAboutFormComplete6(appCtx, userUploader)
		createApprovedMoveWithPPM2(appCtx, userUploader)
		createApprovedMoveWithPPMWeightTicket(appCtx, userUploader)
		createApprovedMoveWithPPMExcessWeight(appCtx, userUploader,
			MoveCreatorInfo{
				UserID:      uuid.Must(uuid.NewV4()),
				Email:       "excessweightsPPM@ppm.approved",
				SmID:        uuid.Must(uuid.NewV4()),
				FirstName:   "PPM",
				LastName:    "ExcessWeights",
				MoveID:      uuid.Must(uuid.NewV4()),
				MoveLocator: "XSWT01",
			})
		createApprovedMoveWithPPMExcessWeightsAnd2WeightTickets(appCtx, userUploader)
		createApprovedMoveWith2PPMShipmentsAndExcessWeights(appCtx, userUploader)
		createApprovedMoveWithPPMAndHHGShipmentsAndExcessWeights(appCtx, userUploader)
		createApprovedMoveWithAllShipmentTypesAndExcessWeights(appCtx, userUploader)
		createApprovedMoveWithPPMMovingExpense(appCtx, nil, userUploader)
		createApprovedMoveWithPPMProgearWeightTicket(appCtx, userUploader)
		createApprovedMoveWithPPMProgearWeightTicket2(appCtx, userUploader)
		createMoveWithPPMShipmentReadyForFinalCloseout(appCtx, userUploader)
		createApprovedMoveWithPPMCloseoutComplete(appCtx, userUploader)
		createApprovedMoveWithPPMCloseoutCompleteMultipleWeightTickets(appCtx, userUploader)
		createApprovedMoveWithPPMCloseoutCompleteWithAllDocTypes(appCtx, userUploader)
		createApprovedMoveWithPPMCloseoutCompleteWithExpenses(appCtx, userUploader)
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
		nts := models.MTOShipmentTypeHHGIntoNTS
		ntsR := models.MTOShipmentTypeHHGOutOfNTSDom

		//Destination Types -- PLEAD, HOR, HOS, OTHER
		plead := models.DestinationTypePlaceEnteredActiveDuty
		hor := models.DestinationTypeHomeOfRecord
		hos := models.DestinationTypeHomeOfSelection
		other := models.DestinationTypeOtherThanAuthorized

		//PCOS - one with nil dest type, 2 others with PLEAD status
		CreateNeedsServicesCounseling(appCtx, pcos, hhg, nil, "NODEST")
		CreateNeedsServicesCounseling(appCtx, pcos, nts, &plead, "PLEAD1")
		CreateNeedsServicesCounseling(appCtx, pcos, nts, &plead, "PLEAD2")

		//Retirees
		CreateNeedsServicesCounseling(appCtx, retirement, hhg, &hor, "RETIR3")
		CreateNeedsServicesCounseling(appCtx, retirement, nts, &hos, "RETIR4")
		CreateNeedsServicesCounseling(appCtx, retirement, ntsR, &other, "RETIR5")
		CreateNeedsServicesCounseling(appCtx, retirement, hhg, &plead, "RETIR6")

		//Separatees
		CreateNeedsServicesCounseling(appCtx, separation, hhg, &hor, "SEPAR3")
		CreateNeedsServicesCounseling(appCtx, separation, nts, &hos, "SEPAR4")
		CreateNeedsServicesCounseling(appCtx, separation, ntsR, &other, "SEPAR5")
		CreateNeedsServicesCounseling(appCtx, separation, ntsR, &plead, "SEPAR6")

		//USMC
		createHHGNeedsServicesCounselingUSMC(appCtx, userUploader)
		createHHGNeedsServicesCounselingUSMC2(appCtx, userUploader)
		createHHGServicesCounselingCompleted(appCtx)
		createHHGNoShipments(appCtx)

		for i := 0; i < 12; i++ {
			validStatuses := []models.MoveStatus{models.MoveStatusNeedsServiceCounseling, models.MoveStatusServiceCounselingCompleted}
			createRandomMove(appCtx, validStatuses, allDutyLocations, originDutyLocationsInGBLOC, false,
				userUploader,
				models.Move{},
				models.MTOShipment{},
				models.Order{},
				models.ServiceMember{},
			)
		}
	}
}

func subScenarioCustomerSupportRemarks(appCtx appcontext.AppContext) func() {
	return func() {
		// Move with a couple of customer support remarks
		remarkMove := factory.BuildMove(appCtx.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "SPTRMK",
					Status:  models.MoveStatusSUBMITTED,
				},
			},
		}, nil)
		_ = factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model:    remarkMove,
				LinkOnly: true,
			},
		}, nil)

		officeUser := factory.BuildOfficeUserWithRoles(appCtx.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		factory.BuildCustomerSupportRemark(appCtx.DB(), []factory.Customization{
			{
				Model:    remarkMove,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model: models.CustomerSupportRemark{
					Content: "This is a customer support remark. It can have text content like this." +
						"This comment has some length to it because sometimes people type a lot of thoughts." +
						"For example during this move the customer perhaps called and explained a unique situation" +
						"that they have to me, leading me to leave this note. Hopefully that could turn into " +
						"some sort of helpful action that leads to a resolution that makes things swell for them." +
						"Here's some more text just to make sure I've gotten all my thoughts out, though I do realize" +
						"how meta this whole thing sounds." +
						"Also Grace Griffin told me to write this.",
				},
			},
		}, nil)
		officeUser2 := factory.BuildOfficeUserWithRoles(appCtx.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

		factory.BuildCustomerSupportRemark(appCtx.DB(), []factory.Customization{
			{
				Model:    remarkMove,
				LinkOnly: true,
			},
			{
				Model:    officeUser2,
				LinkOnly: true,
			},
			{
				Model: models.CustomerSupportRemark{
					Content: "The customer mentioned that there was some damage done to their grandfather clock.",
				},
			},
		}, nil)
	}
}

func subScenarioEvaluationReport(appCtx appcontext.AppContext) func() {
	return func() {
		createQae(appCtx)
		officeUser := models.OfficeUser{}
		email := "qae_role@office.mil"
		err := appCtx.DB().Where("email = ?", email).First(&officeUser)
		if err != nil {
			appCtx.Logger().Panic(fmt.Errorf("failed to query OfficeUser in the DB: %w", err).Error())
		}
		// Move with a few evaluation reports
		move := factory.BuildMove(appCtx.DB(), []factory.Customization{
			{
				Model: models.Move{
					Locator: "EVLRPT",
					Status:  models.MoveStatusSUBMITTED,
				},
			},
		}, nil)
		// Make a transportation office to use as the closeout office
		closeoutOffice := factory.BuildTransportationOffice(appCtx.DB(), []factory.Customization{
			{
				Model: models.TransportationOffice{Name: "Los Angeles AFB"},
			},
		}, nil)

		if *move.Orders.ServiceMember.Affiliation == models.AffiliationARMY || *move.Orders.ServiceMember.Affiliation == models.AffiliationAIRFORCE {
			move.CloseoutOffice = &closeoutOffice
			move.CloseoutOfficeID = &closeoutOffice.ID
			testdatagen.MustSave(appCtx.DB(), &move)
		}

		shipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:                models.MTOShipmentStatusSubmitted,
					ScheduledDeliveryDate: models.TimePointer(time.Now()),
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		factory.BuildPPMShipment(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		storageFacility := factory.BuildStorageFacility(appCtx.DB(), nil, nil)
		ntsShipment := factory.BuildNTSShipment(appCtx.DB(), []factory.Customization{
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ScheduledDeliveryDate: models.TimePointer(time.Now()),
				},
			},
		}, nil)
		factory.BuildNTSRShipment(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ScheduledDeliveryDate: models.TimePointer(time.Now()),
				},
			},
		}, nil)

		submittedTime := time.Now()
		dataReviewInspection := models.EvaluationReportInspectionTypeDataReview
		physicalInspection := models.EvaluationReportInspectionTypePhysical
		virtualInspection := models.EvaluationReportInspectionTypeVirtual
		inspectionTime := time.Now().AddDate(0, 0, -4)
		timeDepart := inspectionTime.Add(time.Hour * 1)
		evalStart := inspectionTime.Add(time.Hour * 3)
		evalEnd := inspectionTime.Add(time.Hour * 5)

		remark := "this is a submitted counseling report"
		location := models.EvaluationReportLocationTypeOrigin

		factory.BuildEvaluationReport(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					SubmittedAt:        &submittedTime,
					InspectionDate:     &submittedTime,
					InspectionType:     &dataReviewInspection,
					Location:           &location,
					ViolationsObserved: models.BoolPointer(false),
					Remarks:            &remark,
				},
			},
		}, nil)

		remark1 := "this is a draft counseling report"
		factory.BuildEvaluationReport(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					Remarks: &remark1,
				},
			},
		}, nil)

		location = models.EvaluationReportLocationTypeDestination
		remark2 := "this is a submitted shipment report"

		factory.BuildEvaluationReport(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					Type:               models.EvaluationReportTypeShipment,
					SubmittedAt:        &submittedTime,
					InspectionDate:     &submittedTime,
					InspectionType:     &virtualInspection,
					Location:           &location,
					ViolationsObserved: models.BoolPointer(true),
					Remarks:            &remark2,
				},
			},
		}, nil)

		remark3 := "this is a draft shipment report"
		factory.BuildEvaluationReport(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					Type:    models.EvaluationReportTypeShipment,
					Remarks: &remark3,
				},
			},
		}, nil)
		location = models.EvaluationReportLocationTypeOrigin
		remark4 := "this is a report with eval times recorded"

		factory.BuildEvaluationReport(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model:    shipment,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					Type:               models.EvaluationReportTypeShipment,
					Remarks:            &remark4,
					InspectionDate:     &submittedTime,
					InspectionType:     &physicalInspection,
					TimeDepart:         &timeDepart,
					EvalStart:          &evalStart,
					EvalEnd:            &evalEnd,
					Location:           &location,
					ViolationsObserved: models.BoolPointer(true),
				},
			},
		}, nil)

		location = models.EvaluationReportLocationTypeOther
		locationDescription := "Route 66 at crash inspection site 3"
		remark = "this is a submitted NTS shipment report"

		factory.BuildEvaluationReport(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model:    ntsShipment,
				LinkOnly: true,
			},
			{
				Model: models.EvaluationReport{
					Type:                models.EvaluationReportTypeShipment,
					SubmittedAt:         &submittedTime,
					InspectionDate:      &submittedTime,
					InspectionType:      &physicalInspection,
					TimeDepart:          &timeDepart,
					EvalStart:           &evalStart,
					EvalEnd:             &evalEnd,
					Location:            &location,
					LocationDescription: &locationDescription,
					ViolationsObserved:  models.BoolPointer(true),
					Remarks:             &remark,
				},
			},
		}, nil)
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
		createQae(appCtx)
		createCustomerServiceRepresentative(appCtx)

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
		nts := models.MTOShipmentTypeHHGIntoNTS
		ntsR := models.MTOShipmentTypeHHGOutOfNTSDom

		//orders type
		retirement := internalmessages.OrdersTypeRETIREMENT
		separation := internalmessages.OrdersTypeSEPARATION

		//Retiree, HOR, HHG
		CreateMoveWithOptions(appCtx, testdatagen.Assertions{
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
		CreateMoveWithOptions(appCtx, testdatagen.Assertions{
			Order: models.Order{
				OrdersType: retirement,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType:       nts,
				DestinationType:    &hor,
				UsesExternalVendor: false,
			},
			Move: models.Move{
				Locator: "R3TNTS",
				Status:  models.MoveStatusSUBMITTED,
			},
			DutyLocation: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		})

		//Retiree, HOS, NTSR
		CreateMoveWithOptions(appCtx, testdatagen.Assertions{
			Order: models.Order{
				OrdersType: retirement,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType:       ntsR,
				DestinationType:    &hos,
				UsesExternalVendor: false,
			},
			Move: models.Move{
				Locator: "R3TNTR",
				Status:  models.MoveStatusSUBMITTED,
			},
			DutyLocation: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		})

		//Separatee, HOS, hhg
		CreateMoveWithOptions(appCtx, testdatagen.Assertions{
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

func subScenarioPaymentRequestCalculations(
	appCtx appcontext.AppContext,
	userUploader *uploader.UserUploader,
	primeUploader *uploader.PrimeUploader,
	moveRouter services.MoveRouter,
	shipmentFetcher services.MTOShipmentFetcher,
) func() {
	if appCtx == nil || userUploader == nil || primeUploader == nil || moveRouter == nil || shipmentFetcher == nil {
		panic("nil argument passed to subScenarioPaymentRequestCalculations")
	}

	return func() {
		if appCtx == nil || userUploader == nil || primeUploader == nil || moveRouter == nil || shipmentFetcher == nil {
			panic("nil argument passed to subScenarioPaymentRequestCalculations")
		}

		createTXO(appCtx)
		createTXOUSMC(appCtx)

		createHHGMoveWithPaymentRequest(appCtx, userUploader, models.AffiliationAIRFORCE,
			models.Move{
				Locator: "SidDLH",
			},
			models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		)

		createHHGWithPaymentServiceItems(appCtx, primeUploader, moveRouter, shipmentFetcher)
		createHHGWithOriginSITServiceItems(appCtx, primeUploader, moveRouter, shipmentFetcher)
		createHHGWithDestinationSITServiceItems(appCtx, primeUploader, moveRouter, shipmentFetcher)
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
		createRandomMove(appCtx, nil, allDutyLocations, originDutyLocationsInGBLOC, true,
			userUploader,
			models.Move{
				Status:             models.MoveStatusAPPROVED,
				Locator:            "APRDVS",
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
			models.MTOShipment{
				Diversion:           true,
				Status:              models.MTOShipmentStatusApproved,
				ApprovedDate:        models.TimePointer(time.Now()),
				ScheduledPickupDate: models.TimePointer(time.Now().AddDate(0, 3, 0)),
			},
			models.Order{},
			models.ServiceMember{},
		)
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
		createReweighWithShipmentEDIErrorPaymentRequest(appCtx, userUploader, primeUploader, moveRouter)
		createReweighWithMixedShipmentStatuses(appCtx, userUploader)
	}
}

func subScenarioSITExtensions(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) func() {
	return func() {
		createTOO(appCtx)
		createMoveWithSITExtensionHistory(appCtx, userUploader)
		createMoveWithFutureSIT(appCtx, userUploader)
		createMoveWithAllPendingTOOActions(appCtx, userUploader, primeUploader)
	}
}

// Create moves with shipment address update requests in each of the three possible states: requested, approved, and rejected
func subScenarioShipmentAddressUpdates(appCtx appcontext.AppContext) func() {
	return func() {
		createTOO(appCtx)

		// Create move CRQST1 with a shipment address update request in requested state
		factory.BuildShipmentAddressUpdate(appCtx.DB(), []factory.Customization{}, []factory.Trait{factory.GetTraitShipmentAddressUpdateRequested})

		// Create move CRQST2 with a shipment address update request in approved state
		factory.BuildShipmentAddressUpdate(appCtx.DB(), []factory.Customization{}, []factory.Trait{factory.GetTraitShipmentAddressUpdateApproved})

		// Create move CRQST3 with a shipment address update request in rejected state
		factory.BuildShipmentAddressUpdate(appCtx.DB(), []factory.Customization{}, []factory.Trait{factory.GetTraitShipmentAddressUpdateRejected})
	}
}

func subScenarioNTSandNTSR(
	appCtx appcontext.AppContext,
	userUploader *uploader.UserUploader,
	moveRouter services.MoveRouter,
	shipmentFetcher services.MTOShipmentFetcher,
) func() {
	return func() {
		pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION

		createTXO(appCtx)
		createTXOServicesCounselor(appCtx)

		createNeedsServicesCounselingSingleHHG(appCtx, pcos, "NTSHHG")
		createNeedsServicesCounselingSingleHHG(appCtx, pcos, "NTSRHG")
		CreateNeedsServicesCounselingMinimalNTSR(appCtx, pcos, "NTSRMN")

		// Create a move with an HHG and NTS prime-handled shipment
		CreateMoveWithHHGAndNTSShipments(appCtx, "PRINTS", false)

		// Create a move with an HHG and NTS external vendor-handled shipment
		CreateMoveWithHHGAndNTSShipments(appCtx, "PRXNTS", true)

		// Create a move with only NTS external vendor-handled shipment
		CreateMoveWithNTSShipment(appCtx, "EXTNTS", true)

		// Create a move with only an NTS external vendor-handled shipment
		CreateMoveWithNTSShipment(appCtx, "NTSNTS", true)

		// Create a move with an HHG and NTS-release prime-handled shipment
		CreateMoveWithHHGAndNTSRShipments(appCtx, "PRINTR", false)

		// Create a move with an HHG and NTS-release external vendor-handled shipment
		CreateMoveWithHHGAndNTSRShipments(appCtx, "PRXNTR", true)

		// Create a move with only an NTS-release external vendor-handled shipment
		createMoveWithNTSRShipment(appCtx, "EXTNTR", true)

		// Create some submitted Moves for TXO users
		createMoveWithHHGAndNTSRMissingInfo(appCtx, moveRouter, shipmentFetcher)
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

		createMoveWithOriginAndDestinationSIT(appCtx, userUploader, "S1TT3R")
		createPaymentRequestsWithPartialSITInvoice(appCtx, primeUploader)
	}
}

func subScenarioPrimeUserAndClientCert(appCtx appcontext.AppContext) func() {
	return func() {
		primeUser := createPrimeUser(appCtx)
		createDevClientCertForUser(appCtx, primeUser)
	}
}

func subScenarioMultipleMoves(appCtx appcontext.AppContext) func() {
	return func() {
		createMultipleMovesTwoMovesHHGAndPPMShipments(appCtx)
		createMultipleMovesThreeMovesHHGPPMNTSShipments(appCtx)
		createMultipleMovesThreeMovesNTSHHGShipments(appCtx)
		createMultipleMovesThreeMovesPPMShipments(appCtx)
	}
}

// Transcom Relational Database Management (TRDM) TGET data
// Active and linked together transportation accounting code and line of accounting
// Creates a LOA and TAC that are active within a date range of 1 year
func createTGETLineOfAccountingAndTransportationAccountingCodeWithActiveDates(appCtx appcontext.AppContext) {
	ordersIssueDate := time.Now()
	startDate := ordersIssueDate.AddDate(-1, 0, 0)
	endDate := ordersIssueDate.AddDate(1, 0, 0)
	tacCode := "GOOD"

	loa := factory.BuildLineOfAccounting(appCtx.DB(), []factory.Customization{
		{
			Model: models.LineOfAccounting{
				LoaBgnDt:               &startDate,
				LoaEndDt:               &endDate,
				LoaSysID:               models.StringPointer("1234567890"),
				LoaHsGdsCd:             models.StringPointer(models.LineOfAccountingHouseholdGoodsCodeOfficer),
				LoaDptID:               models.StringPointer("1"),
				LoaTnsfrDptNm:          models.StringPointer("1"),
				LoaBafID:               models.StringPointer("1"),
				LoaTrsySfxTx:           models.StringPointer("1"),
				LoaMajClmNm:            models.StringPointer("1"),
				LoaOpAgncyID:           models.StringPointer("1"),
				LoaAlltSnID:            models.StringPointer("1"),
				LoaPgmElmntID:          models.StringPointer("1"),
				LoaTskBdgtSblnTx:       models.StringPointer("1"),
				LoaDfAgncyAlctnRcpntID: models.StringPointer("1"),
				LoaJbOrdNm:             models.StringPointer("1"),
				LoaSbaltmtRcpntID:      models.StringPointer("1"),
				LoaWkCntrRcpntNm:       models.StringPointer("1"),
				LoaMajRmbsmtSrcID:      models.StringPointer("1"),
				LoaDtlRmbsmtSrcID:      models.StringPointer("1"),
				LoaCustNm:              models.StringPointer("1"),
				LoaObjClsID:            models.StringPointer("1"),
				LoaSrvSrcID:            models.StringPointer("1"),
				LoaSpclIntrID:          models.StringPointer("1"),
				LoaBdgtAcntClsNm:       models.StringPointer("1"),
				LoaDocID:               models.StringPointer("1"),
				LoaClsRefID:            models.StringPointer("1"),
				LoaInstlAcntgActID:     models.StringPointer("1"),
				LoaLclInstlID:          models.StringPointer("1"),
				LoaFmsTrnsactnID:       models.StringPointer("1"),
				LoaTrnsnID:             models.StringPointer("1"),
				LoaUic:                 models.StringPointer("1"),
				LoaBgFyTx:              models.IntPointer(2023),
				LoaEndFyTx:             models.IntPointer(2025),
			},
		},
	}, nil)
	factory.BuildTransportationAccountingCodeWithoutAttachedLoa(appCtx.DB(), []factory.Customization{
		{
			Model: models.TransportationAccountingCode{
				TAC:               tacCode,
				TrnsprtnAcntBgnDt: &startDate,
				TrnsprtnAcntEndDt: &endDate,
				TacFnBlModCd:      models.StringPointer("1"),
				LoaSysID:          loa.LoaSysID,
			},
		},
	}, nil)
}
func subScenarioTGET(appCtx appcontext.AppContext) func() {
	return func() {
		createTGETLineOfAccountingAndTransportationAccountingCodeWithActiveDates(appCtx)
	}
}
