package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

// ContractCodeLookup looks up the appropriate contract code
type ContractCodeLookup struct {
}

func (c ContractCodeLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	// For now, just default the contract code until we have a process defined
	return ghcrateengine.DefaultContractCode, nil
}
