package ghcv2api

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/ghcapi"
	ghcv2ops "github.com/transcom/mymove/pkg/gen/ghcv2api/ghcv2operations"
	"github.com/transcom/mymove/pkg/handlers"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/move"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
	sitstatus "github.com/transcom/mymove/pkg/services/sit_status"
)

// NewGhcV2APIHandler returns a handler for the GHC V2 API
func NewGhcV2APIHandler(handlerConfig handlers.HandlerConfig) *ghcv2ops.MymoveAPI {
	ghcV2Spec, err := loads.Analyzed(ghcapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	ghcV2API := ghcv2ops.NewMymoveAPI(ghcV2Spec)
	queryBuilder := query.NewQueryBuilder()
	moveRouter := move.NewMoveRouter()
	addressCreator := address.NewAddressCreator()
	paymentRequestRecalculator := paymentrequest.NewPaymentRequestRecalculator(
		paymentrequest.NewPaymentRequestCreator(
			handlerConfig.HHGPlanner(),
			ghcrateengine.NewServiceItemPricer(),
		),
		paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
	)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(paymentRequestRecalculator)

	shipmentSITStatus := sitstatus.NewShipmentSITStatus()

	ghcV2API.ServeError = handlers.ServeCustomError

	mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(
		queryBuilder,
		fetch.NewFetcher(queryBuilder),
		handlerConfig.HHGPlanner(),
		moveRouter,
		move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester()),
		handlerConfig.NotificationSender(),
		paymentRequestShipmentRecalculator,
	)

	addressUpdater := address.NewAddressUpdater()
	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})
	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(ppmEstimator, addressCreator, addressUpdater)
	shipmentUpdater := shipment.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)

	ghcV2API.MtoShipmentUpdateMTOShipmentHandler = UpdateShipmentHandler{
		handlerConfig,
		shipmentUpdater,
		shipmentSITStatus,
	}

	return ghcV2API
}
