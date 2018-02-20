package auth

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/markbates/pop"
)

func TestRegisterProvider(t *testing.T) {
	fmt.Println("hit test register provider")
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
