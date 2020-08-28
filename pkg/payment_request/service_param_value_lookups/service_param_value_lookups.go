package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"github.com/transcom/mymove/pkg/models"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

// ServiceItemParamKeyData contains service item parameter keys
type ServiceItemParamKeyData struct {
	db               *pop.Connection
	planner          route.Planner
	lookups          map[string]ServiceItemParamKeyLookup
	MTOServiceItemID uuid.UUID
	MTOServiceItem   models.MTOServiceItem
	PaymentRequestID uuid.UUID
	MoveTaskOrderID  uuid.UUID
	ContractCode     string
}

// ServiceItemParamKeyLookup does lookup on service item parameter keys
type ServiceItemParamKeyLookup interface {
	lookup(keyData *ServiceItemParamKeyData) (string, error)
}

// ServiceParamLookupInitialize initializes service parameter lookup
func ServiceParamLookupInitialize(
	db *pop.Connection,
	planner route.Planner,
	mtoServiceItemID uuid.UUID,
	paymentRequestID uuid.UUID,
	moveTaskOrderID uuid.UUID,
) (*ServiceItemParamKeyData, error) {

	// Get the MTOServiceItem
	var mtoServiceItem models.MTOServiceItem
	err := db.Eager("ReService").Find(&mtoServiceItem, mtoServiceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, services.NewNotFoundError(mtoServiceItemID, "looking for MTOServiceItemID")
		default:
			return nil, err
		}
	}

	var mtoShipment models.MTOShipment
	var pickupAddress models.Address
	var destinationAddress models.Address
	switch mtoServiceItem.ReService.Code {
	case models.ReServiceCodeCS, models.ReServiceCodeMS:
		// Do nothing, these service items don't use the MTOShipment
	default:
		// Make sure there's an MTOShipment since that's nullable
		if mtoServiceItem.MTOShipmentID == nil {
			return nil, services.NewNotFoundError(uuid.Nil, "looking for MTOShipmentID")
		}
		err := db.Eager("PickupAddress", "DestinationAddress").Find(&mtoShipment, mtoServiceItem.MTOShipmentID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return nil, services.NewNotFoundError(mtoServiceItemID, "looking for MTOServiceItemID")
			default:
				return nil, err
			}
		}

		if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK {
			if mtoShipment.PickupAddressID == nil {
				return nil, services.NewNotFoundError(uuid.Nil, "looking for PickupAddressID")
			}
			pickupAddress = *mtoShipment.PickupAddress
		}

		if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK {
			if mtoShipment.DestinationAddressID == nil {
				return nil, services.NewNotFoundError(uuid.Nil, "looking for DestinationAddressID")
			}
			destinationAddress = *mtoShipment.DestinationAddress
		}
	}

	s := ServiceItemParamKeyData{
		db:               db,
		planner:          planner,
		lookups:          make(map[string]ServiceItemParamKeyLookup),
		MTOServiceItemID: mtoServiceItemID,
		MTOServiceItem:   mtoServiceItem,
		PaymentRequestID: paymentRequestID,
		MoveTaskOrderID:  moveTaskOrderID,
		/*
			DefaultContractCode = TRUSS_TEST is temporarily being used here because the contract
			code is not currently accessible. This is caused by:
				- mtoServiceItem is not linked or associated with a contract record
				- MTO currently has a contractor_id but not a contract_id
			In order for this lookup's query to have accesss to a contract code there must be a contract_code field created on either the mtoServiceItem or the MTO models
			If it'll will be possible for a MTO to contain service items that are associated with different contracts
			then it would be ideal for the mtoServiceItem records to contain a contract code that can then be passed
			to this query. Otherwise the contract_code field could be added to the MTO.
		*/
		ContractCode: ghcrateengine.DefaultContractCode,
	}

	for _, key := range models.ValidServiceItemParamNames {
		s.lookups[key] = NotImplementedLookup{}
	}

	s.lookups[models.ServiceItemParamNameActualPickupDate.String()] = ActualPickupDateLookup{
		MTOShipment: mtoShipment,
	}
	s.lookups[models.ServiceItemParamNameRequestedPickupDate.String()] = RequestedPickupDateLookup{
		MTOShipment: mtoShipment,
	}

	if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK {
		s.lookups[models.ServiceItemParamNameDistanceZip5.String()] = DistanceZip5Lookup{
			PickupAddress:      pickupAddress,
			DestinationAddress: destinationAddress,
		}
		s.lookups[models.ServiceItemParamNameDistanceZip3.String()] = DistanceZip3Lookup{
			PickupAddress:      pickupAddress,
			DestinationAddress: destinationAddress,
		}
	}

	s.lookups[models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier.String()] = FSCWeightBasedDistanceMultiplierLookup{

		MTOShipment: mtoShipment,
	}
	s.lookups[models.ServiceItemParamNameWeightBilledActual.String()] = WeightBilledActualLookup{
		MTOShipment: mtoShipment,
	}
	s.lookups[models.ServiceItemParamNameWeightEstimated.String()] = WeightEstimatedLookup{
		MTOShipment: mtoShipment,
	}
	s.lookups[models.ServiceItemParamNameWeightActual.String()] = WeightActualLookup{
		MTOShipment: mtoShipment,
	}

	if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK {
		s.lookups[models.ServiceItemParamNameZipPickupAddress.String()] = ZipAddressLookup{
			Address: pickupAddress,
		}
	}

	if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK {
		s.lookups[models.ServiceItemParamNameZipDestAddress.String()] = ZipAddressLookup{
			Address: destinationAddress,
		}
	}

	s.lookups[models.ServiceItemParamNameMTOAvailableToPrimeAt.String()] = MTOAvailableToPrimeAtLookup{}

	if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK {
		s.lookups[models.ServiceItemParamNameServiceAreaOrigin.String()] = ServiceAreaLookup{
			Address: pickupAddress,
		}
	}

	if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK {
		s.lookups[models.ServiceItemParamNameServiceAreaDest.String()] = ServiceAreaLookup{
			Address: destinationAddress,
		}
	}
	s.lookups[models.ServiceItemParamNameContractCode.String()] = ContractCodeLookup{}
	s.lookups[models.ServiceItemParamNamePSILinehaulDom.String()] = PSILinehaulDomLookup{
		MTOShipment: mtoShipment,
	}
	s.lookups[models.ServiceItemParamNamePSILinehaulDomPrice.String()] = PSILinehaulDomPriceLookup{
		MTOShipment: mtoShipment,
	}
	s.lookups[models.ServiceItemParamNameEIAFuelPrice.String()] = EIAFuelPriceLookup{
		MTOShipment: mtoShipment,
	}

	if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK {
		s.lookups[models.ServiceItemParamNameServicesScheduleOrigin.String()] = ServicesScheduleLookup{
			Address: pickupAddress,
		}
	}

	if mtoServiceItem.ReService.Code != models.ReServiceCodeDPK {
		s.lookups[models.ServiceItemParamNameServicesScheduleDest.String()] = ServicesScheduleLookup{
			Address: destinationAddress,
		}
	}

	return &s, nil
}

// ServiceParamValue returns a service parameter value from a key
func (s *ServiceItemParamKeyData) ServiceParamValue(key string) (string, error) {
	if lookup, ok := s.lookups[key]; ok {
		value, err := lookup.lookup(s)
		if err != nil {
			return "", fmt.Errorf(" failed ServiceParamValue %sLookup with error %w", key, err)
		}
		return value, nil
	}
	return "", fmt.Errorf("  ServiceParamValue <%sLookup> does not exist for key: <%s>", key, key)
}
