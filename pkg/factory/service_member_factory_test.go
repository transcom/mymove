package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildServiceMember() {
	defaultEmail := "leo_spaceman_sm@example.com"
	defaultAgency := models.AffiliationARMY
	defaultRank := models.ServiceMemberRankE1
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
			Rank:          &defaultRank,
			PersonalEmail: &defaultEmail,
			Affiliation:   &defaultAgency,
		}

		defaultAddress := models.Address{
			StreetAddress1: "123 Any Street",
		}

		defaultUser := models.User{
			LoginGovEmail: "first.last@login.gov.test",
			Active:        false,
		}

		// CALL FUNCTION UNDER TEST
		serviceMember := BuildServiceMember(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultServiceMember.FirstName, serviceMember.FirstName)
		suite.Equal(defaultServiceMember.LastName, serviceMember.LastName)
		suite.Equal(defaultServiceMember.Rank, serviceMember.Rank)
		suite.Equal(defaultServiceMember.PersonalEmail, serviceMember.PersonalEmail)
		suite.Equal(defaultServiceMember.Telephone, serviceMember.Telephone)
		suite.Equal(defaultServiceMember.Affiliation, serviceMember.Affiliation)

		// Check that address was hooked in
		suite.Equal(defaultAddress.StreetAddress1, serviceMember.ResidentialAddress.StreetAddress1)

		// Check that user was hooked in
		suite.Equal(defaultUser.LoginGovEmail, serviceMember.User.LoginGovEmail)
		suite.Equal(defaultUser.Active, serviceMember.User.Active)
	})

	suite.Run("Successful creation of customized ServiceMember", func() {
		// Under test:      BuildServiceMember
		// Set up:          Create a Service Member and pass custom fields
		// Expected outcome:serviceMember should be created with custom fields
		// SETUP
		customRank := models.ServiceMemberRankE3
		customAffiliation := models.AffiliationAIRFORCE

		customServiceMember := models.ServiceMember{
			FirstName:          models.StringPointer("Gregory"),
			LastName:           models.StringPointer("Van der Heide"),
			Telephone:          models.StringPointer("999-999-9999"),
			SecondaryTelephone: models.StringPointer("123-555-9999"),
			PersonalEmail:      models.StringPointer("peyton@example.com"),
			Rank:               &customRank,
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
			LoginGovEmail: "test_email@email.com",
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
		suite.Equal(customServiceMember.Rank, serviceMember.Rank)
		suite.Equal(customServiceMember.Edipi, serviceMember.Edipi)
		suite.Equal(customServiceMember.Affiliation, serviceMember.Affiliation)
		suite.Equal(customServiceMember.Suffix, serviceMember.Suffix)
		suite.Equal(customServiceMember.PhoneIsPreferred, serviceMember.PhoneIsPreferred)
		suite.Equal(customServiceMember.EmailIsPreferred, serviceMember.EmailIsPreferred)

		// Check that address was customized
		suite.Equal(customAddress.StreetAddress1, serviceMember.ResidentialAddress.StreetAddress1)

		// Check that user was customized
		suite.Equal(customUser.LoginGovEmail, serviceMember.User.LoginGovEmail)
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
}
