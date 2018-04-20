package testdatagen

import (
	"log"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeBackupContact creates a single BackupContact and associated service member.
func MakeBackupContact(db *pop.Connection) (models.BackupContact, error) {
	serviceMember, err := MakeServiceMember(db)
	if err != nil {
		return models.BackupContact{}, err
	}

	backupContact := models.BackupContact{
		ServiceMemberID: serviceMember.ID,
		ServiceMember:   serviceMember,
		Name:            "name",
		Email:           "email@example.com",
		Phone:           swag.String("5555555555"),
		Permission:      internalmessages.BackupContactPermissionEDIT,
	}

	verrs, err := db.ValidateAndSave(&backupContact)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	return backupContact, err
}

// MakeBackupContactData created 5 BackupContacts (and in turn a User for each)
func MakeBackupContactData(db *pop.Connection) {
	for i := 0; i < 5; i++ {
		MakeBackupContact(db)
	}
}
