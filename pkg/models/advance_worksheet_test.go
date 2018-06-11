package models_test

import (
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// Make a test that checks a legit advance worksheet struct
func (suite *ModelSuite) Test_MakeAdvanceWorksheetInfo() {
	t := suite.T()

	// make service member
	sm, err := testdatagen.MakeExtendedServiceMember(suite.db)
	if err != nil {
		t.Errorf("Error making a service member")
	}

	// make orders
	ord, err := testdatagen.MakeOrder(suite.db)
	if err != nil {
		t.Errorf("Error making an order")
	}

	// make ppm
	ppm, err := testdatagen.MakePPM(suite.db)
	if err != nil {
		t.Errorf("Error making a PPM")
	}

	// make reimbursement
	reimbursement, err := testdatagen.MakeRequestedReimbursement(suite.db)
	if err != nil {
		t.Errorf("Error making a reimbursement")
	}

	// make backup contact
	bu, err := testdatagen.MakeBackupContact(suite.db, &sm.ID)
	if err != nil {
		t.Errorf("Error making a backup contact")
	}

	advanceWorksheet := models.AdvanceWorksheet{
		FirstName:                            *sm.FirstName,
		LastName:                             *sm.LastName,
		Email:                                sm.PersonalEmail,
		OrderIssueDate:                       ord.IssueDate,
		OrdersType:                           ord.OrdersType,
		NewDutyAssignment:                    "Test Station Name",
		AuthorizedOrigin:                     *ppm.PickupPostalCode,
		AuthorizedDestination:                *ppm.DestinationPostalCode,
		ShipmentPickupDate:                   *ppm.PlannedMoveDate,
		CurrentShipmentStatus:                ppm.Status,
		StorageTotalDays:                     int64(5),
		CurrentPaymentRequestClaim:           reimbursement.ID,
		CurrentPaymentRequestTransactionType: reimbursement.MethodOfReceipt,
		CurrentPaymentAmount:                 reimbursement.RequestedAmount,
		TrustedAgentName:                     bu.Name,
		TrustedAgentAuthorizationDate:        bu.CreatedAt,
		TrustedAgentEmail:                    bu.Email,
		TrustedAgentPhone:                    *bu.Phone,
	}

	fmt.Println(advanceWorksheet)

	if advanceWorksheet.FirstName == "" {
		t.Errorf("Well, that's regrettable and strange.")
	}
}

// Make a test that checks an advance worksheet struct missing something required
// func Test_MakeFlawedAdvanceWorksheetInfo() {

// }
