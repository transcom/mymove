package testharness

import (
	"errors"

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
	"HHGMoveWithServiceItemsAndPaymentRequestsAndFiles": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithServiceItemsAndPaymentRequestsAndFiles(appCtx)
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
		return nil, errors.New("Cannot find builder for action: `" + action + "`")
	}

	return dispatcher(appCtx), nil

}
