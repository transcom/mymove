package adminuser

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type adminUserUpdater struct {
	builder adminUserQueryBuilder
}

func (o *adminUserUpdater) UpdateAdminUser(appCtx appcontext.AppContext, id uuid.UUID, payload *adminmessages.AdminUserUpdate) (*models.AdminUser, *validate.Errors, error) {
	var foundUser models.AdminUser
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(appCtx, &foundUser, filters)

	if err != nil {
		return nil, nil, err
	}

	if payload.FirstName != nil {
		foundUser.FirstName = *payload.FirstName
	}

	if payload.LastName != nil {
		foundUser.LastName = *payload.LastName
	}

	if payload.Active != nil {
		foundUser.Active = *payload.Active
	}

	verrs, err := o.builder.UpdateOne(appCtx, &foundUser, nil)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return &foundUser, nil, nil
}

// NewAdminUserUpdater returns a new admin user updater builder
func NewAdminUserUpdater(builder adminUserQueryBuilder) services.AdminUserUpdater {
	return &adminUserUpdater{builder}
}
