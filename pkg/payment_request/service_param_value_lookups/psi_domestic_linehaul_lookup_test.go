package serviceparamvaluelookups

// func (suite *ServiceParamValueLookupsSuite) setupPSILinehaulTestData(priceCents *unit.Cents, status *models.PaymentServiceItemStatus) (models.PaymentServiceItem, models.PaymentServiceItem) {
// 	code := models.ReServiceCodeDLH
// 	var localPriceCents unit.Cents
// 	if priceCents == nil {
// 		localPriceCents = unit.Cents(102400)
// 	} else {
// 		localPriceCents = *priceCents
// 	}
//
// 	var localStatus models.PaymentServiceItemStatus
// 	if status == nil {
// 		localStatus = models.PaymentServiceItemStatusRequested
// 	} else {
// 		localStatus = *status
// 	}
// 	psiLinehaulDomDLH := testdatagen.MakePaymentServiceItem(suite.DB(),
// 		testdatagen.Assertions{
// 			PaymentServiceItem: models.PaymentServiceItem{
// 				PriceCents: &localPriceCents,
// 				Status:     localStatus,
// 			},
// 			ReService: models.ReService{
// 				Code: code,
// 				Name: string(code),
// 			},
// 		},
// 	)
//
// 	code = models.ReServiceCodeFSC
// 	psiLinehaulDomFSC := testdatagen.MakePaymentServiceItem(suite.DB(),
// 		testdatagen.Assertions{
// 			PaymentServiceItem: models.PaymentServiceItem{
// 				PriceCents: nil,
// 			},
// 			MTOServiceItem: models.MTOServiceItem{
// 				MTOShipmentID: psiLinehaulDomDLH.MTOServiceItem.MTOShipmentID,
// 			},
// 			ReService: models.ReService{
// 				Code: code,
// 				Name: string(code),
// 			},
// 		},
// 	)
//
// 	return psiLinehaulDomFSC, psiLinehaulDomDLH
// }
//
// func (suite *ServiceParamValueLookupsSuite) TestPSILinehaulDomLookup() {
// 	key := models.ServiceItemParamNamePSILinehaulDom.String()
//
// 	suite.Run("Domestic Linehaul Price has been calculated", func() {
//
// 		psiLinehaulDom, expectedPSILinehaulDom := suite.setupPSILinehaulTestData(nil, nil)
// 		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDom.MTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
//
// 		valueStr, err := paramLookup.ServiceParamValue(suite.TestAppContext(), key)
// 		suite.FatalNoError(err)
// 		suite.Equal(expectedPSILinehaulDom.ID.String(), valueStr)
// 	})
//
// 	suite.Run("Domestic Linehaul Price has been calculated twice use latest", func() {
// 		psiLinehaulDom, psiLinehaulDomDLH := suite.setupPSILinehaulTestData(nil, nil)
//
// 		priceCents := unit.Cents(204800)
// 		psiLinehaulDomSecond := testdatagen.MakePaymentServiceItem(suite.DB(),
// 			testdatagen.Assertions{
// 				PaymentServiceItem: models.PaymentServiceItem{
// 					PriceCents: &priceCents,
// 				},
// 				MTOServiceItem: models.MTOServiceItem{
// 					MTOShipmentID: psiLinehaulDomDLH.MTOServiceItem.MTOShipmentID,
// 				},
// 				ReService: models.ReService{
// 					ID: psiLinehaulDomDLH.MTOServiceItem.ReServiceID,
// 				},
// 			},
// 		)
//
// 		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDom.MTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
//
// 		valueStr, err := paramLookup.ServiceParamValue(suite.TestAppContext(), key)
// 		suite.FatalNoError(err)
// 		suite.Equal(psiLinehaulDomSecond.ID.String(), valueStr)
// 	})
//
// 	suite.Run("Domestic Linehaul Price has been calculated and Denied", func() {
// 		price := unit.Cents(102400)
// 		status := models.PaymentServiceItemStatusDenied
// 		psiLinehaulDom, _ := suite.setupPSILinehaulTestData(&price, &status)
//
// 		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDom.MTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
//
// 		_, err := paramLookup.ServiceParamValue(suite.TestAppContext(), key)
// 		suite.Error(err)
// 		expected := fmt.Sprintf(" failed ServiceParamValue PSI_LinehaulDomLookup with error couldn't find PaymentServiceItem for dom linehaul using paymentRequestID: %s and mtoServiceItemID: %s", psiLinehaulDom.PaymentRequestID, psiLinehaulDom.MTOServiceItemID)
// 		suite.Equal(expected, err.Error())
// 	})
//
// 	suite.Run("Invalid MTO Service ID", func() {
// 		psiLinehaulDom, _ := suite.setupPSILinehaulTestData(nil, nil)
//
// 		invalidMTOServiceItemID := uuid.Must(uuid.NewV4())
// 		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, invalidMTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
//
// 		_, err := paramLookup.ServiceParamValue(suite.TestAppContext(), key)
// 		suite.Error(err)
// 		expected := fmt.Sprintf(" failed ServiceParamValue PSI_LinehaulDomLookup with error id: %s not found looking for MTOServiceItemID", invalidMTOServiceItemID)
// 		suite.Equal(expected, err.Error())
// 	})
//
// 	suite.Run("Domestic Linehaul Price has NOT been calculated", func() {
// 		code := models.ReServiceCodeFSC
// 		psiLinehaulDomFSC := testdatagen.MakePaymentServiceItem(suite.DB(),
// 			testdatagen.Assertions{
// 				PaymentServiceItem: models.PaymentServiceItem{
// 					PriceCents: nil,
// 				},
// 				ReService: models.ReService{
// 					Code: code,
// 					Name: string(code),
// 				},
// 			},
// 		)
//
// 		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDomFSC.MTOServiceItemID, psiLinehaulDomFSC.PaymentRequestID, psiLinehaulDomFSC.PaymentRequest.MoveTaskOrderID)
//
// 		_, err := paramLookup.ServiceParamValue(suite.TestAppContext(), key)
// 		suite.Error(err)
// 		expected := fmt.Sprintf(" failed ServiceParamValue PSI_LinehaulDomLookup with error couldn't find PaymentServiceItem for dom linehaul using paymentRequestID: %s and mtoServiceItemID: %s", psiLinehaulDomFSC.PaymentRequestID, psiLinehaulDomFSC.MTOServiceItemID)
// 		suite.Equal(expected, err.Error())
// 	})
// }
//
// func (suite *ServiceParamValueLookupsSuite) TestPSILinehaulDomPriceLookup() {
// 	key := models.ServiceItemParamNamePSILinehaulDomPrice.String()
//
// 	suite.Run("Domestic Linehaul Price has been calculated", func() {
//
// 		psiLinehaulDom, expectedPSILinehaulDom := suite.setupPSILinehaulTestData(nil, nil)
// 		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDom.MTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
//
// 		valueStr, err := paramLookup.ServiceParamValue(suite.TestAppContext(), key)
// 		suite.FatalNoError(err)
// 		suite.Equal(expectedPSILinehaulDom.PriceCents.String(), valueStr)
// 	})
//
// 	suite.Run("Domestic Linehaul Price has been calculated twice use latest", func() {
// 		psiLinehaulDom, psiLinehaulDomDLH := suite.setupPSILinehaulTestData(nil, nil)
//
// 		priceCents := unit.Cents(204800)
// 		psiLinehaulDomSecond := testdatagen.MakePaymentServiceItem(suite.DB(),
// 			testdatagen.Assertions{
// 				PaymentServiceItem: models.PaymentServiceItem{
// 					PriceCents: &priceCents,
// 				},
// 				MTOServiceItem: models.MTOServiceItem{
// 					MTOShipmentID: psiLinehaulDomDLH.MTOServiceItem.MTOShipmentID,
// 				},
// 				ReService: models.ReService{
// 					ID: psiLinehaulDomDLH.MTOServiceItem.ReServiceID,
// 				},
// 			},
// 		)
//
// 		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDom.MTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
//
// 		valueStr, err := paramLookup.ServiceParamValue(suite.TestAppContext(), key)
// 		suite.FatalNoError(err)
// 		suite.Equal(psiLinehaulDomSecond.PriceCents.String(), valueStr)
// 	})
//
// 	suite.Run("Domestic Linehaul Price has been calculated and Denied", func() {
// 		price := unit.Cents(102400)
// 		status := models.PaymentServiceItemStatusDenied
// 		psiLinehaulDom, _ := suite.setupPSILinehaulTestData(&price, &status)
//
// 		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDom.MTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
//
// 		_, err := paramLookup.ServiceParamValue(suite.TestAppContext(), key)
// 		suite.Error(err)
// 		expected := fmt.Sprintf(" failed ServiceParamValue PSI_LinehaulDomPriceLookup with error couldn't find PaymentServiceItem for dom linehaul using paymentRequestID: %s and mtoServiceItemID: %s", psiLinehaulDom.PaymentRequestID, psiLinehaulDom.MTOServiceItemID)
// 		suite.Equal(expected, err.Error())
// 	})
//
// 	suite.Run("Invalid MTO Service ID", func() {
// 		psiLinehaulDom, _ := suite.setupPSILinehaulTestData(nil, nil)
//
// 		invalidMTOServiceItemID := uuid.Must(uuid.NewV4())
// 		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, invalidMTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
//
// 		_, err := paramLookup.ServiceParamValue(suite.TestAppContext(), key)
// 		suite.Error(err)
// 		expected := fmt.Sprintf(" failed ServiceParamValue PSI_LinehaulDomPriceLookup with error id: %s not found looking for MTOServiceItemID", invalidMTOServiceItemID)
// 		suite.Equal(expected, err.Error())
// 	})
//
// 	suite.Run("Domestic Linehaul Price has NOT been calculated", func() {
// 		code := models.ReServiceCodeFSC
// 		psiLinehaulDomFSC := testdatagen.MakePaymentServiceItem(suite.DB(),
// 			testdatagen.Assertions{
// 				PaymentServiceItem: models.PaymentServiceItem{
// 					PriceCents: nil,
// 				},
// 				ReService: models.ReService{
// 					Code: code,
// 					Name: string(code),
// 				},
// 			},
// 		)
//
// 		paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDomFSC.MTOServiceItemID, psiLinehaulDomFSC.PaymentRequestID, psiLinehaulDomFSC.PaymentRequest.MoveTaskOrderID)
//
// 		_, err := paramLookup.ServiceParamValue(suite.TestAppContext(), key)
// 		suite.Error(err)
// 		expected := fmt.Sprintf(" failed ServiceParamValue PSI_LinehaulDomPriceLookup with error couldn't find PaymentServiceItem for dom linehaul using paymentRequestID: %s and mtoServiceItemID: %s", psiLinehaulDomFSC.PaymentRequestID, psiLinehaulDomFSC.MTOServiceItemID)
// 		suite.Equal(expected, err.Error())
// 	})
// }
