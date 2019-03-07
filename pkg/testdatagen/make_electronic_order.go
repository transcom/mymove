package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeElectronicOrder returns a single ElectronicOrder with a single ElectronicOrdersRevision
func MakeElectronicOrder(db *pop.Connection, edipi string, issuer ordersmessages.Issuer, ordersNumber string, affiliation ordersmessages.Affiliation) models.ElectronicOrder {
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
		Paygrade:          ordersmessages.RankE1,
		Status:            ordersmessages.StatusAuthorized,
		DateIssued:        time.Now(),
		NoCostMove:        false,
		TdyEnRoute:        false,
		TourType:          ordersmessages.TourTypeAccompanied,
		OrdersType:        ordersmessages.OrdersTypeSeparation,
		HasDependents:     true,
	}

	mustCreate(db, &rev)

	return order
}
