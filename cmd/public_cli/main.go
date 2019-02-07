package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/apiclient"
	"github.com/transcom/mymove/pkg/gen/apiclient/shipments"
)

func main() {
	// Parse the url for the api
	u, err := url.Parse("http://tsplocal:3000/api/v1")
	if err != nil {
		log.Fatal(err)
	}

	// create the transport
	transport := httptransport.New(u.Host, u.Path, nil)

	// Set the session token
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	var cookies []*http.Cookie
	jwt := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjI1MzM3MDgwODAwMCwiU2Vzc2lvblZhbHVlIjp7IkFwcGxpY2F0aW9uTmFtZSI6IlRTUCIsIkhvc3RuYW1lIjoidHNwbG9jYWwiLCJJRFRva2VuIjoiZGV2bG9jYWwiLCJVc2VySUQiOiJiNmJhODVhYS0xZmVkLTQzNmItYWE5YS1kNDAwMzhhMjA2OGUiLCJFbWFpbCI6IjIwMTkwMTMxMTkxMjMwQGV4YW1wbGUuY29tIiwiRmlyc3ROYW1lIjoiVGVzdHkiLCJNaWRkbGUiOiIiLCJMYXN0TmFtZSI6Ik1jVGVzdGVyIiwiU2VydmljZU1lbWJlcklEIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwiT2ZmaWNlVXNlcklEIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIiwiVHNwVXNlcklEIjoiZWY0MmIxZjctYmQwMy00NDdmLWJmNjEtYzZmMDIwZmNjODc3IiwiRmVhdHVyZXMiOm51bGx9fQ.clBpIhYocazRIvUOguzGwzDGiNO4_tDs7eI982EkDwcOyZ92eP_azzhpn6ILNATPYXy393KKYkZusRRaC-CpHNL-ylCbmnOc4ZZeb0N__luKZ7asAkFFxacPyITrYqI9ZomA2xjG6AdkjScj5xT2qf8w5Uj1P2uNx1FHb2c3wqVtf344sD9if0EOdKULxk0taBkdIn0a4sw23INVFxCIz1ncfLuO6vdbu01tq_iI_XZI6qsUdJcgPmaG49ZbJnXQwhTqJFohbcFf2QRij1BjFURdzO7tjBfZsyB70IY2R1ogrmBZh3Ptf_vtlmBkco4vjiiRmvFLdJlG2wGApXgD1g"
	sessionCookie := &http.Cookie{
		Name:  "tsp_session_token",
		Value: jwt,
	}
	cookies = append(cookies, sessionCookie)
	jar.SetCookies(u, cookies)

	// Attach the cookies to the transport
	transport.Jar = jar

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
