package models

// OktaAccountCreationTemplate is a template of information needed in the okta account creation service
type OktaAccountCreationTemplate struct {
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Login       string `json:"login"`
	Email       string `json:"email"`
	MobilePhone string `json:"mobilePhone"`
	CacEdipi    string `json:"cacedipi"`
	GsaID       string `json:"gsaid"`
}

// Okta account POST Req body profile
type OktaBodyProfile struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Login       string `json:"login"`
	MobilePhone string `json:"mobilePhone"`
	CacEdipi    string `json:"cac_edipi"`
	GsaID       string `json:"gsa_id"`
}

// Okta account POST Req body
type OktaAccountCreationBody struct {
	Profile  OktaBodyProfile `json:"profile"`
	GroupIds []string        `json:"groupIds"`
}
