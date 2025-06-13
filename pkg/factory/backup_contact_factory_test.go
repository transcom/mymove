package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildBackupContact() {
	suite.Run("Successful creation of default BackupContact", func() {
		// Under test:      BuildBackupContact
		// Mocked:          None
		// Set up:          Create a BackupContact with no customizations or traits
		// Expected outcome:BackupContact should be created with default values

		// SETUP
		defaultContact := models.BackupContact{
			Permission: models.BackupContactPermissionEDIT,
			FirstName:  "firstName",
			LastName:   "lastName",
			Email:      "email@example.com",
			Phone:      "555-555-5555",
		}
		defaultServiceMember := models.ServiceMember{
			FirstName: models.StringPointer("Leo"),
			LastName:  models.StringPointer("Spacemen"),
			Telephone: models.StringPointer("212-123-4567"),
		}

		// CALL FUNCTION UNDER TEST
		backupContact := BuildBackupContact(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultContact.Permission, backupContact.Permission)
		suite.Equal(defaultContact.FirstName, backupContact.FirstName)
		suite.Equal(defaultContact.LastName, backupContact.LastName)
		suite.Equal(defaultContact.Email, backupContact.Email)
		suite.Equal(defaultContact.Phone, backupContact.Phone)

		// Check that service member was hooked in
		suite.Equal(*defaultServiceMember.FirstName, *backupContact.ServiceMember.FirstName)
		suite.Equal(*defaultServiceMember.LastName, *backupContact.ServiceMember.LastName)
		suite.Equal(*defaultServiceMember.Telephone, *backupContact.ServiceMember.Telephone)

	})

	suite.Run("Successful creation of customized BackupContact", func() {
		// Under test:      BuildBackupContact
		// Set up:          Create a BackupContact and pass custom fields
		// Expected outcome:BackupContact should be created with custom fields

		// SETUP
		customBackupContact := models.BackupContact{
			ID:         uuid.Must(uuid.NewV4()),
			FirstName:  "Fake",
			LastName:   "Name",
			Email:      "email@example.com",
			Phone:      "555-444-4444",
			Permission: models.BackupContactPermissionVIEW,
		}

		customServiceMember := models.ServiceMember{
			FirstName: models.StringPointer("Jason"),
			LastName:  models.StringPointer("Ash"),
		}

		// CALL FUNCTION UNDER TEST
		backupContact := BuildBackupContact(suite.DB(), []Customization{
			{Model: customBackupContact},
			{Model: customServiceMember},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customBackupContact.ID, backupContact.ID)
		suite.Equal(customBackupContact.FirstName, backupContact.FirstName)
		suite.Equal(customBackupContact.LastName, backupContact.LastName)
		suite.Equal(customBackupContact.Email, backupContact.Email)
		suite.Equal(customBackupContact.Permission, backupContact.Permission)
		suite.Equal(customBackupContact.Phone, backupContact.Phone)

		// Check that the service member was customized
		suite.Equal(*customServiceMember.FirstName, *backupContact.ServiceMember.FirstName)
		suite.Equal(*customServiceMember.LastName, *backupContact.ServiceMember.LastName)
	})

	suite.Run("Successful return of linkOnly BackupContact", func() {
		// Under test:       BuildBackupContact
		// Set up:           Pass in a linkOnly BackupContact
		// Expected outcome: No new BackupContact should be created.

		// Check num BackupContact records
		precount, err := suite.DB().Count(&models.BackupContact{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		backupContact := BuildBackupContact(suite.DB(), []Customization{
			{
				Model: models.BackupContact{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.BackupContact{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, backupContact.ID)

	})
	suite.Run("Successful return of stubbed BackupContact", func() {
		// Under test:       BuildBackupContact
		// Set up:           Create a BackupContact with nil DB
		// Expected outcome: No new BackupContact should be created.

		// Check num BackupContact records
		precount, err := suite.DB().Count(&models.BackupContact{})
		suite.NoError(err)

		// Nil passed in as db
		backupContact := BuildBackupContact(nil, []Customization{
			{
				Model: models.BackupContact{
					FirstName: "Another",
					LastName:  "Name",
				},
			},
		}, nil)
		count, err := suite.DB().Count(&models.BackupContact{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal("Another", backupContact.FirstName)
		suite.Equal("Name", backupContact.LastName)
	})
}
