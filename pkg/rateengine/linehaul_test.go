package rateengine

func (suite *RateEngineSuite) Test_CheckDetermineMileage() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	mileage, err := engine.determineMileage("10024", "18209")
	if err != nil {
		t.Error("Unable to determine mileage: ", err)
	}
	if mileage != 1000 {
		t.Errorf("Determined mileage incorrectly. Expected 1000 got %d", mileage)
	}
}

func (suite *RateEngineSuite) Test_CheckBaseLinehaul() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	mileage := 3200
	cwt := 40

	blh, _ := engine.baseLinehaul(mileage, cwt)

	if blh != 128000 {
		t.Errorf("CWT should have been 12800000 but is %d.", blh)
	}
}

func (suite *RateEngineSuite) Test_CheckLinehaulFactors() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	linehaulFactor, err := engine.linehaulFactors(60, "18209")
	if err != nil {
		t.Error("Unable to determine linehaulFactor: ", err)
	}
	if linehaulFactor != 30.6 {
		t.Errorf("Determined linehaul factor incorrectly. Expected 30.6 got %f", linehaulFactor)
	}
}

func (suite *RateEngineSuite) Test_CheckShorthaulCharge() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	mileage := 799
	cwt := 40

	shc, _ := engine.shorthaulCharge(mileage, cwt)

	if shc != 31960 {
		t.Errorf("Shorthaul charge should have been 31960 but is %f.", shc)
	}
}

func (suite *RateEngineSuite) Test_CheckLinehaulChargeTotal() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	linehaulChargeTotal, err := engine.linehaulChargeTotal("10024", "94103")
	if err != nil {
		t.Error("Unable to determine linehaulChargeTotal: ", err)
	}
	if linehaulChargeTotal != 11812.036000 {
		t.Errorf("Determined linehaul factor incorrectly. Expected 11812.036000 got %f", linehaulChargeTotal)
	}
}
