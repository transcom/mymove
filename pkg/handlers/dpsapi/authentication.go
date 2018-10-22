package dpsapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/gen/dpsapi/dpsoperations/dps"
	"github.com/transcom/mymove/pkg/gen/dpsmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

// GetUserHandler returns user information given an encrypted token
type GetUserHandler struct {
	handlers.HandlerContext
}

var affiliationMap = map[models.ServiceMemberAffiliation]dpsmessages.Affiliation{
	models.AffiliationARMY:       dpsmessages.AffiliationArmy,
	models.AffiliationNAVY:       dpsmessages.AffiliationNavy,
	models.AffiliationMARINES:    dpsmessages.AffiliationMarines,
	models.AffiliationAIRFORCE:   dpsmessages.AffiliationAirForce,
	models.AffiliationCOASTGUARD: dpsmessages.AffiliationCoastGuard,
}

// Handle returns user information given an encrypted token
func (h GetUserHandler) Handle(params dps.GetUserParams) middleware.Responder {
	token := params.Token
	smID, err := dpsauth.CookieToServiceMemberID(token)
	if err != nil {
		h.Logger().Error("Extracting user ID from token", zap.Error(err))
		return dps.NewGetUserInternalServerError()
	}

	payload, err := getPayload(h.DB(), smID)
	if err != nil {
		h.Logger().Error("Fetching user data from user ID", zap.Error(err))
		return dps.NewGetUserInternalServerError()
	}

	return dps.NewGetUserOK().WithPayload(payload)
}

func getPayload(db *pop.Connection, smID string) (*dpsmessages.AuthenticationUserPayload, error) {
	id, err := uuid.FromString(smID)
	if err != nil {
		return nil, err
	}

	sm, err := models.FetchServiceMember(db, id)
	if err != nil {
		return nil, errors.Wrap(err, "Fetching service member")
	}
	user, err := models.GetUser(db, sm.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "Fetching user")
	}

	var affiliation *dpsmessages.Affiliation
	if sm.Affiliation != nil {
		dpsaffiliation := affiliationMap[*sm.Affiliation]
		affiliation = &dpsaffiliation
	}

	ssn, err := getSSNFromIWS(sm.Edipi)
	if err != nil {
		return nil, errors.Wrap(err, "Getting SSN from IWS using EDIPI")
	}

	if sm.FirstName == nil || sm.LastName == nil {
		return nil, errors.New("Service member is missing first and/or last name")
	}

	payload := dpsmessages.AuthenticationUserPayload{
		Affiliation:          affiliation,
		Email:                user.LoginGovEmail,
		FirstName:            *sm.FirstName,
		MiddleName:           sm.MiddleName,
		LastName:             *sm.LastName,
		Suffix:               sm.Suffix,
		LoginGovID:           strfmt.UUID(user.LoginGovUUID.String()),
		SocialSecurityNumber: strfmt.SSN(ssn),
		Telephone:            sm.Telephone,
	}
	return &payload, nil
}

func getSSNFromIWS(edipi *string) (string, error) {
	if edipi == nil {
		return "", errors.New("Service member is missing EDIPI")
	}

	// TODO

	return "", nil
}
