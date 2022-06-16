package customersupportremarks

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *CustomerSupportRemarksSuite) TestCustomerSupportRemarksUpdater() {
	updater := NewCustomerSupportRemarkUpdater()

	suite.Run("Can update customer support remark successfully", func() {
		origRemark := testdatagen.MakeDefaultCustomerSupportRemark(suite.DB())

		remarkEdit := "Test Remark"
		payload := ghcmessages.UpdateCustomerSupportRemarkPayload{Content: &remarkEdit, ID: handlers.FmtUUID(origRemark.ID)}
		updatedCustomerSupportRemark, err := updater.UpdateCustomerSupportRemark(suite.AppContextForTest(), payload)

		suite.Nil(err)
		suite.NotNil(updatedCustomerSupportRemark)
		suite.NotNil(updatedCustomerSupportRemark.MoveID)
		suite.Equal(updatedCustomerSupportRemark.MoveID, origRemark.MoveID)
		suite.NotNil(updatedCustomerSupportRemark.OfficeUserID)
		suite.Equal(updatedCustomerSupportRemark.OfficeUserID, origRemark.OfficeUserID)
		suite.NotNil(updatedCustomerSupportRemark.Content)
		suite.Equal(updatedCustomerSupportRemark.Content, remarkEdit)
		suite.NotEqual(updatedCustomerSupportRemark.Content, origRemark.Content)
	})

	suite.Run("Returns an error when remark is not found", func() {
		badID := uuid.Must(uuid.NewV4())
		remarkEdit := "Test Remark"
		payload := ghcmessages.UpdateCustomerSupportRemarkPayload{Content: &remarkEdit, ID: handlers.FmtUUID(badID)}
		_, err := updater.UpdateCustomerSupportRemark(suite.AppContextForTest(), payload)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

}
