package models_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestMovingExpenseValidation() {
	blankExpenseType := models.MovingExpenseReceiptType("")
	blankStatusType := models.PPMDocumentStatus("")

	validExpenseTypes := strings.Join(models.AllowedExpenseTypes, ", ")

	validStatuses := strings.Join(models.AllowedPPMDocumentStatuses, ", ")

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
		"Missing UUID": {
			movingExpense: models.MovingExpense{},
			expectedErrs: map[string][]string{
				"ppmshipment_id": {"PPMShipmentID can not be blank."},
			},
		},
		"Optional fields are valid": {
			movingExpense: models.MovingExpense{
				PPMShipmentID:      uuid.Must(uuid.NewV4()),
				DocumentID:         uuid.Must(uuid.NewV4()),
				DeletedAt:          models.TimePointer(time.Time{}),
				MovingExpenseType:  &blankExpenseType,
				Description:        models.StringPointer(""),
				Status:             &blankStatusType,
				Reason:             models.StringPointer(""),
				SITStartDate:       models.TimePointer(time.Time{}),
				SITEndDate:         models.TimePointer(time.Time{}),
				TrackingNumber:     models.StringPointer(""),
				WeightShipped:      models.PoundPointer(unit.Pound(-1)),
				ProGearDescription: models.StringPointer(""),
			},
			expectedErrs: map[string][]string{
				"deleted_at":           {"DeletedAt can not be blank."},
				"moving_expense_type":  {fmt.Sprintf("MovingExpenseType is not in the list [%s].", validExpenseTypes)},
				"description":          {"Description can not be blank."},
				"status":               {fmt.Sprintf("Status is not in the list [%s].", validStatuses)},
				"reason":               {"Reason can not be blank."},
				"sitstart_date":        {"SITStartDate can not be blank."},
				"sitend_date":          {"SITEndDate can not be blank."},
				"tracking_number":      {"TrackingNumber can not be blank."},
				"weight_shipped":       {"-1 is less than zero."},
				"pro_gear_description": {"ProGearDescription can not be blank."},
			},
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		suite.Run(name, func() {
			suite.verifyValidationErrors(&tc.movingExpense, tc.expectedErrs, nil)
		})
	}

	suite.Run("Can create a moving expense", func() {
		// This test is meant to be a base smoke test to make sure this model/table works
		ppmShipment := factory.BuildPPMShipment(suite.DB(), nil, []factory.Trait{factory.GetTraitApprovedPPMWithActualInfo})
		serviceMember := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember

		document := factory.BuildDocumentLinkServiceMember(suite.DB(), serviceMember)

		movingExpense := models.MovingExpense{
			PPMShipmentID: ppmShipment.ID,
			DocumentID:    document.ID,
		}

		suite.MustCreate(&movingExpense)
	})
}
func (suite *ModelSuite) TestMovingExpenses_FilterRejected() {
	suite.Run("returns empty slice when input is empty", func() {
		expenses := models.MovingExpenses{}
		filtered := expenses.FilterRejected()
		suite.Equal(0, len(filtered))
	})

	suite.Run("filters out rejected expenses", func() {
		rejectedStatus := models.PPMDocumentStatusRejected
		expenses := models.MovingExpenses{
			{Status: &rejectedStatus},
			{Status: nil},
		}

		filtered := expenses.FilterRejected()
		suite.Equal(1, len(filtered))

		for _, expense := range filtered {
			if expense.Status != nil {
				suite.NotEqual(models.PPMDocumentStatusRejected, *expense.Status)
			}
		}
	})

	suite.Run("keeps all expenses when none are rejected", func() {
		approvedStatus := models.PPMDocumentStatusApproved
		expenses := models.MovingExpenses{
			{Status: &approvedStatus},
			{Status: nil},
		}

		filtered := expenses.FilterRejected()
		suite.Equal(2, len(filtered))
		suite.Equal(expenses, filtered)
	})

	suite.Run("handles all rejected expenses", func() {
		rejectedStatus := models.PPMDocumentStatusRejected
		expenses := models.MovingExpenses{
			{Status: &rejectedStatus},
			{Status: &rejectedStatus},
		}

		filtered := expenses.FilterRejected()
		suite.Equal(0, len(filtered))
	})
}
