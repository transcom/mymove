package user

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type userUpdater struct {
	builder userQueryBuilder
}

// NewUserUpdater returns a new admin user creator builder
func NewUserUpdater(builder userQueryBuilder) services.UserUpdater {
	return &userUpdater{builder}
}

// UpdateUser updates any user
func (o *userUpdater) UpdateUser(id uuid.UUID, user *models.User) (*models.User, *validate.Errors, error) {
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	var foundUser models.User

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

	return &foundUser, nil, nil

}
