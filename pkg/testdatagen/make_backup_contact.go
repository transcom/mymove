package testdatagen

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakeBackupContact creates a single BackupContact and associated service member.
func MakeBackupContact(db *pop.Connection, assertions Assertions) models.BackupContact {
	serviceMember := assertions.BackupContact.ServiceMember
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.BackupContact.ServiceMemberID) {
		serviceMember = MakeServiceMember(db, assertions)
	}

	backupContact := models.BackupContact{
		ServiceMember:   serviceMember,
		ServiceMemberID: serviceMember.ID,
		Name:            "name",
		Email:           "email@example.com",
		Phone:           swag.String("555-555-5555"),
		Permission:      models.BackupContactPermissionEDIT,
	}

	mergeModels(&backupContact, assertions.BackupContact)

	mustCreate(db, &backupContact, assertions.Stub)

	return backupContact
}

// MakeDefaultBackupContact returns a BackupContact with default values
func MakeDefaultBackupContact(db *pop.Connection) models.BackupContact {
	return MakeBackupContact(db, Assertions{})
}
