package auth

import (
	"testing"
)

func TestGenerateNonce(t *testing.T) {
	nonce := generateNonce()

	if (nonce == "") || (len(nonce) < 1) {
		t.Error("No nonce was returned.")
	}
}
