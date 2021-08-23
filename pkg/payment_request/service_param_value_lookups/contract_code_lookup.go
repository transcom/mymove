package serviceparamvaluelookups

import "github.com/transcom/mymove/pkg/appcontext"

// ContractCodeLookup looks up the appropriate contract code
type ContractCodeLookup struct {
}

func (c ContractCodeLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	return keyData.ContractCode, nil
}
