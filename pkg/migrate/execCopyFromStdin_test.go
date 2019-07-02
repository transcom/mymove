package migrate

import (
	"os"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MigrateSuite) TestExecCopyFromStdin() {

	// Create common TSP
	testdatagen.MakeTSP(suite.DB(), testdatagen.Assertions{
		TransportationServiceProvider: models.TransportationServiceProvider{
			ID: uuid.Must(uuid.FromString("231a7b21-346c-4e94-b6bc-672413733f77")),
		},
	})

	// Create TDLs
	testdatagen.MakeTDL(suite.DB(), testdatagen.Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			ID: uuid.Must(uuid.FromString("27f1fbeb-090c-4a91-955c-67899de4d6d6")),
		},
	})

	// Prep the stdin transaction
	tableName := "public.transportation_service_provider_performances"
	listOfColumns := "id, performance_period_start, performance_period_end, traffic_distribution_list_id, quality_band, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate"
	columns := parseColumns(listOfColumns)

	// Load the fixture with the sql example
	f, err := os.Open("./fixtures/copyFromStdin.sql")
	suite.Nil(err)

	in := NewBuffer()
	go ReadInSQL(f, in, true, true, true)

	suite.DB().Transaction(func(tx *pop.Connection) error {
		i := 572
		wait := 30 * time.Second
		lineNum, err := execCopyFromStdin(in, i, tableName, columns, tx, wait)
		suite.Nil(err)
		suite.NotNil(lineNum)
		suite.Equal(lineNum, 840)
		return nil
	})
}
