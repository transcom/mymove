package internalapi

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForServiceAgentModel(s models.ServiceAgent) *internalmessages.ServiceAgent {
	serviceAgentPayload := &internalmessages.ServiceAgent{
		ID:               *handlers.FmtUUID(s.ID),
		ShipmentID:       *handlers.FmtUUID(s.ShipmentID),
		CreatedAt:        strfmt.DateTime(s.CreatedAt),
		UpdatedAt:        strfmt.DateTime(s.UpdatedAt),
		Role:             internalmessages.ServiceAgentRole(s.Role),
		PointOfContact:   handlers.FmtString(s.PointOfContact),
		Email:            s.Email,
		PhoneNumber:      s.PhoneNumber,
		FaxNumber:        s.FaxNumber,
		EmailIsPreferred: s.EmailIsPreferred,
		PhoneIsPreferred: s.PhoneIsPreferred,
		Notes:            s.Notes,
	}
	return serviceAgentPayload
}
