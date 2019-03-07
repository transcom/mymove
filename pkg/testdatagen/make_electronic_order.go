package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeDefaultElectronicOrder returns a single ElectronicOrder with default values and a single ElectronicOrdersRevision
func MakeDefaultElectronicOrder(db *pop.Connection, issuer ordersmessages.Issuer, affiliation ordersmessages.Affiliation) models.ElectronicOrder {
	order := models.ElectronicOrder{
		Edipi:        "1234567890",
		Issuer:       issuer,
		OrdersNumber: "8675309",
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
