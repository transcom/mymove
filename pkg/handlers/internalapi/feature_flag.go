package internalapi

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ffop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/feature_flags"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
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
			// we are only allowing this to be called from the customer app
			// since this is an open route outside of auth, we want to buckle down on validation here
			if !appCtx.Session().IsMilApp() && !appCtx.Session().IsOfficeApp() {
				return ffop.NewBooleanFeatureFlagUnauthenticatedUnauthorized(), apperror.NewSessionError("Request is not coming from the customer app or office app")
			}
			flag, err := h.FeatureFlagFetcher().GetBooleanFlag(
				params.HTTPRequest.Context(), appCtx.Logger(), "customer", params.Key, params.FlagContext)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			flagPayload := internalmessages.FeatureFlagBoolean{
				Entity:    &flag.Entity,
				Key:       &params.Key,
				Match:     &flag.Match,
				Namespace: &flag.Namespace,
			}
			return ffop.NewBooleanFeatureFlagUnauthenticatedOK().WithPayload(&flagPayload), nil
		})
}

// BooleanFeatureFlagsForUserHandler handles evaluating boolean feature flags for
// users
type BooleanFeatureFlagsForUserHandler struct {
	handlers.HandlerConfig
}

// Handle returns the boolean feature flag
func (h BooleanFeatureFlagsForUserHandler) Handle(params ffop.BooleanFeatureFlagForUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			flag, err := h.FeatureFlagFetcher().GetBooleanFlagForUser(
				params.HTTPRequest.Context(), appCtx, params.Key, params.FlagContext)

			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			flagPayload := internalmessages.FeatureFlagBoolean{
				Entity:    &flag.Entity,
				Key:       &params.Key,
				Match:     &flag.Match,
				Namespace: &flag.Namespace,
			}
			return ffop.NewBooleanFeatureFlagForUserOK().WithPayload(&flagPayload), nil
		})
}

// VariantFeatureFlagsForUserHandler handles evaluating variant feature flags for
// users
type VariantFeatureFlagsForUserHandler struct {
	handlers.HandlerConfig
}

// Handle returns the boolean feature flag
func (h VariantFeatureFlagsForUserHandler) Handle(params ffop.VariantFeatureFlagForUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			flag, err := h.FeatureFlagFetcher().GetVariantFlagForUser(
				params.HTTPRequest.Context(), appCtx, params.Key, params.FlagContext)

			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			flagPayload := internalmessages.FeatureFlagVariant{
				Entity:    &flag.Entity,
				Key:       &params.Key,
				Match:     &flag.Match,
				Variant:   &flag.Variant,
				Namespace: &flag.Namespace,
			}
			return ffop.NewVariantFeatureFlagForUserOK().WithPayload(&flagPayload), nil
		})
}
