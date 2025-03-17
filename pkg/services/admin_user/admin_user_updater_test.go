package adminuser

import (
	"fmt"
	"net/http/httptest"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/jarcoal/httpmock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *AdminUserServiceSuite) TestUpdateAdminUser() {
	newUUID, _ := uuid.NewV4()

	firstName := "Leo"
	payload := &adminmessages.AdminUserUpdate{
		FirstName: &firstName,
	}

	oktaProvider := okta.NewOktaProvider(suite.Logger())
	err := oktaProvider.RegisterOktaProvider("adminProvider", "OrgURL", "CallbackURL", "fakeToken", "secret", []string{"openid", "profile", "email"})
	suite.NoError(err)

	// Happy path
	suite.Run("If the user is updated successfully it should be returned", func() {
		fakeUpdateOne := func(appcontext.AppContext, interface{}, *string) (*validate.Errors, error) {
			return nil, nil
		}

		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			return nil
		}

		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeUpdateOne: fakeUpdateOne,
		}
		updater := NewAdminUserUpdater(builder)
		_, verrs, err := updater.UpdateAdminUser(suite.AppContextForTest(), newUUID, payload)
		suite.NoError(err)
		suite.Nil(verrs)
	})

	// Bad organization ID
	suite.Run("If we are provided a organization that doesn't exist, the create should fail", func() {
		fakeUpdateOne := func(appcontext.AppContext, interface{}, *string) (*validate.Errors, error) {
			return nil, nil
		}
		fakeFetchOne := func(appCtx appcontext.AppContext, model interface{}) error {
			return models.ErrFetchNotFound
		}
		builder := &testAdminUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeUpdateOne: fakeUpdateOne,
		}
		updater := NewAdminUserUpdater(builder)
		_, _, err := updater.UpdateAdminUser(suite.AppContextForTest(), newUUID, payload)
		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound.Error(), err.Error())
	})

	suite.Run("updating office user also updates the associated user email", func() {
		adminUser := factory.BuildAdminUser(suite.DB(), []factory.Customization{
			{
				Model: models.AdminUser{
					Active: true,
					Email:  "adminUser@mail.mil",
				},
			},
			{
				Model: models.User{
					Active:    true,
					OktaEmail: "adminUser@mail.mil",
				},
			},
		}, nil)
		queryBuilder := query.NewQueryBuilder()
		updater := NewAdminUserUpdater(queryBuilder)

		mockAndActivateOktaGETEndpointNoError(adminUser.User.OktaID)
		mockAndActivateOktaPOSTEndpointNoError(adminUser.User.OktaID)

		request := httptest.NewRequest("PATCH", fmt.Sprintf("/admin-users/%s", adminUser.UserID.String()), nil)

		session := &auth.Session{
			ApplicationName: auth.AdminApp,
			Hostname:        "adminlocal",
		}

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), session, request)

		payload := &adminmessages.AdminUserUpdate{
			FirstName: &firstName,
			Email:     models.StringPointer("newEmail@mail.mil"),
		}

		updatedAdminUser, verrs, err := updater.UpdateAdminUser(appCtx, adminUser.ID, payload)
		suite.NoError(err)
		suite.Nil(verrs)

		updatedUser := models.User{}
		err = suite.DB().Find(&updatedUser, updatedAdminUser.UserID)
		suite.NoError(err)
		suite.Equal(updatedUser.OktaEmail, updatedAdminUser.Email)
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
