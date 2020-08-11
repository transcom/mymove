package serviceparamvaluelookups

import (
	"fmt"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) setupPSILinehaulTestData(priceCents *unit.Cents, status *models.PaymentServiceItemStatus) (models.PaymentServiceItem, models.PaymentServiceItem) {
	code := models.ReServiceCodeDLH
	var localPriceCents unit.Cents
	if priceCents == nil {
		localPriceCents = unit.Cents(102400)
	} else {
		localPriceCents = *priceCents
	}

	var localStatus models.PaymentServiceItemStatus
	if status == nil {
		localStatus = models.PaymentServiceItemStatusRequested
	} else {
		localStatus = *status
	}
	psiLinehaulDomDLH := testdatagen.MakePaymentServiceItem(suite.DB(),
		testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &localPriceCents,
				Status:     localStatus,
			},
			ReService: models.ReService{
				Code: code,
				Name: string(code),
			},
		},
	)

	code = models.ReServiceCodeFSC
	psiLinehaulDomFSC := testdatagen.MakePaymentServiceItem(suite.DB(),
		testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: nil,
			},
			MTOServiceItem: models.MTOServiceItem{
				MTOShipmentID: psiLinehaulDomDLH.MTOServiceItem.MTOShipmentID,
			},
			ReService: models.ReService{
				Code: code,
				Name: string(code),
			},
		},
	)

	return psiLinehaulDomFSC, psiLinehaulDomDLH
}

func (suite *ServiceParamValueLookupsSuite) TestPSILinehaulDomLookup() {
	key := models.ServiceItemParamNamePSILinehaulDom.String()

	suite.T().Run("Domestic Linehaul Price has been calculated", func(t *testing.T) {

		psiLinehaulDom, expectedPSILinehaulDom := suite.setupPSILinehaulTestData(nil, nil)
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDom.MTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal(expectedPSILinehaulDom.ID.String(), valueStr)
	})

	suite.T().Run("Domestic Linehaul Price has been calculated twice use latest", func(t *testing.T) {
		psiLinehaulDom, psiLinehaulDomDLH := suite.setupPSILinehaulTestData(nil, nil)

		priceCents := unit.Cents(204800)
		psiLinehaulDomSecond := testdatagen.MakePaymentServiceItem(suite.DB(),
			testdatagen.Assertions{
				PaymentServiceItem: models.PaymentServiceItem{
					PriceCents: &priceCents,
				},
				MTOServiceItem: models.MTOServiceItem{
					MTOShipmentID: psiLinehaulDomDLH.MTOServiceItem.MTOShipmentID,
				},
				ReService: models.ReService{
					ID: psiLinehaulDomDLH.MTOServiceItem.ReServiceID,
				},
			},
		)

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDom.MTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal(psiLinehaulDomSecond.ID.String(), valueStr)
	})

	suite.T().Run("Domestic Linehaul Price has been calculated and Denied", func(t *testing.T) {
		price := unit.Cents(102400)
		status := models.PaymentServiceItemStatusDenied
		psiLinehaulDom, _ := suite.setupPSILinehaulTestData(&price, &status)

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDom.MTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
		expected := fmt.Sprintf(" failed ServiceParamValue PSI_LinehaulDomLookup with error couldn't find PaymentServiceItem for dom linehaul using paymentRequestID: %s and mtoServiceItemID: %s", psiLinehaulDom.PaymentRequestID, psiLinehaulDom.MTOServiceItemID)
		suite.Equal(expected, err.Error())
	})

	suite.T().Run("Domestic Linehaul Price has NOT been calculated", func(t *testing.T) {
		code := models.ReServiceCodeFSC
		psiLinehaulDomFSC := testdatagen.MakePaymentServiceItem(suite.DB(),
			testdatagen.Assertions{
				PaymentServiceItem: models.PaymentServiceItem{
					PriceCents: nil,
				},
				ReService: models.ReService{
					Code: code,
					Name: string(code),
				},
			},
		)

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDomFSC.MTOServiceItemID, psiLinehaulDomFSC.PaymentRequestID, psiLinehaulDomFSC.PaymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
		expected := fmt.Sprintf(" failed ServiceParamValue PSI_LinehaulDomLookup with error couldn't find PaymentServiceItem for dom linehaul using paymentRequestID: %s and mtoServiceItemID: %s", psiLinehaulDomFSC.PaymentRequestID, psiLinehaulDomFSC.MTOServiceItemID)
		suite.Equal(expected, err.Error())
	})
}

func (suite *ServiceParamValueLookupsSuite) TestPSILinehaulDomPriceLookup() {
	key := models.ServiceItemParamNamePSILinehaulDomPrice.String()

	suite.T().Run("Domestic Linehaul Price has been calculated", func(t *testing.T) {

		psiLinehaulDom, expectedPSILinehaulDom := suite.setupPSILinehaulTestData(nil, nil)
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDom.MTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal(expectedPSILinehaulDom.PriceCents.String(), valueStr)
	})

	suite.T().Run("Domestic Linehaul Price has been calculated twice use latest", func(t *testing.T) {
		psiLinehaulDom, psiLinehaulDomDLH := suite.setupPSILinehaulTestData(nil, nil)

		priceCents := unit.Cents(204800)
		psiLinehaulDomSecond := testdatagen.MakePaymentServiceItem(suite.DB(),
			testdatagen.Assertions{
				PaymentServiceItem: models.PaymentServiceItem{
					PriceCents: &priceCents,
				},
				MTOServiceItem: models.MTOServiceItem{
					MTOShipmentID: psiLinehaulDomDLH.MTOServiceItem.MTOShipmentID,
				},
				ReService: models.ReService{
					ID: psiLinehaulDomDLH.MTOServiceItem.ReServiceID,
				},
			},
		)

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDom.MTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal(psiLinehaulDomSecond.PriceCents.String(), valueStr)
	})

	suite.T().Run("Domestic Linehaul Price has been calculated and Denied", func(t *testing.T) {
		price := unit.Cents(102400)
		status := models.PaymentServiceItemStatusDenied
		psiLinehaulDom, _ := suite.setupPSILinehaulTestData(&price, &status)

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDom.MTOServiceItemID, psiLinehaulDom.PaymentRequestID, psiLinehaulDom.PaymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
		expected := fmt.Sprintf(" failed ServiceParamValue PSI_LinehaulDomPriceLookup with error couldn't find PaymentServiceItem for dom linehaul using paymentRequestID: %s and mtoServiceItemID: %s", psiLinehaulDom.PaymentRequestID, psiLinehaulDom.MTOServiceItemID)
		suite.Equal(expected, err.Error())
	})

	suite.T().Run("Domestic Linehaul Price has NOT been calculated", func(t *testing.T) {
		code := models.ReServiceCodeFSC
		psiLinehaulDomFSC := testdatagen.MakePaymentServiceItem(suite.DB(),
			testdatagen.Assertions{
				PaymentServiceItem: models.PaymentServiceItem{
					PriceCents: nil,
				},
				ReService: models.ReService{
					Code: code,
					Name: string(code),
				},
			},
		)

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, psiLinehaulDomFSC.MTOServiceItemID, psiLinehaulDomFSC.PaymentRequestID, psiLinehaulDomFSC.PaymentRequest.MoveTaskOrderID)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
		expected := fmt.Sprintf(" failed ServiceParamValue PSI_LinehaulDomPriceLookup with error couldn't find PaymentServiceItem for dom linehaul using paymentRequestID: %s and mtoServiceItemID: %s", psiLinehaulDomFSC.PaymentRequestID, psiLinehaulDomFSC.MTOServiceItemID)
		suite.Equal(expected, err.Error())
	})
}
