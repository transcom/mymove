package officeuser

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type officeUserUpdater struct {
	builder officeUserQueryBuilder
}

// UpdateOfficeUser updates an office user
func (o *officeUserUpdater) UpdateOfficeUser(appCtx appcontext.AppContext, id uuid.UUID, payload *adminmessages.OfficeUserUpdate) (*models.OfficeUser, *validate.Errors, error) {
	var foundUser models.OfficeUser
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(appCtx, &foundUser, filters)

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

	transportationOfficeID := payload.TransportationOfficeID.String()
	if transportationOfficeID != uuid.Nil.String() && transportationOfficeID != "" {
		transportationIDFilter := []services.QueryFilter{
			query.NewQueryFilter("id", "=", transportationOfficeID),
		}
		// Use FetchOne to see if we have a transportation office that matches the provided id
		var transportationOffice models.TransportationOffice
		fetchErr := o.builder.FetchOne(appCtx, &transportationOffice, transportationIDFilter)

		if fetchErr != nil {
			return nil, nil, fetchErr
		}

		foundUser.TransportationOfficeID = uuid.FromStringOrNil(transportationOfficeID)
	}

	verrs, err := o.builder.UpdateOne(appCtx, &foundUser, nil)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return &foundUser, nil, nil
}

// NewOfficeUserUpdater returns a new office user updater
func NewOfficeUserUpdater(builder officeUserQueryBuilder) services.OfficeUserUpdater {
	return &officeUserUpdater{builder}
}
