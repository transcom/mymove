package migrate

import (
	"github.com/gobuffalo/pop"
)

func (suite *MigrateSuite) TestPrepareCopyFromStdinWithSchema() {
	tableName := "public.transportation_service_provider_performances"
	listOfColumns := "id, performance_period_start, performance_period_end, traffic_distribution_list_id, quality_band, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate"
	columns := parseColumns(listOfColumns)

	suite.DB().Transaction(func(tx *pop.Connection) error {
		stmt, err := prepareCopyFromStdin(tableName, columns, tx)
		suite.Nil(err)
		suite.NotNil(stmt)
		suite.Nil(stmt.Close())
		return nil
	})
}

func (suite *MigrateSuite) TestPrepareCopyFromStdinWithoutSchema() {
	tableName := "transportation_service_provider_performances"
	listOfColumns := "id, performance_period_start, performance_period_end, traffic_distribution_list_id, quality_band, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate"
	columns := parseColumns(listOfColumns)

	suite.DB().Transaction(func(tx *pop.Connection) error {
		stmt, err := prepareCopyFromStdin(tableName, columns, tx)
		suite.Nil(err)
		suite.NotNil(stmt)
		suite.Nil(stmt.Close())
		return nil
	})
}
