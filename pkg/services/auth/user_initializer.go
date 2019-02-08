package auth

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/markbates/goth"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// userInitializer is a service object to deliver and price a Shipment
type userInitializer struct {
	DB *pop.Connection
}

// NewUserInitializer returns a new userInitializer struct
func NewUserInitializer(db *pop.Connection) services.UserInitializer {
	return &userInitializer{DB: db}
}

func (c *userInitializer) createBaseUser(userID string, email string, officeUser *models.OfficeUser, tspUser *models.TspUser) (*models.User, *validate.Errors, error) {
	var newUser *models.User
	var responseError error
	responseVErrors := validate.NewErrors()

	c.DB.Transaction(func(tx *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		user, verrs, err := models.CreateUser(tx, userID, email)
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

// InitializeUser returns a collection of user data, generating a base user as necessary
func (c *userInitializer) InitializeUser(detector services.AppDetector, openIDUser goth.User) (response services.InitializeUserResponse, verrs *validate.Errors, err error) {
	verrs = validate.NewErrors()

	userIdentity, identityErr := models.FetchUserIdentity(c.DB, openIDUser.UserID)
	if identityErr != nil && identityErr != models.ErrFetchNotFound {
		// An unknown error
		err = errors.Wrap(identityErr, "Unknown error while fetching user identity")
		return
	}

	// We found a user identity
	if identityErr == nil {
		// If we already have the relevant user info then we're done
		if detector.IsOfficeApp() && userIdentity.OfficeUserID != nil {
			response.OfficeUserID = *userIdentity.OfficeUserID
			return
		} else if detector.IsTspApp() && userIdentity.TspUserID != nil {
			response.TspUserID = *userIdentity.TspUserID
			return
		}
	}

	// Else we'll need to look up user information by email address
	var officeUser *models.OfficeUser
	var tspUser *models.TspUser
	if detector.IsOfficeApp() {
		officeUser, err = models.FetchOfficeUserByEmail(c.DB, openIDUser.Email)
		if err != nil {
			err = errors.Wrap(err, "Error while fetching office user")
			return
		}
		response.OfficeUserID = officeUser.ID
	} else if detector.IsTspApp() {
		tspUser, err = models.FetchTspUserByEmail(c.DB, openIDUser.Email)
		if err != nil {
			err = errors.Wrap(err, "Error while fetching TSP user")
			return
		}
		response.TspUserID = tspUser.ID
	}

	// If we never got an identity, we need to generate a base user.
	// Also pass in office/tsp user models so they can be associated
	if identityErr == models.ErrFetchNotFound {
		_, verrs, err = c.createBaseUser(openIDUser.UserID, openIDUser.Email, officeUser, tspUser)
		if verrs.HasAny() || err != nil {
			err = errors.Wrap(err, "Error while creating base user")
			return
		}
		// Fetch the userIdentity again, which should now have a base user
		userIdentity, err = models.FetchUserIdentity(c.DB, openIDUser.UserID)
		if err != nil {
			err = errors.Wrap(err, "A user identity could not be found after creating a user")
			return
		}
	}
	response.UserID = userIdentity.ID
	if userIdentity.ServiceMemberID != nil {
		response.ServiceMemberID = *userIdentity.ServiceMemberID
	}

	response.FirstName = userIdentity.FirstName()
	response.LastName = userIdentity.LastName()
	response.Middle = userIdentity.Middle()

	return response, validate.NewErrors(), nil
}
