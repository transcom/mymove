package models_test

import (
	"github.com/transcom/mymove/pkg/auth"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestInvoiceValidations() {
	invoice := &Invoice{}

	expErrors := map[string][]string{
		"status":         {"Status can not be blank."},
		"approver_id":    {"ApproverID can not be blank."},
		"invoice_number": {"InvoiceNumber can not be blank."},
		"invoiced_date":  {"InvoicedDate can not be blank."},
		"shipment_id":    {"ShipmentID can not be blank."},
	}

	suite.verifyValidationErrors(invoice, expErrors)
}

func (suite *ModelSuite) TestFetchInvoice() {
	// Data setup
	numTspUsers := 1
	numShipments := 1
	numShipmentOfferSplit := []int{1}
	status := []ShipmentStatus{ShipmentStatusDELIVERED}
	tspUsers, shipments, _, err := testdatagen.CreateShipmentOfferData(suite.db, numTspUsers, numShipments, numShipmentOfferSplit, status)
	suite.NoError(err)

	shipment := shipments[0]

	authedTspUser := tspUsers[0]
	unverifiedTspUser := testdatagen.MakeDefaultTspUser(suite.db)
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)

	// Invoice tied to a shipment of authed tsp user
	invoice := testdatagen.MakeInvoice(suite.db, testdatagen.Assertions{
		Invoice: Invoice{
			ShipmentID: shipment.ID,
		},
	})

	// When: office user tries to access
	session := &auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *officeUser.UserID,
		OfficeUserID:    officeUser.ID,
	}

	// Then: invoice is returned
	extantInvoice, err := FetchInvoice(suite.db, session, invoice.ID)
	suite.Nil(err)
	if suite.NoError(err) {
		suite.Equal(extantInvoice.ID, invoice.ID)
	}
	// When: Unverified TSP tries to access
	session = &auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *unverifiedTspUser.UserID,
		TspUserID:       unverifiedTspUser.ID,
	}
	// Then: fetch forbidden returned
	extantInvoice, err = FetchInvoice(suite.db, session, invoice.ID)
	suite.Equal("FETCH_FORBIDDEN", err.Error())

	// When: authed TSP tries to access
	session = &auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *authedTspUser.UserID,
		TspUserID:       authedTspUser.ID,
	}
	// Then: invoice is returned
	extantInvoice, err = FetchInvoice(suite.db, session, invoice.ID)
	suite.Nil(err)
	if suite.NoError(err) {
		suite.Equal(extantInvoice.ID, invoice.ID)
	}
	// When: Service Member tries to access
	session = &auth.Session{
		ApplicationName: auth.MyApp,
		UserID:          *authedTspUser.UserID,
		ServiceMemberID: authedTspUser.ID,
	}
	// Then: Fetch Forbidden returned
	extantInvoice, err = FetchInvoice(suite.db, session, invoice.ID)
	suite.Equal("FETCH_FORBIDDEN", err.Error())

}

func (suite *ModelSuite) TestFetchInvoicesForShipment() {
	invoice1 := testdatagen.MakeDefaultInvoice(suite.db)
	testdatagen.MakeDefaultInvoice(suite.db)

	// Then: invoice is returned
	extantInvoices, err := FetchInvoicesForShipment(suite.db, invoice1.ShipmentID)
	if suite.NoError(err) {
		suite.Len(extantInvoices, 1)
		suite.Equal(extantInvoices[0].ID, invoice1.ID)
	}
}
