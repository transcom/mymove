package publicapi

import (

	// "os"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	shipmentop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/shipments"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/paperwork"
	uploaderpkg "github.com/transcom/mymove/pkg/uploader"
	"go.uber.org/zap"
)

func payloadForShipmentModel(s models.Shipment) *apimessages.Shipment {
	shipmentpayload := &apimessages.Shipment{
		ID:                                  *handlers.FmtUUID(s.ID),
		TrafficDistributionList:             payloadForTrafficDistributionListModel(s.TrafficDistributionList),
		ServiceMember:                       payloadForServiceMemberModel(s.ServiceMember),
		PickupDate:                          *handlers.FmtDateTimePtr(s.PickupDate),
		DeliveryDate:                        *handlers.FmtDateTimePtr(s.DeliveryDate),
		CreatedAt:                           strfmt.DateTime(s.CreatedAt),
		UpdatedAt:                           strfmt.DateTime(s.UpdatedAt),
		SourceGbloc:                         apimessages.GBLOC(*s.SourceGBLOC),
		DestinationGbloc:                    apimessages.GBLOC(*s.DestinationGBLOC),
		Market:                              apimessages.ShipmentMarket(*s.Market),
		BookDate:                            *handlers.FmtDatePtr(s.BookDate),
		RequestedPickupDate:                 *handlers.FmtDateTimePtr(s.RequestedPickupDate),
		Move:                                payloadForMoveModel(&s.Move),
		Status:                              apimessages.ShipmentStatus(s.Status),
		EstimatedPackDays:                   handlers.FmtInt64(*s.EstimatedPackDays),
		EstimatedTransitDays:                handlers.FmtInt64(*s.EstimatedTransitDays),
		PickupAddress:                       payloadForAddressModel(s.PickupAddress),
		HasSecondaryPickupAddress:           *handlers.FmtBool(s.HasSecondaryPickupAddress),
		SecondaryPickupAddress:              payloadForAddressModel(s.SecondaryPickupAddress),
		HasDeliveryAddress:                  *handlers.FmtBool(s.HasDeliveryAddress),
		DeliveryAddress:                     payloadForAddressModel(s.DeliveryAddress),
		HasPartialSitDeliveryAddress:        *handlers.FmtBool(s.HasPartialSITDeliveryAddress),
		PartialSitDeliveryAddress:           payloadForAddressModel(s.PartialSITDeliveryAddress),
		WeightEstimate:                      handlers.FmtInt64(s.WeightEstimate.Int64()),
		ProgearWeightEstimate:               handlers.FmtInt64(s.ProgearWeightEstimate.Int64()),
		SpouseProgearWeightEstimate:         handlers.FmtInt64(s.SpouseProgearWeightEstimate.Int64()),
		ActualWeight:                        handlers.FmtPoundPtr(s.ActualWeight),
		PmSurveyPlannedPackDate:             handlers.FmtDatePtr(s.PmSurveyPlannedPackDate),
		PmSurveyPlannedPickupDate:           handlers.FmtDatePtr(s.PmSurveyPlannedPickupDate),
		PmSurveyPlannedDeliveryDate:         handlers.FmtDatePtr(s.PmSurveyPlannedDeliveryDate),
		PmSurveyWeightEstimate:              handlers.FmtPoundPtr(s.PmSurveyWeightEstimate),
		PmSurveyProgearWeightEstimate:       handlers.FmtPoundPtr(s.PmSurveyProgearWeightEstimate),
		PmSurveySpouseProgearWeightEstimate: handlers.FmtPoundPtr(s.PmSurveySpouseProgearWeightEstimate),
		PmSurveyNotes:                       s.PmSurveyNotes,
		PmSurveyMethod:                      s.PmSurveyMethod,
	}
	return shipmentpayload
}

// IndexShipmentsHandler returns a list of shipments
type IndexShipmentsHandler struct {
	handlers.HandlerContext
}

// Handle retrieves a list of all shipments
func (h IndexShipmentsHandler) Handle(params shipmentop.IndexShipmentsParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewIndexShipmentsForbidden()
	}

	shipments, err := models.FetchShipmentsByTSP(h.DB(), tspUser.TransportationServiceProviderID,
		params.Status, params.OrderBy, params.Limit, params.Offset)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewIndexShipmentsBadRequest()
	}

	isp := make(apimessages.IndexShipments, len(shipments))
	for i, s := range shipments {
		isp[i] = payloadForShipmentModel(s)
	}
	return shipmentop.NewIndexShipmentsOK().WithPayload(isp)
}

// GetShipmentHandler returns a particular shipment
type GetShipmentHandler struct {
	handlers.HandlerContext
}

// Handle returns a specified shipment
func (h GetShipmentHandler) Handle(params shipmentop.GetShipmentParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewGetShipmentForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewGetShipmentBadRequest()
	}

	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewGetShipmentOK().WithPayload(sp)
}

// AcceptShipmentHandler allows a TSP to accept a particular shipment
type AcceptShipmentHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h AcceptShipmentHandler) Handle(params shipmentop.AcceptShipmentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	// TODO: (cgilmer 2018_08_22) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewAcceptShipmentForbidden()
	}

	// Accept the shipment
	shipment, shipmentOffer, verrs, err := models.AcceptShipmentForTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil || verrs.HasAny() {
		if err == models.ErrFetchNotFound {
			h.Logger().Error("DB Query", zap.Error(err))
			return shipmentop.NewAcceptShipmentBadRequest()
		} else if err == models.ErrInvalidTransition {
			h.Logger().Info("Attempted to accept shipment, got invalid transition", zap.Error(err), zap.String("shipment_status", string(shipment.Status)))
			h.Logger().Info("Attempted to accept shipment offer, got invalid transition", zap.Error(err), zap.Bool("shipment_offer_accepted", *shipmentOffer.Accepted))
			return shipmentop.NewAcceptShipmentConflict()
		} else {
			h.Logger().Error("Unknown Error", zap.Error(err))
			return handlers.ResponseForVErrors(h.Logger(), verrs, err)
		}
	}

	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewAcceptShipmentOK().WithPayload(sp)
}

// RejectShipmentHandler allows a TSP to refuse a particular shipment
type RejectShipmentHandler struct {
	handlers.HandlerContext
}

// Handle refuses the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h RejectShipmentHandler) Handle(params shipmentop.RejectShipmentParams) middleware.Responder {
	// set reason, set thing
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	// TODO: (cgilmer 2018_08_22) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewRejectShipmentForbidden()
	}

	// Reject the shipment
	shipment, shipmentOffer, verrs, err := models.RejectShipmentForTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID, *params.Payload.Reason)
	if err != nil || verrs.HasAny() {
		if err == models.ErrFetchNotFound {
			h.Logger().Error("DB Query", zap.Error(err))
			return shipmentop.NewRejectShipmentBadRequest()
		} else if err == models.ErrInvalidTransition {
			h.Logger().Info("Attempted to reject shipment, got invalid transition", zap.Error(err), zap.String("shipment_status", string(shipment.Status)))
			h.Logger().Info("Attempted to reject shipment offer, got invalid transition", zap.Error(err), zap.Bool("shipment_offer_accepted", *shipmentOffer.Accepted))
			return shipmentop.NewRejectShipmentConflict()
		} else {
			h.Logger().Error("Unknown Error", zap.Error(err))
			return handlers.ResponseForVErrors(h.Logger(), verrs, err)
		}
	}

	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewRejectShipmentOK().WithPayload(sp)
}

func patchShipmentWithPayload(shipment *models.Shipment, payload *apimessages.Shipment) {
	// Premove Survey values entered by TSP agent
	requiredValue := payload.PmSurveyPlannedPackDate

	// If any PmSurvey data was sent, update all fields
	// This takes advantage of the fact that all PmSurvey data is updated at once and allows us to null out optional fields
	if requiredValue != nil {
		shipment.PmSurveyPlannedDeliveryDate = (*time.Time)(payload.PmSurveyPlannedDeliveryDate)
		shipment.PmSurveyNotes = payload.PmSurveyNotes
		shipment.PmSurveyMethod = payload.PmSurveyMethod
		shipment.PmSurveyPlannedPackDate = (*time.Time)(payload.PmSurveyPlannedPackDate)
		shipment.PmSurveyPlannedPickupDate = (*time.Time)(payload.PmSurveyPlannedPickupDate)
		shipment.PmSurveyProgearWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveyProgearWeightEstimate)
		shipment.PmSurveySpouseProgearWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveySpouseProgearWeightEstimate)
		shipment.PmSurveyWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveyWeightEstimate)
	}

	if payload.ActualWeight != nil {
		shipment.ActualWeight = handlers.PoundPtrFromInt64Ptr(payload.ActualWeight)
	}
}

// PatchShipmentHandler allows a TSP to refuse a particular shipment
type PatchShipmentHandler struct {
	handlers.HandlerContext
}

// Handle updates the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h PatchShipmentHandler) Handle(params shipmentop.PatchShipmentParams) middleware.Responder {

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewPatchShipmentForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewPatchShipmentBadRequest()
	}

	patchShipmentWithPayload(shipment, params.Update)
	verrs, err := models.SaveShipmentAndAddresses(h.DB(), shipment)

	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	shipmentPayload := payloadForShipmentModel(*shipment)
	return shipmentop.NewPatchShipmentOK().WithPayload(shipmentPayload)
}

// CreateGovBillOfLadingHandler creates a GBL move document and associates it to a shipment
type CreateGovBillOfLadingHandler struct {
	handlers.HandlerContext
}

// Handle generates the GBL PDF & uploads it as a document associated to a move doc, shipment and move
func (h CreateGovBillOfLadingHandler) Handle(params shipmentop.CreateGovBillOfLadingParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// Verify that the logged in TSP user exists
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewCreateGovBillOfLadingUnauthorized()
	}

	// Verify that TSP user is authorized to generate GBL
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewCreateGovBillOfLadingForbidden()
	}

	// Don't allow GBL generation for shipments that already have a GBL move document
	// extantGBLS, _ := models.FetchMoveDocumentsByTypeForShipment(h.DB(), session, models.MoveDocumentTypeGOVBILLOFLADING, shipmentID)
	// if len(extantGBLS) > 0 {
	// 	h.Logger().Error("There are already GBLs for this shipment.")
	// 	return shipmentop.NewCreateGovBillOfLadingBadRequest()
	// }

	// Create PDF for GBL
	gbl, err := models.FetchGovBillOfLadingExtractor(h.DB(), shipmentID)
	if err != nil {
		h.Logger().Error("Failed retrieving the GBL data.", zap.Error(err))
		return shipmentop.NewCreateGovBillOfLadingExpectationFailed()
	}
	formLayout := paperwork.Form1203Layout

	// Read in bytes from Asset pkg
	data, err := assets.Asset(formLayout.TemplateImagePath)
	if err != nil {
		h.Logger().Error("Error reading template file", zap.Error(err))
		return shipmentop.NewCreateGovBillOfLadingInternalServerError()
	}
	f, err := h.FileStorer().FileSystem().Create("something.png")
	_, err = f.Write(data)
	if err != nil {
		h.Logger().Error("Error writing template bytes to file", zap.Error(err))
		return shipmentop.NewCreateGovBillOfLadingInternalServerError()
	}
	f.Seek(0, 0)

	form, err := paperwork.NewTemplateForm(f, formLayout.FieldsLayout)
	if err != nil {
		h.Logger().Error("Error initializing GBL template form.", zap.Error(err))
		return shipmentop.NewCreateGovBillOfLadingInternalServerError()
	}

	// Populate form fields with GBL data
	err = form.DrawData(gbl)
	if err != nil {
		h.Logger().Error("Failure writing GBL data to form.", zap.Error(err))
		return shipmentop.NewCreateGovBillOfLadingInternalServerError()
	}

	aFile, err := h.FileStorer().FileSystem().Create("some name")
	if err != nil {
		h.Logger().Error("Error creating a new afero file for GBL form.", zap.Error(err))
		return shipmentop.NewCreateGovBillOfLadingInternalServerError()
	}

	err = form.Output(aFile)
	if err != nil {
		h.Logger().Error("Failure exporting GBL form to file.", zap.Error(err))
		return shipmentop.NewCreateGovBillOfLadingInternalServerError()
	}

	uploader := uploaderpkg.NewUploader(h.DB(), h.Logger(), h.FileStorer())
	upload, verrs, err := uploader.CreateUpload(nil, *tspUser.UserID, aFile)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	uploads := []models.Upload{*upload}

	// Create GBL move document associated to the shipment
	_, verrs, err = shipment.Move.CreateMoveDocument(h.DB(),
		uploads,
		&shipmentID,
		models.MoveDocumentTypeGOVBILLOFLADING,
		string("Government Bill Of Lading"),
		swag.String(""),
	)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}
	url, err := uploader.PresignedURL(upload)
	if err != nil {
		h.Logger().Error("failed to get presigned url", zap.Error(err))
		return shipmentop.NewCreateGovBillOfLadingInternalServerError()
	}

	uploadPayload := &apimessages.UploadPayload{
		ID:          handlers.FmtUUID(upload.ID),
		Filename:    swag.String(upload.Filename),
		ContentType: swag.String(upload.ContentType),
		URL:         handlers.FmtURI(url),
		Bytes:       &upload.Bytes,
		CreatedAt:   handlers.FmtDateTime(upload.CreatedAt),
		UpdatedAt:   handlers.FmtDateTime(upload.UpdatedAt),
	}
	return shipmentop.NewCreateGovBillOfLadingCreated().WithPayload(uploadPayload)

}

// GetShipmentContactDetailsHandler allows a TSP to accept a particular shipment
type GetShipmentContactDetailsHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h GetShipmentContactDetailsHandler) Handle(p shipmentop.GetShipmentContactDetailsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// GetShipmentClaimsHandler allows a TSP to accept a particular shipment
type GetShipmentClaimsHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h GetShipmentClaimsHandler) Handle(p shipmentop.GetShipmentClaimsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// GetShipmentDocumentsHandler allows a TSP to accept a particular shipment
type GetShipmentDocumentsHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h GetShipmentDocumentsHandler) Handle(p shipmentop.GetShipmentDocumentsParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}

// CreateShipmentDocumentHandler allows a TSP to accept a particular shipment
type CreateShipmentDocumentHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h CreateShipmentDocumentHandler) Handle(p shipmentop.CreateShipmentDocumentParams) middleware.Responder {
	return middleware.NotImplemented("operation .shipmentContactDetails has not yet been implemented")
}
