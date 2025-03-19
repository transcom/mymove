package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
)

type OktaUserPayload struct {
	Profile  OktaProfile `json:"profile"`
	GroupIds []string    `json:"groupIds"`
}

type OktaProfile struct {
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	Email       string  `json:"email"`
	Login       string  `json:"login"`
	MobilePhone string  `json:"mobilePhone"`
	CacEdipi    string  `json:"cac_edipi"`
	GsaID       *string `json:"gsa_id"`
}

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
		FirstName   string  `json:"firstName"`
		LastName    string  `json:"lastName"`
		MobilePhone string  `json:"mobilePhone"`
		SecondEmail string  `json:"secondEmail"`
		Login       string  `json:"login"`
		Email       string  `json:"email"`
		CacEdipi    *string `json:"cac_edipi"`
		GsaID       *string `json:"gsa_id"`
	} `json:"profile"`
}

// ensures a valid email address
func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// ensures edipi is 10 digits
func isValidEdipi(edipi string) bool {
	edipiRegex := `^\d{10}$`
	re := regexp.MustCompile(edipiRegex)
	return re.MatchString(edipi)
}

// OKTA ACCOUNT FETCHING //
// we need to first check if there is an existing okta user before creating one
// email and edipi are unique in okta, so searching for those should be enough to ensure there isn't an existing account
// gsaID is used for office users that do not use the typical EDIPI - this will be nil when searching for existing customers
func SearchForExistingOktaUsers(appCtx appcontext.AppContext, provider *okta.Provider, apiKey, oktaEmail string, oktaEdipi *string, gsaID *string) ([]CreatedOktaUser, error) {
	if oktaEmail == "" {
		return nil, fmt.Errorf("email is required and cannot be empty")
	}
	if !isValidEmail(oktaEmail) {
		return nil, fmt.Errorf("invalid email format: %s", oktaEmail)
	}

	if oktaEdipi != nil && *oktaEdipi != "" {
		if !isValidEdipi(*oktaEdipi) {
			return nil, fmt.Errorf("invalid EDIPI format: %s", *oktaEdipi)
		}
	}

	searchFilter := fmt.Sprintf(`profile.email eq "%s"`, oktaEmail)
	if oktaEdipi != nil && *oktaEdipi != "" {
		searchFilter += fmt.Sprintf(` or profile.cac_edipi eq "%s"`, *oktaEdipi)
	}
	if gsaID != nil && *gsaID != "" {
		searchFilter += fmt.Sprintf(` or profile.gsa_id eq "%s"`, *gsaID)
	}

	u, err := url.Parse(provider.GetUsersURL())
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("search", searchFilter)
	u.RawQuery = q.Encode()

	// making HTTP request to Okta Users API to list all users
	// this is done via a GET request for fetching all users based on the provided search parameters
	// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/User/#tag/User/operation/listUsers
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		appCtx.Logger().Error("could not create GET request when fetching existing okta users for registration", zap.Error(err))
		return nil, err
	}
	req.Header.Add("Authorization", "SSWS "+apiKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		appCtx.Logger().Error("could not execute GET request when fetching existing okta users for registration", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		appCtx.Logger().Error("could not read GET response when fetching existing okta users for registration", zap.Error(err))
		return nil, err
	}

	var users []CreatedOktaUser
	if err := json.Unmarshal(response, &users); err != nil {
		appCtx.Logger().Error("could not unmarshal GET response when fetching existing okta users for registration", zap.Error(err))
		return nil, err
	}
	return users, nil
}

// OKTA ACCOUNT CREATION //
// we have validated an existing account doesn't exist, so it is now safe to create one
// assigning to the customer group so they can successfully get through the okta validations
func CreateOktaUser(appCtx appcontext.AppContext, provider *okta.Provider, apiKey string, payload OktaUserPayload) (*CreatedOktaUser, error) {
	activate := "true"
	baseURL := provider.GetCreateUserURL(activate)
	body, err := json.Marshal(payload)
	if err != nil {
		appCtx.Logger().Error("error marshaling payload", zap.Error(err))
		return nil, err
	}

	// making HTTP request to Okta Users API to create a user
	// this is done via a POST request for creating a user that sends an activation email (when activate=true)
	// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/User/#tag/User/operation/createUser
	req, err := http.NewRequest("POST", baseURL, bytes.NewReader(body))
	if err != nil {
		appCtx.Logger().Error("could not create POST request", zap.Error(err))
		return nil, err
	}
	req.Header.Add("Authorization", "SSWS "+apiKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		appCtx.Logger().Error("could not execute POST request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	postResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		appCtx.Logger().Error("could not read POST response", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(postResponse))
	}

	var createdUser CreatedOktaUser
	if err = json.Unmarshal(postResponse, &createdUser); err != nil {
		appCtx.Logger().Error("could not unmarshal POST response when creating Okta user", zap.Error(err))
		return nil, err
	}
	return &createdUser, nil
}
