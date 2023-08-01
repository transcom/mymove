package models_test

import (
	"reflect"
	"testing"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

// This test will check if any field names are not being mapped properly. This test is not finalized.
func TestTransportationAccountingCodeMapForUnusedFields(t *testing.T) {
	t.Skip("Skipping this test until the fields and usecase has been finalized.")

	// Example of TransportationAccountingCodeTrdmFileRecord
	tacFileRecord := models.TransportationAccountingCodeTrdmFileRecord{
		TRNSPRTN_ACNT_CD:        "4EVR",
		TAC_SYS_ID:              "3080819",
		LOA_SYS_ID:              "55555555",
		TAC_FY_TXT:              "2023",
		TAC_FN_BL_MOD_CD:        "W",
		ORG_GRP_DFAS_CD:         "HS",
		TAC_MVT_DSG_ID:          "",
		TAC_TY_CD:               "O",
		TAC_USE_CD:              "N",
		TAC_MAJ_CLMT_ID:         "012345",
		TAC_BILL_ACT_TXT:        "123456",
		TAC_COST_CTR_NM:         "012345",
		BUIC:                    "",
		TAC_HIST_CD:             "",
		TAC_STAT_CD:             "I",
		TRNSPRTN_ACNT_TX:        "For the purpose of MacDill AFB transporting to Scott AFB",
		TRNSPRTN_ACNT_BGN_DT:    "2022-10-01 00:00:00",
		TRNSPRTN_ACNT_END_DT:    "2023-09-30 00:00:00",
		DD_ACTVTY_ADRS_ID:       "A12345",
		TAC_BLLD_ADD_FRST_LN_TX: "MacDill",
		TAC_BLLD_ADD_SCND_LN_TX: "Second Address Line",
		TAC_BLLD_ADD_THRD_LN_TX: "",
		TAC_BLLD_ADD_FRTH_LN_TX: "TAMPA FL 33621",
		TAC_FNCT_POC_NM:         "THISISNOTAREALPERSON@USCG.MIL",
	}

	mappedStruct := models.MapTransportationAccountingCodeFileRecordToInternalStruct(tacFileRecord)

	reflectedMappedStruct := reflect.TypeOf(mappedStruct)
	reflectedTacFileRecord := reflect.TypeOf(tacFileRecord)

	// Iterate through each field in the tacRecord struct for the comparison
	for i := 0; i < reflectedTacFileRecord.NumField(); i++ {
		fieldName := reflectedTacFileRecord.Field(i).Name

		// Check if this field exists in the reflectedMappedStruct
		_, exists := reflectedMappedStruct.FieldByName(fieldName)

		// Error if the field isn't found in the reflectedMappedStruct
		if !exists {
			t.Errorf("Field '%s' in TransportationAccountingCodeTrdmFileRecord is not used in MapTransportationAccountingCodeFileRecordToInternalStruct function", fieldName)
		}
	}
}

// This function will test the receival of a parsed TAC that has undergone the pipe delimited .txt file parser. It will test
// that the received values correctly map to our internal TAC struct. For example, our Transporation Accounting Code is called
// "TAC" in its struct, however when it is received in pipe delimited format it will be received as "TRNSPRTN_ACNT_CD".
// This function makes sure it gets connected properly.
func TestTransportationAccountingCodeMapToInternal(t *testing.T) {

	tacFileRecord := models.TransportationAccountingCodeTrdmFileRecord{
		TRNSPRTN_ACNT_CD: "4EVR",
	}

	mappedTacFileRecord := models.MapTransportationAccountingCodeFileRecordToInternalStruct(tacFileRecord)

	// Check that the TRNSPRTN_ACNT_CD field in the original struct was correctly
	// mapped to the TAC field in the resulting struct
	if mappedTacFileRecord.TAC != tacFileRecord.TRNSPRTN_ACNT_CD {
		t.Errorf("Expected TAC to be '%s', got '%s'", tacFileRecord.TRNSPRTN_ACNT_CD, mappedTacFileRecord.TAC)
	}
}

func (suite *ModelSuite) TestCanSaveValidTac() {
	tac := models.TransportationAccountingCode{
		TAC: "Tac1",
	}

	suite.MustCreate(&tac)
}

func (suite *ModelSuite) TestInvalidTac() {
	tac := models.TransportationAccountingCode{}

	expErrors := map[string][]string{
		"tac": {"TAC can not be blank."},
	}

	verrs, err := suite.DB().ValidateAndSave(&tac)

	suite.Equal(expErrors, verrs.Errors)
	suite.NoError(err)
}

func (suite *ModelSuite) TestCanSaveAndFetchTac() {
	// Can save
	tac := factory.BuildFullTransportationAccountingCode(suite.DB())

	suite.MustSave(&tac)

	// Can fetch tac with associations
	var fetchedTac models.TransportationAccountingCode
	err := suite.DB().Where("tac = $1", tac.TAC).Eager("LineOfAccounting").First(&fetchedTac)

	suite.NoError(err)
	suite.Equal(tac.TAC, fetchedTac.TAC)
	suite.NotNil(*fetchedTac.LineOfAccounting.LoaSysID)
}
