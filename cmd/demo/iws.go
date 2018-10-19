package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/namsral/flag"
	"github.com/transcom/mymove/pkg/iws"
)

func main() {
	host := flag.String("iws_rbs_host", "", "hostname of the IWS RBS environment")
	custNum := flag.String("iws_rbs_cust_num", "", "customer number to present when connecting to IWS RBS")
	moveMilDODCACert := flag.String("move_mil_dod_ca_cert", "", "The DoD CA certificate used to sign the move.mil TLS certificates.")
	moveMilDODTLSCert := flag.String("move_mil_dod_tls_cert", "", "The DoD-signed TLS certificate for various move.mil services.")
	moveMilDODTLSKey := flag.String("move_mil_dod_tls_key", "", "The private key for the DoD-signed TLS certificate for various move.mil services.")
	edipi := flag.Uint64("edipi", 0, "10-digit EDIPI to look up (op=edi)")
	ssn := flag.String("ssn", "", "9-digit SSN to look up, without dashes (op=pids)")
	lastName := flag.String("last", "", "Last Name to look up (op=pids)")
	firstName := flag.String("first", "", "First Name to look up (op=pids) [optional]")
	workEmail := flag.String("email", "", "Work e-mail address to look up (op=wkEma)")

	flag.Parse()

	// Load client cert
	cert, err := tls.X509KeyPair([]byte(*moveMilDODTLSCert), []byte(*moveMilDODTLSKey))
	if err != nil {
		log.Fatal(err)
	}

	// Load CA certs
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(*moveMilDODCACert))

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	client := http.Client{Transport: transport}
	var retcode int

	if *edipi != 0 {
		retcode = edi(client, *host, *custNum, *edipi)
	} else if *ssn != "" {
		retcode = pids(client, *host, *custNum, *ssn, *lastName, *firstName)
	} else if *workEmail != "" {
		retcode = wkEma(client, *host, *custNum, *workEmail)
	} else {
		flag.Usage()
		retcode = -1
	}

	os.Exit(retcode)
}

func edi(client http.Client, host string, custNum string, edipi uint64) int {
	fmt.Printf("Identity Web Services: Real-Time Broker Service (REST)\nHost: %s\nCustomer Number: %s\nOperation: edi\nEDIPI: %d\n", host, custNum, edipi)
	person, personnel, err := iws.GetPersonUsingEDIPI(client, host, custNum, edipi)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return -1
	}

	if person == nil {
		fmt.Println("No match")
		return 0
	}

	fmt.Printf("Person: %+v\nPersonnel: %+v\n", person, personnel)
	return 0
}

func pids(client http.Client, host string, custNum string, ssn string, lastName string, firstName string) int {
	fmt.Printf("Identity Web Services: Real-Time Broker Service (REST)\nHost: %s\nCustomer Number: %s\nOperation: pids-P\nSSN: %s\nLast Name: %s\nFirst Name: %s\n", host, custNum, ssn, lastName, firstName)

	params := iws.GetPersonUsingSSNParams{
		Ssn:       ssn,
		LastName:  lastName,
		FirstName: firstName,
	}
	reason, edipi, person, personnel, err := iws.GetPersonUsingSSN(client, host, custNum, params)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return -1
	}

	if reason == iws.MatchReasonCodeNone {
		fmt.Println("No match")
		return 0
	}

	fmt.Printf("Match Reason: %s\nEDIPI: %d\nPerson: %+v\nPersonnel: %+v\n", reason, edipi, person, personnel)
	return 0
}

func wkEma(client http.Client, host string, custNum string, email string) int {
	fmt.Printf("Identity Web Services: Real-Time Broker Service (REST)\nHost: %s\nCustomer Number: %s\nOperation: wkEma\nWork E-mail: %s\n", host, custNum, email)
	edipi, person, personnel, err := iws.GetPersonUsingWorkEmail(client, host, custNum, email)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return -1
	}

	if edipi == 0 {
		fmt.Println("No match")
		return 0
	}

	fmt.Printf("EDIPI: %d\nPerson: %+v\nPersonnel: %+v\n", edipi, person, personnel)
	return 0
}
