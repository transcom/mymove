package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeElectronicOrder returns a single ElectronicOrder with a single ElectronicOrdersRevision
func MakeElectronicOrder(db *pop.Connection, assertions Assertions) models.ElectronicOrder {
	//func MakeElectronicOrder(db *pop.Connection, edipi string, issuer models.Issuer, ordersNumber string, affiliation models.ElectronicOrdersAffiliation) models.ElectronicOrder {

	order := models.ElectronicOrder{
		Edipi:        "1234567890",
		Issuer:       models.IssuerAirForce,
		OrdersNumber: "8675309",
	}

	mergeModels(&order, assertions.ElectronicOrder)
	mustCreate(db, &order)

	rev := models.ElectronicOrdersRevision{
		ElectronicOrderID: order.ID,
		ElectronicOrder:   order,
		SeqNum:            0,
		GivenName:         "First",
		FamilyName:        "Last",
		Affiliation:       models.ElectronicOrdersAffiliationAirForce,
		Paygrade:          models.PaygradeE1,
		Status:            models.ElectronicOrdersStatusAuthorized,
		DateIssued:        time.Now(),
		NoCostMove:        false,
		TdyEnRoute:        false,
		TourType:          models.TourTypeAccompanied,
		OrdersType:        models.ElectronicOrdersTypeSeparation,
		HasDependents:     true,
	}

	mergeModels(&rev, assertions.ElectronicOrdersRevision)
	mustCreate(db, &rev)

	return order
}

// MakeDefaultElectronicOrder return an ElectronicOrder with default values (including a default ElectronicOrdersRevision)
func MakeDefaultElectronicOrder(db *pop.Connection) models.ElectronicOrder {
	return MakeElectronicOrder(db, Assertions{})
}
