package primeapiv2

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/cli"
	mtoserviceitemops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/primev2messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi"
	primeapipayloads "github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/featureflag"
)

// CreateableServiceItemMap is a map of MTOServiceItemModelTypes and their allowed statuses
// THIS WILL NEED TO BE UPDATED AS WE CONTINUE TO ADD MORE SERVICE ITEMS.
// We will eventually remove this when all service items are added.
var CreateableServiceItemMap = map[primev2messages.MTOServiceItemModelType]bool{
	primev2messages.MTOServiceItemModelTypeMTOServiceItemOriginSIT:            true,
	primev2messages.MTOServiceItemModelTypeMTOServiceItemDestSIT:              true,
	primev2messages.MTOServiceItemModelTypeMTOServiceItemShuttle:              true,
	primev2messages.MTOServiceItemModelTypeMTOServiceItemDomesticCrating:      true,
	primev2messages.MTOServiceItemModelTypeMTOServiceItemInternationalCrating: true,
}

// CreateMTOServiceItemHandler is the handler to create MTO service items
type CreateMTOServiceItemHandler struct {
	handlers.HandlerConfig
	mtoServiceItemCreator  services.MTOServiceItemCreator
	mtoAvailabilityChecker services.MoveTaskOrderChecker
}

// Handle handler that creates a mto service item
func (h CreateMTOServiceItemHandler) Handle(params mtoserviceitemops.CreateMTOServiceItemParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			/** Feature Flag - Alaska **/
			isAlaskaEnabled := false
			featureFlagName := "enable_alaska"
			flag, err := h.FeatureFlagFetcher().GetBooleanFlag(params.HTTPRequest.Context(), appCtx.Logger(), "", featureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
			} else {
				isAlaskaEnabled = flag.Match
			}

			/** Turn on/off international crating/uncrating service items **/
			if !isAlaskaEnabled {
				delete(CreateableServiceItemMap, primev2messages.MTOServiceItemModelTypeMTOServiceItemInternationalCrating)
			}

			// restrict creation to a list
			if _, ok := primeapi.CreateableServiceItemMap[params.Body.ModelType()]; !ok {
				// throw error if modelType() not on the list
				mapKeys := primeapi.GetMapKeys(primeapi.CreateableServiceItemMap)
				detailErr := fmt.Sprintf("MTOServiceItem modelType() not allowed: %s ", params.Body.ModelType())
				verrs := validate.NewErrors()
				verrs.Add("modelType", fmt.Sprintf("allowed modelType() %v", mapKeys))

				appCtx.Logger().Error("primeapi.CreateMTOServiceItemHandler error", zap.Error(verrs))
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(primeapipayloads.ValidationError(
					detailErr, h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
			}

			// validation errors passed back if any
			mtoServiceItem, verrs := primeapipayloads.MTOServiceItemModel(params.Body)

			if verrs != nil && verrs.HasAny() {
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(primeapipayloads.ValidationError(
					"Invalid input found in service item", h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
			} else if mtoServiceItem == nil {
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(
					primeapipayloads.ValidationError("Unable to process service item", h.GetTraceIDFromRequest(params.HTTPRequest), nil)), verrs
			}

			moveTaskOrderID := uuid.FromStringOrNil(mtoServiceItem.MoveTaskOrderID.String())
			mtoAvailableToPrime, err := h.mtoAvailabilityChecker.MTOAvailableToPrime(appCtx, moveTaskOrderID)
			var mtoServiceItems *models.MTOServiceItems

			if mtoAvailableToPrime {
				v := viper.New()
				featureFlagFetcher, err := featureflag.NewFeatureFlagFetcher(cli.GetFliptFetcherConfig(v))
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("Error setting up feature flag fetcher: %s", err))
					return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(primeapipayloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}

				featureFlagValues, err := handlers.GetAllDomesticMHFlags(appCtx, featureFlagFetcher)
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("Error fetching mobile home feature flags: %s", err))
					return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(primeapipayloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
				mtoServiceItem.Status = models.MTOServiceItemStatusSubmitted
				mtoServiceItems, verrs, err = h.mtoServiceItemCreator.CreateMTOServiceItem(appCtx, mtoServiceItem, featureFlagValues)
				if err != nil {
					appCtx.Logger().Error(fmt.Sprintf("Error creating MTO Service Item: %s", err))
					return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(primeapipayloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			} else if err == nil {
				primeErr := apperror.NewNotFoundError(moveTaskOrderID, "primeapi.CreateMTOServiceItemHandler error - MTO is not available to Prime")
				appCtx.Logger().Error(primeErr.Error())
				return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(primeapipayloads.ClientError(
					handlers.NotFoundMessage, fmt.Sprintf("id: %s not found for moveTaskOrder", moveTaskOrderID), h.GetTraceIDFromRequest(params.HTTPRequest))), primeErr
			}

			if verrs != nil && verrs.HasAny() {
				return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(primeapipayloads.ValidationError(
					verrs.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
			}

			// Could be the error from MTOAvailableToPrime or CreateMTOServiceItem:
			if err != nil {
				appCtx.Logger().Error("primeapi.CreateMTOServiceItemHandler error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoserviceitemops.NewCreateMTOServiceItemNotFound().WithPayload(primeapipayloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					return mtoserviceitemops.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(primeapipayloads.ValidationError(e.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.ConflictError:
					return mtoserviceitemops.NewCreateMTOServiceItemConflict().WithPayload(primeapipayloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("primeapi.CreateMTOServiceItemHandler query error", zap.Error(e.Unwrap()))
					}
					return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(primeapipayloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoserviceitemops.NewCreateMTOServiceItemInternalServerError().WithPayload(primeapipayloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			mtoServiceItemsPayload := *primeapipayloads.MTOServiceItems(mtoServiceItems)
			return mtoserviceitemops.NewCreateMTOServiceItemOK().WithPayload(mtoServiceItemsPayload), nil
		})
}
