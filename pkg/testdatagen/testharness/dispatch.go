package testharness

import (
	"errors"
	"sort"
	"sync"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type testHarnessResponse interface{}

type actionFunc func(appCtx appcontext.AppContext) testHarnessResponse

var actionDispatcher = map[string]actionFunc{
	"DefaultAdminUser": func(appCtx appcontext.AppContext) testHarnessResponse {
		return factory.BuildDefaultAdminUser(appCtx.DB())
	},
	"SuperAdminUser": func(appCtx appcontext.AppContext) testHarnessResponse {
		return factory.BuildDefaultSuperAdminUser(appCtx.DB())
	},
	"DefaultMove": func(appCtx appcontext.AppContext) testHarnessResponse {
		return factory.BuildMove(appCtx.DB(), nil, nil)
	},
	"MoveWithOrders": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeMoveWithOrders(appCtx.DB())
	},
	"SpouseProGearMove": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeSpouseProGearMove(appCtx.DB())
	},
	"WithShipmentMove": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeWithShipmentMove(appCtx)
	},
	"HHGMoveWithNTSAndNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithNTSAndNeedsSC(appCtx)
	},
	"MobileHomeMoveNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeMobileHomeMoveNeedsSC(appCtx)
	},
	"GoodTACAndLoaCombination": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeGoodTACAndLoaCombination(appCtx)
	},
	"MoveWithMinimalNTSRNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeMoveWithMinimalNTSRNeedsSC(appCtx)
	},
	"HHGMoveNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveNeedsSC(appCtx)
	},
	"IntlHHGMoveNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveNeedsSC(appCtx)
	},
	"HHGMoveNeedsSCOtherGBLOC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveNeedsSCOtherGBLOC(appCtx)
	},
	"HHGMoveAsUSMCNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveNeedsServicesCounselingUSMC(appCtx)
	},
	"HHGMoveWithAmendedOrders": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithAmendedOrders(appCtx)
	},
	"HHGMoveForSeparationNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveForSeparationNeedsSC(appCtx)
	},
	"HHGMoveForRetireeNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveForRetireeNeedsSC(appCtx)
	},
	"HHGMoveInSIT": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSIT(appCtx)
	},
	"HHGMoveWithPastSITs": func(appCtx appcontext.AppContext) testHarnessResponse {
		return HHGMoveWithPastSITs(appCtx)
	},
	"HHGMoveInSITNoDestinationSITOutDate": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSITNoDestinationSITOutDate(appCtx)
	},
	"HHGMoveInSITNoExcessWeight": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSITNoExcessWeight(appCtx)
	},
	"HHGMoveInSITWithPendingExtension": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSITWithPendingExtension(appCtx)
	},
	"HHGMoveInSITEndsToday": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSITEndsToday(appCtx)
	},
	"HHGMoveInSITEndsTomorrow": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSITEndsTomorrow(appCtx)
	},
	"HHGMoveInSITEndsYesterday": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSITEndsYesterday(appCtx)
	},
	"HHGMoveInSITDeparted": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSITDeparted(appCtx)
	},
	"HHGMoveInSITStartsInFuture": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSITStartsInFuture(appCtx)
	},
	"HHGMoveInSITNotApproved": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSITNotApproved(appCtx)
	},
	"HHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO(appCtx)
	},
	"HHGMoveInTerminatedStatus": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInTerminatedStatus(appCtx)
	},
	"HHGMoveWithIntlCratingServiceItemsTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithIntlCratingServiceItemsTOO(appCtx)
	},
	"HHGMoveWithIntlShuttleServiceItemsTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithIntlShuttleServiceItemsTOO(appCtx)
	},
	"HHGMoveForTOOAfterActualPickupDate": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveForTOOAfterActualPickupDate(appCtx)
	},
	"HHGMoveWithRetireeForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithRetireeForTOO(appCtx)
	},
	"HHGMoveWithNTSShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithNTSShipmentsForTOO(appCtx)
	},
	"MoveWithNTSShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeMoveWithNTSShipmentsForTOO(appCtx)
	},
	"HHGMoveWithPPMShipmentsReadyForCloseout": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithPPMShipmentsForTOO(appCtx, true)
	},
	"HHGMoveWithPPMShipmentsReadyForCounseling": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithPPMShipmentsForTOO(appCtx, false)
	},
	"HHGMoveWithExternalNTSShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithExternalNTSShipmentsForTOO(appCtx)
	},
	"HHGMoveWithApprovedNTSShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithApprovedNTSShipmentsForTOO(appCtx)
	},
	"HHGMoveWithNTSRShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithNTSRShipmentsForTOO(appCtx)
	},
	"HHGMoveWithApprovedNTSRShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithApprovedNTSRShipmentsForTOO(appCtx)
	},
	"HHGMoveWithExternalNTSRShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithExternalNTSRShipmentsForTOO(appCtx)
	},
	"HHGMoveWithServiceItemsandPaymentRequestsForTIO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithServiceItemsandPaymentRequestsForTIO(appCtx)
	},
	"HHGMoveWithServiceItemsandPaymentRequestReviewedForQAE": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithServiceItemsandPaymentRequestReviewedForQAE(appCtx)
	},
	"HHGMoveWithServiceItemsandPaymentRequestWithDocsReviewedForQAE": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithServiceItemsandPaymentRequestWithDocsReviewedForQAE(appCtx)
	},
	"HHGMoveInSITWithAddressChangeRequestOver50Miles": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSITWithAddressChangeRequestOver50Miles(appCtx)
	},
	"HHGMoveInSITWithAddressChangeRequestUnder50Miles": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSITWithAddressChangeRequestUnder50Miles(appCtx)
	},
	"NTSRMoveWithPaymentRequest": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeNTSRMoveWithPaymentRequest(appCtx)
	},
	"NTSRMoveWithServiceItemsAndPaymentRequest": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeNTSRMoveWithServiceItemsAndPaymentRequest(appCtx)
	},
	"PrimeSimulatorMoveNeedsShipmentUpdate": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakePrimeSimulatorMoveNeedsShipmentUpdate(appCtx)
	},
	"MakePrimeSimulatorMoveSameBasePointCity": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakePrimeSimulatorMoveSameBasePointCity(appCtx)
	},
	"NeedsOrdersUser": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeNeedsOrdersUser(appCtx.DB())
	},
	"MoveWithPPMShipmentReadyForFinalCloseout": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeMoveWithPPMShipmentReadyForFinalCloseout(appCtx)
	},
	"MoveWithPPMShipmentReadyForFinalCloseoutWithSIT": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeMoveWithPPMShipmentReadyForFinalCloseoutWithSIT(appCtx)
	},
	"PPMMoveWithCloseout": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakePPMMoveWithCloseout(appCtx)
	},
	"PPMMoveWithCloseoutOffice": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakePPMMoveWithCloseoutOffice(appCtx)
	},
	"ApprovedMoveWithPPM": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPM(appCtx)
	},
	"SubmittedMoveWithPPMShipmentForSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeSubmittedMoveWithPPMShipmentForSC(appCtx)
	},
	"ApprovedMoveWithSubmittedPPMShipmentForSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithSubmittedPPMShipmentForSC(appCtx)
	},
	"UnSubmittedMoveWithPPMShipmentThroughEstimatedWeights": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights(appCtx)
	},
	"ApprovedMoveWithPPMWithAboutFormComplete": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMWithAboutFormComplete(appCtx)
	},
	"UnsubmittedMoveWithMultipleFullPPMShipmentComplete": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeUnsubmittedMoveWithMultipleFullPPMShipmentComplete(appCtx)
	},
	"ApprovedMoveWithPPMProgearWeightTicket": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMProgearWeightTicket(appCtx)
	},
	"ApprovedMoveWithPPMProgearWeightTicketOffice": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMProgearWeightTicketOffice(appCtx)
	},
	"ApprovedMoveWithPPMWithMultipleProgearWeightTicketsOffice": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMWithMultipleProgearWeightTicketsOffice(appCtx)
	},
	"ApprovedMoveWithPPMWithMultipleProgearWeightTicketsOffice2": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMWithMultipleProgearWeightTicketsOffice2(appCtx)
	},
	"ApprovedMoveWithPPMProgearWeightTicketOfficeCivilian": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMProgearWeightTicketOfficeCivilian(appCtx)
	},
	"ApprovedMoveWithPPMGunSafeWeightTicketOffice": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMGunSafeWeightTicketOffice(appCtx)
	},
	"ApprovedMoveWithPPMWeightTicketOffice": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMWeightTicketOffice(appCtx)
	},
	"ApprovedMoveWithPPMWeightTicketOfficeWithHHG": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMWeightTicketOfficeWithHHG(appCtx)
	},
	"ApprovedMoveWithPPMMovingExpense": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMMovingExpense(appCtx)
	},
	"ApprovedMoveWithPPMMovingExpenseOffice": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMMovingExpenseOffice(appCtx)
	},
	"ApprovedMoveWithPPMAllDocTypesOffice": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMAllDocTypesOffice(appCtx)
	},
	"DraftMoveWithPPMWithDepartureDate": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeDraftMoveWithPPMWithDepartureDate(appCtx)
	},
	"OfficeUserWithTOOAndTIO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeOfficeUserWithTOOAndTIO(appCtx)
	},
	"OfficeUserWithMultirole": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeOfficeUserWithMultirole(appCtx)
	},
	"RequestedOfficeUser": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeRequestedOfficeUserWithTOO(appCtx)
	},
	"RejectedOfficeUser": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeRejectedOfficeUserWithTOO(appCtx)
	},
	"WebhookSubscription": func(appCtx appcontext.AppContext) testHarnessResponse {
		return testdatagen.MakeWebhookSubscription(appCtx.DB(), testdatagen.Assertions{})
	},
	"ApprovedMoveWithPPMShipmentAndExcessWeight": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMShipmentAndExcessWeight(appCtx)
	},
	"HHGMoveWithAddressChangeRequest": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithAddressChangeRequest(appCtx)
	},
	"MakeHHGMoveWithAddressChangeRequestAndUnknownDeliveryAddress": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithAddressChangeRequestAndUnknownDeliveryAddress(appCtx)
	},
	"MakeHHGMoveWithAddressChangeRequestAndSecondDeliveryLocation": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithAddressChangeRequestAndSecondDeliveryLocation(appCtx)
	},
	"NTSRMoveWithAddressChangeRequest": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeNTSRMoveWithAddressChangeRequest(appCtx)
	},
	"MakeMoveReadyForEDI": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeMoveReadyForEDI(appCtx)
	},
	"MakeCoastGuardMoveReadyForEDI": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeCoastGuardMoveReadyForEDI(appCtx)
	},
	"BoatHaulAwayMoveNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeBoatHaulAwayMoveNeedsSC(appCtx)
	},
	"BoatHaulAwayMoveNeedsTOOApproval": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeBoatHaulAwayMoveNeedsTOOApproval(appCtx)
	},
	"OfficeUserWithCustomer": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeOfficeUserWithCustomer(appCtx)
	},
	"OfficeUserWithContractingOfficer": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeOfficeUserWithContractingOfficer(appCtx)
	},
	"OfficeUserWithPrimeSimulator": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeOfficeUserWithPrimeSimulator(appCtx)
	},
	"OfficeUserWithGSR": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeOfficeUserWithGSR(appCtx)
	},
	"InternationalAlaskaBasicHHGMoveForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeInternationalAlaskaBasicHHGMoveForTOO(appCtx)
	},
	"InternationalHHGIntoInternationalNTSMoveWithServiceItemsandPaymentRequestsForTIO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeBasicInternationalHHGMoveWithServiceItemsandPaymentRequestsForTIO(appCtx, true)
	},
	"InternationalHHGMoveWithServiceItemsandPaymentRequestsForTIO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeBasicInternationalHHGMoveWithServiceItemsandPaymentRequestsForTIO(appCtx, false)
	},
	"IntlHHGMoveWithCratingUncratingServiceItemsAndPaymentRequestsForTIO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveWithCratingUncratingServiceItemsAndPaymentRequestsForTIO(appCtx)
	},
	// basic iHHG move with CONUS -> AK needing TOO approval
	"IntlHHGMoveDestAKZone1Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAKZone1Army(appCtx)
	},
	"IntlHHGMoveDestAKZone2Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAKZone2Army(appCtx)
	},
	"IntlHHGMoveDestAKZone1AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAKZone1AirForce(appCtx)
	},
	"IntlHHGMoveDestAKZone2AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAKZone2AirForce(appCtx)
	},
	"IntlHHGMoveDestAKZone1SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAKZone1SpaceForce(appCtx)
	},
	"IntlHHGMoveDestAKZone2SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAKZone2SpaceForce(appCtx)
	},
	"IntlHHGMoveDestAKZone1USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAKZone1USMC(appCtx)
	},
	"IntlHHGMoveDestAKZone2USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAKZone2USMC(appCtx)
	},
	// basic iHHG move with AK -> CONUS needing TOO approval
	"IntlHHGMovePickupAKZone1Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMovePickupAKZone1Army(appCtx)
	},
	"IntlHHGMovePickupAKZone2Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMovePickupAKZone2Army(appCtx)
	},
	"IntlHHGMovePickupAKZone1AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMovePickupAKZone1AirForce(appCtx)
	},
	"IntlHHGMovePickupAKZone2AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMovePickupAKZone2AirForce(appCtx)
	},
	"IntlHHGMovePickupAKZone1SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMovePickupAKZone1SpaceForce(appCtx)
	},
	"IntlHHGMovePickupAKZone2SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMovePickupAKZone2SpaceForce(appCtx)
	},
	"IntlHHGMovePickupAKZone1USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMovePickupAKZone1USMC(appCtx)
	},
	"IntlHHGMovePickupAKZone2USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMovePickupAKZone2USMC(appCtx)
	},
	// iHHG with international origin SIT in SUBMITTED status
	"IntlHHGMoveOriginSITRequestedAKZone1Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITRequestedAKZone1Army(appCtx)
	},
	"IntlHHGMoveOriginSITRequestedAKZone2Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITRequestedAKZone2Army(appCtx)
	},
	"IntlHHGMoveOriginSITRequestedAKZone1AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITRequestedAKZone1AirForce(appCtx)
	},
	"IntlHHGMoveOriginSITRequestedAKZone2AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITRequestedAKZone2AirForce(appCtx)
	},
	"IntlHHGMoveOriginSITRequestedAKZone1SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITRequestedAKZone1SpaceForce(appCtx)
	},
	"IntlHHGMoveOriginSITRequestedAKZone2SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITRequestedAKZone2SpaceForce(appCtx)
	},
	"IntlHHGMoveOriginSITRequestedAKZone1USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITRequestedAKZone1USMC(appCtx)
	},
	"IntlHHGMoveOriginSITRequestedAKZone2USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITRequestedAKZone2USMC(appCtx)
	},
	// iHHG with international destination SIT in SUBMITTED status
	"IntlHHGMoveDestSITRequestedAKZone1Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITRequestedAKZone1Army(appCtx)
	},
	"IntlHHGMoveDestSITRequestedAKZone2Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITRequestedAKZone2Army(appCtx)
	},
	"IntlHHGMoveDestSITRequestedAKZone1AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITRequestedAKZone1AirForce(appCtx)
	},
	"IntlHHGMoveDestSITRequestedAKZone2AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITRequestedAKZone2AirForce(appCtx)
	},
	"IntlHHGMoveDestSITRequestedAKZone1SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITRequestedAKZone1SpaceForce(appCtx)
	},
	"IntlHHGMoveDestSITRequestedAKZone2SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITRequestedAKZone2SpaceForce(appCtx)
	},
	"IntlHHGMoveDestSITRequestedAKZone1USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITRequestedAKZone1USMC(appCtx)
	},
	"IntlHHGMoveDestSITRequestedAKZone2USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITRequestedAKZone2USMC(appCtx)
	},
	// iHHG with BOTH international origin & destination SIT in SUBMITTED status
	"IntlHHGMoveBothSITRequestedAKZone1Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothSITRequestedAKZone1Army(appCtx)
	},
	"IntlHHGMoveBothSITRequestedAKZone2Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothSITRequestedAKZone2Army(appCtx)
	},
	"IntlHHGMoveBothSITRequestedAKZone1AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothSITRequestedAKZone1AirForce(appCtx)
	},
	"IntlHHGMoveBothSITRequestedAKZone2AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothSITRequestedAKZone2AirForce(appCtx)
	},
	"IntlHHGMoveBothSITRequestedAKZone1SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothSITRequestedAKZone1SpaceForce(appCtx)
	},
	"IntlHHGMoveBothSITRequestedAKZone2SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothSITRequestedAKZone2SpaceForce(appCtx)
	},
	"IntlHHGMoveBothSITRequestedAKZone1USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothSITRequestedAKZone1USMC(appCtx)
	},
	"IntlHHGMoveBothSITRequestedAKZone2USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothSITRequestedAKZone2USMC(appCtx)
	},
	// iHHG with international origin shuttle in SUBMITTED status
	"IntlHHGMoveOriginShuttleRequestedAKZone1Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginShuttleRequestedAKZone1Army(appCtx)
	},
	"IntlHHGMoveOriginShuttleRequestedAKZone2Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginShuttleRequestedAKZone2Army(appCtx)
	},
	"IntlHHGMoveOriginShuttleRequestedAKZone1AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginShuttleRequestedAKZone1AirForce(appCtx)
	},
	"IntlHHGMoveOriginShuttleRequestedAKZone2AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginShuttleRequestedAKZone2AirForce(appCtx)
	},
	"IntlHHGMoveOriginShuttleRequestedAKZone1SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginShuttleRequestedAKZone1SpaceForce(appCtx)
	},
	"IntlHHGMoveOriginShuttleRequestedAKZone2SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginShuttleRequestedAKZone2SpaceForce(appCtx)
	},
	"IntlHHGMoveOriginShuttleRequestedAKZone1USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginShuttleRequestedAKZone1USMC(appCtx)
	},
	"IntlHHGMoveOriginShuttleRequestedAKZone2USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginShuttleRequestedAKZone2USMC(appCtx)
	},
	// iHHG with international destination shuttle in SUBMITTED status
	"IntlHHGMoveDestShuttleRequestedAKZone1Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestShuttleRequestedAKZone1Army(appCtx)
	},
	"IntlHHGMoveDestShuttleRequestedAKZone2Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestShuttleRequestedAKZone2Army(appCtx)
	},
	"IntlHHGMoveDestShuttleRequestedAKZone1AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestShuttleRequestedAKZone1AirForce(appCtx)
	},
	"IntlHHGMoveDestShuttleRequestedAKZone2AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestShuttleRequestedAKZone2AirForce(appCtx)
	},
	"IntlHHGMoveDestShuttleRequestedAKZone1SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestShuttleRequestedAKZone1SpaceForce(appCtx)
	},
	"IntlHHGMoveDestShuttleRequestedAKZone2SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestShuttleRequestedAKZone2SpaceForce(appCtx)
	},
	"IntlHHGMoveDestShuttleRequestedAKZone1USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestShuttleRequestedAKZone1USMC(appCtx)
	},
	"IntlHHGMoveDestShuttleRequestedAKZone2USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestShuttleRequestedAKZone2USMC(appCtx)
	},
	// iHHG with BOTH international origin & destination shuttle in SUBMITTED status
	"IntlHHGMoveBothShuttleRequestedAKZone1Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothShuttleRequestedAKZone1Army(appCtx)
	},
	"IntlHHGMoveBothShuttleRequestedAKZone2Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothShuttleRequestedAKZone2Army(appCtx)
	},
	"IntlHHGMoveBothShuttleRequestedAKZone1AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothShuttleRequestedAKZone1AirForce(appCtx)
	},
	"IntlHHGMoveBothShuttleRequestedAKZone2AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothShuttleRequestedAKZone2AirForce(appCtx)
	},
	"IntlHHGMoveBothShuttleRequestedAKZone1SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothShuttleRequestedAKZone1SpaceForce(appCtx)
	},
	"IntlHHGMoveBothShuttleRequestedAKZone2SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothShuttleRequestedAKZone2SpaceForce(appCtx)
	},
	"IntlHHGMoveBothShuttleRequestedAKZone1USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothShuttleRequestedAKZone1USMC(appCtx)
	},
	"IntlHHGMoveBothShuttleRequestedAKZone2USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveBothShuttleRequestedAKZone2USMC(appCtx)
	},
	// iHHG with a destination address request in REQUESTED status
	"IntlHHGMoveDestAddressRequestedAKZone1Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAddressRequestedAKZone1Army(appCtx)
	},
	"IntlHHGMoveDestAddressRequestedAKZone2Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAddressRequestedAKZone2Army(appCtx)
	},
	"IntlHHGMoveDestAddressRequestedAKZone1AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAddressRequestedAKZone1AirForce(appCtx)
	},
	"IntlHHGMoveDestAddressRequestedAKZone2AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAddressRequestedAKZone2AirForce(appCtx)
	},
	"IntlHHGMoveDestAddressRequestedAKZone1SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAddressRequestedAKZone1SpaceForce(appCtx)
	},
	"IntlHHGMoveDestAddressRequestedAKZone2SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAddressRequestedAKZone2SpaceForce(appCtx)
	},
	"IntlHHGMoveDestAddressRequestedAKZone1USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAddressRequestedAKZone1USMC(appCtx)
	},
	"IntlHHGMoveDestAddressRequestedAKZone2USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestAddressRequestedAKZone2USMC(appCtx)
	},
	// iHHG with a PENDING SIT extension request containing origin SIT
	"IntlHHGMoveOriginSITExtensionRequestedAKZone1Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITExtensionRequestedAKZone1Army(appCtx)
	},
	"IntlHHGMoveOriginSITExtensionRequestedAKZone2Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITExtensionRequestedAKZone2Army(appCtx)
	},
	"IntlHHGMoveOriginSITExtensionRequestedAKZone1AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITExtensionRequestedAKZone1AirForce(appCtx)
	},
	"IntlHHGMoveOriginSITExtensionRequestedAKZone2AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITExtensionRequestedAKZone2AirForce(appCtx)
	},
	"IntlHHGMoveOriginSITExtensionRequestedAKZone1SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITExtensionRequestedAKZone1SpaceForce(appCtx)
	},
	"IntlHHGMoveOriginSITExtensionRequestedAKZone2SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITExtensionRequestedAKZone2SpaceForce(appCtx)
	},
	"IntlHHGMoveOriginSITExtensionRequestedAKZone1USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITExtensionRequestedAKZone1USMC(appCtx)
	},
	"IntlHHGMoveOriginSITExtensionRequestedAKZone2USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveOriginSITExtensionRequestedAKZone2USMC(appCtx)
	},
	// iHHG with a PENDING SIT extension request containing destination SIT
	"IntlHHGMoveDestSITExtensionRequestedAKZone1Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITExtensionRequestedAKZone1Army(appCtx)
	},
	"IntlHHGMoveDestSITExtensionRequestedAKZone2Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITExtensionRequestedAKZone2Army(appCtx)
	},
	"IntlHHGMoveDestSITExtensionRequestedAKZone1AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITExtensionRequestedAKZone1AirForce(appCtx)
	},
	"IntlHHGMoveDestSITExtensionRequestedAKZone2AirForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITExtensionRequestedAKZone2AirForce(appCtx)
	},
	"IntlHHGMoveDestSITExtensionRequestedAKZone1SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITExtensionRequestedAKZone1SpaceForce(appCtx)
	},
	"IntlHHGMoveDestSITExtensionRequestedAKZone2SpaceForce": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITExtensionRequestedAKZone2SpaceForce(appCtx)
	},
	"IntlHHGMoveDestSITExtensionRequestedAKZone1USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITExtensionRequestedAKZone1USMC(appCtx)
	},
	"IntlHHGMoveDestSITExtensionRequestedAKZone2USMC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveDestSITExtensionRequestedAKZone2USMC(appCtx)
	},
	// iHHG with a PENDING excess weight notification
	"IntlHHGMoveExcessWeightAKZone1Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveExcessWeightAKZone1Army(appCtx)
	},
	"IntlHHGMoveExcessWeightAKZone2Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlHHGMoveExcessWeightAKZone2Army(appCtx)
	},
	// iUB with a PENDING excess UB weight notification
	"IntlUBMoveExcessWeightAKZone1Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlUBMoveExcessWeightAKZone1Army(appCtx)
	},
	"IntlUBMoveExcessWeightAKZone2Army": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeIntlUBMoveExcessWeightAKZone2Army(appCtx)
	},
}

func Actions() []string {
	actions := make([]string, 0, len(actionDispatcher))
	for k := range actionDispatcher {
		actions = append(actions, k)
	}
	sort.Strings(actions)
	return actions
}

var mutex sync.Mutex

func Dispatch(appCtx appcontext.AppContext, action string) (testHarnessResponse, error) {

	// ensure only one dispatch is running at a time in a heavy handed
	// way to prevent multiple setup functions from stomping on each
	// other when creating shared data (like duty locations)
	mutex.Lock()
	defer mutex.Unlock()

	dispatcher, ok := actionDispatcher[action]
	if !ok {
		appCtx.Logger().Error("Cannot find testharness dispatcher", zap.Any("action", action))
		return nil, errors.New("Cannot find testharness dispatcher for action: `" + action + "`")
	}

	appCtx.Logger().Info("Found testharness dispatcher", zap.Any("action", action))
	return dispatcher(appCtx), nil

}
