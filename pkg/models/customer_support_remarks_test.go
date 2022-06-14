package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestCustomerSupportRemarkCreation() {
	move := testdatagen.MakeDefaultMove(suite.DB())
	suite.NotNil(move)

	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	suite.NotNil(officeUser)

	suite.Run("test valid office remark", func() {
		customerSupportRemark := "This is a note that's saying something about the move."
		validCustomerSupportRemark := models.CustomerSupportRemark{
			Content:      customerSupportRemark,
			OfficeUser:   officeUser,
			OfficeUserID: officeUser.ID,
			Move:         move,
			MoveID:       move.ID,
		}

		suite.MustSave(&validCustomerSupportRemark)
		suite.NotNil(validCustomerSupportRemark.ID)
		suite.NotEqual(uuid.Nil, validCustomerSupportRemark.ID)
		suite.Equal(move.ID, validCustomerSupportRemark.MoveID)
		suite.Equal(customerSupportRemark, validCustomerSupportRemark.Content)
		suite.Equal(officeUser.ID, validCustomerSupportRemark.OfficeUserID)
	})
}
