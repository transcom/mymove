package publicapi

import (
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForAccessorialModels(s []models.Accessorial) apimessages.Accessorials {
	payloads := make(apimessages.Accessorials, len(s))

	for i, acc := range s {
		payloads[i] = payloadForAccessorialModel(&acc)
	}

	return payloads
}

func payloadForAccessorialModel(a *models.Accessorial) *apimessages.Accessorial {
	if a == nil {
		return nil
	}

	return &apimessages.Accessorial{
		ID:           *handlers.FmtUUID(a.ID),
		Code:         *handlers.FmtString(a.Code),
		DiscountType: *handlers.FmtString(string(a.DiscountType)),
		Item:         *handlers.FmtString(a.Item),
		Location:     apimessages.AccessorialLocation(string(a.AllowedLocation)),
		RefCode:      *handlers.FmtString(string(a.RateRefCode)),
		Uom1:         *handlers.FmtString(string(a.MeasurementUnit1)),
		Uom2:         *handlers.FmtString(string(a.MeasurementUnit2)),
		CreatedAt:    *handlers.FmtDateTime(a.CreatedAt),
		UpdatedAt:    *handlers.FmtDateTime(a.UpdatedAt),
	}
}
