package models_test

import (
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestGetEntitlementWithValidValues() {
	E1 := models.ServiceMemberGradeE1
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION

	suite.Run("E1 with dependents", func() {
		E1FullLoad := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(8000, E1FullLoad.TotalWeightSelfPlusDependents)
	})

	suite.Run("E1 without dependents", func() {
		E1Solo := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(5000, E1Solo.TotalWeightSelf)
	})

	suite.Run("E1 Pro Gear", func() {
		E1ProGear := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(2000, E1ProGear.ProGearWeight)
	})

	suite.Run("E1 Pro Gear Spouse", func() {
		E1ProGearSpouse := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(500, E1ProGearSpouse.ProGearWeightSpouse)
	})
}

func (suite *ModelSuite) TestGetEntitlementByOrdersTypeWithValidValues() {
	E1 := models.ServiceMemberGradeE1
	ordersType := internalmessages.OrdersTypeSTUDENTTRAVEL

	suite.Run("Student Travel with dependents", func() {
		STFullLoad := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(350, STFullLoad.TotalWeightSelfPlusDependents)
	})

	suite.Run("Student Travel without dependents", func() {
		STSolo := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(350, STSolo.TotalWeightSelf)
	})

	suite.Run("Student Travel Pro Gear", func() {
		STProGear := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(0, STProGear.ProGearWeight)
	})

	suite.Run("Student Travel Pro Gear Spouse", func() {
		STProGearSpouse := models.GetWeightAllotment(E1, ordersType)
		suite.Assertions.Equal(0, STProGearSpouse.ProGearWeightSpouse)
	})
}
