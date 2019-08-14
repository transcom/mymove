package adminapi

import (
	"github.com/transcom/mymove/pkg/models"
)

func payloadForPhoneLines(officePhoneLines models.OfficePhoneLines) []string {
	var phoneLines []string
	for _, phoneLine := range officePhoneLines {
		if phoneLine.Type == "voice" {
			phoneLines = append(phoneLines, phoneLine.Number)
		}
	}

	return phoneLines
}
