package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildServiceMember() {
	defaultEmail := "leo_spaceman_sm@example.com"
	defaultAgency := models.AffiliationARMY
	suite.Run("Successful creation of default ServiceMember", func() {
		// Under test:      BuildServiceMember
		// Mocked:          None
		// Set up:          Create a service member with no customizations or traits
		// Expected outcome:serviceMember should be created with default values

		// SETUP
		defaultServiceMember := models.ServiceMember{
			FirstName:     models.StringPointer("Leo"),
			LastName:      models.StringPointer("Spacemen"),
			Telephone:     models.StringPointer("212-123-4567"),
			PersonalEmail: &defaultEmail,
			Affiliation:   &defaultAgency,
		}

		defaultAddress := models.Address{
			StreetAddress1: "123 Any Street",
		}

		defaultUser := models.User{
			OktaEmail: "first.last@okta.mil",

			Active: false,
		}

		// CALL FUNCTION UNDER TEST
		serviceMember := BuildServiceMember(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultServiceMember.FirstName, serviceMember.FirstName)
		suite.Equal(defaultServiceMember.LastName, serviceMember.LastName)
		suite.Equal(defaultServiceMember.PersonalEmail, serviceMember.PersonalEmail)
		suite.Equal(defaultServiceMember.Telephone, serviceMember.Telephone)
		suite.Equal(defaultServiceMember.Affiliation, serviceMember.Affiliation)

		// Check that address was hooked in
		suite.Equal(defaultAddress.StreetAddress1, serviceMember.ResidentialAddress.StreetAddress1)

		// Check that user was hooked in
		suite.Equal(defaultUser.OktaEmail, serviceMember.User.OktaEmail)
		suite.Equal(defaultUser.Active, serviceMember.User.Active)
	})

	suite.Run("Successful creation of customized ServiceMember", func() {
		// Under test:      BuildServiceMember
		// Set up:          Create a Service Member and pass custom fields
		// Expected outcome:serviceMember should be created with custom fields
		// SETUP
		customAffiliation := models.AffiliationAIRFORCE

		customServiceMember := models.ServiceMember{
			FirstName:          models.StringPointer("Gregory"),
			LastName:           models.StringPointer("Van der Heide"),
			Telephone:          models.StringPointer("999-999-9999"),
			SecondaryTelephone: models.StringPointer("123-555-9999"),
			PersonalEmail:      models.StringPointer("peyton@example.com"),
			Edipi:              models.StringPointer("1000011111"),
			Affiliation:        &customAffiliation,
			Suffix:             models.StringPointer("Random suffix string"),
			PhoneIsPreferred:   models.BoolPointer(false),
			EmailIsPreferred:   models.BoolPointer(false),
		}

		customAddress := models.Address{
			StreetAddress1: "987 Another Street",
		}

		customUser := models.User{
			OktaEmail: "test_email@email.com",
		}

		// CALL FUNCTION UNDER TEST
		serviceMember := BuildServiceMember(suite.DB(), []Customization{
			{Model: customServiceMember},
			{Model: customAddress},
			{Model: customUser},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customServiceMember.FirstName, serviceMember.FirstName)
		suite.Equal(customServiceMember.LastName, serviceMember.LastName)
		suite.Equal(customServiceMember.Telephone, serviceMember.Telephone)
		suite.Equal(customServiceMember.SecondaryTelephone, serviceMember.SecondaryTelephone)
		suite.Equal(customServiceMember.PersonalEmail, serviceMember.PersonalEmail)
		suite.Equal(customServiceMember.Edipi, serviceMember.Edipi)
		suite.Equal(customServiceMember.Affiliation, serviceMember.Affiliation)
		suite.Equal(customServiceMember.Suffix, serviceMember.Suffix)
		suite.Equal(customServiceMember.PhoneIsPreferred, serviceMember.PhoneIsPreferred)
		suite.Equal(customServiceMember.EmailIsPreferred, serviceMember.EmailIsPreferred)

		// Check that address was customized
		suite.Equal(customAddress.StreetAddress1, serviceMember.ResidentialAddress.StreetAddress1)

		// Check that user was customized
		suite.Equal(customUser.OktaEmail, serviceMember.User.OktaEmail)
	})

	suite.Run("Successful creation of service member with customized residential and backup mailing address", func() {
		// Under test:      BuildServiceMember
		// Set up:          Create a Service Member with unique residential address and backup mailing address
		// Expected outcome:serviceMember should be created with custom residential address different from address attached backup mailing address

		// SETUP
		customAffiliation := models.AffiliationAIRFORCE

		customResidentialAddress := models.Address{
			StreetAddress1: "123 Any Street",
		}

		customBackupMailingAddress := models.Address{
			StreetAddress1: "456 Something Else Street",
		}

		customServiceMember := models.ServiceMember{
			FirstName:          models.StringPointer("Gregory"),
			LastName:           models.StringPointer("Van der Heide"),
			Telephone:          models.StringPointer("999-999-9999"),
			SecondaryTelephone: models.StringPointer("123-555-9999"),
			PersonalEmail:      models.StringPointer("peyton@example.com"),
			Edipi:              models.StringPointer("1000011111"),
			Affiliation:        &customAffiliation,
			Suffix:             models.StringPointer("Random suffix string"),
			PhoneIsPreferred:   models.BoolPointer(false),
			EmailIsPreferred:   models.BoolPointer(false),
		}

		// CALL FUNCTION UNDER TEST
		serviceMember := BuildServiceMember(suite.DB(), []Customization{
			{Model: customServiceMember},
			{Model: customResidentialAddress, Type: &Addresses.ResidentialAddress},
			{Model: customBackupMailingAddress, Type: &Addresses.BackupMailingAddress},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customServiceMember.FirstName, serviceMember.FirstName)
		suite.Equal(customServiceMember.LastName, serviceMember.LastName)
		suite.Equal(customServiceMember.Telephone, serviceMember.Telephone)
		suite.Equal(customServiceMember.SecondaryTelephone, serviceMember.SecondaryTelephone)
		suite.Equal(customServiceMember.PersonalEmail, serviceMember.PersonalEmail)
		suite.Equal(customServiceMember.Edipi, serviceMember.Edipi)
		suite.Equal(customServiceMember.Affiliation, serviceMember.Affiliation)
		suite.Equal(customServiceMember.Suffix, serviceMember.Suffix)
		suite.Equal(customServiceMember.PhoneIsPreferred, serviceMember.PhoneIsPreferred)
		suite.Equal(customServiceMember.EmailIsPreferred, serviceMember.EmailIsPreferred)

		// Check that Residential Address was customized
		suite.Equal(customResidentialAddress.StreetAddress1, serviceMember.ResidentialAddress.StreetAddress1)

		// Check that Residential Address & Backup Mailing Address are different
		suite.NotEqual(serviceMember.ResidentialAddress.StreetAddress1, serviceMember.BackupMailingAddress.StreetAddress1)

		// no BackupContact for regular ServiceMember
		suite.Equal(0, len(serviceMember.BackupContacts))
	})

	suite.Run("Successful return of linkOnly ServiceMember", func() {
		// Under test:       BuildServiceMember
		// Set up:           Pass in a linkOnly serviceMember
		// Expected outcome: No new ServiceMember should be created.

		// Check num ServiceMember records
		precount, err := suite.DB().Count(&models.ServiceMember{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		serviceMember := BuildServiceMember(suite.DB(), []Customization{
			{
				Model: models.ServiceMember{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.ServiceMember{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, serviceMember.ID)
	})

	suite.Run("Successful return of stubbed ServiceMember", func() {
		// Under test:       BuildServiceMember
		// Set up:           Pass in a linkOnly serviceMember
		// Expected outcome: No new ServiceMember should be created.

		// Check num ServiceMember records
		precount, err := suite.DB().Count(&models.ServiceMember{})
		suite.NoError(err)

		customFirstName := models.StringPointer("Gregory")
		customLastName := models.StringPointer("Van der Heide")
		customTelephone := models.StringPointer("999-999-9999")

		// Nil passed in as db
		serviceMember := BuildServiceMember(nil, []Customization{
			{
				Model: models.ServiceMember{
					FirstName: customFirstName,
					LastName:  customLastName,
					Telephone: customTelephone,
				},
			},
		}, nil)

		count, err := suite.DB().Count(&models.ServiceMember{})
		suite.Equal(precount, count)
		suite.NoError(err)

		suite.Equal(customFirstName, serviceMember.FirstName)
		suite.Equal(customLastName, serviceMember.LastName)
		suite.Equal(customTelephone, serviceMember.Telephone)
	})

	suite.Run("Successful creation of customized ExtendedServiceMember", func() {
		// Under test:      BuildExtendedServiceMember
		//
		// Set up: Create a Service Member with residential address,
		// backup mailing address, dutyLocation & backupContact
		//
		// Expected outcome:serviceMember should be created with
		// custom residential address and backup mailing address.
		// The orders dutyLocation should be used for extended service
		// members and a backupContact should be created

		// SETUP
		customAffiliation := models.AffiliationAIRFORCE

		customServiceMember := models.ServiceMember{
			FirstName:          models.StringPointer("Gregory"),
			LastName:           models.StringPointer("Van der Heide"),
			Telephone:          models.StringPointer("999-999-9999"),
			SecondaryTelephone: models.StringPointer("123-555-9999"),
			PersonalEmail:      models.StringPointer("peyton@example.com"),
			Edipi:              models.StringPointer("1000011111"),
			Affiliation:        &customAffiliation,
			Suffix:             models.StringPointer("Random suffix string"),
		}

		customUser := models.User{
			OktaEmail: "custom.email@example.com",
			Active:    true,
		}
		customResidentialAddress := GetTraitAddress2()[0].Model.(models.Address)
		customBackupAddress := GetTraitAddress3()[0].Model.(models.Address)

		// CALL FUNCTION UNDER TEST
		serviceMember := BuildExtendedServiceMember(suite.DB(), []Customization{
			{
				Model: customServiceMember,
			},
			{
				Model: customUser,
			},
			{
				Model: customResidentialAddress,
				Type:  &Addresses.ResidentialAddress,
			},
			{
				Model: customBackupAddress,
				Type:  &Addresses.BackupMailingAddress,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customServiceMember.FirstName, serviceMember.FirstName)
		suite.Equal(customServiceMember.LastName, serviceMember.LastName)
		suite.Equal(customServiceMember.Telephone, serviceMember.Telephone)
		suite.Equal(customServiceMember.SecondaryTelephone, serviceMember.SecondaryTelephone)
		suite.Equal(customServiceMember.PersonalEmail, serviceMember.PersonalEmail)
		suite.Equal(customServiceMember.Edipi, serviceMember.Edipi)
		suite.Equal(customServiceMember.Affiliation, serviceMember.Affiliation)
		suite.Equal(customServiceMember.Suffix, serviceMember.Suffix)
		// extended service member defaults to email is preferred
		suite.NotNil(serviceMember.EmailIsPreferred)
		suite.True(*serviceMember.EmailIsPreferred)

		// custom user
		suite.Equal(customUser.Active, serviceMember.User.Active)
		suite.Equal(customUser.OktaEmail, serviceMember.User.OktaEmail)

		// custom residential address
		suite.Equal(customResidentialAddress.StreetAddress1,
			serviceMember.ResidentialAddress.StreetAddress1)
		suite.Equal(customResidentialAddress.StreetAddress2,
			serviceMember.ResidentialAddress.StreetAddress2)
		suite.Equal(customResidentialAddress.StreetAddress3,
			serviceMember.ResidentialAddress.StreetAddress3)
		suite.Equal(customResidentialAddress.City,
			serviceMember.ResidentialAddress.City)
		suite.Equal(customResidentialAddress.State,
			serviceMember.ResidentialAddress.State)
		suite.Equal(customResidentialAddress.PostalCode,
			serviceMember.ResidentialAddress.PostalCode)

		// custom backup mailing address
		suite.Equal(customBackupAddress.StreetAddress1,
			serviceMember.BackupMailingAddress.StreetAddress1)
		suite.Equal(customBackupAddress.StreetAddress2,
			serviceMember.BackupMailingAddress.StreetAddress2)
		suite.Equal(customBackupAddress.StreetAddress3,
			serviceMember.BackupMailingAddress.StreetAddress3)
		suite.Equal(customBackupAddress.City,
			serviceMember.BackupMailingAddress.City)
		suite.Equal(customBackupAddress.State,
			serviceMember.BackupMailingAddress.State)
		suite.Equal(customBackupAddress.PostalCode,
			serviceMember.BackupMailingAddress.PostalCode)

		// Check that backup contact was made and appended to service member
		suite.Equal(1, len(serviceMember.BackupContacts))
		suite.Equal(models.BackupContactPermissionEDIT, serviceMember.BackupContacts[0].Permission)
	})
	suite.Run("Successful return of stubbed ExtendedServiceMember", func() {
		// Under test:       BuildServiceMember
		// Set up:           Pass in a linkOnly serviceMember
		// Expected outcome: No new ServiceMember should be created.

		// Check num ServiceMember records
		precount, err := suite.DB().Count(&models.ServiceMember{})
		suite.NoError(err)

		BuildExtendedServiceMember(nil, nil, nil)

		count, err := suite.DB().Count(&models.ServiceMember{})
		suite.Equal(precount, count)
		suite.NoError(err)
	})
}
