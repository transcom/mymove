package factory_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/factory"
)

type FactorySuite struct {
	suite.Suite
}

func TestFactorySuite(t *testing.T) {
	suite.Run(t, new(FactorySuite))
}

func (suite *FactorySuite) TestBuildOktaProvider() {
	name := "TestProvider"
	provider, err := factory.BuildOktaProvider(name)

	suite.NoError(err)
	suite.Equal(factory.DummyOktaOrgURL, provider.GetOrgURL())
	suite.Equal(factory.DummyClientID, provider.GetClientID())
}
