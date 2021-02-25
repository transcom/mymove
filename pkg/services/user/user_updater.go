package user

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type userUpdater struct {
	builder           userQueryBuilder
	officeUserUpdater services.OfficeUserUpdater
	adminUserUpdater  services.AdminUserUpdater
}

// NewUserUpdater returns a new admin user creator builder
func NewUserUpdater(builder userQueryBuilder, officeUserUpdater services.OfficeUserUpdater, adminUserUpdater services.AdminUserUpdater) services.UserUpdater {
	return &userUpdater{
		builder,
		officeUserUpdater,
		adminUserUpdater,
	}
}

// UpdateUser updates any user
func (o *userUpdater) UpdateUser(id uuid.UUID, user *models.User) (*models.User, *validate.Errors, error) {
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	var foundUser models.User
	var logger Logger

	if user == nil {
		return nil, nil, nil
	}
	// Find the existing user to update
	err := o.builder.FetchOne(&foundUser, filters)

	if err != nil {
		return nil, nil, err
	}

	// Update user's new status for Active
	if &user.Active != nil {
		foundUser.Active = user.Active
	}

	verrs, err := o.builder.UpdateOne(&foundUser, nil)
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
		err = o.builder.FetchOne(&foundOfficeUser, filters)

		// If we find a matching Office User, update their status
		if err == nil {
			payload := adminmessages.OfficeUserUpdatePayload{
				Active: &user.Active,
			}
			_, verrs, err = o.officeUserUpdater.UpdateOfficeUser(foundOfficeUser.ID, &payload)

			if verrs != nil {
				logger.Info("Could not update office user", zap.Error(verrs))
			} else if err != nil {
				logger.Info("Could not update office user", zap.Error(err))
			}
		}

		// Check for Admin User
		foundAdminUser := models.AdminUser{}
		err = o.builder.FetchOne(&foundAdminUser, filters)
		// If we find a matching Admin User, update their status
		if err == nil {
			payload := adminmessages.AdminUserUpdatePayload{
				Active: &user.Active,
			}
			_, verrs, err = o.adminUserUpdater.UpdateAdminUser(foundAdminUser.ID, &payload)
			if verrs != nil {
				logger.Info("Could not update admin user", zap.Error(verrs))
			} else if err != nil {
				logger.Info("Could not update admin user", zap.Error(err))
			}
		}
	}

	return &foundUser, nil, nil

}
