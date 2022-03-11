package nullable

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestString_UnmarshalJSON(t *testing.T) {
	goodString := "string"
	tests := []struct {
		name      string
		buf       *bytes.Buffer
		expect    String
		expectErr bool
	}{
		{
			name: "null value",
			buf:  bytes.NewBufferString(`{"value":null}`),
			expect: String{
				Present: true,
			},
			expectErr: false,
		},
		{
			name: "valid value",
			buf:  bytes.NewBufferString(`{"value":"string"}`),
			expect: String{
				Present: true,
				Value:   &goodString,
			},
			expectErr: false,
		},
		{
			name:      "invalid value",
			buf:       bytes.NewBufferString(`{"value":5}`),
			expect:    String{},
			expectErr: true,
		},
		{
			name:      "empty",
			buf:       bytes.NewBufferString(`{}`),
			expect:    String{},
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := struct {
				Value String `json:"value"`
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
