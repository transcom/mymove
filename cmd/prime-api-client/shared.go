package main

const (
	// FilenameFlag is the name of the file being passed in
	FilenameFlag string = "filename"
	// ETagFlag is the etag for the mto shipment being updated
	ETagFlag string = "etag"
	// PaymentRequestID is the ID of payment request to use
	PaymentRequestID string = "paymentRequestID"
)

func containsDash(args []string) bool {
	for _, arg := range args {
		if arg == "-" {
			return true
		}
	}
	return false
}
