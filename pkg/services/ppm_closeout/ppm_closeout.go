package ppmcloseout

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
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
	ppmShipment, err := p.GetPPMShipment(appCtx, ppmShipmentID)
	if err != nil {
		return nil, err
	}

	actualWeight := p.GetActualWeight(*ppmShipment)
	proGearWeightCustomer, proGearWeightSpouse := p.GetProGearWeights(*ppmShipment)

	ppmCloseoutObj.ID = &ppmShipmentID
	ppmCloseoutObj.PlannedMoveDate = &ppmShipment.ExpectedDepartureDate
	ppmCloseoutObj.ActualMoveDate = ppmShipment.ActualMoveDate
	ppmCloseoutObj.Miles = (*int)(ppmShipment.Shipment.Distance)
	ppmCloseoutObj.EstimatedWeight = ppmShipment.EstimatedWeight
	ppmCloseoutObj.ActualWeight = &actualWeight
	ppmCloseoutObj.ProGearWeightCustomer = &proGearWeightCustomer
	ppmCloseoutObj.ProGearWeightSpouse = &proGearWeightSpouse
	ppmCloseoutObj.GrossIncentive = ppmShipment.FinalIncentive
	ppmCloseoutObj.GCC = nil
	ppmCloseoutObj.AOA = ppmShipment.AdvanceAmountReceived
	ppmCloseoutObj.RemainingIncentive = nil
	ppmCloseoutObj.HaulPrice = nil
	ppmCloseoutObj.HaulFSC = nil
	ppmCloseoutObj.DOP = nil
	ppmCloseoutObj.DDP = nil
	ppmCloseoutObj.PackPrice = nil
	ppmCloseoutObj.UnpackPrice = nil
	ppmCloseoutObj.SITReimbursement = nil

	return &ppmCloseoutObj, nil
}

/*
* returns
* customer pro gear weight, spouse pro gear weight
 */
func (p *ppmCloseoutFetcher) GetProGearWeights(ppmShipment models.PPMShipment) (unit.Pound, unit.Pound) {
	var totalWeightCustomer unit.Pound
	var totalWeightSpouse unit.Pound
	if len(ppmShipment.ProgearWeightTickets) >= 1 {
		for _, weightTicket := range ppmShipment.ProgearWeightTickets {
			if weightTicket.Status != nil && *weightTicket.BelongsToSelf && weightTicket.Weight != nil && *weightTicket.Status != models.PPMDocumentStatusRejected {
				totalWeightCustomer += *weightTicket.Weight
			} else if !*weightTicket.BelongsToSelf {
				totalWeightSpouse += *weightTicket.Weight
			}
		}
	}
	return totalWeightCustomer, totalWeightSpouse
}

func (p *ppmCloseoutFetcher) GetPPMShipment(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.PPMShipment, error) {
	var ppmShipment models.PPMShipment
	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"ID",
			"Shipment",
			"ExpectedDepartureDate",
			"ActualMoveDate",
			"EstimatedWeight",
			"WeightTickets",
			"ProgearWeightTickets",
			"FinalIncentive",
			"AdvanceAmountReceived",
			"Shipment.Distance",
		).
		Find(&ppmShipment, ppmShipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(ppmShipmentID, "while looking for PPMShipment")
		default:
			return nil, apperror.NewQueryError("PPMShipment", err, "unable to find PPMShipment")
		}
	}

	// Check if PPM shipment is in "NEEDS_PAYMENT_APPROVAL" status, if not, it's not ready for closeout
	if ppmShipment.Status != models.PPMShipmentStatusNeedsPaymentApproval {
		return nil, apperror.NewPPMNotReadyForCloseoutError(ppmShipmentID, "")
	}

	return &ppmShipment, err
}

func (p *ppmCloseoutFetcher) GetActualWeight(ppmShipment models.PPMShipment) unit.Pound {
	var totalWeight unit.Pound
	if len(ppmShipment.WeightTickets) >= 1 {
		for _, weightTicket := range ppmShipment.WeightTickets {
			if weightTicket.Status != nil && weightTicket.FullWeight != nil && weightTicket.EmptyWeight != nil && *weightTicket.Status != models.PPMDocumentStatusRejected {
				totalWeight += *weightTicket.FullWeight - *weightTicket.EmptyWeight
			}
		}
	}
	return totalWeight
}

func (p *ppmCloseoutFetcher) GetExpenseStoragePrice(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (unit.Cents, error) {
	var expenseItems []models.MovingExpense
	var storageExpensePrice unit.Cents
	err := appCtx.DB().Where("ppm_shipment_id = ?", ppmShipmentID).All(&expenseItems)
	if err != nil {
		return unit.Cents(0), err
	}

	for _, movingExpense := range expenseItems {
		if movingExpense.MovingExpenseType != nil && *movingExpense.MovingExpenseType == models.MovingExpenseReceiptTypeStorage {
			storageExpensePrice += *movingExpense.Amount
		}
	}
	return storageExpensePrice, err
}

func (p *ppmCloseoutFetcher) GetEntitlement(appCtx appcontext.AppContext, moveID uuid.UUID) (*models.Entitlement, error) {
	var moveModel models.Move
	err := appCtx.DB().EagerPreload(
		"OrdersID",
	).Find(&moveModel, moveID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(moveID, "while looking for Move")
		default:
			return nil, apperror.NewQueryError("Move", err, "unable to find Move")
		}
	}

	var order models.Order
	orderID := &moveModel.OrdersID
	errOrder := appCtx.DB().EagerPreload(
		"EntitlementID",
	).Find(&order, orderID)

	if errOrder != nil {
		switch errOrder {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(*orderID, "while looking for Order")
		default:
			return nil, apperror.NewQueryError("Order", errOrder, "unable to find Order")
		}
	}

	var entitlement models.Entitlement
	entitlementID := order.EntitlementID
	errEntitlement := appCtx.DB().EagerPreload(
		"DBAuthorizedWeight",
	).Find(&entitlement, entitlementID)

	if errEntitlement != nil {
		switch errEntitlement {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(*entitlementID, "while looking for Entitlement")
		default:
			return nil, apperror.NewQueryError("Entitlement", errEntitlement, "unable to find Entitlement")
		}
	}
	return &entitlement, nil
}
