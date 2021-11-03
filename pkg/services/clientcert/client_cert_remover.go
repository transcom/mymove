package clientcert

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type clientCertRemover struct {
	builder clientCertQueryBuilder
	sender  notifications.NotificationSender
}

func (o *clientCertRemover) RemoveClientCert(appCtx appcontext.AppContext, id uuid.UUID) (*models.ClientCert, *validate.Errors, error) {
	var foundClientCert models.ClientCert
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(appCtx, &foundClientCert, filters)

	if err != nil {
		return nil, nil, err
	}

	err = o.builder.DeleteOne(appCtx, &foundClientCert)
	if err != nil {
		return nil, nil, err
	}

	session := appCtx.Session()
	if session == nil {
		return nil, nil, apperror.NewContextError("Unable to find Session in Context")
	}
	email, emailErr := notifications.NewClientCertRemoved(
		notifications.GetSysAdminEmail(o.sender), foundClientCert.ID, foundClientCert.UpdatedAt, session.UserID, session.Hostname)
	if emailErr != nil {
		return nil, nil, emailErr
	}
	userActivityEmail := notifications.Notification(email)
	err = o.sender.SendNotification(appCtx, userActivityEmail)
	if err != nil {
		return nil, nil, err
	}

	return &foundClientCert, nil, nil
}

// NewClientCertRemover returns a new admin user updater builder
func NewClientCertRemover(builder clientCertQueryBuilder, sender notifications.NotificationSender) services.ClientCertRemover {
	return &clientCertRemover{builder, sender}
}
