package ghcimport

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

const testContractCode = "TEST"
const testContractCode2 = "TEST2"
const testContractName = "Test Contract"

var testContractStartDate = time.Date(2021, time.February, 01, 0, 0, 0, 0, time.UTC)

var tablesToTruncate = []string{
	"re_contract_years",
	"re_contracts",
	"re_domestic_accessorial_prices",
	"re_domestic_linehaul_prices",
	"re_domestic_other_prices",
	"re_domestic_service_area_prices",
	"re_domestic_service_areas",
	"re_intl_accessorial_prices",
	"re_intl_other_prices",
	"re_intl_prices",
	"re_rate_areas",
	"re_services",
	"re_shipment_type_prices",
	"re_task_order_fees",
	"re_zip3s",
}

type GHCRateEngineImportSuite struct {
	*testingsuite.PopTestSuite
}

func (suite *GHCRateEngineImportSuite) SetupTest() {
	// Clean up only the rate engine tables we're going to be inserting into for the tests.
	err := suite.Truncate(tablesToTruncate)
	suite.NoError(err)

	// setup re_services which is normally a migration in other environments
	suite.helperSetupReServicesTable()
}

func (suite *GHCRateEngineImportSuite) SetupSuite() {
	suite.helperSetupStagingTables()
}

func (suite *GHCRateEngineImportSuite) TearDownSuite() {
	suite.PopTestSuite.TearDown()
}

func (suite *GHCRateEngineImportSuite) helperLoadSQLFixture(fileName string) {
	path := filepath.Join("fixtures", fileName)
	_, err := os.Stat(path)
	suite.NoError(err)

	c, ioErr := os.ReadFile(filepath.Clean(path))
	suite.NoError(ioErr)

	sql := string(c)
	err = suite.DB().RawQuery(sql).Exec()
	suite.NoError(err)
}

func (suite *GHCRateEngineImportSuite) helperSetupStagingTables() {
	suite.helperLoadSQLFixture("stage_ghc_pricing.sql")
}

func (suite *GHCRateEngineImportSuite) helperSetupReServicesTable() {
	suite.helperLoadSQLFixture("re_services_data.sql")
}

func TestGHCRateEngineImportSuite(t *testing.T) {
	hs := &GHCRateEngineImportSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}

	suite.Run(t, hs)
}

func (suite *GHCRateEngineImportSuite) TestGHCRateEngineImporter_Import() {
	tests := []struct {
		name    string
		gre     *GHCRateEngineImporter
		wantErr bool
	}{
		{
			name: "Run GHC Rate Engine Importer",
			gre: &GHCRateEngineImporter{
				ContractCode:      testContractCode,
				ContractName:      testContractName,
				ContractStartDate: testContractStartDate,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if err := tt.gre.Import(suite.AppContextForTest()); (err != nil) != tt.wantErr {
				t.Errorf("GHCRateEngineImporter.Import() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
