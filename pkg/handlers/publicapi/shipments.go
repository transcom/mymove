package publicapi

import (
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/awardqueue"
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
		ID: *handlers.FmtUUID(s.ID),
		TrafficDistributionList:             payloadForTrafficDistributionListModel(s.TrafficDistributionList),
		ServiceMember:                       payloadForServiceMemberModel(&s.ServiceMember),
		ActualPickupDate:                    handlers.FmtDatePtr(s.ActualPickupDate),
		ActualPackDate:                      handlers.FmtDatePtr(s.ActualPackDate),
		ActualDeliveryDate:                  handlers.FmtDatePtr(s.ActualDeliveryDate),
		CreatedAt:                           strfmt.DateTime(s.CreatedAt),
		UpdatedAt:                           strfmt.DateTime(s.UpdatedAt),
		SourceGbloc:                         apimessages.GBLOC(*s.SourceGBLOC),
		DestinationGbloc:                    apimessages.GBLOC(*s.DestinationGBLOC),
		GblNumber:                           s.GBLNumber,
		Market:                              apimessages.ShipmentMarket(*s.Market),
		BookDate:                            *handlers.FmtDatePtr(s.BookDate),
		RequestedPickupDate:                 *handlers.FmtDatePtr(s.RequestedPickupDate),
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
		WeightEstimate:                      handlers.FmtPoundPtr(s.WeightEstimate),
		ProgearWeightEstimate:               handlers.FmtPoundPtr(s.ProgearWeightEstimate),
		SpouseProgearWeightEstimate:         handlers.FmtPoundPtr(s.SpouseProgearWeightEstimate),
		NetWeight:                           handlers.FmtPoundPtr(s.NetWeight),
		GrossWeight:                         handlers.FmtPoundPtr(s.GrossWeight),
		TareWeight:                          handlers.FmtPoundPtr(s.TareWeight),
		PmSurveyConductedDate:               handlers.FmtDatePtr(s.PmSurveyConductedDate),
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

	go awardqueue.NewAwardQueue(h.DB(), h.Logger()).Run()

	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewRejectShipmentOK().WithPayload(sp)
}

// TransportShipmentHandler allows a TSP to start transporting a particular shipment
type TransportShipmentHandler struct {
	handlers.HandlerContext
}

// Handle accepts the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h TransportShipmentHandler) Handle(params shipmentop.TransportShipmentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewTransportShipmentForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewTransportShipmentBadRequest()
	}

	actualPickupDate := (time.Time)(*params.Payload.ActualPickupDate)

	err = shipment.Transport(actualPickupDate)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	verrs, err := h.DB().ValidateAndUpdate(shipment)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}
	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewTransportShipmentOK().WithPayload(sp)
}

// DeliverShipmentHandler allows a TSP to start transporting a particular shipment
type DeliverShipmentHandler struct {
	handlers.HandlerContext
}

// Handle delivers the shipment - checks that currently logged in user is authorized to act for the TSP assigned the shipment
func (h DeliverShipmentHandler) Handle(params shipmentop.DeliverShipmentParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	// TODO: (cgilmer 2018_07_25) This is an extra query we don't need to run on every request. Put the
	// TransportationServiceProviderID into the session object after refactoring the session code to be more readable.
	// See original commits in https://github.com/transcom/mymove/pull/802
	tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewDeliverShipmentForbidden()
	}

	shipment, err := models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return shipmentop.NewDeliverShipmentBadRequest()
	}

	actualDeliveryDate := (time.Time)(*params.Payload.ActualDeliveryDate)

	err = shipment.Deliver(actualDeliveryDate)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}
	verrs, err := h.DB().ValidateAndUpdate(shipment)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}
	sp := payloadForShipmentModel(*shipment)
	return shipmentop.NewDeliverShipmentOK().WithPayload(sp)
}

func patchShipmentWithPayload(shipment *models.Shipment, payload *apimessages.Shipment) {

	// PM Survey fields may be updated individually in the Dates panel and so cannot be lumped into one update
	if payload.PmSurveyConductedDate != nil {
		shipment.PmSurveyConductedDate = (*time.Time)(payload.PmSurveyConductedDate)
	}

	if payload.PmSurveyPlannedDeliveryDate != nil {
		shipment.PmSurveyPlannedDeliveryDate = (*time.Time)(payload.PmSurveyPlannedDeliveryDate)
	}

	if payload.PmSurveyMethod != "" {
		shipment.PmSurveyMethod = payload.PmSurveyMethod
	}

	if payload.PmSurveyPlannedPackDate != nil {
		shipment.PmSurveyPlannedPackDate = (*time.Time)(payload.PmSurveyPlannedPackDate)
	}

	if payload.PmSurveyPlannedPickupDate != nil {
		shipment.PmSurveyPlannedPickupDate = (*time.Time)(payload.PmSurveyPlannedPickupDate)
	}

	if payload.PmSurveyNotes != nil {
		shipment.PmSurveyNotes = payload.PmSurveyNotes
	}

	if payload.PmSurveyProgearWeightEstimate != nil {
		shipment.PmSurveyProgearWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveyProgearWeightEstimate)
	}

	if payload.PmSurveySpouseProgearWeightEstimate != nil {
		shipment.PmSurveySpouseProgearWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveySpouseProgearWeightEstimate)
	}

	if payload.PmSurveyWeightEstimate != nil {
		shipment.PmSurveyWeightEstimate = handlers.PoundPtrFromInt64Ptr(payload.PmSurveyWeightEstimate)
	}

	if payload.NetWeight != nil {
		shipment.NetWeight = handlers.PoundPtrFromInt64Ptr(payload.NetWeight)
	}

	if payload.GrossWeight != nil {
		shipment.GrossWeight = handlers.PoundPtrFromInt64Ptr(payload.GrossWeight)
	}

	if payload.TareWeight != nil {
		shipment.TareWeight = handlers.PoundPtrFromInt64Ptr(payload.TareWeight)
	}

	if payload.ActualPickupDate != nil {
		shipment.ActualPickupDate = (*time.Time)(payload.ActualPickupDate)
	}

	if payload.ActualPackDate != nil {
		shipment.ActualPackDate = (*time.Time)(payload.ActualPackDate)
	}

	if payload.ActualDeliveryDate != nil {
		shipment.ActualDeliveryDate = (*time.Time)(payload.ActualDeliveryDate)
	}

	if payload.PickupAddress != nil {
		if shipment.PickupAddress == nil {
			shipment.PickupAddress = addressModelFromPayload(payload.PickupAddress)
		} else {
			updateAddressWithPayload(shipment.PickupAddress, payload.PickupAddress)
		}
	}
	if payload.HasSecondaryPickupAddress == false {
		shipment.SecondaryPickupAddress = nil
	} else if payload.HasSecondaryPickupAddress == true {
		if payload.SecondaryPickupAddress != nil {
			if shipment.SecondaryPickupAddress == nil {
				shipment.SecondaryPickupAddress = addressModelFromPayload(payload.SecondaryPickupAddress)
			} else {
				updateAddressWithPayload(shipment.SecondaryPickupAddress, payload.SecondaryPickupAddress)
			}
		}
	}
	shipment.HasSecondaryPickupAddress = payload.HasSecondaryPickupAddress
	if payload.HasDeliveryAddress == false {
		shipment.DeliveryAddress = nil
	} else if payload.HasDeliveryAddress == true {
		if payload.DeliveryAddress != nil {
			if shipment.DeliveryAddress == nil {
				shipment.DeliveryAddress = addressModelFromPayload(payload.DeliveryAddress)
			} else {
				updateAddressWithPayload(shipment.DeliveryAddress, payload.DeliveryAddress)
			}
		}
	}
	shipment.HasDeliveryAddress = payload.HasDeliveryAddress
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

// CreateGovBillOfLadingHandler creates a GBL PDF & uploads it as a document associated to a move doc, shipment and move
type CreateGovBillOfLadingHandler struct {
	handlers.HandlerContext
}

// Handle generates the GBL PDF & uploads it as a document associated to a move doc, shipment and move
func (h CreateGovBillOfLadingHandler) Handle(params shipmentop.CreateGovBillOfLadingParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// Verify that the TSP user is authorized to update move doc
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())
	tspUser, shipment, err := models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentID)
	if err != nil {
		if err.Error() == "USER_UNAUTHORIZED" {
			h.Logger().Error("DB Query", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
		if err.Error() == "FETCH_FORBIDDEN" {
			h.Logger().Error("DB Query", zap.Error(err))
			return handlers.ResponseForError(h.Logger(), err)
		}
		return handlers.ResponseForError(h.Logger(), err)
	}
	// Don't allow GBL generation for shipments that already have a GBL move document
	extantGBLS, _ := models.FetchMoveDocumentsByTypeForShipment(h.DB(), session, models.MoveDocumentTypeGOVBILLOFLADING, shipmentID)
	if len(extantGBLS) > 0 {
		h.Logger().Error("There are already GBLs for this shipment.")
		return shipmentop.NewCreateGovBillOfLadingBadRequest()
	}

	// Create PDF for GBL
	gbl, err := models.FetchGovBillOfLadingExtractor(h.DB(), shipmentID)
	if err != nil {
		// TODO: (andrea) Pass info of exactly what is missing in custom error message
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
		string(apimessages.SelectedMoveTypeHHG),
	)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}
	url, err := uploader.PresignedURL(upload)
	if err != nil {
		h.Logger().Error("failed to get presigned url", zap.Error(err))
		return shipmentop.NewCreateGovBillOfLadingInternalServerError()
	}

	// TODO: (andrea) Return a document payload instead, once the HHG document is defined in public swagger
	// This one is copy pasted from internal.yaml to api.yaml :/
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
