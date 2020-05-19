package serviceparamvaluelookups

// NotImplementedLookup is the default for unimplemented service item param keys in look-ups
type NotImplementedLookup struct {
}

func (r NotImplementedLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	var value string

	value = "NOT IMPLEMENTED"

	return value, nil
}
