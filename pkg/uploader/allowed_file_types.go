package uploader

// AllowedFileTypes contains a list of content types
type AllowedFileTypes []string

var (
	// AllowedTypesServiceMember are the content types we allow service members to upload
	AllowedTypesServiceMember AllowedFileTypes = []string{"image/jpeg", "image/png", "application/pdf"}

	// AllowedTypesText accepts text files
	AllowedTypesText AllowedFileTypes = []string{"text/plain", "text/plain; charset=utf-8"}

	// AllowedTypesPDF accepts PDF files
	AllowedTypesPDF AllowedFileTypes = []string{"application/pdf"}

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
