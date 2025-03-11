package officeuser

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type officeUserUpdater struct {
	builder officeUserQueryBuilder
}

// UpdateOfficeUser updates an office user
func (o *officeUserUpdater) UpdateOfficeUser(appCtx appcontext.AppContext, id uuid.UUID, payload *adminmessages.OfficeUserUpdate, primaryTransportationOfficeID uuid.UUID) (*models.OfficeUser, *validate.Errors, error) {
	updateUserAndOkta := false
	var foundUser models.OfficeUser
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(appCtx, &foundUser, filters)
	if err != nil {
		return nil, nil, err
	}

	if payload.Email != nil {
		foundUser.Email = *payload.Email
		updateUserAndOkta = true
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

	transportationOfficeID := primaryTransportationOfficeID.String()
	if primaryTransportationOfficeID != uuid.Nil && transportationOfficeID != uuid.Nil.String() && transportationOfficeID != "" {
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

	transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := o.builder.UpdateOne(txnCtx, &foundUser, nil)
		if err != nil {
			return err
		}
		if verrs != nil && verrs.HasAny() {
			return verrs
		}

		// if the email is being updated, we need to also update the user email & okta email
		if updateUserAndOkta {
			var existingUser models.User
			filters := []services.QueryFilter{query.NewQueryFilter("id", "=", foundUser.UserID.String())}
			err := o.builder.FetchOne(appCtx, &existingUser, filters)
			if err != nil {
				return err
			}

			existingUser.OktaEmail = foundUser.Email
			verrs, err := o.builder.UpdateOne(txnCtx, &existingUser, nil)
			if err != nil {
				return err
			}
			if verrs != nil && verrs.HasAny() {
				return verrs
			}

			apiKey := models.GetOktaAPIKey()
			oktaID := existingUser.OktaID
			req := appCtx.Session().HTTPRequest
			if req == nil {
				return fmt.Errorf("failed to retrieve HTTP request from session")
			}

			// Use the HTTP request to get the Okta provider
			provider, err := okta.GetOktaProviderForRequest(req)
			if err != nil {
				return fmt.Errorf("error retrieving Okta provider: %w", err)
			}

			// verifying the okta user exists but we also need all the okta profile info prior to updating
			existingOktaUser, err := models.GetOktaUser(appCtx, provider, oktaID, apiKey)
			if err != nil {
				return fmt.Errorf("error getting Okta user prior to updating: %w", err)
			}
			if existingOktaUser == nil {
				return fmt.Errorf("okta user cannot be nil before updating okta email of office user")
			}
			existingOktaUser.Profile.Email = existingUser.OktaEmail
			existingOktaUser.Profile.Login = existingUser.OktaEmail

			_, err = models.UpdateOktaUser(appCtx, provider, oktaID, apiKey, *existingOktaUser)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if transactionError != nil {
		return nil, nil, transactionError
	}

	return &foundUser, nil, nil
}

// NewOfficeUserUpdater returns a new office user updater
func NewOfficeUserUpdater(builder officeUserQueryBuilder) services.OfficeUserUpdater {
	return &officeUserUpdater{builder}
}
