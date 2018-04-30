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
	var order models.Order

	sm, err := MakeServiceMember(db)
	if err != nil {
		return order, err
	}

	station := MakeAnyDutyStation(db)
	if err != nil {
		return order, err
	}

	order = models.Order{
		ServiceMemberID:  sm.ID,
		ServiceMember:    sm,
		NewDutyStationID: station.ID,
		NewDutyStation:   station,
		IssueDate:        time.Date(2018, time.March, 15, 0, 0, 0, 0, time.UTC),
		ReportByDate:     time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC),
		OrdersType:       internalmessages.OrdersTypeAccession,
		HasDependents:    true,
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
