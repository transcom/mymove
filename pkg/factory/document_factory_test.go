package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildDocument() {
	suite.Run("Successful creation of default Document", func() {
		// Under test:      BuildDocument
		// Mocked:          None
		// Set up:          Create a Document with no customizations or traits
		// Expected outcome:Document should be created with default values

		// SETUP
		defaultServiceMember := BuildServiceMember(nil, nil, nil)
		defaultDocument := models.Document{
			ServiceMember: defaultServiceMember,
		}

		// CALL FUNCTION UNDER TEST
		document := BuildDocument(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.NotNil(document.ServiceMember)
		suite.NotNil(document.CreatedAt)

		// Check that service member was hooked in
		suite.Equal(*defaultDocument.ServiceMember.FirstName, *document.ServiceMember.FirstName)
		suite.Equal(*defaultDocument.ServiceMember.LastName, *document.ServiceMember.LastName)
		suite.Equal(*defaultDocument.ServiceMember.Telephone, *document.ServiceMember.Telephone)

	})

	suite.Run("Successful creation of customized Document", func() {
		// Under test:      BuildDocument
		// Set up:          Create a Document and pass custom fields
		// Expected outcome:Document should be created with custom fields

		// SETUP
		customDocument := models.Document{
			ID: uuid.Must(uuid.NewV4()),
		}

		customServiceMember := models.ServiceMember{
			FirstName: models.StringPointer("Jason"),
			LastName:  models.StringPointer("Ash"),
		}

		// CALL FUNCTION UNDER TEST
		document := BuildDocument(suite.DB(), []Customization{
			{Model: customDocument},
			{Model: customServiceMember},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customDocument.ID, document.ID)

		// Check that the service member was customized
		suite.Equal(*customServiceMember.FirstName, *document.ServiceMember.FirstName)
		suite.Equal(*customServiceMember.LastName, *document.ServiceMember.LastName)
	})

	suite.Run("Successful return of linkOnly Document", func() {
		// Under test:       BuildDocument
		// Set up:           Pass in a linkOnly Document
		// Expected outcome: No new Document should be created.

		// Check num Document records
		precount, err := suite.DB().Count(&models.Document{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		document := BuildDocument(suite.DB(), []Customization{
			{
				Model: models.Document{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)
		count, err := suite.DB().Count(&models.Document{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, document.ID)

	})
	suite.Run("Successful return of stubbed Document", func() {
		// Under test:       BuildDocument
		// Set up:           Create a Document with nil DB
		// Expected outcome: No new Document should be created.

		// Check num Document records
		precount, err := suite.DB().Count(&models.Document{})
		suite.NoError(err)

		// Nil passed in as db
		id := uuid.Must(uuid.NewV4())
		document := BuildDocument(nil, []Customization{
			{
				Model: models.Document{
					ID: id,
				},
			},
		}, nil)
		count, err := suite.DB().Count(&models.Document{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(id, document.ID)
	})

	suite.Run("Successful creation of Document using BuildDocumentLinkServiceMember", func() {
		// Under test:       BuildDocumentLinkServiceMember
		// Set up:           Create a Document
		// Expected outcome: No new Document should be created.

		serviceMember := BuildServiceMember(suite.DB(), []Customization{
			{
				Model: models.ServiceMember{
					FirstName: models.StringPointer("Jason"),
					LastName:  models.StringPointer("Ash"),
				},
			},
		}, nil)

		// CALL FUNCTION UNDER TEST
		document := BuildDocumentLinkServiceMember(suite.DB(), serviceMember)

		// Check that the service member was hooked in
		suite.Equal(*serviceMember.FirstName, *document.ServiceMember.FirstName)
		suite.Equal(*serviceMember.LastName, *document.ServiceMember.LastName)
	})
}
