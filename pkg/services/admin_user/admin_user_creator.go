package adminuser

import (
	"strings"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type adminUserCreator struct {
	builder adminUserQueryBuilder
}

// CreateAdminUser creates admin user
func (o *adminUserCreator) CreateAdminUser(admin *models.AdminUser, organizationIDFilter []services.QueryFilter) (*models.AdminUser, *validate.Errors, error) {
	// Use FetchOne to see if we have an organization that matches the provided id
	var organization models.Organization
	err := o.builder.FetchOne(&organization, organizationIDFilter)

	if err != nil {
		return nil, nil, err
	}

	user := &models.User{
		LoginGovEmail: strings.ToLower(admin.Email),
		Active:        true,
	}

	verrs, err := o.builder.CreateOne(user)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	admin.UserID = &user.ID
	admin.User = *user

	verrs, err = o.builder.CreateOne(admin)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return admin, nil, nil
}

// NewAdminUserCreator returns a new admin user creator builder
func NewAdminUserCreator(builder adminUserQueryBuilder) services.AdminUserCreator {
	return &adminUserCreator{builder}
}
