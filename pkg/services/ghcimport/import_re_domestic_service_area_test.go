package ghcimport

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) Test_importREDomesticServiceArea() {
	suite.T().Run("import success", func(t *testing.T) {
		gre := &GHCRateEngineImporter{
			Logger:       suite.logger,
			ContractCode: testContractCode,
		}

		err := gre.importREDomesticServiceArea(suite.DB())
		suite.NoError(err)
		suite.helperVerifyServiceAreaCount()
		suite.NotNil(gre.serviceAreaToIDMap)

		// Spot check a service area
		suite.helperCheckServiceAreaValue()
	})

	suite.T().Run("Run a second time with one changed row, should still succeed", func(t *testing.T) {
		gre := &GHCRateEngineImporter{
			Logger:       suite.logger,
			ContractCode: testContractCode,
		}

		// Change a service area and remove a zip and see if they return as they were before.
		var serviceArea models.ReDomesticServiceArea
		err := suite.DB().Where("service_area = '452'").First(&serviceArea)
		suite.NoError(err)
		serviceArea.BasePointCity = "New City"
		suite.MustSave(&serviceArea)
		var zip3 models.ReZip3
		err = suite.DB().Where("zip3 = '647'").First(&zip3)
		suite.NoError(err)
		suite.MustDestroy(&zip3)

		err = gre.importREDomesticServiceArea(suite.DB())
		suite.NoError(err)
		suite.helperVerifyServiceAreaCount()
		suite.NotNil(gre.serviceAreaToIDMap)

		// Check to see if data changed above has been reverted.
		suite.helperCheckServiceAreaValue()
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyServiceAreaCount() {
	// Domestic service areas count
	count, err := suite.DB().Count(&models.ReDomesticServiceArea{})
	suite.NoError(err)
	suite.Equal(4, count)

	// Zip3s count
	count, err = suite.DB().Count(&models.ReZip3s{})
	suite.NoError(err)
	suite.Equal(18, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckServiceAreaValue() {
	var serviceArea models.ReDomesticServiceArea
	err := suite.DB().Where("service_area = '452'").First(&serviceArea)
	suite.NoError(err)
	suite.Equal("Butler", serviceArea.BasePointCity)
	suite.Equal("MO", serviceArea.State)
	suite.Equal(1, serviceArea.ServicesSchedule)
	suite.Equal(3, serviceArea.SITPDSchedule)

	expectedZip3s := []string{"647", "648", "656", "657", "658"}

	var zip3s models.ReZip3s
	err = suite.DB().Where("domestic_service_area_id = ?", serviceArea.ID).All(&zip3s)
	suite.NoError(err)
	for _, zip3 := range zip3s {
		found := false
		for _, expectedZip3 := range expectedZip3s {
			if zip3.Zip3 == expectedZip3 {
				found = true
			}
		}
		suite.True(found, "Could not find zip3 [%s]", zip3.Zip3)
	}

	suite.Len(zip3s, len(expectedZip3s), "Zip3 lengths did not match")
}
