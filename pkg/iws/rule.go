package iws

// Rule repeats back the ruleset used to query DMDC's Identity Web Services: Real-time Broker Service REST API
type Rule struct {
	Customer      uint32 `xml:"customer"`
	SchemaName    string `xml:"schemaName"`
	SchemaVersion string `xml:"schemaVersion"`
}
