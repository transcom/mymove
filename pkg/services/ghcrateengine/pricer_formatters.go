package ghcrateengine

import (
	"time"
)

// FormatTimestamp returns a formatted timestamp to display to the TXO
func FormatTimestamp(value time.Time) (string, error) {
	valueString := value.Format(TimestampParamFormat)
	return valueString, nil
}

// FormatDate returns a formatted date to display to the TXO
func FormatDate(value time.Time) (string, error) {
	valueString := value.Format(DateParamFormat)
	return valueString, nil
}
