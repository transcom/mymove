package nullable

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestInt_UnmarshalJSON(t *testing.T) {
	valueOne := int64(1)

	tests := []struct {
		name      string
		buf       *bytes.Buffer
		expect    Int
		expectErr bool
	}{
		{
			name: "null value",
			buf:  bytes.NewBufferString(`{"value":null}`),
			expect: Int{
				Present: true,
				Value:   nil,
			},
			expectErr: false,
		},
		{
			name: "valid value",
			buf:  bytes.NewBufferString(`{"value":1}`),
			expect: Int{
				Present: true,
				Value:   &valueOne,
			},
			expectErr: false,
		},
		{
			name:      "invalid value",
			buf:       bytes.NewBufferString(`{"value":"definitely not an integer"}`),
			expect:    Int{},
			expectErr: true,
		},
		{
			name:      "empty",
			buf:       bytes.NewBufferString(`{}`),
			expect:    Int{},
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := struct {
				Value Int `json:"value"`
			}{}
			err := json.Unmarshal(tt.buf.Bytes(), &str)
			if tt.expectErr && err == nil {
				t.Fatalf("expected error")
			}
			if err != nil {
				if !tt.expectErr {
					t.Fatalf("unexpected unmarshaling error: %s", err)
				}
			}

			if !tt.expectErr {
				got := str.Value
				valuesMatch := false
				if tt.expect.Value == nil && got.Value == nil {
					valuesMatch = true
				}
				if tt.expect.Value != nil && got.Value != nil {
					valuesMatch = *tt.expect.Value == *got.Value
				}

				if got.Present != tt.expect.Present || !valuesMatch {
					t.Errorf("expected value to be %#v got %#v", tt.expect, got)
				}
			}
		})
	}
}
