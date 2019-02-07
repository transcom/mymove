package iws

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/transcom/mymove/pkg/server"
)

// RBSPersonLookup handles requests to the Real-Time Broker Service
type RBSPersonLookup struct {
	Client http.Client
	Host   string
}

// GetPersonUsingSSNParams contains person-specific query parameters for GetPidsUsingSSN
type GetPersonUsingSSNParams struct {
	Ssn       string
	LastName  string
	FirstName string
}

var myMoveCustNum = "2675"
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&’*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$")
var ssnRegex = regexp.MustCompile("^\\d{9}$")

// GetPersonUsingEDIPI retrieves personal information through the IWS:RBS REST API using that person's EDIPI (aka DOD ID number).
// If matched succesfully, it returns the full name and SSN information, as well as the personnel information for each of the organizations the person belongs to
func (r RBSPersonLookup) GetPersonUsingEDIPI(edipi uint64) (*Person, []Personnel, error) {
	url, err := buildEdiURL(r.Host, myMoveCustNum, edipi)
	if err != nil {
		return nil, []Personnel{}, err
	}

	response, err := r.sendGetRequest(url)
	if err != nil {
		return nil, []Personnel{}, err
	}

	return parseEdiResponse(response)
}

// GetPersonUsingSSN retrieves personal information (including EDIPI) through the IWS:RBS REST API using a SSN, last name, and optionally a first name
// If matched succesfully, it returns the EDIPI, the full name and SSN information, and the personnel information for each of the organizations the person belongs to
func (r RBSPersonLookup) GetPersonUsingSSN(params GetPersonUsingSSNParams) (MatchReasonCode, uint64, *Person, []Personnel, error) {
	url, err := buildPidsURL(r.Host, myMoveCustNum, params.Ssn, params.LastName, params.FirstName)
	if err != nil {
		return MatchReasonCodeNone, 0, nil, []Personnel{}, err
	}

	response, err := r.sendGetRequest(url)
	if err != nil {
		return MatchReasonCodeNone, 0, nil, []Personnel{}, err
	}

	return parsePidsResponse(response)
}

// GetPersonUsingWorkEmail retrieves personal information (including SSN and EDIPI) through the IWS:RBS REST API using a work e-mail address.
// If matched succesfully, it returns the EDIPI, the full name and SSN information, and the personnel information for each of the organizations the person belongs to
func (r RBSPersonLookup) GetPersonUsingWorkEmail(workEmail string) (uint64, *Person, []Personnel, error) {
	url, err := buildWkEmaURL(r.Host, myMoveCustNum, workEmail)
	if err != nil {
		return 0, nil, []Personnel{}, err
	}

	response, err := r.sendGetRequest(url)
	if err != nil {
		return 0, nil, []Personnel{}, err
	}

	return parseWkEmaResponse(response)
}

// NewRBSPersonLookup creates a new instance of RBSPersonLookup. This should
// only be instantiated once
func NewRBSPersonLookup(host string, dodCACertPackage string, certString string, keyString string) (*RBSPersonLookup, error) {
	if host == "" {
		return nil, errors.New("IWS host is not set")
	}

	// Load client cert
	cert, err := tls.X509KeyPair([]byte(certString), []byte(keyString))
	if err != nil {
		return nil, err
	}

	// Load CA certs
	pkcs7Package, err := ioutil.ReadFile(filepath.Clean(dodCACertPackage)) // to placate GOSEC
	if err != nil {
		return nil, err
	}
	caCertPool, err := server.LoadCertPoolFromPkcs7Package(pkcs7Package)
	if err != nil {
		return nil, err
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return &RBSPersonLookup{
		Client: http.Client{Transport: transport},
		Host:   host,
	}, nil
}

func (r RBSPersonLookup) sendGetRequest(url string) ([]byte, error) {
	var data []byte
	resp, err := r.Client.Get(url)
	// Interesting fact: RBS responds 200 OK, not 404 Not Found, if there are no matches
	if err != nil {
		return data, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func buildEdiURL(host string, custNum string, edipi uint64) (string, error) {
	if edipi > 9999999999 {
		return "", errors.New("Invalid EDIPI")
	}

	return fmt.Sprintf(
		"https://%s/appj/rbs/rest/op=edi/customer=%s/schemaName=get_cac_data/schemaVersion=1.0/DOD_EDI_PN_ID=%d",
		host, custNum, edipi), nil
}

func buildWkEmaURL(host string, custNum string, workEmail string) (string, error) {
	if !emailRegex.MatchString(workEmail) {
		return "", errors.New("Invalid e-mail address")
	}

	// e-mail addresses are limited to 80 characters
	l := len(workEmail)
	if l > 80 {
		l = 80
	}

	return fmt.Sprintf(
		"https://%s/appj/rbs/rest/op=wkEma/customer=%s/schemaName=get_cac_data/schemaVersion=1.0/EMA_TX=%s",
		host, custNum, workEmail[:l]), nil
}

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

func parseEdiResponse(data []byte) (*Person, []Personnel, error) {
	rec, err := recordFromResponse(data)
	if err != nil {
		return nil, []Personnel{}, err
	}

	return rec.AdrRecord.Person, rec.AdrRecord.Personnel, nil
}

func parseWkEmaResponse(data []byte) (uint64, *Person, []Personnel, error) {
	rec, err := recordFromResponse(data)
	if err != nil {
		return 0, nil, []Personnel{}, err
	}

	// Not found
	if rec.AdrRecord.WorkEmail == nil {
		return 0, nil, []Personnel{}, nil
	}

	return rec.AdrRecord.WorkEmail.Edipi, rec.AdrRecord.Person, rec.AdrRecord.Personnel, nil
}

func parsePidsResponse(data []byte) (MatchReasonCode, uint64, *Person, []Personnel, error) {
	rec, err := recordFromResponse(data)
	if err != nil {
		return MatchReasonCodeNone, 0, nil, []Personnel{}, err
	}

	reason := rec.AdrRecord.PidsRecord.MtchRsnCd
	if reason == MatchReasonCodeNone {
		return reason, 0, nil, []Personnel{}, nil
	}

	return reason, rec.AdrRecord.PidsRecord.Edipi, rec.AdrRecord.Person, rec.AdrRecord.Personnel, nil
}

func recordFromResponse(data []byte) (Record, error) {
	rec := Record{}
	unmarshalErr := xml.Unmarshal(data, &rec)
	if unmarshalErr != nil {
		// Couldn't unmarshal as a record, try as an RbsError next
		rbsError := RbsError{}
		unmarshalErr = xml.Unmarshal(data, &rbsError)
		if unmarshalErr == nil {
			return rec, &rbsError
		}
		return rec, unmarshalErr
	}
	return rec, nil
}
