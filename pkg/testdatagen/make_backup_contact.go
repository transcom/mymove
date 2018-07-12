package testdatagen

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeBackupContact creates a single BackupContact and associated service member.
func MakeBackupContact(db *pop.Connection, assertions Assertions) models.BackupContact {
	serviceMember := assertions.BackupContact.ServiceMember
	if isZeroUUID(assertions.BackupContact.ServiceMemberID) {
		serviceMember = MakeServiceMember(db, assertions)
	}

	backupContact := models.BackupContact{
		ServiceMember:   serviceMember,
		ServiceMemberID: serviceMember.ID,
		Name:            "name",
		Email:           "email@example.com",
		Phone:           swag.String("555-555-5555"),
		Permission:      internalmessages.BackupContactPermissionEDIT,
	}

	mergeModels(&backupContact, assertions.BackupContact)

	mustCreate(db, &backupContact)

	return backupContact
}

// MakeDefaultBackupContact returns a BackupContact with default values
func MakeDefaultBackupContact(db *pop.Connection) models.BackupContact {
	return MakeBackupContact(db, Assertions{})
}

// MakeBackupContactData created 5 BackupContacts (and in turn a User for each)
func MakeBackupContactData(db *pop.Connection) {
	for i := 0; i < 5; i++ {
		MakeDefaultBackupContact(db)
	}
}
