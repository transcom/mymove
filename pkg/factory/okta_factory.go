package factory

import (
	"github.com/markbates/goth"
	gothOkta "github.com/markbates/goth/providers/okta"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
)

const DummyOktaOrgURL = "https://dummy.okta.com"
const DummyOktaCallbackURL = "https://dummy.okta.com/auth/callback"
const DummyClientID = "dummyClientID"
const DummySecret = "dummySecret"

var DummyOIDCScope = []string{"openid", "profile", "email"}

type ProviderConfig struct {
	Name        string
	OrgURL      string
	CallbackURL string
	ClientID    string
	Secret      string
	Scope       []string
	Logger      *zap.Logger
}

func CreateAndWrapProvider(config ProviderConfig) (*okta.Provider, error) {
	// Create a new Okta provider with goth
	provider := gothOkta.New(config.ClientID, config.Secret, config.OrgURL, config.CallbackURL, config.Scope...)
	// Set the name in goth
	provider.SetName(config.Name)

	// For return, the gothOkta provider wrapped with our own provider struct
	wrappedProvider := okta.WrapOktaProvider(provider, config.OrgURL, config.ClientID, config.Secret, config.CallbackURL, config.Logger)
	// Assign to the active goth providers
	goth.UseProviders(wrappedProvider)

	return wrappedProvider, nil
}

func BuildOktaProvider(name string) (*okta.Provider, error) {
	logger, _ := zap.NewDevelopment()

	providerConfig := ProviderConfig{
		Name:        name,
		OrgURL:      DummyOktaOrgURL,
		CallbackURL: DummyOktaCallbackURL,
		ClientID:    DummyClientID,
		Secret:      DummySecret,
		Scope:       DummyOIDCScope,
		Logger:      logger,
	}

	return CreateAndWrapProvider(providerConfig)
}
