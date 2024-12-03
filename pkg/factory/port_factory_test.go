package factory

import "github.com/transcom/mymove/pkg/models"

func (suite *FactorySuite) TestFetchPort() {

	const defaultPortType = models.PortTypeAir
	const defaultPortCode = "PDX"
	const defaultPortName = "PORTLAND INTL"

	suite.Run("Successful fetch of default Port", func() {

		defaultPortLocation := FetchPort(nil, nil, nil)

		suite.Equal(defaultPortLocation.PortType, defaultPortType)
		suite.Equal(defaultPortLocation.PortCode, defaultPortCode)
		suite.Equal(defaultPortLocation.PortName, defaultPortName)
	})

	suite.Run("Successful fetch of Port using port code", func() {

		customPort := FetchPortLocation(suite.DB(), []Customization{
			{
				Model: models.Port{
					PortCode: "SEA",
				},
			},
		}, nil)

		suite.Equal("A", customPort.Port.PortType.String())
		suite.Equal("SEA", customPort.Port.PortCode)
		suite.Equal("SEATTLE TACOMA INTL", customPort.Port.PortName)
	})
}
