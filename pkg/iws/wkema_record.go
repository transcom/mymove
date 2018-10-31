package iws

// WkEmaRecord contains the EDIPI and work e-mail address of an individual
type WkEmaRecord struct {
	Edipi uint64 `xml:"DOD_EDI_PN_ID"`
	Email string `xml:"EMA_TX"`
}
