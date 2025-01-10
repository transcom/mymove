package primeapiv3

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/primev3api"
	"github.com/transcom/mymove/pkg/gen/primev3api/primev3operations"
	"github.com/transcom/mymove/pkg/handlers"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/services/address"
	boatshipment "github.com/transcom/mymove/pkg/services/boat_shipment"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	mobilehomeshipment "github.com/transcom/mymove/pkg/services/mobile_home_shipment"
	"github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
	signedcertification "github.com/transcom/mymove/pkg/services/signed_certification"
)

// NewPrimeAPI returns the Prime API
func NewPrimeAPI(handlerConfig handlers.HandlerConfig) *primev3operations.MymoveAPI {
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	queryBuilder := query.NewQueryBuilder()
	moveRouter := move.NewMoveRouter()

	primeSpec, err := loads.Analyzed(primev3api.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	primeAPIV3 := primev3operations.NewMymoveAPI(primeSpec)

	primeAPIV3.ServeError = handlers.ServeCustomError

	addressCreator := address.NewAddressCreator()
	addressUpdater := address.NewAddressUpdater()

	primeAPIV3.MoveTaskOrderGetMoveTaskOrderHandler = GetMoveTaskOrderHandler{
		handlerConfig,
		movetaskorder.NewMoveTaskOrderFetcher(),
		mtoshipment.NewMTOShipmentRateAreaFetcher(),
	}

	signedCertificationCreator := signedcertification.NewSignedCertificationCreator()
	signedCertificationUpdater := signedcertification.NewSignedCertificationUpdater()
	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})

	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		queryBuilder,
		mtoserviceitem.NewMTOServiceItemCreator(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
		moveRouter, signedCertificationCreator, signedCertificationUpdater, ppmEstimator,
	)

	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreatorV2(builder, fetcher, moveRouter, addressCreator)
	ppmShipmentCreator := ppmshipment.NewPPMShipmentCreator(ppmEstimator, addressCreator)
	boatShipmentCreator := boatshipment.NewBoatShipmentCreator()
	mobileHomeShipmentCreator := mobilehomeshipment.NewMobileHomeShipmentCreator()
	shipmentRouter := mtoshipment.NewShipmentRouter()

	shipmentCreator := shipment.NewShipmentCreator(mtoShipmentCreator, ppmShipmentCreator, boatShipmentCreator, mobileHomeShipmentCreator, shipmentRouter, moveTaskOrderUpdater)

	primeAPIV3.MtoShipmentCreateMTOShipmentHandler = CreateMTOShipmentHandler{
		handlerConfig,
		shipmentCreator,
		movetaskorder.NewMoveTaskOrderChecker(),
	}
	paymentRequestRecalculator := paymentrequest.NewPaymentRequestRecalculator(
		paymentrequest.NewPaymentRequestCreator(
			handlerConfig.HHGPlanner(),
			ghcrateengine.NewServiceItemPricer(),
		),
		paymentrequest.NewPaymentRequestStatusUpdater(queryBuilder),
	)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(paymentRequestRecalculator)
	moveWeights := move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())
	mtoShipmentUpdater := mtoshipment.NewPrimeMTOShipmentUpdater(
		builder,
		fetcher,
		handlerConfig.HHGPlanner(),
		moveRouter,
		moveWeights,
		handlerConfig.NotificationSender(),
		paymentRequestShipmentRecalculator,
		addressUpdater,
		addressCreator,
	)

	mtoServiceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(
		handlerConfig.HHGPlanner(),
		queryBuilder,
		moveRouter,
		ghcrateengine.NewDomesticUnpackPricer(),
		ghcrateengine.NewDomesticPackPricer(),
		ghcrateengine.NewDomesticLinehaulPricer(),
		ghcrateengine.NewDomesticShorthaulPricer(),
		ghcrateengine.NewDomesticOriginPricer(),
		ghcrateengine.NewDomesticDestinationPricer(),
		ghcrateengine.NewFuelSurchargePricer())

	ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(ppmEstimator, addressCreator, addressUpdater)
	boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
	mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()
	shipmentUpdater := shipment.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater, mtoServiceItemCreator)
	primeAPIV3.MtoShipmentUpdateMTOShipmentHandler = UpdateMTOShipmentHandler{
		handlerConfig,
		shipmentUpdater,
		handlerConfig.DTODPlanner(),
	}

	return primeAPIV3
}
