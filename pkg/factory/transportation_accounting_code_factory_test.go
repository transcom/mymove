package factory

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
)

var defaultTac = models.TransportationAccountingCode{
	TAC: "E01A",
}

func (suite *FactorySuite) TestBuildTransportationAccountingCode() {

	suite.Run("Successful creation of default TAC", func() {
		// Under test:      BuildTransportationAccountingCode
		// Set up:          Create a TAC with no customizations or traits
		// Expected outcome:TAC should be created with default values
		tac := BuildTransportationAccountingCode(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultTac.TAC, tac.TAC)

		suite.Equal(*tac.LoaID, tac.LineOfAccounting.ID)
		suite.Equal((*int)(nil), tac.LineOfAccounting.LoaSysID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaDptID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaTnsfrDptNm)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaBafID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaTrsySfxTx)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaMajClmNm)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaOpAgncyID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaAlltSnID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaPgmElmntID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaTskBdgtSblnTx)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaDfAgncyAlctnRcpntID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaJbOrdNm)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaSbaltmtRcpntID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaWkCntrRcpntNm)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaMajRmbsmtSrcID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaDtlRmbsmtSrcID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaCustNm)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaObjClsID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaSrvSrcID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaSpclIntrID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaBdgtAcntClsNm)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaDocID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaClsRefID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaInstlAcntgActID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaLclInstlID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaFmsTrnsactnID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaDscTx)
		suite.Equal((*time.Time)(nil), tac.LineOfAccounting.LoaBgnDt)
		suite.Equal((*time.Time)(nil), tac.LineOfAccounting.LoaEndDt)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaFnctPrsNm)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaStatCd)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaHistStatCd)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaHsGdsCd)
		suite.Equal((*string)(nil), tac.LineOfAccounting.OrgGrpDfasCd)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaUic)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaTrnsnID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaSubAcntID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaBetCd)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaFndTyFgCd)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaBgtLnItmID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaScrtyCoopImplAgncCd)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaScrtyCoopDsgntrCd)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaScrtyCoopLnItmID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaAgncDsbrCd)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaAgncAcntngCd)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaFndCntrID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaCstCntrID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaPrjID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaActvtyID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaCstCd)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaWrkOrdID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaFnclArID)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaScrtyCoopCustCd)
		suite.Equal((*int)(nil), tac.LineOfAccounting.LoaEndFyTx)
		suite.Equal((*int)(nil), tac.LineOfAccounting.LoaBgFyTx)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaBgtRstrCd)
		suite.Equal((*string)(nil), tac.LineOfAccounting.LoaBgtSubActCd)

		suite.Equal((*int)(nil), tac.TacSysID)
		suite.Equal((*int)(nil), tac.LoaSysID)
		suite.Equal((*int)(nil), tac.TacFyTxt)
		suite.Equal((*string)(nil), tac.TacFnBlModCd)
		suite.Equal((*string)(nil), tac.OrgGrpDfasCd)
		suite.Equal((*string)(nil), tac.TacMvtDsgID)
		suite.Equal((*string)(nil), tac.TacTyCd)
		suite.Equal((*string)(nil), tac.TacUseCd)
		suite.Equal((*string)(nil), tac.TacMajClmtID)
		suite.Equal((*string)(nil), tac.TacBillActTxt)
		suite.Equal((*string)(nil), tac.TacCostCtrNm)
		suite.Equal((*string)(nil), tac.Buic)
		suite.Equal((*string)(nil), tac.TacHistCd)
		suite.Equal((*string)(nil), tac.TacStatCd)
		suite.Equal((*string)(nil), tac.TrnsprtnAcntTx)
		suite.Equal((*time.Time)(nil), tac.TrnsprtnAcntBgnDt)
		suite.Equal((*time.Time)(nil), tac.TrnsprtnAcntEndDt)
		suite.Equal((*string)(nil), tac.DdActvtyAdrsID)
		suite.Equal((*string)(nil), tac.TacBlldAddFrstLnTx)
		suite.Equal((*string)(nil), tac.TacBlldAddScndLnTx)
		suite.Equal((*string)(nil), tac.TacBlldAddThrdLnTx)
		suite.Equal((*string)(nil), tac.TacBlldAddFrthLnTx)
		suite.Equal((*string)(nil), tac.TacFnctPocNm)
	})

	suite.Run("Successful creation of a TAC with customization", func() {
		// Under test:      BuildTransportationAccountingCode
		// Set up:          Create a TAC with no customizations or traits
		// Expected outcome:TAC should be created with custom value
		tac := BuildTransportationAccountingCode(suite.DB(), []Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC: "1234",
				},
			},
			{
				Model: models.LineOfAccounting{
					LoaSysID: models.IntPointer(4321),
				},
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal("1234", tac.TAC)
		suite.Equal(4321, *tac.LineOfAccounting.LoaSysID)
	})
}
