package main

import (
	"fmt"
	"log"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/apiclient"
	"github.com/transcom/mymove/pkg/gen/apiclient/shipments"
)

func main() {
	// create the transport
	transport := httptransport.New("tsplocal:3000", "/api/v1", nil)

	// create the API client, with the transport
	client := apiclient.New(transport, strfmt.Default)

	// make the request to get all items
	params := &shipments.IndexShipmentsParams{}
	params.WithTimeout(time.Second * 10)

	resp, err := client.Shipments.IndexShipments(params)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", resp.Payload)
}
