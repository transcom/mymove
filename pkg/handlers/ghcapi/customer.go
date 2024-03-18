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

// Handle updates a customer from a request payload
func (h CreateCustomerWithOktaOptionHandler) Handle(params customercodeop.CreateCustomerWithOktaOptionParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			payload := params.Body
			email := payload.PersonalEmail

			// delcaring okta values outside of if statements so we can use them later
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

			// creating a user and populating okta values (for now these can be null)
			user, err := models.CreateUser(appCtx.DB(), oktaSub, *email)
			if err != nil {
				appCtx.Logger().Error("error creating user", zap.Error(err))
				return customercodeop.NewCreateCustomerWithOktaOptionBadRequest(), err
			}

			// now we will take all the data we have and build the service member
			userID := user.ID
			residentialAddress := addressModelFromPayload(&payload.ResidentialAddress.Address)
			backupMailingAddress := addressModelFromPayload(&payload.BackupMailingAddress.Address)

			// Create a new serviceMember using the userID
			newServiceMember := models.ServiceMember{
				UserID:               userID,
				Edipi:                payload.Edipi,
				Affiliation:          (*models.ServiceMemberAffiliation)(payload.Affiliation),
				FirstName:            &payload.FirstName,
				MiddleName:           payload.MiddleName,
				LastName:             &payload.LastName,
				Suffix:               payload.Suffix,
				Telephone:            payload.Telephone,
				SecondaryTelephone:   payload.SecondaryTelephone,
				PersonalEmail:        payload.PersonalEmail,
				PhoneIsPreferred:     &payload.PhoneIsPreferred,
				EmailIsPreferred:     &payload.EmailIsPreferred,
				ResidentialAddress:   residentialAddress,
				BackupMailingAddress: backupMailingAddress,
			}
			// create the service member and save to the db
			smVerrs, err := models.SaveServiceMember(appCtx, &newServiceMember)
			if smVerrs.HasAny() || err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			// creating backup contact associated with service member
			// default permission of EDIT since we want them to be able to change this info
			defaultPermission := models.BackupContactPermissionEDIT
			backupContact, verrs, err := newServiceMember.CreateBackupContact(appCtx.DB(),
				*payload.BackupContact.Name,
				*payload.BackupContact.Email,
				payload.BackupContact.Phone,
				models.BackupContactPermission(defaultPermission))
			if err != nil || verrs.HasAny() {
				return handlers.ResponseForVErrors(appCtx.Logger(), verrs, err), err
			}

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

func createOktaProfile(appCtx appcontext.AppContext, params customercodeop.CreateCustomerWithOktaOptionParams) (*models.CreatedOktaUser, error) {

	payload := params.Body
	oktaEmail := payload.PersonalEmail
	oktaFirstName := payload.FirstName
	oktaLastName := payload.LastName
	oktaPhone := payload.Telephone

	// Creating the Profile struct
	profile := models.Profile{
		FirstName:   oktaFirstName,
		LastName:    oktaLastName,
		Email:       *oktaEmail,
		Login:       *oktaEmail,
		MobilePhone: *oktaPhone,
	}

	// Creating the OktaUserPayload struct
	oktaPayload := models.OktaUserPayload{
		Profile:  profile,
		GroupIds: []string{"00g3ja8t0dwKG8Mmi0k6"},
	}

	// setting viper so we can access the api key in the env vars
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	apiKey := v.GetString(cli.OktaAPIKeyFlag)

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
	// this is done via a POST request for creating a user that sends an activation email
	// https://developer.okta.com/docs/reference/api/users/#create-user-without-credentials
	req, _ := http.NewRequest("POST", baseURL, bytes.NewReader(body))
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
