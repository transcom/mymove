package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_AllFieldsOptionalCanSave() {
	loa := &models.LineOfAccounting{
		ID: uuid.Must(uuid.NewV4()),
	}

	err := suite.DB().Save(loa)
	suite.NoError(err)
}

func (suite *ModelSuite) Test_AllFieldsPresentCanSave() {
	loa := &models.LineOfAccounting{
		ID:                     uuid.Must(uuid.NewV4()),
		LoaSysID:               models.IntPointer(123456),
		LoaDptID:               models.StringPointer("12"),
		LoaTnsfrDptNm:          models.StringPointer("1234"),
		LoaBafID:               models.StringPointer("1234"),
		LoaTrsySfxTx:           models.StringPointer("1234"),
		LoaMajClmNm:            models.StringPointer("1234"),
		LoaOpAgncyID:           models.StringPointer("1234"),
		LoaAlltSnID:            models.StringPointer("12345"),
		LoaPgmElmntID:          models.StringPointer("123456789012"),
		LoaTskBdgtSblnTx:       models.StringPointer("88888888"),
		LoaDfAgncyAlctnRcpntID: models.StringPointer("1234"),
		LoaJbOrdNm:             models.StringPointer("1234567890"),
		LoaSbaltmtRcpntID:      models.StringPointer("1"),
		LoaWkCntrRcpntNm:       models.StringPointer("123456"),
		LoaMajRmbsmtSrcID:      models.StringPointer("1"),
		LoaDtlRmbsmtSrcID:      models.StringPointer("123"),
		LoaCustNm:              models.StringPointer("123456"),
		LoaObjClsID:            models.StringPointer("123456"),
		LoaSrvSrcID:            models.StringPointer("1"),
		LoaSpclIntrID:          models.StringPointer("12"),
		LoaBdgtAcntClsNm:       models.StringPointer("12345678"),
		LoaDocID:               models.StringPointer("123456789012345"),
		LoaClsRefID:            models.StringPointer("12"),
		LoaInstlAcntgActID:     models.StringPointer("123456"),
		LoaLclInstlID:          models.StringPointer("123456789012345678"),
		LoaFmsTrnsactnID:       models.StringPointer("123456789012"),
		LoaDscTx:               models.StringPointer("LoaDscTx"),
		LoaBgnDt:               models.TimePointer(time.Now()),
		LoaEndDt:               models.TimePointer(time.Now()),
		LoaFnctPrsNm:           models.StringPointer("LoaFnctPrsNm"),
		LoaStatCd:              models.StringPointer("1"),
		LoaHistStatCd:          models.StringPointer("1"),
		LoaHsGdsCd:             models.StringPointer("12"),
		OrgGrpDfasCd:           models.StringPointer("12"),
		LoaUic:                 models.StringPointer("123456"),
		LoaTrnsnID:             models.StringPointer("LoaTrnsnID"),
		LoaSubAcntID:           models.StringPointer("123"),
		LoaBetCd:               models.StringPointer("1234"),
		LoaFndTyFgCd:           models.StringPointer("1"),
		LoaBgtLnItmID:          models.StringPointer("12345678"),
		LoaScrtyCoopImplAgncCd: models.StringPointer("1"),
		LoaScrtyCoopDsgntrCd:   models.StringPointer("1234"),
		LoaScrtyCoopLnItmID:    models.StringPointer("123"),
		LoaAgncDsbrCd:          models.StringPointer("123456"),
		LoaAgncAcntngCd:        models.StringPointer("123456"),
		LoaFndCntrID:           models.StringPointer("123456789012"),
		LoaCstCntrID:           models.StringPointer("1234567890123456"),
		LoaPrjID:               models.StringPointer("123456789012"),
		LoaActvtyID:            models.StringPointer("12345678901"),
		LoaCstCd:               models.StringPointer("1234567890123456"),
		LoaWrkOrdID:            models.StringPointer("1234567890123456"),
		LoaFnclArID:            models.StringPointer("123456"),
		LoaScrtyCoopCustCd:     models.StringPointer("12"),
		LoaEndFyTx:             models.IntPointer(time.Now().Year()),
		LoaBgFyTx:              models.IntPointer(time.Now().Year()),
		LoaBgtRstrCd:           models.StringPointer("1"),
		LoaBgtSubActCd:         models.StringPointer("1234"),
	}

	err := suite.DB().Save(loa)
	suite.NoError(err)
}
