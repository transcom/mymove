package publicapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	tspop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/transportation_service_provider"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForTransportationServiceProviderModel(t models.TransportationServiceProvider) *apimessages.TransportationServiceProvider {
	transportationServiceProviderPayload := &apimessages.TransportationServiceProvider{
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
	return transportationServiceProviderPayload
}

// GetTransportationServiceProviderHandler returns a TSP for a shipment
type GetTransportationServiceProviderHandler struct {
	handlers.HandlerContext
}

// Handle getting the tsp for a shipment
func (h GetTransportationServiceProviderHandler) Handle(params tspop.GetTransportationServiceProviderParams) middleware.Responder {
	var shipment *models.Shipment
	var err error
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)
	shipmentID, _ := uuid.FromString(params.ShipmentID.String())

	if session.IsTspUser() {
		tspUser, fetchTSPByUserErr := models.FetchTspUserByID(h.DB(), session.TspUserID)
		if fetchTSPByUserErr != nil {
			logger.Error("DB Query", zap.Error(fetchTSPByUserErr))
			return tspop.NewGetTransportationServiceProviderForbidden()
		}

		shipment, err = models.FetchShipmentByTSP(h.DB(), tspUser.TransportationServiceProviderID, shipmentID)
		if err != nil {
			return tspop.NewGetTransportationServiceProviderNotFound()
		}
	} else if session.IsOfficeUser() {
		shipment, err = models.FetchShipment(h.DB(), session, shipmentID)
		if err != nil {
			return tspop.NewGetTransportationServiceProviderNotFound()
		}
	} else if session.IsServiceMember() {
		shipment, err = models.FetchShipment(h.DB(), session, shipmentID)
		if err != nil {
			if err == models.ErrFetchForbidden {
				return tspop.NewGetTransportationServiceProviderForbidden()
			}
			return tspop.NewGetTransportationServiceProviderNotFound()
		}
	} else {
		return tspop.NewGetTransportationServiceProviderForbidden()
	}

	// Office Users and Service Members aren't guaranteed that the shipment has been awarded
	// TSP Users that reach this point are because otherwise they are forbidden from viewing
	// If the Shipment is not yet awarded then the TSP ID will be a nil UUID
	transportationServiceProviderID := shipment.CurrentTransportationServiceProviderID()
	if transportationServiceProviderID == uuid.Nil {
		return tspop.NewGetTransportationServiceProviderNotFound()
	}

	transportationServiceProvider, err := models.FetchTransportationServiceProvider(h.DB(), transportationServiceProviderID)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	transportationServiceProviderPayload := payloadForTransportationServiceProviderModel(*transportationServiceProvider)
	return tspop.NewGetTransportationServiceProviderOK().WithPayload(transportationServiceProviderPayload)
}
