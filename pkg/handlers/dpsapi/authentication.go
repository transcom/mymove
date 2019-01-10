package dpsapi

import (
	"fmt"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/gen/dpsapi/dpsoperations/dps"
	"github.com/transcom/mymove/pkg/gen/dpsmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/iws"
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
	clientCert := authentication.ClientCertFromRequestContext(params.HTTPRequest)
	if clientCert == nil || !clientCert.AllowDpsAuthAPI {
		h.Logger().Info("Client certificate is not authorized to access this API")
		return dps.NewGetUserUnauthorized()
	}

	token := params.Token
	loginGovID, err := dpsauth.CookieToLoginGovID(token)
	if err != nil {
		h.Logger().Error("Extracting user ID from token", zap.Error(err))

		switch err.(type) {
		case *dpsauth.ErrInvalidCookie:
			return dps.NewGetUserUnprocessableEntity()
		}
		return dps.NewGetUserInternalServerError()
	}

	payload, err := getPayload(h.DB(), loginGovID, h.IWSRealTimeBrokerService())
	if err != nil {
		switch e := err.(type) {
		case *errUserMissingData:
			h.Logger().Error("Fetching user data from user ID", zap.Error(err), zap.String("user", e.userID.String()))
		default:
			h.Logger().Error("Fetching user data from user ID", zap.Error(err))
		}

		return dps.NewGetUserInternalServerError()
	}

	return dps.NewGetUserOK().WithPayload(payload)
}

func getPayload(db *pop.Connection, loginGovID string, rbs iws.RealTimeBrokerService) (*dpsmessages.AuthenticationUserPayload, error) {
	userIdentity, err := models.FetchUserIdentity(db, loginGovID)
	if err != nil {
		return nil, errors.Wrap(err, "Fetching user identity")
	}

	if userIdentity.ServiceMemberID == nil {
		return nil, &errUserMissingData{
			userID:     userIdentity.ID,
			errMessage: fmt.Sprintf("User %s is missing a service member ID", userIdentity.ID.String()),
		}
	}

	sm, err := models.FetchServiceMember(db, *userIdentity.ServiceMemberID)
	if err != nil {
		return nil, errors.Wrap(err, "Fetching service member")
	}

	var affiliation *dpsmessages.Affiliation
	if sm.Affiliation != nil {
		dpsaffiliation := affiliationMap[*sm.Affiliation]
		affiliation = &dpsaffiliation
	}

	if sm.Edipi == nil {
		return nil, &errUserMissingData{
			userID:     userIdentity.ID,
			errMessage: fmt.Sprintf("User %s is missing EDIPI", userIdentity.ID.String()),
		}
	}
	ssn, err := getSSNFromIWS(*sm.Edipi, rbs)
	if err != nil {
		return nil, errors.Wrap(err, "Getting SSN from IWS using EDIPI")
	}

	if sm.FirstName == nil || sm.LastName == nil {
		return nil, &errUserMissingData{
			userID:     userIdentity.ID,
			errMessage: fmt.Sprintf("User %s is missing first and/or last name", userIdentity.ID.String()),
		}
	}

	payload := dpsmessages.AuthenticationUserPayload{
		Affiliation:          affiliation,
		Email:                userIdentity.Email,
		FirstName:            *sm.FirstName,
		MiddleName:           sm.MiddleName,
		LastName:             *sm.LastName,
		Suffix:               sm.Suffix,
		LoginGovID:           strfmt.UUID(loginGovID),
		SocialSecurityNumber: strfmt.SSN(ssn),
		Telephone:            sm.Telephone,
	}
	return &payload, nil
}

func getSSNFromIWS(edipi string, rbs iws.RealTimeBrokerService) (string, error) {
	edipiInt, err := strconv.ParseUint(edipi, 10, 64)
	if err != nil {
		return "", errors.Wrap(err, "Converting EDIPI from string to int")
	}

	person, _, err := rbs.GetPersonUsingEDIPI(edipiInt)
	if err != nil {
		return "", errors.Wrap(err, "Using IWS")
	}

	if person.TypeCode != iws.PersonTypeCodeSSN {
		return "", errors.New("Person from IWS does not have SSN TypeCode")
	}

	return person.ID, nil
}
