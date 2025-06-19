package factory

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildBackupContact creates a single BackupContact.
// Also creates, if not provided
// - ServiceMember
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildBackupContact(db *pop.Connection, customs []Customization, traits []Trait) models.BackupContact {
	customs = setupCustomizations(customs, traits)

	// Find BackupContact assertion and convert to models.BackupContact
	var cBackupContact models.BackupContact
	if result := findValidCustomization(customs, BackupContact); result != nil {
		cBackupContact = result.Model.(models.BackupContact)
		if result.LinkOnly {
			return cBackupContact
		}
	}

	// Find/create the ServiceMember model
	serviceMember := BuildServiceMember(db, customs, traits)

	// Create backupContact
	backupContact := models.BackupContact{
		ServiceMemberID: serviceMember.ID,
		ServiceMember:   serviceMember,
		Permission:      models.BackupContactPermissionEDIT,
		FirstName:       "firstName",
		LastName:        "lastName",
		Email:           "email@example.com",
		Phone:           "555-555-5555",
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&backupContact, cBackupContact)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &backupContact)
	}
	return backupContact
}
