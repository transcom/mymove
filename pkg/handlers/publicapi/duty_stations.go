package publicapi

import (
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForDutyStationModel(station models.DutyStation) *apimessages.DutyStation {
	payload := apimessages.DutyStation{
		ID:        handlers.FmtUUID(station.ID),
		CreatedAt: handlers.FmtDateTime(station.CreatedAt),
		UpdatedAt: handlers.FmtDateTime(station.UpdatedAt),
		Name:      swag.String(station.Name),
		Address:   payloadForAddressModel(&station.Address),
	}

	return &payload
}
