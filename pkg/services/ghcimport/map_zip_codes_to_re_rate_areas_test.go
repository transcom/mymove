package ghcimport

import (
	"github.com/transcom/mymove/pkg/models"
	"testing"
)

// Fixtures and test data should have the minimal amount of data needed

func (suite *GHCRateEngineImportSuite) Test_mapZipsToRateAreas() {
	gre := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode,
	}
//
	suite.T().Run("map zip3s and zip5s to rate areas", func(t *testing.T) {
		rezip3s := []models.ReZip3 {
			models.ReZip3{Zip3: "735", ContractID: gre.ContractID, BasePointCity: "Snyder", State: "OK", DomesticServiceAreaID: gre.serviceAreaToIDMap["636"]},
			models.ReZip3{Zip3: "820", ContractID: gre.ContractID, BasePointCity: "Laramie", State: "WY", DomesticServiceAreaID: gre.serviceAreaToIDMap["880"]},
			models.ReZip3{Zip3: "833", ContractID: gre.ContractID, BasePointCity: "Shoshone", State: "ID", DomesticServiceAreaID: gre.serviceAreaToIDMap["244"]},
			models.ReZip3{Zip3: "850", ContractID: gre.ContractID, BasePointCity: "Phoenix", State: "AZ", DomesticServiceAreaID: gre.serviceAreaToIDMap["028"]},
			models.ReZip3{Zip3: "923", ContractID: gre.ContractID, BasePointCity: "Barstow", State: "CA", DomesticServiceAreaID: gre.serviceAreaToIDMap["072"]},
		}

		for _, zip3 := range rezip3s {
			err := suite.DB().Save(&zip3)
			if err != nil {
				suite.Error(err)
			}
		}

		err := gre.mapZipsToRateAreas(suite.DB())
		suite.NoError(err)
//
//		// call our function
//		// assert that zip3s were associated with the correct rate areas
//		// assert that zip5s were created with the correct rate areas
	})
}
