package adminapi

import (
	"encoding/json"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// generateQueryFilters is helper to convert filter params from a json string into a slice of services.QueryFilter
//
// The convertFns is a mapping from the filter parameter name(s) to the function that will operate on the
// associated non-empty filter parameter value
func generateQueryFilters(log handlers.Logger, params *string, convertFns map[string]func(string) []services.QueryFilter) (results []services.QueryFilter) {
	if params == nil {
		return results
	}

	input := map[string]string{}
	if err := json.Unmarshal([]byte(*params), &input); err != nil {
		log.Warn(
			"unable to decode param",
			zap.Error(err),
			zap.String("filters", *params),
		)
		return results
	}

	for key, fn := range convertFns {
		content := input[key]
		if content != "" {
			results = append(results, fn(content)...)
		}
	}
	return results
}
