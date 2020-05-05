package serviceparamvaluelookups

// NotImplementedLookup is the weight estimated lookup
type NotImplementedLookup struct {
}

func (r NotImplementedLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	var value string

	value = "NOT IMPLEMENTED"

	return value, nil
}
