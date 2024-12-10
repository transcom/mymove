package primeapi

import (
	"log"

	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/primeapi"
	"github.com/transcom/mymove/pkg/gen/primeapi/primeoperations"
	"github.com/transcom/mymove/pkg/handlers"
	paperwork "github.com/transcom/mymove/pkg/paperwork"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoagent "github.com/transcom/mymove/pkg/services/mto_agent"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	order "github.com/transcom/mymove/pkg/services/order"
	paperwork_service "github.com/transcom/mymove/pkg/services/paperwork"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/services/reweigh"
	shipmentaddressupdate "github.com/transcom/mymove/pkg/services/shipment_address_update"
	signedcertification "github.com/transcom/mymove/pkg/services/signed_certification"
	sitextension "github.com/transcom/mymove/pkg/services/sit_extension"
	"github.com/transcom/mymove/pkg/services/upload"
	"github.com/transcom/mymove/pkg/uploader"
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
	addressCreator := address.NewAddressCreator()
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
	moveWeights := move.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())
	uploadCreator := upload.NewUploadCreator(handlerConfig.FileStorer())
	ppmEstimator := ppmshipment.NewEstimatePPM(handlerConfig.DTODPlanner(), &paymentrequesthelper.RequestPaymentHelper{})
	serviceItemUpdater := mtoserviceitem.NewMTOServiceItemUpdater(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, shipmentFetcher, addressCreator, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer())

	userUploader, err := uploader.NewUserUploader(handlerConfig.FileStorer(), uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		log.Fatalln(err)
	}

	pdfGenerator, err := paperwork.NewGenerator(userUploader.Uploader())
	if err != nil {
		log.Fatalln(err)
	}
	primeDownloadMoveUploadPDFGenerator, err := paperwork_service.NewMoveUserUploadToPDFDownloader(pdfGenerator)
	if err != nil {
		log.Fatalln(err)
	}

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
		mtoserviceitem.NewMTOServiceItemCreator(handlerConfig.HHGPlanner(), builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
		movetaskorder.NewMoveTaskOrderChecker(),
	}

	primeAPI.MtoServiceItemUpdateMTOServiceItemHandler = UpdateMTOServiceItemHandler{
		handlerConfig,
		serviceItemUpdater,
	}

	primeAPI.MtoServiceItemCreateServiceRequestDocumentUploadHandler = CreateServiceRequestDocumentUploadHandler{
		handlerConfig,
		mtoserviceitem.NewServiceRequestDocumentUploadCreator(handlerConfig.FileStorer()),
	}

	primeAPI.MtoShipmentUpdateShipmentDestinationAddressHandler = UpdateShipmentDestinationAddressHandler{
		handlerConfig,
		shipmentaddressupdate.NewShipmentAddressUpdateRequester(handlerConfig.HHGPlanner(), addressCreator, moveRouter),
	}

	addressUpdater := address.NewAddressUpdater()

	signedCertificationCreator := signedcertification.NewSignedCertificationCreator()
	signedCertificationUpdater := signedcertification.NewSignedCertificationUpdater()
	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		queryBuilder,
		mtoserviceitem.NewMTOServiceItemCreator(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
		moveRouter, signedCertificationCreator, signedCertificationUpdater, ppmEstimator,
	)

	primeAPI.MtoShipmentDeleteMTOShipmentHandler = DeleteMTOShipmentHandler{
		handlerConfig,
		mtoshipment.NewPrimeShipmentDeleter(moveTaskOrderUpdater),
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
			mtoserviceitem.NewMTOServiceItemCreator(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			moveRouter, signedCertificationCreator, signedCertificationUpdater, ppmEstimator,
		),
		movetaskorder.NewMoveTaskOrderChecker(),
	}

	primeAPI.MtoShipmentUpdateMTOShipmentAddressHandler = UpdateMTOShipmentAddressHandler{
		handlerConfig,
		mtoshipment.NewMTOShipmentAddressUpdater(handlerConfig.HHGPlanner(), addressCreator, addressUpdater),
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
		mtoshipment.NewPrimeMTOShipmentUpdater(builder, fetcher, handlerConfig.HHGPlanner(), moveRouter, moveWeights, handlerConfig.NotificationSender(), paymentRequestShipmentRecalculator, addressUpdater, addressCreator),
		mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder,
			mtoserviceitem.NewMTOServiceItemCreator(handlerConfig.HHGPlanner(), queryBuilder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()), handlerConfig.HHGPlanner()),
	}

	primeAPI.MtoShipmentUpdateReweighHandler = UpdateReweighHandler{
		handlerConfig,
		reweigh.NewReweighUpdater(movetaskorder.NewMoveTaskOrderChecker(), paymentRequestShipmentRecalculator),
	}

	primeAPI.MtoShipmentCreateSITExtensionHandler = CreateSITExtensionHandler{
		handlerConfig,
		sitextension.NewSitExtensionCreator(moveRouter),
	}

	primeAPI.MoveTaskOrderDownloadMoveOrderHandler = DownloadMoveOrderHandler{
		handlerConfig,
		move.NewMoveSearcher(),
		order.NewOrderFetcher(),
		primeDownloadMoveUploadPDFGenerator,
	}

	return primeAPI
}
