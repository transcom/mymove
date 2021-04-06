package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/unit"
)

// FormatTimestamp returns a formatted timestamp to display to the TXO
func FormatTimestamp(value time.Time) string {
	valueString := value.Format(TimestampParamFormat)
	return valueString
}

// FormatDate returns a formatted date to display to the TXO
func FormatDate(value time.Time) string {
	valueString := value.Format(DateParamFormat)
	return valueString
}

// FormatCents returns a formatted dollar value, without a $, to display to the TXO
func FormatCents(value unit.Cents) string {
	valueFloat := value.ToDollarFloat()
	return fmt.Sprintf("%.2f", valueFloat)
}
