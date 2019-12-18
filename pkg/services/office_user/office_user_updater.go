package officeuser

import (
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type officeUserUpdater struct {
	builder officeUserQueryBuilder
}

func (o *officeUserUpdater) UpdateOfficeUser(id uuid.UUID, payload *adminmessages.OfficeUserUpdatePayload) (*models.OfficeUser, *validate.Errors, error) {
	var foundUser models.OfficeUser
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(&foundUser, filters)

	if err != nil {
		return nil, nil, err
	}

	if payload.FirstName != nil {
		foundUser.FirstName = *payload.FirstName
	}

	if payload.MiddleInitials != nil {
		foundUser.MiddleInitials = payload.MiddleInitials
	}

	if payload.LastName != nil {
		foundUser.LastName = *payload.LastName
	}

	if payload.Telephone != nil {
		foundUser.Telephone = *payload.Telephone
	}

	if payload.Active != nil {
		foundUser.Active = *payload.Active
	}

	verrs, err := o.builder.UpdateOne(&foundUser)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return &foundUser, nil, nil
}

func NewOfficeUserUpdater(builder officeUserQueryBuilder) services.OfficeUserUpdater {
	return &officeUserUpdater{builder}
}
