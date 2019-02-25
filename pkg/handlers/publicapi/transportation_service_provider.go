package publicapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	tspop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/transportation_service_provider"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForTransportationServiceProviderModel(t models.TransportationServiceProvider) *apimessages.TransportationServiceProvider {
	transporationServiceProviderPayload := &apimessages.TransportationServiceProvider{
		ID:                       *handlers.FmtUUID(t.ID),
		StandardCarrierAlphaCode: handlers.FmtString(t.StandardCarrierAlphaCode),
		CreatedAt:                strfmt.DateTime(t.CreatedAt),
		UpdatedAt:                strfmt.DateTime(t.UpdatedAt),
		Enrolled:                 t.Enrolled,
		Name:                     t.Name,
		PocGeneralName:           t.PocGeneralName,
		PocGeneralEmail:          t.PocGeneralEmail,
		PocGeneralPhone:          t.PocGeneralPhone,
		PocClaimsName:            t.PocClaimsName,
		PocClaimsEmail:           t.PocClaimsEmail,
		PocClaimsPhone:           t.PocClaimsPhone,
	}
	return transporationServiceProviderPayload
}

// GetTransportationServiceProviderHandler returns a TSP for a shipment
type GetTransportationServiceProviderHandler struct {
	handlers.HandlerContext
}

// Handle getting the tsp for a shipment
func (h GetTransportationServiceProviderHandler) Handle(params tspop.GetTransportationServiceProviderParams) middleware.Responder {
	var shipment *models.Shipment
	var err error
	session := auth.SessionFromRequestContext(params.HTTPRequest)
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	if session.IsTspUser() {
		// TODO (2018_08_27 cgilmer): Find a way to check Shipment belongs to TSP without 2 queries
		tspUser, err := models.FetchTspUserByID(h.DB(), session.TspUserID)
		if err != nil {
			h.Logger().Error("DB Query", zap.Error(err))
			// return tspop.NewGetTransporationServiceProviderForbidden()
		}

		shipment, err = models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
		if err != nil {
			handlers.ResponseForError(h.Logger(), err)
			return tspop.NewGetTransportationServiceProviderBadRequest()
		}
	} else if session.IsOfficeUser() {
		shipment, err = models.FetchShipment(h.DB(), session, shipmentID)
		if err != nil {
			handlers.ResponseForError(h.Logger(), err)
			return tspop.NewGetTransportationServiceProviderBadRequest()
		}
	} else if session.IsServiceMember() {
		shipment, err = models.FetchShipment(h.DB(), session, shipmentID)
		if err != nil {
			handlers.ResponseForError(h.Logger(), err)
			if err == models.ErrFetchForbidden {
				return tspop.NewGetTransportationServiceProviderForbidden()
			}
			return tspop.NewGetTransportationServiceProviderBadRequest()
		}
	} else {
		return tspop.NewGetTransportationServiceProviderForbidden()
	}

	transportationServiceProviderID := shipment.CurrentTransportationServiceProviderID()
	if transportationServiceProviderID == uuid.Nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	transportationServiceProvider, err := models.FetchTransportationServiceProvider(h.DB(), transportationServiceProviderID)
	if err != nil {
		return handlers.ResponseForError(h.Logger(), err)
	}

	transportationServiceProviderPayload := payloadForTransportationServiceProviderModel(*transportationServiceProvider)
	return tspop.NewGetTransportationServiceProviderOK().WithPayload(transportationServiceProviderPayload)
}
