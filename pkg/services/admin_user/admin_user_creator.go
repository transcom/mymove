package adminuser

import (
	"strings"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type adminUserCreator struct {
	builder adminUserQueryBuilder
	sender  notifications.NotificationSender
}

// CreateAdminUser creates admin user
func (o *adminUserCreator) CreateAdminUser(
	appCtx appcontext.AppContext,
	admin *models.AdminUser,
	organizationIDFilter []services.QueryFilter,
) (*models.AdminUser, *validate.Errors, error) {
	// Use FetchOne to see if we have an organization that matches the provided id
	var organization models.Organization
	fetchErr := o.builder.FetchOne(appCtx, &organization, organizationIDFilter)

	if fetchErr != nil {
		return nil, nil, fetchErr
	}

	// A user may already exist with that email from a previous user (office, service member, ...)
	var user models.User
	userEmailFilter := query.NewQueryFilter("login_gov_email", "=", admin.Email)
	fetchErr = o.builder.FetchOne(appCtx, &user, []services.QueryFilter{userEmailFilter})

	if fetchErr != nil {
		user = models.User{
			LoginGovEmail: strings.ToLower(admin.Email),
			Active:        true,
		}
	}

	var verrs *validate.Errors
	var err error
	var userActivityEmail notifications.Notification
	// We don't want to be left with a user record and no admin user so setup a transaction to rollback
	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if user.ID == uuid.Nil {
			verrs, err = o.builder.CreateOne(txnAppCtx, &user)
			if verrs != nil || err != nil {
				return err
			}

			email, emailErr := notifications.NewUserAccountCreated(
				appCtx, notifications.GetSysAdminEmail(o.sender), user.ID, user.UpdatedAt)
			if emailErr != nil {
				return emailErr
			}
			userActivityEmail = notifications.Notification(email)
		}

		admin.UserID = &user.ID
		admin.User = user

		verrs, err = o.builder.CreateOne(txnAppCtx, admin)
		if verrs != nil || err != nil {
			return err
		}

		return nil
	})

	if verrs != nil || txErr != nil {
		return nil, verrs, txErr
	}

	if userActivityEmail != nil {
		err = o.sender.SendNotification(appCtx, userActivityEmail)
		if err != nil {
			return nil, nil, err
		}
	}

	return admin, nil, nil
}

// NewAdminUserCreator returns a new admin user creator builder
func NewAdminUserCreator(builder adminUserQueryBuilder, sender notifications.NotificationSender) services.AdminUserCreator {
	return &adminUserCreator{builder, sender}
}
