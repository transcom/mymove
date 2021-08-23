package customer

import (
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *CustomerServiceSuite) TestCustomerFetcher() {
	customer := testdatagen.MakeDefaultServiceMember(suite.DB())
	mtoFetcher := NewCustomerFetcher()

	actualCustomer, err := mtoFetcher.FetchCustomer(suite.TestAppContext(), customer.ID)
	suite.NoError(err)

	suite.Equal(customer.ID, actualCustomer.ID)
	suite.Equal(*customer.Edipi, *actualCustomer.Edipi)
	suite.Equal(customer.UserID, actualCustomer.UserID)
	suite.NotNil(customer.User)
}
