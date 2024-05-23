package factory

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func BuildTransportationAccountingCode(db *pop.Connection, customs []Customization, traits []Trait) models.TransportationAccountingCode {
	customs = setupCustomizations(customs, traits)

	var cTransportationAccountingCode models.TransportationAccountingCode
	if result := findValidCustomization(customs, TransportationAccountingCode); result != nil {
		cTransportationAccountingCode = result.Model.(models.TransportationAccountingCode)
		if result.LinkOnly {
			return cTransportationAccountingCode
		}
	}

	// Create the associated LineOfAccounting model
	lineOfAccounting := BuildLineOfAccounting(db, customs, nil)

	transportationAccountingCode := models.TransportationAccountingCode{
		TAC:              "E01A",
		LineOfAccounting: &lineOfAccounting,
	}

	testdatagen.MergeModels(&transportationAccountingCode, cTransportationAccountingCode)

	if db != nil {
		mustCreate(db, &transportationAccountingCode)
	}

	return transportationAccountingCode
}

func BuildDefaultTransportationAccountingCode(db *pop.Connection) models.TransportationAccountingCode {
	return BuildTransportationAccountingCode(db, nil, nil)
}

func BuildFullTransportationAccountingCode(db *pop.Connection) models.TransportationAccountingCode {
	// Creating as a stub since the Line of Accounting will be Created by the TAC factory
	LineOfAccounting := BuildFullLineOfAccounting(nil, nil, nil)

	defaultCustoms := []Customization{
		{
			Model: LineOfAccounting,
		},
		{
			Model: models.TransportationAccountingCode{
				TacSysID:           models.StringPointer("1234"),
				LoaSysID:           models.StringPointer("1234"),
				TacFyTxt:           models.StringPointer("1234"),
				TacFnBlModCd:       models.StringPointer("1"),
				OrgGrpDfasCd:       models.StringPointer("12"),
				TacMvtDsgID:        models.StringPointer("TacMvtDsgID"),
				TacTyCd:            models.StringPointer("1"),
				TacUseCd:           models.StringPointer("12"),
				TacMajClmtID:       models.StringPointer("123456"),
				TacBillActTxt:      models.StringPointer("123456"),
				TacCostCtrNm:       models.StringPointer("123456"),
				Buic:               models.StringPointer("123456"),
				TacHistCd:          models.StringPointer("1"),
				TacStatCd:          models.StringPointer("1"),
				TrnsprtnAcntTx:     models.StringPointer("TrnsprtnAcntTx"),
				TrnsprtnAcntBgnDt:  models.TimePointer(time.Now()),
				TrnsprtnAcntEndDt:  models.TimePointer(time.Now().AddDate(1, 0, 0)),
				DdActvtyAdrsID:     models.StringPointer("123456"),
				TacBlldAddFrstLnTx: models.StringPointer("TacBlldAddFrstLnTx"),
				TacBlldAddScndLnTx: models.StringPointer("TacBlldAddScndLnTx"),
				TacBlldAddThrdLnTx: models.StringPointer("TacBlldAddThrdLnTx"),
				TacBlldAddFrthLnTx: models.StringPointer("TacBlldAddFrthLnTx"),
				TacFnctPocNm:       models.StringPointer("TacFnctPocNm"),
			},
		},
	}

	return BuildTransportationAccountingCode(db, defaultCustoms, nil)
}

func BuildTransportationAccountingCodeWithoutAttachedLoa(db *pop.Connection, customs []Customization, traits []Trait) models.TransportationAccountingCode {
	customs = setupCustomizations(customs, traits)

	var cTransportationAccountingCode models.TransportationAccountingCode
	if result := findValidCustomization(customs, TransportationAccountingCode); result != nil {
		cTransportationAccountingCode = result.Model.(models.TransportationAccountingCode)
		if result.LinkOnly {
			return cTransportationAccountingCode
		}
	}

	transportationAccountingCode := models.TransportationAccountingCode{
		TAC: "E01A",
	}

	testdatagen.MergeModels(&transportationAccountingCode, cTransportationAccountingCode)

	if db != nil {
		mustCreate(db, &transportationAccountingCode)
	}

	return transportationAccountingCode
}
