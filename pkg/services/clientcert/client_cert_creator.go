package clientcert

import (
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/notifications"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type clientCertCreator struct {
	builder clientCertQueryBuilder
	sender  notifications.NotificationSender
}

// CreateClientCert creates admin user
func (o *clientCertCreator) CreateClientCert(
	appCtx appcontext.AppContext,
	cert *models.ClientCert,
) (*models.ClientCert, *validate.Errors, error) {

	var verrs *validate.Errors
	var err error
	var userActivityEmail notifications.Notification
	// We don't want to be left with a user record and no admin user so setup a transaction to rollback
	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {

		verrs, err = o.builder.CreateOne(txnAppCtx, cert)
		if verrs != nil || err != nil {
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
		err = o.sender.SendNotification(appCtx, userActivityEmail)
		if err != nil {
			return nil, nil, err
		}
	}

	return cert, nil, nil
}

// NewClientCertCreator returns a new admin user creator builder
func NewClientCertCreator(builder clientCertQueryBuilder, sender notifications.NotificationSender) services.ClientCertCreator {
	return &clientCertCreator{builder, sender}
}
