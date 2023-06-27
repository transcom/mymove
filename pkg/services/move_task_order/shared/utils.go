package shared

import (
	"database/sql"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func SetMTOQueryFilters(query *pop.Query, searchParams *services.MoveTaskOrderFetcherParams) {
	// Always exclude hidden moves by default:
	if searchParams == nil {
		query.Where("show = TRUE")
	} else {
		if searchParams.IsAvailableToPrime {
			query.Where("available_to_prime_at IS NOT NULL")
		}

		// This value defaults to false - we want to make sure including hidden moves needs to be explicitly requested.
		if !searchParams.IncludeHidden {
			query.Where("show = TRUE")
		}

		if searchParams.Since != nil {
			query.Where("updated_at > ?", *searchParams.Since)
		}
	}
	// No return since this function uses pointers to modify the referenced query directly
}

// FetchReweigh retrieves a reweigh for a given shipment id
func FetchReweigh(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*models.Reweigh, error) {
	reweigh := &models.Reweigh{}
	err := appCtx.DB().
		Where("shipment_id = ?", shipmentID).
		First(reweigh)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.Reweigh{}, nil
		default:
			return &models.Reweigh{}, err
		}
	}
	return reweigh, nil
}
