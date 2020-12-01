package serviceparamvaluelookups

// ContractCodeLookup looks up the appropriate contract code
type ContractCodeLookup struct {
}

func (c ContractCodeLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	return keyData.ContractCode, nil
}
