package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildNotification creates a single Notification.
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildNotification(db *pop.Connection, customs []Customization, traits []Trait) models.Notification {
	customs = setupCustomizations(customs, traits)

	var cNotification models.Notification
	if result := findValidCustomization(customs, Notification); result != nil {
		cNotification = result.Model.(models.Notification)
		if result.LinkOnly {
			return cNotification
		}
	}

	serviceMember := BuildServiceMember(db, customs, traits)

	// Create default Notification
	notification := models.Notification{
		ServiceMemberID: serviceMember.ID,
		ServiceMember:   serviceMember,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&notification, cNotification)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &notification)
	}

	return notification
}
