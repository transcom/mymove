package iws

// RbsError is the XML root tag for error replies from DMDC's Identity Web Services: Real-time Broker Service REST API
type RbsError struct {
	FaultCode    uint64 `xml:"faultCode"`
	FaultMessage string `xml:"faultMessage"`
}
