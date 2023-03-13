package customer

import (
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *CustomerServiceSuite) TestCustomerFetcher() {
	customer := factory.BuildServiceMember(suite.DB(), nil, nil)
	mtoFetcher := NewCustomerFetcher()

	actualCustomer, err := mtoFetcher.FetchCustomer(suite.AppContextForTest(), customer.ID)
	suite.NoError(err)

	suite.Equal(customer.ID, actualCustomer.ID)
	suite.Equal(*customer.Edipi, *actualCustomer.Edipi)
	suite.Equal(customer.UserID, actualCustomer.UserID)
	suite.NotNil(customer.User)
}
