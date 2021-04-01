package ghcrateengine

import (
	"time"
)

func FormatTimestamp(value time.Time) (string, error) {
	valueString := value.Format(TimestampParamFormat)
	return valueString, nil
}

func FormatDate(value time.Time) (string, error) {
	valueString := value.Format(DateParamFormat)
	return valueString, nil
}

