package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeElectronicOrder returns a single ElectronicOrder with a single ElectronicOrdersRevision
func MakeElectronicOrder(db *pop.Connection, edipi string, issuer models.Issuer, ordersNumber string, affiliation models.ElectronicOrdersAffiliation) models.ElectronicOrder {
	order := models.ElectronicOrder{
		Edipi:        edipi,
		Issuer:       issuer,
		OrdersNumber: ordersNumber,
	}

	mustCreate(db, &order)

	rev := models.ElectronicOrdersRevision{
		ElectronicOrderID: order.ID,
		ElectronicOrder:   order,
		SeqNum:            0,
		GivenName:         "First",
		FamilyName:        "Last",
		Affiliation:       affiliation,
		Paygrade:          models.PaygradeE1,
		Status:            models.ElectronicOrdersStatusAuthorized,
		DateIssued:        time.Now(),
		NoCostMove:        false,
		TdyEnRoute:        false,
		TourType:          models.TourTypeAccompanied,
		OrdersType:        models.ElectronicOrdersTypeSeparation,
		HasDependents:     true,
	}

	mustCreate(db, &rev)

	return order
}
