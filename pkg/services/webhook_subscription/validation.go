package webhooksubscription

import (
	"database/sql"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

// checks business requirements
// webhookSubscriptionValidator describes a method for checking business requirements
type webhookSubscriptionValidator interface {
	// Validate checks the newWebhookSubscription for adherence to business rules.
	// It is safe to return a *validate.Errors with zero added errors as
	// a success case.
	Validate(a appcontext.AppContext, newWebhookSubscription models.WebhookSubscription) error
}

// validateWebhookSubscription checks a webhookSubscription against a passed-in set of business rule checks
func validateWebhookSubscription(
	appCtx appcontext.AppContext,
	newWebhookSubscription models.WebhookSubscription,
	checks ...webhookSubscriptionValidator,
) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, newWebhookSubscription); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				// accumulate validation errors
				verrs.Append(e)
			default:
				// non-validation errors have priority,
				// and short-circuit doing any further checks
				return err
			}
		}
	}
	if verrs.HasAny() {
		result = apperror.NewInvalidInputError(newWebhookSubscription.ID, nil, verrs, "Invalid input found while validating the webhookSubscription.")
	}
	return result
}

// webhookSubscriptionValidatorFunc is an adapter type for converting a function into an implementation of webhookSubscriptionValidator
type webhookSubscriptionValidatorFunc func(appcontext.AppContext, models.WebhookSubscription) error

// Validate fulfills the webhookSubscriptionValidator interface
func (fn webhookSubscriptionValidatorFunc) Validate(appCtx appcontext.AppContext, newer models.WebhookSubscription) error {
	return fn(appCtx, newer)
}

// checkSubscriberExists checks that a webhook's subscriber corresponds to an existing contractor
func checkSubscriberExists(builder webhookSubscriptionQueryBuilder) webhookSubscriptionValidator {
	return webhookSubscriptionValidatorFunc(func(a appcontext.AppContext, newWebhookSubscription models.WebhookSubscription) error {
		subscriberIDFilter := []services.QueryFilter{
			query.NewQueryFilter("id", "=", newWebhookSubscription.SubscriberID),
		}
		var contractor models.Contractor
		fetchErr := builder.FetchOne(a, &contractor, subscriberIDFilter)
		if fetchErr != nil {
			switch fetchErr {
			case sql.ErrNoRows:
				return apperror.NewNotFoundError(newWebhookSubscription.SubscriberID, "while looking for SubscriberID")
			default:
				return apperror.NewQueryError("Contractor", fetchErr, "")
			}
		}
		return nil
	})
}
