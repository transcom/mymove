package testdatagen

import (
	"log"
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeOrder creates a single Move and associated User.
func MakeOrder(db *pop.Connection) (models.Order, error) {
	sm, err := MakeExtendedServiceMember(db)
	if err != nil {
		return models.Order{}, err
	}

	return MakeOrderForServiceMember(db, sm)
}

// MakeOrderForServiceMember makes an order for a given service member
func MakeOrderForServiceMember(db *pop.Connection, sm models.ServiceMember) (models.Order, error) {
	var order models.Order

	station := MakeAnyDutyStation(db)

	document, _ := MakeDocument(db, &sm, models.UploadedOrdersDocumentName)

	order = models.Order{
		ServiceMemberID:  sm.ID,
		ServiceMember:    sm,
		NewDutyStationID: station.ID,
		NewDutyStation:   station,
		IssueDate:        time.Date(2018, time.March, 15, 0, 0, 0, 0, time.UTC),
		ReportByDate:     time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC),
		OrdersType:       internalmessages.OrdersTypePERMANENTCHANGEOFSTATION,
		HasDependents:    true,
		SpouseHasProGear: true,
		UploadedOrdersID: document.ID,
		UploadedOrders:   document,
		Status:           models.OrderStatusDRAFT,
	}

	verrs, err := db.ValidateAndSave(&order)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() != 0 {
		log.Panic(verrs.Error())
	}

	return order, err
}
