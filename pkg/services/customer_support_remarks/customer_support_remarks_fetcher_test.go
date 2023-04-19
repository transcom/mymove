package customersupportremarks

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *CustomerSupportRemarksSuite) setupTestData() (models.CustomerSupportRemarks, models.Move) {

	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	move := factory.BuildMove(suite.DB(), nil, nil)

	var customerSupportRemarks models.CustomerSupportRemarks
	for i := 0; i < 5; i++ {
		remark := factory.BuildCustomerSupportRemark(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model: models.CustomerSupportRemark{
					Content: fmt.Sprintln("This is remark number: %i", i),
				},
			},
		}, nil)
		customerSupportRemarks = append(customerSupportRemarks, remark)
	}
	return customerSupportRemarks, move
}

func (suite *CustomerSupportRemarksSuite) setupTestDataMultipleUsers() (models.CustomerSupportRemarks, models.Move) {
	move := factory.BuildMove(suite.DB(), nil, nil)

	var officeUsers models.OfficeUsers
	var customerSupportRemarks models.CustomerSupportRemarks
	for i := 0; i < 3; i++ {
		officeUsers = append(officeUsers, factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO}))
		for x := 0; x < 2; x++ {
			remark := factory.BuildCustomerSupportRemark(suite.DB(), []factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    officeUsers[i],
					LinkOnly: true,
				},
				{
					Model: models.CustomerSupportRemark{
						Content: fmt.Sprintln("This is remark number: %i", i),
					},
				},
			}, nil)
			customerSupportRemarks = append(customerSupportRemarks, remark)
		}
	}
	return customerSupportRemarks, move
}

func (suite *CustomerSupportRemarksSuite) TestCustomerSupportRemarksListFetcher() {
	fetcher := NewCustomerSupportRemarks()

	suite.Run("Can fetch office move remarks successfully", func() {
		createdCustomerSupportRemarks, move := suite.setupTestData()
		customerSupportRemarks, err := fetcher.ListCustomerSupportRemarks(suite.AppContextForTest(), move.Locator)
		suite.NoError(err)
		suite.NotNil(customerSupportRemarks)

		customerSupportRemarkValues := *customerSupportRemarks
		suite.Len(customerSupportRemarkValues, len(createdCustomerSupportRemarks))
	})

	suite.Run("Can fetch office move remarks involving multiple users properly", func() {
		createdCustomerSupportRemarks, move := suite.setupTestDataMultipleUsers()
		customerSupportRemarks, err := fetcher.ListCustomerSupportRemarks(suite.AppContextForTest(), move.Locator)
		suite.NoError(err)
		suite.NotNil(customerSupportRemarks)

		customerSupportRemarkValues := *customerSupportRemarks
		suite.Len(createdCustomerSupportRemarks, len(customerSupportRemarkValues))
	})

	suite.Run("Office move remarks aren't found", func() {
		_, _ = suite.setupTestData()
		incorrectLocator := "ZZZZZZZZ"
		_, err := fetcher.ListCustomerSupportRemarks(suite.AppContextForTest(), incorrectLocator)
		suite.Error(models.ErrFetchNotFound, err)
	})

	suite.Run("Soft deleted remarks should not be returned", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		move := factory.BuildMove(suite.DB(), nil, nil)
		remark := factory.BuildCustomerSupportRemark(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model: models.CustomerSupportRemark{
					Content: "this is a remark",
				},
			},
		}, nil)

		deletedTime := time.Now()
		factory.BuildCustomerSupportRemark(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
			{
				Model: models.CustomerSupportRemark{
					Content:   "this is a deleted remark",
					DeletedAt: &deletedTime,
				},
			},
		}, nil)
		customerSupportRemarks, err := fetcher.ListCustomerSupportRemarks(suite.AppContextForTest(), move.Locator)
		suite.NoError(err)
		suite.NotNil(customerSupportRemarks)

		customerSupportRemarkValues := *customerSupportRemarks
		suite.Len(customerSupportRemarkValues, 1)

		suite.Equal(remark.Content, customerSupportRemarkValues[0].Content)
	})
}
