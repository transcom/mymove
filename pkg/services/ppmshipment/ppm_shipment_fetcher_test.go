package ppmshipment

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *PPMShipmentSuite) TestFetchPPMShipment() {
	suite.Run("FindPPMShipmentWithDocument - document belongs to weight ticket", func() {
		weightTicket := testdatagen.MakeDefaultWeightTicket(suite.DB())

		err := FindPPMShipmentWithDocument(suite.AppContextForTest(), weightTicket.PPMShipmentID, weightTicket.EmptyDocumentID)
		suite.NoError(err, "expected to find PPM Shipment for empty weight document")

		err = FindPPMShipmentWithDocument(suite.AppContextForTest(), weightTicket.PPMShipmentID, weightTicket.FullDocumentID)
		suite.NoError(err, "expected to find PPM Shipment for full weight document")

		err = FindPPMShipmentWithDocument(suite.AppContextForTest(), weightTicket.PPMShipmentID, weightTicket.FullDocumentID)
		suite.NoError(err, "expected to find PPM Shipment for trailer ownership document")
	})

	suite.Run("FindPPMShipmentWithDocument - document belongs to pro gear", func() {
		proGear := testdatagen.MakeDefaultProgearWeightTicket(suite.DB())

		err := FindPPMShipmentWithDocument(suite.AppContextForTest(), proGear.PPMShipmentID, proGear.DocumentID)
		suite.NoError(err, "expected to find PPM Shipment for full weight document")
	})

	suite.Run("FindPPMShipmentWithDocument - document belongs to moving expenses", func() {
		movingExpense := testdatagen.MakeDefaultMovingExpense(suite.DB())

		err := FindPPMShipmentWithDocument(suite.AppContextForTest(), movingExpense.PPMShipmentID, movingExpense.DocumentID)
		suite.NoError(err, "expected to find PPM Shipment for moving expense document")
	})

	suite.Run("FindPPMShipmentWithDocument - document not found", func() {
		weightTicket := testdatagen.MakeDefaultWeightTicket(suite.DB())
		testdatagen.MakeProgearWeightTicket(suite.DB(), testdatagen.Assertions{PPMShipment: weightTicket.PPMShipment})
		testdatagen.MakeMovingExpense(suite.DB(), testdatagen.Assertions{PPMShipment: weightTicket.PPMShipment})

		documentID := uuid.Must(uuid.NewV4())
		err := FindPPMShipmentWithDocument(suite.AppContextForTest(), weightTicket.PPMShipmentID, documentID)
		suite.Error(err, "expected to return not found error for unknown document id")
		suite.Equal(fmt.Sprintf("ID: %s not found document does not exist for the given shipment", documentID), err.Error())
	})
}
