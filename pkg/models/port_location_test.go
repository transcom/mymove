package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
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

func (suite *ModelSuite) TestFetchPortLocationByCode() {
	suite.Run("Port location can be fetched when it exists", func() {

		portLocation := factory.FetchPortLocation(suite.DB(), []factory.Customization{
			{
				Model: models.Port{
					PortCode: "SEA",
				},
			},
		}, nil)
		suite.NotNil(portLocation)

		result, err := models.FetchPortLocationByCode(suite.AppContextForTest().DB(), "SEA")
		suite.NotNil(result)
		suite.NoError(err)
		suite.Equal(portLocation.ID, result.ID)
	})

	suite.Run("Sends back an error when it does not exist", func() {
		result, err := models.FetchPortLocationByCode(suite.AppContextForTest().DB(), "123")
		suite.Nil(result)
		suite.Error(err)
	})
}
