package migrate

import (
	"fmt"
	"log"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MigrateSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *MigrateSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestMigrateSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	ms := &MigrateSuite{
		PopTestSuite: testingsuite.NewPopTestSuite("migrate"),
		logger:       logger,
	}
	suite.Run(t, ms)
}

func TestCopyStdinPattern(t *testing.T) {
	tableName := "public.transportation_service_provider_performances"
	listOfColumns := "id, performance_period_start, performance_period_end, traffic_distribution_list_id, quality_band, offer_count, best_value_score, transportation_service_provider_id, created_at, updated_at, rate_cycle_start, rate_cycle_end, linehaul_rate, sit_rate"
	// #nosec This is sql formatting only for testing
	stmtString := fmt.Sprintf("COPY %s (%s) FROM stdin;", tableName, listOfColumns)
	match := copyStdinPattern.FindStringSubmatch(stmtString)

	assert.NotNil(t, match, "Match not found")
	// 0 : Full Line
	assert.Equal(t, match[0], stmtString)
	// 1 : Whitespace Prefix
	assert.Equal(t, match[1], "")
	// 2 : COPY
	assert.Equal(t, match[2], "COPY")
	// 3 : Whitespace
	assert.Equal(t, match[3], " ")
	// 4 : table name
	assert.Equal(t, match[4], tableName)
	// 5 : whitespace
	assert.Equal(t, match[5], " ")
	// 6 : list of columns
	assert.Equal(t, match[6], listOfColumns)
	// 7 : whitespace
	assert.Equal(t, match[7], " ")
	// 8 : FROM
	assert.Equal(t, match[8], "FROM")
	// 9 : whitespace
	assert.Equal(t, match[9], " ")
	// 10 : stdin
	assert.Equal(t, match[10], "stdin")
	// 11 : whitespace
	assert.Equal(t, match[11], "")
	// 12 : ;
	assert.Equal(t, match[12], ";")
	// 12 : whitespace
	assert.Equal(t, match[13], "")

	// preparedStmt := pq.CopyInSchema(parts[0], parts[1], columns...)
	// assert.Equal(t, preparedStmt, stmtString)
}
