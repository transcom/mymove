package rateengine

func (suite *RateEngineSuite) Test_originServiceFee() {
	t := suite.T()
	engine := NewRateEngine(suite.db, suite.logger)

	fee, err := engine.originServiceFee(2500, "18209")
	if err != nil {
		t.Fatalf("failed to calculate origin service fee: %s", err)
	}

	expected := 1.35
	if fee != expected {
		t.Errorf("wrong origin service fee: expected %f, got %f", expected, fee)
	}
}
