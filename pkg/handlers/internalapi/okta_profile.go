package internalapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	oktaop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/okta_profile"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
)

// GetOktaProfileHandler gets Okta Profile via GET /okta-profile
type GetOktaProfileHandler struct {
	handlers.HandlerConfig
}

// Handle performs a GET request from Okta API, returns values in profile object from response
// Could  not use data from sessions since access token data does not change when profile is updated
func (h GetOktaProfileHandler) Handle(params oktaop.ShowOktaInfoParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// if the "Local Sign In" is clicked we are going to send back dummy values
			sess := appCtx.Session()
			if sess.IDToken == "devlocal" {
				oktaUserPayload := internalmessages.OktaUserProfileData{
					Login:     "devlocal",
					Email:     "devlocal",
					FirstName: "devlocal",
					LastName:  "devlocal",
					CacEdipi:  nil,
					Sub:       "devlocal",
				}
				return oktaop.NewShowOktaInfoOK().WithPayload(&oktaUserPayload), nil
			}

			// getting okta id of user from session, to be used for api call
			oktaUserID := appCtx.Session().OktaSessionInfo.Sub

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

			// need to pull this payload since it is wrapped in a profile object so resp
			// body can populate accurately
			user := internalmessages.UpdateOktaUserProfileData{}

			// getting the api call url from provider.go
			baseURL := provider.GetUserURL(oktaUserID)

			// making HTTP request to Okta Users API to update user
			// this is done via a POST request for partial profile updates
			// https://developer.okta.com/docs/reference/api/users/#update-current-user-s-profile
			req, _ := http.NewRequest("GET", baseURL, bytes.NewReader([]byte("")))
			h := req.Header
			h.Add("Authorization", "SSWS "+apiKey)
			h.Add("Accept", "application/json")
			h.Add("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				appCtx.Logger().Error("could not execute request", zap.Error(err))
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				appCtx.Logger().Error("could not read response body", zap.Error(err))
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			if resp.StatusCode != http.StatusOK {
				if resp.StatusCode == http.StatusInternalServerError {
					return oktaop.NewShowOktaInfoInternalServerError(), err
				}
				if resp.StatusCode == http.StatusBadRequest {
					return oktaop.NewShowOktaInfoBadRequest(), err
				}
				if resp.StatusCode == http.StatusForbidden {
					return oktaop.NewShowOktaInfoForbidden(), err
				}
			}

			defer resp.Body.Close()

			err = json.Unmarshal(body, &user)
			if err != nil {
				appCtx.Logger().Error("could not unmarshal body", zap.Error(err))
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			// the return value has to be of type OktaUserPayload
			// our initial objet was of type UpdateOktaUserPayload, so needs to be changed
			// OktaUserPayload is not wrapped in a profile object
			oktaUserPayload := internalmessages.OktaUserProfileData{
				Login:     user.Profile.Login,
				Email:     user.Profile.Email,
				FirstName: user.Profile.FirstName,
				LastName:  user.Profile.LastName,
				CacEdipi:  user.Profile.CacEdipi,
				Sub:       user.Profile.Sub,
			}
			appCtx.Logger().Info("getting okta profile", zap.Any("okta profile", oktaUserPayload))

			return oktaop.NewShowOktaInfoOK().WithPayload(&oktaUserPayload), nil
		})
}

// GetOktaProfileHandler gets Okta Profile via GET /okta-profile
type UpdateOktaProfileHandler struct {
	handlers.HandlerConfig
}

type ErrorResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorSummary string `json:"errorSummary"`
	ErrorLink    string `json:"errorLink"`
	ErrorID      string `json:"errorId"`
	ErrorCauses  []struct {
		ErrorSummary string `json:"errorSummary"`
	} `json:"errorCauses"`
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
			apiKey := v.GetString(cli.OktaAPIKeyFlag)

			// getting okta domain url for post request
			provider, err := okta.GetOktaProviderForRequest(params.HTTPRequest)
			if err != nil {
				return nil, err
			}

			// payload is what is submitted from frontend, should contain
			// {email, login, firstName, lastName, cac_edipi}
			payload := params.UpdateOktaUserProfileData

			// getting the api call url from provider.go
			baseURL := provider.GetUserURL(oktaUserID)

			body, _ := json.Marshal(payload)

			// making HTTP request to Okta Users API to update user
			// this is done via a POST request for partial profile updates
			// https://developer.okta.com/docs/reference/api/users/#update-current-user-s-profile
			req, _ := http.NewRequest("POST", baseURL, bytes.NewReader(body))
			h := req.Header
			h.Add("Authorization", "SSWS "+apiKey)
			h.Add("Accept", "application/json")
			h.Add("Content-Type", "application/json")

			client := &http.Client{}

			resp, err := client.Do(req)
			if err != nil {
				appCtx.Logger().Error("could not execute request", zap.Error(err))
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			body, err = io.ReadAll(resp.Body)
			if err != nil {
				appCtx.Logger().Error("could not read response body", zap.Any("returned status", resp.Status))
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			defer resp.Body.Close()

			// we are going to check for an okta error response
			var response ErrorResponse
			err = json.Unmarshal(body, &response)
			if err != nil {
				appCtx.Logger().Error("could not unmarshal body", zap.Error(err))
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}
			// if there's an error code, we will see what the error will be and display it to the user
			if response.ErrorCode != "" {
				errorSummary := response.ErrorSummary
				fieldError := strings.TrimSpace(strings.SplitN(errorSummary, "Api validation failed:", 2)[1])
				// changing json fields to something more easier to read
				switch fieldError {
				case "cac_edipi":
					fieldError = "Dod ID | EDIPI"
				case "firstName":
					fieldError = "Last Name"
				case "lastName":
					fieldError = "First Name"
				case "email":
					fieldError = "Email"
				default:
					fieldError = ""
				}
				// extracting the part of the response that we want
				errorDescription := response.ErrorCauses[0].ErrorSummary
				extractedDescription := strings.TrimSpace(strings.SplitN(errorDescription, ":", 2)[1])
				errPayload := internalmessages.ValidationError{}
				errPayload.Detail = handlers.FmtString(string(fieldError + ": " + extractedDescription))
				return oktaop.NewUpdateOktaInfoUnprocessableEntity().WithPayload(errPayload), err
			}

			if resp.StatusCode != http.StatusOK {
				if resp.StatusCode == http.StatusInternalServerError {
					return oktaop.NewShowOktaInfoInternalServerError(), err
				}
				if resp.StatusCode == http.StatusBadRequest {
					return oktaop.NewUpdateOktaInfoBadRequest(), err
				}
				if resp.StatusCode == http.StatusForbidden {
					return oktaop.NewShowOktaInfoForbidden(), err
				}
			}

			err = json.Unmarshal(body, &payload)
			if err != nil {
				appCtx.Logger().Error("could not unmarshal body", zap.Error(err))
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			// when calling Okta, we have to have the body wrapped in a JSON profile object
			// here we will take the repsonse and convert it to a struct that doesn't have profile wrap
			oktaUserPayload := internalmessages.OktaUserProfileData{
				Login:     payload.Profile.Login,
				Email:     payload.Profile.Email,
				FirstName: payload.Profile.FirstName,
				LastName:  payload.Profile.LastName,
				CacEdipi:  payload.Profile.CacEdipi,
				Sub:       oktaUserID,
			}
			appCtx.Logger().Info("updated okta profile", zap.Any("okta profile", oktaUserPayload))

			// setting app context values with updated values so frontend can update
			appCtx.Session().OktaSessionInfo.Login = oktaUserPayload.Login
			appCtx.Session().OktaSessionInfo.Email = oktaUserPayload.Email
			appCtx.Session().OktaSessionInfo.FirstName = oktaUserPayload.FirstName
			appCtx.Session().OktaSessionInfo.LastName = oktaUserPayload.LastName
			appCtx.Session().OktaSessionInfo.Edipi = *oktaUserPayload.CacEdipi

			return oktaop.NewUpdateOktaInfoOK().WithPayload(&oktaUserPayload), nil
		})
}
