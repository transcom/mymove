package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildCustomerSupportRemark() {
	defaultContent := "This is an office remark."
	suite.Run("Successful creation of default customerSupportRemark", func() {
		// Under test:      BuildCustomerSupportRemark
		// Mocked:          None
		// Set up:          Create an CustomerSupportRemark with no customizations or traits
		// Expected outcome:CustomerSupportRemark should be created with default values

		customerSupportRemark := BuildCustomerSupportRemark(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultContent, customerSupportRemark.Content)
		suite.False(customerSupportRemark.MoveID.IsNil())
		suite.False(customerSupportRemark.OfficeUserID.IsNil())
	})

	suite.Run("Successful creation of an customerSupportRemark with customization", func() {
		// Under test:      BuildCustomerSupportRemark
		// Set up:          Create an CustomerSupportRemark with
		// custom content
		// Expected outcome:CustomerSupportRemark should be created
		// with custom content

		customRemark := models.CustomerSupportRemark{
			Content: "my content",
		}
		customMove := models.Move{
			Locator: "ABC999",
		}
		customOfficeUser := models.OfficeUser{
			FirstName: "Csr",
			LastName:  "Test",
		}
		customerSupportRemark := BuildCustomerSupportRemark(suite.DB(), []Customization{
			{
				Model: customRemark,
			},
			{
				Model: customMove,
			},
			{
				Model: customOfficeUser,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customRemark.Content, customerSupportRemark.Content)
		suite.False(customerSupportRemark.MoveID.IsNil())
		suite.False(customerSupportRemark.OfficeUserID.IsNil())
		suite.Equal(customMove.Locator, customerSupportRemark.Move.Locator)
		suite.Equal(customOfficeUser.FirstName, customerSupportRemark.OfficeUser.FirstName)
		suite.Equal(customOfficeUser.LastName, customerSupportRemark.OfficeUser.LastName)
	})

	suite.Run("Successful creation of stubbed customerSupportRemark", func() {
		// Under test:      BuildCustomerSupportRemark
		// Set up:          Create a customized customerSupportRemark, but don't pass in a db
		// Expected outcome:CustomerSupportRemark should be created with email
		//                  No customerSupportRemark should be created in database
		precount, err := suite.DB().Count(&models.CustomerSupportRemark{})
		suite.NoError(err)

		customRemark := models.CustomerSupportRemark{
			Content: "my content",
		}
		customerSupportRemark := BuildCustomerSupportRemark(nil, []Customization{
			{
				Model: customRemark,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customRemark.Content, customerSupportRemark.Content)
		suite.True(customerSupportRemark.MoveID.IsNil())
		suite.True(customerSupportRemark.OfficeUserID.IsNil())

		// Count how many customerSupportRemarks are in the DB, no new customerSupportRemarks should have been created
		count, err := suite.DB().Count(&models.CustomerSupportRemark{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("Successful creation of customerSupportRemark with linked customerSupportRemark", func() {
		// Under test:       BuildCustomerSupportRemark
		// Set up:           Create an customerSupportRemark and pass in a linkOnly customerSupportRemark
		// Expected outcome: No new customerSupportRemark should be created.

		// Check num customerSupportRemarks
		precount, err := suite.DB().Count(&models.CustomerSupportRemark{})
		suite.NoError(err)

		customRemark := models.CustomerSupportRemark{
			ID:      uuid.Must(uuid.NewV4()),
			Content: "nodb content",
		}
		customerSupportRemark := BuildCustomerSupportRemark(suite.DB(), []Customization{
			{
				Model:    customRemark,
				LinkOnly: true,
			},
		}, nil)

		count, err := suite.DB().Count(&models.CustomerSupportRemark{})
		suite.Equal(precount, count)
		suite.NoError(err)

		// VALIDATE RESULTS
		suite.Equal(customRemark.ID, customerSupportRemark.ID)
		suite.Equal(customRemark.Content, customerSupportRemark.Content)
	})
}
