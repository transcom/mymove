package factory

import (
	gothOkta "github.com/markbates/goth/providers/okta"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
)

type ProviderConfig struct {
	Name        string
	OrgURL      string
	CallbackURL string
	ClientID    string
	Secret      string
	Scope       []string
	Logger      *zap.Logger
}

func NewProviderFactory(config ProviderConfig) (*okta.Provider, error) {
	// Create a new Okta provider with goth
	provider := gothOkta.New(config.ClientID, config.Secret, config.OrgURL, config.CallbackURL, config.Scope...)
	// Set the name in goth
	provider.SetName(config.Name)

	// Return the gothOkta provider wrapped with out our own provider struct
	return okta.WrapOktaProvider(provider, config.OrgURL, config.ClientID, config.Secret, config.CallbackURL, config.Logger), nil
}

func DummyProviderFactory(name string) (*okta.Provider, error) {
	logger, _ := zap.NewDevelopment()

	// TODO: replace with consts
	dummyConfig := ProviderConfig{
		Name:        name,
		OrgURL:      "https://dummy.okta.com",
		CallbackURL: "https://dummy-callback.com",
		ClientID:    "dummyClientID",
		Secret:      "dummySecret",
		Scope:       []string{"openid", "profile", "email"},
		Logger:      logger,
	}

	return NewProviderFactory(dummyConfig)
}
