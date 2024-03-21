package models

type OktaUserPayload struct {
	Profile  Profile  `json:"profile"`
	GroupIds []string `json:"groupIds"`
}

type Profile struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Login       string `json:"login"`
	MobilePhone string `json:"mobilePhone"`
}
