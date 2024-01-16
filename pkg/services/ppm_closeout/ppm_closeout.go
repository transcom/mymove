package ppmcloseout

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type ppmCloseoutFetcher struct{}

func NewPPMCloseoutFetcher() services.PPMCloseoutFetcher {
	return &ppmCloseoutFetcher{}
}

func (p *ppmCloseoutFetcher) GetPPMCloseout(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.PPMCloseout, error) {
	var ppmCloseoutObj models.PPMCloseout
	var ppmShipment models.PPMShipment

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"ID",
		).
		Find(&ppmShipment, ppmShipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(ppmShipmentID, "while looking for PPMShipment")
		default:
			return nil, apperror.NewQueryError("PPMShipment", err, "unable to find PPMShipment")
		}
	}
	ppmCloseoutObj.ID = &ppmShipmentID
	ppmCloseoutObj.PlannedMoveDate = nil
	ppmCloseoutObj.ActualMoveDate = nil
	ppmCloseoutObj.Miles = nil
	ppmCloseoutObj.EstimatedWeight = nil
	ppmCloseoutObj.ActualWeight = nil
	ppmCloseoutObj.ProGearWeightCustomer = nil
	ppmCloseoutObj.ProGearWeightSpouse = nil
	ppmCloseoutObj.GrossIncentive = nil
	ppmCloseoutObj.GCC = nil
	ppmCloseoutObj.AOA = nil
	ppmCloseoutObj.RemainingReimbursementOwed = nil
	ppmCloseoutObj.HaulPrice = nil
	ppmCloseoutObj.HaulFSC = nil
	ppmCloseoutObj.DOP = nil
	ppmCloseoutObj.DDP = nil
	ppmCloseoutObj.PackPrice = nil
	ppmCloseoutObj.UnpackPrice = nil
	ppmCloseoutObj.SITReimbursement = nil

	return &ppmCloseoutObj, nil
}
