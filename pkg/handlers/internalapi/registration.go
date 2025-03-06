package internalapi

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/cli"
	registrationop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/registration"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/handlers/internalapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
)

// CustomerRegistrationHandler creates a MilMove and Okta profile allowing for self registration of service members
type CustomerRegistrationHandler struct {
	handlers.HandlerConfig
}

func (h CustomerRegistrationHandler) Handle(params registrationop.CustomerRegistrationParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			if !appCtx.Session().IsMilApp() {
				return registrationop.NewCustomerRegistrationUnprocessableEntity(), apperror.NewSessionError("Request is not from the customer app")
			}

			oktaUser, oktaErr := fetchOrCreateOktaProfile(appCtx, params)
			if oktaErr != nil || oktaUser == nil {
				appCtx.Logger().Error("error creating okta profile", zap.Error(oktaErr))
				errPayload := payloads.ValidationError(
					oktaErr.Error(),
					h.GetTraceIDFromRequest(params.HTTPRequest),
					nil,
				)
				return registrationop.NewCustomerRegistrationUnprocessableEntity().WithPayload(errPayload), apperror.NewSessionError("Error")
			}

			transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
				oktaSub := oktaUser.ID
				payload := params.Registration

				var user *models.User
				user, userErr := models.GetUserFromOktaID(appCtx.DB(), oktaSub)
				if userErr != nil {
					appCtx.Logger().Error("error fetching user", zap.Error(userErr))
					return userErr
				}

				// if user doesn't exist, we need to create one
				if user == nil {
					user, userErr = models.CreateUser(appCtx.DB(), oktaSub, payload.Email)
					if userErr != nil {
						appCtx.Logger().Error("error creating user", zap.Error(userErr))
						return userErr
					}
				}

				userID := user.ID

				// now we need to see if the service member exists based off of the user id we have now
				existingServiceMember, smErr := models.FetchServiceMemberByUserID(appCtx.DB(), userID.String())
				if smErr != sql.ErrNoRows && smErr != nil {
					appCtx.Logger().Error("error creating service member", zap.Error(smErr))
					return smErr
				}

				// evaluating feature flag to see if we need to check if the DODID exists already
				var dodidUniqueFeatureFlag bool
				featureFlagName := "dodid_unique"
				flag, err := h.FeatureFlagFetcher().GetBooleanFlag(params.HTTPRequest.Context(), appCtx.Logger(), "customer", featureFlagName, map[string]string{})
				if err != nil {
					appCtx.Logger().Error("Error fetching dodid_unique feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
					dodidUniqueFeatureFlag = false
				} else {
					dodidUniqueFeatureFlag = flag.Match
				}

				var serviceMembers []models.ServiceMember
				if dodidUniqueFeatureFlag {
					query := `SELECT service_members.edipi
								FROM service_members
								WHERE service_members.edipi = $1`
					err := appCtx.DB().RawQuery(query, payload.Edipi).All(&serviceMembers)
					if err != nil {
						errorMsg := apperror.NewBadDataError("error when checking for existing service member")
						return errorMsg
					} else if len(serviceMembers) > 0 {
						errorMsg := apperror.NewConflictError(h.GetTraceIDFromRequest(params.HTTPRequest), "Service member with this DODID already exists. Please use a different DODID number.")
						return errorMsg
					}
				}

				// if we do not have a service member, we can now create one
				if existingServiceMember == nil {
					serviceMember := models.ServiceMember{
						UserID:             userID,
						Edipi:              payload.Edipi,
						Emplid:             payload.Emplid,
						Affiliation:        (*models.ServiceMemberAffiliation)(payload.Affiliation),
						FirstName:          &payload.FirstName,
						MiddleName:         payload.MiddleInitial,
						LastName:           &payload.LastName,
						Telephone:          &payload.Telephone,
						SecondaryTelephone: &payload.SecondaryTelephone,
						PersonalEmail:      &payload.Email,
						PhoneIsPreferred:   &payload.PhoneIsPreferred,
						EmailIsPreferred:   &payload.EmailIsPreferred,
					}
					smVerrs, smErr := models.SaveServiceMember(appCtx, &serviceMember)
					if smVerrs.HasAny() || smErr != nil {
						appCtx.Logger().Error("error updating service member", zap.Error(smErr))
						return smErr
					}

					return nil
				}

				return nil
			})

			if transactionError != nil {
				appCtx.Logger().Error("error occurred while service member tried to register an account", zap.Error(transactionError))
				errPayload := payloads.ValidationError(
					transactionError.Error(),
					h.GetTraceIDFromRequest(params.HTTPRequest),
					nil,
				)
				switch transactionError.(type) {
				case *pq.Error:
					return registrationop.NewCustomerRegistrationUnprocessableEntity().WithPayload(errPayload), transactionError
				default:
					return registrationop.NewCustomerRegistrationInternalServerError(), transactionError
				}
			}

			return registrationop.NewCustomerRegistrationCreated(), nil
		})
}

// fetchOrCreateOktaProfile send some requests to the Okta Users API
// handles seeing if an okta user already exists with the form data, if not - it will then create one
// this creates a user in Okta assigned to the customer group (allowing access to the customer application)
func fetchOrCreateOktaProfile(appCtx appcontext.AppContext, params registrationop.CustomerRegistrationParams) (*models.CreatedOktaUser, error) {
	// setting viper so we can access the api key in the env vars
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	apiKey := v.GetString(cli.OktaAPIKeyFlag)
	customerGroupID := v.GetString(cli.OktaCustomerGroupIDFlag)

	// taking all the data that we'll need for the okta profile creation
	payload := params.Registration
	oktaEmail := payload.Email
	oktaFirstName := payload.FirstName
	oktaLastName := payload.LastName
	oktaPhone := payload.Telephone
	oktaEdipi := payload.Edipi

	// getting the right okta provider
	provider, err := okta.GetOktaProviderForRequest(params.HTTPRequest)
	if err != nil {
		return nil, err
	}

	// OKTA ACCOUNT FETCHING //
	// build the search filter according to Okta's syntax
	searchFilter := fmt.Sprintf(`profile.email eq "%s" or profile.cac_edipi eq "%s"`, oktaEmail, *oktaEdipi)
	getUsersURL := provider.GetUsersURL()
	u, err := url.Parse(getUsersURL)
	if err != nil {
		return nil, err
	}

	// adding the search parameter
	q := u.Query()
	q.Set("search", searchFilter)
	u.RawQuery = q.Encode()

	// making HTTP request to Okta Users API to list all users
	// this is done via a GET request for creating a user that sends an activation email (when activate=true)
	// https://developer.okta.com/docs/api/openapi/okta-management/management/tag/User/#tag/User/operation/listUsers
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		appCtx.Logger().Error("could not execute request", zap.Error(err))
		return nil, err
	}
	h := req.Header
	h.Add("Authorization", "SSWS "+apiKey)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		appCtx.Logger().Error("could not execute request", zap.Error(err))
		return nil, err
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		appCtx.Logger().Error("could not read response body", zap.Error(err))
		return nil, err
	}

	// we get an array back from Okta, so we need to unmarshal their response into our structs
	var users []models.CreatedOktaUser
	if err := json.Unmarshal([]byte(response), &users); err != nil {
		appCtx.Logger().Error("could not unmarshal body", zap.Error(err))
		return nil, err
	}

	// checking if we have existing or conflicts okta users in our organization based on submitted form values
	// existing edipi & email that match -> send back that okta user, we don't need to create
	// existing email but edipi doesn't match that profile -> send back an error
	// existing edipi but email doesn't match that profile -> send back an error
	if len(users) > 0 {
		var oktaUser *models.CreatedOktaUser
		exactMatch := false
		emailMatch := false
		edipiMatch := false
		for i, user := range users {
			if oktaEmail != "" && oktaEdipi != nil && user.Profile.Email != "" && user.Profile.CacEdipi != nil {
				if user.Profile.Email == oktaEmail && *user.Profile.CacEdipi == *oktaEdipi {
					exactMatch = true
					oktaUser = &users[i]
					break
				}
			}
			if oktaEmail != "" && user.Profile.Email != "" {
				if user.Profile.Email == oktaEmail {
					emailMatch = true
				}
			}
			if oktaEdipi != nil && user.Profile.CacEdipi != nil {
				if *user.Profile.CacEdipi == *oktaEdipi {
					edipiMatch = true
				}
			}
		}

		if exactMatch {
			return oktaUser, nil
		}

		// if we get more than one result, we need to handle the error returns differently than if we just have one existing okta user but not an exact match
		if emailMatch && !edipiMatch && len(users) > 1 {
			return nil, fmt.Errorf("email and DoD IDs match different users - please open up a help desk ticket")
		} else if emailMatch && !edipiMatch && len(users) == 1 {
			return nil, fmt.Errorf("there is an existing okta account with that email - please update the DoD ID (EDIPI) in your okta profile to match your registration DoD ID and try registering again")
		}

		if !emailMatch && edipiMatch && len(users) > 1 {
			return nil, fmt.Errorf("email and DoD IDs match different users - please open up a help desk ticket")
		} else if !emailMatch && edipiMatch && len(users) == 1 {
			return nil, fmt.Errorf("there is an existing okta account with that DoD ID (EDIPI) - please update the email in your okta profile to match your registration email and try registering again")
		}

		// if we get an email & edipi match on two different users and NOT an exact match, we need them to open a HDT
		if emailMatch && edipiMatch && len(users) > 1 {
			return nil, fmt.Errorf("there are multiple okta accounts with that email and DoD ID - please open up a help desk ticket")
		}

		return nil, fmt.Errorf("okta account creation error - please open up a help desk ticket")
	}

	// OKTA ACCOUNT CREATION //
	// now we will create the okta account since we now know it doesn't exist with our unique values (email/edipi)
	// active = true meanas that the user will get an activation email from Okta
	activate := "true"
	baseURL := provider.GetCreateUserURL(activate)

	profile := models.Profile{
		FirstName:   oktaFirstName,
		LastName:    oktaLastName,
		Email:       oktaEmail,
		Login:       oktaEmail,
		MobilePhone: oktaPhone,
		CacEdipi:    *oktaEdipi,
	}

	// okta needs a certain structure in the request, assigning the user to the customer app
	oktaPayload := models.OktaUserPayload{
		Profile:  profile,
		GroupIds: []string{customerGroupID},
	}

	body, err := json.Marshal(oktaPayload)
	if err != nil {
		appCtx.Logger().Error("error marshaling payload", zap.Error(err))
		return nil, err
	}

	// making HTTP request to Okta Users API to create a user
	// this is done via a POST request for creating a user that sends an activation email (when activate=true)
	// https://developer.okta.com/docs/reference/api/users/#create-user-without-credentials
	req, err = http.NewRequest("POST", baseURL, bytes.NewReader(body))
	if err != nil {
		appCtx.Logger().Error("could not execute request", zap.Error(err))
		return nil, err
	}
	h = req.Header
	h.Add("Authorization", "SSWS "+apiKey)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		appCtx.Logger().Error("could not execute request", zap.Error(err))
		return nil, err
	}

	response, err = io.ReadAll(resp.Body)
	if err != nil {
		appCtx.Logger().Error("could not read response body", zap.Error(err))
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		apiErr := fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(response))
		if resp.StatusCode == http.StatusInternalServerError {
			return nil, apiErr
		}
		if resp.StatusCode == http.StatusBadRequest {
			return nil, apiErr
		}
		if resp.StatusCode == http.StatusForbidden {
			return nil, apiErr
		}
	}

	user := models.CreatedOktaUser{}
	err = json.Unmarshal(response, &user)
	if err != nil {
		appCtx.Logger().Error("could not unmarshal body", zap.Error(err))
		return nil, err
	}

	defer resp.Body.Close()

	return &user, nil
}
