package accesscode

import (
	"time"

	"github.com/gobuffalo/validate/v3"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// accessCodeClaimer is a service object to validate an access code.
type accessCodeClaimer struct {
}

// NewAccessCodeClaimer creates a new struct with the service dependencies
func NewAccessCodeClaimer() services.AccessCodeClaimer {
	return &accessCodeClaimer{}
}

// fetchAccessCodeForUpdate gets an access code based upon the code given to determine whether or not it is a used code
func fetchAccessCodeForUpdate(appCtx appcontext.AppContext, code string) (*models.AccessCode, error) {
	ac := models.AccessCode{}

	err := appCtx.DB().RawQuery(`
    SELECT access_codes.claimed_at,
		access_codes.code,
		access_codes.created_at,
		access_codes.id,
		access_codes.move_type,
		access_codes.service_member_id
	FROM access_codes AS access_codes
	WHERE code = $1 FOR UPDATE`, code).First(&ac)

	if err != nil {
		return &ac, err
	}

	return &ac, nil
}

// ClaimAccessCode validates an access code based upon the code and move type. A valid access
// code is assumed to have no `service_member_id`
func (v accessCodeClaimer) ClaimAccessCode(appCtx appcontext.AppContext, code string, serviceMemberID uuid.UUID) (*models.AccessCode, *validate.Errors, error) {
	var accessCode *models.AccessCode
	var err error
	verrs := validate.NewErrors()

	transactionErr := appCtx.NewTransaction(func(txnAppCfg appcontext.AppContext) error {
		accessCode, err = fetchAccessCodeForUpdate(txnAppCfg, code)

		if err != nil {
			return errors.Wrap(err, "Unable to find access code")
		}

		if accessCode.ServiceMemberID != nil {
			return errors.New("Access code already claimed")
		}

		claimedAtTime := time.Now()
		accessCode.ClaimedAt = &claimedAtTime
		accessCode.ServiceMemberID = &serviceMemberID

		verrs, err = txnAppCfg.DB().ValidateAndSave(accessCode)
		if err != nil || verrs.HasAny() {
			return errors.New("error claiming access code")
		}

		return nil
	})

	if transactionErr != nil || verrs.HasAny() {
		return accessCode, verrs, transactionErr
	}

	return accessCode, verrs, nil
}
