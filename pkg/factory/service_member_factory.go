package factory

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type serviceMemberBuildType byte

const (
	serviceMemberBuildBasic serviceMemberBuildType = iota
	serviceMemberBuildExtended
)

// buildServiceMemberWithBuildType does the actual work
// if buildType is basic, it builds
//   - User
//   - ResidentialAddress
//
// if buildType is extended, it builds
//   - User,
//   - ResidentialAddress
//   - BackupMailingAddress
//   - DutyLocation
//   - BackupContact
//
// basic build type will build User, ResidentialAddress,
// BackupMailingAddress and DutyLocation if and only if a
// customization is provided
//
// a BackupContact is only built for extended
func buildServiceMemberWithBuildType(db *pop.Connection, customs []Customization, traits []Trait, buildType serviceMemberBuildType) models.ServiceMember {
	customs = setupCustomizations(customs, traits)

	// Find ServiceMember customization and extract the custom ServiceMember
	var cServiceMember models.ServiceMember
	if result := findValidCustomization(customs, ServiceMember); result != nil {
		cServiceMember = result.Model.(models.ServiceMember)
		if result.LinkOnly {
			return cServiceMember
		}
	}

	// Find/create the ResidentialAddress
	tempResAddressCustoms := customs
	result := findValidCustomization(customs, Addresses.ResidentialAddress)
	if result != nil {
		tempResAddressCustoms = convertCustomizationInList(tempResAddressCustoms, Addresses.ResidentialAddress, Address)
	}

	resAddress := BuildAddress(db, tempResAddressCustoms, traits)

	// Find/create the user model
	user := BuildUser(db, customs, traits)

	email := "leo_spaceman_sm@example.com"
	agency := models.AffiliationARMY

	// Create random edipi
	randomEdipi := RandomEdipi()

	serviceMember := models.ServiceMember{
		UserID:               user.ID,
		User:                 user,
		Edipi:                models.StringPointer(randomEdipi),
		Affiliation:          &agency,
		FirstName:            models.StringPointer("Leo"),
		LastName:             models.StringPointer("Spacemen"),
		Telephone:            models.StringPointer("212-123-4567"),
		PersonalEmail:        &email,
		ResidentialAddressID: &resAddress.ID,
		ResidentialAddress:   &resAddress,
		CacValidated:         false,
	}

	backupAddressResult := findValidCustomization(customs, Addresses.BackupMailingAddress)
	// Find/create the BackupMailingAddress if customization is
	// provided
	if backupAddressResult != nil {
		backupAddressCustoms := convertCustomizationInList(customs, Addresses.BackupMailingAddress, Address)

		backupAddress := BuildAddress(db, backupAddressCustoms, traits)
		serviceMember.BackupMailingAddressID = &backupAddress.ID
		serviceMember.BackupMailingAddress = &backupAddress
	}

	if buildType == serviceMemberBuildExtended {
		serviceMember.EmailIsPreferred = models.BoolPointer(true)

		// ensure extended service member has backup mailing address,
		// even if customization is not provided
		if serviceMember.BackupMailingAddressID == nil {
			backupAddress := BuildAddress(db, customs, traits)
			serviceMember.BackupMailingAddressID = &backupAddress.ID
			serviceMember.BackupMailingAddress = &backupAddress
		}
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&serviceMember, cServiceMember)

	if db != nil {
		mustCreate(db, &serviceMember)
	}

	// Extended service members also have a backup contact. Backup
	// contacts need a service member, and so need to wait until after
	// creating the service member before building a backup contact.
	if buildType == serviceMemberBuildExtended {
		backupContactResult := findValidCustomization(customs, BackupContact)

		// before building the backup contact, create a link only
		// customization for the newly created service member if
		// saving to the db
		backupCustoms := []Customization{}
		if db != nil {
			backupCustoms = append(backupCustoms, Customization{
				Model:    serviceMember,
				LinkOnly: true,
			})
		}
		// include any backup contact customizations
		if backupContactResult != nil {
			backupCustoms = append(backupCustoms, *backupContactResult)
		}
		backupContact := BuildBackupContact(db, backupCustoms, traits)
		serviceMember.BackupContacts = append(serviceMember.BackupContacts, backupContact)
		if db != nil {
			mustSave(db, &serviceMember)
		}
	}

	return serviceMember
}

// BuildServiceMember creates a single ServiceMember
// Also creates, if not provided:
// - Residential Address of the ServiceMember
// - User
//
// Will also build User, ResidentialAddress, BackupMailingAddress and
// DutyLocation if and only if a customization is provided
//
// Will never build a BackupContact
func BuildServiceMember(db *pop.Connection, customs []Customization, traits []Trait) models.ServiceMember {
	return buildServiceMemberWithBuildType(db, customs, traits, serviceMemberBuildBasic)
}

// BuildExtendedServiceMember creates a single ServiceMember
// If not provided it will also create an associated
//   - User,
//   - ResidentialAddress
//   - BackupMailingAddress
//   - DutyLocation
//   - BackupContact
func BuildExtendedServiceMember(db *pop.Connection, customs []Customization, traits []Trait) models.ServiceMember {
	return buildServiceMemberWithBuildType(db, customs, traits, serviceMemberBuildExtended)
}

// GetTraitServiceMemberSetIDs is a sample GetTraitFunc
// that sets ids for both ServiceMember and User models
func GetTraitServiceMemberSetIDs() []Customization {
	return []Customization{
		{
			Model: models.ServiceMember{
				ID: uuid.Must(uuid.NewV4()),
			},
		},
		{
			Model: models.User{
				ID: uuid.Must(uuid.NewV4()),
			},
		},
	}
}

// GetTraitServiceMemberUserActive sets the User as Active
func GetTraitActiveServiceMemberUser() []Customization {
	return []Customization{
		{
			Model: models.User{
				Active: true,
			},
		},
	}
}
