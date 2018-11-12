package authentication

// LoginGovConfig contains values needed to register login.gov callbacks for the sites
type LoginGovConfig struct {
	Host             string
	Secret           string
	MyClientID       string
	OfficeClientID   string
	TspClientID      string
	CallbackProtocol string
	CallbackPort     int
}
