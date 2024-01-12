package ppmcloseout

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type ppmCloseoutFetcher struct{}

func NewPPMCloseoutFetcher() services.PPMCloseoutFetcher {
	return &ppmCloseoutFetcher{}
}

func (p *ppmCloseoutFetcher) GetPPMCloseout(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.PPMCloseout, error) {
	var ppmCloseoutObj models.PPMCloseout
	var results []map[string]any

	err := appCtx.DB().Q().RawQuery("SELECT ms.distance, ms.priime_actual_weight, ps.* FROM ppm_shipments ps join on ps.shipment_id = ms.id WHERE ps.id = ?", ppmShipmentID).All(&results)

	if err != nil {
		if len(results) == 0 {
			return nil, apperror.NewNotFoundError(ppmShipmentID, "while looking for PPMShipment")
		}
		return nil, apperror.NewQueryError("PPMShipment", err, "unable to find PPMShipment")
	}

	ppmCloseoutObj.ID = results[0]["id"].(*uuid.UUID)
	ppmCloseoutObj.PlannedMoveDate = results[0]["expected_departure_date"].(*time.Time)
	ppmCloseoutObj.ActualMoveDate = results[0]["actual_move_date"].(*time.Time)
	ppmCloseoutObj.Miles = results[0]["distance"].(*int)
	ppmCloseoutObj.EstimatedWeight = results[0]["estimated_weight"].(*unit.Pound)
	ppmCloseoutObj.ActualWeight = results[0]["prime_actual_weight"].(*unit.Pound)
	ppmCloseoutObj.ProGearWeightCustomer = results[0]["pro_gear_weight"].(*unit.Pound)
	ppmCloseoutObj.ProGearWeightSpouse = results[0]["pro_gear_weight_spouse"].(*unit.Pound)
	ppmCloseoutObj.GrossIncentive = results[0]["final_incentive"].(*unit.Cents)

	return &ppmCloseoutObj, nil
}
