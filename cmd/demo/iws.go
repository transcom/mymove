package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/namsral/flag"
	"github.com/transcom/mymove/pkg/iws"
	"github.com/transcom/mymove/pkg/server"
)

func main() {
	host := flag.String("iws_rbs_host", "", "hostname of the IWS RBS environment")
	custNum := flag.String("iws_rbs_cust_num", "", "customer number to present when connecting to IWS RBS")
	dodCaCertPackage := flag.String("dod_ca_package", "", "Path to PKCS7 package containing all DoD Certificate Authority certificates.")
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
		os.Exit(1)
	}

	// Load CA certs
	pkcs7Package, err := ioutil.ReadFile(*dodCaCertPackage)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	caCertPool, err := server.LoadCertPoolFromPkcs7Package(pkcs7Package)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	client := http.Client{Transport: transport}

	rbs := iws.RealTimeBrokerService{
		Client:  client,
		Host:    *host,
		CustNum: *custNum,
	}
	var retcode int

	if *edipi != 0 {
		retcode = edi(rbs, *edipi)
	} else if *ssn != "" {
		retcode = pids(rbs, *ssn, *lastName, *firstName)
	} else if *workEmail != "" {
		retcode = wkEma(rbs, *workEmail)
	} else {
		flag.Usage()
		retcode = -1
	}

	os.Exit(retcode)
}

func edi(rbs iws.RealTimeBrokerService, edipi uint64) int {
	fmt.Printf("Identity Web Services: Real-Time Broker Service (REST)\nHost: %s\nCustomer Number: %s\nOperation: edi\nEDIPI: %d\n", rbs.Host, rbs.CustNum, edipi)
	person, personnel, err := rbs.GetPersonUsingEDIPI(edipi)

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

func pids(rbs iws.RealTimeBrokerService, ssn string, lastName string, firstName string) int {
	fmt.Printf("Identity Web Services: Real-Time Broker Service (REST)\nHost: %s\nCustomer Number: %s\nOperation: pids-P\nSSN: %s\nLast Name: %s\nFirst Name: %s\n", rbs.Host, rbs.CustNum, ssn, lastName, firstName)

	params := iws.GetPersonUsingSSNParams{
		Ssn:       ssn,
		LastName:  lastName,
		FirstName: firstName,
	}
	reason, edipi, person, personnel, err := rbs.GetPersonUsingSSN(params)

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

func wkEma(rbs iws.RealTimeBrokerService, email string) int {
	fmt.Printf("Identity Web Services: Real-Time Broker Service (REST)\nHost: %s\nCustomer Number: %s\nOperation: wkEma\nWork E-mail: %s\n", rbs.Host, rbs.CustNum, email)
	edipi, person, personnel, err := rbs.GetPersonUsingWorkEmail(email)

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
