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

type requestedOfficeUserUpdater struct {
	builder requestedOfficeUserQueryBuilder
}

func (o *requestedOfficeUserUpdater) UpdateRequestedOfficeUser(appCtx appcontext.AppContext, id uuid.UUID, payload *adminmessages.RequestedOfficeUserUpdate) (*models.OfficeUser, *validate.Errors, error) {
	var officeUser models.OfficeUser
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(appCtx, &officeUser, filters)

	if err != nil {
		return nil, nil, err
	}

	if payload.FirstName != nil {
		officeUser.FirstName = *payload.FirstName
	}

	if payload.MiddleInitials != nil {
		officeUser.MiddleInitials = payload.MiddleInitials
	}

	if payload.LastName != nil {
		officeUser.LastName = *payload.LastName
	}

	if payload.Telephone != nil {
		officeUser.Telephone = *payload.Telephone
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

		officeUser.TransportationOfficeID = uuid.FromStringOrNil(transportationOfficeID)
	}

	if payload.Edipi != "" {
		officeUser.EDIPI = &payload.Edipi
	}

	if payload.OtherUniqueID != "" {
		officeUser.OtherUniqueID = &payload.OtherUniqueID
	}

	if payload.RejectionReason != "" {
		officeUser.RejectionReason = &payload.RejectionReason
	}

	if payload.Status != "" {
		officeUser.Status = (*models.OfficeUserStatus)(&payload.Status)
		if *officeUser.Status == models.OfficeUserStatusAPPROVED {
			officeUser.Active = true
		}
	}

	verrs, err := o.builder.UpdateOne(appCtx, &officeUser, nil)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return &officeUser, nil, nil
}

// NewRequestedOfficeUserUpdater returns a new requested office user updater builder
func NewRequestedOfficeUserUpdater(builder requestedOfficeUserQueryBuilder) services.RequestedOfficeUserUpdater {
	return &requestedOfficeUserUpdater{builder}
}
