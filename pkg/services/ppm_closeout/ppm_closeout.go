package ppmcloseout

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type ppmCloseoutFetcher struct{}

func NewPPMCloseoutFetcher() services.PPMCloseoutFetcher {
	return &ppmCloseoutFetcher{}
}

func (p *ppmCloseoutFetcher) GetPPMCloseout(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.PPMCloseout, error) {
	var ppmCloseoutObj models.PPMCloseout
	var ppmShipmentQueryResult models.PPMShipment
	var mtoShipmentQueryResult models.MTOShipment

	ppmShipmentQueryErr := appCtx.DB().Find(&ppmShipmentQueryResult, ppmShipmentID)
	if ppmShipmentQueryErr != nil {
		if ppmShipmentQueryResult.ID != uuid.Nil {
			return nil, apperror.NewQueryError("PPMShipment", ppmShipmentQueryErr, "unable to find PPMShipment")
		}
		return nil, apperror.NewNotFoundError(ppmShipmentID, "while looking for PPMShipment")
	}

	mtoShipmentQueryResultErr := appCtx.DB().Find(&mtoShipmentQueryResult, &ppmShipmentQueryResult.ShipmentID)
	if mtoShipmentQueryResultErr != nil {
		if mtoShipmentQueryResult.ID != uuid.Nil {
			return nil, apperror.NewQueryError("MTOShipment", mtoShipmentQueryResultErr, "unable to find associated MTOShipment")
		}
		return nil, apperror.NewNotFoundError(ppmShipmentQueryResult.ShipmentID, "while looking for MTOShipment")
	}

	ppmCloseoutObj.ID = ppmShipmentQueryResult.ID
	ppmCloseoutObj.PlannedMoveDate = &ppmShipmentQueryResult.ExpectedDepartureDate
	ppmCloseoutObj.ActualMoveDate = ppmShipmentQueryResult.ActualMoveDate
	ppmCloseoutObj.Miles = mtoShipmentQueryResult.Distance
	ppmCloseoutObj.EstimatedWeight = ppmShipmentQueryResult.EstimatedWeight
	ppmCloseoutObj.ActualWeight = mtoShipmentQueryResult.PrimeActualWeight
	ppmCloseoutObj.ProGearWeightCustomer = ppmShipmentQueryResult.ProGearWeight
	ppmCloseoutObj.ProGearWeightSpouse = ppmShipmentQueryResult.SpouseProGearWeight
	ppmCloseoutObj.GrossIncentive = ppmShipmentQueryResult.FinalIncentive

	return &ppmCloseoutObj, nil
}
