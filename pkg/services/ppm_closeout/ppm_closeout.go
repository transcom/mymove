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

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope(&ppmCloseoutObj)).
		EagerPreload(
			"ID",
			"PlannedMoveDate",
			"ActualMoveDate",
			"Miles",
			"EstimatedWeight",
			"ActualWeight",
			"ProGearWeightCustomer",
			"ProGearWeightSpouse",
			"GrossIncentive",
		).
		Join("mto_shipments", "mto_shipments.id = ppm_shipments.shipment_id").
		Find(&ppmShipment, ppmShipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(ppmShipmentID, "while looking for PPMShipment")
		default:
			return nil, apperror.NewQueryError("PPMShipment", err, "unable to find PPMShipment")
		}
	}

	return &ppmCloseoutObj, nil
}
