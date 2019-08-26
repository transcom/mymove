package storageintransit

import (
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

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
