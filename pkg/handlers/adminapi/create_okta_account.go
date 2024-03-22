package adminapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/okta"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
)

func payloadToOktaAccountCreationModel(payload *adminmessages.CreateOktaAccount) models.OktaAccountCreationTemplate {
	return models.OktaAccountCreationTemplate{
		FirstName:   *payload.FirstName,
		LastName:    *payload.LastName,
		Login:       *payload.Login,
		Email:       *payload.Email,
		CacEdipi:    payload.CacEdipi,
		MobilePhone: *payload.MobilePhone,
		GsaID:       payload.GsaID,
	}
}

func CreateAccountOkta(appCtx appcontext.AppContext, params userop.CreateOktaAccountParams) (*http.Response, error) {

	// Payload to OktaAccountCreationTemplate
	oktaAccountInformation := payloadToOktaAccountCreationModel(params.CreateOktaAccountPayload)

	// Get Okta provider
	provider, err := okta.GetOktaProviderForRequest(params.HTTPRequest)
	if err != nil {
		appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf(" error getting okta provider - okta account not created")))
		return nil, err
	}

	// Setting viper so we can access the api key in the env vars
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// Okta api key
	apiKey := v.GetString(cli.OktaAPIKeyFlag)

	// Okta createUser url
	activate := "true"
	baseURL := provider.GetCreateAccountURL(activate)

	// Build okta profile body
	oktaProfileBody := models.OktaBodyProfile{
		FirstName:   oktaAccountInformation.FirstName,
		LastName:    oktaAccountInformation.LastName,
		Login:       oktaAccountInformation.Login,
		Email:       oktaAccountInformation.Email,
		MobilePhone: oktaAccountInformation.MobilePhone,
		CacEdipi:    oktaAccountInformation.CacEdipi,
		GsaID:       oktaAccountInformation.GsaID,
	}

	// Build Post request body
	body := models.OktaAccountCreationBody{
		Profile:  oktaProfileBody,
		GroupIds: []string{},
	}

	// Get Okta Office Group Id and add it to the request
	oktaOfficeGroupID := v.GetString(cli.OktaOfficeGroupIDFlag)
	body.GroupIds = append(body.GroupIds, oktaOfficeGroupID)

	// Marshall Post request body
	marshalledBody, err := json.Marshal(body)
	if err != nil {
		appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf(" error marshalling okta post request body - okta account not created")))
		return nil, err
	}

	// Create POST request
	userPostReq, err := http.NewRequest("POST", baseURL, bytes.NewReader(marshalledBody))
	if err != nil {
		appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf(" error creating okta post request - okta account not created")))
		return nil, err
	}

	// Set POST request header
	userPostReq.Header.Add("Authorization", "SSWS "+apiKey)
	userPostReq.Header.Add("Accept", "application/json")
	userPostReq.Header.Add("Content-Type", "application/json")

	// Execute POST request
	client := &http.Client{}
	res, err := client.Do(userPostReq)
	if err != nil {
		appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf(" error with okta account creation post request")))
		return res, err
	}

	return res, nil
}

// CreateOktaAccount Handler creates okta accounts
type CreateOktaAccount struct {
	handlers.HandlerConfig
}

// Handle creates an okta account
func (h CreateOktaAccount) Handle(params userop.CreateOktaAccountParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			res, err := CreateAccountOkta(appCtx, params)
			if err != nil {
				if res.StatusCode == http.StatusInternalServerError {
					appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf("okta returned internal server error")))
					return userop.NewCreateOktaAccountInternalServerError(), err
				}
				if res.StatusCode == http.StatusForbidden {
					appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf("okta returned status forbidden error")))
					return userop.NewCreateOktaAccountForbidden(), err
				}
				if res.StatusCode == http.StatusBadRequest {
					appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf("okta returned status bad request")))
					return userop.NewCreateOktaAccountBadRequest(), err
				}

				return userop.NewCreateOktaAccountInternalServerError(), err
			}

			// If account creation is successful
			if res.StatusCode == http.StatusOK {

				response, err := io.ReadAll(res.Body)
				if err != nil {
					appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf(" could not read response body")))
					return userop.NewCreateOktaAccountOK(), err
				}

				oktaAccountInfo := new(adminmessages.OktaAccountInfoResponse)
				err = json.Unmarshal(response, &oktaAccountInfo)
				if err != nil {
					appCtx.Logger().Error("could not unmarshal body", zap.Error(err))
					return userop.NewCreateOktaAccountOK(), err
				}

				defer res.Body.Close()

				appCtx.Logger().Info("Okta account successfully created")
				return userop.NewCreateOktaAccountOK().WithPayload(oktaAccountInfo), err
			}

			appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf(" unkown error")))
			return userop.NewCreateOktaAccountInternalServerError(), nil
		})
}
