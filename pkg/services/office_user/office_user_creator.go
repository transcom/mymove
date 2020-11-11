package officeuser

import (
	"strings"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeUserCreator struct {
	db      *pop.Connection
	builder officeUserQueryBuilder
}

// CreateOfficeUser creates office users
func (o *officeUserCreator) CreateOfficeUser(officeUser *models.OfficeUser, transportationIDFilter []services.QueryFilter) (*models.OfficeUser, *validate.Errors, error) {
	// Use FetchOne to see if we have a transportation office that matches the provided id
	var transportationOffice models.TransportationOffice
	fetchErr := o.builder.FetchOne(&transportationOffice, transportationIDFilter)

	if fetchErr != nil {
		return nil, nil, fetchErr
	}

	user := &models.User{
		LoginGovEmail: strings.ToLower(officeUser.Email),
		Active:        true,
	}
	var verrs *validate.Errors
	var err error
	// We don't want to be left with a user record and no office user so setup a transaction to rollback
	txErr := o.db.Transaction(func(connection *pop.Connection) error {

		verrs, err = o.builder.CreateOne(user)
		if verrs != nil || err != nil {
			return err
		}

		officeUser.UserID = &user.ID
		officeUser.User = *user

		verrs, err = o.builder.CreateOne(officeUser)
		if verrs != nil || err != nil {
			return err
		}

		return nil
	})

	if verrs != nil || txErr != nil {
		return nil, verrs, txErr
	}

	return officeUser, nil, nil
}

// NewOfficeUserCreator returns a new office user creator
func NewOfficeUserCreator(db *pop.Connection, builder officeUserQueryBuilder) services.OfficeUserCreator {
	return &officeUserCreator{db, builder}
}
