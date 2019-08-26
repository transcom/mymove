package storageintransit

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/apimessages"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

func authorizeStorageInTransitHTTPRequest(db *pop.Connection, session *auth.Session, shipmentID uuid.UUID, allowOffice bool) (isUserAuthorized bool, err error) {
	if session.IsTspUser() {
		_, _, err := models.FetchShipmentForVerifiedTSPUser(db, session.TspUserID, shipmentID)

		if err != nil {
			return false, err
		}
		return true, nil
	}

	if session.IsOfficeUser() {
		if allowOffice {
			return true, nil
		}
	}
	return false, models.ErrFetchForbidden
}

// Stole this local function from publicapi/addresses.go
// TODO: How should we handle this on the services side? We should figure out how these shared resources should operate.
func updateAddressWithPayload(a *models.Address, payload *apimessages.Address) {
	a.StreetAddress1 = *payload.StreetAddress1
	a.StreetAddress2 = payload.StreetAddress2
	a.StreetAddress3 = payload.StreetAddress3
	a.City = *payload.City
	a.State = *payload.State
	a.PostalCode = *payload.PostalCode
	a.Country = payload.Country
}
