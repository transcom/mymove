package factory

import (
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
		suite.Nil(tac.LineOfAccounting.LoaSysID)
		suite.Nil(tac.LineOfAccounting.LoaDptID)
		suite.Nil(tac.LineOfAccounting.LoaTnsfrDptNm)
		suite.Nil(tac.LineOfAccounting.LoaBafID)
		suite.Nil(tac.LineOfAccounting.LoaTrsySfxTx)
		suite.Nil(tac.LineOfAccounting.LoaMajClmNm)
		suite.Nil(tac.LineOfAccounting.LoaOpAgncyID)
		suite.Nil(tac.LineOfAccounting.LoaAlltSnID)
		suite.Nil(tac.LineOfAccounting.LoaPgmElmntID)
		suite.Nil(tac.LineOfAccounting.LoaTskBdgtSblnTx)
		suite.Nil(tac.LineOfAccounting.LoaDfAgncyAlctnRcpntID)
		suite.Nil(tac.LineOfAccounting.LoaJbOrdNm)
		suite.Nil(tac.LineOfAccounting.LoaSbaltmtRcpntID)
		suite.Nil(tac.LineOfAccounting.LoaWkCntrRcpntNm)
		suite.Nil(tac.LineOfAccounting.LoaMajRmbsmtSrcID)
		suite.Nil(tac.LineOfAccounting.LoaDtlRmbsmtSrcID)
		suite.Nil(tac.LineOfAccounting.LoaCustNm)
		suite.Nil(tac.LineOfAccounting.LoaObjClsID)
		suite.Nil(tac.LineOfAccounting.LoaSrvSrcID)
		suite.Nil(tac.LineOfAccounting.LoaSpclIntrID)
		suite.Nil(tac.LineOfAccounting.LoaBdgtAcntClsNm)
		suite.Nil(tac.LineOfAccounting.LoaDocID)
		suite.Nil(tac.LineOfAccounting.LoaClsRefID)
		suite.Nil(tac.LineOfAccounting.LoaInstlAcntgActID)
		suite.Nil(tac.LineOfAccounting.LoaLclInstlID)
		suite.Nil(tac.LineOfAccounting.LoaFmsTrnsactnID)
		suite.Nil(tac.LineOfAccounting.LoaDscTx)
		suite.Nil(tac.LineOfAccounting.LoaBgnDt)
		suite.Nil(tac.LineOfAccounting.LoaEndDt)
		suite.Nil(tac.LineOfAccounting.LoaFnctPrsNm)
		suite.Nil(tac.LineOfAccounting.LoaStatCd)
		suite.Nil(tac.LineOfAccounting.LoaHistStatCd)
		suite.Nil(tac.LineOfAccounting.LoaHsGdsCd)
		suite.Nil(tac.LineOfAccounting.OrgGrpDfasCd)
		suite.Nil(tac.LineOfAccounting.LoaUic)
		suite.Nil(tac.LineOfAccounting.LoaTrnsnID)
		suite.Nil(tac.LineOfAccounting.LoaSubAcntID)
		suite.Nil(tac.LineOfAccounting.LoaBetCd)
		suite.Nil(tac.LineOfAccounting.LoaFndTyFgCd)
		suite.Nil(tac.LineOfAccounting.LoaBgtLnItmID)
		suite.Nil(tac.LineOfAccounting.LoaScrtyCoopImplAgncCd)
		suite.Nil(tac.LineOfAccounting.LoaScrtyCoopDsgntrCd)
		suite.Nil(tac.LineOfAccounting.LoaScrtyCoopLnItmID)
		suite.Nil(tac.LineOfAccounting.LoaAgncDsbrCd)
		suite.Nil(tac.LineOfAccounting.LoaAgncAcntngCd)
		suite.Nil(tac.LineOfAccounting.LoaFndCntrID)
		suite.Nil(tac.LineOfAccounting.LoaCstCntrID)
		suite.Nil(tac.LineOfAccounting.LoaPrjID)
		suite.Nil(tac.LineOfAccounting.LoaActvtyID)
		suite.Nil(tac.LineOfAccounting.LoaCstCd)
		suite.Nil(tac.LineOfAccounting.LoaWrkOrdID)
		suite.Nil(tac.LineOfAccounting.LoaFnclArID)
		suite.Nil(tac.LineOfAccounting.LoaScrtyCoopCustCd)
		suite.Nil(tac.LineOfAccounting.LoaEndFyTx)
		suite.Nil(tac.LineOfAccounting.LoaBgFyTx)
		suite.Nil(tac.LineOfAccounting.LoaBgtRstrCd)
		suite.Nil(tac.LineOfAccounting.LoaBgtSubActCd)

		suite.Nil(tac.TacSysID)
		suite.Nil(tac.LoaSysID)
		suite.Nil(tac.TacFyTxt)
		suite.Nil(tac.TacFnBlModCd)
		suite.Nil(tac.OrgGrpDfasCd)
		suite.Nil(tac.TacMvtDsgID)
		suite.Nil(tac.TacTyCd)
		suite.Nil(tac.TacUseCd)
		suite.Nil(tac.TacMajClmtID)
		suite.Nil(tac.TacBillActTxt)
		suite.Nil(tac.TacCostCtrNm)
		suite.Nil(tac.Buic)
		suite.Nil(tac.TacHistCd)
		suite.Nil(tac.TacStatCd)
		suite.Nil(tac.TrnsprtnAcntTx)
		suite.Nil(tac.TrnsprtnAcntBgnDt)
		suite.Nil(tac.TrnsprtnAcntEndDt)
		suite.Nil(tac.DdActvtyAdrsID)
		suite.Nil(tac.TacBlldAddFrstLnTx)
		suite.Nil(tac.TacBlldAddScndLnTx)
		suite.Nil(tac.TacBlldAddThrdLnTx)
		suite.Nil(tac.TacBlldAddFrthLnTx)
		suite.Nil(tac.TacFnctPocNm)
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
					LoaSysID: models.StringPointer("4321"),
				},
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal("1234", tac.TAC)
		suite.Equal("4321", *tac.LineOfAccounting.LoaSysID)
	})

	suite.Run("Successful creation of a fully-filled TAC", func() {
		// Under test:      BuildFullTransportationAccountingCode
		// Set up:          Create a TAC with no customizations or traits
		// Expected outcome:TAC should be created with custom value
		tac := BuildFullTransportationAccountingCode(suite.DB())

		// VALIDATE RESULTS
		suite.Equal(defaultTac.TAC, tac.TAC)

		suite.Equal(*tac.LoaID, tac.LineOfAccounting.ID)
		suite.NotNil(tac.LineOfAccounting.LoaSysID)
		suite.NotNil(tac.LineOfAccounting.LoaDptID)
		suite.NotNil(tac.LineOfAccounting.LoaTnsfrDptNm)
		suite.NotNil(tac.LineOfAccounting.LoaBafID)
		suite.NotNil(tac.LineOfAccounting.LoaTrsySfxTx)
		suite.NotNil(tac.LineOfAccounting.LoaMajClmNm)
		suite.NotNil(tac.LineOfAccounting.LoaOpAgncyID)
		suite.NotNil(tac.LineOfAccounting.LoaAlltSnID)
		suite.NotNil(tac.LineOfAccounting.LoaPgmElmntID)
		suite.NotNil(tac.LineOfAccounting.LoaTskBdgtSblnTx)
		suite.NotNil(tac.LineOfAccounting.LoaDfAgncyAlctnRcpntID)
		suite.NotNil(tac.LineOfAccounting.LoaJbOrdNm)
		suite.NotNil(tac.LineOfAccounting.LoaSbaltmtRcpntID)
		suite.NotNil(tac.LineOfAccounting.LoaWkCntrRcpntNm)
		suite.NotNil(tac.LineOfAccounting.LoaMajRmbsmtSrcID)
		suite.NotNil(tac.LineOfAccounting.LoaDtlRmbsmtSrcID)
		suite.NotNil(tac.LineOfAccounting.LoaCustNm)
		suite.NotNil(tac.LineOfAccounting.LoaObjClsID)
		suite.NotNil(tac.LineOfAccounting.LoaSrvSrcID)
		suite.NotNil(tac.LineOfAccounting.LoaSpclIntrID)
		suite.NotNil(tac.LineOfAccounting.LoaBdgtAcntClsNm)
		suite.NotNil(tac.LineOfAccounting.LoaDocID)
		suite.NotNil(tac.LineOfAccounting.LoaClsRefID)
		suite.NotNil(tac.LineOfAccounting.LoaInstlAcntgActID)
		suite.NotNil(tac.LineOfAccounting.LoaLclInstlID)
		suite.NotNil(tac.LineOfAccounting.LoaFmsTrnsactnID)
		suite.NotNil(tac.LineOfAccounting.LoaDscTx)
		suite.NotNil(tac.LineOfAccounting.LoaBgnDt)
		suite.NotNil(tac.LineOfAccounting.LoaEndDt)
		suite.NotNil(tac.LineOfAccounting.LoaFnctPrsNm)
		suite.NotNil(tac.LineOfAccounting.LoaStatCd)
		suite.NotNil(tac.LineOfAccounting.LoaHistStatCd)
		suite.NotNil(tac.LineOfAccounting.LoaHsGdsCd)
		suite.NotNil(tac.LineOfAccounting.OrgGrpDfasCd)
		suite.NotNil(tac.LineOfAccounting.LoaUic)
		suite.NotNil(tac.LineOfAccounting.LoaTrnsnID)
		suite.NotNil(tac.LineOfAccounting.LoaSubAcntID)
		suite.NotNil(tac.LineOfAccounting.LoaBetCd)
		suite.NotNil(tac.LineOfAccounting.LoaFndTyFgCd)
		suite.NotNil(tac.LineOfAccounting.LoaBgtLnItmID)
		suite.NotNil(tac.LineOfAccounting.LoaScrtyCoopImplAgncCd)
		suite.NotNil(tac.LineOfAccounting.LoaScrtyCoopDsgntrCd)
		suite.NotNil(tac.LineOfAccounting.LoaScrtyCoopLnItmID)
		suite.NotNil(tac.LineOfAccounting.LoaAgncDsbrCd)
		suite.NotNil(tac.LineOfAccounting.LoaAgncAcntngCd)
		suite.NotNil(tac.LineOfAccounting.LoaFndCntrID)
		suite.NotNil(tac.LineOfAccounting.LoaCstCntrID)
		suite.NotNil(tac.LineOfAccounting.LoaPrjID)
		suite.NotNil(tac.LineOfAccounting.LoaActvtyID)
		suite.NotNil(tac.LineOfAccounting.LoaCstCd)
		suite.NotNil(tac.LineOfAccounting.LoaWrkOrdID)
		suite.NotNil(tac.LineOfAccounting.LoaFnclArID)
		suite.NotNil(tac.LineOfAccounting.LoaScrtyCoopCustCd)
		suite.NotNil(tac.LineOfAccounting.LoaEndFyTx)
		suite.NotNil(tac.LineOfAccounting.LoaBgFyTx)
		suite.NotNil(tac.LineOfAccounting.LoaBgtRstrCd)
		suite.NotNil(tac.LineOfAccounting.LoaBgtSubActCd)

		suite.NotNil(tac.TacSysID)
		suite.NotNil(tac.LoaSysID)
		suite.NotNil(tac.TacFyTxt)
		suite.NotNil(tac.TacFnBlModCd)
		suite.NotNil(tac.OrgGrpDfasCd)
		suite.NotNil(tac.TacMvtDsgID)
		suite.NotNil(tac.TacTyCd)
		suite.NotNil(tac.TacUseCd)
		suite.NotNil(tac.TacMajClmtID)
		suite.NotNil(tac.TacBillActTxt)
		suite.NotNil(tac.TacCostCtrNm)
		suite.NotNil(tac.Buic)
		suite.NotNil(tac.TacHistCd)
		suite.NotNil(tac.TacStatCd)
		suite.NotNil(tac.TrnsprtnAcntTx)
		suite.NotNil(tac.TrnsprtnAcntBgnDt)
		suite.NotNil(tac.TrnsprtnAcntEndDt)
		suite.NotNil(tac.DdActvtyAdrsID)
		suite.NotNil(tac.TacBlldAddFrstLnTx)
		suite.NotNil(tac.TacBlldAddScndLnTx)
		suite.NotNil(tac.TacBlldAddThrdLnTx)
		suite.NotNil(tac.TacBlldAddFrthLnTx)
		suite.NotNil(tac.TacFnctPocNm)
	})
}
