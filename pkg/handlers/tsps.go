package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/markbates/pop"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/restapi/apioperations"
)

// NewCreateIssueHandler returns a new CreateIssueHandler
func NewTSPIndexHandler(db *pop.Connection, logger *zap.Logger) apioperations.IndexTSPsHandler {
	return apioperations.IndexTSPsHandlerFunc(func(params apioperations.IndexTSPsParams) middleware.Responder {
		return middleware.NotImplemented("operation .IndexTSPs has not yet been implemented")
	})
}

func NewTSPShipmentsHandler(db *pop.Connection, logger *zap.Logger) apioperations.TspShipmentsHandler {
	return apioperations.TspShipmentsHandlerFunc(func(params apioperations.TspShipmentsParams) middleware.Responder {
		return middleware.NotImplemented("operation .TspShipments has not yet been implemented")
	})

}
