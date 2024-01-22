package ppmcloseout

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

type ppmCloseoutFetcher struct{}

func NewPPMCloseoutFetcher() services.PPMCloseoutFetcher {
	return &ppmCloseoutFetcher{}
}

func (p *ppmCloseoutFetcher) GetPPMCloseout(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.PPMCloseout, error) {
	var ppmCloseoutObj models.PPMCloseout
	var ppmShipment models.PPMShipment
	var mtoShipment models.MTOShipment

	errPPM := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
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
			"AdvanceAmountReceived",
		).
		Find(&ppmShipment, ppmShipmentID)

	if errPPM != nil {
		switch errPPM {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(ppmShipmentID, "while looking for PPMShipment")
		default:
			return nil, apperror.NewQueryError("PPMShipment", errPPM, "unable to find PPMShipment")
		}
	}

	mtoShipmentID := &ppmShipment.ShipmentID
	errMTO := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"ID",
			"ScheduledPickupDate",
			"ActualPickupDate",
			"Distance",
			"PrimeActualWeight",
		).
		Find(&mtoShipment, mtoShipmentID)

	if errMTO != nil {
		switch errMTO {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(*mtoShipmentID, "while looking for MTOShipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", errMTO, "unable to find MTOShipment")
		}
	}

	// get all mtoServiceItems IDs that share a mtoShipmentID
	var mtoServiceItems models.MTOServiceItems
	errTest4 := appCtx.DB().Eager("ID").Where("mto_service_items.mto_shipment_id = ?", &mtoShipmentID).All(&mtoServiceItems)

	if errTest4 != nil {
		return nil, errTest4
	}

	var paymentServiceItem models.PaymentServiceItem
	var remainingIncentive unit.Cents
	var test = &mtoServiceItems

	for _, element := range *test {
		errTest5 := appCtx.DB().Eager("PriceCents", "staus", "paid_at").Where("payment_service_items.mto_service_item_id = ?", element.ID).All(&paymentServiceItem)

		if paymentServiceItem.Status == models.PaymentServiceItemStatusApproved && paymentServiceItem.PaidAt == nil {
			remainingIncentive = remainingIncentive.AddCents(*paymentServiceItem.PriceCents)
		}

		if errTest5 != nil {
			return nil, errTest5
		}
	}

	// paymentServiceItemID := uuid.FromStringOrNil("4730fee7-663d-4b09-9d09-0dab2c22f5d8")
	// errTest := appCtx.DB().Find(&paymentServiceItem, paymentServiceItemID)

	// if errTest != nil {
	// 	switch errTest {
	// 	case sql.ErrNoRows:
	// 		return nil, apperror.NewNotFoundError(paymentServiceItemID, "while looking for paymentServiceItem")
	// 	default:
	// 		return nil, apperror.NewQueryError("paymentServiceItem", errTest, "unable to find paymentServiceItem")
	// 	}
	// }

	//services.DomesticLinehaulPricer

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
	ppmCloseoutObj.AOA = ppmShipment.AdvanceAmountReceived
	ppmCloseoutObj.RemainingReimbursementOwed = &remainingIncentive //paymentServiceItem.PriceCents
	ppmCloseoutObj.HaulPrice = nil
	ppmCloseoutObj.HaulFSC = nil
	ppmCloseoutObj.DOP = nil
	ppmCloseoutObj.DDP = nil
	ppmCloseoutObj.PackPrice = nil
	ppmCloseoutObj.UnpackPrice = nil
	ppmCloseoutObj.SITReimbursement = nil

	return &ppmCloseoutObj, nil
}
