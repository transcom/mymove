package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestElectronicOrderValidate() {
	newOrder := &models.ElectronicOrder{
		Edipi:        "1234567890",
		Issuer:       models.IssuerArmy,
		OrdersNumber: "8675309",
	}

	verrs, err := newOrder.Validate(nil)
	suite.NoError(err)
	suite.NoVerrs(verrs)
}

func (suite *ModelSuite) TestElectronicOrderValidations() {
	order := &models.ElectronicOrder{}

	var expErrors = map[string][]string{
		"orders_number": {"OrdersNumber can not be blank."},
		"edipi":         {"Edipi can not be blank.", "Edipi does not match the expected format."},
		"issuer":        {"Issuer is not in the list [air-force, army, coast-guard, marine-corps, navy]."},
	}

	suite.verifyValidationErrors(order, expErrors, nil)

	order.Edipi = "wrongfmt"

	expErrors["edipi"] = []string{"Edipi does not match the expected format."}

	suite.verifyValidationErrors(order, expErrors, nil)
}

func (suite *ModelSuite) TestCreateElectronicOrder() {
	newOrder := models.ElectronicOrder{
		Edipi:        "1234567890",
		Issuer:       models.IssuerArmy,
		OrdersNumber: "8675309",
	}

	verrs, err := models.CreateElectronicOrder(suite.DB(), &newOrder)
	suite.NoError(err)
	suite.NoVerrs(verrs)
}

func (suite *ModelSuite) TestCreateElectronicOrderWithRevision() {
	newOrder := models.ElectronicOrder{
		Edipi:        "1234567890",
		Issuer:       models.IssuerArmy,
		OrdersNumber: "8675309",
	}
	rev := models.ElectronicOrdersRevision{
		SeqNum:        0,
		GivenName:     "First",
		FamilyName:    "Last",
		Affiliation:   models.ElectronicOrdersAffiliationArmy,
		Paygrade:      models.PaygradeE1,
		Status:        models.ElectronicOrdersStatusAuthorized,
		DateIssued:    time.Now(),
		NoCostMove:    false,
		TdyEnRoute:    false,
		TourType:      models.TourTypeAccompanied,
		OrdersType:    models.ElectronicOrdersTypeSeparation,
		HasDependents: true,
	}
	verrs, err := models.CreateElectronicOrderWithRevision(suite.DB(), &newOrder, &rev)
	suite.NoError(err)
	suite.NoVerrs(verrs)
}

func (suite *ModelSuite) TestFetchElectronicOrderByID() {
	newOrder := models.ElectronicOrder{
		Edipi:        "1234567890",
		Issuer:       models.IssuerArmy,
		OrdersNumber: "8675309",
	}

	verrs, err := models.CreateElectronicOrder(suite.DB(), &newOrder)
	suite.NoError(err)
	suite.NoVerrs(verrs)

	retrievedOrder, err := models.FetchElectronicOrderByID(suite.DB(), newOrder.ID)
	suite.NoError(err)
	suite.Equal(newOrder.ID, retrievedOrder.ID)
	suite.Equal(newOrder.OrdersNumber, retrievedOrder.OrdersNumber)
	suite.Equal(newOrder.Edipi, retrievedOrder.Edipi)
	suite.Equal(newOrder.Issuer, retrievedOrder.Issuer)
}

func (suite *ModelSuite) TestFetchElectronicOrderByIssuerAndOrdersNum() {
	newOrder := models.ElectronicOrder{
		Edipi:        "1234567890",
		Issuer:       models.IssuerArmy,
		OrdersNumber: "8675309",
	}

	verrs, err := models.CreateElectronicOrder(suite.DB(), &newOrder)
	suite.NoError(err)
	suite.NoVerrs(verrs)

	retrievedOrder, err := models.FetchElectronicOrderByIssuerAndOrdersNum(suite.DB(), string(models.IssuerArmy), newOrder.OrdersNumber)
	suite.NoError(err)
	suite.Equal(newOrder.ID, retrievedOrder.ID)
	suite.Equal(newOrder.OrdersNumber, retrievedOrder.OrdersNumber)
	suite.Equal(newOrder.Edipi, retrievedOrder.Edipi)
	suite.Equal(newOrder.Issuer, retrievedOrder.Issuer)
}

func (suite *ModelSuite) TestFetchElectronicOrdersByEdipiAndIssuers() {
	edipi := "1234567890"
	newOrder1 := models.ElectronicOrder{
		Edipi:        edipi,
		Issuer:       models.IssuerArmy,
		OrdersNumber: "8675309",
	}

	verrs, err := models.CreateElectronicOrder(suite.DB(), &newOrder1)
	suite.NoError(err)
	suite.NoVerrs(verrs)

	newOrder2 := models.ElectronicOrder{
		Edipi:        edipi,
		Issuer:       models.IssuerAirForce,
		OrdersNumber: "5551234",
	}

	verrs, err = models.CreateElectronicOrder(suite.DB(), &newOrder2)
	suite.NoError(err)
	suite.NoVerrs(verrs)

	retrievedOrders, err := models.FetchElectronicOrdersByEdipiAndIssuers(suite.DB(), edipi, []string{string(models.IssuerArmy), string(models.IssuerAirForce)})
	suite.NoError(err)
	suite.Len(retrievedOrders, 2)
	ordersnumbers := []string{newOrder1.OrdersNumber, newOrder2.OrdersNumber}
	suite.Contains(ordersnumbers, retrievedOrders[0].OrdersNumber)
	suite.Contains(ordersnumbers, retrievedOrders[1].OrdersNumber)
	suite.NotEqual(retrievedOrders[0].OrdersNumber, retrievedOrders[1].OrdersNumber)
}
