package adminuser

import (
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type requestedOfficeUserUpdater struct {
	builder requestedOfficeUserQueryBuilder
}

func (o *requestedOfficeUserUpdater) UpdateRequestedOfficeUser(appCtx appcontext.AppContext, id uuid.UUID, payload *adminmessages.RequestedOfficeUserUpdate) (*models.OfficeUser, *validate.Errors, error) {
	updateUserEmail := false
	var officeUser models.OfficeUser
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(appCtx, &officeUser, filters)
	if err != nil {
		return nil, nil, err
	}

	if payload.Email != nil {
		officeUser.Email = *payload.Email
		updateUserEmail = true
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

	rejectedOn := time.Now()
	if payload.RejectionReason != "" {
		officeUser.RejectionReason = &payload.RejectionReason
		officeUser.RejectedOn = &rejectedOn
	}

	if payload.Status != "" {
		officeUser.Status = (*models.OfficeUserStatus)(&payload.Status)
		if *officeUser.Status == models.OfficeUserStatusAPPROVED {
			officeUser.Active = true
		}
	}

	transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := o.builder.UpdateOne(appCtx, &officeUser, nil)
		if err != nil {
			return err
		}
		if verrs != nil && verrs.HasAny() {
			return verrs
		}

		// if the email is being updated, we need to also update the user email & okta email
		if updateUserEmail {
			var existingUser models.User
			filters := []services.QueryFilter{query.NewQueryFilter("id", "=", officeUser.UserID.String())}
			err := o.builder.FetchOne(appCtx, &existingUser, filters)
			if err != nil {
				return err
			}

			existingUser.OktaEmail = officeUser.Email
			verrs, err := o.builder.UpdateOne(txnCtx, &existingUser, nil)
			if err != nil {
				return err
			}
			if verrs != nil && verrs.HasAny() {
				return verrs
			}

			// requested office users will likely not have Okta accounts yet, but we still need to check the edge case
			if existingUser.OktaID != "" && appCtx.Session().IDToken != "devlocal" {
				apiKey := models.GetOktaAPIKey()
				oktaID := existingUser.OktaID
				req := appCtx.HTTPRequest()
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
					return fmt.Errorf("okta user cannot be nil before updating okta email of requested office user")
				}
				existingOktaUser.Profile.Email = existingUser.OktaEmail
				existingOktaUser.Profile.Login = existingUser.OktaEmail

				_, err = models.UpdateOktaUser(appCtx, provider, oktaID, apiKey, *existingOktaUser)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if transactionError != nil {
		return nil, nil, transactionError
	}

	return &officeUser, nil, nil
}

// NewRequestedOfficeUserUpdater returns a new requested office user updater builder
func NewRequestedOfficeUserUpdater(builder requestedOfficeUserQueryBuilder) services.RequestedOfficeUserUpdater {
	return &requestedOfficeUserUpdater{builder}
}
