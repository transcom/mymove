package user

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
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

// UpdateUser updates any user
func (o *userUpdater) UpdateUser(appCtx appcontext.AppContext, id uuid.UUID, user *models.User) (*models.User, *validate.Errors, error) {
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	var foundUser models.User
	var userActivityEmail notifications.Notification

	if user == nil {
		return nil, nil, nil
	}
	// Find the existing user to update
	err := o.builder.FetchOne(appCtx, &foundUser, filters)

	if err != nil {
		return nil, nil, err
	}

	// Update user's new status for Active
	userWasActive := foundUser.Active
	foundUser.Active = user.Active

	verrs, err := o.builder.UpdateOne(appCtx, &foundUser, nil)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	// If the update was successful and we are deactivating the user,
	// update the office and admin statuses to match.

	// Check if we are deactivating the user
	if !user.Active {

		// Check for Office User
		foundOfficeUser := models.OfficeUser{}
		filters = []services.QueryFilter{query.NewQueryFilter("user_id", "=", id.String())}
		err = o.builder.FetchOne(appCtx, &foundOfficeUser, filters)

		// If we find a matching Office User, update their status
		if err == nil {
			payload := adminmessages.OfficeUserUpdatePayload{
				Active: &user.Active,
			}
			_, verrs, err = o.officeUserUpdater.UpdateOfficeUser(appCtx, foundOfficeUser.ID, &payload)

			if verrs != nil {
				appCtx.Logger().Error("Could not update office user", zap.Error(verrs))
			} else if err != nil {
				appCtx.Logger().Error("Could not update office user", zap.Error(err))
			}
		}

		// Check for Admin User
		foundAdminUser := models.AdminUser{}
		err = o.builder.FetchOne(appCtx, &foundAdminUser, filters)
		// If we find a matching Admin User, update their status
		if err == nil {
			payload := adminmessages.AdminUserUpdatePayload{
				Active: &user.Active,
			}
			_, verrs, err = o.adminUserUpdater.UpdateAdminUser(appCtx, foundAdminUser.ID, &payload)
			if verrs != nil {
				appCtx.Logger().Error("Could not update admin user", zap.Error(verrs))
			} else if err != nil {
				appCtx.Logger().Error("Could not update admin user", zap.Error(err))
			}
		}

		if userWasActive {
			email, emailErr := notifications.NewUserAccountDeactivated(
				appCtx, notifications.GetSysAdminEmail(o.sender), foundUser.ID, foundUser.UpdatedAt)
			if emailErr != nil {
				appCtx.Logger().Error("Could not send user deactivation email", zap.Error(emailErr))
			} else {
				userActivityEmail = notifications.Notification(email)
			}
		}
	} else if !userWasActive {
		email, emailErr := notifications.NewUserAccountActivated(
			appCtx, notifications.GetSysAdminEmail(o.sender), foundUser.ID, foundUser.UpdatedAt)
		if emailErr != nil {
			appCtx.Logger().Error("Could not send user activation email", zap.Error(emailErr))
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
