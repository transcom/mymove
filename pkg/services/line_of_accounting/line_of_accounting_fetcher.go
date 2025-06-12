package lineofaccounting

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	edi "github.com/transcom/mymove/pkg/edi/segment"
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
func (f linesOfAccountingFetcher) FetchLongLinesOfAccounting(departmentIndicator models.DepartmentIndicator, ordersIssueDate time.Time, tacCode string, appCtx appcontext.AppContext) ([]models.LineOfAccounting, error) {
	// Fetch the TACs associated with this order and department indicator
	// Our lines of accounting will be present here
	tacs, err := f.FetchOrderTransportationAccountingCodes(departmentIndicator, ordersIssueDate, tacCode, appCtx)
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
	for currLoaIndex, loa := range linesOfAccounting {
		// if LOA Household Goods Program Code is null, invalid
		// This is NOT part of the DFAS elements, but is still required
		// for validity
		if loa.LoaHsGdsCd == nil {
			validHhgProgramCodeForLoa = false
		} else {
			validHhgProgramCodeForLoa = true
		}
		linesOfAccounting[currLoaIndex].ValidHhgProgramCodeForLoa = &validHhgProgramCodeForLoa

		missingRequiredDfasFields := validateDFASFields(loa)
		valid := len(missingRequiredDfasFields) == 0
		if !valid {
			var errMessage string
			for _, missingField := range missingRequiredDfasFields {
				errMessage += missingField + ", "
			}
			// If any LOA DFAS elements are missing, log it for informational purposes
			appCtx.Logger().Info("LOA with ID "+loa.ID.String()+" missing information: "+errMessage, zap.Error(err))
		}

		linesOfAccounting[currLoaIndex].ValidLoaForTac = &valid
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

// Helper function to run the LOA through the DFAS validator
// and return whether or not it is valid as well as which required
// DFAS fields it may be missing.
// NOTE: This only returns REQUIRED missing fields, not optional.
func validateDFASFields(loa models.LineOfAccounting) []string {
	dfasFields := newDfasValidator()

	// Only track the missing codes that ARE "REQUIRED"
	var missingRequiredCodes []string

	// Pull the DFAS fields from the LOA with
	// special DFAS getters
	for _, field := range dfasFields {
		val := field.Getter(&loa)
		// In the past we've had a field be empty strings
		// It shouldn't be possible anymore, but just in case
		if val == nil || strings.TrimSpace(*val) == "" {
			// Field is missing, check if it's a required field
			// to be considered invalid
			if field.Required {
				// Required field is missing!
				missingRequiredCodes = append(missingRequiredCodes, field.Code)
			}
		}
	}

	// If we have any required fields missing, this LOA is invalid
	return missingRequiredCodes
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

// dfasField Allows us to map a DFAS code in a LOA
// to its getter and if it is required or not
type dfasField struct {
	Code     string
	Required bool
	Getter   func(*models.LineOfAccounting) *string
}

// Helper function for accessing DFAS fields from the LOA, pairing them
// with a special "getter". This is needed because certain fields such as A3
// require special concat logic and comparison of 2 different fields
//
// "Required" is set by the MilMove government representatives.
// This is due to how it is decided to mark a field as valid or invalid.
// LOA data is a mess, and if it's too strict then everything is invalid,
// so this is our middle ground... for now. Subject to all kinds of change in
// the future. It is NOT the official "DFAS Valid/Invalid Guideline"
func newDfasValidator() []dfasField {
	// Return all known DFAS fields mapping their code
	// to their getter and whether or not it's required
	return []dfasField{
		{
			Code:     edi.FA2DetailCodeA1.String(),
			Required: true,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaDptID
			},
		},
		{
			Code:     edi.FA2DetailCodeA2.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaTnsfrDptNm
			},
		},
		{
			Code:     edi.FA2DetailCodeA3.String(),
			Required: true,
			Getter: func(loa *models.LineOfAccounting) *string {
				// A3 is a concatenation of the fiscal years
				if loa.LoaBgFyTx == nil || loa.LoaEndFyTx == nil {
					return nil
				}
				concatFy := fmt.Sprintf("%d%d", *loa.LoaBgFyTx, *loa.LoaEndFyTx)
				return &concatFy
			},
		},
		{
			Code:     edi.FA2DetailCodeA4.String(),
			Required: true,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaBafID
			},
		},
		{
			Code:     edi.FA2DetailCodeA5.String(),
			Required: true,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaTrsySfxTx
			},
		},
		{
			Code:     edi.FA2DetailCodeA6.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaMajClmNm
			},
		},
		{
			Code:     edi.FA2DetailCodeB1.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaOpAgncyID
			},
		},
		{
			Code:     edi.FA2DetailCodeB2.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaAlltSnID
			},
		},
		{
			Code:     edi.FA2DetailCodeB3.String(),
			Required: true,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaUic
			},
		},
		{
			Code:     edi.FA2DetailCodeC1.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaPgmElmntID
			},
		},
		{
			Code:     edi.FA2DetailCodeC2.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaTskBdgtSblnTx
			},
		},
		{
			Code:     edi.FA2DetailCodeD1.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaDfAgncyAlctnRcpntID
			},
		},
		{
			Code:     edi.FA2DetailCodeD4.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaJbOrdNm
			},
		},
		{
			Code:     edi.FA2DetailCodeD6.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaSbaltmtRcpntID
			},
		},
		{
			Code:     edi.FA2DetailCodeD7.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaWkCntrRcpntNm
			},
		},
		{
			Code:     edi.FA2DetailCodeE1.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaMajRmbsmtSrcID
			},
		},
		{
			Code:     edi.FA2DetailCodeE2.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaDtlRmbsmtSrcID
			},
		},
		{
			Code:     edi.FA2DetailCodeE3.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaCustNm
			},
		},
		{
			Code:     edi.FA2DetailCodeF1.String(),
			Required: true,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaObjClsID
			},
		},
		{
			Code:     edi.FA2DetailCodeF3.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaSrvSrcID
			},
		},
		{
			Code:     edi.FA2DetailCodeG2.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaSpclIntrID
			},
		},
		{
			Code:     edi.FA2DetailCodeI1.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaBdgtAcntClsNm
			},
		},
		{
			Code:     edi.FA2DetailCodeJ1.String(),
			Required: true,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaDocID
			},
		},
		{
			Code:     edi.FA2DetailCodeK6.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaClsRefID
			},
		},
		{
			Code:     edi.FA2DetailCodeL1.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaInstlAcntgActID
			},
		},
		{
			Code:     edi.FA2DetailCodeM1.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaLclInstlID
			},
		},
		{
			Code:     edi.FA2DetailCodeN1.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaTrnsnID
			},
		},
		{
			Code:     edi.FA2DetailCodeP5.String(),
			Required: false,
			Getter: func(loa *models.LineOfAccounting) *string {
				return loa.LoaFmsTrnsactnID
			},
		},
	}
}
