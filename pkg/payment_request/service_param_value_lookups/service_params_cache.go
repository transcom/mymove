package serviceparamvaluelookups

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

/*
ServiceParamsCache can be used to cache service item param values per a given Payment Request creation.

Service Item Param Lookup functions can cache values and check if the value exist before having to potentially repeat
a db query or calculation to determine the value for a given MTOShipment/ParamKey value.

This struct also allows for a quick way to determine given a service item code, which param keys are needed.
*/

/******
Param Key Value Cache Map
maps shipment ID to param cache values, guarding against
multiple shipments having varying param key values
******/

// ParamKeyValue param cache value
type ParamKeyValue struct {
	value *string
}

// ParamKeyValueCacheMap maps service param key string to the param cached value
type ParamKeyValueCacheMap map[models.ServiceItemParamName]ParamKeyValue

//MTOShipmentParamKeyMap maps an MTOShipmentID to a map of param caches for that shipment
type MTOShipmentParamKeyMap map[uuid.UUID]ParamKeyValueCacheMap

/******
Service Item Param Key maps,
tells us if a service item code needs a param key based on the
param key string
******/

// keyExistMap if the param key name is present, return true
type keyExistMap map[models.ServiceItemParamName]bool

// NeedsParamKeyMap param key maps, maps an ReService.Code to a list of param keys used for pricing
type NeedsParamKeyMap map[models.ReServiceCode]keyExistMap

// ServiceParamsCache contains service item parameter keys
type ServiceParamsCache struct {
	paramsCache   MTOShipmentParamKeyMap
	needsParamKey NeedsParamKeyMap
}

// NewServiceParamsCache creates a ServiceParamCache with initialized fields
func NewServiceParamsCache() ServiceParamsCache {
	return ServiceParamsCache{
		needsParamKey: NeedsParamKeyMap{},
		paramsCache:   MTOShipmentParamKeyMap{},
	}
}

func (spc *ServiceParamsCache) addParamValue(mtoShipmentID uuid.UUID, paramKey models.ServiceItemParamName, value string) {
	if paramValueCacheMap, ok := spc.paramsCache[mtoShipmentID]; ok {
		paramValueCacheMap[paramKey] = ParamKeyValue{
			value: &value,
		}
	} else {
		spc.paramsCache[mtoShipmentID] = ParamKeyValueCacheMap{}
		spc.paramsCache[mtoShipmentID][paramKey] = ParamKeyValue{
			value: &value,
		}
	}
}

// ParamValue returns the caches param value for the given MTOShipmentID
func (spc *ServiceParamsCache) ParamValue(mtoShipmentID uuid.UUID, paramKey models.ServiceItemParamName) *string {
	if paramValueCacheMap, ok := spc.paramsCache[mtoShipmentID]; ok {
		if keyValue, paramOK := paramValueCacheMap[paramKey]; paramOK {
			return keyValue.value
		}
	}

	return nil
}

func (spc *ServiceParamsCache) setNeedsParamKeyMap(appCtx appcontext.AppContext, code models.ReServiceCode) error {
	// build up service item code paramkey map if it doesn't yet exist
	if _, ok := spc.needsParamKey[code]; !ok {
		type ParamKeys []string
		paramKeys := ParamKeys{}
		query := `
		SELECT key FROM service_item_param_keys
		LEFT JOIN service_params sp on service_item_param_keys.id = sp.service_item_param_key_id
		LEFT JOIN re_services rs on sp.service_id = rs.id
		WHERE rs.code = $1
    `
		codeStr := string(code)
		err := appCtx.DB().RawQuery(query, codeStr).All(&paramKeys)
		if err != nil {
			return err
		}

		tmpMap := keyExistMap{}
		for _, pk := range paramKeys {
			tmpMap[models.ServiceItemParamName(pk)] = true
		}
		spc.needsParamKey[code] = tmpMap
	}
	return nil
}

func (spc *ServiceParamsCache) paramKeyExist(paramKey models.ServiceItemParamName, keyMap keyExistMap) bool {
	if val, paramKeyOK := keyMap[paramKey]; paramKeyOK {
		return val
	}

	return false
}

// ServiceItemNeedsParamKey returns true/false if the ReServiceCode uses the particular ServiceItemParamKey
// for calculating the service item  price
func (spc *ServiceParamsCache) ServiceItemNeedsParamKey(appCtx appcontext.AppContext, code models.ReServiceCode, paramKey models.ServiceItemParamName) (bool, error) {
	var err error

	if keyMap, codeOK := spc.needsParamKey[code]; codeOK {
		return spc.paramKeyExist(paramKey, keyMap), nil
	}

	err = spc.setNeedsParamKeyMap(appCtx, code)
	if err == nil {
		if keyMap, codeOK := spc.needsParamKey[code]; codeOK {
			return spc.paramKeyExist(paramKey, keyMap), nil
		}
	}

	return false, fmt.Errorf("ServiceParamsCache.needsParamKey failed to retrieve NeedsParamKeyMap for service item code %s with error: %w", code, err)

}
