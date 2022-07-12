package customersupportremarks

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *CustomerSupportRemarksSuite) TestCustomerSupportRemarkDeleter() {
	setupTestData := func() (services.CustomerSupportRemarkDeleter, models.CustomerSupportRemark, appcontext.AppContext) {
		deleter := NewCustomerSupportRemarkDeleter()
		remark := testdatagen.MakeDefaultCustomerSupportRemark(suite.DB())

		appCtx := suite.AppContextWithSessionForTest(&auth.Session{
			OfficeUserID: remark.OfficeUserID,
		})

		return deleter, remark, appCtx
	}

	suite.Run("delete existing remark", func() {
		deleter, remark, appCtx := setupTestData()

		suite.NoError(deleter.DeleteCustomerSupportRemark(appCtx, remark.ID))
		var dbRemark models.CustomerSupportRemark
		err := suite.DB().Find(&dbRemark, remark.ID)
		suite.NoError(err)
		suite.NotNil(dbRemark.DeletedAt)
	})

	suite.Run("Returns an error when delete non-existent remark", func() {
		deleter, _, appCtx := setupTestData()

		uuid := uuid.Must(uuid.NewV4())
		err := deleter.DeleteCustomerSupportRemark(appCtx, uuid)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)

	})

	suite.Run("Returns an error when attempting to delete an already deleted remark", func() {
		deleter, remark, appCtx := setupTestData()

		suite.NoError(deleter.DeleteCustomerSupportRemark(appCtx, remark.ID))
		var dbRemark models.CustomerSupportRemark
		err := suite.DB().Find(&dbRemark, remark.ID)
		suite.NoError(err)
		suite.NotNil(dbRemark.DeletedAt)

		err = deleter.DeleteCustomerSupportRemark(appCtx, remark.ID)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
