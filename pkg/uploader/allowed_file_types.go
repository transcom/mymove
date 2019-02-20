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

// AllowsAny returned true if any file type is acceptable
func (aft AllowedFileTypes) AllowsAny() bool {
	if len(aft) == 1 && aft[0] == "*" {
		return true
	}
	return false
}
