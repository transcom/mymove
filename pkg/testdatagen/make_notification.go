package testdatagen

import (
	"fmt"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeNotification creates a single notification and associated ServiceMember
func MakeNotification(db *pop.Connection, assertions Assertions) models.Notification {
	serviceMember := assertions.Notification.ServiceMember
	// There's a uniqueness constraint on notification emails so add some randomness
	email := fmt.Sprintf("leo_spaceman_office_%s@example.com", makeRandomString(5))

	if isZeroUUID(assertions.Notification.ServiceMemberID) {
		if assertions.User.LoginGovEmail == "" {
			assertions.User.LoginGovEmail = email
		}
		serviceMember = MakeServiceMember(db, assertions)
	}

	notification := models.Notification{
		ServiceMemberID: serviceMember.ID,
		ServiceMember:   serviceMember,
	}

	mergeModels(&notification, assertions.Notification)

	mustCreate(db, &notification, assertions.Stub)

	return notification
}

// MakeDefaultNotification makes an Notification with default values
func MakeDefaultNotification(db *pop.Connection) models.Notification {
	return MakeNotification(db, Assertions{})
}
