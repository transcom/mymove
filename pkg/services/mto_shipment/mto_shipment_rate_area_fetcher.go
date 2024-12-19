package mtoshipment

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/exp/slices"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type mtoShipmentRateAreaFetcher struct {
}

// NewMTOShipmentFetcher creates a new MTOShipmentFetcher struct that supports ListMTOShipments
func NewMTOShipmentRateAreaFetcher() services.ShipmentRateAreaFinder {
	return &mtoShipmentRateAreaFetcher{}
}

func (f mtoShipmentRateAreaFetcher) GetPrimeMoveShipmentOconusRateArea(appCtx appcontext.AppContext, moveTaskOrder models.Move) (*[]services.ShipmentPostalCodeRateArea, error) {
	if moveTaskOrder.AvailableToPrimeAt == nil {
		return nil, apperror.NewUnprocessableEntityError("Move not available to the Prime, unable to retrieve move shipment oconus rateArea")
	}

	contract, err := fetchContract(appCtx, *moveTaskOrder.AvailableToPrimeAt)
	if err != nil {
		return nil, err
	}

	// build set of postalCodes to fetch rateArea for
	var postalCodes = make([]string, 0)
	for _, shipment := range moveTaskOrder.MTOShipments {
		if shipment.PickupAddress != nil {
			if !slices.Contains(postalCodes, shipment.PickupAddress.PostalCode) {
				postalCodes = append(postalCodes, shipment.PickupAddress.PostalCode)
			}
		}
		if shipment.DestinationAddress != nil {
			if !slices.Contains(postalCodes, shipment.DestinationAddress.PostalCode) {
				postalCodes = append(postalCodes, shipment.DestinationAddress.PostalCode)
			}
		}
		if shipment.PPMShipment != nil {
			if shipment.PPMShipment.PickupAddress != nil {
				if !slices.Contains(postalCodes, shipment.PPMShipment.PickupAddress.PostalCode) {
					postalCodes = append(postalCodes, shipment.PPMShipment.PickupAddress.PostalCode)
				}
			}
			if shipment.PPMShipment.DestinationAddress != nil {
				if !slices.Contains(postalCodes, shipment.PPMShipment.DestinationAddress.PostalCode) {
					postalCodes = append(postalCodes, shipment.PPMShipment.DestinationAddress.PostalCode)
				}
			}
		}
	}

	ra, err := fetchRateArea(appCtx, contract.ID, postalCodes)
	if err != nil {
		return nil, err
	}

	return ra, nil
}

func fetchRateArea(appCtx appcontext.AppContext, contractId uuid.UUID, postalCode []string) (*[]services.ShipmentPostalCodeRateArea, error) {
	var rateArea = make([]services.ShipmentPostalCodeRateArea, 0)
	for _, code := range postalCode {
		ra, err := fetchOconusRateAreaByPostalCode(appCtx, contractId, code)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, apperror.NewQueryError("GetRateArea", err, fmt.Sprintf("error retrieving rateArea for contractId:%s, postalCode:%s", contractId, code))
			}
		} else {
			rateArea = append(rateArea, services.ShipmentPostalCodeRateArea{PostalCode: code, RateArea: ra})
		}
	}
	return &rateArea, nil
}

func fetchOconusRateAreaByPostalCode(appCtx appcontext.AppContext, contractId uuid.UUID, postalCode string) (*models.ReRateArea, error) {
	var area models.ReRateArea

	err := appCtx.DB().Q().RawQuery(`select
  		re_rate_areas.*
			from v_locations
  			join re_oconus_rate_areas on re_oconus_rate_areas.us_post_region_cities_id = v_locations.uprc_id
  			join re_rate_areas on re_oconus_rate_areas.rate_area_id = re_rate_areas.id  and v_locations.uspr_zip_id = ?
			and re_rate_areas.contract_id = ?`,
		postalCode, contractId).First(&area)

	appCtx.Logger().Info(fmt.Sprintf("fetchOconusRateAreaByPostalCode sql=%s", fmt.Sprintf(`select
  		re_rate_areas.*
			from v_locations
  			join re_oconus_rate_areas on re_oconus_rate_areas.us_post_region_cities_id = v_locations.uprc_id
  			join re_rate_areas on re_oconus_rate_areas.rate_area_id = re_rate_areas.id  and v_locations.uspr_zip_id = '%s'
			and re_rate_areas.contract_id = '%s'`, postalCode, contractId)))

	if err != nil {
		return nil, err
	}

	return &area, err
}

func fetchContract(appCtx appcontext.AppContext, date time.Time) (*models.ReContract, error) {
	var contractYear models.ReContractYear
	err := appCtx.DB().EagerPreload("Contract").Where("? between start_date and end_date", date).
		First(&contractYear)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewNotFoundError(uuid.Nil, fmt.Sprintf("no contract year found for %s", date.String()))
		}
		return nil, err
	}

	return &contractYear.Contract, nil
}
