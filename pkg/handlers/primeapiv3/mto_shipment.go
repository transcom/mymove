package primeapiv3

import (
	"context"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primev3api/primev3operations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primev3messages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi"
	"github.com/transcom/mymove/pkg/handlers/primeapiv3/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
)

// CreateMTOShipmentHandler is the handler to create MTO shipments
type CreateMTOShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentCreator
	mtoAvailabilityChecker services.MoveTaskOrderChecker
	services.VLocation
}

// Handle creates the mto shipment
func (h CreateMTOShipmentHandler) Handle(params mtoshipmentops.CreateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.Body
			if payload == nil {
				err := apperror.NewBadDataError("the MTO Shipment request body cannot be empty")
				appCtx.Logger().Error(err.Error())
				return mtoshipmentops.NewCreateMTOShipmentBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			/** Feature Flag - Boat Shipment **/
			const featureFlagName = "boat"
			isBoatFeatureOn := false
			flag, err := h.FeatureFlagFetcher().GetBooleanFlag(params.HTTPRequest.Context(), appCtx.Logger(), "", featureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
			} else {
				isBoatFeatureOn = flag.Match
			}

			// Return an error if boat shipment is sent while the feature flag is turned off.
			if !isBoatFeatureOn && (*params.Body.ShipmentType == primev3messages.MTOShipmentTypeBOATHAULAWAY || *params.Body.ShipmentType == primev3messages.MTOShipmentTypeBOATTOWAWAY) {
				return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(
					"Boat shipment type was used but the feature flag is not enabled.", h.GetTraceIDFromRequest(params.HTTPRequest), nil)), nil
			}

			/** Feature Flag - Mobile Home Shipment **/
			const featureFlagMobileHome = "mobile_home"
			isMobileHomeFeatureOn := false
			flagMH, err := h.FeatureFlagFetcher().GetBooleanFlag(params.HTTPRequest.Context(), appCtx.Logger(), "", featureFlagMobileHome, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flagMH", zap.String("featureFlagKey", featureFlagMobileHome), zap.Error(err))
			} else {
				isMobileHomeFeatureOn = flagMH.Match
			}

			// Return an error if mobile home shipment is sent while the feature flag is turned off.
			if !isMobileHomeFeatureOn && (*params.Body.ShipmentType == primev3messages.MTOShipmentTypeMOBILEHOME) {
				return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(
					"Mobile Home shipment type was used but the feature flag is not enabled.", h.GetTraceIDFromRequest(params.HTTPRequest), nil)), nil
			}

			/** Feature Flag - UB Shipment **/
			const featureFlagNameUB = "unaccompanied_baggage"
			isUBFeatureOn := false
			flag, err = h.FeatureFlagFetcher().GetBooleanFlag(params.HTTPRequest.Context(), appCtx.Logger(), "", featureFlagNameUB, map[string]string{})

			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagNameUB), zap.Error(err))
			} else {
				isUBFeatureOn = flag.Match
			}

			// Return an error if UB shipment is sent while the feature flag is turned off.
			if !isUBFeatureOn && (*params.Body.ShipmentType == primev3messages.MTOShipmentTypeUNACCOMPANIEDBAGGAGE) {
				return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(
					"Unaccompanied baggage shipments can't be created unless the unaccompanied_baggage feature flag is enabled.", h.GetTraceIDFromRequest(params.HTTPRequest), nil)), nil
			}

			/** Feature Flag - Alaska - Determines if AK can be included/excluded **/
			isAlaskaEnabled := false
			akFeatureFlagName := "enable_alaska"
			flag, err = h.FeatureFlagFetcher().GetBooleanFlagForUser(context.TODO(), appCtx, akFeatureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", akFeatureFlagName), zap.Error(err))
			} else {
				isAlaskaEnabled = flag.Match
			}

			/** Feature Flag - Hawaii - Determines if HI can be included/excluded **/
			isHawaiiEnabled := false
			hiFeatureFlagName := "enable_hawaii"
			flag, err = h.FeatureFlagFetcher().GetBooleanFlagForUser(context.TODO(), appCtx, hiFeatureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", hiFeatureFlagName), zap.Error(err))
			} else {
				isHawaiiEnabled = flag.Match
			}

			// build states to exlude filter list
			statesToExclude := make([]string, 0)
			if !isAlaskaEnabled {
				statesToExclude = append(statesToExclude, "AK")
			}
			if !isHawaiiEnabled {
				statesToExclude = append(statesToExclude, "HI")
			}

			for _, mtoServiceItem := range params.Body.MtoServiceItems() {
				// restrict creation to a list
				if _, ok := CreateableServiceItemMap[mtoServiceItem.ModelType()]; !ok {
					// throw error if modelType() not on the list
					mapKeys := primeapi.GetMapKeys(primeapi.CreateableServiceItemMap)
					detailErr := fmt.Sprintf("MTOServiceItem modelType() not allowed: %s ", mtoServiceItem.ModelType())
					verrs := validate.NewErrors()
					verrs.Add("modelType", fmt.Sprintf("allowed modelType() %v", mapKeys))

					appCtx.Logger().Error("primeapiv3.CreateMTOShipmentHandler error", zap.Error(verrs))
					return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(
						detailErr, h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
				}
			}

			mtoShipment, verrs := payloads.MTOShipmentModelFromCreate(payload)
			if verrs != nil && verrs.HasAny() {
				appCtx.Logger().Error("Error validating mto shipment object: ", zap.Error(verrs))

				return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(
					"The MTO shipment object is invalid. "+verrs.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), nil)), verrs
			}

			mtoShipment.Status = models.MTOShipmentStatusSubmitted
			mtoServiceItemsList, verrs := payloads.MTOServiceItemModelListFromCreate(payload)

			if verrs != nil && verrs.HasAny() {
				appCtx.Logger().Error("Error validating mto service item list: ", zap.Error(verrs))

				return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(
					"The MTO service item list is invalid.", h.GetTraceIDFromRequest(params.HTTPRequest), nil)), verrs
			}

			mtoShipment.MTOServiceItems = mtoServiceItemsList

			moveTaskOrderID := uuid.FromStringOrNil(payload.MoveTaskOrderID.String())
			mtoAvailableToPrime, err := h.mtoAvailabilityChecker.MTOAvailableToPrime(appCtx, moveTaskOrderID)

			if mtoAvailableToPrime {
				// check each address prior to creating the shipment to ensure only valid addresses are being used to create the shipment
				var addresses []models.Address

				if mtoShipment.ShipmentType != models.MTOShipmentTypePPM {
					if mtoShipment.PickupAddress != nil {
						addresses = append(addresses, *mtoShipment.PickupAddress)
					}

					if mtoShipment.DestinationAddress != nil {
						addresses = append(addresses, *mtoShipment.DestinationAddress)
					}

					if mtoShipment.SecondaryPickupAddress != nil {
						addresses = append(addresses, *mtoShipment.SecondaryPickupAddress)
					}

					if mtoShipment.TertiaryPickupAddress != nil {
						addresses = append(addresses, *mtoShipment.TertiaryPickupAddress)
					}

					if mtoShipment.SecondaryDeliveryAddress != nil {
						addresses = append(addresses, *mtoShipment.SecondaryDeliveryAddress)
					}

					if mtoShipment.TertiaryDeliveryAddress != nil {
						addresses = append(addresses, *mtoShipment.TertiaryDeliveryAddress)
					}
				} else {
					if mtoShipment.PPMShipment.PickupAddress != nil {
						addresses = append(addresses, *mtoShipment.PPMShipment.PickupAddress)
					}

					if mtoShipment.PPMShipment.DestinationAddress != nil {
						addresses = append(addresses, *mtoShipment.PPMShipment.DestinationAddress)
					}

					if mtoShipment.PPMShipment.SecondaryPickupAddress != nil {
						addresses = append(addresses, *mtoShipment.PPMShipment.SecondaryPickupAddress)
					}

					if mtoShipment.PPMShipment.TertiaryPickupAddress != nil {
						addresses = append(addresses, *mtoShipment.PPMShipment.TertiaryPickupAddress)
					}

					if mtoShipment.PPMShipment.SecondaryDestinationAddress != nil {
						addresses = append(addresses, *mtoShipment.PPMShipment.SecondaryDestinationAddress)
					}

					if mtoShipment.PPMShipment.TertiaryDestinationAddress != nil {
						addresses = append(addresses, *mtoShipment.PPMShipment.TertiaryDestinationAddress)
					}
				}

				for _, address := range addresses {
					addressSearch := address.City + ", " + address.State + " " + address.PostalCode
					err := checkValidAddress(h.VLocation, appCtx, statesToExclude, addressSearch)

					if err != nil {
						appCtx.Logger().Error("primeapi.UpdateMTOShipmentHandler error", zap.Error(err))
						switch e := err.(type) {
						case apperror.UnprocessableEntityError:
							payload := payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), nil)
							return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payload), err
						default:
							errStr := e.Error() // we do this because InternalServerError wants a *string
							payload := payloads.InternalServerError(&errStr, h.GetTraceIDFromRequest(params.HTTPRequest))
							return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(payload), e
						}
					}
				}

				mtoShipment, err = h.ShipmentCreator.CreateShipment(appCtx, mtoShipment)
			} else if err == nil {
				appCtx.Logger().Error("primeapiv3.CreateMTOShipmentHandler error - MTO is not available to Prime")
				return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(payloads.ClientError(
					handlers.NotFoundMessage, fmt.Sprintf("id: %s not found for moveTaskOrder", moveTaskOrderID), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			// Could be the error from MTOAvailableToPrime or CreateMTOShipment:
			if err != nil {
				appCtx.Logger().Error("primeapiv3.CreateMTOShipmentHandler error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.EventError:
					return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(
						payloads.ClientError(handlers.InternalServerErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(
						payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("primeapiv3.CreateMTOShipmentHandler query error", zap.Error(e.Unwrap()))
					}
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}
			returnPayload := payloads.MTOShipment(mtoShipment)
			return mtoshipmentops.NewCreateMTOShipmentOK().WithPayload(returnPayload), nil
		})
}

func checkValidAddress(vLocation services.VLocation, appCtx appcontext.AppContext, statesToExclude []string, addressSearch string) error {
	locationList, err := vLocation.GetLocationsByZipCityState(appCtx, addressSearch, statesToExclude, true)

	if err != nil {
		serverError := apperror.NewInternalServerError("Error searching for address")
		return serverError
	} else if len(*locationList) == 0 {
		unprocessableErr := apperror.NewUnprocessableEntityError(
			fmt.Sprintf("primeapi.UpdateShipmentDestinationAddress: could not find the provided location: %s", addressSearch))
		return unprocessableErr
	}

	return nil
}

// UpdateMTOShipmentHandler is the handler to update MTO shipments
type UpdateMTOShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentUpdater
	planner route.Planner
	services.VLocation
}

// Handle handler that updates a mto shipment
func (h UpdateMTOShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			mtoShipment := payloads.MTOShipmentModelFromUpdate(params.Body, params.MtoShipmentID)

			dbShipment, err := mtoshipment.FindShipment(appCtx, mtoShipment.ID,
				"DestinationAddress",
				"SecondaryPickupAddress",
				"SecondaryDeliveryAddress",
				"TertiaryPickupAddress",
				"TertiaryDeliveryAddress",
				"StorageFacility",
				"PPMShipment",
				"MTOServiceItems.PODLocation",
				"MTOServiceItems.POELocation",
			)
			if err != nil {
				return mtoshipmentops.NewUpdateMTOShipmentNotFound().WithPayload(
					payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			var agents []models.MTOAgent
			err = appCtx.DB().Scope(utilities.ExcludeDeletedScope()).Where("mto_shipment_id = ?", mtoShipment.ID).All(&agents)
			if err != nil {
				return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
					payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}
			dbShipment.MTOAgents = agents

			// Validate further prime restrictions on model
			mtoShipment.ShipmentType = dbShipment.ShipmentType
			appCtx.Logger().Info("primeapi.UpdateMTOShipmentHandler info", zap.String("pointOfContact", params.Body.PointOfContact))

			/** Feature Flag - Alaska - Determines if AK can be included/excluded **/
			isAlaskaEnabled := false
			akFeatureFlagName := "enable_alaska"
			flag, err := h.FeatureFlagFetcher().GetBooleanFlagForUser(context.TODO(), appCtx, akFeatureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", akFeatureFlagName), zap.Error(err))
			} else {
				isAlaskaEnabled = flag.Match
			}

			/** Feature Flag - Hawaii - Determines if HI can be included/excluded **/
			isHawaiiEnabled := false
			hiFeatureFlagName := "enable_hawaii"
			flag, err = h.FeatureFlagFetcher().GetBooleanFlagForUser(context.TODO(), appCtx, hiFeatureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", hiFeatureFlagName), zap.Error(err))
			} else {
				isHawaiiEnabled = flag.Match
			}

			// build states to exlude filter list
			statesToExclude := make([]string, 0)
			if !isAlaskaEnabled {
				statesToExclude = append(statesToExclude, "AK")
			}
			if !isHawaiiEnabled {
				statesToExclude = append(statesToExclude, "HI")
			}

			// check each address prior to updating the shipment to ensure only valid addresses are being used
			var addresses []models.Address

			if mtoShipment.ShipmentType != models.MTOShipmentTypePPM {
				if mtoShipment.PickupAddress != nil {
					addresses = append(addresses, *mtoShipment.PickupAddress)
				}

				if mtoShipment.SecondaryPickupAddress != nil {
					addresses = append(addresses, *mtoShipment.SecondaryPickupAddress)
				}

				if mtoShipment.TertiaryPickupAddress != nil {
					addresses = append(addresses, *mtoShipment.TertiaryPickupAddress)
				}

				if mtoShipment.DestinationAddress != nil {
					addresses = append(addresses, *mtoShipment.DestinationAddress)
				}

				if mtoShipment.SecondaryDeliveryAddress != nil {
					addresses = append(addresses, *mtoShipment.SecondaryDeliveryAddress)
				}

				if mtoShipment.TertiaryDeliveryAddress != nil {
					addresses = append(addresses, *mtoShipment.TertiaryDeliveryAddress)
				}
			} else {
				if mtoShipment.PPMShipment.PickupAddress != nil {
					addresses = append(addresses, *mtoShipment.PPMShipment.PickupAddress)
				}

				if mtoShipment.PPMShipment.SecondaryPickupAddress != nil {
					addresses = append(addresses, *mtoShipment.PPMShipment.SecondaryPickupAddress)
				}

				if mtoShipment.PPMShipment.TertiaryPickupAddress != nil {
					addresses = append(addresses, *mtoShipment.PPMShipment.TertiaryPickupAddress)
				}

				if mtoShipment.PPMShipment.DestinationAddress != nil {
					addresses = append(addresses, *mtoShipment.PPMShipment.DestinationAddress)
				}

				if mtoShipment.PPMShipment.SecondaryDestinationAddress != nil {
					addresses = append(addresses, *mtoShipment.PPMShipment.SecondaryDestinationAddress)
				}

				if mtoShipment.PPMShipment.TertiaryDestinationAddress != nil {
					addresses = append(addresses, *mtoShipment.PPMShipment.TertiaryDestinationAddress)
				}
			}

			for _, address := range addresses {
				addressSearch := address.City + ", " + address.State + " " + address.PostalCode
				err := checkValidAddress(h.VLocation, appCtx, statesToExclude, addressSearch)

				if err != nil {
					appCtx.Logger().Error("primeapi.UpdateMTOShipmentHandler error", zap.Error(err))
					switch e := err.(type) {
					case apperror.UnprocessableEntityError:
						payload := payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), nil)
						return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(payload), e
					default:
						errStr := e.Error() // we do this because InternalServerError wants a *string
						payload := payloads.InternalServerError(&errStr, h.GetTraceIDFromRequest(params.HTTPRequest))
						return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(payload), e
					}
				}
			}

			mtoShipment, err = h.ShipmentUpdater.UpdateShipment(appCtx, mtoShipment, params.IfMatch, "prime-v3")
			if err != nil {
				appCtx.Logger().Error("primeapi.UpdateMTOShipmentHnadler error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.EventError:
					return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(
						payloads.ClientError(handlers.InternalServerErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					payload := payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(
						payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}
			mtoShipmentPayload := payloads.MTOShipment(mtoShipment)
			return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(mtoShipmentPayload), nil
		})
}
