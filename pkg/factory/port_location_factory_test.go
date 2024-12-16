package factory

import "github.com/transcom/mymove/pkg/models"

func (suite *FactorySuite) TestFetchPortLocation() {

	const defaultPortType = models.PortTypeAir
	const defaultPortCode = "PDX"
	const defaultPortName = "PORTLAND INTL"
	const defaultCityName = "PORTLAND"
	const defaultCountyName = "MULTNOMAH"
	const defaultZip = "97220"
	const defaultStateName = "OREGON"
	const defaultCountryName = "UNITED STATES"

	suite.Run("Successful fetch of default Port Location", func() {

		defaultPortLocation := models.PortLocation{
			Port: models.Port{
				PortType: defaultPortType,
				PortCode: defaultPortCode,
				PortName: defaultPortName,
			},
			City: models.City{
				CityName: defaultCityName,
			},
			UsPostRegionCity: models.UsPostRegionCity{
				UsprcCountyNm: defaultCountyName,
				UsprZipID:     defaultZip,
				UsPostRegion: models.UsPostRegion{
					State: models.State{
						StateName: defaultStateName,
					},
				},
			},
			Country: models.Country{
				CountryName: defaultCountryName,
			},
		}

		portLocation := FetchPortLocation(suite.DB(), nil, nil)

		suite.Equal(defaultPortLocation.Port.PortType.String(), portLocation.Port.PortType.String())
		suite.Equal(defaultPortLocation.Port.PortCode, portLocation.Port.PortCode)
		suite.Equal(defaultPortLocation.Port.PortName, portLocation.Port.PortName)
		suite.Equal(defaultPortLocation.City.CityName, portLocation.City.CityName)
		suite.Equal(defaultPortLocation.UsPostRegionCity.UsprcCountyNm, portLocation.UsPostRegionCity.UsprcCountyNm)
		suite.Equal(defaultPortLocation.UsPostRegionCity.UsprZipID, portLocation.UsPostRegionCity.UsprZipID)
		suite.Equal(defaultPortLocation.UsPostRegionCity.UsPostRegion.State.StateName, portLocation.UsPostRegionCity.UsPostRegion.State.StateName)
		suite.Equal(defaultPortLocation.Country.CountryName, portLocation.Country.CountryName)
	})

	suite.Run("Successful fetch of Port Location using port code", func() {

		customPortLocation := FetchPortLocation(suite.DB(), []Customization{
			{
				Model: models.Port{
					PortCode: "SEA",
				},
			},
		}, nil)

		suite.Equal("A", customPortLocation.Port.PortType.String())
		suite.Equal("SEA", customPortLocation.Port.PortCode)
		suite.Equal("SEATTLE TACOMA INTL", customPortLocation.Port.PortName)
		suite.Equal("SEATTLE", customPortLocation.City.CityName)
		suite.Equal("KING", customPortLocation.UsPostRegionCity.UsprcCountyNm)
		suite.Equal("98158", customPortLocation.UsPostRegionCity.UsprZipID)
		suite.Equal("WASHINGTON", customPortLocation.UsPostRegionCity.UsPostRegion.State.StateName)
		suite.Equal("UNITED STATES", customPortLocation.Country.CountryName)
	})
}
