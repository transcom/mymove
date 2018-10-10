package iws

import "encoding/xml"

// RbsError is the XML root tag for error replies from DMDC's Identity Web Services: Real-time Broker Service REST API
type RbsError struct {
	XMLName      xml.Name `xml:"RbsError"`
	FaultCode    uint64   `xml:"faultCode"`
	FaultMessage string   `xml:"faultMessage"`
}

func (e *RbsError) Error() string {
	return e.FaultMessage
}
