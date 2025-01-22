package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
)

func (suite *ModelSuite) TestPortLocation() {
	suite.Run("Port location has correct fields", func() {

		portLocation := factory.FetchPortLocation(suite.DB(), nil, nil)

		suite.NotNil(portLocation)
		suite.Equal(portLocation.PortId, portLocation.Port.ID)
		suite.Equal(portLocation.CitiesId, portLocation.City.ID)
		suite.Equal(portLocation.UsPostRegionCitiesId, portLocation.UsPostRegionCity.ID)
		suite.Equal(portLocation.CountryId, portLocation.Country.ID)
	})

	suite.Run("Port location table name is correct", func() {

		portLocation := factory.FetchPortLocation(suite.DB(), nil, nil)

		suite.NotNil(portLocation)
		suite.Equal("port_locations", portLocation.TableName())
	})
}
