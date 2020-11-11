package adminuser

import (
	"strings"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type adminUserCreator struct {
	db      *pop.Connection
	builder adminUserQueryBuilder
}

// CreateAdminUser creates admin user
func (o *adminUserCreator) CreateAdminUser(admin *models.AdminUser, organizationIDFilter []services.QueryFilter) (*models.AdminUser, *validate.Errors, error) {
	// Use FetchOne to see if we have an organization that matches the provided id
	var organization models.Organization
	fetchErr := o.builder.FetchOne(&organization, organizationIDFilter)

	if fetchErr != nil {
		return nil, nil, fetchErr
	}

	user := &models.User{
		LoginGovEmail: strings.ToLower(admin.Email),
		Active:        true,
	}

	var verrs *validate.Errors
	var err error
	//
	txErr := o.db.Transaction(func(connection *pop.Connection) error {
		verrs, err = o.builder.CreateOne(user)
		if verrs != nil || err != nil {
			return err
		}

		admin.UserID = &user.ID
		admin.User = *user

		verrs, err = o.builder.CreateOne(admin)
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
func NewAdminUserCreator(db *pop.Connection, builder adminUserQueryBuilder) services.AdminUserCreator {
	return &adminUserCreator{db, builder}
}
