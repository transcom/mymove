package ghcimport

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

const testContractCode = "TEST"
const testContractName = "Test Contract"

type GHCRateEngineImportSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *GHCRateEngineImportSuite) SetupTest() {
	suite.DB().TruncateAll()
	suite.helperSetupStagingTables()
}

func (suite *GHCRateEngineImportSuite) TearDownSuite() {
	suite.PopTestSuite.TearDown()
}

func (suite *GHCRateEngineImportSuite) helperSetupStagingTables() {
	path := filepath.Join("fixtures", "stage_ghc_pricing.sql")
	c, ioErr := ioutil.ReadFile(path)
	suite.NoError(ioErr)

	sql := string(c)
	err := suite.DB().RawQuery(sql).Exec()
	suite.NoError(err)
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
