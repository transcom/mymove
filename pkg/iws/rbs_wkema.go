package iws

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

//func GetInfoFromWorkEmail(client http.Client, host string, custNum string, workEmailAddr string) (uint64, Person, Personnel, error) {
// https://pkict.dmdc.osd.mil/appj/rbs/rest/op=wkEma/customer=2675/schemaName=schema_name/schemaVersion=1.0/EMA_TX=nobody_here@mail.mil
//}

// GetPersonUsingWorkEmail retrieves personal information (including SSN and EDIPI) through the IWS:RBS REST API using a work e-mail address.
// If matched succesfully, it returns the EDIPI, the full name and SSN information, and the personnel information for each of the organizations the person belongs to
func GetPersonUsingWorkEmail(client http.Client, host string, custNum string, workEmail string) (uint64, *Person, []Personnel, error) {
	url := fmt.Sprintf("https://%s"+
		"/appj/rbs/rest/op=wkEma/customer=%s"+
		"/schemaName=get_cac_data/schemaVersion=1.0/EMA_TX=%s",
		host, custNum, workEmail)

	resp, getErr := client.Get(url)
	if getErr != nil {
		return 0, nil, []Personnel{}, getErr
	}

	defer resp.Body.Close()
	data, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return 0, nil, []Personnel{}, readErr
	}

	rec := Record{}
	unmarshalErr := xml.Unmarshal([]byte(data), &rec)
	if unmarshalErr != nil {
		// Couldn't unmarshal as a record, try as an RbsError next
		rbsError := RbsError{}
		unmarshalErr = xml.Unmarshal([]byte(data), &rbsError)
		if unmarshalErr == nil {
			return 0, nil, []Personnel{}, &rbsError
		}
		return 0, nil, []Personnel{}, unmarshalErr
	}

	// Not found
	if rec.AdrRecord.WorkEmail == nil {
		return 0, nil, []Personnel{}, nil
	}

	return rec.AdrRecord.WorkEmail.Edipi, rec.AdrRecord.Person, rec.AdrRecord.Personnel, nil
}
