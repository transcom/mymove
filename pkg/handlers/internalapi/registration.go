package internalapi

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
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

			// evaluating feature flag to see if we need to check if the DODID exists already
			// this is to prevent duplicate service_member accounts
			var dodidUniqueFeatureFlag bool
			featureFlagName := "dodid_unique"
			flag, err := h.FeatureFlagFetcher().GetBooleanFlag(params.HTTPRequest.Context(), appCtx.Logger(), "customer", featureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching dodid_unique feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
				dodidUniqueFeatureFlag = false
			} else {
				dodidUniqueFeatureFlag = flag.Match
			}

			transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
				oktaSub := oktaUser.ID
				payload := params.Registration

				var user *models.User
				user, userErr := models.GetUserFromOktaID(appCtx.DB(), oktaSub)
				if userErr != sql.ErrNoRows && userErr != nil {
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

				// now we need to see if the service member exists based off of the user id we have now
				existingServiceMember, smErr := models.FetchServiceMemberByUserID(appCtx.DB(), user.ID.String())
				if smErr != sql.ErrNoRows && smErr != nil {
					appCtx.Logger().Error("error creating service member", zap.Error(smErr))
					return smErr
				}

				// if we couldn't find an existing service member with the okta_id
				// we need to ensure we don't have an existing SM with the same edipi
				// this will only be checked if dodid_unique flag is on
				var serviceMembers []models.ServiceMember
				if existingServiceMember == nil && dodidUniqueFeatureFlag {
					query := `SELECT service_members.edipi
								FROM service_members
								WHERE service_members.edipi = $1`
					err := appCtx.DB().RawQuery(query, payload.Edipi).All(&serviceMembers)
					if err != nil {
						errorMsg := apperror.NewBadDataError("error when checking for existing service member")
						return errorMsg
					} else if len(serviceMembers) > 0 {
						errorMsg := fmt.Errorf("there is already an existing MilMove user with this DoD ID - an Okta account has also been found or created, please try signing into MilMove instead")
						return errorMsg
					}
				}

				// if we do not have a service member, we can now create one
				if existingServiceMember == nil {
					serviceMember := models.ServiceMember{
						UserID:             user.ID,
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
				return registrationop.NewCustomerRegistrationUnprocessableEntity().WithPayload(errPayload), transactionError
			}

			return registrationop.NewCustomerRegistrationCreated(), nil
		})
}

func getCustomerGroupID() (apiKey, customerGroupID string) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	return v.GetString(cli.OktaAPIKeyFlag), v.GetString(cli.OktaCustomerGroupIDFlag)
}

// fetchOrCreateOktaProfile send some requests to the Okta Users API
// handles seeing if an okta user already exists with the form data, if not - it will then create one
// this creates a user in Okta assigned to the customer group (allowing access to the customer application)
func fetchOrCreateOktaProfile(appCtx appcontext.AppContext, params registrationop.CustomerRegistrationParams) (*models.CreatedOktaUser, error) {
	apiKey, customerGroupID := getCustomerGroupID()

	payload := params.Registration
	oktaEmail := payload.Email
	oktaFirstName := payload.FirstName
	oktaLastName := payload.LastName
	oktaPhone := payload.Telephone
	oktaEdipi := payload.Edipi

	provider, err := okta.GetOktaProviderForRequest(params.HTTPRequest)
	if err != nil {
		return nil, err
	}

	users, err := models.SearchForExistingOktaUsers(appCtx, provider, apiKey, oktaEmail, oktaEdipi, nil)
	if err != nil {
		return nil, err
	}

	// checking if we have existing and/or mismatched okta users in our organization based on submitted form values
	// existing edipi & email that match -> send back that okta user, we don't need to create
	// existing email but edipi doesn't match that profile -> send back an error
	// existing edipi but email doesn't match that profile -> send back an error
	if len(users) > 0 {
		var oktaUser *models.CreatedOktaUser
		var exactMatch, emailMatch, edipiMatch bool
		for i, user := range users {
			if oktaEmail != "" && oktaEdipi != nil && user.Profile.Email != "" && user.Profile.CacEdipi != nil {
				if user.Profile.Email == oktaEmail && *user.Profile.CacEdipi == *oktaEdipi {
					exactMatch = true
					oktaUser = &users[i]
					break
				}
			}
			if oktaEmail != "" && user.Profile.Email == oktaEmail {
				emailMatch = true
			}
			if oktaEdipi != nil && user.Profile.CacEdipi != nil && *user.Profile.CacEdipi == *oktaEdipi {
				edipiMatch = true
			}
		}

		if exactMatch {
			return oktaUser, nil
		}
		if emailMatch && !edipiMatch && len(users) > 1 {
			return nil, fmt.Errorf("email and DoD IDs match different users - please open up a help desk ticket")
		} else if emailMatch && !edipiMatch && len(users) == 1 {
			return nil, fmt.Errorf("there is an existing Okta account with that email - please update the DoD ID (EDIPI) in your Okta profile to match your registration DoD ID and try registering again")
		}

		if !emailMatch && edipiMatch && len(users) > 1 {
			return nil, fmt.Errorf("email and DoD IDs match different users - please open up a help desk ticket")
		} else if !emailMatch && edipiMatch && len(users) == 1 {
			return nil, fmt.Errorf("there is an existing Okta account with that DoD ID (EDIPI) - please update the email in your Okta profile to match your registration email and try registering again")
		}

		// if we get an email & edipi match on two different users and NOT an exact match, we need them to open a HDT
		if emailMatch && edipiMatch && len(users) > 1 {
			return nil, fmt.Errorf("there are multiple Okta accounts with that email and DoD ID - please open up a help desk ticket")
		}
	}

	profile := models.OktaProfile{
		FirstName:   oktaFirstName,
		LastName:    oktaLastName,
		Email:       oktaEmail,
		Login:       oktaEmail,
		MobilePhone: oktaPhone,
		CacEdipi:    *oktaEdipi,
	}
	oktaPayload := models.OktaUserPayload{
		Profile:  profile,
		GroupIds: []string{customerGroupID},
	}

	return models.CreateOktaUser(appCtx, provider, apiKey, oktaPayload)
}
