package iws

import "encoding/xml"

// Record is the XML root tag for responses from DMDC's Identity Web Services: Real-time Broker Service REST API.
type Record struct {
	XMLName    xml.Name   `xml:"record"`
	Rule       Rule       `xml:"rule"`
	Identifier Identifier `xml:"identifier"`
	AdrRecord  AdrRecord  `xml:"adrRecord"`
}
