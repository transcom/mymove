package migrate

import (
	"github.com/gobuffalo/pop/v6"
)

func (suite *MigrateSuite) TestPrepareCopyFromStdinWithSchema() {
	tableName := "public.transportation_service_provider_performances"
	listOfColumns := "id, performance_period_start, performance_period_end, traffic_distribution_list_id, quality_band, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate"
	columns := parseColumns(listOfColumns)

	err := suite.DB().Transaction(func(tx *pop.Connection) error {
		stmt, err := prepareCopyFromStdin(tableName, columns, tx)
		suite.Nil(err)
		suite.NotNil(stmt)
		suite.Nil(stmt.Close())
		return nil
	})
	suite.NoError(err)
}

func (suite *MigrateSuite) TestPrepareCopyFromStdinWithBadSchema() {
	tableName := "public.bad.transportation_service_provider_performances"
	listOfColumns := "id, performance_period_start, performance_period_end, traffic_distribution_list_id, quality_band, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate"
	columns := parseColumns(listOfColumns)

	err := suite.DB().Transaction(func(tx *pop.Connection) error {
		stmt, err := prepareCopyFromStdin(tableName, columns, tx)
		suite.NotNil(err)
		suite.Nil(stmt)
		return err
	})
	suite.NotNil(err)
	// TODO: Fix this DB error string literal comparison when we move the COPY-related functionality to jackc/pgx.
	suite.Equal("error preparing copy from stdin statement: pq: relation \"public.bad.transportation_service_provider_performances\" does not exist", err.Error())
}

func (suite *MigrateSuite) TestPrepareCopyFromStdinWithoutSchema() {
	tableName := "transportation_service_provider_performances"
	listOfColumns := "id, performance_period_start, performance_period_end, traffic_distribution_list_id, quality_band, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate"
	columns := parseColumns(listOfColumns)

	err := suite.DB().Transaction(func(tx *pop.Connection) error {
		stmt, err := prepareCopyFromStdin(tableName, columns, tx)
		suite.Nil(err)
		suite.NotNil(stmt)
		suite.Nil(stmt.Close())
		return nil
	})
	suite.NoError(err)
}

func (suite *MigrateSuite) TestPrepareCopyFromStdinWithoutSchemaFail() {
	tableName := "bad_transportation_service_provider_performances"
	listOfColumns := "id, performance_period_start, performance_period_end, traffic_distribution_list_id, quality_band, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate"
	columns := parseColumns(listOfColumns)

	err := suite.DB().Transaction(func(tx *pop.Connection) error {
		stmt, err := prepareCopyFromStdin(tableName, columns, tx)
		suite.NotNil(err)
		suite.Nil(stmt)
		return err
	})
	suite.NotNil(err)
	// TODO: Fix this DB error string literal comparison when we move the COPY-related functionality to jackc/pgx.
	suite.Equal("error preparing copy from stdin statement: pq: relation \"bad_transportation_service_provider_performances\" does not exist", err.Error())
}
