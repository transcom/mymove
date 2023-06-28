package serviceparamvaluelookups

import "github.com/transcom/mymove/pkg/appcontext"

// NotImplementedLookup is the default for unimplemented service item param keys in look-ups
type NotImplementedLookup struct {
}

func (r NotImplementedLookup) lookup(_ appcontext.AppContext, _ *ServiceItemParamKeyData) (string, error) {
	value := "NOT IMPLEMENTED"

	return value, nil
}
