package internalapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/appcontext"
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

// Handle implements okta_profile.ShowOktaInfoHandler.
func (h UpdateOktaProfileHandler) Handle(params oktaop.UpdateOktaInfoParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			oktaUserID := appCtx.Session().OktaSessionInfo.Sub
			provider, err := okta.GetOktaProviderForRequest(params.HTTPRequest)
			if err != nil {
				return nil, err
			}

			payload := params.UpdateOktaUserPayload

			url := provider.GetUserURL(oktaUserID)
			body, _ := json.Marshal(payload)

			req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
			h := req.Header
			h.Add("Authorization", "Bearer "+appCtx.Session().AccessToken)
			h.Add("Accept", "application/json")

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

			return oktaop.NewUpdateOktaInfoOK().WithPayload(payload), nil
		})
}
