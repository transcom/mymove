package ghcimport

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREDomesticServiceAreaPrices() {
	suite.T().Run("successfully run import re_domestic_service_area_prices table empty", func(t *testing.T) {
		suite.helperSetupReServicesTable()
		gre := &GHCRateEngineImporter{
			Logger:       suite.logger,
			ContractCode: testContractCode,
		}
		err := gre.importREContract(suite.DB())
		suite.NoError(err)
		suite.NotNil(gre.contractID)
		err = gre.importREDomesticServiceArea(suite.DB())
		suite.NoError(err)
		err = gre.importREDomesticServiceAreaPrices(suite.DB())
		suite.NoError(err)

		var serviceAreaPrices models.ReDomesticServiceAreaPrices
		err = suite.DB().Eager("DomesticServiceArea", "Service").All(&serviceAreaPrices)
		suite.NoError(err)
		suite.Equal(32, len(serviceAreaPrices))

		suite.Equal("DSH", serviceAreaPrices[0].Service.Code)
		suite.Equal("AL", serviceAreaPrices[0].DomesticServiceArea.State)
		suite.Equal("Birmingham", serviceAreaPrices[0].DomesticServiceArea.BasePointCity)
		suite.Equal(false, serviceAreaPrices[0].IsPeakPeriod)
		suite.Equal(unit.Cents(127), serviceAreaPrices[0].PriceCents)

		suite.Equal("DODP", serviceAreaPrices[1].Service.Code)
		suite.Equal("AL", serviceAreaPrices[1].DomesticServiceArea.State)
		suite.Equal("Birmingham", serviceAreaPrices[1].DomesticServiceArea.BasePointCity)
		suite.Equal(false, serviceAreaPrices[1].IsPeakPeriod)
		suite.Equal(unit.Cents(689), serviceAreaPrices[1].PriceCents)

		suite.Equal("DFSIT", serviceAreaPrices[2].Service.Code)
		suite.Equal("AL", serviceAreaPrices[2].DomesticServiceArea.State)
		suite.Equal("Birmingham", serviceAreaPrices[2].DomesticServiceArea.BasePointCity)
		suite.Equal(false, serviceAreaPrices[2].IsPeakPeriod)
		suite.Equal(unit.Cents(1931), serviceAreaPrices[2].PriceCents)

		suite.Equal("DASIT", serviceAreaPrices[3].Service.Code)
		suite.Equal("AL", serviceAreaPrices[3].DomesticServiceArea.State)
		suite.Equal("Birmingham", serviceAreaPrices[3].DomesticServiceArea.BasePointCity)
		suite.Equal(false, serviceAreaPrices[3].IsPeakPeriod)
		suite.Equal(unit.Cents(68), serviceAreaPrices[3].PriceCents)
	})
}

func (suite *GHCRateEngineImportSuite) Test_importREDomesticServiceAreaPricesFailures() {
	suite.T().Run("stage_domestic_service_area_prices table missing", func(t *testing.T) {
		// drop a staging table that we are depending on to do import
		dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s;", "stage_domestic_service_area_prices")
		dropErr := suite.DB().RawQuery(dropQuery).Exec()
		suite.NoError(dropErr)

		suite.helperSetupReServicesTable()

		gre := &GHCRateEngineImporter{
			Logger:       suite.logger,
			ContractCode: testContractCode,
		}

		err := gre.importREContract(suite.DB())
		suite.NoError(err)
		suite.NotNil(gre.contractID)

		err = gre.importREDomesticServiceAreaPrices(suite.DB())
		if suite.Error(err) {
			suite.Equal("Error looking up StageDomesticServiceAreaPrice data: unable to fetch records: pq: relation \"stage_domestic_service_area_prices\" does not exist", err.Error())
		}
	})
}

func (suite *GHCRateEngineImportSuite) helperSetupReServicesTable() {
	path := filepath.Join("../../../migrations", "20191101201107_create-re-services-table-with-values.up.sql")
	c, ioErr := ioutil.ReadFile(path)
	suite.NoError(ioErr)

	sql := string(c)
	err := suite.DB().RawQuery(sql).Exec()
	suite.NoError(err)
}
