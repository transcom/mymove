package internalapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/appcontext"
	oktaop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/okta_profile"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
)

// GetOktaProfileHandler gets Okta Profile via GET /okta-profile
type GetOktaProfileHandler struct {
	handlers.HandlerConfig
}

// Handle retrieves a service member in the system belonging to the logged in user given service member ID
func (h GetOktaProfileHandler) Handle(params oktaop.ShowOktaInfoParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			oktaProfile := appCtx.Session().OktaSessionInfo

			oktaUserPayload := internalmessages.OktaUserPayload{
				Username:  oktaProfile.Username,
				Email:     oktaProfile.Email,
				FirstName: oktaProfile.FirstName,
				LastName:  oktaProfile.LastName,
				Edipi:     &oktaProfile.Edipi,
				Sub:       oktaProfile.Sub,
			}

			// this is going to check to see if the Okta profile data is present in the session
			if oktaUserPayload.Sub == "" {
				appCtx.Logger().Error("Session does not contain Okta values")
			}

			return oktaop.NewShowOktaInfoOK().WithPayload(&oktaUserPayload), nil
		})
}
