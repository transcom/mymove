package officeuser

import (
	"database/sql"
	"fmt"
	"net/http/httptest"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/jarcoal/httpmock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers/adminapi/payloads"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *OfficeUserServiceSuite) TestUpdateOfficeUser() {
	queryBuilder := query.NewQueryBuilder()
	updater := NewOfficeUserUpdater(queryBuilder)

	oktaProvider := okta.NewOktaProvider(suite.Logger())
	err := oktaProvider.RegisterOktaProvider("adminProvider", "OrgURL", "CallbackURL", "fakeToken", "secret", []string{"openid", "profile", "email"})
	suite.NoError(err)

	// Happy path
	suite.Run("If the user is updated successfully it should be returned", func() {
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		transportationOffice := factory.BuildDefaultTransportationOffice(suite.DB())
		primaryOffice := true

		firstName := "Lea"
		middleInitials := "L"
		payload := &adminmessages.OfficeUserUpdate{
			FirstName:      &firstName,
			MiddleInitials: &middleInitials,
			TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
				{
					TransportationOfficeID: strfmt.UUID(transportationOffice.ID.String()),
					PrimaryOffice:          &primaryOffice,
				},
			},
			Active: models.BoolPointer(true),
		}

		officeUserUpdatesModel := payloads.OfficeUserModelFromUpdate(payload, &officeUser)

		updatedOfficeUser, verrs, err := updater.UpdateOfficeUser(suite.AppContextForTest(), officeUser.ID, officeUserUpdatesModel, uuid.FromStringOrNil(transportationOffice.ID.String()))
		suite.NoError(err)
		suite.Nil(verrs)
		suite.Equal(updatedOfficeUser.ID.String(), officeUser.ID.String())
		suite.Equal(updatedOfficeUser.TransportationOfficeID.String(), transportationOffice.ID.String())
		suite.NotEqual(updatedOfficeUser.TransportationOfficeID.String(), officeUser.TransportationOffice.ID.String())
		suite.Equal(updatedOfficeUser.FirstName, firstName)
		suite.Equal(updatedOfficeUser.LastName, officeUser.LastName)
		suite.Equal(updatedOfficeUser.MiddleInitials, payload.MiddleInitials)
		suite.Equal(updatedOfficeUser.Active, *payload.Active)
	})

	// Bad office user ID
	suite.Run("If we are provided an office user that doesn't exist, the create should fail", func() {
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		payload := &adminmessages.OfficeUserUpdate{}
		officeUserUpdatesModel := payloads.OfficeUserModelFromUpdate(payload, &officeUser)

		_, _, err := updater.UpdateOfficeUser(suite.AppContextForTest(), uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"), officeUserUpdatesModel, uuid.Nil)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows.Error(), err.Error())
	})

	// Bad transportation office ID
	suite.Run("If we are provided a transportation office that doesn't exist, the create should fail", func() {
		officeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.Country{
					Country:     "US",
					CountryName: "UNITED STATES",
				},
			},
		}, nil)
		primaryOffice := true

		payload := &adminmessages.OfficeUserUpdate{
			TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
				{
					TransportationOfficeID: strfmt.UUID("00000000-0000-0000-0000-000000000001"),
					PrimaryOffice:          &primaryOffice,
				},
			},
		}

		officeUserUpdatesModel := payloads.OfficeUserModelFromUpdate(payload, &officeUser)

		_, _, err := updater.UpdateOfficeUser(suite.AppContextForTest(), officeUser.ID, officeUserUpdatesModel, uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"))
		suite.Error(err)
		suite.Equal(sql.ErrNoRows.Error(), err.Error())
	})

	suite.Run("updating office user also updates the associated user email", func() {
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
		transportationOffice := factory.BuildDefaultTransportationOffice(suite.DB())
		primaryOffice := true

		mockAndActivateOktaGETEndpointNoError(officeUser.User.OktaID)
		mockAndActivateOktaPOSTEndpointNoError(officeUser.User.OktaID)

		request := httptest.NewRequest("PATCH", fmt.Sprintf("/office-users/%s", officeUser.UserID.String()), nil)

		session := &auth.Session{
			ApplicationName: auth.AdminApp,
			Hostname:        "adminlocal",
		}

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), session, request)

		payload := &adminmessages.OfficeUserUpdate{
			Email: models.StringPointer("newEmail@mail.mil"),
			TransportationOfficeAssignments: []*adminmessages.OfficeUserTransportationOfficeAssignment{
				{
					TransportationOfficeID: strfmt.UUID(transportationOffice.ID.String()),
					PrimaryOffice:          &primaryOffice,
				},
			},
		}

		officeUserUpdatesModel := payloads.OfficeUserModelFromUpdate(payload, &officeUser)

		updatedOfficeUser, verrs, err := updater.UpdateOfficeUser(appCtx, officeUser.ID, officeUserUpdatesModel, uuid.Nil)
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
