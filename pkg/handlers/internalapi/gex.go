package internalapi

import (
	"bytes"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	gexop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/gex"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
)

// SendGexRequestHandler sends a request to GEX
type SendGexRequestHandler struct {
	handlers.HandlerContext
}

// Handle sends a request to GEX
func (h SendGexRequestHandler) Handle(params gexop.SendGexRequestParams) middleware.Responder {
	transactionName := *params.SendGexRequestPayload.TransactionName
	transactionBody := *params.SendGexRequestPayload.TransactionBody

	// Ensure that the transaction body ends with a newline, otherwise the GEX
	// EDI parser will fail silently
	transactionBody = strings.TrimSpace(transactionBody) + "\n"

	resp, err := h.GexSender().Call(transactionBody, transactionName)
	if err != nil {
		h.Logger().Error("Sending Invoice to Gex", zap.Error(err))
		return gexop.NewSendGexRequestInternalServerError()
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	responseBody := buf.String()

	responsePayload := internalmessages.GexResponsePayload{
		GexResponse: resp.Status + "; " + responseBody,
	}
	return gexop.NewSendGexRequestOK().WithPayload(&responsePayload)
}
