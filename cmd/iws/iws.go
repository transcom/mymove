package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/iws"
)

const (
	edipiFlag = "edipi"
	ssnFlag   = "ssn"
	lastFlag  = "last"
	firstFlag = "first"
	emailFlag = "email"
)

func initFlags(flag *pflag.FlagSet) {

	cli.InitIWSFlags(flag)
	cli.InitCertFlags(flag)

	flag.Uint64(edipiFlag, 0, "10-digit EDIPI to look up (op=edi)")
	flag.String(ssnFlag, "", "9-digit SSN to look up, without dashes (op=pids)")
	flag.String(lastFlag, "", "Last Name to look up (op=pids)")
	flag.String(firstFlag, "", "First Name to look up (op=pids) [optional]")
	flag.String(emailFlag, "", "Work e-mail address to look up (op=wkEma)")

	// Don't sort flags
	flag.SortFlags = false
}

func checkConfig(v *viper.Viper) error {

	if err := cli.CheckIWS(v); err != nil {
		return err
	}

	if err := cli.CheckCert(v); err != nil {
		return err
	}

	return nil
}

func main() {

	flag := pflag.CommandLine
	initFlags(flag)
	err := flag.Parse(os.Args[1:])
	if err != nil {
		log.Fatal("Could not parse flags", err)
	}

	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		log.Fatal("Could not find flags", err)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	err = checkConfig(v)
	if err != nil {
		log.Fatal("Cannot validate config", err)
	}

	// Load client cert
	rbs, err := iws.NewRBSPersonLookup(
		v.GetString(cli.IWSRBSHostFlag),
		v.GetStringSlice(cli.DoDCAPackageFlag),
		v.GetString(cli.MoveMilDoDTLSCertFlag),
		v.GetString(cli.MoveMilDoDTLSKeyFlag))

	if err != nil {
		log.Fatal("Cannot initialize rbs person lookup", err)
	}

	var retcode int

	edipi := v.GetUint64(edipiFlag)
	ssn := v.GetString(ssnFlag)
	lastName := v.GetString(lastFlag)
	firstName := v.GetString(firstFlag)
	workEmail := v.GetString(emailFlag)

	if edipi != 0 {
		retcode = edi(*rbs, edipi)
	} else if ssn != "" {
		retcode = pids(*rbs, ssn, lastName, firstName)
	} else if workEmail != "" {
		retcode = wkEma(*rbs, workEmail)
	} else {
		flag.Usage()
		retcode = -1
	}

	os.Exit(retcode)
}

func edi(rbs iws.RBSPersonLookup, edipi uint64) int {
	fmt.Printf("Identity Web Services: Real-Time Broker Service (REST)\nHost: %s\nOperation: edi\nEDIPI: %d\n", rbs.Host, edipi)
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

func pids(rbs iws.RBSPersonLookup, ssn string, lastName string, firstName string) int {
	fmt.Printf("Identity Web Services: Real-Time Broker Service (REST)\nHost: %s\nOperation: pids-P\nSSN: %s\nLast Name: %s\nFirst Name: %s\n", rbs.Host, ssn, lastName, firstName)

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

func wkEma(rbs iws.RBSPersonLookup, email string) int {
	fmt.Printf("Identity Web Services: Real-Time Broker Service (REST)\nHost: %s\nOperation: wkEma\nWork E-mail: %s\n", rbs.Host, email)
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
