package adminapi

import (
	"github.com/transcom/mymove/pkg/models"
)

func payloadForPhoneLines(OfficePhoneLines models.OfficePhoneLines) []string {
	var phoneLines []string
	for _, phoneLine := range OfficePhoneLines {
		if phoneLine.Type == "voice" {
			phoneLines = append(phoneLines, phoneLine.Number)
		}
	}

	return phoneLines
}
