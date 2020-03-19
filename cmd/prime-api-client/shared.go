package main

const (
	// FilenameFlag is the name of the file being passed in
	FilenameFlag string = "filename"
	// ETagFlag is the etag for the mto shipment being updated
	ETagFlag string = "etag"
)

func containsDash(args []string) bool {
	for _, arg := range args {
		if arg == "-" {
			return true
		}
	}
	return false
}
