package auth

import (
	"log"
	"os"
	"testing"

	"github.com/markbates/pop"
)

func TestGenerateNonce(t *testing.T) {
	nonce := generateNonce()

	if (nonce == "") || (len(nonce) < 1) {
		t.Error("No nonce was returned.")
	}
}

var dbConnection *pop.Connection

func setupDBConnection() {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	conn, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	dbConnection = conn
}

func TestMain(m *testing.M) {
	setupDBConnection()
	os.Exit(m.Run())
}
