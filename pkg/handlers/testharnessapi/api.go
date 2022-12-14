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

func buildDefault(appCtx appcontext.AppContext, r *http.Request) (testHarnessResponse, error) {

	params := mux.Vars(r)
	action := params["action"]

	var response interface{}
	switch action {
	case "DefaultAdminUser":
		response = factory.BuildDefaultAdminUser(appCtx.DB())
	case "DefaultMove":
		response = testdatagen.MakeDefaultMove(appCtx.DB())
	case "SpouseProGearMove":
		response = testharness.MakeSpouseProGearMove(appCtx.DB())
	case "NeedsOrdersUser":
		response = testharness.MakeNeedsOrdersUser(appCtx.DB())
	default:
		return nil, errors.New("Cannot find builder for action: `" + action + "`")
	}

	return response, nil
}
