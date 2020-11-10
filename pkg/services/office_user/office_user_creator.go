package officeuser

import (
	"strings"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeUserCreator struct {
	builder officeUserQueryBuilder
}

// CreateOfficeUser creates office users
func (o *officeUserCreator) CreateOfficeUser(officeUser *models.OfficeUser, transportationIDFilter []services.QueryFilter) (*models.OfficeUser, *validate.Errors, error) {
	// Use FetchOne to see if we have a transportation office that matches the provided id
	var transportationOffice models.TransportationOffice
	err := o.builder.FetchOne(&transportationOffice, transportationIDFilter)

	if err != nil {
		return nil, nil, err
	}

	user := &models.User{
		LoginGovEmail: strings.ToLower(officeUser.Email),
		Active:        true,
	}
	verrs, err := o.builder.CreateOne(user)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	officeUser.UserID = &user.ID
	officeUser.User = *user

	verrs, err = o.builder.CreateOne(officeUser)
	if verrs != nil || err != nil {
		return nil, verrs, err
	}

	return officeUser, nil, nil
}

// NewOfficeUserCreator returns a new office user creator
func NewOfficeUserCreator(builder officeUserQueryBuilder) services.OfficeUserCreator {
	return &officeUserCreator{builder}
}
