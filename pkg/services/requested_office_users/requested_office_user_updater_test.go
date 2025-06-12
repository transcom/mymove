package adminuser

import (
	"database/sql"
	"fmt"
	"net/http/httptest"

	"github.com/gofrs/uuid"
	"github.com/jarcoal/httpmock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *RequestedOfficeUsersServiceSuite) TestUpdateRequestedOfficeUser() {
	queryBuilder := query.NewQueryBuilder()
	updater := NewRequestedOfficeUserUpdater(queryBuilder)
	setupTestData := func() models.OfficeUser {
		officeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Active: true,
					Email:  "officeUser@mail.mil",
				},
			},
			{
				Model: models.User{
					Active:    true,
					OktaEmail: "officeUser@mail.mil",
				},
			},
		}, nil)
		return officeUser
	}

	oktaProvider := okta.NewOktaProvider(suite.Logger())
	err := oktaProvider.RegisterOktaProvider("adminProvider", "OrgURL", "CallbackURL", "fakeToken", "secret", []string{"openid", "profile", "email"})
	suite.NoError(err)

	// Happy path
	suite.Run("If the user is updated successfully it should be returned", func() {
		officeUser := setupTestData()
		transportationOffice := factory.BuildDefaultTransportationOffice(suite.DB())

		firstName := "Jimmy"
		lastName := "Jim"
		status := "APPROVED"
		payload := &adminmessages.RequestedOfficeUserUpdate{
			FirstName:              &firstName,
			LastName:               &lastName,
			TransportationOfficeID: handlers.FmtUUID(transportationOffice.ID),
			Status:                 status,
		}
		updatedOfficeUser, verrs, err := updater.UpdateRequestedOfficeUser(suite.AppContextForTest(), officeUser.ID, payload)
		suite.NoError(err)
		suite.Nil(verrs)
		suite.Equal(updatedOfficeUser.ID.String(), officeUser.ID.String())
		suite.Equal(updatedOfficeUser.TransportationOfficeID.String(), transportationOffice.ID.String())
		suite.NotEqual(updatedOfficeUser.TransportationOfficeID.String(), officeUser.TransportationOffice.ID.String())
		suite.Equal(updatedOfficeUser.FirstName, firstName)
		suite.Equal(updatedOfficeUser.LastName, lastName)
		suite.Equal(updatedOfficeUser.Active, true)
	})

	// Bad office user ID
	suite.Run("If we are provided an office user that doesn't exist, the create should fail", func() {
		payload := &adminmessages.RequestedOfficeUserUpdate{}

		_, _, err := updater.UpdateRequestedOfficeUser(suite.AppContextForTest(), uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"), payload)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows.Error(), err.Error())
	})

	// Bad transportation office ID
	suite.Run("If we are provided a transportation office that doesn't exist, the create should fail", func() {
		officeUser := setupTestData()
		badID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
		payload := &adminmessages.RequestedOfficeUserUpdate{
			TransportationOfficeID: handlers.FmtUUID(badID),
		}

		_, _, err := updater.UpdateRequestedOfficeUser(suite.AppContextForTest(), officeUser.ID, payload)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows.Error(), err.Error())
	})

	suite.Run("updating requested office user also updates the associated user email", func() {
		officeUser := setupTestData()
		transportationOffice := factory.BuildDefaultTransportationOffice(suite.DB())
		mockAndActivateOktaGETEndpointNoError(officeUser.User.OktaID)
		mockAndActivateOktaPOSTEndpointNoError(officeUser.User.OktaID)

		request := httptest.NewRequest("PATCH", fmt.Sprintf("/requested-office-users/%s", officeUser.UserID.String()), nil)

		session := &auth.Session{
			ApplicationName: auth.AdminApp,
			Hostname:        "adminlocal",
		}

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), session, request)

		payload := &adminmessages.RequestedOfficeUserUpdate{
			Email:                  models.StringPointer("newEmail@mail.mil"),
			TransportationOfficeID: handlers.FmtUUID(transportationOffice.ID),
		}

		updatedOfficeUser, verrs, err := updater.UpdateRequestedOfficeUser(appCtx, officeUser.ID, payload)
		suite.NoError(err)
		suite.Nil(verrs)

		updatedUser := models.User{}
		err = suite.DB().Find(&updatedUser, updatedOfficeUser.UserID)
		suite.NoError(err)
		suite.Equal(updatedUser.OktaEmail, updatedOfficeUser.Email)
	})
}

func mockAndActivateOktaGETEndpointNoError(oktaID string) {
	httpmock.Activate()
	getUsersEndpoint := "OrgURL/api/v1/users/" + oktaID
	response := fmt.Sprintf(`{
			"id": "%s",
			"status": "ACTIVE",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "email@email.com",
				"email": "email@email.com",
				"cac_edipi": "1234567890"
			}
		}`, oktaID)

	httpmock.RegisterResponder("GET", getUsersEndpoint,
		httpmock.NewStringResponder(200, response))
}

func mockAndActivateOktaPOSTEndpointNoError(oktaID string) {
	httpmock.Activate()
	updateUsersEndpoint := "OrgURL/api/v1/users/" + oktaID
	response := fmt.Sprintf(`{
			"id": "%s",
			"status": "ACTIVE",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "email@email.com",
				"email": "email@email.com",
				"cac_edipi": "1234567890"
			}
		}`, oktaID)

	httpmock.RegisterResponder("POST", updateUsersEndpoint,
		httpmock.NewStringResponder(200, response))
}
