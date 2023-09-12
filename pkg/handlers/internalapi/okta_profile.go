package internalapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/spf13/viper"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	oktaop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/okta_profile"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"go.uber.org/zap"
)

// GetOktaProfileHandler gets Okta Profile via GET /okta-profile
type GetOktaProfileHandler struct {
	handlers.HandlerConfig
}

// Handle retrieves a service member in the system belonging to the logged in user given service member ID
func (h GetOktaProfileHandler) Handle(params oktaop.ShowOktaInfoParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// getting okta id of user from session, to be used for api call
			oktaUserID := appCtx.Session().OktaSessionInfo.Sub

			// setting viper so we can access the api key in the env vars
			v := viper.New()
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv()
			apiKey := v.GetString(cli.OktaApiKeyFlag)

			// getting okta domain url for request
			provider, err := okta.GetOktaProviderForRequest(params.HTTPRequest)
			if err != nil {
				return nil, err
			}

			// need to pull this payload since it is wrapped in a profile object
			user := internalmessages.UpdateOktaUserPayload{}

			// getting the api call url from provider.go
			baseUrl := provider.GetUserURL(oktaUserID)

			body, _ := json.Marshal(user)

			// making HTTP request to Okta Users API to update user
			// this is done via a POST request for partial profile updates
			// https://developer.okta.com/docs/reference/api/users/#update-current-user-s-profile
			req, _ := http.NewRequest("GET", baseUrl, bytes.NewReader([]byte("")))
			h := req.Header
			h.Add("Authorization", "SSWS "+apiKey)
			h.Add("Accept", "application/json")
			h.Add("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				appCtx.Logger().Error("could not execute request", zap.Error(err))
			}

			body, err = io.ReadAll(resp.Body)
			if err != nil {
				appCtx.Logger().Error("could not read response body", zap.Error(err))
			}

			defer resp.Body.Close()

			err = json.Unmarshal(body, &user)
			if err != nil {
				appCtx.Logger().Error("could not unmarshal body", zap.Error(err))
			}

			oktaUserPayload := internalmessages.OktaUserPayload{
				Login:     user.Profile.Login,
				Email:     user.Profile.Email,
				FirstName: user.Profile.FirstName,
				LastName:  user.Profile.LastName,
				CacEdipi:  user.Profile.CacEdipi,
				Sub:       user.Profile.Sub,
			}

			return oktaop.NewShowOktaInfoOK().WithPayload(&oktaUserPayload), nil
		})
}

// GetOktaProfileHandler gets Okta Profile via GET /okta-profile
type UpdateOktaProfileHandler struct {
	handlers.HandlerConfig
}

type ProfileStruct struct {
	internalmessages.OktaUserPayload
}

// Handle implements okta_profile.UpdateOktaInfoHandler
// following the API call docs here: https://developer.okta.com/docs/reference/api/oidc/#client-authentication-methods
func (h UpdateOktaProfileHandler) Handle(params oktaop.UpdateOktaInfoParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// getting okta id of user from session, to be used for api call
			oktaUserID := appCtx.Session().OktaSessionInfo.Sub

			// setting viper so we can access the api key in the env vars
			v := viper.New()
			v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			v.AutomaticEnv()
			apiKey := v.GetString(cli.OktaApiKeyFlag)

			// getting okta domain url for post request
			provider, err := okta.GetOktaProviderForRequest(params.HTTPRequest)
			if err != nil {
				return nil, err
			}

			// payload is what is submitted from frontend, should contain
			// {email, login, firstName, lastName, cac_edipi}
			payload := params.UpdateOktaUserPayload

			// getting the api call url from provider.go
			baseUrl := provider.GetUserURL(oktaUserID)

			body, _ := json.Marshal(payload)

			// making HTTP request to Okta Users API to update user
			// this is done via a POST request for partial profile updates
			// https://developer.okta.com/docs/reference/api/users/#update-current-user-s-profile
			req, _ := http.NewRequest("POST", baseUrl, bytes.NewReader(body))
			h := req.Header
			h.Add("Authorization", "SSWS "+apiKey)
			h.Add("Accept", "application/json")
			h.Add("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				appCtx.Logger().Error("could not execute request", zap.Error(err))
			}

			body, err = io.ReadAll(resp.Body)
			if err != nil {
				appCtx.Logger().Error("could not read response body", zap.Error(err))
			}

			defer resp.Body.Close()

			err = json.Unmarshal(body, &payload)
			if err != nil {
				appCtx.Logger().Error("could not unmarshal body", zap.Error(err))
			}

			// when calling Okta, we have to have the body wrapped in a JSON profile object
			// here we will take the repsonse and convert it to a struct that doesn't have profile wrap
			oktaUserPayload := internalmessages.OktaUserPayload{
				Login:     payload.Profile.Login,
				Email:     payload.Profile.Email,
				FirstName: payload.Profile.FirstName,
				LastName:  payload.Profile.LastName,
				CacEdipi:  payload.Profile.CacEdipi,
				Sub:       oktaUserID,
			}

			// setting app context values with updated values so frontend can update
			appCtx.Session().OktaSessionInfo.Login = oktaUserPayload.Login
			appCtx.Session().OktaSessionInfo.Email = oktaUserPayload.Email
			appCtx.Session().OktaSessionInfo.FirstName = oktaUserPayload.FirstName
			appCtx.Session().OktaSessionInfo.LastName = oktaUserPayload.LastName
			appCtx.Session().OktaSessionInfo.Edipi = *oktaUserPayload.CacEdipi

			appCtx.Logger().Debug("API POST request return values to session")
			appCtx.Logger().Debug(appCtx.Session().OktaSessionInfo.Login)
			appCtx.Logger().Debug(appCtx.Session().OktaSessionInfo.Email)
			appCtx.Logger().Debug(appCtx.Session().OktaSessionInfo.FirstName)
			appCtx.Logger().Debug(appCtx.Session().OktaSessionInfo.LastName)
			appCtx.Logger().Debug(appCtx.Session().OktaSessionInfo.Edipi)

			return oktaop.NewUpdateOktaInfoOK().WithPayload(&oktaUserPayload), nil
		})
}
