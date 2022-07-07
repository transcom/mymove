package customersupportremarks

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *CustomerSupportRemarksSuite) setupTestData() (models.CustomerSupportRemarks, models.Move) {

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	move := testdatagen.MakeDefaultMove(suite.DB())

	var customerSupportRemarks models.CustomerSupportRemarks
	for i := 0; i < 5; i++ {
		remark := testdatagen.MakeCustomerSupportRemark(suite.DB(),
			testdatagen.Assertions{
				CustomerSupportRemark: models.CustomerSupportRemark{
					Content:      fmt.Sprintln("This is remark number: %i", i),
					OfficeUserID: officeUser.ID,
					MoveID:       move.ID,
				}})
		customerSupportRemarks = append(customerSupportRemarks, remark)
	}
	return customerSupportRemarks, move
}

func (suite *CustomerSupportRemarksSuite) setupTestDataMultipleUsers() (models.CustomerSupportRemarks, models.Move) {
	move := testdatagen.MakeDefaultMove(suite.DB())

	var officeUsers models.OfficeUsers
	var customerSupportRemarks models.CustomerSupportRemarks
	for i := 0; i < 3; i++ {
		officeUsers = append(officeUsers, testdatagen.MakeDefaultOfficeUser(suite.DB()))
		for x := 0; x < 2; x++ {
			remark := testdatagen.MakeCustomerSupportRemark(suite.DB(),
				testdatagen.Assertions{
					CustomerSupportRemark: models.CustomerSupportRemark{
						Content:      fmt.Sprintln("This is remark number: %i", i),
						OfficeUserID: officeUsers[i].ID,
						MoveID:       move.ID,
					}})
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
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		move := testdatagen.MakeDefaultMove(suite.DB())
		remark := testdatagen.MakeCustomerSupportRemark(suite.DB(),
			testdatagen.Assertions{
				CustomerSupportRemark: models.CustomerSupportRemark{
					Content:      "this is a remark",
					OfficeUserID: officeUser.ID,
					MoveID:       move.ID,
				}})

		deletedTime := time.Now()
		testdatagen.MakeCustomerSupportRemark(suite.DB(),
			testdatagen.Assertions{
				CustomerSupportRemark: models.CustomerSupportRemark{
					Content:      "this is a deleted remark",
					OfficeUserID: officeUser.ID,
					MoveID:       move.ID,
					DeletedAt:    &deletedTime,
				}})
		customerSupportRemarks, err := fetcher.ListCustomerSupportRemarks(suite.AppContextForTest(), move.Locator)
		suite.NoError(err)
		suite.NotNil(customerSupportRemarks)

		customerSupportRemarkValues := *customerSupportRemarks
		suite.Len(customerSupportRemarkValues, 1)

		suite.Equal(remark.Content, customerSupportRemarkValues[0].Content)
	})
}
