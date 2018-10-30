package dpsapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/gen/dpsapi/dpsoperations/dps"
	"github.com/transcom/mymove/pkg/gen/dpsmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"go.uber.org/zap"
)

// GetUserHandler returns user information given an encrypted token
type GetUserHandler struct {
	handlers.HandlerContext
}

// Handle returns user information given an encrypted token
func (h GetUserHandler) Handle(params dps.GetUserParams) middleware.Responder {
	token := params.Token
	userID, err := dpsauth.CookieToUserID(token)
	if err != nil {
		h.Logger().Error("Extracting user ID from token", zap.Error(err))
		return dps.NewGetUserInternalServerError()
	}

	return dps.NewGetUserOK().WithPayload(getPayload(userID))
}

func getPayload(userID string) *dpsmessages.AuthenticationUserPayload {
	// TODO: Add real data
	affiliation := dpsmessages.AffiliationArmy
	middleName := "M"
	suffix := "III"
	telephone := "(555) 555-5555"
	payload := dpsmessages.AuthenticationUserPayload{
		Affiliation:          &affiliation,
		Email:                "test@example.com",
		FirstName:            "Jane",
		MiddleName:           &middleName,
		LastName:             "Doe",
		Suffix:               &suffix,
		LoginGovID:           strfmt.UUID(userID),
		SocialSecurityNumber: "666555555",
		Telephone:            &telephone,
	}
	return &payload
}
