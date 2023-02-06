package testharnessapi

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/handlers"
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

func NewDefaultBuilder(handlerConfig handlers.HandlerConfig) http.Handler {
	return handlerConfig.AuditableAppContextFromRequestBasicHandler(
		func(appCtx appcontext.AppContext, w http.ResponseWriter, r *http.Request) error {
			params := mux.Vars(r)
			action := params["action"]

			response, err := testharness.Dispatch(appCtx, action)
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
