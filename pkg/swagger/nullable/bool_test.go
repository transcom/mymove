package nullable

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestBool_UnmarshalJSON(t *testing.T) {
	valueTrue := true
	valueFalse := false

	tests := []struct {
		name      string
		expectErr bool
		buf       *bytes.Buffer
		expect    Bool
	}{
		{
			name: "null value",
			buf:  bytes.NewBufferString(`{"value":null}`),
			expect: Bool{
				Present: true,
				Value:   nil,
			},
			expectErr: false,
		},
		{
			name: "valid value true",
			buf:  bytes.NewBufferString(`{"value":true}`),
			expect: Bool{
				Present: true,
				Value:   &valueTrue,
			},
			expectErr: false,
		},
		{
			name: "valid value false",
			buf:  bytes.NewBufferString(`{"value":false}`),
			expect: Bool{
				Present: true,
				Value:   &valueFalse,
			},
			expectErr: false,
		},
		{
			name:      "invalid value",
			buf:       bytes.NewBufferString(`{"value":0}`),
			expect:    Bool{},
			expectErr: true,
		},
		{
			name:      "empty",
			buf:       bytes.NewBufferString(`{}`),
			expect:    Bool{},
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := struct {
				Value Bool `json:"value"`
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
