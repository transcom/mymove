package testharnessapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/testharness"
)

type InternalServerError struct {
	// The error string
	//
	// Required: true
	Error string `json:"error"`
}

type BaseTestHarnessHandler struct {
	handlers.HandlerConfig
}

type testHarnessResponse interface{}

func NewDefaultBuilder(handlerConfig handlers.HandlerConfig) http.Handler {
	return handlerConfig.AuditableAppContextFromRequestBasicHandler(
		func(appCtx appcontext.AppContext, w http.ResponseWriter, r *http.Request) error {
			response, err := buildDefault(appCtx, r)
			if err != nil {
				appCtx.Logger().Error("Testharness error", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				response = InternalServerError{
					Error: err.Error(),
				}
			}

			w.Header().Set("content-type", "application/json")
			return json.NewEncoder(w).Encode(response)
		})
}

type actionFunc func(appCtx appcontext.AppContext) testHarnessResponse

var actionDispatcher = map[string]actionFunc{
	"DefaultAdminUser": func(appCtx appcontext.AppContext) testHarnessResponse {
		return factory.BuildDefaultAdminUser(appCtx.DB())
	},
	"DefaultMove": func(appCtx appcontext.AppContext) testHarnessResponse {
		return testdatagen.MakeDefaultMove(appCtx.DB())
	},
	"MoveWithOrders": func(appCtx appcontext.AppContext) testHarnessResponse {
		return testharness.MakeMoveWithOrders(appCtx.DB())
	},
	"SpouseProGearMove": func(appCtx appcontext.AppContext) testHarnessResponse {
		return testharness.MakeSpouseProGearMove(appCtx.DB())
	},
	"WithShipmentMove": func(appCtx appcontext.AppContext) testHarnessResponse {
		return testharness.MakeWithShipmentMove(appCtx)
	},
	"HHGMoveWithServiceItemsAndPaymentRequestsAndFiles": func(appCtx appcontext.AppContext) testHarnessResponse {
		return testharness.MakeHHGMoveWithServiceItemsAndPaymentRequestsAndFiles(appCtx)
	},
	"PrimeSimulatorMoveNeedsShipmentUpdate": func(appCtx appcontext.AppContext) testHarnessResponse {
		return testharness.MakePrimeSimulatorMoveNeedsShipmentUpdate(appCtx)
	},
	"NeedsOrdersUser": func(appCtx appcontext.AppContext) testHarnessResponse {
		return testharness.MakeNeedsOrdersUser(appCtx.DB())
	},
	"PPMInProgressMove": func(appCtx appcontext.AppContext) testHarnessResponse {
		return testharness.MakePPMInProgressMove(appCtx)
	},
	"OfficeUserWithTOOAndTIO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return testharness.MakeOfficeUserWithTOOAndTIO(appCtx)
	},
}

func buildDefault(appCtx appcontext.AppContext, r *http.Request) (testHarnessResponse, error) {

	params := mux.Vars(r)
	action := params["action"]

	dispatcher, ok := actionDispatcher[action]
	if !ok {
		return nil, errors.New("Cannot find builder for action: `" + action + "`")
	}

	return dispatcher(appCtx), nil
}
