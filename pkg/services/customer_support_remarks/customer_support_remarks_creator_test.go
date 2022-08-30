package customersupportremarks

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *CustomerSupportRemarksSuite) TestCustomerSupportRemarksCreator() {
	creator := NewCustomerSupportRemarksCreator()

	suite.Run("Can create customer support remark successfully", func() {

		move := testdatagen.MakeDefaultMove(suite.DB())
		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		remark := &models.CustomerSupportRemark{Content: "Test Remark", OfficeUserID: officeUser.ID}
		createdCustomerSupportRemark, err := creator.CreateCustomerSupportRemark(suite.AppContextForTest(), remark, move.Locator)

		suite.Nil(err)
		suite.NotNil(createdCustomerSupportRemark)
		suite.NotNil(createdCustomerSupportRemark.MoveID)
		suite.Equal(createdCustomerSupportRemark.MoveID, move.ID)
		suite.NotNil(createdCustomerSupportRemark.OfficeUserID)
		suite.Equal(createdCustomerSupportRemark.OfficeUserID, officeUser.ID)
		suite.NotNil(createdCustomerSupportRemark.Content)
		suite.Equal(createdCustomerSupportRemark.Content, remark.Content)
	})

	suite.Run("Remark requires valid move", func() {

		officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
		remark := &models.CustomerSupportRemark{Content: "Bad Move Remark", OfficeUserID: officeUser.ID}
		createdCustomerSupportRemark, err := creator.CreateCustomerSupportRemark(suite.AppContextForTest(), remark, "0")

		suite.Error(err)
		suite.Nil(createdCustomerSupportRemark)
		suite.Equal("FETCH_NOT_FOUND", err.Error())
	})
}
