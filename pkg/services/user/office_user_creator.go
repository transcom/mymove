package user

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeUserCreator struct {
	db      *pop.Connection
	builder officeUserQueryBuilder
}

func (o *officeUserCreator) CreateOfficeUser(user *models.OfficeUser, transporationIDFilter []services.QueryFilter) (*models.OfficeUser, *validate.Errors, error) {
	// Use FetchOne to see if we have a transportation office that matches the provided id
	var transporationOffice models.TransportationOffice
	err := o.builder.FetchOne(&transporationOffice, transporationIDFilter)

	if err != nil {
		return nil, nil, err
	}

	verrs, err := o.db.ValidateAndCreate(user)
	if err != nil || verrs.HasAny() {
		return nil, verrs, err
	}

	return user, nil, nil
}

func NewOfficeUserCreator(db *pop.Connection, builder officeUserQueryBuilder) services.OfficeUserCreator {
	return &officeUserCreator{db, builder}
}