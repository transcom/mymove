package models_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestMovingExpenseValidation() {
	blankExpenseType := models.MovingExpenseReceiptType("")
	blankStatusType := models.PPMDocumentStatus("")
	validExpenseTypes := strings.Join([]string{
		string(models.MovingExpenseReceiptTypeContractedExpense),
		string(models.MovingExpenseReceiptTypeOil),
		string(models.MovingExpenseReceiptTypePackingMaterials),
		string(models.MovingExpenseReceiptTypeRentalEquipment),
		string(models.MovingExpenseReceiptTypeStorage),
		string(models.MovingExpenseReceiptTypeTolls),
		string(models.MovingExpenseReceiptTypeWeighingFees),
		string(models.MovingExpenseReceiptTypeOther),
	}, ", ")

	validStatuses := strings.Join([]string{
		string(models.PPMDocumentStatusApproved),
		string(models.PPMDocumentStatusExcluded),
		string(models.PPMDocumentStatusRejected),
	}, ", ")

	testCases := map[string]struct {
		movingExpense models.MovingExpense
		expectedErrs  map[string][]string
	}{
		"Successful create": {
			movingExpense: models.MovingExpense{
				PPMShipmentID: uuid.Must(uuid.NewV4()),
				DocumentID:    uuid.Must(uuid.NewV4()),
			},
			expectedErrs: nil,
		},
		"Missing UUIDs": {
			movingExpense: models.MovingExpense{},
			expectedErrs: map[string][]string{
				"ppmshipment_id": {"PPMShipmentID can not be blank."},
				"document_id":    {"DocumentID can not be blank."},
			},
		},
		"Optional fields are valid": {
			movingExpense: models.MovingExpense{
				PPMShipmentID:     uuid.Must(uuid.NewV4()),
				DocumentID:        uuid.Must(uuid.NewV4()),
				DeletedAt:         models.TimePointer(time.Time{}),
				MovingExpenseType: &blankExpenseType,
				Description:       models.StringPointer(""),
				Status:            &blankStatusType,
				SITStartDate:      models.TimePointer(time.Time{}),
				SITEndDate:        models.TimePointer(time.Time{}),
			},
			expectedErrs: map[string][]string{
				"deleted_at":          {"DeletedAt can not be blank."},
				"moving_expense_type": {fmt.Sprintf("MovingExpenseType is not in the list [%s].", validExpenseTypes)},
				"description":         {"Description can not be blank."},
				"status":              {fmt.Sprintf("Status is not in the list [%s].", validStatuses)},
				"sitstart_date":       {"SITStartDate can not be blank."},
				"sitend_date":         {"SITEndDate can not be blank."},
			},
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		suite.Run(name, func() {
			suite.verifyValidationErrors(&tc.movingExpense, tc.expectedErrs)
		})
	}

	suite.Run("Can create a moving expense", func() {
		// This test is meant to be a base smoke test to make sure this model/table works
		ppmShipment := testdatagen.MakeApprovedPPMShipmentWithActualInfo(suite.DB(), testdatagen.Assertions{})

		serviceMember := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		document := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{
			Document: models.Document{
				ServiceMemberID: serviceMember.ID,
				ServiceMember:   serviceMember,
			},
		})

		movingExpense := models.MovingExpense{
			PPMShipmentID: ppmShipment.ID,
			DocumentID:    document.ID,
		}

		suite.MustCreate(&movingExpense)
	})
}
