package internal

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	transportationofficeop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/transportation_offices"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers/utils"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForTransportationOfficeModel(office models.TransportationOffice) *internalmessages.TransportationOffice {
	var phoneLines []string
	for _, phoneLine := range office.PhoneLines {
		if phoneLine.Type == "voice" {
			phoneLines = append(phoneLines, phoneLine.Number)
		}
	}

	payload := &internalmessages.TransportationOffice{
		ID:         utils.FmtUUID(office.ID),
		CreatedAt:  utils.FmtDateTime(office.CreatedAt),
		UpdatedAt:  utils.FmtDateTime(office.UpdatedAt),
		Name:       swag.String(office.Name),
		Gbloc:      office.Gbloc,
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
		return utils.ResponseForError(h.logger, err)
	}
	transportationOfficePayload := payloadForTransportationOfficeModel(transportationOffice)

	return transportationofficeop.NewShowDutyStationTransportationOfficeOK().WithPayload(transportationOfficePayload)
}
