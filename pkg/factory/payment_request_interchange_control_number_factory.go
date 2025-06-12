package factory

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/db/sequence"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildPaymentRequestToInterchangeControlNumber creates a single PaymentRequest and PaymentRequestToInterchangeControlNumber
func BuildPaymentRequestToInterchangeControlNumber(db *pop.Connection, customs []Customization, traits []Trait) models.PaymentRequestToInterchangeControlNumber {
	customs = setupCustomizations(customs, traits)

	// Find paymentRequestInterchangeControlNumber customization and extract the custom paymentRequestInterchangeControlNumber
	var cPr2icn models.PaymentRequestToInterchangeControlNumber
	if result := findValidCustomization(customs, PaymentRequestToInterchangeControlNumber); result != nil {
		cPr2icn = result.Model.(models.PaymentRequestToInterchangeControlNumber)
		if result.LinkOnly {
			return cPr2icn
		}
	}

	paymentRequest := BuildPaymentRequest(db, customs, traits)

	icnSequencer, err := sequence.NewRandomSequencer(ediinvoice.ICNRandomMin, ediinvoice.ICNRandomMax)
	if err != nil {
		log.Panic(fmt.Errorf("errors encountered creating random sequencer: %v", err))
	}

	// for now, hack together an appcontext, so we don't have to change all of testdatagen
	appCtx := appcontext.NewAppContext(db, nil, nil, nil)
	icn, err := icnSequencer.NextVal(appCtx)
	if err != nil {
		log.Panic(fmt.Errorf("errors encountered getting random interchange control number: %v", err))
	}

	pr2icn := models.PaymentRequestToInterchangeControlNumber{
		PaymentRequestID:         paymentRequest.ID,
		InterchangeControlNumber: int(icn),
		EDIType:                  models.EDIType858,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&pr2icn, cPr2icn)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &pr2icn)
	}

	return pr2icn

}
