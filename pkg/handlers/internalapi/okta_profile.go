package internalapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

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

			oktaUser := appCtx.Session().OktaSessionInfo

			oktaUserPayload := internalmessages.OktaUserPayload{
				Username:  oktaUser.Username,
				Email:     oktaUser.Email,
				FirstName: oktaUser.FirstName,
				LastName:  oktaUser.LastName,
				Edipi:     &oktaUser.Edipi,
				Sub:       oktaUser.Sub,
			}

			// this is going to check to see if the Okta profile data is present in the session
			if oktaUserPayload.Sub == "" {
				appCtx.Logger().Error("Session does not contain Okta values")
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
// following the docs here: https://developer.okta.com/docs/reference/api/oidc/#client-authentication-methods
func (h UpdateOktaProfileHandler) Handle(params oktaop.UpdateOktaInfoParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// getting okta id of user from session, to be used for api call
			// oktaUserID := appCtx.Session().OktaSessionInfo.Sub
			v := viper.New()
			apiKey := v.GetString(cli.OktaApiKeyFlag)
			a := os.Getenv("OKTA_CUSTOMER_SECRET_KEY")
			appCtx.Logger().Debug(a)

			// getting okta domain url for post request
			provider, err := okta.GetOktaProviderForRequest(params.HTTPRequest)
			if err != nil {
				return nil, err
			}

			// payload is what is submitted from FE, should contain
			// {email, username, first_name, last_naame, edipi, sub}
			payload := params.UpdateOktaUserPayload

			// getting the api call url from provider.go
			baseUrl := provider.GetUserURL()

			body, _ := json.Marshal(payload)

			// making HTTP request to Okta Users API
			req, _ := http.NewRequest("GET", baseUrl, bytes.NewReader([]byte("")))
			h := req.Header
			h.Add("Authorization", "Bearer "+appCtx.Session().AccessToken)
			h.Add("Accept", "application/json; okta-version=1.0.0")
			h.Add("scope", apiKey)

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
			err = json.Unmarshal(body, payload)
			if err != nil {
				appCtx.Logger().Error("could not unmarshal body", zap.Error(err))
			}

			return oktaop.NewUpdateOktaInfoOK().WithPayload(nil), nil
		})
}
