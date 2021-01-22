package user

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
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
func (o *userUpdater) UpdateUser(id uuid.UUID, payload *adminmessages.UserUpdatePayload) (*models.User, *validate.Errors, error) {
	var foundUser models.User

	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(&foundUser, filters)

	if err != nil {
		return nil, nil, err
	}

	if payload.Active != nil {
		foundUser.Active = *payload.Active
	}

	// If we are changing Active to False, we are
	// deactivating the user. We also need to revoke all
	// sessions for this user.

	verrs, err := o.builder.UpdateOne(&foundUser, nil)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return &foundUser, nil, nil

}
