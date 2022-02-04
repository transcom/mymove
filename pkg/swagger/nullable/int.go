package nullable

import (
	"bytes"
	"encoding/json"
)

// Int represents an int that may be null or not
// present in json at all.
type Int struct {
	Present bool // Present is true if key is present in json
	Value   *int64
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (i *Int) UnmarshalJSON(data []byte) error {
	i.Present = true

	if bytes.Equal(data, null) {
		return nil
	}

	if err := json.Unmarshal(data, &i.Value); err != nil {
		return err
	}

	return nil
}
