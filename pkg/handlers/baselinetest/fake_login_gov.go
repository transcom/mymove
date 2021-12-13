package baselinetest

import (
	"net/url"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"

	"github.com/transcom/mymove/pkg/handlers/authentication"
)

type FakeGothSession struct {
	authURL string
	state   string
}

func (s *FakeGothSession) GetAuthURL() (string, error) {
	if s.state != "" {
		v := url.Values{}
		v.Set("state", s.state)
		return s.authURL + "?" + v.Encode(), nil
	}
	return s.authURL, nil
}

func (s *FakeGothSession) Marshal() string {
	return ""
}

func (s *FakeGothSession) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	return "", nil
}

type FakeGothProvider struct {
	name        string
	clientKey   string
	baseAuthURL string
}

func (p *FakeGothProvider) Name() string {
	return p.name
}

func (p *FakeGothProvider) SetName(name string) {
	p.name = name
}

func (p *FakeGothProvider) BeginAuth(state string) (goth.Session, error) {
	return &FakeGothSession{
		authURL: p.baseAuthURL,
		state:   state,
	}, nil
}

func (p *FakeGothProvider) UnmarshalSession(string) (goth.Session, error) {
	return &FakeGothSession{
		authURL: p.baseAuthURL,
		state:   "",
	}, nil
}
func (p *FakeGothProvider) FetchUser(goth.Session) (goth.User, error) {
	return goth.User{}, nil
}
func (p *FakeGothProvider) Debug(bool) {}
func (p *FakeGothProvider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	return nil, nil
}
func (p *FakeGothProvider) RefreshTokenAvailable() bool {
	return true
}

func (p *FakeGothProvider) ClientKey() string {
	return p.clientKey
}

func (suite *BaselineSuite) initFakeLoginGovProvider() authentication.LoginGovProvider {
	fakeLoginGovHost := "fake-login-gov.example.com"
	fakeLoginGovURL := "http://" + fakeLoginGovHost
	p := authentication.NewLoginGovProvider(fakeLoginGovHost, "secret_key", suite.Logger())

	milProvider := &FakeGothProvider{
		baseAuthURL: fakeLoginGovURL,
		clientKey:   "milClientKey",
	}
	officeProvider := &FakeGothProvider{
		baseAuthURL: fakeLoginGovURL,
		clientKey:   "officeClientKey",
	}
	adminProvider := &FakeGothProvider{
		baseAuthURL: fakeLoginGovURL,
		clientKey:   "adminClientKey",
	}
	authentication.SetLoginGovProviders(milProvider, officeProvider, adminProvider)
	return p
}
