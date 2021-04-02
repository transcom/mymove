package models

// EDIType represents types of EDI Responses
type EDIType string

const (
	// EDIType810 captures enum value "810"
	EDIType810 EDIType = "810"
	// EDIType824 captures enum value "824"
	EDIType824 EDIType = "824"
	// EDIType858 captures enum value "858"
	EDIType858 EDIType = "858"
	// EDIType997 captures enum value "997"
	EDIType997 EDIType = "997"
)

var allowedEDITypes = []string{
	string(EDIType810),
	string(EDIType824),
	string(EDIType858),
	string(EDIType997),
}

// String returns a string representation of the admin role
func (e EDIType) String() string {
	return string(e)
}
