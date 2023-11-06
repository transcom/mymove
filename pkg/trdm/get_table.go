package trdm

import (
	"time"

	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// const successResponseString = "Successful"
const lineOfAccounting = "LN_OF_ACCT"
const transportationAccountingCode = "TRNSPRTN_ACNT"

// Fetch Transportation Accounting Codes from DB and return the list of records if the updated_at field is < the returned LastTableUpdate updated time.
// Ex:
//
//	LastTableUpdate : 2023-08-30 15:24:13.19931
//	updated_at: 2023-08-29 15:24:13.19931
//
// Because updated_at is before LastTableUpdate the DB will return records that match this case.
//
//	returns []models.TransportationAccountingCode, error
func FetchTACRecordsByTime(appcontext appcontext.AppContext, time time.Time) ([]models.TransportationAccountingCode, error) {
	var tacCodes []models.TransportationAccountingCode
	err := appcontext.DB().Select("*").Where("updated_at < $1", time).All(&tacCodes)

	if err != nil {
		return tacCodes, errors.Wrap(err, "Fetch line items query failed")
	}

	return tacCodes, nil
}

// Fetch Line Of Accounting records from DB and return the list of records if the updated_at field is < the returned LastTableUpdate updated time.
// Ex:
//
//	LastTableUpdate : 2023-08-30 15:24:13.19931
//	updated_at: 2023-08-29 15:24:13.19931
//
// Because updated_at is before LastTableUpdate the DB will return records that match this case.
//
//	returns []models.LineOfAccounting, error
func FetchLOARecordsByTime(appcontext appcontext.AppContext, time time.Time) ([]models.LineOfAccounting, error) {
	var loa []models.LineOfAccounting
	err := appcontext.DB().Select("*").Where("updated_at < $1", time).All(&loa)
	if err != nil {
		return loa, errors.Wrap(err, "Fetch line items query failed")
	}
	return loa, nil
}

// Determines if call is needed to be made
// If the DB does not return any records we do not need make a call to GetTable and update our local mapping
//   - appCtx: Application Context
//   - physicalName: Table Name (Will be either TAC or LOA)
//   - lastUpdate: Returned date time from LastTableUpdate Soap Request
//
// returns error
func GetTable(appCtx appcontext.AppContext, physicalName string, lastUpdate time.Time) error {

	switch physicalName {
	case lineOfAccounting:
		loaRecords, loaFetchErr := FetchLOARecordsByTime(appCtx, lastUpdate)

		if loaFetchErr != nil {
			return loaFetchErr
		}

		if len(loaRecords) > 0 {
			// TODO: Send off the call
			return loaFetchErr // Remove me
		}
	case transportationAccountingCode:
		tacRecords, fetchErr := FetchTACRecordsByTime(appCtx, lastUpdate)
		if fetchErr != nil {
			return fetchErr
		}
		if len(tacRecords) > 0 {
			// TODO: Send off the call
			return fetchErr // Remove me
		}
	}

	return nil
}

// Parses pipedelimited file attachment from GetTable webservice and saves records to database
//
//	returns error
// TODO: Impelement again
/*
func parseGetTableResponse(appcontext appcontext.AppContext, response *gosoap.Response, physicalName string) error {
	reader := bytes.NewReader(response.Payload)
	switch physicalName {
	case lineOfAccounting:
		loaCodes, err := loa.Parse(reader)
		if err != nil {
			return err
		}
		saveErr := saveLoaCodes(appcontext, loaCodes)
		if saveErr != nil {
			return saveErr
		}
	case transportationAccountingCode:
		tacCodes, err := tac.Parse(reader)
		consolidatedTacs := tac.ConsolidateDuplicateTACsDesiredFromTRDM(tacCodes)
		if err != nil {
			return err
		}
		if saveErr := saveTacCodes(appcontext, consolidatedTacs); saveErr != nil {
			return saveErr
		}
	default:
		return nil
	}
	return nil
}
*/

// Saves TAC Code slice to DB and updates records
// TODO: Implement again
/*
func saveTacCodes(appcontext appcontext.AppContext, tacCodes []models.TransportationAccountingCode) error {
	saveErr := appcontext.DB().Update(tacCodes)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

// Saves LOA Code slice to DB and updates records
func saveLoaCodes(appcontext appcontext.AppContext, loa []models.LineOfAccounting) error {
	saveErr := appcontext.DB().Update(loa)
	if saveErr != nil {
		return saveErr
	}
	return nil
}
*/
