package adminuser

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type adminUserUpdater struct {
	builder adminUserQueryBuilder
}

func (o *adminUserUpdater) UpdateAdminUser(appCtx appcontext.AppContext, id uuid.UUID, payload *adminmessages.AdminUserUpdate) (*models.AdminUser, *validate.Errors, error) {
	updateUserAndOkta := false
	var foundUser models.AdminUser
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(appCtx, &foundUser, filters)
	if err != nil {
		return nil, nil, err
	}

	if payload.Email != nil && payload.Email != &foundUser.Email {
		foundUser.Email = *payload.Email
		updateUserAndOkta = true
	}

	if payload.FirstName != nil {
		foundUser.FirstName = *payload.FirstName
	}

	if payload.LastName != nil {
		foundUser.LastName = *payload.LastName
	}

	if payload.Active != nil {
		foundUser.Active = *payload.Active
	}

	transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := o.builder.UpdateOne(txnCtx, &foundUser, nil)
		if err != nil {
			return err
		}
		if verrs != nil && verrs.HasAny() {
			return verrs
		}

		// if the email is being updated, we need to also update the user email & okta email
		if updateUserAndOkta {
			var existingUser models.User
			filters := []services.QueryFilter{query.NewQueryFilter("id", "=", foundUser.UserID.String())}
			err := o.builder.FetchOne(appCtx, &existingUser, filters)
			if err != nil {
				return err
			}

			existingUser.OktaEmail = foundUser.Email
			verrs, err := o.builder.UpdateOne(txnCtx, &existingUser, nil)
			if err != nil {
				return err
			}
			if verrs != nil && verrs.HasAny() {
				return verrs
			}

			if existingUser.OktaID != "" && appCtx.Session().IDToken != "devlocal" {
				apiKey := models.GetOktaAPIKey()
				oktaID := existingUser.OktaID
				req := appCtx.HTTPRequest()
				if req == nil {
					return fmt.Errorf("failed to retrieve HTTP request from session")
				}

				// Use the HTTP request to get the Okta provider
				provider, err := okta.GetOktaProviderForRequest(req)
				if err != nil {
					return fmt.Errorf("error retrieving Okta provider: %w", err)
				}

				// verifying the okta user exists but we also need all the okta profile info prior to updating
				existingOktaUser, err := models.GetOktaUser(appCtx, provider, oktaID, apiKey)
				if err != nil {
					return fmt.Errorf("error getting Okta user prior to updating: %w", err)
				}
				if existingOktaUser == nil {
					return fmt.Errorf("okta user cannot be nil before updating okta email of admin user")
				}
				existingOktaUser.Profile.Email = existingUser.OktaEmail
				existingOktaUser.Profile.Login = existingUser.OktaEmail

				_, err = models.UpdateOktaUser(appCtx, provider, oktaID, apiKey, *existingOktaUser)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if transactionError != nil {
		return nil, nil, transactionError
	}

	return &foundUser, nil, nil
}

// NewAdminUserUpdater returns a new admin user updater builder
func NewAdminUserUpdater(builder adminUserQueryBuilder) services.AdminUserUpdater {
	return &adminUserUpdater{builder}
}
