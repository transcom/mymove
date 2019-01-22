package iws

// Identifier repeats back the terms used to query DMDC's Identity Web Services: Real-time Broker Service REST API
type Identifier struct {
	Edipi *uint64 `xml:"DOD_EDI_PN_ID,omitempty"`
	Pids  *Person `xml:"pids,omitempty"`
}
