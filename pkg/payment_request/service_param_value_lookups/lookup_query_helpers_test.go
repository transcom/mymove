package serviceparamvaluelookups

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestLookupQueryHelpers() {
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

}
