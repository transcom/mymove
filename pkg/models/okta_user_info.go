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
	UpdatedAt         string `json:"updated_at"`
	EmailVerified     string `json:"email_verified"`
	Edipi             string `json:"cac_edipi"`
}
