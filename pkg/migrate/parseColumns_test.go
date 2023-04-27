package migrate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseColumns(t *testing.T) {
	listOfColumns := "id, performance_period_start, performance_period_end, traffic_distribution_list_id, quality_band, offer_count, best_value_score, some_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate"
	columns := parseColumns(listOfColumns)
	assert.Equal(t, len(columns), 14)
	assert.Equal(t, columns[0], "id")
	assert.Equal(t, columns[13], "sit_rate")
}
