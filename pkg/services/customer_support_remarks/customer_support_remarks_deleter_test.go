package customersupportremarks

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *CustomerSupportRemarksSuite) TestCustomerSupportRemarkDeleter() {
	deleter := NewCustomerSupportRemarkDeleter()

	suite.Run("delete existing remark", func() {
		remark := testdatagen.MakeDefaultCustomerSupportRemark(suite.DB())
		suite.NoError(deleter.DeleteCustomerSupportRemark(suite.AppContextForTest(), remark.ID))
		var dbRemark models.CustomerSupportRemark
		err := suite.DB().Find(&dbRemark, remark.ID)
		suite.NoError(err)
		suite.NotNil(dbRemark.DeletedAt)
	})

	suite.Run("Returns an error when delete non-existent remark", func() {
		uuid := uuid.Must(uuid.NewV4())
		err := deleter.DeleteCustomerSupportRemark(suite.AppContextForTest(), uuid)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)

	})

	suite.Run("Returns an error when attempting to delete an already deleted remark", func() {
		remark := testdatagen.MakeDefaultCustomerSupportRemark(suite.DB())
		suite.NoError(deleter.DeleteCustomerSupportRemark(suite.AppContextForTest(), remark.ID))
		var dbRemark models.CustomerSupportRemark
		err := suite.DB().Find(&dbRemark, remark.ID)
		suite.NoError(err)
		suite.NotNil(dbRemark.DeletedAt)
		err = deleter.DeleteCustomerSupportRemark(suite.AppContextForTest(), remark.ID)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
