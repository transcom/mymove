package dpsauth

// Params contains configuration params for DPS authentication
type Params struct {
	SDDCProtocol   string
	SDDCHostname   string
	SDDCPort       string
	SecretKey      string
	DPSRedirectURL string
	CookieName     string
	CookieDomain   string
	CookieSecret   []byte
	CookieExpires  int
}
