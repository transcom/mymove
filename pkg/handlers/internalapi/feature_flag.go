package internalapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	ffop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/feature_flags"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
)

// ShowLoggedInUserHandler returns the logged in user
type FeatureFlagsForUserHandler struct {
	handlers.HandlerConfig
}

// Handle returns the logged in user
func (h FeatureFlagsForUserHandler) Handle(params ffop.FeatureFlagForUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			flag, err := h.FeatureFlagFetcher().GetFlagForUser(
				params.HTTPRequest.Context(), appCtx, params.Key, params.FlagContext)

			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			flagPayload := internalmessages.FeatureFlag{
				Entity:    &flag.Entity,
				Key:       &flag.Key,
				Enabled:   &flag.Enabled,
				Value:     &flag.Value,
				Namespace: &flag.Namespace,
			}
			return ffop.NewFeatureFlagForUserOK().WithPayload(&flagPayload), nil
		})
}
