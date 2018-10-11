package iws

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GetPersonUsingEDIPI retrieves personal information through the IWS:RBS REST API using that person's EDIPI (aka DOD ID number).
// If matched succesfully, it returns the full name and SSN information, as well as the personnel information for each of the organizations the person belongs to
func GetPersonUsingEDIPI(client http.Client, host string, custNum string, edipi uint64) (*Person, []Personnel, error) {
	url := fmt.Sprintf("https://%s"+
		"/appj/rbs/rest/op=edi/customer=%s"+
		"/schemaName=get_cac_data/schemaVersion=1.0/DOD_EDI_PN_ID=%d",
		host, custNum, edipi)

	resp, getErr := client.Get(url)
	// Interesting fact: RBS responds 200 OK, not 404 Not Found, if there are no matches
	if getErr != nil {
		return nil, []Personnel{}, getErr
	}

	defer resp.Body.Close()
	data, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, []Personnel{}, readErr
	}

	rec := Record{}
	unmarshalErr := xml.Unmarshal([]byte(data), &rec)
	if unmarshalErr != nil {
		// Couldn't unmarshal as a record, try as an RbsError next
		rbsError := RbsError{}
		unmarshalErr = xml.Unmarshal([]byte(data), &rbsError)
		if unmarshalErr == nil {
			return nil, []Personnel{}, &rbsError
		}
		return nil, []Personnel{}, unmarshalErr
	}

	return rec.AdrRecord.Person, rec.AdrRecord.Personnel, nil
}
