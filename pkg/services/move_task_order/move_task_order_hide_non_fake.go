package movetaskorder

import (
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"

	fakedata "github.com/transcom/mymove/pkg/fakedata_approved"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveTaskOrderHider struct {
	db *pop.Connection
}

// NewMoveTaskOrderHider creates a new struct with the service dependencies
func NewMoveTaskOrderHider(db *pop.Connection) services.MoveTaskOrderHider {
	return &moveTaskOrderHider{db}
}

// Hide hides any MTO that isn't using valid fake data
func (o *moveTaskOrderHider) Hide() (models.Moves, error) {
	var mtos models.Moves
	err := o.db.Q().
		Where("show = ?", swag.Bool(true)).
		All(&mtos)
	if err != nil {
		return nil, services.NewQueryError("Moves", err, fmt.Sprintf("Could not find move task orders: %s", err))
	}

	var invalidFakeMoves models.Moves
	for _, mto := range mtos {
		isValid, _ := isValidFakeServiceMember(mto.Orders.ServiceMember)
		if !isValid {
			dontShow := false
			mto.Show = &dontShow
			invalidFakeMoves = append(invalidFakeMoves, mto)
		}
	}

	return invalidFakeMoves, nil
}

// isValidFakeServiceMember - checks if the contact info
// of a service member is fake
func isValidFakeServiceMember(sm models.ServiceMember) (bool, error) {
	email := sm.PersonalEmail
	if email != nil {
		isValidFakeEmail, _ := fakedata.IsValidFakeDataEmail(*email)
		if !isValidFakeEmail {
			return false, nil
		}
	}
	phone := sm.Telephone
	if phone != nil {
		isValidFakePhone, _ := fakedata.IsValidFakeDataPhone(*phone)
		if isValidFakePhone == false {
			return false, nil
		}
	}
	secondaryPhone := sm.SecondaryTelephone
	if secondaryPhone != nil {
		isValidFakeSecondaryPhone, _ := fakedata.IsValidFakeDataPhone(*secondaryPhone)
		if !isValidFakeSecondaryPhone {
			return false, nil
		}
	}
	if sm.ResidentialAddress != nil {
		address := sm.ResidentialAddress.StreetAddress1
		if address != "" {
			isValidFakeAddress, _ := fakedata.IsValidFakeDataAddress(sm.ResidentialAddress.StreetAddress1)
			if !isValidFakeAddress {
				return false, nil
			}
		}
	}
	if sm.BackupMailingAddress != nil {
		backupAddress := sm.BackupMailingAddress.StreetAddress1
		if backupAddress != "" {
			isValidFakeBackupAddress, _ := fakedata.IsValidFakeDataAddress(sm.BackupMailingAddress.StreetAddress1)
			if !isValidFakeBackupAddress {
				return false, nil
			}
		}
	}

	fName := sm.FirstName
	lName := sm.LastName
	if fName != nil && lName != nil {
		isValidFakeName, _ := fakedata.IsValidFakeDataFullName(*fName, *lName)
		if isValidFakeName == false {
			return false, nil
		}
	}
	return true, nil
}
