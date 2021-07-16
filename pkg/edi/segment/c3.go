package edisegment

import "fmt"

// C3 represents the C3 EDI segment (Currency)
type C3 struct {
	// http://www.iso.org/iso/en/prods-services/popstds/currencycodeslist.html is the URL
	// referenced in the 858 standard, but that URL is not valid
	// https://en.wikipedia.org/wiki/ISO_4217 seems to be a good list to reference
	CurrencyCodeC301 string `validate:"omitempty,max=3"`
	ExchangeRate     string `validate:"omitempty"`
	CurrencyCodeC303 string `validate:"omitempty,max=3"`
	CurrencyCodeC304 string `validate:"omitempty,max=3"`
}

// StringArray converts C3 to an array of strings
func (s *C3) StringArray() []string {
	return []string{
		"C3",
		s.CurrencyCodeC301,
		s.ExchangeRate,
		s.CurrencyCodeC303,
		s.CurrencyCodeC304,
	}
}

// Parse parses an C3 string that's split into an array into the C3 struct
func (s *C3) Parse(elements []string) error {
	expectedMinNumElements := 1
	expectedMaxNumElements := 4
	numElements := len(elements)
	if numElements < expectedMinNumElements || numElements > expectedMaxNumElements {
		return fmt.Errorf("C3: Wrong number of fields, expected min %d and max %d, got %d", expectedMinNumElements, expectedMaxNumElements, len(elements))
	}

	s.CurrencyCodeC301 = elements[0]
	if numElements > 1 {
		s.ExchangeRate = elements[1]
	}
	if numElements > 2 {
		s.CurrencyCodeC303 = elements[2]
	}
	if numElements > 3 {
		s.CurrencyCodeC304 = elements[3]
	}

	return nil
}
