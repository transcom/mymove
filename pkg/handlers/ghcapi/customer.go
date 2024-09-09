package ghcapi

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/cli"
	customercodeop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/customer"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func addressModelFromPayload(rawAddress *ghcmessages.Address) *models.Address {
	if rawAddress == nil {
		return nil
	}
	return &models.Address{
		StreetAddress1: *rawAddress.StreetAddress1,
		StreetAddress2: rawAddress.StreetAddress2,
		StreetAddress3: rawAddress.StreetAddress3,
		City:           *rawAddress.City,
		State:          *rawAddress.State,
		PostalCode:     *rawAddress.PostalCode,
		Country:        rawAddress.Country,
	}
}

// GetCustomerHandler fetches the information of a specific customer
type GetCustomerHandler struct {
	handlers.HandlerConfig
	services.CustomerFetcher
}

// Handle getting the information of a specific customer
func (h GetCustomerHandler) Handle(params customercodeop.GetCustomerParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			customerID, _ := uuid.FromString(params.CustomerID.String())
			customer, err := h.FetchCustomer(appCtx, customerID)
			if err != nil {
				appCtx.Logger().Error("Loading Customer Info", zap.Error(err))
				switch err {
				case sql.ErrNoRows:
					return customercodeop.NewGetCustomerNotFound(), err
				default:
					return customercodeop.NewGetCustomerInternalServerError(), err
				}
			}
			customerInfoPayload := payloads.Customer(customer)
			return customercodeop.NewGetCustomerOK().WithPayload(customerInfoPayload), nil
		})
}

type SearchCustomersHandler struct {
	handlers.HandlerConfig
	services.CustomerSearcher
}

func (h SearchCustomersHandler) Handle(params customercodeop.SearchCustomersParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			searchCustomersParams := services.SearchCustomersParams{
				DodID:        params.Body.DodID,
				CustomerName: params.Body.CustomerName,
				Page:         params.Body.Page,
				PerPage:      params.Body.PerPage,
				Sort:         params.Body.Sort,
				Order:        params.Body.Order,
			}

			customers, totalCount, err := h.CustomerSearcher.SearchCustomers(appCtx, &searchCustomersParams)

			if err != nil {
				appCtx.Logger().Error("Error searching for customer", zap.Error(err))
				switch err.(type) {
				case apperror.ForbiddenError:
					return customercodeop.NewSearchCustomersForbidden(), err
				default:
					return customercodeop.NewSearchCustomersInternalServerError(), err
				}
			}

			searchCustomers := payloads.SearchCustomers(customers)
			payload := &ghcmessages.SearchCustomersResult{
				Page:            1,
				PerPage:         20,
				TotalCount:      int64(totalCount),
				SearchCustomers: *searchCustomers,
			}
			return customercodeop.NewSearchCustomersOK().WithPayload(payload), nil
		})
}

// UpdateCustomerHandler updates a customer via PATCH /customer/{customerId}
type UpdateCustomerHandler struct {
	handlers.HandlerConfig
	customerUpdater services.CustomerUpdater
}

// Handle updates a customer from a request payload
func (h UpdateCustomerHandler) Handle(params customercodeop.UpdateCustomerParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			customerID, err := uuid.FromString(params.CustomerID.String())
			if err != nil {
				appCtx.Logger().Error("unable to parse customer id param to uuid", zap.Error(err))
				return customercodeop.NewUpdateCustomerBadRequest(), err
			}

			newCustomer := payloads.CustomerToServiceMember(*params.Body)
			newCustomer.ID = customerID

			updatedCustomer, err := h.customerUpdater.UpdateCustomer(appCtx, params.IfMatch, newCustomer)

			if err != nil {
				appCtx.Logger().Error("error updating customer", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return customercodeop.NewGetCustomerNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return customercodeop.NewUpdateCustomerUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return customercodeop.NewUpdateCustomerPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return customercodeop.NewUpdateCustomerInternalServerError(), err
				}
			}

			customerPayload := payloads.Customer(updatedCustomer)

			return customercodeop.NewUpdateCustomerOK().WithPayload(customerPayload), nil
		})
}

type CreateCustomerWithOktaOptionHandler struct {
	handlers.HandlerConfig
}

// Handle creates a customer/serviceMember from a request payload
func (h CreateCustomerWithOktaOptionHandler) Handle(params customercodeop.CreateCustomerWithOktaOptionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.Body
			var err error
			var serviceMembers []models.ServiceMember
			var dodidUniqueFeatureFlag bool

			// evaluating feature flag to see if we need to check if the DODID exists already
			featureFlagName := "dodid_unique"
			flag, err := h.FeatureFlagFetcher().GetBooleanFlagForUser(params.HTTPRequest.Context(), appCtx, featureFlagName, map[string]string{})
			if err != nil {
				appCtx.Logger().Error("Error fetching dodid_unique feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
				dodidUniqueFeatureFlag = false
			} else {
				dodidUniqueFeatureFlag = flag.Match
			}

			if dodidUniqueFeatureFlag {
				query := `SELECT service_members.edipi
								FROM service_members
								WHERE service_members.edipi = $1`
				err := appCtx.DB().RawQuery(query, payload.Edipi).All(&serviceMembers)
				if err != nil {
					errorMsg := apperror.NewBadDataError("error when checking for existing service member")
					payload := payloadForValidationError("Unable to create a customer", errorMsg.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return customercodeop.NewCreateCustomerWithOktaOptionUnprocessableEntity().WithPayload(payload), errorMsg
				} else if len(serviceMembers) > 0 {
					errorMsg := apperror.NewConflictError(h.GetTraceIDFromRequest(params.HTTPRequest), "Service member with this DODID already exists. Please use a different DODID number.")
					payload := payloadForValidationError("Unable to create a customer", errorMsg.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return customercodeop.NewCreateCustomerWithOktaOptionUnprocessableEntity().WithPayload(payload), errorMsg
				}
			}

			// Endpoint specific EDIPI and EMPLID check
			// The following validation currently is only intended for the customer creation
			// conducted by an office user such as the Service Counselor
			if payload.Affiliation != nil && *payload.Affiliation == ghcmessages.AffiliationCOASTGUARD {
				// EMPLID cannot be null
				if payload.Emplid == nil {
					errorMsg := apperror.NewConflictError(h.GetTraceIDFromRequest(params.HTTPRequest), "Service members from the Coast Guard require an EMPLID for creation.")
					payload := payloadForValidationError("Unable to create a customer", errorMsg.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return customercodeop.NewCreateCustomerWithOktaOptionUnprocessableEntity().WithPayload(payload), errorMsg
				}
			}

			var newServiceMember models.ServiceMember
			var backupContact models.BackupContact

			email := payload.PersonalEmail
			if email == "" {
				badDataError := apperror.NewBadDataError("missing personal email")
				payload := payloadForValidationError("Unable to create a customer", badDataError.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
				return customercodeop.NewCreateCustomerWithOktaOptionUnprocessableEntity().WithPayload(payload), badDataError
			}

			// declaring okta values outside of if statements so we can use them later
			var oktaSub string
			oktaUser := &models.CreatedOktaUser{}

			// if the office user checked "yes", then we will create an okta account for the user
			// this will add the user to the okta customer application and send an activation email
			if payload.CreateOktaAccount {
				var oktaErr error
				oktaUser, oktaErr = createOktaProfile(appCtx, params)
				if oktaErr != nil {
					appCtx.Logger().Error("error creating okta profile", zap.Error(oktaErr))
					return customercodeop.NewCreateCustomerWithOktaOptionBadRequest(), oktaErr
				}
				oktaSub = oktaUser.ID
			}

			transactionError := appCtx.NewTransaction(func(_ appcontext.AppContext) error {
				// if the office user checked "no" to indicate the customer does NOT have a CAC, set cac_validated
				// to true so that the customer can log in without having to authenticate with a CAC
				var cacValidated = false
				if !payload.CacUser {
					cacValidated = true
				}
				var verrs *validate.Errors
				// creating a user and populating okta values (for now these can be null)
				user, userErr := models.CreateUser(appCtx.DB(), oktaSub, email)
				if userErr != nil {
					appCtx.Logger().Error("error creating user", zap.Error(err))
					return userErr
				}

				// now we will take all the data we have and build the service member
				userID := user.ID
				residentialAddress := addressModelFromPayload(&payload.ResidentialAddress.Address)
				backupMailingAddress := addressModelFromPayload(&payload.BackupMailingAddress.Address)

				// Create a new serviceMember using the userID
				newServiceMember = models.ServiceMember{
					UserID:               userID,
					Edipi:                &payload.Edipi,
					Emplid:               payload.Emplid,
					Affiliation:          (*models.ServiceMemberAffiliation)(payload.Affiliation),
					FirstName:            &payload.FirstName,
					MiddleName:           payload.MiddleName,
					LastName:             &payload.LastName,
					Suffix:               payload.Suffix,
					Telephone:            payload.Telephone,
					SecondaryTelephone:   payload.SecondaryTelephone,
					PersonalEmail:        &payload.PersonalEmail,
					PhoneIsPreferred:     &payload.PhoneIsPreferred,
					EmailIsPreferred:     &payload.EmailIsPreferred,
					ResidentialAddress:   residentialAddress,
					BackupMailingAddress: backupMailingAddress,
					CacValidated:         cacValidated,
				}

				// create the service member and save to the db
				smVerrs, smErr := models.SaveServiceMember(appCtx, &newServiceMember)
				if smVerrs.HasAny() || smErr != nil {
					appCtx.Logger().Error("error creating service member", zap.Error(smErr))
					return smErr
				}

				// creating backup contact associated with service member since this is done separately
				// default permission of EDIT since we want them to be able to change this info
				defaultPermission := models.BackupContactPermissionEDIT
				backupContact, verrs, err = newServiceMember.CreateBackupContact(appCtx.DB(),
					*payload.BackupContact.Name,
					*payload.BackupContact.Email,
					payload.BackupContact.Phone,
					models.BackupContactPermission(defaultPermission))
				if err != nil || verrs.HasAny() {
					appCtx.Logger().Error("error creating backup contact", zap.Error(err))
					return err
				}
				return nil
			})

			if transactionError != nil {
				switch transactionError.(type) {
				case *pq.Error:
					// handle duplicate key error for emplid
					return customercodeop.NewCreateCustomerWithOktaOptionConflict(), transactionError
				default:
					return customercodeop.NewCreateCustomerWithOktaOptionBadRequest(), transactionError
				}
			}

			// covering error returns
			if err != nil {
				appCtx.Logger().Error("error creating customer", zap.Error(err))
				switch err.(type) {
				case apperror.NotFoundError:
					return customercodeop.NewCreateCustomerWithOktaOptionNotFound(), err
				case apperror.InvalidInputError:
					payload := payloadForValidationError("Unable to complete request", err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), validate.NewErrors())
					return customercodeop.NewCreateCustomerWithOktaOptionUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return customercodeop.NewCreateCustomerWithOktaOptionPreconditionFailed().WithPayload(&ghcmessages.Error{Message: handlers.FmtString(err.Error())}), err
				default:
					return customercodeop.NewUpdateCustomerInternalServerError(), err
				}
			}

			customerPayload := payloads.CreatedCustomer(&newServiceMember, oktaUser, &backupContact)

			return customercodeop.NewCreateCustomerWithOktaOptionOK().WithPayload(customerPayload), nil
		})
}

// createOktaProfile sends a request to the Okta Users API
// this creates a user in Okta assigned to the customer group (allowing access to the customer application)
func createOktaProfile(appCtx appcontext.AppContext, params customercodeop.CreateCustomerWithOktaOptionParams) (*models.CreatedOktaUser, error) {
	// setting viper so we can access the api key in the env vars
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	apiKey := v.GetString(cli.OktaAPIKeyFlag)
	customerGroupID := v.GetString(cli.OktaCustomerGroupIDFlag)

	// taking all the data that we'll need for the okta profile creation
	payload := params.Body
	oktaEmail := payload.PersonalEmail
	oktaFirstName := payload.FirstName
	oktaLastName := payload.LastName
	oktaPhone := payload.Telephone

	// Creating the Profile struct
	profile := models.Profile{
		FirstName:   oktaFirstName,
		LastName:    oktaLastName,
		Email:       oktaEmail,
		Login:       oktaEmail,
		MobilePhone: *oktaPhone,
	}

	// Creating the OktaUserPayload struct
	oktaPayload := models.OktaUserPayload{
		Profile:  profile,
		GroupIds: []string{customerGroupID},
	}

	// getting okta domain url for request
	provider, err := okta.GetOktaProviderForRequest(params.HTTPRequest)
	if err != nil {
		return nil, err
	}

	// getting the api call url from provider.go
	activate := "true"
	baseURL := provider.GetCreateUserURL(activate)

	body, err := json.Marshal(oktaPayload)
	if err != nil {
		appCtx.Logger().Error("error marshaling payload", zap.Error(err))
		return nil, err
	}

	// making HTTP request to Okta Users API to create a user
	// this is done via a POST request for creating a user that sends an activation email (when activate=true)
	// https://developer.okta.com/docs/reference/api/users/#create-user-without-credentials
	req, err := http.NewRequest("POST", baseURL, bytes.NewReader(body))
	if err != nil {
		appCtx.Logger().Error("could not execute request", zap.Error(err))
		return nil, err
	}
	h := req.Header
	h.Add("Authorization", "SSWS "+apiKey)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/json")

	// now let the client send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		appCtx.Logger().Error("could not execute request", zap.Error(err))
		return nil, err
	}

	// if all is well, should have a 200 response
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		appCtx.Logger().Error("could not read response body", zap.Error(err))
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusInternalServerError {
			return nil, err
		}
		if resp.StatusCode == http.StatusBadRequest {
			return nil, err
		}
		if resp.StatusCode == http.StatusForbidden {
			return nil, err
		}
	}

	// now we will take the response and parse it into our Go struct
	user := models.CreatedOktaUser{}
	err = json.Unmarshal(response, &user)
	if err != nil {
		appCtx.Logger().Error("could not unmarshal body", zap.Error(err))
		return nil, err
	}

	defer resp.Body.Close()

	return &user, nil
}
