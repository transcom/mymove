package ghcrateengine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/unit"
)

// FormatTimestamp returns a formatted timestamp to display to the TXO
func FormatTimestamp(value time.Time) string {
	return value.Format(TimestampParamFormat)
}

// FormatDate returns a formatted date to display to the TXO
func FormatDate(value time.Time) string {
	return value.Format(DateParamFormat)
}

// FormatCents returns a formatted dollar value, without a $, to display to the TXO
func FormatCents(value unit.Cents) string {
	valueFloat := value.ToDollarFloat()
	return fmt.Sprintf("%.2f", valueFloat)
}

// FormatBool returns a formatted boolean value to display to the TXO
func FormatBool(value bool) string {
	return strconv.FormatBool(value)
}

// FormatFloat returns a formatted float value to display to the TXO
func FormatFloat(value float64, precision int) string {
	return strconv.FormatFloat(value, 'f', precision, 64)
}
