package customersupportremarks

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *CustomerSupportRemarksSuite) setupTestData() models.CustomerSupportRemarks {

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	move := testdatagen.MakeDefaultMove(suite.DB())

	var customerSupportRemarks models.CustomerSupportRemarks
	for i := 0; i < 5; i++ {
		remark := testdatagen.MakeOfficeMoveRemark(suite.DB(),
			testdatagen.Assertions{
				OfficeMoveRemark: models.CustomerSupportRemark{
					Content:      fmt.Sprintln("This is remark number: %i", i),
					OfficeUserID: officeUser.ID,
					MoveID:       move.ID,
				}})
		customerSupportRemarks = append(customerSupportRemarks, remark)
	}
	return customerSupportRemarks
}

func (suite *CustomerSupportRemarksSuite) setupTestDataMultipleUsers() models.CustomerSupportRemarks {
	move := testdatagen.MakeDefaultMove(suite.DB())

	var officeUsers models.OfficeUsers
	var customerSupportRemarks models.CustomerSupportRemarks
	for i := 0; i < 3; i++ {
		officeUsers = append(officeUsers, testdatagen.MakeDefaultOfficeUser(suite.DB()))
		for x := 0; x < 2; x++ {
			remark := testdatagen.MakeOfficeMoveRemark(suite.DB(),
				testdatagen.Assertions{
					OfficeMoveRemark: models.CustomerSupportRemark{
						Content:      fmt.Sprintln("This is remark number: %i", i),
						OfficeUserID: officeUsers[i].ID,
						MoveID:       move.ID,
					}})
			customerSupportRemarks = append(customerSupportRemarks, remark)
		}
	}
	return customerSupportRemarks
}

func (suite *CustomerSupportRemarksSuite) TestCustomerSupportRemarksListFetcher() {
	fetcher := NewCustomerSupportRemarks()

	suite.Run("Can fetch office move remarks successfully", func() {
		createdCustomerSupportRemarks := suite.setupTestData()
		customerSupportRemarks, err := fetcher.ListCustomerSupportRemarks(suite.AppContextForTest(), createdCustomerSupportRemarks[0].MoveID)
		suite.NoError(err)
		suite.NotNil(customerSupportRemarks)

		customerSupportRemarkValues := *customerSupportRemarks
		suite.Len(customerSupportRemarkValues, 5)
	})

	suite.Run("Can fetch office move remarks involving multiple users properly", func() {
		createdCustomerSupportRemarks := suite.setupTestDataMultipleUsers()
		customerSupportRemarks, err := fetcher.ListCustomerSupportRemarks(suite.AppContextForTest(), createdCustomerSupportRemarks[0].MoveID)
		suite.NoError(err)
		suite.NotNil(customerSupportRemarks)

		customerSupportRemarkValues := *customerSupportRemarks
		suite.Len(createdCustomerSupportRemarks, len(customerSupportRemarkValues))
	})

	suite.Run("Office move remarks aren't found", func() {
		_ = suite.setupTestData()
		randomUUID, _ := uuid.NewV4()
		_, err := fetcher.ListCustomerSupportRemarks(suite.AppContextForTest(), randomUUID)
		suite.Error(models.ErrFetchNotFound, err)
	})

}
