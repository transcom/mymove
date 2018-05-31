package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	transportationofficeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/transportation_offices"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForTransportationOfficeModel(office *models.TransportationOffice) *internalmessages.TransportationOffice {
	var phoneLines []string
	for _, phoneLine := range office.PhoneLines {
		if phoneLine.Type == "voice" {
			phoneLines = append(phoneLines, phoneLine.Number)
		}
	}

	payload := &internalmessages.TransportationOffice{
		ID:         fmtUUID(office.ID),
		CreatedAt:  fmtDateTime(office.CreatedAt),
		UpdatedAt:  fmtDateTime(office.UpdatedAt),
		Name:       swag.String(office.Name),
		Address:    payloadForAddressModel(&office.Address),
		PhoneLines: phoneLines,
	}
	return payload
}

// ShowDutyStationTransportationOfficeHandler returns the transportation office for a duty station ID
type ShowDutyStationTransportationOfficeHandler HandlerContext

// Handle retrieves the transportation office in the system for a given duty station ID
func (h ShowDutyStationTransportationOfficeHandler) Handle(params transportationofficeop.ShowDutyStationTransportationOfficeParams) middleware.Responder {
	dutyStationID, _ := uuid.FromString(params.DutyStationID.String())
	transportationOffice, err := models.FetchDutyStationTransportationOffice(h.db, dutyStationID)
	if err != nil {
		return responseForError(h.logger, err)
	}
	transportationOfficePayload := payloadForTransportationOfficeModel(transportationOffice)

	return transportationofficeop.NewShowDutyStationTransportationOfficeOK().WithPayload(transportationOfficePayload)
}
