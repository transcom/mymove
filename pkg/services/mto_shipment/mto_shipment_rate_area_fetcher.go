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

func (f mtoShipmentRateAreaFetcher) GetPrimeMoveShipmentRateAreas(appCtx appcontext.AppContext, moveTaskOrder models.Move) (*[]services.ShipmentPostalCodeRateArea, error) {
	if moveTaskOrder.AvailableToPrimeAt == nil {
		return nil, apperror.NewUnprocessableEntityError("Move not available to the Prime, unable to retrieve move shipment oconus rateArea")
	}

	contract, err := fetchContract(appCtx, *moveTaskOrder.AvailableToPrimeAt)
	if err != nil {
		return nil, err
	}

	// build set of postalCodes to fetch rateArea for
	var oconusPostalCodes = make([]string, 0)
	var conusPostalCodes = make([]string, 0)
	for _, shipment := range moveTaskOrder.MTOShipments {
		// B-22767: We want both domestic and international rate area info but only for international shipments
		if shipment.MarketCode == models.MarketCodeInternational {
			if shipment.PickupAddress != nil {
				if !slices.Contains(oconusPostalCodes, shipment.PickupAddress.PostalCode) &&
					shipment.PickupAddress.IsOconus != nil && *shipment.PickupAddress.IsOconus {
					oconusPostalCodes = append(oconusPostalCodes, shipment.PickupAddress.PostalCode)
				} else if !slices.Contains(conusPostalCodes, shipment.PickupAddress.PostalCode) {
					conusPostalCodes = append(conusPostalCodes, shipment.PickupAddress.PostalCode)
				}
			}
			if shipment.DestinationAddress != nil {
				if !slices.Contains(oconusPostalCodes, shipment.DestinationAddress.PostalCode) &&
					shipment.DestinationAddress.IsOconus != nil && *shipment.DestinationAddress.IsOconus {
					oconusPostalCodes = append(oconusPostalCodes, shipment.DestinationAddress.PostalCode)
				} else if !slices.Contains(conusPostalCodes, shipment.DestinationAddress.PostalCode) {
					conusPostalCodes = append(conusPostalCodes, shipment.DestinationAddress.PostalCode)
				}
			}
			if shipment.PPMShipment != nil {
				if shipment.PPMShipment.PickupAddress != nil {
					if !slices.Contains(oconusPostalCodes, shipment.PPMShipment.PickupAddress.PostalCode) &&
						shipment.PPMShipment.PickupAddress.IsOconus != nil && *shipment.PPMShipment.PickupAddress.IsOconus {
						oconusPostalCodes = append(oconusPostalCodes, shipment.PPMShipment.PickupAddress.PostalCode)
					} else if !slices.Contains(conusPostalCodes, shipment.PPMShipment.PickupAddress.PostalCode) {
						conusPostalCodes = append(conusPostalCodes, shipment.PPMShipment.PickupAddress.PostalCode)
					}
				}
				if shipment.PPMShipment.DestinationAddress != nil {
					if !slices.Contains(oconusPostalCodes, shipment.PPMShipment.DestinationAddress.PostalCode) &&
						shipment.PPMShipment.DestinationAddress.IsOconus != nil && *shipment.PPMShipment.DestinationAddress.IsOconus {
						oconusPostalCodes = append(oconusPostalCodes, shipment.PPMShipment.DestinationAddress.PostalCode)
					} else if !slices.Contains(conusPostalCodes, shipment.PPMShipment.DestinationAddress.PostalCode) {
						conusPostalCodes = append(conusPostalCodes, shipment.PPMShipment.DestinationAddress.PostalCode)
					}
				}
			}
		}
	}

	ora, err := fetchOconusRateAreas(appCtx, contract.ID, oconusPostalCodes)
	if err != nil {
		return nil, err
	}

	cra, err := fetchConusRateAreas(appCtx, contract.ID, conusPostalCodes)
	if err != nil {
		return nil, err
	}

	ra := append(*ora, *cra...)
	return &ra, nil
}

func fetchOconusRateAreas(appCtx appcontext.AppContext, contractId uuid.UUID, postalCodes []string) (*[]services.ShipmentPostalCodeRateArea, error) {
	var rateAreasMap = make([]services.ShipmentPostalCodeRateArea, 0)
	for _, postalCode := range postalCodes {
		ra, err := fetchOconusRateAreaByPostalCode(appCtx, contractId, postalCode)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, apperror.NewQueryError("GetRateArea", err, fmt.Sprintf("error retrieving rateArea for contractId:%s, postalCode:%s", contractId, postalCode))
			}
		} else {
			rateAreasMap = append(rateAreasMap, services.ShipmentPostalCodeRateArea{PostalCode: postalCode, RateArea: ra})
		}
	}
	return &rateAreasMap, nil
}

func fetchConusRateAreas(appCtx appcontext.AppContext, contractId uuid.UUID, postalCodes []string) (*[]services.ShipmentPostalCodeRateArea, error) {
	var rateAreasMap = make([]services.ShipmentPostalCodeRateArea, 0)
	for _, postalCode := range postalCodes {
		ra, err := models.FetchConusRateAreaByPostalCode(appCtx.DB(), postalCode, contractId)
		if err != nil {
			if err != sql.ErrNoRows {
				return nil, apperror.NewQueryError("GetRateArea", err, fmt.Sprintf("error retrieving rateArea for contractId:%s, postalCode:%s", contractId, postalCode))
			}
		} else {
			rateAreasMap = append(rateAreasMap, services.ShipmentPostalCodeRateArea{PostalCode: postalCode, RateArea: ra})
		}
	}
	return &rateAreasMap, nil
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
