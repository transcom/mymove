package iws

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var ssnRegex = regexp.MustCompile("^\\d{9}$")

func buildPidsURL(host string, custNum string, ssn string, lastName string, firstName string) (string, error) {
	if !ssnRegex.MatchString(ssn) {
		return "", errors.New("SSN must be exactly 9 digits")
	}

	var urlBuilder strings.Builder
	baseURL := fmt.Sprintf("https://%s"+
		"/appj/rbs/rest/op=pids-P/customer=%s"+
		"/schemaName=get_cac_data/schemaVersion=1.0/PN_ID=%s"+
		"/PN_ID_TYP_CD=S/PN_LST_NM=", host, custNum, ssn)
	urlBuilder.WriteString(baseURL)
	l := len(lastName)
	if l > 26 {
		// Last names are limited to 26 characters in IWS
		l = 26
	}
	urlBuilder.WriteString(lastName[:l])

	// The first name is optional
	l = len(firstName)
	if l > 0 {
		urlBuilder.WriteString("/PN_1ST_NM=")
		if l > 20 {
			l = 20
		}
		// First names are limited to 20 characters in IWS
		urlBuilder.WriteString(firstName[:l])
	}

	return urlBuilder.String(), nil
}

func parsePidsResponse(data []byte) (MatchReasonCode, uint64, *Person, []Personnel, error) {
	rec := Record{}
	unmarshalErr := xml.Unmarshal(data, &rec)
	if unmarshalErr != nil {
		// Couldn't unmarshal as a record, try as an RbsError next
		rbsError := RbsError{}
		unmarshalErr = xml.Unmarshal(data, &rbsError)
		if unmarshalErr == nil {
			return MatchReasonCodeNone, 0, nil, []Personnel{}, &rbsError
		}
		return MatchReasonCodeNone, 0, nil, []Personnel{}, unmarshalErr
	}

	reason := rec.AdrRecord.PidsRecord.MtchRsnCd
	if reason == MatchReasonCodeNone {
		return reason, 0, nil, []Personnel{}, nil
	}

	return reason, rec.AdrRecord.PidsRecord.Edipi, rec.AdrRecord.Person, rec.AdrRecord.Personnel, nil
}

// GetPersonUsingSSNParams contains person-specific query parameters for GetPidsUsingSSN
type GetPersonUsingSSNParams struct {
	Ssn       string
	LastName  string
	FirstName string
}

// GetPersonUsingSSN retrieves personal information (including EDIPI) through the IWS:RBS REST API using a SSN, last name, and optionally a first name
// If matched succesfully, it returns the EDIPI, the full name and SSN information, and the personnel information for each of the organizations the person belongs to
func GetPersonUsingSSN(client http.Client, host string, custNum string, params GetPersonUsingSSNParams) (MatchReasonCode, uint64, *Person, []Personnel, error) {
	url, err := buildPidsURL(host, custNum, params.Ssn, params.LastName, params.FirstName)
	if err != nil {
		return MatchReasonCodeNone, 0, nil, []Personnel{}, err
	}

	resp, getErr := client.Get(url)
	// Interesting fact: RBS responds 200 OK, not 404 Not Found, if there are no matches
	if getErr != nil {
		return MatchReasonCodeNone, 0, nil, []Personnel{}, getErr
	}

	defer resp.Body.Close()
	data, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return MatchReasonCodeNone, 0, nil, []Personnel{}, readErr
	}

	return parsePidsResponse(data)
}
