package transportationaccountingcode

import (
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
func (f transportationAccountingCodeFetcher) FetchOrderTransportationAccountingCodes(serviceMemberAffiliation models.ServiceMemberAffiliation, ordersIssueDate time.Time, tacCode string, appCtx appcontext.AppContext) ([]models.TransportationAccountingCode, error) {
	var tacs []models.TransportationAccountingCode
	var err error
	switch serviceMemberAffiliation {
	case models.AffiliationCOASTGUARD:
		// If a service member is in the Coast Guard don't filter out the household goods code of 'HS' because that is
		// primarily how their TGET records are coded along with 'HT' and 'HC' infrequently. If this changes in the future
		// then this can be revisited to weight the different LOAs similar to the other services.
		err = appCtx.DB().Q().
			Join("lines_of_accounting loa", "loa.loa_sys_id = transportation_accounting_codes.loa_sys_id").
			Where("transportation_accounting_codes.tac = ?", tacCode).
			Where("? between transportation_accounting_codes.trnsprtn_acnt_bgn_dt and transportation_accounting_codes.trnsprtn_acnt_end_dt", ordersIssueDate).
			Where("? between loa.loa_bgn_dt and loa.loa_end_dt", ordersIssueDate).
			Where("transportation_accounting_codes.tac_fn_bl_mod_cd != 'P'").
			All(&tacs)
	default:
		// For all other service members, filter out LineOfAccountingHouseholdGoodsCodeNTS "HS"
		err = appCtx.DB().Q().
			Join("lines_of_accounting loa", "loa.loa_sys_id = transportation_accounting_codes.loa_sys_id").
			Where("transportation_accounting_codes.tac = ?", tacCode).
			Where("? between transportation_accounting_codes.trnsprtn_acnt_bgn_dt and transportation_accounting_codes.trnsprtn_acnt_end_dt", ordersIssueDate).
			Where("? between loa.loa_bgn_dt and loa.loa_end_dt", ordersIssueDate).
			Where("transportation_accounting_codes.tac_fn_bl_mod_cd != 'P'").
			Where("loa.loa_hs_gds_cd != ?", models.LineOfAccountingHouseholdGoodsCodeNTS).
			All(&tacs)
	}
	if err != nil {
		return []models.TransportationAccountingCode{}, err
	}
	return tacs, nil
}
