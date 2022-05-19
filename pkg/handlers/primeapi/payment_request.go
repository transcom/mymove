package primeapi

import (
	"fmt"
	"sort"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_request"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// CreatePaymentRequestHandler is the handler for creating payment requests
type CreatePaymentRequestHandler struct {
	handlers.HandlerConfig
	services.PaymentRequestCreator
}

// Handle creates the payment request
func (h CreatePaymentRequestHandler) Handle(params paymentrequestop.CreatePaymentRequestParams) middleware.Responder {
	// TODO: authorization to create payment request

	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			payload := params.Body
			if payload == nil {
				err := apperror.NewBadDataError("Invalid payment request: params Body is nil")
				errPayload := payloads.ClientError(handlers.SQLErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))
				appCtx.Logger().Error(err.Error(), zap.Any("payload", errPayload))
				return paymentrequestop.NewCreatePaymentRequestBadRequest().WithPayload(errPayload), err
			}

			appCtx.Logger().Info("primeapi.CreatePaymentRequestHandler info", zap.String("pointOfContact", params.Body.PointOfContact))

			moveTaskOrderIDString := payload.MoveTaskOrderID.String()
			mtoID, err := uuid.FromString(moveTaskOrderIDString)
			if err != nil {
				appCtx.Logger().Error("Invalid payment request: params MoveTaskOrderID cannot be converted to a UUID",
					zap.String("MoveTaskOrderID", moveTaskOrderIDString), zap.Error(err))
				// create a custom verrs for returning a 422
				verrs :=
					&validate.Errors{Errors: map[string][]string{
						"move_id": {"id cannot be converted to UUID"},
					},
					}
				errPayload := payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), verrs)
				return paymentrequestop.NewCreatePaymentRequestUnprocessableEntity().WithPayload(errPayload), err
			}

			isFinal := false
			if payload.IsFinal != nil {
				isFinal = *payload.IsFinal
			}

			paymentRequest := models.PaymentRequest{
				IsFinal:         isFinal,
				MoveTaskOrderID: mtoID,
			}

			// Build up the paymentRequest.PaymentServiceItems using the incoming payload to offload Swagger data coming
			// in from the API. These paymentRequest.PaymentServiceItems will be used as a temp holder to process the incoming API data
			var verrs *validate.Errors
			paymentRequest.PaymentServiceItems, verrs, err = h.buildPaymentServiceItems(appCtx, payload)

			if err != nil || verrs.HasAny() {

				appCtx.Logger().Error("could not build service items", zap.Error(err))
				// TODO: do not bail out before creating the payment request, we need the failed record
				//       we should create the failed record and store it as failed with a rejection
				errPayload := payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), verrs)
				return paymentrequestop.NewCreatePaymentRequestUnprocessableEntity().WithPayload(errPayload), err
			}

			createdPaymentRequest, err := h.PaymentRequestCreator.CreatePaymentRequestCheck(appCtx, &paymentRequest)
			if err != nil {
				appCtx.Logger().Error("Error creating payment request", zap.Error(err))
				switch e := err.(type) {
				case apperror.InvalidCreateInputError:
					verrs := e.ValidationErrors
					detail := err.Error()
					payload := payloads.ValidationError(detail, h.GetTraceIDFromRequest(params.HTTPRequest), verrs)

					appCtx.Logger().Error("Payment Request",
						zap.Any("payload", payload))
					return paymentrequestop.NewCreatePaymentRequestUnprocessableEntity().WithPayload(payload), err

				case apperror.NotFoundError:
					payload := payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))

					appCtx.Logger().Error("Payment Request",
						zap.Any("payload", payload))
					return paymentrequestop.NewCreatePaymentRequestNotFound().WithPayload(payload), err
				case apperror.ConflictError:
					payload := payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))

					appCtx.Logger().Error("Payment Request",
						zap.Any("payload", payload))
					return paymentrequestop.NewCreatePaymentRequestConflict().WithPayload(payload), err
				case apperror.InvalidInputError:
					payload := payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), &validate.Errors{})

					appCtx.Logger().Error("Payment Request",
						zap.Any("payload", payload))
					return paymentrequestop.NewCreatePaymentRequestUnprocessableEntity().WithPayload(payload), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("primeapi.CreatePaymentRequestHandler query error", zap.Error(e.Unwrap()))
					}
					return paymentrequestop.NewCreatePaymentRequestInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err

				case *apperror.BadDataError:
					payload := payloads.ClientError(handlers.BadRequestErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))

					appCtx.Logger().Error("Payment Request",
						zap.Any("payload", payload))
					return paymentrequestop.NewCreatePaymentRequestBadRequest().WithPayload(payload), err
				default:
					appCtx.Logger().Error("Payment Request",
						zap.Any("payload", payload))
					return paymentrequestop.NewCreatePaymentRequestInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			returnPayload := payloads.PaymentRequest(createdPaymentRequest)
			appCtx.Logger().Info("Successful payment request creation for mto ID", zap.String("moveID", moveTaskOrderIDString))
			return paymentrequestop.NewCreatePaymentRequestCreated().WithPayload(returnPayload), nil
		})
}

func (h CreatePaymentRequestHandler) buildPaymentServiceItems(appCtx appcontext.AppContext, payload *primemessages.CreatePaymentRequest) (models.PaymentServiceItems, *validate.Errors, error) {

	var paymentServiceItems models.PaymentServiceItems
	verrs := validate.NewErrors()
	for _, payloadServiceItem := range payload.ServiceItems {
		mtoServiceItemID, err := uuid.FromString(payloadServiceItem.ID.String())
		if err != nil {
			// create a custom verrs for returning a 422
			verrs = &validate.Errors{Errors: map[string][]string{
				"payment_service_item_id": {"id cannot be converted to UUID"},
			},
			}
			return nil, verrs, fmt.Errorf("could not convert service item ID [%v] to UUID: %w", payloadServiceItem.ID, err)
		}

		// Find the ReService model that maps to to the MTOServiceItem
		var mtoServiceItem models.MTOServiceItem
		err = appCtx.DB().Eager("ReService").Find(&mtoServiceItem, mtoServiceItemID)
		if err != nil {
			return nil, verrs, fmt.Errorf("could not find RE (rate engine) service item for MTO Service Item with UUID %s with error: %w", mtoServiceItemID, err)
		}

		paymentServiceItem := models.PaymentServiceItem{
			// The rest of the model will be filled in when the payment request is created
			MTOServiceItemID: mtoServiceItemID,
			MTOServiceItem:   mtoServiceItem,
		}

		paymentServiceItem.PaymentServiceItemParams, err = h.buildPaymentServiceItemParams(mtoServiceItem.ReService.Code, payloadServiceItem)
		if err != nil {
			return nil, verrs, err
		}

		paymentServiceItems = append(paymentServiceItems, paymentServiceItem)
	}

	sort.SliceStable(paymentServiceItems, func(i, j int) bool {
		return paymentServiceItems[i].MTOServiceItem.ReService.Priority < paymentServiceItems[j].MTOServiceItem.ReService.Priority
	})

	return paymentServiceItems, verrs, nil
}

func (h CreatePaymentRequestHandler) buildPaymentServiceItemParams(reServiceCode models.ReServiceCode, payloadMTOServiceItem *primemessages.ServiceItem) (models.PaymentServiceItemParams, error) {
	var paymentServiceItemParams models.PaymentServiceItemParams

	for _, payloadServiceItemParam := range payloadMTOServiceItem.Params {
		if !AllowedParamKeysPaymentRequest.Contains(reServiceCode, payloadServiceItemParam.Key) {
			return models.PaymentServiceItemParams{}, fmt.Errorf("the parameter %s is either invalid or cannot be passed while creating a payment request for a %s service item", payloadServiceItemParam.Key, reServiceCode)
		}
		paymentServiceItemParam := models.PaymentServiceItemParam{
			// ID and PaymentServiceItemID to be filled in when payment request is created
			IncomingKey: payloadServiceItemParam.Key,
			Value:       payloadServiceItemParam.Value,
		}

		paymentServiceItemParams = append(paymentServiceItemParams, paymentServiceItemParam)
	}

	return paymentServiceItemParams, nil
}
