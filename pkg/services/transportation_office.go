package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

//go:generate mockery --name TransportationOfficesFetcher
type TransportationOfficesFetcher interface {
	GetTransportationOffices(appCtx appcontext.AppContext, search string, forPpm bool, forAdminOfficeUserReqFilter bool) (*models.TransportationOffices, error)
	GetTransportationOffice(appCtx appcontext.AppContext, transportationOfficeID uuid.UUID, includeOnlyPPMCloseoutOffices bool) (*models.TransportationOffice, error)
	GetAllGBLOCs(appCtx appcontext.AppContext) (*models.GBLOCs, error)
	GetCounselingOffices(appCtx appcontext.AppContext, dutyLocationID uuid.UUID, serviceMemberID uuid.UUID) (*models.TransportationOffices, error)
	FindCounselingOfficeForPrimeCounseled(appCtx appcontext.AppContext, dutyLocationID uuid.UUID, serviceMemberID uuid.UUID) (*models.TransportationOffice, error)
}
