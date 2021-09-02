package mtoshipment

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestApproveSITExtension() {
	suite.T().Run("Returns an error when shipment is not found", func(t *testing.T) {
		sitExtensionApprover := NewSITExtensionApprover()
		nonexistentUUID := uuid.Must(uuid.NewV4())
		approvedDays := int(20)
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionApprover.ApproveSITExtension(suite.TestAppContext(), nonexistentUUID, nonexistentUUID, &approvedDays, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when SIT extension is not found", func(t *testing.T) {
		sitExtensionApprover := NewSITExtensionApprover()
		nonexistentUUID := uuid.Must(uuid.NewV4())
		mtoShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		approvedDays := int(20)
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionApprover.ApproveSITExtension(suite.TestAppContext(), mtoShipment.ID, nonexistentUUID, &approvedDays, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when etag does not match", func(t *testing.T) {
		sitExtensionApprover := NewSITExtensionApprover()
		mtoShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		sitExtension := testdatagen.MakePendingSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		approvedDays := int(20)
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionApprover.ApproveSITExtension(suite.TestAppContext(), mtoShipment.ID, sitExtension.ID, &approvedDays, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), mtoShipment.ID.String())
	})
}
