package rateengine

func (suite *RateEngineSuite) Test_CheckDetermineMileage() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	mileage, err := engine.determineMileage("10024", "18209")
	if err != nil {
		t.Error("Unable to determine mileage: ", err)
	}
	expected := 1000
	if mileage != expected {
		t.Errorf("Determined mileage incorrectly. Expected %d, got %d", expected, mileage)
	}
}

func (suite *RateEngineSuite) Test_CheckBaseLinehaul() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	mileage := 3200
	cwt := 40

	blh, _ := engine.baseLinehaul(mileage, cwt)
	expected := 128000
	if blh != 128000 {
		t.Errorf("CWT should have been %d but is %d.", expected, blh)
	}
}

func (suite *RateEngineSuite) Test_CheckLinehaulFactors() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	linehaulFactor, err := engine.linehaulFactors(60, "18209")
	if err != nil {
		t.Error("Unable to determine linehaulFactor: ", err)
	}
	expected := 3060
	if linehaulFactor != expected {
		t.Errorf("Determined linehaul factor incorrectly. Expected %d, got %d", expected, linehaulFactor)
	}
}

func (suite *RateEngineSuite) Test_CheckShorthaulCharge() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	mileage := 799
	cwt := 40

	shc, _ := engine.shorthaulCharge(mileage, cwt)
	expected := 31960
	if shc != expected {
		t.Errorf("Shorthaul charge should have been %d, but is %d.", expected, shc)
	}
}

func (suite *RateEngineSuite) Test_CheckLinehaulChargeTotal() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)
	linehaulChargeTotal, err := engine.linehaulChargeTotal("10024", "94103")
	if err != nil {
		t.Error("Unable to determine linehaulChargeTotal: ", err)
	}
	expected := 13003
	if linehaulChargeTotal != expected {
		t.Errorf("Determined linehaul factor incorrectly. Expected %d, got %d", expected, linehaulChargeTotal)
	}
}
