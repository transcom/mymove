package ppmshipment

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/copier"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

// ppmShipmentReviewDocuments implements the services.PPMShipmentReviewDocuments interface
type ppmShipmentReviewDocuments struct {
	services.PPMShipmentRouter
	services.SignedCertificationCreator
	services.SignedCertificationUpdater
	services.SSWPPMComputer
}

// NewPPMShipmentReviewDocuments creates a new ppmShipmentReviewDocuments
func NewPPMShipmentReviewDocuments(
	ppmShipmentRouter services.PPMShipmentRouter,
	signedCertificationCreator services.SignedCertificationCreator,
	signedCertificationUpdater services.SignedCertificationUpdater,
	sswPPMComputer services.SSWPPMComputer,
) services.PPMShipmentReviewDocuments {
	return &ppmShipmentReviewDocuments{
		PPMShipmentRouter:          ppmShipmentRouter,
		SignedCertificationCreator: signedCertificationCreator,
		SignedCertificationUpdater: signedCertificationUpdater,
		SSWPPMComputer:             sswPPMComputer,
	}
}

// SubmitReviewedDocuments saves a new customer signature for PPM documentation agreement and routes PPM shipment
func (p *ppmShipmentReviewDocuments) SubmitReviewedDocuments(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.PPMShipment, error) {
	if ppmShipmentID.IsNil() {
		return nil, apperror.NewBadDataError("PPM ID is required")
	}

	ppmShipment, err := FindPPMShipment(appCtx, ppmShipmentID)

	if err != nil {
		return nil, err
	}

	var updatedPPMShipment models.PPMShipment

	err = copier.CopyWithOption(&updatedPPMShipment, ppmShipment, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, err
	}

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		err = p.PPMShipmentRouter.SubmitReviewedDocuments(txnAppCtx, &updatedPPMShipment)

		if err != nil {
			return err
		}
		updatedPPMShipment.Shipment.MoveTaskOrder.SCCloseoutAssignedID = nil
		err = validatePPMShipment(appCtx, updatedPPMShipment, ppmShipment, &ppmShipment.Shipment, PPMShipmentUpdaterChecks...)

		if err != nil {
			return err
		}

		verrs, err := txnAppCtx.DB().ValidateAndUpdate(&updatedPPMShipment)

		if verrs.HasAny() {
			return apperror.NewInvalidInputError(ppmShipment.ID, err, verrs, "unable to validate PPMShipment")
		} else if err != nil {
			return apperror.NewQueryError("PPMShipment", err, "unable to update PPMShipment")
		}

		verrs, err = appCtx.DB().ValidateAndSave(&updatedPPMShipment.Shipment.MoveTaskOrder)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(updatedPPMShipment.ID, nil, verrs, "")
		}
		if err != nil {
			return err
		}

		err = p.signCertificationPPMCloseout(appCtx, updatedPPMShipment.Shipment.MoveTaskOrderID, updatedPPMShipment.ID)

		if err != nil {
			return err
		}

		// write the SSW calculated values out to the ppm_closeouts table
		ppmCloseoutSummary, err := p.convertSSWValuesToPPMCloseoutSummary(appCtx, updatedPPMShipment.ID)

		if err != nil {
			return err
		}

		verrs, err = appCtx.DB().ValidateAndCreate(ppmCloseoutSummary)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(ppmCloseoutSummary.ID, nil, verrs, "")
		}

		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return &updatedPPMShipment, nil
}

// Fetch SSW Data and populate the data into the PPM Closeout table
func (p *ppmShipmentReviewDocuments) convertSSWValuesToPPMCloseoutSummary(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.PPMCloseoutSummary, error) {
	var ppmCloseoutSummary models.PPMCloseoutSummary
	ssfd, err := p.SSWPPMComputer.FetchDataShipmentSummaryWorksheetFormData(appCtx, appCtx.Session(), ppmShipmentID)

	if err != nil {
		return nil, err
	}

	// we currently only need data from pages 1 and 2 and not page 3
	page1Data, page2Data, _, err := p.SSWPPMComputer.FormatValuesShipmentSummaryWorksheet(appCtx, *ssfd, true)

	if err != nil {
		return nil, err
	}

	// write the values to model then the ppm_closeouts table
	ppmCloseoutSummary.ID = uuid.Must(uuid.NewV4())
	ppmCloseoutSummary.PPMShipmentID = ppmShipmentID

	// values are in dollar format with $ need to convert to cents without $
	if page1Data.MaxObligationGCCMaxAdvance != "" && page1Data.MaxObligationGCCMaxAdvance != "Advance not available." {
		maxAdvance, err := priceToCents(page1Data.MaxObligationGCCMaxAdvance)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.MaxAdvance = (*unit.Cents)(&maxAdvance)
	}

	if page2Data.ContractedExpenseMemberPaid != "" {
		memberExpense, err := priceToCents(page2Data.ContractedExpenseMemberPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.MemberPaidContractedExpense = (*unit.Cents)(&memberExpense)
	}

	if page2Data.ContractedExpenseGTCCPaid != "" {
		gtccExpense, err := priceToCents(page2Data.ContractedExpenseGTCCPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.GTCCPaidContractedExpense = (*unit.Cents)(&gtccExpense)
	}

	if page2Data.PackingMaterialsMemberPaid != "" {
		memberPackingMaterials, err := priceToCents(page2Data.PackingMaterialsMemberPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.MemberPaidPackingMaterials = (*unit.Cents)(&memberPackingMaterials)
	}

	if page2Data.PackingMaterialsGTCCPaid != "" {
		gtccPackingMaterials, err := priceToCents(page2Data.PackingMaterialsGTCCPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.GTCCPaidPackingMaterials = (*unit.Cents)(&gtccPackingMaterials)
	}

	if page2Data.WeighingFeesMemberPaid != "" {
		memberWeighingFee, err := priceToCents(page2Data.WeighingFeesMemberPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.MemberPaidWeighingFee = (*unit.Cents)(&memberWeighingFee)
	}

	if page2Data.WeighingFeesGTCCPaid != "" {
		gtccWeighingFee, err := priceToCents(page2Data.WeighingFeesGTCCPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.GTCCPaidWeighingFee = (*unit.Cents)(&gtccWeighingFee)
	}

	if page2Data.RentalEquipmentMemberPaid != "" {
		memberRental, err := priceToCents(page2Data.RentalEquipmentMemberPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.MemberPaidRentalEquipment = (*unit.Cents)(&memberRental)
	}

	if page2Data.RentalEquipmentGTCCPaid != "" {
		gtccRental, err := priceToCents(page2Data.RentalEquipmentGTCCPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.GTCCPaidRentalEquipment = (*unit.Cents)(&gtccRental)
	}

	if page2Data.TollsMemberPaid != "" {
		memberTolls, err := priceToCents(page2Data.TollsMemberPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.MemberPaidTolls = (*unit.Cents)(&memberTolls)
	}

	if page2Data.TollsGTCCPaid != "" {
		gtccTolls, err := priceToCents(page2Data.TollsGTCCPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.GTCCPaidTolls = (*unit.Cents)(&gtccTolls)
	}

	if page2Data.OilMemberPaid != "" {
		memberOil, err := priceToCents(page2Data.OilMemberPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.MemberPaidOil = (*unit.Cents)(&memberOil)
	}

	if page2Data.OilGTCCPaid != "" {
		gtccOil, err := priceToCents(page2Data.OilGTCCPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.GTCCPaidOil = (*unit.Cents)(&gtccOil)
	}

	if page2Data.OtherMemberPaid != "" {
		memberOther, err := priceToCents(page2Data.OtherMemberPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.MemberPaidOther = (*unit.Cents)(&memberOther)
	}

	if page2Data.OtherGTCCPaid != "" {
		gtccOther, err := priceToCents(page2Data.OtherGTCCPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.GTCCPaidOther = (*unit.Cents)(&gtccOther)
	}

	if page2Data.TotalMemberPaid != "" {
		totalMember, err := priceToCents(page2Data.TotalMemberPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.TotalMemberPaidExpenses = (*unit.Cents)(&totalMember)
	}

	if page2Data.TotalGTCCPaid != "" {
		totalGtcc, err := priceToCents(page2Data.TotalGTCCPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.TotalGTCCPaidExpenses = (*unit.Cents)(&totalGtcc)
	}

	if page2Data.TotalMemberPaidSIT != "" {
		memberSIT, err := priceToCents(page2Data.TotalMemberPaidSIT)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.MemberPaidSIT = (*unit.Cents)(&memberSIT)
	}

	if page2Data.TotalGTCCPaidSIT != "" {
		gtccSIT, err := priceToCents(page2Data.TotalGTCCPaidSIT)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.GTCCPaidSIT = (*unit.Cents)(&gtccSIT)
	}

	if page2Data.SmallPackageExpenseGTCCPaid != "" {
		gtccSmallPackage, err := priceToCents(page2Data.SmallPackageExpenseGTCCPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.GTCCPaidSmallPackage = (*unit.Cents)(&gtccSmallPackage)
	}

	if page2Data.SmallPackageExpenseMemberPaid != "" {
		memberSmallPackage, err := priceToCents(page2Data.SmallPackageExpenseMemberPaid)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.MemberPaidSmallPackage = (*unit.Cents)(&memberSmallPackage)
	}

	if page2Data.PPMRemainingEntitlement != "" {
		remainingIncentive, err := priceToCents(page2Data.PPMRemainingEntitlement)

		if err != nil {
			return nil, err
		}

		ppmCloseoutSummary.RemainingIncentive = (*unit.Cents)(&remainingIncentive)
	}

	if page2Data.Disbursement != "" {
		if page2Data.Disbursement != "N/A" {
			// GTCC and Member disbursement are displayed as one value on the SSW
			// example: "GTCC: $1,000.00\nMember: $300.00"
			// we need to parse each value out of the Disbursement string.
			disbursement := strings.Split(page2Data.Disbursement, "\n")
			gtccDisbursementStr := strings.Split(disbursement[0], "GTCC: ")
			memberDisbursementStr := strings.Split(disbursement[1], "Member: ")

			memberDisbursement, err := priceToCents(memberDisbursementStr[1])

			if err != nil {
				return nil, err
			}

			ppmCloseoutSummary.MemberDisbursement = (*unit.Cents)(&memberDisbursement)

			gtccDisbursement, err := priceToCents(gtccDisbursementStr[1])

			if err != nil {
				return nil, err
			}

			ppmCloseoutSummary.GTCCDisbursement = (*unit.Cents)(&gtccDisbursement)
		}
	}

	return &ppmCloseoutSummary, err
}

func getPriceParts(rawPrice string, expectedDecimalPlaces int) (int, int, error) {
	// Get rid of a dollar sign if there is one.
	basePrice := strings.Replace(rawPrice, "$", "", -1)
	basePrice = strings.Replace(basePrice, ",", "", -1)

	// Split the string on the decimal point.
	priceParts := strings.Split(basePrice, ".")
	if len(priceParts) != 2 {
		return 0, 0, fmt.Errorf("expected 2 price parts but found %d for price [%s]", len(priceParts), rawPrice)
	}

	integerPart, err := strconv.Atoi(priceParts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("could not convert integer part of price [%s]", rawPrice)
	}

	if len(priceParts[1]) != expectedDecimalPlaces {
		return 0, 0, fmt.Errorf("expected %d decimal places but found %d for price [%s]", expectedDecimalPlaces,
			len(priceParts[1]), rawPrice)
	}

	fractionalPart, err := strconv.Atoi(priceParts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("could not convert fractional part of price [%s]", rawPrice)
	}

	return integerPart, fractionalPart, nil
}

func priceToCents(rawPrice string) (int, error) {
	s := strings.TrimSpace(rawPrice)
	if !strings.Contains(s, "$") {
		return 0, nil
	}
	integerPart, fractionalPart, err := getPriceParts(rawPrice, 2)
	if err != nil {
		return 0, fmt.Errorf("could not parse price [%s]: %w", rawPrice, err)
	}

	cents := (integerPart * 100) + fractionalPart
	return cents, nil
}

func (p *ppmShipmentReviewDocuments) signCertificationPPMCloseout(appCtx appcontext.AppContext, moveID uuid.UUID, ppmShipmentID uuid.UUID) error {
	// Retrieve if PPM has certificate
	signedCertifications, err := models.FetchSignedCertificationPPMByType(appCtx.DB(), appCtx.Session(), moveID, ppmShipmentID, models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT)
	if err != nil {
		return err
	}

	signatureText := fmt.Sprintf("%s %s", appCtx.Session().FirstName, appCtx.Session().LastName)

	if len(signedCertifications) == 0 {
		// Add new certificate
		certType := models.SignedCertificationTypeCloseoutReviewedPPMPAYMENT
		now := time.Now()
		signedCertification := models.SignedCertification{
			SubmittingUserID:  appCtx.Session().UserID,
			MoveID:            moveID,
			PpmID:             models.UUIDPointer(ppmShipmentID),
			CertificationType: &certType,
			CertificationText: "Confirmed: Reviewed Closeout PPM PAYMENT ",
			Signature:         signatureText,
			Date:              now,
		}
		cert, err := p.SignedCertificationCreator.CreateSignedCertification(appCtx, signedCertification)
		fmt.Println(cert)
		if err != nil {
			return err
		}
	} else {
		// Update existing certificate. Note, reviews can occur N times.
		eTag := etag.GenerateEtag(signedCertifications[0].UpdatedAt)
		// Update with current counselor information
		signedCertifications[0].SubmittingUserID = appCtx.Session().UserID
		signedCertifications[0].Signature = signatureText
		_, err := p.SignedCertificationUpdater.UpdateSignedCertification(appCtx, *signedCertifications[0], eTag)
		if err != nil {
			return err
		}
	}

	return nil
}
