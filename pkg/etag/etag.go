package etag

import (
	"encoding/base64"
	"time"
)

// GenerateEtag creates a string identifier generated from a timestamp
func GenerateEtag(timestamp time.Time) string {
	return base64.StdEncoding.EncodeToString([]byte(timestamp.Format(time.RFC3339Nano)))
}
