package serviceparamvaluelookups

import "github.com/transcom/mymove/pkg/appconfig"

// ContractCodeLookup looks up the appropriate contract code
type ContractCodeLookup struct {
}

func (c ContractCodeLookup) lookup(appCfg appconfig.AppConfig, keyData *ServiceItemParamKeyData) (string, error) {
	return keyData.ContractCode, nil
}
