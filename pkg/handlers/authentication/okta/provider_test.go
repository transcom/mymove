package okta_test

import (
	"testing"

	"github.com/markbates/goth"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type OktaSuite struct {
	handlers.BaseHandlerTestSuite
}

func TestAuthSuite(t *testing.T) {
	hs := &OktaSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("milmovelocal"), testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

// TODO: Flesh out with actually accessing the data. See error that can pop up when okta provider is not wrapped correctly.
func (suite *OktaSuite) TestConvertGothProviderToOktaProvider() {
	oktaProvider := &okta.Provider{}
	gothProvider := goth.Provider(oktaProvider)

	provider, err := okta.ConvertGothProviderToOktaProvider(gothProvider)
	suite.NoError(err)
	suite.Equal(oktaProvider, provider)
}
