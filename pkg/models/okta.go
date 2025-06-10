package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
)

type OktaUserPayload struct {
	Profile  OktaProfile `json:"profile"`
	GroupIds []string    `json:"groupIds"`
}

type OktaUpdateProfile struct {
	Profile OktaProfile `json:"profile"`
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

type OktaError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorSummary string `json:"errorSummary"`
	ErrorLink    string `json:"errorLink"`
	ErrorId      string `json:"errorId"`
	ErrorCauses  []struct {
		ErrorSummary string `json:"errorSummary"`
	} `json:"errorCauses"`
}

type OktaGroupProfile struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OktaGroup struct {
	ID      string           `json:"id"`
	Profile OktaGroupProfile `json:"profile"`
}

type OktaStatus string

const (
	OktaStatusActive        OktaStatus = "ACTIVE"
	OktaStatusDeprovisioned OktaStatus = "DEPROVISIONED"
	OktaStatusProvisioned   OktaStatus = "PROVISIONED"
	OktaStatusSuspended     OktaStatus = "SUSPENDED"
)

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

func GetOktaAPIKey() (key string) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	return v.GetString(cli.OktaAPIKeyFlag)
}

// OKTA USER FETCH //
// handles getting a single okta user by their okta id
func GetOktaUser(appCtx appcontext.AppContext, provider *okta.Provider, oktaID string, apiKey string) (*CreatedOktaUser, error) {
	baseURL := provider.GetUserURL(oktaID)

	// making HTTP request to Okta Users API to get a user
	// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/User/#tag/User/operation/getUser
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		appCtx.Logger().Error("could not create GET request", zap.Error(err))
		return nil, err
	}
	req.Header.Add("Authorization", "SSWS "+apiKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		appCtx.Logger().Error("could not execute GET request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	postResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		appCtx.Logger().Error("could not read GET response", zap.Error(err))
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

// OKTA ACCOUNT FETCHING SEVERAL USERS //
// fetching existing users by email/edipi
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
// this should only be used after validating a user doesn't exist with the email/edipi values
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

// OKTA ACCOUNT UPDATE //
// handles updating an existing okta user by providing their okta id and new profile information
// this is done via post so it is important to include all profile data by fetching first
func UpdateOktaUser(appCtx appcontext.AppContext, provider *okta.Provider, oktaID string, apiKey string, profile CreatedOktaUser) (*CreatedOktaUser, error) {
	baseURL := provider.GetUserURL(oktaID)
	body, err := json.Marshal(profile)
	if err != nil {
		appCtx.Logger().Error("error marshaling payload", zap.Error(err))
		return nil, err
	}

	// making HTTP request to Okta Users API to get a user
	// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/User/#tag/User/operation/updateUser
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

// OKTA USER GROUP ASSOCIATIONS //
// this func handles showing all groups a user is a part of
func GetOktaUserGroups(appCtx appcontext.AppContext, provider *okta.Provider, apiKey, userID string) ([]OktaGroup, error) {
	u, err := url.Parse(provider.GetUserGroupsURL(userID))
	if err != nil {
		return nil, err
	}

	// this is done via a GET request for fetching all groups associated with a user
	// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/UserResources/#tag/UserResources/operation/listUserGroups
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		appCtx.Logger().Error("could not create GET request when fetching user groups", zap.Error(err))
		return nil, err
	}
	req.Header.Add("Authorization", "SSWS "+apiKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		appCtx.Logger().Error("could not execute GET request when fetching user groups", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		appCtx.Logger().Error("could not read GET response when fetching user groups", zap.Error(err))
		return nil, err
	}

	var groups []OktaGroup
	if err := json.Unmarshal(response, &groups); err != nil {
		appCtx.Logger().Error("could not unmarshal GET response when fetching user groups", zap.Error(err))
		return nil, err
	}
	return groups, nil
}

// OKTA ADDING USER TO GROUP //
// this func handles adding a user to the group ID that is provided
func AddOktaUserToGroup(appCtx appcontext.AppContext, provider *okta.Provider, apiKey, groupID string, userID string) error {
	u, err := url.Parse(provider.AddUserToGroupURL(groupID, userID))
	if err != nil {
		return err
	}

	// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/Group/#tag/Group/operation/assignUserToGroup
	req, err := http.NewRequest("PUT", u.String(), nil)
	if err != nil {
		appCtx.Logger().Error("could not create PUT request when adding user to Okta group", zap.Error(err))
		return err
	}
	req.Header.Add("Authorization", "SSWS "+apiKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		appCtx.Logger().Error("could not execute PUT request when adding user to Okta group", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		appCtx.Logger().Error("could not read PUT response when adding user to Okta group", zap.Error(err))
		return err
	}

	// this means we were successful since Okta sends back a 204
	if len(response) == 0 {
		return nil
	}

	var oktaErr OktaError
	if err := json.Unmarshal(response, &oktaErr); err != nil {
		appCtx.Logger().Error("could not unmarshal Okta error response", zap.Error(err))
		return err
	}

	// if we can see the error summary, we will send that back
	if oktaErr.ErrorSummary != "" {
		return errors.New(oktaErr.ErrorSummary)
	}

	return nil
}

// Deletes the Okta account tied to the provided oktaID and logs any errors
func DeleteOktaUserHandled(appCtx appcontext.AppContext, oktaID string) {
	if oktaID != "" {
		req := appCtx.HTTPRequest()
		if req == nil {
			appCtx.Logger().Error("failed to retrieve HTTP request from session")
			return
		}
		provider, err := okta.GetOktaProviderForRequest(req)
		if err != nil {
			appCtx.Logger().Error("error retrieving Okta provider: %w", zap.Error(err))
			return
		}
		apiKey := GetOktaAPIKey()
		err = DeleteOktaUser(appCtx, provider, oktaID, apiKey)
		if err != nil {
			appCtx.Logger().Error("error deleting user from okta: %w", zap.Error(err))
			return
		}
	}
}

// Deletes the Okta account tied to the provided oktaID
func DeleteOktaUser(appCtx appcontext.AppContext, provider *okta.Provider, oktaID string, apiKey string) error {
	if len(oktaID) == 0 {
		return fmt.Errorf("DeleteOktaUser was called with an empty oktaID")
	}

	baseURL := provider.GetUserURL(oktaID)

	// verify the okta user exists before we attempt to delete the account
	existingOktaUser, err := GetOktaUser(appCtx, provider, oktaID, apiKey)
	if err != nil {
		return fmt.Errorf("error getting Okta user prior to deletion: %w", err)
	}
	if existingOktaUser == nil {
		return fmt.Errorf("okta user cannot be nil when preparing to delete the account")
	}

	// Okta will only let you delete a user that is in DEPROVISIONED status.
	// Calling delete on a user that is in any status other than DEPROVISIONED will result in the account being deactivated (DEPROVISIONED).
	// Therefore, in order to actually delete an ACTIVE user (or any status other than DEPROVISIONED), we will need to call delete twice.
	deleteAttempts := 1
	if OktaStatus(existingOktaUser.Status) != OktaStatusDeprovisioned {
		deleteAttempts = 2
	}

	client := &http.Client{}

	for i := 1; i <= deleteAttempts; i++ {

		// making HTTP request to Okta Users API to delete a user
		// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/User/#tag/User/operation/deleteUser
		req, err := http.NewRequest("DELETE", baseURL, nil)

		if err != nil {
			appCtx.Logger().Error(fmt.Sprintf("could not create DELETE request on delete attempt #%d", i), zap.Error(err))
			return err
		}
		req.Header.Add("Authorization", "SSWS "+apiKey)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			appCtx.Logger().Error(fmt.Sprintf("could not execute the request when attempting to delete existing okta user on delete attempt #%d", i), zap.Error(err))
			return err
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			appCtx.Logger().Error(fmt.Sprintf("could not read the response when attempting deleting existing okta user on delete attempt #%d", i), zap.Error(err))
			return err
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("API error (status %d): %s  on delete attempt #%d", resp.StatusCode, string(responseBody), i)
		}
	}
	appCtx.Logger().Info(fmt.Sprintf("Successfully deleted Okta Account for oktaID %s", oktaID))
	return nil
}
