package serviceparamvaluelookups

import (
	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
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

func (suite *ServiceParamValueLookupsSuite) TestFetchRateArea() {
	suite.Run("Successful", func() {
		service := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIOPSIT)
		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		address := factory.BuildAddress(suite.DB(), nil, nil)
		ra, err := fetchRateArea(suite.AppContextForTest(), service.ID, address.ID, contract.ID)
		suite.FatalNoError(err)
		suite.True(len(ra.Code) > 0)
	})

	suite.Run("failure", func() {
		service := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIOPSIT)
		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		invalidAddressID := uuid.Must(uuid.NewV4())
		_, err := fetchRateArea(suite.AppContextForTest(), service.ID, invalidAddressID, contract.ID)
		suite.NotNil(err)
		suite.Contains(err.Error(), "error fetching rate area id for shipment")
	})
}
