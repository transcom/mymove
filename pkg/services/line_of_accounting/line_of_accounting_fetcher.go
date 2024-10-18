package lineofaccounting

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type linesOfAccountingFetcher struct {
	services.TransportationAccountingCodeFetcher
}

// NewLinesOfAccountingFetcher creates a new linesOfAccountingFetcher service
func NewLinesOfAccountingFetcher(tacFetcher services.TransportationAccountingCodeFetcher) services.LineOfAccountingFetcher {
	return &linesOfAccountingFetcher{tacFetcher}
}

// This function returns all applicable long lines of accounting that can be sorted later according to business rules
func (f linesOfAccountingFetcher) FetchLongLinesOfAccounting(serviceMemberAffiliation models.ServiceMemberAffiliation, ordersIssueDate time.Time, tacCode string, appCtx appcontext.AppContext) ([]models.LineOfAccounting, error) {
	// Fetch the TACs associated with this order and service member affiliation
	// Our lines of accounting will be present here
	tacs, err := f.FetchOrderTransportationAccountingCodes(serviceMemberAffiliation, ordersIssueDate, tacCode, appCtx)
	if err != nil {
		return []models.LineOfAccounting{}, err
	}
	// Now that we have our TACs and LOAs, we need to sort accordingly
	linesOfAccounting := sortTransportationAccountingCodesAndLinesOfAccounting(tacs)

	linesOfAccounting, err = checkForValidHhgProgramCodeForLoaAndValidLoaForTac(linesOfAccounting, appCtx)
	if err != nil {
		return []models.LineOfAccounting{}, err
	}

	return linesOfAccounting, nil
}

// Sort the lines of accounting according to business logic. See confluence articles, for example here is one article (Of many)
// https://dp3.atlassian.net/wiki/spaces/MT/pages/2232680449/2023-08-14+TTV+EDI+Sync
// The result of the sorted lines of accounting will allow us to select the first item in the array
// as the first item post sorting is the correct line of counting to utilize for proper EDI858 invoicing
func sortTransportationAccountingCodesAndLinesOfAccounting(tacs []models.TransportationAccountingCode) []models.LineOfAccounting {

	// Sort TACs in memory
	sort.SliceStable(tacs, func(i, j int) bool {
		// Order by tac_fn_bl_mod_cd ascending
		if *tacs[i].TacFnBlModCd != *tacs[j].TacFnBlModCd {
			return *tacs[i].TacFnBlModCd < *tacs[j].TacFnBlModCd
		}
		// On a tie break, sort by loa_bgn_dt descending
		if tacs[i].LineOfAccounting.LoaBgnDt != tacs[j].LineOfAccounting.LoaBgnDt {
			return tacs[i].LineOfAccounting.LoaBgnDt.After(*tacs[j].LineOfAccounting.LoaBgnDt)
		}
		// On a tie break, sort by tac_fy_txt descending
		return *tacs[i].TacFyTxt > *tacs[j].TacFyTxt
	})

	// Extract the LOAs from the sorted TACs
	var linesOfAccounting []models.LineOfAccounting
	for _, tac := range tacs {
		linesOfAccounting = append(linesOfAccounting, *tac.LineOfAccounting)
	}

	return linesOfAccounting
}

func checkForValidHhgProgramCodeForLoaAndValidLoaForTac(linesOfAccounting []models.LineOfAccounting, appCtx appcontext.AppContext) ([]models.LineOfAccounting, error) {
	var err error
	var validHhgProgramCodeForLoa bool
	var validLoaForTac bool
	for currLoaIndex, loa := range linesOfAccounting {

		// if LOA Household Goods Program Code is null, invalid
		if loa.LoaHsGdsCd == nil {
			validHhgProgramCodeForLoa = false
		} else {
			validHhgProgramCodeForLoa = true
		}
		linesOfAccounting[currLoaIndex].ValidHhgProgramCodeForLoa = &validHhgProgramCodeForLoa

		// if any LOA DFAS elements are missing, invalid
		var missingLoaFields []string
		if loa.LoaSysID == nil {
			missingLoaFields = append(missingLoaFields, "loa.LoaSysID")
		}
		if loa.LoaDptID == nil { // Department Indicator (A1)
			missingLoaFields = append(missingLoaFields, "loa.LoaDptID")
		}
		if loa.LoaTnsfrDptNm == nil { // Transfer from Department (A2)
			missingLoaFields = append(missingLoaFields, "loa.LoaTnsfrDptNm")
		}
		if loa.LoaBgFyTx == nil || loa.LoaEndFyTx == nil { // Ending Fiscal Year Indicator (A3)
			// A3 is a concatenation of both LoaBgFyTx and LoaEndFyTx
			if loa.LoaBgFyTx == nil {
				missingLoaFields = append(missingLoaFields, "loa.LoaBgFyTx")
			}
			if loa.LoaEndFyTx == nil {
				missingLoaFields = append(missingLoaFields, "loa.LoaEndFyTx")
			}
		}
		if loa.LoaBafID == nil { // Basic Symbol Number (A4)
			missingLoaFields = append(missingLoaFields, "loa.LoaBafID")
		}
		if loa.LoaTrsySfxTx == nil { // Subhead/Limit (A5)
			missingLoaFields = append(missingLoaFields, "loa.LoaTrsySfxTx")
		}
		if loa.LoaMajClmNm == nil { // Fund Code/MC (A6)
			missingLoaFields = append(missingLoaFields, "loa.LoaMajClmNm")
		}
		if loa.LoaOpAgncyID == nil { // Operating Agency Code/Fund Admin (B1)
			missingLoaFields = append(missingLoaFields, "loa.LoaOpAgncyID")
		}
		if loa.LoaAlltSnID == nil { // Allotment Serial Number (B2)
			missingLoaFields = append(missingLoaFields, "loa.LoaAlltSnID")
		}
		if loa.LoaUic == nil { // Activity Address Code/UIC (B3)
			missingLoaFields = append(missingLoaFields, "loa.LoaUic")
		}
		if loa.LoaPgmElmntID == nil { // Program Element (C1)
			missingLoaFields = append(missingLoaFields, "loa.LoaPgmElmntID")
		}
		if loa.LoaTskBdgtSblnTx == nil { // Project Task or Budget Sub line (C2)
			missingLoaFields = append(missingLoaFields, "loa.LoaTskBdgtSblnTx")
		}
		if loa.LoaDfAgncyAlctnRcpntID == nil { // Defense Agency Allocation Recipient (D1)
			missingLoaFields = append(missingLoaFields, "loa.LoaDfAgncyAlctnRcpntID")
		}
		if loa.LoaJbOrdNm == nil { // Job Order/Work Order Code (D4)
			missingLoaFields = append(missingLoaFields, "loa.LoaJbOrdNm")
		}
		if loa.LoaSbaltmtRcpntID == nil { // Sub-allotment Recipient (D6)
			missingLoaFields = append(missingLoaFields, "loa.LoaSbaltmtRcpntID")
		}
		if loa.LoaWkCntrRcpntNm == nil { // Work Center Recipient (D7)
			missingLoaFields = append(missingLoaFields, "loa.LoaWkCntrRcpntNm")
		}
		if loa.LoaMajRmbsmtSrcID == nil { // Major Reimbursement Source Code (E1)
			missingLoaFields = append(missingLoaFields, "loa.LoaMajRmbsmtSrcID")
		}
		if loa.LoaDtlRmbsmtSrcID == nil { // Detail Reimbursement Source Code (E2)
			missingLoaFields = append(missingLoaFields, "loa.LoaDtlRmbsmtSrcID")
		}
		if loa.LoaCustNm == nil { // Customer Indicator/MPC (E3)
			missingLoaFields = append(missingLoaFields, "loa.LoaCustNm")
		}
		if loa.LoaObjClsID == nil { // Object Class (F1)
			missingLoaFields = append(missingLoaFields, "loa.LoaObjClsID")
		}
		if loa.LoaSrvSrcID == nil { // Government or Public Sector Identifier (F3)
			missingLoaFields = append(missingLoaFields, "loa.LoaSrvSrcID")
		}
		if loa.LoaSpclIntrID == nil { // Special Interest Code or Special Program Cost Code (G2)
			missingLoaFields = append(missingLoaFields, "loa.LoaSpclIntrID")
		}
		if loa.LoaBdgtAcntClsNm == nil { // Abbreviated Department of Defense (DoD) Budget and Accounting (I1)
			missingLoaFields = append(missingLoaFields, "loa.LoaBdgtAcntClsNm")
		}
		if loa.LoaDocID == nil { // (SDN) Document or Record Reference Number (J1)
			missingLoaFields = append(missingLoaFields, "loa.LoaDocID")
		}
		if loa.LoaClsRefID == nil { // (ACRN) Accounting Classification Reference Code (K6)
			missingLoaFields = append(missingLoaFields, "loa.LoaClsRefID")
		}
		if loa.LoaInstlAcntgActID == nil { // Accounting Installation Number (L1)
			missingLoaFields = append(missingLoaFields, "loa.LoaInstlAcntgActID")
		}
		if loa.LoaLclInstlID == nil { // Local Installation Data (M1)
			missingLoaFields = append(missingLoaFields, "loa.LoaLclInstlID")
		}
		if loa.LoaTrnsnID == nil { // Transaction ID (N1)
			missingLoaFields = append(missingLoaFields, "loa.LoaTrnsnID")
		}
		if loa.LoaFmsTrnsactnID == nil { // Transaction Type (P5)
			missingLoaFields = append(missingLoaFields, "loa.LoaFmsTrnsactnID")
		}

		if missingLoaFields != nil {
			validLoaForTac = false

			var errMessage string
			if len(missingLoaFields) == 1 {
				errMessage += missingLoaFields[0]
			} else {
				for i := range missingLoaFields {
					errMessage += missingLoaFields[i] + ", "
				}
			}
			// If any LOA DFAS elements are missing, log it for informational purposes
			appCtx.Logger().Info("LOA with ID "+loa.ID.String()+" missing information: "+errMessage, zap.Error(err))
		} else {
			validLoaForTac = true
		}

		linesOfAccounting[currLoaIndex].ValidLoaForTac = &validLoaForTac
	}

	transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		for currLoaIndex := range linesOfAccounting {
			// update line of accounting validHHGProgramCodeForLOA and ValidLoaForTac fields
			verrs, err := txnCtx.DB().ValidateAndUpdate(&linesOfAccounting[currLoaIndex])
			if verrs != nil && verrs.HasAny() {
				return apperror.NewInvalidInputError(linesOfAccounting[currLoaIndex].ID, err, verrs, "invalid input found while updating ValidLoaForTac or ValidHhgProgramCodeForLoa for LOA")
			} else if err != nil {
				return apperror.NewQueryError("LOA", err, "")
			}
		}
		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return linesOfAccounting, err
}

func (f linesOfAccountingFetcher) BuildFullLineOfAccountingString(loa models.LineOfAccounting) string {
	emptyString := ""
	var loaFyTx string
	if fmt.Sprint(*loa.LoaBgFyTx) != "" && fmt.Sprint(*loa.LoaEndFyTx) != "" {
		loaFyTx = fmt.Sprint(*loa.LoaBgFyTx) + fmt.Sprint(*loa.LoaEndFyTx)
	} else {
		loaFyTx = ""
	}

	if loa.LoaDptID == nil {
		loa.LoaDptID = &emptyString
	}
	if loa.LoaTnsfrDptNm == nil {
		loa.LoaTnsfrDptNm = &emptyString
	}
	if loa.LoaBafID == nil {
		loa.LoaBafID = &emptyString
	}
	if loa.LoaTrsySfxTx == nil {
		loa.LoaTrsySfxTx = &emptyString
	}
	if loa.LoaMajClmNm == nil {
		loa.LoaMajClmNm = &emptyString
	}
	if loa.LoaOpAgncyID == nil {
		loa.LoaOpAgncyID = &emptyString
	}
	if loa.LoaAlltSnID == nil {
		loa.LoaAlltSnID = &emptyString
	}
	if loa.LoaUic == nil {
		loa.LoaUic = &emptyString
	}
	if loa.LoaPgmElmntID == nil {
		loa.LoaPgmElmntID = &emptyString
	}
	if loa.LoaTskBdgtSblnTx == nil {
		loa.LoaTskBdgtSblnTx = &emptyString
	}
	if loa.LoaDfAgncyAlctnRcpntID == nil {
		loa.LoaDfAgncyAlctnRcpntID = &emptyString
	}
	if loa.LoaJbOrdNm == nil {
		loa.LoaJbOrdNm = &emptyString
	}
	if loa.LoaSbaltmtRcpntID == nil {
		loa.LoaSbaltmtRcpntID = &emptyString
	}
	if loa.LoaWkCntrRcpntNm == nil {
		loa.LoaWkCntrRcpntNm = &emptyString
	}
	if loa.LoaMajRmbsmtSrcID == nil {
		loa.LoaMajRmbsmtSrcID = &emptyString
	}
	if loa.LoaDtlRmbsmtSrcID == nil {
		loa.LoaDtlRmbsmtSrcID = &emptyString
	}
	if loa.LoaCustNm == nil {
		loa.LoaCustNm = &emptyString
	}
	if loa.LoaObjClsID == nil {
		loa.LoaObjClsID = &emptyString
	}
	if loa.LoaSrvSrcID == nil {
		loa.LoaSrvSrcID = &emptyString
	}
	if loa.LoaSpclIntrID == nil {
		loa.LoaSpclIntrID = &emptyString
	}
	if loa.LoaBdgtAcntClsNm == nil {
		loa.LoaBdgtAcntClsNm = &emptyString
	}
	if loa.LoaDocID == nil {
		loa.LoaDocID = &emptyString
	}
	if loa.LoaClsRefID == nil {
		loa.LoaClsRefID = &emptyString
	}
	if loa.LoaInstlAcntgActID == nil {
		loa.LoaInstlAcntgActID = &emptyString
	}
	if loa.LoaLclInstlID == nil {
		loa.LoaLclInstlID = &emptyString
	}
	if loa.LoaTrnsnID == nil {
		loa.LoaTrnsnID = &emptyString
	}
	if loa.LoaFmsTrnsactnID == nil {
		loa.LoaFmsTrnsactnID = &emptyString
	}

	// commented names are Navy PPTAS equivalents
	LineOfAccountingDfasElementOrder := []string{
		*loa.LoaDptID,               // "LoaDptID"
		*loa.LoaTnsfrDptNm,          // "LoaTnsfrDptNm",
		loaFyTx,                     // "LoaEndFyTx",
		*loa.LoaBafID,               // "LoaBafID",
		*loa.LoaTrsySfxTx,           // "LoaTrsySfxTx",
		*loa.LoaMajClmNm,            // "LoaMajClmNm",
		*loa.LoaOpAgncyID,           // "LoaOpAgncyID",
		*loa.LoaAlltSnID,            // "LoaAlltSnID",
		*loa.LoaUic,                 // "LoaUic",
		*loa.LoaPgmElmntID,          // "LoaPgmElmntID",
		*loa.LoaTskBdgtSblnTx,       // "LoaTskBdgtSblnTx",
		*loa.LoaDfAgncyAlctnRcpntID, // "LoaDfAgncyAlctnRcpntID",
		*loa.LoaJbOrdNm,             // "LoaJbOrdNm",
		*loa.LoaSbaltmtRcpntID,      // "LoaSbaltmtRcpntID",
		*loa.LoaWkCntrRcpntNm,       // "LoaWkCntrRcpntNm",
		*loa.LoaMajRmbsmtSrcID,      // "LoaMajRmbsmtSrcID",
		*loa.LoaDtlRmbsmtSrcID,      // "LoaDtlRmbsmtSrcID",
		*loa.LoaCustNm,              // "LoaCustNm",
		*loa.LoaObjClsID,            // "LoaObjClsID",
		*loa.LoaSrvSrcID,            // "LoaSrvSrcID",
		*loa.LoaSpclIntrID,          // "LoaSpcLIntrID",
		*loa.LoaBdgtAcntClsNm,       // "LoaBdgtAcntCLsNm",
		*loa.LoaDocID,               // "LoaDocID",
		*loa.LoaClsRefID,            // "LoaCLsRefID",
		*loa.LoaInstlAcntgActID,     // "LoaInstLAcntgActID",
		*loa.LoaLclInstlID,          // "LoaLcLInstLID",
		*loa.LoaTrnsnID,             // "LoaTrnsnID",
		*loa.LoaFmsTrnsactnID,       // "LoaFmsTrnsactnID",
	}

	longLoa := strings.Join(LineOfAccountingDfasElementOrder, "*")
	longLoa = strings.ReplaceAll(longLoa, " *", "*")

	return longLoa
}
