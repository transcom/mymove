package testharness

import (
	"errors"

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
		return testdatagen.MakeDefaultMove(appCtx.DB())
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
	"HHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO(appCtx)
	},
	"HHGMoveWithRetireeForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithRetireeForTOO(appCtx)
	},
	"HHGMoveWithServiceItemsandPaymentRequestsForTIO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithServiceItemsandPaymentRequestsForTIO(appCtx)
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
	"NeedsOrdersUser": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeNeedsOrdersUser(appCtx.DB())
	},
	"PPMInProgressMove": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakePPMInProgressMove(appCtx)
	},
	"OfficeUserWithTOOAndTIO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeOfficeUserWithTOOAndTIO(appCtx)
	},
	"WebhookSubscription": func(appCtx appcontext.AppContext) testHarnessResponse {
		return testdatagen.MakeWebhookSubscription(appCtx.DB(), testdatagen.Assertions{})
	},
}

func Dispatch(appCtx appcontext.AppContext, action string) (testHarnessResponse, error) {
	dispatcher, ok := actionDispatcher[action]
	if !ok {
		appCtx.Logger().Error("Cannot find testharness dispatcher", zap.Any("action", action))
		return nil, errors.New("Cannot find testharness dispatcher for action: `" + action + "`")
	}

	appCtx.Logger().Info("Found testharness dispatcher", zap.Any("action", action))
	return dispatcher(appCtx), nil

}
