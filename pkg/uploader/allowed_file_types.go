package uploader

// AllowedFileTypes contains a list of content types
type AllowedFileTypes []string

var (
	// AllowedTypesServiceMember are the content types we allow service members to upload for orders
	AllowedTypesServiceMember AllowedFileTypes = []string{"image/jpeg", "image/png", "application/pdf", "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}

	// AllowedTypesPPMDocuments are the content types we allow service members to upload for PPM shipment closeout documentation
	AllowedTypesPPMDocuments AllowedFileTypes = []string{"image/jpeg", "image/png", "application/pdf", "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"}

	// AllowedTypesPaymentRequest are the content types we allow prime to upload
	AllowedTypesPaymentRequest AllowedFileTypes = []string{"image/jpeg", "image/png", "application/pdf"}

	// AllowedTypesText accepts text files
	AllowedTypesText AllowedFileTypes = []string{"text/plain", "text/plain; charset=utf-8"}

	// AllowedTypesPDF accepts PDF files
	AllowedTypesPDF AllowedFileTypes = []string{"application/pdf"}

	// AllowedTypesPDFImages accepts PDF files and images
	AllowedTypesPDFImages AllowedFileTypes = []string{"image/jpeg", "image/png", "application/pdf"}

	// AllowedTypesAny accepts any file type
	AllowedTypesAny AllowedFileTypes = []string{"*"}
)

// Contains checks to see if the provided file type is acceptable
func (aft AllowedFileTypes) Contains(fileType string) bool {
	if len(aft) == 1 && aft[0] == "*" {
		return true
	}
	for _, fType := range aft {
		if fType == fileType {
			return true
		}
	}
	return false
}

// Contents returns the allowed types as a slice of strings
func (aft AllowedFileTypes) Contents() []string {
	return aft
}
