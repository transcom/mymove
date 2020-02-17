package ghcimport

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

const testContractCode = "TEST"
const testContractCode2 = "TEST2"
const testContractName = "Test Contract"

var tablesToTruncate = [...]string{
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
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *GHCRateEngineImportSuite) SetupTest() {
	// Clean up only the rate engine tables we're going to be inserting into for the tests.
	for _, table := range tablesToTruncate {
		sql := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		err := suite.DB().RawQuery(sql).Exec()
		suite.NoError(err)
	}
	// setup re_services which is normally a migration in other environments
	suite.helperSetupReServicesTable()
}

func (suite *GHCRateEngineImportSuite) SetupSuite() {
	suite.helperSetupStagingTables()
}

func (suite *GHCRateEngineImportSuite) TearDownSuite() {
	suite.PopTestSuite.TearDown()
}

func (suite *GHCRateEngineImportSuite) helperSetupStagingTables() {
	fmt.Print("Importing stage data...")
	path := filepath.Join("fixtures", "stage_ghc_pricing.sql")
	c, ioErr := ioutil.ReadFile(path)
	suite.NoError(ioErr)

	sql := string(c)
	err := suite.DB().RawQuery(sql).Exec()
	if suite.NoError(err) {
		fmt.Println(" success")
	}
}

func (suite *GHCRateEngineImportSuite) helperSetupReServicesTable() {
	fmt.Print("Importing re_services data...")
	path := filepath.Join("../../../migrations/app/schema", "20191101201107_create-re-services-table-with-values.up.sql")
	c, ioErr := ioutil.ReadFile(path)
	suite.NoError(ioErr)

	sql := string(c)
	err := suite.DB().RawQuery(sql).Exec()
	if suite.NoError(err) {
		// read second migration
		path = filepath.Join("../../../migrations/app/schema", "20191126160639_update-and-add-values-for-reservices-table.up.sql")
		c, ioErr = ioutil.ReadFile(path)
		suite.NoError(ioErr)

		sql := string(c)
		err := suite.DB().RawQuery(sql).Exec()
		if suite.NoError(err) {
			fmt.Println(" success")
		}
	}
}

func TestGHCRateEngineImportSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &GHCRateEngineImportSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
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
				Logger:       suite.logger,
				ContractCode: testContractCode,
				ContractName: testContractName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if err := tt.gre.Import(suite.DB()); (err != nil) != tt.wantErr {
				t.Errorf("GHCRateEngineImporter.Import() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
