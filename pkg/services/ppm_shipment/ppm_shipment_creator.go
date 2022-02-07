package ppm_shipment

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
)

builder := createMTOSHipmentQueryBuilder
fetcher := services.Fetcher
moveRouter := services.moveRouter

func createPPMShipment(appCtx appcontext.AppContext) error {
	// Make sure we do this whole process in a transaction so partial changes do not get made committed
	// in the event of an error.
	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
	// Start a transaction that will create a Shipment, then create a PPM
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		var err error
		createShipment, err := mtoShipmentCreator.CreateMTOShipment(txnAppCtx) // PASS IN ALL THE OBJECTS CreateMTOShipment needs
		// Check that mtoshipment is created. If not, bail out.
		// CREATE THE PPM_SHIPEMNT. Use the same txnAppCtx
		// Create the model object for ppm_shipment and save that data to the DB
		// Check that ppm_shipment is created. If not, bail out adn return the error object

		return err
	})
	if transactionError != nil {
		return transactionError
	}
	return nil
}
