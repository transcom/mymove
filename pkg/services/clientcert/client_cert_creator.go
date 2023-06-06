package clientcert

import (
	"database/sql"
	"strings"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type clientCertCreator struct {
	builder clientCertQueryBuilder
	services.UserRoleAssociator
	sender notifications.NotificationSender
}

// CreateClientCert creates admin user
func (o *clientCertCreator) CreateClientCert(
	appCtx appcontext.AppContext,
	email string,
	cert *models.ClientCert,
) (*models.ClientCert, *validate.Errors, error) {
	var verrs *validate.Errors
	var userActivityEmail notifications.Notification
	// We don't want to be left with a user record and no admin user so setup a transaction to rollback
	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {

		// A user may already exist with the email from a previous user
		// (admin, office, ...) if we are creating a CAC login certificate
		var user models.User
		userEmailFilter := query.NewQueryFilter("login_gov_email", "=", email)
		err := o.builder.FetchOne(txnAppCtx, &user, []services.QueryFilter{userEmailFilter})

		if err != nil && err != sql.ErrNoRows {
			return err
		}

		// if the fetch failed, the user doesn't exist, so we will
		// need to create one
		// This logic is similar to what is used when creating office users
		if err == sql.ErrNoRows {
			user = models.User{
				LoginGovEmail: strings.ToLower(email),
				Active:        true,
			}

			verrs, err = o.builder.CreateOne(txnAppCtx, &user)
			if verrs != nil {
				return apperror.NewInvalidCreateInputError(verrs, "Invalid user params")
			}
			if err != nil {
				return err
			}
		}

		// ensure this user has the prime role if necessary
		err = updatePrimeRoleForUser(appCtx, user.ID, cert, o.builder,
			o.UserRoleAssociator)
		if err != nil {
			return err
		}

		// assign the userID to the cert
		cert.UserID = user.ID
		verrs, err = o.builder.CreateOne(txnAppCtx, cert)
		if verrs != nil {
			return apperror.NewInvalidCreateInputError(verrs, "Invalid cert params")
		}
		if err != nil {
			return err
		}

		session := txnAppCtx.Session()
		if session == nil {
			return apperror.NewContextError("Unable to find Session in Context")
		}
		email, emailErr := notifications.NewClientCertCreated(
			notifications.GetSysAdminEmail(o.sender), cert.ID, cert.UpdatedAt, session.UserID, session.Hostname)
		if emailErr != nil {
			return emailErr
		}
		userActivityEmail = notifications.Notification(email)

		return nil
	})

	if verrs != nil || txErr != nil {
		return nil, verrs, txErr
	}

	if userActivityEmail != nil {
		err := o.sender.SendNotification(appCtx, userActivityEmail)
		if err != nil {
			return nil, nil, err
		}
	}

	return cert, nil, nil
}

// NewClientCertCreator returns a new admin user creator builder
func NewClientCertCreator(builder clientCertQueryBuilder, userRoleAssociator services.UserRoleAssociator, sender notifications.NotificationSender) services.ClientCertCreator {
	return &clientCertCreator{builder, userRoleAssociator, sender}
}
