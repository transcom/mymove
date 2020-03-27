package ghcimport

import (
	"github.com/transcom/mymove/pkg/testdatagen"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) Test_mapZipsToRateAreas() {
	gre := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode,
	}

	suite.T().Run("Successfully map rezip3s and zip5s to rate areas", func(t *testing.T) {
		rezip3s := []models.ReZip3 {
			models.ReZip3{Zip3: "735"},
			models.ReZip3{Zip3: "820"},
			models.ReZip3{Zip3: "833"},
			models.ReZip3{Zip3: "850"},
			models.ReZip3{Zip3: "923"},
		}

		for _, zip3 := range rezip3s {
			err := suite.DB().Save(&zip3)
			if err != nil {
				suite.Error(err)
			}
		}

		// call our function
		// assert that zip3s were associated with the correct rate areas
		// assert that zip5s were created with the correct rate areas

		//err := gre.loadServiceMap(suite.DB())
		//suite.NoError(err)
		//
		//suite.NotNil(gre.serviceToIDMap)
		//
		//count, err := suite.DB().Count(&models.ReService{})
		//suite.NoError(err)
		//suite.Greater(count, 0)
		//suite.Equal(count, len(gre.serviceToIDMap))
		//
		//// Spot-check a service code
		//testServiceCode := "DOASIT"
		//if suite.Contains(gre.serviceToIDMap, testServiceCode) {
		//	suite.NotEqual(uuid.Nil, gre.serviceToIDMap[testServiceCode])
		//}
	})
}
