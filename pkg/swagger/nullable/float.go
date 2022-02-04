package nullable

import (
	"bytes"
	"encoding/json"
)

// Float represents a float that may be null or not
// present in json at all.
type Float struct {
	Present bool // Present is true if key is present in json
	Value   *float64
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (f *Float) UnmarshalJSON(data []byte) error {
	f.Present = true

	if bytes.Equal(data, null) {
		return nil
	}

	if err := json.Unmarshal(data, &f.Value); err != nil {
		return err
	}

	return nil
}
