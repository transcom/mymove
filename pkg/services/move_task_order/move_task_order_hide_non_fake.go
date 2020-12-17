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
		// Note: We may be able to save some queries if we load on demand, but we'll need to
		// refactor the methods that check for valid fake data to pass in the DB connection.
		Eager(
			"Orders.ServiceMember.ResidentialAddress",
			"Orders.ServiceMember.BackupMailingAddress",
			"Orders.ServiceMember.BackupContacts",
			"MTOShipments.PickupAddress",
			"MTOShipments.DestinationAddress",
			"MTOShipments.SecondaryPickupAddress",
			"MTOShipments.SecondaryDeliveryAddress",
			"MTOShipments.MTOAgents",
		).
		Where("show = ?", swag.Bool(true)).
		All(&mtos)
	if err != nil {
		return nil, services.NewQueryError("Moves", err, fmt.Sprintf("Could not find move task orders: %s", err))
	}

	var invalidFakeMoves models.Moves
	for _, mto := range mtos {
		// TODO: what should we do if there is an error?
		isValid, _ := isValidFakeModelServiceMember(mto.Orders.ServiceMember)
		if isValid {
			// TODO: what should we do if there is an error?
			isValid, _ = isValidFakeModelMTOShipments(mto.MTOShipments)
		}

		if !isValid {
			mto.Show = swag.Bool(false)
			invalidFakeMoves = append(invalidFakeMoves, mto)
		}
	}

	// TODO: Should we be doing this in a transaction?
	for i := range invalidFakeMoves {
		// Take the address of the slice element to avoid implicit memory aliasing of items from a range statement.
		mto := invalidFakeMoves[i]
		verrs, updateErr := o.db.ValidateAndUpdate(&mto)
		if verrs != nil && verrs.HasAny() {
			return nil, services.NewInvalidInputError(mto.ID, err, verrs, "")
		}
		if updateErr != nil {
			return nil, services.NewQueryError("Move", err, fmt.Sprintf("Unexpected error when saving move: %v", err))
		}
	}

	return invalidFakeMoves, nil
}

func isValidFakeModelAddress(a *models.Address) (bool, error) {
	if a != nil {
		return fakedata.IsValidFakeDataAddress(a.StreetAddress1)
	}
	return true, nil
}

func isValidFakeModelMTOAgent(a models.MTOAgent) (bool, error) {
	if a.FirstName != nil && a.LastName != nil {
		ok, err := fakedata.IsValidFakeDataFullName(*a.FirstName, *a.LastName)
		if err != nil {
			return false, err
		}
		if ok == false {
			return false, nil
		}
	}

	if a.Phone != nil {
		ok, err := fakedata.IsValidFakeDataPhone(*a.Phone)
		if err != nil {
			return false, err
		}
		if ok == false {
			return false, nil
		}
	}

	if a.Email != nil {
		ok, err := fakedata.IsValidFakeDataEmail(*a.Email)
		if err != nil {
			return false, err
		}
		if ok == false {
			return false, nil
		}
	}

	return true, nil
}

func isValidFakeModelMTOShipments(shipments models.MTOShipments) (bool, error) {
	for _, shipment := range shipments {
		ok, err := isValidFakeModelMTOShipment(shipment)
		if err != nil {
			return false, err
		}
		if ok == false {
			return false, nil
		}
	}
	return true, nil
}

func isValidFakeModelMTOShipment(s models.MTOShipment) (bool, error) {
	if s.PickupAddress != nil {
		ok, err := isValidFakeModelAddress(s.PickupAddress)
		if err != nil {
			return false, err
		}
		if ok == false {
			return false, nil
		}
	}

	if s.SecondaryPickupAddress != nil {
		ok, err := isValidFakeModelAddress(s.SecondaryPickupAddress)
		if err != nil {
			return false, err
		}
		if ok == false {
			return false, nil
		}
	}

	if s.DestinationAddress != nil {
		ok, err := isValidFakeModelAddress(s.DestinationAddress)
		if err != nil {
			return false, err
		}
		if ok == false {
			return false, nil
		}
	}

	if s.SecondaryDeliveryAddress != nil {
		ok, err := isValidFakeModelAddress(s.SecondaryDeliveryAddress)
		if err != nil {
			return false, err
		}
		if ok == false {
			return false, nil
		}
	}

	// may have to load MTOAgents
	for _, agent := range s.MTOAgents {
		ok, err := isValidFakeModelMTOAgent(agent)
		if err != nil {
			return false, err
		}
		if ok == false {
			return false, nil
		}
	}

	return true, nil
}

func isValidFakeModelBackupContact(bc models.BackupContact) (bool, error) {
	ok, err := fakedata.IsValidFakeDataName(bc.Name)
	if err != nil {
		return false, err
	}
	if ok == false {
		return false, nil
	}

	ok, err = fakedata.IsValidFakeDataEmail(bc.Email)
	if err != nil {
		return false, err
	}
	if ok == false {
		return false, nil
	}

	if bc.Phone != nil {
		ok, err := fakedata.IsValidFakeDataPhone(*bc.Phone)
		if err != nil {
			return false, err
		}
		if ok == false {
			return false, nil
		}
	}

	return true, nil
}

// isValidFakeModelServiceMember - checks if the contact info
// of a service member is fake
func isValidFakeModelServiceMember(sm models.ServiceMember) (bool, error) {
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
	ok, err := isValidFakeModelAddress(sm.ResidentialAddress)
	if err != nil {
		return false, err
	}
	if ok == false {
		return false, nil
	}

	ok, err = isValidFakeModelAddress(sm.BackupMailingAddress)
	if err != nil {
		return false, err
	}
	if ok == false {
		return false, nil
	}

	fName := sm.FirstName
	lName := sm.LastName
	if fName != nil && lName != nil {
		isValidFakeName, _ := fakedata.IsValidFakeDataFullName(*fName, *lName)
		if isValidFakeName == false {
			return false, nil
		}
	}

	// might need to load BackupContacts
	for _, backupContact := range sm.BackupContacts {
		ok, err = isValidFakeModelBackupContact(backupContact)
		if err != nil {
			return false, err
		}
		if ok == false {
			return false, nil
		}
	}

	return true, nil
}
