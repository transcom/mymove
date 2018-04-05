package rateengine

func (suite *RateEngineSuite) Test_CheckServiceFee() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)

	fee, err := engine.serviceFee(25, "18209")
	if err != nil {
		t.Fatalf("failed to calculate service fee: %s", err)
	}

	expected := float64(97)
	if fee != expected {
		t.Errorf("wrong service fee: expected %f, got %f", expected, fee)
	}
}

func (suite *RateEngineSuite) Test_CheckFullPack() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)

	fee, err := engine.fullPack(25, "18209")
	if err != nil {
		t.Fatalf("failed to calculate full pack fee: %s", err)
	}

	expected := float64(1375)
	if fee != expected {
		t.Errorf("wrong full pack fee: expected %f, got %f", expected, fee)
	}
}

func (suite *RateEngineSuite) Test_CheckFullUnpack() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)

	fee, err := engine.fullUnpack(25, "18209")
	if err != nil {
		t.Fatalf("failed to calculate full pack fee: %s", err)
	}

	expected := float64(125)
	if fee != expected {
		t.Errorf("wrong full pack fee: expected %f, got %f", expected, fee)
	}
}

func (suite *RateEngineSuite) Test_CheckNonLinehaulChargeTotal() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)

	fee, err := engine.nonLinehaulChargeTotal("10024", "18209", 0.5)
	if err != nil {
		t.Fatalf("failed to calculate full pack fee: %s", err)
	}
	// (155.2 + 155.2 + 2200 + 200) * .5
	expected := float64(1355.2)
	if fee != expected {
		t.Errorf("wrong non-linehaul charge total: expected %f, got %f", expected, fee)
	}
}
