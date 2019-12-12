package customer

import (
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *CustomerServiceSuite) TestMoveTaskOrderFetcher() {
	customer := testdatagen.MakeDefaultCustomer(suite.DB())
	mtoFetcher := NewCustomerFetcher(suite.DB())

	actualCustomer, err := mtoFetcher.FetchCustomer(customer.ID)
	suite.NoError(err)

	suite.Equal(customer.ID, actualCustomer.ID)
	suite.Equal(customer.DODID, actualCustomer.DODID)
	suite.Equal(customer.UserID, actualCustomer.UserID)
	suite.NotNil(customer.User)
}
