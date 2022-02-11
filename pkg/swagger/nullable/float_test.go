package nullable

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestFloat_UnmarshalJSON(t *testing.T) {
	valueTwelve := 12.0

	tests := []struct {
		name      string
		buf       *bytes.Buffer
		expect    Float
		expectErr bool
	}{
		{
			name: "null value",
			buf:  bytes.NewBufferString(`{"value":null}`),
			expect: Float{
				Present: true,
				Value:   nil,
			},
			expectErr: false,
		},
		{
			name: "valid value",
			buf:  bytes.NewBufferString(`{"value":12.0}`),
			expect: Float{
				Present: true,
				Value:   &valueTwelve,
			},
			expectErr: false,
		},
		{
			name:      "invalid value",
			buf:       bytes.NewBufferString(`{"value":"definitely not a float"}`),
			expect:    Float{},
			expectErr: true,
		},
		{
			name:      "empty",
			buf:       bytes.NewBufferString(`{}`),
			expect:    Float{},
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := struct {
				Value Float `json:"value"`
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
