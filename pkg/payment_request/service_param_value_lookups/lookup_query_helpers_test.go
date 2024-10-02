package serviceparamvaluelookups

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestLookupQueryHelpers() {

	setupInternationalRateAreaTest := func() (models.ReRateArea, models.ReZip5RateArea) {
		internationalRateArea := testdatagen.FetchOrMakeReRateArea(suite.DB(), testdatagen.Assertions{
			ReRateArea: models.ReRateArea{
				Name: "InternationalRateArea",
			},
			ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
		})

		zip5 := testdatagen.MakeReZip5RateArea(suite.DB(), testdatagen.Assertions{
			ReZip5RateArea: models.ReZip5RateArea{
				Contract:   internationalRateArea.Contract,
				ContractID: internationalRateArea.ContractID,
				RateArea:   internationalRateArea,
				Zip5:       "35035",
			},
		})
		return internationalRateArea, zip5
	}

	suite.Run("can lookup domestic service area by zip3", func() {
		domesticServiceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea:      "004",
				ServicesSchedule: 2,
			},
			ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
		})

		zip3 := testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            domesticServiceArea.Contract,
				ContractID:          domesticServiceArea.ContractID,
				DomesticServiceArea: domesticServiceArea,
				Zip3:                "350",
			},
		})
		dsa, err := fetchDomesticServiceArea(suite.AppContextForTest(), zip3.Contract.Code, zip3.Zip3)
		suite.FatalNoError(err)
		suite.Equal(strfmt.UUID(domesticServiceArea.ID.String()), strfmt.UUID(dsa.ID.String()))
	})

	suite.Run("can lookup international rate area by zip5", func() {
		internationalRateArea, zip5 := setupInternationalRateAreaTest()
		rateArea, err := fetchInternationalRateArea(suite.AppContextForTest(), zip5.Contract.Code, zip5.Zip5)
		suite.FatalNoError(err)
		suite.Equal(strfmt.UUID(internationalRateArea.ID.String()), strfmt.UUID(rateArea.ID.String()))
	})

	suite.Run("can fail to lookup international rate area by zip5", func() {
		_, err := fetchInternationalRateArea(suite.AppContextForTest(), "test", "12345")
		suite.Error(err)
	})
}
