package transportationaccountingcode

import (
	"database/sql"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type transportationAccountingCodeFetcher struct {
}

// NewTransportationAccountingCodeFetcher creates a new transportationAccountingCodeFetcher service
func NewTransportationAccountingCodeFetcher() services.TransportationAccountingCodeFetcher {
	return &transportationAccountingCodeFetcher{}
}

// FetchOrderTransportationAccountingCodes returns all applicable transportation accounting codes
// alongside their associated lines of accounting for an order
func (f transportationAccountingCodeFetcher) FetchOrderTransportationAccountingCodes(departmentIndicator models.DepartmentIndicator, ordersIssueDate time.Time, tacCode string, appCtx appcontext.AppContext) ([]models.TransportationAccountingCode, error) {
	var tacs []models.TransportationAccountingCode
	var err error

	// If a service member is in the Coast Guard don't filter out the household goods code of 'HS' because that is
	// primarily how their TGET records are coded along with 'HT' and 'HC' infrequently. If this changes in the future
	// then this can be revisited to weight the different LOAs similar to the other services.
	query := appCtx.DB().Q().
		EagerPreload("LineOfAccounting").
		Join("lines_of_accounting loa", "loa.loa_sys_id = transportation_accounting_codes.loa_sys_id").
		Where("transportation_accounting_codes.tac = ?", tacCode).
		Where("? BETWEEN transportation_accounting_codes.trnsprtn_acnt_bgn_dt AND transportation_accounting_codes.trnsprtn_acnt_end_dt", ordersIssueDate).
		Where("? BETWEEN loa.loa_bgn_dt AND loa.loa_end_dt", ordersIssueDate).
		Where("transportation_accounting_codes.tac_fn_bl_mod_cd != 'P'")

	// For all other service members, filter out LineOfAccountingHouseholdGoodsCodeNTS "HS"
	if departmentIndicator != models.DepartmentIndicatorCOASTGUARD {
		query = query.Where("loa.loa_hs_gds_cd != ?", models.LineOfAccountingHouseholdGoodsCodeNTS)
	}
	err = query.All(&tacs)
	if err != nil {
		return []models.TransportationAccountingCode{}, err
	}
	// Grab the associated LOAs
	for memoryIterationOfTac := range tacs {
		var loa models.LineOfAccounting
		// Find the LOA for this TAC's loa_sys_id
		if tacs[memoryIterationOfTac].LoaSysID != nil {
			err = appCtx.DB().Where("loa_sys_id = ?", tacs[memoryIterationOfTac].LoaSysID).First(&loa)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					continue
				default:
					return []models.TransportationAccountingCode{}, err
				}
			}
			// If this TAC is getting a LOA from POP, that's because it's from manually imported TGET data prior to September 2023
			// So, we need to have a conditional memory assignment here
			if tacs[memoryIterationOfTac].LineOfAccounting == nil {
				// This is new TGET data with no key references in the databases (Future overhaul should switch to using LoaSysId instead of ID)
				// Since this is new, we need to assign the memory for a loa object before linking the LOA to the TAC
				tacs[memoryIterationOfTac].LineOfAccounting = &loa
			} else {
				// Else, LineOfAccounting is not nil, meaning a LOA is being returned with this TAC.
				// We want to override this with the LoaSysId pulled LOA (The loas we manually looked up)
				*tacs[memoryIterationOfTac].LineOfAccounting = loa
			}
		}
	}
	return tacs, nil
}
