package accesscode

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// claimAccessCode is a service object to validate an access code.
type claimAccessCode struct {
	db *pop.Connection
}

// NewAccessCodeClaimer creates a new struct with the service dependencies
func NewAccessCodeClaimer(db *pop.Connection) services.AccessCodeClaimer {
	return &claimAccessCode{db}
}

// fetchAccessCode gets an access code based upon the code given to determine whether or not it is a used code
func (v claimAccessCode) fetchAccessCode(code string) (*models.AccessCode, error) {
	ac := models.AccessCode{}

	err := v.db.
		Where("code = ?", code).
		First(&ac)

	if err != nil {
		return &ac, err
	}

	return &ac, nil
}

// ClaimAccessCode validates an access code based upon the code and move type. A valid access
// code is assumed to have no `service_member_id`
func (v claimAccessCode) ClaimAccessCode(code string, serviceMemberID uuid.UUID) (*models.AccessCode, error) {
	accessCode, err := v.fetchAccessCode(code)

	if err != nil {
		return accessCode, errors.Wrap(err, "Unable to find access code")
	}

	if accessCode.ServiceMemberID != nil {
		return accessCode, errors.New("Access code already claimed")
	}

	transactionErr := v.db.Transaction(func(connection *pop.Connection) error {
		claimedAtTime := time.Now()
		accessCode.ClaimedAt = &claimedAtTime
		accessCode.ServiceMemberID = &serviceMemberID

		verrs, err := connection.ValidateAndSave(accessCode)
		if err != nil || verrs.HasAny() {
			return errors.New("error claiming access code")
		}

		return nil
	})

	if transactionErr != nil {
		return accessCode, errors.Wrap(transactionErr, "Unable to claim access code")
	}

	return accessCode, nil
}
