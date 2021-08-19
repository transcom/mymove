package officeuser

import (
	"context"
	"strings"

	"github.com/transcom/mymove/pkg/notifications"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type officeUserCreator struct {
	db                 *pop.Connection
	builder            officeUserQueryBuilder
	notificationSender notifications.NotificationSender
}

// CreateOfficeUser creates office users
func (o *officeUserCreator) CreateOfficeUser(
	ctx context.Context,
	officeUser *models.OfficeUser,
	transportationIDFilter []services.QueryFilter,
) (*models.OfficeUser, *validate.Errors, error) {
	// Use FetchOne to see if we have a transportation office that matches the provided id
	var transportationOffice models.TransportationOffice
	fetchErr := o.builder.FetchOne(&transportationOffice, transportationIDFilter)

	if fetchErr != nil {
		return nil, nil, fetchErr
	}

	// A user may already exist with that email from a previous user (admin, service member, ...)
	var user models.User
	userEmailFilter := query.NewQueryFilter("login_gov_email", "=", officeUser.Email)
	fetchErr = o.builder.FetchOne(&user, []services.QueryFilter{userEmailFilter})

	if fetchErr != nil {
		user = models.User{
			LoginGovEmail: strings.ToLower(officeUser.Email),
			Active:        true,
		}
	}

	var verrs *validate.Errors
	var err error
	var userActivityEmail notifications.Notification
	// We don't want to be left with a user record and no office user so setup a transaction to rollback
	txErr := o.db.Transaction(func(connection *pop.Connection) error {
		if user.ID == uuid.Nil {
			verrs, err = o.builder.CreateOne(&user)
			if verrs != nil || err != nil {
				return err
			}

			email, emailErr := notifications.NewUserAccountCreated(ctx, "sandy@truss.works", user.ID, user.UpdatedAt)
			if emailErr != nil {
				return emailErr
			}
			userActivityEmail = notifications.Notification(email)
		}

		officeUser.UserID = &user.ID
		officeUser.User = user

		verrs, err = o.builder.CreateOne(officeUser)
		if verrs != nil || err != nil {
			return err
		}

		return nil
	})

	if verrs != nil || txErr != nil {
		return nil, verrs, txErr
	}

	if userActivityEmail != nil {
		err = o.notificationSender.SendNotification(userActivityEmail)
		if err != nil {
			return nil, nil, err
		}
	}

	return officeUser, nil, nil
}

// NewOfficeUserCreator returns a new office user creator
func NewOfficeUserCreator(
	db *pop.Connection,
	builder officeUserQueryBuilder,
	notificationSender notifications.NotificationSender,
) services.OfficeUserCreator {
	return &officeUserCreator{db, builder, notificationSender}
}
