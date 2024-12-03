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
	"HHGMoveWithIntlCratingServiceItemsTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithIntlCratingServiceItemsTOO(appCtx)
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
	"ApprovedMoveWithPPMProgearWeightTicketOfficeCivilian": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMProgearWeightTicketOfficeCivilian(appCtx)
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
