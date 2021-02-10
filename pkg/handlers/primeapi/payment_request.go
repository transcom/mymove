package primeapi

import (
	"fmt"
	"sort"

	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"

	"github.com/gobuffalo/validate/v3"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_request"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// CreatePaymentRequestHandler is the handler for creating payment requests
type CreatePaymentRequestHandler struct {
	handlers.HandlerContext
	services.PaymentRequestCreator
}

// Handle creates the payment request
func (h CreatePaymentRequestHandler) Handle(params paymentrequestop.CreatePaymentRequestParams) middleware.Responder {
	// TODO: authorization to create payment request

	logger := h.LoggerFromRequest(params.HTTPRequest)

	payload := params.Body
	if payload == nil {
		errPayload := payloads.ClientError(handlers.SQLErrMessage, "Invalid payment request: params Body is nil", h.GetTraceID())
		logger.Error("Invalid payment request: params Body is nil", zap.Any("payload", errPayload))
		return paymentrequestop.NewCreatePaymentRequestBadRequest().WithPayload(errPayload)
	}

	logger.Info("primeapi.CreatePaymentRequestHandler info", zap.String("pointOfContact", params.Body.PointOfContact))

	moveTaskOrderIDString := payload.MoveTaskOrderID.String()
	mtoID, err := uuid.FromString(moveTaskOrderIDString)
	if err != nil {
		logger.Error("Invalid payment request: params MoveTaskOrderID cannot be converted to a UUID",
			zap.String("MoveTaskOrderID", moveTaskOrderIDString), zap.Error(err))
		// create a custom verrs for returning a 422
		verrs :=
			&validate.Errors{Errors: map[string][]string{
				"move_id": {"id cannot be converted to UUID"},
			},
			}
		errPayload := payloads.ValidationError(err.Error(), h.GetTraceID(), verrs)
		return paymentrequestop.NewCreatePaymentRequestUnprocessableEntity().WithPayload(errPayload)
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
	verrs := validate.NewErrors()
	paymentRequest.PaymentServiceItems, verrs, err = h.buildPaymentServiceItems(payload)

	if err != nil || verrs.HasAny() {

		logger.Error("could not build service items", zap.Error(err))
		// TODO: do not bail out before creating the payment request, we need the failed record
		//       we should create the failed record and store it as failed with a rejection
		errPayload := payloads.ValidationError(err.Error(), h.GetTraceID(), verrs)
		return paymentrequestop.NewCreatePaymentRequestUnprocessableEntity().WithPayload(errPayload)
	}

	createdPaymentRequest, err := h.PaymentRequestCreator.CreatePaymentRequest(&paymentRequest)
	if err != nil {
		logger.Error("Error creating payment request", zap.Error(err))
		switch e := err.(type) {
		case services.InvalidCreateInputError:
			verrs := e.ValidationErrors
			detail := err.Error()
			payload := payloads.ValidationError(detail, h.GetTraceID(), verrs)

			logger.Error("Payment Request",
				zap.Any("payload", payload))
			return paymentrequestop.NewCreatePaymentRequestUnprocessableEntity().WithPayload(payload)

		case services.NotFoundError:
			payload := payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceID())

			logger.Error("Payment Request",
				zap.Any("payload", payload))
			return paymentrequestop.NewCreatePaymentRequestNotFound().WithPayload(payload)
		case services.ConflictError:
			payload := payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceID())

			logger.Error("Payment Request",
				zap.Any("payload", payload))
			return paymentrequestop.NewCreatePaymentRequestConflict().WithPayload(payload)
		case services.InvalidInputError:
			payload := payloads.ValidationError(err.Error(), h.GetTraceID(), &validate.Errors{})

			logger.Error("Payment Request",
				zap.Any("payload", payload))
			return paymentrequestop.NewCreatePaymentRequestUnprocessableEntity().WithPayload(payload)
		case services.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				logger.Error("primeapi.CreatePaymentRequestHandler query error", zap.Error(e.Unwrap()))
			}
			return paymentrequestop.NewCreatePaymentRequestInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))

		case *services.BadDataError:
			payload := payloads.ClientError(handlers.BadRequestErrMessage, err.Error(), h.GetTraceID())

			logger.Error("Payment Request",
				zap.Any("payload", payload))
			return paymentrequestop.NewCreatePaymentRequestBadRequest().WithPayload(payload)
		default:
			logger.Error("Payment Request",
				zap.Any("payload", payload))
			return paymentrequestop.NewCreatePaymentRequestInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceID()))
		}
	}

	returnPayload := payloads.PaymentRequest(createdPaymentRequest)
	logger.Info("Successful payment request creation for mto ID", zap.String("moveID", moveTaskOrderIDString))
	return paymentrequestop.NewCreatePaymentRequestCreated().WithPayload(returnPayload)
}

func (h CreatePaymentRequestHandler) buildPaymentServiceItems(payload *primemessages.CreatePaymentRequest) (models.PaymentServiceItems, *validate.Errors, error) {

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
		err = h.DB().Eager("ReService").Find(&mtoServiceItem, mtoServiceItemID)
		if err != nil {
			return nil, verrs, fmt.Errorf("could not find RE (rate engine) service item for MTO Service Item with UUID %s with error: %w", mtoServiceItemID, err)
		}

		paymentServiceItem := models.PaymentServiceItem{
			// The rest of the model will be filled in when the payment request is created
			MTOServiceItemID: mtoServiceItemID,
			MTOServiceItem:   mtoServiceItem,
		}

		paymentServiceItem.PaymentServiceItemParams, err = h.buildPaymentServiceItemParams(payloadServiceItem, mtoServiceItem.ReService)
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

func (h CreatePaymentRequestHandler) buildPaymentServiceItemParams(payloadMTOServiceItem *primemessages.ServiceItem, reService models.ReService) (models.PaymentServiceItemParams, error) {
	/************
	  ServiceItem.params is set to readOnly = true currently in prime.yaml. Therefore we are only checking if
	  there were params sent. If there were params in via the create payment request then we will error out.

	  Currently not expecting the prime to provide any params. This might change as we continue adding service items
	  for billing and then we'll have to adjust which service items allow incoming params at that time.
	***********/

	if len(payloadMTOServiceItem.Params) > 0 {
		// if not in this function it can also be done up top
		return models.PaymentServiceItemParams{}, fmt.Errorf("updating service item params not allowed for service item [%s] with MTO Service UUID: %s", reService.Name, payloadMTOServiceItem.ID)

	}

	return models.PaymentServiceItemParams{}, nil
}
