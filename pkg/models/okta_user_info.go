package models

type OktaUser struct {
	Sub               string `json:"sub"`
	Name              string `json:"name"`
	Locale            string `json:"locale"`
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	FamilyName        string `json:"family_name"`
	GivenName         string `json:"given_name"`
	ZoneInfo          string `json:"zoneinfo"`
	UpdatedAt         int    `json:"updated_at"`
	EmailVerified     bool   `json:"email_verified"`
	Edipi             string `json:"cac_edipi"`
}

type CreatedOktaUser struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Created   string `json:"created"`
	Activated string `json:"activated"`
	Profile   struct {
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		MobilePhone string `json:"mobilePhone"`
		SecondEmail string `json:"secondEmail"`
		Login       string `json:"login"`
		Email       string `json:"email"`
	} `json:"profile"`
}
