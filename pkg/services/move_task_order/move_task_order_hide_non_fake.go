package movetaskorder

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	fakedata "github.com/transcom/mymove/pkg/fakedata_approved"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type invalidReasonsType map[string]string

type moveTaskOrderHider struct {
}

// NewMoveTaskOrderHider creates a new struct with the service dependencies
func NewMoveTaskOrderHider() services.MoveTaskOrderHider {
	return &moveTaskOrderHider{}
}

// Hide hides any MTO that isn't using valid fake data
func (o *moveTaskOrderHider) Hide(appCtx appcontext.AppContext) (services.HiddenMoves, error) {
	invalidMoves := services.HiddenMoves{}
	var mtos models.Moves
	err := appCtx.DB().Q().
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
		return nil, apperror.NewQueryError("Moves", err, fmt.Sprintf("Could not find move task orders: %s", err))
	}

	var invalidFakeMoves models.Moves
	for _, mto := range mtos {
		hide := services.HiddenMove{}
		// TODO: what should we do if there is an error?
		isValid, invalidReasons, _ := isValidFakeModelServiceMember(mto.Orders.ServiceMember)
		if isValid {
			// TODO: what should we do if there is an error?
			var invalidReasons2 invalidReasonsType
			isValid, invalidReasons2, _ = isValidFakeModelMTOShipments(mto.MTOShipments)
			invalidReasons = mergeReasonsMap(invalidReasons, invalidReasons2)
		}

		if !isValid {
			mto.Show = swag.Bool(false)
			invalidFakeMoves = append(invalidFakeMoves, mto)
			reasonsJSONString, jsonErr := json.Marshal(invalidReasons)
			hide.MTOID = mto.ID
			if jsonErr != nil {
				hide.Reason = "json.Marshal to string failed"
			} else {
				hide.Reason = string(reasonsJSONString)
			}
			invalidMoves = append(invalidMoves, hide)
		}
	}

	// TODO: Should we be doing this in a transaction?
	for i := range invalidFakeMoves {
		// Take the address of the slice element to avoid implicit memory aliasing of items from a range statement.
		mto := invalidFakeMoves[i]

		verrs, updateErr := appCtx.DB().ValidateAndUpdate(&mto)
		if verrs != nil && verrs.HasAny() {
			return nil, apperror.NewInvalidInputError(mto.ID, err, verrs, "")
		}
		if updateErr != nil {
			return nil, apperror.NewQueryError("Move", err, fmt.Sprintf("Unexpected error when saving move: %v", err))
		}
	}

	return invalidMoves, nil
}

func mergeReasonsMap(r1 invalidReasonsType, r2 invalidReasonsType) invalidReasonsType {
	for k, v := range r2 {
		r1[k] = v
	}
	return r1
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
		if !ok {
			return false, nil
		}
	}

	if a.Phone != nil {
		ok, err := fakedata.IsValidFakeDataPhone(*a.Phone)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	if a.Email != nil {
		ok, err := fakedata.IsValidFakeDataEmail(*a.Email)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}

func isValidFakeModelMTOShipments(shipments models.MTOShipments) (bool, invalidReasonsType, error) {
	invalidReasons := invalidReasonsType{}

	for _, shipment := range shipments {
		var ok bool
		var err error
		ok, invalidReasons, err = isValidFakeModelMTOShipment(shipment)
		if err != nil {
			return false, invalidReasons, err
		}
		if !ok {
			return false, invalidReasons, nil
		}
	}
	return true, invalidReasons, nil
}

func isValidFakeModelMTOShipment(s models.MTOShipment) (bool, invalidReasonsType, error) {
	invalidReasons := invalidReasonsType{}
	if s.PickupAddress != nil {
		ok, err := isValidFakeModelAddress(s.PickupAddress)
		if err != nil {
			return false, invalidReasons, err
		}
		if !ok {
			invalidReasons["mtoshipment.pickupaddress"] = s.PickupAddress.StreetAddress1
			return false, invalidReasons, nil
		}
	}

	if s.SecondaryPickupAddress != nil {
		ok, err := isValidFakeModelAddress(s.SecondaryPickupAddress)
		if err != nil {
			return false, invalidReasons, err
		}
		if !ok {
			invalidReasons["mtoshipment.secondarypickupaddress"] = s.SecondaryPickupAddress.StreetAddress1
			return false, invalidReasons, nil
		}
	}

	if s.DestinationAddress != nil {
		ok, err := isValidFakeModelAddress(s.DestinationAddress)
		if err != nil {
			return false, invalidReasons, err
		}
		if !ok {
			invalidReasons["mtoshipment.destinationaddress"] = s.DestinationAddress.StreetAddress1
			return false, invalidReasons, nil
		}
	}

	if s.SecondaryDeliveryAddress != nil {
		ok, err := isValidFakeModelAddress(s.SecondaryDeliveryAddress)
		if err != nil {
			return false, invalidReasons, err
		}
		if !ok {
			invalidReasons["mtoshipment.secondarydeliveryaddress"] = s.SecondaryDeliveryAddress.StreetAddress1
			return false, invalidReasons, nil
		}
	}

	// may have to load MTOAgents
	for _, agent := range s.MTOAgents {
		ok, err := isValidFakeModelMTOAgent(agent)
		if err != nil {
			return false, invalidReasons, err
		}
		if !ok {
			invalidReasons["mtoshipment.agent"] = "agent first name and/or last name is invalid"
			return false, invalidReasons, nil
		}
	}

	return true, invalidReasons, nil
}

func isValidFakeModelBackupContact(bc models.BackupContact) (bool, error) {
	ok, err := fakedata.IsValidFakeDataName(bc.Name)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	ok, err = fakedata.IsValidFakeDataEmail(bc.Email)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	if bc.Phone != nil {
		ok, err := fakedata.IsValidFakeDataPhone(*bc.Phone)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}

// isValidFakeModelServiceMember - checks if the contact info
// of a service member is fake
func isValidFakeModelServiceMember(sm models.ServiceMember) (bool, invalidReasonsType, error) {
	invalidReasons := invalidReasonsType{}
	email := sm.PersonalEmail
	if email != nil {
		isValidFakeEmail, _ := fakedata.IsValidFakeDataEmail(*email)
		if !isValidFakeEmail {
			invalidReasons["servicemember.email"] = *email
			return false, invalidReasons, nil
		}
	}
	phone := sm.Telephone
	if phone != nil {
		isValidFakePhone, _ := fakedata.IsValidFakeDataPhone(*phone)
		if !isValidFakePhone {
			invalidReasons["servicemember.phone"] = *phone
			return false, invalidReasons, nil
		}
	}
	secondaryPhone := sm.SecondaryTelephone
	if secondaryPhone != nil {
		isValidFakeSecondaryPhone, _ := fakedata.IsValidFakeDataPhone(*secondaryPhone)
		if !isValidFakeSecondaryPhone {
			invalidReasons["servicemember.phone2"] = *secondaryPhone
			return false, invalidReasons, nil
		}
	}
	ok, err := isValidFakeModelAddress(sm.ResidentialAddress)
	if err != nil {
		return false, invalidReasons, err
	}
	if !ok {
		invalidReasons["servicemember.residentialaddress"] = sm.ResidentialAddress.StreetAddress1
		return false, invalidReasons, nil
	}

	ok, err = isValidFakeModelAddress(sm.BackupMailingAddress)
	if err != nil {
		return false, invalidReasons, err
	}
	if !ok {
		invalidReasons["servicemember.backupmailingaddress"] = sm.BackupMailingAddress.StreetAddress1
		return false, invalidReasons, nil
	}

	fName := sm.FirstName
	lName := sm.LastName
	if fName != nil && lName != nil {
		isValidFakeName, _ := fakedata.IsValidFakeDataFullName(*fName, *lName)
		if !isValidFakeName {
			invalidReasons["servicemember.fullname"] = *fName + " " + *lName
			return false, invalidReasons, nil
		}
	}

	// might need to load BackupContacts
	for _, backupContact := range sm.BackupContacts {
		ok, err = isValidFakeModelBackupContact(backupContact)
		if err != nil {
			return false, invalidReasons, err
		}
		if !ok {
			invalidReasons["servicemember.backupcontact"] = "name, email, or phone found to be invalid"
			return false, invalidReasons, nil
		}
	}

	return true, invalidReasons, nil
}
