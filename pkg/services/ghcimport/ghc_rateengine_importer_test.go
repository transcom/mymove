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

type GHCRateEngineImportSuite struct {
	*testingsuite.PopTestSuite
}

func (suite *GHCRateEngineImportSuite) SetupSuite() {
	suite.PreloadData(func() {
		suite.helperSetupStagingTables()
		suite.helperSetupReServicesTable()
	})
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
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
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
		suite.Run(tt.name, func() {
			if err := tt.gre.Import(suite.AppContextForTest()); (err != nil) != tt.wantErr {
				suite.T().Errorf("GHCRateEngineImporter.Import() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
