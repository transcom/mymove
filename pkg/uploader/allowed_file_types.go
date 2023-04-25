package uploader

const (
	// FileTypeJPEG is the content type for JPEG images
	FileTypeJPEG = "image/jpeg"
	// FileTypePNG is the content type for PNG images
	FileTypePNG = "image/png"
	// FileTypePDF is the content type for PDF documents
	FileTypePDF = "application/pdf"
	// FileTypeText is the content type for text files
	FileTypeText = "text/plain"
	// FileTypeTextUTF8 is the content type for text files with UTF-8 encoding
	FileTypeTextUTF8 = "text/plain; charset=utf-8"
	// FileTypeExcel is the content type for Excel files
	FileTypeExcel = "application/vnd.ms-excel"
	// FileTypeExcelXLSX is the content type for Excel xlsx files
	FileTypeExcelXLSX = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
)

const AccessiblePDFFormat = "PDF/A-1a"

// AllowedFileTypes contains a list of content types
type AllowedFileTypes []string

var (
	// AllowedTypesServiceMember are the content types we allow service members to upload for orders
	AllowedTypesServiceMember AllowedFileTypes = []string{FileTypeJPEG, FileTypePNG, FileTypePDF}

	// AllowedTypesPPMDocuments are the content types we allow service members to upload for PPM shipment closeout documentation
	AllowedTypesPPMDocuments AllowedFileTypes = []string{FileTypeJPEG, FileTypePNG, FileTypePDF, FileTypeExcel, FileTypeExcelXLSX}

	// AllowedTypesPaymentRequest are the content types we allow prime to upload
	AllowedTypesPaymentRequest AllowedFileTypes = []string{FileTypeJPEG, FileTypePNG, FileTypePDF}

	// AllowedTypesText accepts text files
	AllowedTypesText AllowedFileTypes = []string{FileTypeText, FileTypeTextUTF8}

	// AllowedTypesPDF accepts PDF files
	AllowedTypesPDF AllowedFileTypes = []string{FileTypePDF}

	// AllowedTypesPDFImages accepts PDF files and images
	AllowedTypesPDFImages AllowedFileTypes = []string{FileTypeJPEG, FileTypePNG, FileTypePDF}

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
