package ppmcloseout

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

type PPMCloseout struct{}

func NewPPMCloseout() *PPMCloseout {
	return &PPMCloseout{}
}

func (p *PPMCloseout) GetPPMCloseout(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) *models.PPMCloseout {
	ppmCloseoutObj := &models.PPMCloseout{}
	var queryResults []map[string]any

	err := appCtx.DB().Q().Select("mto_shipments.distance", "mto_shipments.prime_actual_weight", "ppm_shipments.*").Join("ppm_shipments", "ppm_shipments.shipment_id = mto_shipments.id").Find(&queryResults, ppmShipmentID)
	if err != nil {
		return nil
	}

	// Build the ppmCloseoutObj object
	for _, result := range queryResults {
		ppmCloseoutObj.ID = result["id"].(*uuid.UUID)
		ppmCloseoutObj.PlannedMoveDate = result["expected_departure_date"].(*time.Time)
		ppmCloseoutObj.ActualMoveDate = result["actual_move_date"].(*time.Time)
		ppmCloseoutObj.Miles = result["distance"].(*unit.Miles)
		ppmCloseoutObj.EstimatedWeight = result["estimated_weight"].(*unit.Pound)
		ppmCloseoutObj.ActualWeight = result["prime_actual_weight"].(*unit.Pound)
		ppmCloseoutObj.ProGearWeightCustomer = result[""].(*unit.Pound)
		ppmCloseoutObj.ProGearWeightSpouse = result[""].(*unit.Pound)
		ppmCloseoutObj.GrossIncentive = result[""].(*unit.Cents)
		ppmCloseoutObj.GCC = result[""].(*unit.Cents)
		ppmCloseoutObj.AOA = result[""].(*unit.Cents)
		ppmCloseoutObj.RemainingReimbursementOwed = result[""].(*unit.Cents)
		ppmCloseoutObj.HaulPrice = result[""].(*unit.Cents)
		ppmCloseoutObj.HaulFSC = result[""].(*unit.Cents)
		ppmCloseoutObj.DOP = result[""].(*unit.Cents)
		ppmCloseoutObj.DDP = result[""].(*unit.Cents)
		ppmCloseoutObj.PackPrice = result[""].(*unit.Cents)
		ppmCloseoutObj.UnpackPrice = result[""].(*unit.Cents)
		ppmCloseoutObj.SITReimbursement = result[""].(*unit.Cents)
	}

	return ppmCloseoutObj
}
