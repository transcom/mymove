package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	progearops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/ppm"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateProgearWeightTicketHandler
type UpdateProgearWeightTicketHandler struct {
	handlers.HandlerConfig
	progearUpdater services.ProgearWeightTicketUpdater
}

func (h UpdateProgearWeightTicketHandler) Handle(params progearops.UpdateProGearWeightTicketParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.UpdateProGearWeightTicket

			progearWeightTicket := payloads.ProgearWeightTicketModelFromUpdate(payload)

			progearWeightTicket.ID = uuid.FromStringOrNil(params.ProGearWeightTicketID.String())

			updatedProgearWeightTicket, _ := h.progearUpdater.UpdateProgearWeightTicket(appCtx, *progearWeightTicket, params.IfMatch)

			returnPayload := payloads.ProGearWeightTicket(h.FileStorer(), updatedProgearWeightTicket)
			return progearops.NewUpdateProGearWeightTicketOK().WithPayload(returnPayload), nil
		})
}
