package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ffop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/feature_flags"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
)

// BooleanFeatureFlagsUnauthenticatedHandler handles evaluating boolean feature flags outside of authentication
type BooleanFeatureFlagsUnauthenticatedHandler struct {
	handlers.HandlerConfig
}

// Handle returns the boolean feature flag for an unauthenticated user
func (h BooleanFeatureFlagsUnauthenticatedHandler) Handle(params ffop.BooleanFeatureFlagUnauthenticatedParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			// we are only allowing this to be called from the office app
			// since this is an open route outside of auth, we want to buckle down on validation here
			if !appCtx.Session().IsOfficeApp() {
				return ffop.NewBooleanFeatureFlagUnauthenticatedUnauthorized(), apperror.NewSessionError("Request is not from the office app")
			}
			flag, err := h.FeatureFlagFetcher().GetBooleanFlag(
				params.HTTPRequest.Context(), appCtx.Logger(), "office", params.Key, params.FlagContext)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			flagPayload := ghcmessages.FeatureFlagBoolean{
				Entity:    &flag.Entity,
				Key:       &params.Key,
				Match:     &flag.Match,
				Namespace: &flag.Namespace,
			}
			return ffop.NewBooleanFeatureFlagUnauthenticatedOK().WithPayload(&flagPayload), nil
		})
}
