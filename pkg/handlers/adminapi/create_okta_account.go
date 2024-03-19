package adminapi

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Login:       payload.Login,
		Email:       payload.Email,
		CacEdipi:    payload.CacEdipi,
		MobilePhone: payload.MobilePhone,
		GsaID:       payload.GsaID,
		GroupIds:    payload.GroupID,
	}
}

// CreateOktaAccount Handler creats okta accounts
type CreateOktaAccount struct {
	handlers.HandlerConfig
}

// Handle creates an okta account
func (h CreateOktaAccount) Handle(params userop.CreateOktaAccountParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// Payload to OktaAccountCreationTemplate
			oktaAccountInformation := payloadToOktaAccountCreationModel(params.CreateOktaAccountPayload)

			// Get Okta provider
			provider, err := okta.GetOktaProviderForRequest(params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf(" error getting okta provider - okta account not created")))
				return userop.NewCreateOktaAccountInternalServerError(), err
			}

			// Get the Okta Domain from the Okta provider
			oktaDomain := provider.GetOrgURL()

			// setting viper so we can access the api key in the env vars
			v := viper.New()
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv()

			// Okta api key
			apiKey := v.GetString(cli.OktaAPIKeyFlag)

			// Okta getUser url
			baseURL := oktaDomain + "/api/v1/users"

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
				GroupIds: oktaAccountInformation.GroupIds,
			}

			// Marshall Post request body
			marshalledBody, err := json.Marshal(body)
			if err != nil {
				appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf(" error marshalling okta post request body - okta account not created")))
				return userop.NewCreateOktaAccountInternalServerError(), err
			}

			// Create POST request
			userPostReq, err := http.NewRequest("POST", baseURL, bytes.NewReader(marshalledBody))
			if err != nil {
				appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf(" error creating okta post request - okta account not created")))
				return userop.NewCreateOktaAccountInternalServerError(), err
			}

			// Add url params
			urlParams := userPostReq.URL.Query()
			urlParams.Add("activate", "false")

			// Set POST request header
			userPostReq.Header.Add("Authorization", "SSWS "+apiKey)
			userPostReq.Header.Add("Accept", "application/json")
			userPostReq.Header.Add("Content-Type", "application/json")

			// // Execute POST request
			client := &http.Client{}
			res, err := client.Do(userPostReq)
			if err != nil {
				appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf(" error with okta account creation post request")))
				return userop.NewCreateOktaAccountInternalServerError(), err
			}

			// If account creation is success
			if res.StatusCode == http.StatusOK {
				appCtx.Logger().Info("Okta account successfully created")
				return userop.NewCreateOktaAccountOK().WithPayload(params.CreateOktaAccountPayload), err
			}

			if res.StatusCode == http.StatusInternalServerError {
				appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf("okta returned internal server error")))
			}
			if res.StatusCode == http.StatusForbidden {
				appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf("okta returned status forbidden error")))
			}
			if res.StatusCode == http.StatusBadRequest {
				appCtx.Logger().Error("oktaAccountCreator Error", zap.Error(fmt.Errorf("okta returned status bad request")))
			}

			return userop.NewCreateOktaAccountInternalServerError(), err
		})
}
