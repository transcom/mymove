package factory

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func BuildLineOfAccounting(db *pop.Connection, customs []Customization, traits []Trait) models.LineOfAccounting {
	customs = setupCustomizations(customs, traits)

	var cLineOfAccounting models.LineOfAccounting
	if result := findValidCustomization(customs, LineOfAccounting); result != nil {
		cLineOfAccounting = result.Model.(models.LineOfAccounting)
		if result.LinkOnly {
			return cLineOfAccounting
		}
	}

	lineOfAccounting := models.LineOfAccounting{}

	testdatagen.MergeModels(&lineOfAccounting, cLineOfAccounting)

	if db != nil {
		mustCreate(db, &lineOfAccounting)
	}

	return lineOfAccounting
}

func BuildDefaultLineOfAccounting(db *pop.Connection) models.LineOfAccounting {
	return BuildLineOfAccounting(db, nil, nil)
}

func BuildFullLineOfAccounting(db *pop.Connection, customs []Customization, traits []Trait) models.LineOfAccounting {
	customs = setupCustomizations(customs, traits)

	var cLineOfAccounting models.LineOfAccounting
	if result := findValidCustomization(customs, LineOfAccounting); result != nil {
		cLineOfAccounting = result.Model.(models.LineOfAccounting)
		if result.LinkOnly {
			return cLineOfAccounting
		}
	}

	now := time.Now()
	later := now.AddDate(1, 0, 0)

	lineOfAccounting := models.LineOfAccounting{
		LoaSysID:               models.StringPointer(MakeRandomString(20)),
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
		LoaInstlAcntgActID:     models.StringPointer("102"),
		LoaLclInstlID:          models.StringPointer("123456789012345678"),
		LoaFmsTrnsactnID:       models.StringPointer("123456789012"),
		LoaDscTx:               models.StringPointer("LoaDscTx"),
		LoaBgnDt:               models.TimePointer(now),
		LoaEndDt:               models.TimePointer(later),
		LoaFnctPrsNm:           models.StringPointer("LoaFnctPrsNm"),
		LoaStatCd:              models.StringPointer("1"),
		LoaHistStatCd:          models.StringPointer("1"),
		LoaHsGdsCd:             models.StringPointer("12"),
		OrgGrpDfasCd:           models.StringPointer("12"),
		LoaUic:                 models.StringPointer("123456"),
		LoaTrnsnID:             models.StringPointer("123"),
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
		LoaEndFyTx:             models.IntPointer(later.Year()),
		LoaBgFyTx:              models.IntPointer(now.Year()),
		LoaBgtRstrCd:           models.StringPointer("1"),
		LoaBgtSubActCd:         models.StringPointer("1234"),
	}

	testdatagen.MergeModels(&lineOfAccounting, cLineOfAccounting)

	if db != nil {
		mustCreate(db, &lineOfAccounting)
	}

	return lineOfAccounting
}
