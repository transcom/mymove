package etag

import (
	"encoding/base64"
	"time"
)

// GenerateEtag creates a string identifier generated from a timestamp
func GenerateEtag(timestamp time.Time) string {
	return base64.StdEncoding.EncodeToString([]byte(timestamp.Format(time.RFC3339Nano)))
}

// DecodeEtag decodes the eTag and parses it into a timestamp
func DecodeEtag(eTag string) (time.Time, error) {
	timeStr, err := base64.StdEncoding.DecodeString(eTag)
	if err != nil {
		return time.Time{}, err
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, string(timeStr))
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
