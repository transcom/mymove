package officeuser

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/random"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type officeUserCreator struct {
	builder officeUserQueryBuilder
	sender  notifications.NotificationSender
}

// CreateOfficeUser creates office users
func (o *officeUserCreator) CreateOfficeUser(
	appCtx appcontext.AppContext,
	officeUser *models.OfficeUser,
	transportationIDFilter []services.QueryFilter,
) (*models.OfficeUser, *validate.Errors, error) {
	// Use FetchOne to see if we have a transportation office that matches the provided id
	var transportationOffice models.TransportationOffice
	fetchErr := o.builder.FetchOne(appCtx, &transportationOffice, transportationIDFilter)

	if fetchErr != nil {
		return nil, nil, fetchErr
	}

	// Check and update rejected office user if necessary
	existingOfficeUser, rejectedUserVerrs, rejectedUserErr := o.checkAndUpdateRejectedOfficeUser(appCtx, officeUser)
	if rejectedUserVerrs != nil || rejectedUserErr != nil {
		return nil, rejectedUserVerrs, rejectedUserErr
	}

	// If an existing rejected user was updated, return the updated user
	if existingOfficeUser != nil {
		return existingOfficeUser, nil, nil
	}

	// A user may already exist with that email from a previous user (admin, service member, ...)
	var user models.User
	userEmailFilter := query.NewQueryFilter("okta_email", "=", officeUser.Email)
	fetchErr = o.builder.FetchOne(appCtx, &user, []services.QueryFilter{userEmailFilter})

	if fetchErr != nil {
		user = models.User{
			OktaEmail: strings.ToLower(officeUser.Email),
			Active:    true,
		}

		sess := appCtx.Session()

		if sess.IDToken == "devlocal" || appCtx.Session().Hostname == "officelocal" {
			// in devlocal we generate a random okta_id for accounts to use
			user.OktaID = GenerateFakeOktaID()
		}
	}

	var verrs *validate.Errors
	var err error
	var userActivityEmail notifications.Notification
	// We don't want to be left with a user record and no office user so setup a transaction to rollback
	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if user.ID == uuid.Nil {
			verrs, err = o.builder.CreateOne(txnAppCtx, &user)
			if verrs != nil || err != nil {
				return err
			}

			email, emailErr := notifications.NewUserAccountCreated(
				appCtx, notifications.GetSysAdminEmail(o.sender), user.ID, user.UpdatedAt)
			if emailErr != nil {
				return emailErr
			}
			userActivityEmail = notifications.Notification(email)
		}

		officeUser.UserID = &user.ID
		officeUser.User = user

		verrs, err = o.builder.CreateOne(txnAppCtx, officeUser)
		if verrs != nil || err != nil {
			if err != nil {
				if verrs == nil {
					verrs = validate.NewErrors()
				}

				switch err.Error() {
				// If these cases are hit, it is not a true internal server error. Instead, verrs should be appended
				case models.UniqueConstraintViolationOfficeUserEmailErrorString:
					verrs.Add("email", fmt.Sprintf("The email %s is already in use.", officeUser.Email))
					return err

				case models.UniqueConstraintViolationOfficeUserEdipiErrorString:
					// Nil check
					if officeUser.EDIPI != nil {
						verrs.Add("edipi", fmt.Sprintf("The DODID# %s is already in use.", *officeUser.EDIPI))
					} else {
						verrs.Add("edipi", "The DODID# is required, not provided, and appears to already exist in our database.")
					}
					return err

				case models.UniqueConstraintViolationOfficeUserOtherUniqueIDErrorString:
					// Nil check
					if officeUser.OtherUniqueID != nil {
						verrs.Add("other_unique_id", fmt.Sprintf("The other unique ID %s is already in use.", *officeUser.OtherUniqueID))
					} else {
						verrs.Add("other_unique_id", "The other unique ID is required, not provided, and appears to already exist in our database.")
					}
					return err
				}
			}

			return err
		}

		return nil
	})

	if verrs != nil || txErr != nil {
		return nil, verrs, txErr
	}

	if userActivityEmail != nil {
		err = o.sender.SendNotification(appCtx, userActivityEmail)
		if err != nil {
			return nil, nil, err
		}
	}

	return officeUser, nil, nil
}

func GenerateFakeOktaID() string {
	const ID_LEN = 20
	const CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fakeOktaID := ""

	for range [ID_LEN]int{} {
		index, err := random.GetRandomInt(len(CHARSET))

		if err != nil {
			return ""
		}

		fakeOktaID += string(CHARSET[index])
	}

	return fakeOktaID
}

// NewOfficeUserCreator returns a new office user creator
func NewOfficeUserCreator(builder officeUserQueryBuilder, sender notifications.NotificationSender) services.OfficeUserCreator {
	return &officeUserCreator{builder, sender}
}

func (o *officeUserCreator) checkAndUpdateRejectedOfficeUser(
	appCtx appcontext.AppContext,
	officeUser *models.OfficeUser,
) (*models.OfficeUser, *validate.Errors, error) {

	// checking if the office user currently exists and has a previous status of rejected
	var requestedOfficeUser models.OfficeUser
	previouslyRejectedCheck := query.NewQueryFilter("email", "=", officeUser.Email)
	fetchErr := o.builder.FetchOne(appCtx, &requestedOfficeUser, []services.QueryFilter{previouslyRejectedCheck})
	if fetchErr != nil && fetchErr != sql.ErrNoRows {
		return nil, nil, fetchErr // Return the actual error if it's not a "no rows" error
	} else if fetchErr == sql.ErrNoRows {
		// If no rows were found, then we can skip this check
		return nil, nil, nil
	}

	// If the office user exists and was previously rejected, update the status to REQUESTED as well as any new info
	if requestedOfficeUser.ID != uuid.Nil && *requestedOfficeUser.Status == models.OfficeUserStatusREJECTED {
		// Except do not reflect rejection reason or status, we want to manage those via this service func
		reflectedRequestedOfficeUser := reflect.ValueOf(&requestedOfficeUser).Elem()
		reflectedCurrentOfficeUser := reflect.ValueOf(officeUser).Elem()
		requestedOfficeUserType := reflectedRequestedOfficeUser.Type() // Retrieve type from reflect val
		// Iterate over struct fields
		for i := 0; i < reflectedRequestedOfficeUser.NumField(); i++ {
			fieldName := requestedOfficeUserType.Field(i).Name
			// Retrieve the fields
			requestedField := reflectedRequestedOfficeUser.Field(i)
			if !requestedField.CanSet() {
				continue
			}
			currentField := reflectedCurrentOfficeUser.Field(i)
			// Fields we don't want to change
			if fieldName == "ID" || fieldName == "UserID" || fieldName == "User" || fieldName == "TransportationOffice" || fieldName == "CreatedAt" || fieldName == "UpdatedAt" {
				continue
			}
			// Set custom vals for status and rejection reason
			if fieldName == "Status" {
				requestedStatus := models.OfficeUserStatusREQUESTED
				requestedField.Set(reflect.ValueOf(&requestedStatus))
				continue
			}
			if fieldName == "RejectionReason" {
				// Reset the rejection reason
				requestedField.Set(reflect.Zero(requestedField.Type()))
				continue
			}
			if !reflect.DeepEqual(requestedField.Interface(), currentField.Interface()) {
				// A mismatched field has been found, set it to the original office user field
				requestedField.Set(currentField)
			}
		}
		verrs, err := o.builder.UpdateOne(appCtx, &requestedOfficeUser, nil)
		if verrs != nil || err != nil {
			return nil, verrs, err
		}
		return &requestedOfficeUser, nil, nil
	}

	return nil, nil, nil
}
