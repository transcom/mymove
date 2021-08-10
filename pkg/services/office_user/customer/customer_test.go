package customer

import (
	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *CustomerServiceSuite) TestCustomerFetcher() {
	customer := testdatagen.MakeDefaultServiceMember(suite.DB())
	mtoFetcher := NewCustomerFetcher()

	appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
	actualCustomer, err := mtoFetcher.FetchCustomer(appCfg, customer.ID)
	suite.NoError(err)

	suite.Equal(customer.ID, actualCustomer.ID)
	suite.Equal(*customer.Edipi, *actualCustomer.Edipi)
	suite.Equal(customer.UserID, actualCustomer.UserID)
	suite.NotNil(customer.User)
}
