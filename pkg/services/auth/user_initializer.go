package auth

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/markbates/goth"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type userInitializer struct {
	DB *pop.Connection
}

// NewUserInitializer returns a new userInitializer struct
func NewUserInitializer(db *pop.Connection) services.UserInitializer {
	return &userInitializer{DB: db}
}

func (c *userInitializer) associateUserModels(userID string, email string, officeUser *models.OfficeUser, tspUser *models.TspUser) (*models.User, *validate.Errors, error) {
	var newUser *models.User
	var responseError error
	responseVErrors := validate.NewErrors()

	c.DB.Transaction(func(tx *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		user, verrs, err := models.CreateUserIfNotExists(tx, userID, email)
		if verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error creating base user")
			return transactionError
		}
		newUser = user

		if officeUser != nil {
			officeUser.UserID = &user.ID
			if verrs, err := tx.ValidateAndUpdate(officeUser); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error updating office user")
				return transactionError
			}
		}

		if tspUser != nil {
			tspUser.UserID = &user.ID
			if verrs, err := tx.ValidateAndUpdate(tspUser); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error updating TSP user")
				return transactionError
			}
		}

		return nil
	})

	return newUser, responseVErrors, responseError
}

// InitializeUser creates a User that is associated to existing Office and TSP user identities
func (c *userInitializer) InitializeUser(openIDUser goth.User) (*models.UserIdentity, error) {
	officeUser, err := models.FetchOfficeUserByEmail(c.DB, openIDUser.Email)
	if err != nil && err != models.ErrFetchNotFound {
		err = errors.Wrap(err, "Error while fetching office user")
		return nil, err
	}

	tspUser, err := models.FetchTspUserByEmail(c.DB, openIDUser.Email)
	if err != nil && err != models.ErrFetchNotFound {
		err = errors.Wrap(err, "Error while fetching TSP user")
		return nil, err
	}

	_, verrs, err := c.associateUserModels(openIDUser.UserID, openIDUser.Email, officeUser, tspUser)
	if err != nil {
		err = errors.Wrap(err, "Error while creating base user")
		return nil, err
	} else if verrs.HasAny() {
		return nil, verrs
	}

	userIdentity, err := models.FetchUserIdentity(c.DB, openIDUser.UserID)
	if err != nil {
		err = errors.Wrap(err, "A user identity could not be found after creating a user")
		return nil, err
	}

	return userIdentity, nil
}
