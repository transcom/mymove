package testdatagen

import (
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// makeBackupContact creates a single BackupContact and associated
// service member.
//
// Deprecated: use factory.BuildBackupContact
func makeBackupContact(db *pop.Connection, assertions Assertions) models.BackupContact {
	serviceMember := assertions.BackupContact.ServiceMember
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.BackupContact.ServiceMemberID) {
		serviceMember = makeServiceMember(db, assertions)
	}

	backupContact := models.BackupContact{
		ServiceMember:   serviceMember,
		ServiceMemberID: serviceMember.ID,
		FirstName:       "firstName",
		LastName:        "lastName",
		Email:           "email@example.com",
		Phone:           "555-555-5555",
		Permission:      models.BackupContactPermissionEDIT,
	}

	mergeModels(&backupContact, assertions.BackupContact)

	mustCreate(db, &backupContact, assertions.Stub)

	return backupContact
}
