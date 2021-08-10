package serviceparamvaluelookups

import (
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestLookupQueryHelpers() {
	domesticServiceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
		ReDomesticServiceArea: models.ReDomesticServiceArea{
			ServiceArea:      "004",
			ServicesSchedule: 2,
		},
	})

	zip3 := testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
		ReZip3: models.ReZip3{
			Contract:            domesticServiceArea.Contract,
			DomesticServiceArea: domesticServiceArea,
			Zip3:                "350",
		},
	})
	appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
	dsa, err := fetchDomesticServiceArea(appCfg, zip3.Contract.Code, zip3.Zip3)
	suite.FatalNoError(err)
	suite.Equal(strfmt.UUID(domesticServiceArea.ID.String()), strfmt.UUID(dsa.ID.String()))

}
