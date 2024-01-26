package ppmcloseout

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	serviceparamvaluelookups "github.com/transcom/mymove/pkg/payment_request/service_param_value_lookups"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/unit"
)

type ppmCloseoutFetcher struct {
	planner              route.Planner
	paymentRequestHelper paymentrequesthelper.Helper
}

func NewPPMCloseoutFetcher(planner route.Planner, paymentRequestHelper paymentrequesthelper.Helper) services.PPMCloseoutFetcher {
	return &ppmCloseoutFetcher{
		planner:              planner,
		paymentRequestHelper: paymentRequestHelper,
	}
}

func (p *ppmCloseoutFetcher) GetPPMCloseout(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.PPMCloseout, error) {
	var ppmCloseoutObj models.PPMCloseout
	var ppmShipment models.PPMShipment
	var mtoShipment models.MTOShipment
	var err error

	err = appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"ID",
			"ShipmentID",
			"ExpectedDepartureDate",
			"ActualMoveDate",
			"EstimatedWeight",
			"HasProGear",
			"ProGearWeight",
			"SpouseProGearWeight",
			"FinalIncentive",
		).
		Find(&ppmShipment, ppmShipmentID)

	// Check if PPM shipment is in "NEEDS_PAYMENT_APPROVAL" status, if not, it's not ready for closeout, so return
	if ppmShipment.Status != models.PPMShipmentStatusNeedsPaymentApproval {
		return nil, apperror.NewPPMNotReadyForCloseoutError(ppmShipmentID, "")
	}

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(ppmShipmentID, "while looking for PPMShipment")
		default:
			return nil, apperror.NewQueryError("PPMShipment", err, "unable to find PPMShipment")
		}
	}

	var expenseItems []models.MovingExpense
	storageExpensePrice := unit.Cents(0)

	err = appCtx.DB().Where("ppm_shipment_id = ?", ppmShipmentID).All(&expenseItems)
	if err != nil {
		return nil, err
	}

	for _, movingExpense := range expenseItems {
		if movingExpense.MovingExpenseType != nil && *movingExpense.MovingExpenseType == models.MovingExpenseReceiptTypeStorage {
			storageExpensePrice += *movingExpense.Amount
		}
	}

	mtoShipmentID := &ppmShipment.ShipmentID
	err = appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"ID",
			"ScheduledPickupDate",
			"ActualPickupDate",
			"Distance",
			"PrimeActualWeight",
		).
		Find(&mtoShipment, mtoShipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(*mtoShipmentID, "while looking for MTOShipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "unable to find MTOShipment")
		}
	}

	// Get all DLH, FSC, DOP, DDP, DPK, and DUPK service items for the shipment
	var serviceItemsToPrice []models.MTOServiceItem
	logger := appCtx.Logger()
	idString := ppmShipment.ShipmentID.String()
	fmt.Print(idString)
	err = appCtx.DB().Where("mto_shipment_id = ?", ppmShipment.ShipmentID).All(&serviceItemsToPrice)
	if err != nil {
		return nil, err
	}
	serviceItemsToPrice = ppmshipment.BaseServiceItems(ppmShipment.ShipmentID)
	logger.Debug(fmt.Sprintf("serviceItemsToPrice %+v", serviceItemsToPrice))
	// itemPricer := ghcrateengine.NewServiceItemPricer()
	contractDate := ppmShipment.ExpectedDepartureDate
	contract, err := serviceparamvaluelookups.FetchContract(appCtx, contractDate)
	if err != nil {
		return nil, err
	}

	paramsForServiceItems, paramErr := p.paymentRequestHelper.FetchServiceParamsForServiceItems(appCtx, serviceItemsToPrice)
	if paramErr != nil {
		return nil, paramErr
	}
	var totalPrice, packPrice, unpackPrice, destinationPrice, originPrice unit.Cents
	var totalWeight unit.Pound
	var ppmToMtoShipment models.MTOShipment

	if len(ppmShipment.WeightTickets) >= 1 {
		for _, weightTicket := range ppmShipment.WeightTickets {
			if weightTicket.Status != nil && *weightTicket.Status == models.PPMDocumentStatusRejected {
				totalWeight += 0
			} else if weightTicket.AdjustedNetWeight != nil {
				totalWeight += *weightTicket.AdjustedNetWeight
			} else if weightTicket.FullWeight != nil && weightTicket.EmptyWeight != nil {
				totalWeight += *weightTicket.FullWeight - *weightTicket.EmptyWeight
			}
		}
	}
	if totalWeight > 0 {
		// Reassign ppm shipment fields to their expected location on the mto shipment for dates, addresses, weights ...
		ppmToMtoShipment = ppmshipment.MapPPMShipmentFinalFields(ppmShipment, *ppmShipment.EstimatedWeight)
	} else {
		// Reassign ppm shipment fields to their expected location on the mto shipment for dates, addresses, weights ...
		ppmToMtoShipment = ppmshipment.MapPPMShipmentEstimatedFields(ppmShipment)
	}

	for _, serviceItem := range serviceItemsToPrice {
		pricer, err := ghcrateengine.PricerForServiceItem(serviceItem.ReService.Code)
		if err != nil {
			logger.Error("unable to find pricer for service item", zap.Error(err))
			return nil, err
		}

		// For the non-accessorial service items there isn't any initialization that is going to change between lookups
		// for the same param. However, this is how the payment request does things and we'd want to know if it breaks
		// rather than optimizing I think.
		serviceItemLookups := serviceparamvaluelookups.InitializeLookups(ppmToMtoShipment, serviceItem)

		// This is the struct that gets passed to every param lookup() method that was initialized above
		keyData := serviceparamvaluelookups.NewServiceItemParamKeyData(p.planner, serviceItemLookups, serviceItem, mtoShipment, contract.Code)

		// The distance value gets saved to the mto shipment model to reduce repeated api calls.
		var shipmentWithDistance models.MTOShipment
		err = appCtx.DB().Find(&shipmentWithDistance, mtoShipment.ID)
		if err != nil {
			logger.Error("could not find shipment in the database")
			return nil, err
		}
		serviceItem.MTOShipment = shipmentWithDistance
		// set this to avoid potential eTag errors because the MTOShipment.Distance field was likely updated
		ppmShipment.Shipment = shipmentWithDistance

		var paramValues models.PaymentServiceItemParams
		for _, param := range paramsForServiceCode(serviceItem.ReService.Code, paramsForServiceItems) {
			paramKey := param.ServiceItemParamKey
			// This is where the lookup() method of each service item param is actually evaluated
			paramValue, serviceParamErr := keyData.ServiceParamValue(appCtx, paramKey.Key) // Fails with "DistanceZip" param?
			if serviceParamErr != nil {
				logger.Error("could not calculate param value lookup", zap.Error(serviceParamErr))
				return nil, serviceParamErr
			}

			// Gather all the param values for the service item to pass to the pricer's Price() method
			paymentServiceItemParam := models.PaymentServiceItemParam{
				// Some pricers like Fuel Surcharge try to requery the shipment through the service item, this is a
				// workaround to avoid a not found error because our PPM shipment has no service items saved in the db.
				// I think the FSC service item should really be relying on one of the zip distance params.
				PaymentServiceItem: models.PaymentServiceItem{
					MTOServiceItem: serviceItem,
				},
				ServiceItemParamKey: paramKey,
				Value:               paramValue,
			}
			paramValues = append(paramValues, paymentServiceItemParam)
		}

		if len(paramValues) == 0 {
			return nil, fmt.Errorf("no params were found for service item %s", serviceItem.ReService.Code)
		}

		// Middle var here can give you info on payment params like FSC multiplier, price rate/factor, etc. if needed.
		centsValue, _, err := pricer.PriceUsingParams(appCtx, paramValues)
		if err != nil {
			return nil, err
		}

		totalPrice = totalPrice.AddCents(centsValue)

		switch serviceItem.ReService.Code {
		case "DPK":
			packPrice += centsValue
		case "DUPK":
			unpackPrice += centsValue
		case "DOP":
			originPrice += centsValue
		case "DDP":
			destinationPrice += centsValue
			// TODO: Others (Konstance?) can put cases here for FSC (fuel surcharge), DLH (domestic linehaul), etc. here.
		}
	}

	ppmCloseoutObj.ID = &ppmShipmentID
	ppmCloseoutObj.PlannedMoveDate = mtoShipment.ScheduledPickupDate
	ppmCloseoutObj.ActualMoveDate = mtoShipment.ActualPickupDate
	ppmCloseoutObj.Miles = (*int)(mtoShipment.Distance)
	ppmCloseoutObj.EstimatedWeight = ppmShipment.EstimatedWeight
	ppmCloseoutObj.ActualWeight = mtoShipment.PrimeActualWeight
	ppmCloseoutObj.ProGearWeightCustomer = ppmShipment.ProGearWeight
	ppmCloseoutObj.ProGearWeightSpouse = ppmShipment.SpouseProGearWeight
	ppmCloseoutObj.GrossIncentive = ppmShipment.FinalIncentive
	ppmCloseoutObj.GCC = nil
	ppmCloseoutObj.AOA = nil
	ppmCloseoutObj.RemainingReimbursementOwed = nil
	ppmCloseoutObj.HaulPrice = nil
	ppmCloseoutObj.HaulFSC = nil
	ppmCloseoutObj.DOP = &originPrice
	ppmCloseoutObj.DDP = &destinationPrice
	ppmCloseoutObj.PackPrice = &packPrice
	ppmCloseoutObj.UnpackPrice = &unpackPrice
	ppmCloseoutObj.SITReimbursement = &storageExpensePrice

	return &ppmCloseoutObj, nil
}

func paramsForServiceCode(code models.ReServiceCode, serviceParams models.ServiceParams) models.ServiceParams {
	var serviceItemParams models.ServiceParams
	for _, serviceParam := range serviceParams {
		if serviceParam.Service.Code == code {
			serviceItemParams = append(serviceItemParams, serviceParam)
		}
	}
	return serviceItemParams
}
