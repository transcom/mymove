package lineofaccounting

import (
	"sort"
	"time"

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

	var validHHGProgramCodeForLOA bool
	for currLoaIndex, loa := range linesOfAccounting {
		// if LOA Household Goods Program Code isn't null, it's valid
		if loa.LoaHsGdsCd != nil {
			validHHGProgramCodeForLOA = true
		} else {
			validHHGProgramCodeForLOA = false
		}
		linesOfAccounting[currLoaIndex].ValidHhgProgramCodeForLoa = &validHHGProgramCodeForLOA
	}

	transactionError := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		for currLoaIndex := range linesOfAccounting {
			// update line of accounting validHHGProgramCodeForLOA field in the database
			verrs, err := txnCtx.DB().ValidateAndUpdate(&linesOfAccounting[currLoaIndex])
			if verrs != nil && verrs.HasAny() {
				return apperror.NewInvalidInputError(linesOfAccounting[currLoaIndex].ID, err, verrs, "invalid input found while updating final destination address of service item")
			} else if err != nil {
				return apperror.NewQueryError("Service item", err, "")
			}
		}
		return nil
	})

	if transactionError != nil {
		return nil, transactionError
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
