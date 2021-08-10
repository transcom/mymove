package adminuser

import (
	"strings"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type adminUserCreator struct {
	builder adminUserQueryBuilder
}

// CreateAdminUser creates admin user
func (o *adminUserCreator) CreateAdminUser(appCfg appconfig.AppConfig, admin *models.AdminUser, organizationIDFilter []services.QueryFilter) (*models.AdminUser, *validate.Errors, error) {
	// Use FetchOne to see if we have an organization that matches the provided id
	var organization models.Organization
	fetchErr := o.builder.FetchOne(appCfg, &organization, organizationIDFilter)

	if fetchErr != nil {
		return nil, nil, fetchErr
	}

	// A user may already exist with that email from a previous user (office, service member, ...)
	var user models.User
	userEmailFilter := query.NewQueryFilter("login_gov_email", "=", admin.Email)
	fetchErr = o.builder.FetchOne(appCfg, &user, []services.QueryFilter{userEmailFilter})

	if fetchErr != nil {
		user = models.User{
			LoginGovEmail: strings.ToLower(admin.Email),
			Active:        true,
		}
	}

	var verrs *validate.Errors
	var err error
	// We don't want to be left with a user record and no admin user so setup a transaction to rollback
	txErr := appCfg.NewTransaction(func(txnAppCfg appconfig.AppConfig) error {
		if user.ID == uuid.Nil {
			verrs, err = o.builder.CreateOne(txnAppCfg, &user)
			if verrs != nil || err != nil {
				return err
			}
		}

		admin.UserID = &user.ID
		admin.User = user

		verrs, err = o.builder.CreateOne(txnAppCfg, admin)
		if verrs != nil || err != nil {
			return err
		}

		return nil
	})

	if verrs != nil || txErr != nil {
		return nil, verrs, txErr
	}

	return admin, nil, nil
}

// NewAdminUserCreator returns a new admin user creator builder
func NewAdminUserCreator(builder adminUserQueryBuilder) services.AdminUserCreator {
	return &adminUserCreator{builder}
}
