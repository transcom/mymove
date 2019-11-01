package adminuser

import (
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type adminUserUpdater struct {
	builder adminUserQueryBuilder
}

func (o *adminUserUpdater) UpdateAdminUser(user *models.AdminUser) (*models.AdminUser, *validate.Errors, error) {
	var foundUser models.AdminUser
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", user.ID.String())}
	err := o.builder.FetchOne(&foundUser, filters)

	if err != nil {
		return nil, nil, err
	}

	foundUser.FirstName = user.FirstName
	foundUser.LastName = user.LastName
	foundUser.Active = user.Active

	verrs, err := o.builder.UpdateOne(&foundUser)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return &foundUser, nil, nil
}

func NewAdminUserUpdater(builder adminUserQueryBuilder) services.AdminUserUpdater {
	return &adminUserUpdater{builder}
}
