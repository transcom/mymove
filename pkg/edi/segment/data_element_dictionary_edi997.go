package edisegment

import (
	"fmt"
)

// Helpful Data Element Dictionary definitions found in the 997 Function Acknowledgment spec
// for use with Syncada application

// 717 Transaction SEt Acknowledgement Code
// Code indicating accept or reject condition based on the syntax editing of the transaction set
var transactionSetAckCode717 = map[string]string{
	"A": "Accepted",
	"E": "Accepted, But Errors Were Noted.",
	"M": "Rejected, Message Authentication Code (MAC) Failed",
	"P": "Partially Accepted, At Least One Transaction Set Was Rejected",
	"R": "Rejected",
	"W": "Rejected, Assurance Failed Validity Tests",
	"X": "Rejected, Content After Decryption Could Not Be Analyzed",
}

var transactionSetAckCode717Accepted = []string{
	"A",
	"E",
}

type dataElement717 struct {
}

func (de dataElement717) Accepted(code string) bool {
	for _, ackCode := range transactionSetAckCode717Accepted {
		if code == ackCode {
			return true
		}
	}
	return false
}

func (de dataElement717) Description(code string) (string, error) {
	description, ok := transactionSetAckCode717[code]
	if ok {
		return description, nil
	}
	return "", fmt.Errorf("code %s not found", code)
}
