package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/transcom/mymove/pkg/edi/gex"
)

// Call this from command line with go run cmd/send_to_gex/main.go -edi <filepath>
func main() {
	ediFile := flag.String("edi", "", "The filepath to an edi file to send to GEX")
	transactionName := flag.String("transactionName", "test", "The required name sent in the url of the gex api request")
	flag.Parse()
	if *ediFile == "" {
		log.Fatal("Usage: go run cmd/send_to_gex/main.go  --edi <edi filepath> --transactionName <name>")
	}

	file, err := os.Open(*ediFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	edi, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	ediString := string(edi[:])
	// make sure edi ends in new line
	ediString = strings.TrimSpace(ediString) + "\n"

	fmt.Println(ediString)
	resp, err := gex.SendInvoiceToGex(ediString, *transactionName)

	fmt.Println("Sending to GEX. . .")
	fmt.Printf("status code: %v, error: %v \n", resp.StatusCode, err)
}
