package factory

import "github.com/transcom/mymove/pkg/models"

func (suite *FactorySuite) TestBuildLineOfAccounting() {

	suite.Run("Successful creation of default TAC", func() {
		// Under test:      BuildLineOfAccounting
		// Set up:          Create a TAC with no customizations or traits
		// Expected outcome:Line of accounting should be created with empty rows, since all fields are optional
		loa := BuildLineOfAccounting(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Nil(loa.LoaSysID)
		suite.Nil(loa.LoaDptID)
		suite.Nil(loa.LoaTnsfrDptNm)
		suite.Nil(loa.LoaBafID)
		suite.Nil(loa.LoaTrsySfxTx)
		suite.Nil(loa.LoaMajClmNm)
		suite.Nil(loa.LoaOpAgncyID)
		suite.Nil(loa.LoaAlltSnID)
		suite.Nil(loa.LoaPgmElmntID)
		suite.Nil(loa.LoaTskBdgtSblnTx)
		suite.Nil(loa.LoaDfAgncyAlctnRcpntID)
		suite.Nil(loa.LoaJbOrdNm)
		suite.Nil(loa.LoaSbaltmtRcpntID)
		suite.Nil(loa.LoaWkCntrRcpntNm)
		suite.Nil(loa.LoaMajRmbsmtSrcID)
		suite.Nil(loa.LoaDtlRmbsmtSrcID)
		suite.Nil(loa.LoaCustNm)
		suite.Nil(loa.LoaObjClsID)
		suite.Nil(loa.LoaSrvSrcID)
		suite.Nil(loa.LoaSpclIntrID)
		suite.Nil(loa.LoaBdgtAcntClsNm)
		suite.Nil(loa.LoaDocID)
		suite.Nil(loa.LoaClsRefID)
		suite.Nil(loa.LoaInstlAcntgActID)
		suite.Nil(loa.LoaLclInstlID)
		suite.Nil(loa.LoaFmsTrnsactnID)
		suite.Nil(loa.LoaDscTx)
		suite.Nil(loa.LoaBgnDt)
		suite.Nil(loa.LoaEndDt)
		suite.Nil(loa.LoaFnctPrsNm)
		suite.Nil(loa.LoaStatCd)
		suite.Nil(loa.LoaHistStatCd)
		suite.Nil(loa.LoaHsGdsCd)
		suite.Nil(loa.OrgGrpDfasCd)
		suite.Nil(loa.LoaUic)
		suite.Nil(loa.LoaTrnsnID)
		suite.Nil(loa.LoaSubAcntID)
		suite.Nil(loa.LoaBetCd)
		suite.Nil(loa.LoaFndTyFgCd)
		suite.Nil(loa.LoaBgtLnItmID)
		suite.Nil(loa.LoaScrtyCoopImplAgncCd)
		suite.Nil(loa.LoaScrtyCoopDsgntrCd)
		suite.Nil(loa.LoaScrtyCoopLnItmID)
		suite.Nil(loa.LoaAgncDsbrCd)
		suite.Nil(loa.LoaAgncAcntngCd)
		suite.Nil(loa.LoaFndCntrID)
		suite.Nil(loa.LoaCstCntrID)
		suite.Nil(loa.LoaPrjID)
		suite.Nil(loa.LoaActvtyID)
		suite.Nil(loa.LoaCstCd)
		suite.Nil(loa.LoaWrkOrdID)
		suite.Nil(loa.LoaFnclArID)
		suite.Nil(loa.LoaScrtyCoopCustCd)
		suite.Nil(loa.LoaEndFyTx)
		suite.Nil(loa.LoaBgFyTx)
		suite.Nil(loa.LoaBgtRstrCd)
		suite.Nil(loa.LoaBgtSubActCd)
	})

	suite.Run("Successful creation of a Line of Accounting with customization", func() {
		// Under test:      BuildLineOfAccounting
		// Set up:          Create a TAC with no customizations or traits
		// Expected outcome:Line of accounting should be created with custom value
		loa := BuildLineOfAccounting(suite.DB(), []Customization{
			{
				Model: models.LineOfAccounting{
					LoaSysID: models.StringPointer("4321"),
				},
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal("4321", *loa.LoaSysID)
	})

	suite.Run("Successful creation of a fully-filled Line of Accounting", func() {
		// Under test:      BuildFullLineOfAccounting
		// Set up:          Create a Line of Accounting with no customizations or traits
		// Expected outcome:TAC should be created with custom value
		loa := BuildFullLineOfAccounting(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.NotNil(loa.LoaSysID)
		suite.NotNil(loa.LoaDptID)
		suite.NotNil(loa.LoaTnsfrDptNm)
		suite.NotNil(loa.LoaBafID)
		suite.NotNil(loa.LoaTrsySfxTx)
		suite.NotNil(loa.LoaMajClmNm)
		suite.NotNil(loa.LoaOpAgncyID)
		suite.NotNil(loa.LoaAlltSnID)
		suite.NotNil(loa.LoaPgmElmntID)
		suite.NotNil(loa.LoaTskBdgtSblnTx)
		suite.NotNil(loa.LoaDfAgncyAlctnRcpntID)
		suite.NotNil(loa.LoaJbOrdNm)
		suite.NotNil(loa.LoaSbaltmtRcpntID)
		suite.NotNil(loa.LoaWkCntrRcpntNm)
		suite.NotNil(loa.LoaMajRmbsmtSrcID)
		suite.NotNil(loa.LoaDtlRmbsmtSrcID)
		suite.NotNil(loa.LoaCustNm)
		suite.NotNil(loa.LoaObjClsID)
		suite.NotNil(loa.LoaSrvSrcID)
		suite.NotNil(loa.LoaSpclIntrID)
		suite.NotNil(loa.LoaBdgtAcntClsNm)
		suite.NotNil(loa.LoaDocID)
		suite.NotNil(loa.LoaClsRefID)
		suite.NotNil(loa.LoaInstlAcntgActID)
		suite.NotNil(loa.LoaLclInstlID)
		suite.NotNil(loa.LoaFmsTrnsactnID)
		suite.NotNil(loa.LoaDscTx)
		suite.NotNil(loa.LoaBgnDt)
		suite.NotNil(loa.LoaEndDt)
		suite.NotNil(loa.LoaFnctPrsNm)
		suite.NotNil(loa.LoaStatCd)
		suite.NotNil(loa.LoaHistStatCd)
		suite.NotNil(loa.LoaHsGdsCd)
		suite.NotNil(loa.OrgGrpDfasCd)
		suite.NotNil(loa.LoaUic)
		suite.NotNil(loa.LoaTrnsnID)
		suite.NotNil(loa.LoaSubAcntID)
		suite.NotNil(loa.LoaBetCd)
		suite.NotNil(loa.LoaFndTyFgCd)
		suite.NotNil(loa.LoaBgtLnItmID)
		suite.NotNil(loa.LoaScrtyCoopImplAgncCd)
		suite.NotNil(loa.LoaScrtyCoopDsgntrCd)
		suite.NotNil(loa.LoaScrtyCoopLnItmID)
		suite.NotNil(loa.LoaAgncDsbrCd)
		suite.NotNil(loa.LoaAgncAcntngCd)
		suite.NotNil(loa.LoaFndCntrID)
		suite.NotNil(loa.LoaCstCntrID)
		suite.NotNil(loa.LoaPrjID)
		suite.NotNil(loa.LoaActvtyID)
		suite.NotNil(loa.LoaCstCd)
		suite.NotNil(loa.LoaWrkOrdID)
		suite.NotNil(loa.LoaFnclArID)
		suite.NotNil(loa.LoaScrtyCoopCustCd)
		suite.NotNil(loa.LoaEndFyTx)
		suite.NotNil(loa.LoaBgFyTx)
		suite.NotNil(loa.LoaBgtRstrCd)
		suite.NotNil(loa.LoaBgtSubActCd)
	})
}
