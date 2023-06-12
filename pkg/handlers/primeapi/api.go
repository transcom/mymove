package primeapi

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/primeapi"
	"github.com/transcom/mymove/pkg/gen/primeapi/primeoperations"
	"github.com/transcom/mymove/pkg/handlers"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoagent "github.com/transcom/mymove/pkg/services/mto_agent"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/services/reweigh"
	sitaddressupdate "github.com/transcom/mymove/pkg/services/sit_address_update"
	sitextension "github.com/transcom/mymove/pkg/services/sit_extension"
	"github.com/transcom/mymove/pkg/services/upload"
)

// NewPrimeAPI returns the Prime API
func NewPrimeAPI(handlerConfig handlers.HandlerConfig) *primeoperations.MymoveAPI {
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)

	primeSpec, err := loads.Analyzed(primeapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	primeAPI := primeoperations.NewMymoveAPI(primeSpec)
	queryBuilder := query.NewQueryBuilder()
	moveRouter := move.NewMoveRouter()
	moveWeights := move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())
	uploadCreator := upload.NewUploadCreator(handlerConfig.FileStorer())

	paymentRequestRecalculator := paymentrequest.NewPaymentRequestRecalculator(
		paymentrequest.NewPaymentRequestCreator(
			handlerConfig.HHGPlanner(),
			ghcrateengine.NewServiceItemPricer(),
		),
		paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
	)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(paymentRequestRecalculator)

	primeAPI.ServeError = handlers.ServeCustomError

	primeAPI.MoveTaskOrderListMovesHandler = ListMovesHandler{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}

	primeAPI.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandler{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderFetcher(),
	}

	primeAPI.MoveTaskOrderCreateExcessWeightRecordHandler = CreateExcessWeightRecordHandler{
		handlerConfig,
		move.NewPrimeMoveExcessWeightUploader(uploadCreator),
	}

	primeAPI.MtoServiceItemCreateMTOServiceItemHandler = CreateMTOServiceItemHandler{
		handlerConfig,
		mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter),
		movetaskorder.NewMoveTaskOrderChecker(),
	}

	primeAPI.MtoServiceItemUpdateMTOServiceItemHandler = UpdateMTOServiceItemHandler{
		handlerConfig,
		mtoserviceitem.NewMTOServiceItemUpdater(builder, moveRouter),
	}

	primeAPI.MtoServiceItemCreateServiceRequestDocumentUploadHandler = CreateServiceRequestDocumentUploadHandler{
		handlerConfig,
		mtoserviceitem.NewServiceRequestDocumentUploadCreator(handlerConfig.FileStorer()),
	}

	mtoShipmentUpdater := mtoshipment.NewPrimeMTOShipmentUpdater(
		builder,
		fetcher,
		handlerConfig.HHGPlanner(),
		moveRouter,
		moveWeights,
		handlerConfig.NotificationSender(),
		paymentRequestShipmentRecalculator,
	)

	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})
	addressCreator := address.NewAddressCreator()
	addressUpdater := address.NewAddressUpdater()
	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(ppmEstimator, addressCreator, addressUpdater)
	shipmentUpdater := shipment.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater)

	primeAPI.MtoShipmentUpdateMTOShipmentHandler = UpdateMTOShipmentHandler{
		handlerConfig,
		shipmentUpdater,
	}

	primeAPI.MtoShipmentDeleteMTOShipmentHandler = DeleteMTOShipmentHandler{
		handlerConfig,
		mtoshipment.NewPrimeShipmentDeleter(),
	}

	primeAPI.PaymentRequestCreatePaymentRequestHandler = CreatePaymentRequestHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestCreator(
			handlerConfig.HHGPlanner(),
			ghcrateengine.NewServiceItemPricer(),
		),
	}

	primeAPI.PaymentRequestCreateUploadHandler = CreateUploadHandler{
		handlerConfig,
		paymentrequest.NewPaymentRequestUploadCreator(handlerConfig.FileStorer()),
	}

	primeAPI.MoveTaskOrderUpdateMTOPostCounselingInformationHandler = UpdateMTOPostCounselingInformationHandler{
		handlerConfig,
		fetch.NewFetcher(queryBuilder),
		movetaskorder.NewMoveTaskOrderUpdater(
			queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
			moveRouter,
		),
		movetaskorder.NewMoveTaskOrderChecker(),
	}

	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
	ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(ppmEstimator)
	shipmentRouter := mtoshipment.NewShipmentRouter()
	shipmentCreator := shipment.NewShipmentCreator(mtoShipmentCreator, ppmShipmentCreator, shipmentRouter)

	primeAPI.MtoShipmentCreateMTOShipmentHandler = CreateMTOShipmentHandler{
		handlerConfig,
		shipmentCreator,
		movetaskorder.NewMoveTaskOrderChecker(),
	}

	primeAPI.MtoShipmentUpdateMTOShipmentAddressHandler = UpdateMTOShipmentAddressHandler{
		handlerConfig,
		mtoshipment.NewMTOShipmentAddressUpdater(),
	}

	primeAPI.MtoShipmentCreateMTOAgentHandler = CreateMTOAgentHandler{
		handlerConfig,
		mtoagent.NewMTOAgentCreator(movetaskorder.NewMoveTaskOrderChecker()),
	}

	primeAPI.MtoShipmentUpdateMTOAgentHandler = UpdateMTOAgentHandler{
		handlerConfig,
		mtoagent.NewMTOAgentUpdater(movetaskorder.NewMoveTaskOrderChecker()),
	}

	primeAPI.MtoShipmentUpdateMTOShipmentStatusHandler = UpdateMTOShipmentStatusHandler{
		handlerConfig,
		mtoshipment.NewPrimeMTOShipmentUpdater(builder, fetcher, handlerConfig.HHGPlanner(), moveRouter, moveWeights, handlerConfig.NotificationSender(), paymentRequestShipmentRecalculator),
		mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter), handlerConfig.HHGPlanner()),
	}

	primeAPI.MtoShipmentUpdateReweighHandler = UpdateReweighHandler{
		handlerConfig,
		reweigh.NewReweighUpdater(movetaskorder.NewMoveTaskOrderChecker(), paymentRequestShipmentRecalculator),
	}

	primeAPI.MtoShipmentCreateSITExtensionHandler = CreateSITExtensionHandler{
		handlerConfig,
		sitextension.NewSitExtensionCreator(moveRouter),
	}

	primeAPI.SitAddressUpdateCreateSITAddressUpdateRequestHandler = CreateSITAddressUpdateRequestHandler{
		handlerConfig,
		sitaddressupdate.NewSITAddressUpdateRequestCreator(handlerConfig.HHGPlanner(), addressCreator, moveRouter),
	}
	return primeAPI
}
