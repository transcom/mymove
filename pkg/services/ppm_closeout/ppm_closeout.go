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
	var mtoShipment models.MTOShipment

	errPPM := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"ID",
			"ShipmentID",
			"ExpectedDepartureDate",
			"ActualMoveDate",
			"EstimatedWeight",
			"HasProGear",
			"ProGearWeight",
			"SpouseProGearWeight",
			"FinalIncentive",
		).
		Find(&ppmShipment, ppmShipmentID)

	if errPPM != nil {
		switch errPPM {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(ppmShipmentID, "while looking for PPMShipment")
		default:
			return nil, apperror.NewQueryError("PPMShipment", errPPM, "unable to find PPMShipment")
		}
	}

	mtoShipmentID := &ppmShipment.ShipmentID
	errMTO := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"ID",
			"ScheduledPickupDate",
			"ActualPickupDate",
			"Distance",
			"PrimeActualWeight",
		).
		Find(&mtoShipment, mtoShipmentID)

	if errMTO != nil {
		switch errMTO {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(*mtoShipmentID, "while looking for MTOShipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", errMTO, "unable to find MTOShipment")
		}
	}

	ppmCloseoutObj.ID = &ppmShipmentID
	ppmCloseoutObj.PlannedMoveDate = mtoShipment.ScheduledPickupDate
	ppmCloseoutObj.ActualMoveDate = mtoShipment.ActualPickupDate
	ppmCloseoutObj.Miles = (*int)(mtoShipment.Distance)
	ppmCloseoutObj.EstimatedWeight = ppmShipment.EstimatedWeight
	ppmCloseoutObj.ActualWeight = mtoShipment.PrimeActualWeight
	ppmCloseoutObj.ProGearWeightCustomer = ppmShipment.ProGearWeight
	ppmCloseoutObj.ProGearWeightSpouse = ppmShipment.SpouseProGearWeight
	ppmCloseoutObj.GrossIncentive = ppmShipment.FinalIncentive
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
