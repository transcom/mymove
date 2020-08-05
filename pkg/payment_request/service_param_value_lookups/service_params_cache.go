package serviceparamvaluelookups

import (
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/models"

	"github.com/gofrs/uuid"

)

/*
ServiceParamsCache can be used to cache service item param values per a given Payment Request creation.

Service Item Param Lookup functions can cache values and check if the value exist before having to potentially repeat
a db query or calculation to determine the value for a given MTOShipment/ParamKey value.

This struct also allows for a quick way to determine given a service item code, which param keys are needed.
 */


/*
Param Key Value Cache Map
maps shipment ID to param cache values, guarding against
multiple shipments having varying param key values
 */
type ParamKeyValue struct {
	value *string
}
type ParamKeyValueCacheMap map[string]ParamKeyValue
type MTOShipmentParamKeyMap map[uuid.UUID]ParamKeyValueCacheMap

/*
Service Item Param Key maps,
tells us if a service item code needs a param key based on the
param key string
 */
type keyExistMap map[models.ServiceItemParamName]bool
type NeedsParamKeyMap map[models.ReServiceCode]keyExistMap

// ServiceParamsCache contains service item parameter keys
type ServiceParamsCache struct {
	db               *pop.Connection
	paramsCache      MTOShipmentParamKeyMap
	needsParamKey    NeedsParamKeyMap
}

func (spc *ServiceParamsCache) addParamValue (mtoShipmentID uuid.UUID, paramKey string, value string) {
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

func (spc *ServiceParamsCache) ParamValue (mtoShipmentID uuid.UUID, paramKey string) *string {
	if paramValueCacheMap, ok := spc.paramsCache[mtoShipmentID]; ok {
		if keyValue, paramOK := paramValueCacheMap[paramKey]; paramOK {
			return keyValue.value
		}
	}

	return nil
}


func (spc *ServiceParamsCache) setNeedsParamKeyMap (code models.ReServiceCode) error {

	// build up service item code paramkey map if it doesn't yet exist
	if _, ok := spc.needsParamKey[code]; !ok {
		var paramKeys []string
		query := `
        SELECT sipk.key FROM re_services
		LEFT JOIN service_params sp on re_services.id = sp.service_id
		LEFT JOIN service_item_param_keys sipk on sp.service_item_param_key_id = sipk.id
        WHERE
            re_services.code = $1
    `
		err := spc.db.RawQuery(query, code).All(paramKeys)
		if err != nil {
			return err
		} else {
			var tmpMap keyExistMap
			for _, pk := range paramKeys {
				tmpMap[models.ServiceItemParamName(pk)] = true
			}
			spc.needsParamKey[code] = tmpMap
		}
		return  nil
	}
	return nil
}

func (spc *ServiceParamsCache) paramKeyExist (paramKey models.ServiceItemParamName, keyMap keyExistMap) bool {
	if val, paramKeyOK := keyMap[paramKey]; paramKeyOK {
		return val
	} else {
		return false
	}
}

func (spc *ServiceParamsCache) ServiceItemNeedsParamKey (code models.ReServiceCode, paramKey models.ServiceItemParamName) (bool, error) {
	var err error
	if keyMap, codeOK := spc.needsParamKey[code]; codeOK {
		return spc.paramKeyExist(paramKey, keyMap), nil
	} else {
		err = spc.setNeedsParamKeyMap(code)
		if err == nil {
			return spc.paramKeyExist(paramKey, keyMap), nil
		} else {
			return false, fmt.Errorf("ServiceParamsCache.needsParamKey failed to retrieve NeedsParamKeyMap for service item code %s with error: %w", code, err)
		}
	}
}