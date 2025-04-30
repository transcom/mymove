package user

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type userUpdater struct {
	builder           userQueryBuilder
	officeUserUpdater services.OfficeUserUpdater
	adminUserUpdater  services.AdminUserUpdater
	sender            notifications.NotificationSender
}

// NewUserUpdater returns a new admin user creator builder
func NewUserUpdater(
	builder userQueryBuilder,
	officeUserUpdater services.OfficeUserUpdater,
	adminUserUpdater services.AdminUserUpdater,
	sender notifications.NotificationSender,
) services.UserUpdater {
	return &userUpdater{
		builder,
		officeUserUpdater,
		adminUserUpdater,
		sender,
	}
}

func (o *userUpdater) UpdateUser(appCtx appcontext.AppContext, id uuid.UUID, user *models.User) (*models.User, *validate.Errors, error) {
	if user == nil {
		return nil, nil, fmt.Errorf("user cannot be nil")
	}

	var foundUser models.User
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(appCtx, &foundUser, filters)
	if err != nil {
		return nil, nil, err
	}

	// determine if we are updating email, activating, or deactivating the user
	updatingEmail := foundUser.OktaEmail != user.OktaEmail
	activatingUser := user.Active && !foundUser.Active
	deactivatingUser := !user.Active && foundUser.Active

	foundUser.OktaEmail = user.OktaEmail
	foundUser.Active = user.Active

	transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := o.builder.UpdateOne(txnCtx, &foundUser, nil)
		if err != nil || (verrs != nil && verrs.HasAny()) {
			appCtx.Logger().Error("could not update user as admin", zap.Error(err), zap.Error(verrs))
			return err
		}

		// finding existing users based off the
		var existingOfficeUser models.OfficeUser
		var existingAdminUser models.AdminUser
		var existingServiceMemberUser models.ServiceMember
		userFilters := []services.QueryFilter{query.NewQueryFilter("user_id", "=", user.ID.String())}

		_ = o.builder.FetchOne(txnCtx, &existingOfficeUser, userFilters)
		_ = o.builder.FetchOne(txnCtx, &existingAdminUser, userFilters)
		_ = o.builder.FetchOne(txnCtx, &existingServiceMemberUser, userFilters)

		// update office user if found
		if existingOfficeUser.ID != uuid.Nil {
			existingOfficeUser.Active = user.Active
			existingOfficeUser.Email = user.OktaEmail
			verrs, err := o.builder.UpdateOne(txnCtx, &existingOfficeUser, nil)
			if err != nil || (verrs != nil && verrs.HasAny()) {
				appCtx.Logger().Error("could not update existing office user as admin", zap.Error(err), zap.Error(verrs))
				return err
			}
		}

		// update admin user if found
		if existingAdminUser.ID != uuid.Nil {
			existingAdminUser.Active = user.Active
			existingAdminUser.Email = user.OktaEmail
			verrs, err := o.builder.UpdateOne(txnCtx, &existingAdminUser, nil)
			if err != nil || (verrs != nil && verrs.HasAny()) {
				appCtx.Logger().Error("could not update existing admin user as admin", zap.Error(err), zap.Error(verrs))
				return err
			}
		}

		// update service member user if found
		// only need to update if the email is being updated
		if existingServiceMemberUser.ID != uuid.Nil && updatingEmail {
			existingServiceMemberUser.PersonalEmail = &user.OktaEmail
			verrs, err := o.builder.UpdateOne(txnCtx, &existingServiceMemberUser, nil)
			if err != nil || (verrs != nil && verrs.HasAny()) {
				appCtx.Logger().Error("could not update existing service member user as admin", zap.Error(err), zap.Error(verrs))
				return err
			}
		}

		// if the user email is being updated, we need to also update the Okta profile
		if updatingEmail && foundUser.OktaEmail != "" && appCtx.Session().IDToken != "devlocal" {
			req := appCtx.HTTPRequest()
			if req == nil {
				return fmt.Errorf("failed to retrieve HTTP request from session")
			}
			provider, err := okta.GetOktaProviderForRequest(req)
			if err != nil {
				return fmt.Errorf("error retrieving Okta provider: %w", err)
			}

			apiKey := models.GetOktaAPIKey()
			oktaID := foundUser.OktaID

			existingOktaUser, err := models.GetOktaUser(txnCtx, provider, oktaID, apiKey)
			if err != nil || existingOktaUser == nil {
				return fmt.Errorf("failed to fetch Okta user before update")
			}

			existingOktaUser.Profile.Email = foundUser.OktaEmail
			existingOktaUser.Profile.Login = foundUser.OktaEmail

			_, err = models.UpdateOktaUser(txnCtx, provider, oktaID, apiKey, *existingOktaUser)
			if err != nil {
				return fmt.Errorf("error updating Okta user: %w", err)
			}
		}

		return nil
	})

	if transactionError != nil {
		return nil, nil, transactionError
	}

	// sending notifications based on if the user is being activated/deactivated
	var userActivityEmail notifications.Notification
	if activatingUser && !deactivatingUser {
		email, emailErr := notifications.NewUserAccountActivated(
			appCtx, notifications.GetSysAdminEmail(o.sender), foundUser.ID, foundUser.UpdatedAt)
		if emailErr != nil {
			appCtx.Logger().Error("Could not send user activation email", zap.Error(emailErr))
		} else {
			userActivityEmail = notifications.Notification(email)
		}
	} else if !activatingUser && deactivatingUser {
		email, emailErr := notifications.NewUserAccountDeactivated(
			appCtx, notifications.GetSysAdminEmail(o.sender), foundUser.ID, foundUser.UpdatedAt)
		if emailErr != nil {
			appCtx.Logger().Error("Could not send user deactivation email", zap.Error(emailErr))
		} else {
			userActivityEmail = notifications.Notification(email)
		}
	}

	if userActivityEmail != nil {
		err = o.sender.SendNotification(appCtx, userActivityEmail)
		if err != nil {
			return nil, nil, err
		}
	}

	return &foundUser, nil, nil
}
