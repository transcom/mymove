package iampostgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"go.uber.org/zap"
)

type RDSUTest struct {
	passes []string
}

func (r RDSUTest) GetToken(endpoint string, region string, user string, iamcreds *credentials.Credentials) (string, error) {
	pass := r.passes[0]
	r.passes = append(r.passes[:0], r.passes[1:]...)

	return pass, nil
}

func TestEnableIamNilCreds(t *testing.T) {
	assert := assert.New(t)

	rdsu := RDSUTest{}
	logger, _ := zap.NewProduction()

	tmr := time.NewTicker(1 * time.Second)

	shouldQuitChan := make(chan bool)

	EnableIAM("server", "8080", "us-east-1", "dbuser", "***",
		nil,
		rdsu,
		tmr,
		logger,
		shouldQuitChan)
	time.Sleep(2 * time.Second)

	iamConfig.currentPassMutex.Lock()
	t.Logf("Current password: %s", iamConfig.currentIamPass)
	assert.Equal(iamConfig.currentIamPass, "")
	iamConfig.currentPassMutex.Unlock()
	tmr.Stop()

}

func TestGetCurrentPassword(t *testing.T) {
	assert := assert.New(t)

	rdsu := RDSUTest{}
	rdsu.passes = append(rdsu.passes, "abc")
	logger, _ := zap.NewProduction()
	iamConfig.currentIamPass = "" // ensure iamConfig is in new state

	tmr := time.NewTicker(2 * time.Second)

	shouldQuitChan := make(chan bool)

	EnableIAM("server", "8080", "us-east-1", "dbuser", "***",
		credentials.NewStaticCredentials("id", "pass", "token"),
		rdsu,
		tmr,
		logger,
		shouldQuitChan)

	// this should block for ~ 250ms and then continue
	currentPass := GetCurrentPass()
	assert.Equal(currentPass, "abc")
	shouldQuitChan <- true

	tmr.Stop()

}

func TestGetCurrentPasswordFail(t *testing.T) {
	// This tests when the timeout is hit

	assert := assert.New(t)

	rdsu := RDSUTest{}
	rdsu.passes = append(rdsu.passes, "") // set mocked pass to empty to simulate failed cred generation
	logger, _ := zap.NewProduction()
	iamConfig.currentIamPass = ""

	tmr := time.NewTicker(1 * time.Second)

	shouldQuitChan := make(chan bool)

	EnableIAM("server", "8080", "us-east-1", "dbuser", "***",
		credentials.NewStaticCredentials("id", "pass", "token"),
		rdsu,
		tmr,
		logger,
		shouldQuitChan)

	// this should block for 30s then return empty string
	currentPass := GetCurrentPass()
	assert.Equal(currentPass, "")
	shouldQuitChan <- true
	tmr.Stop()

}

func TestEnableIAMNormal(t *testing.T) {
	assert := assert.New(t)

	testData := []string{"abc", "123", "xyz", "999"}
	rdsu := RDSUTest{}
	rdsu.passes = append(rdsu.passes, testData...)
	logger, _ := zap.NewProduction()

	tmr := time.NewTicker(2 * time.Second)

	shouldQuitChan := make(chan bool)

	// We use 2 second timer since that builds in a buffer so we have stable tests
	EnableIAM("server", "8080", "us-east-1", "dbuser", "***",
		credentials.NewStaticCredentials("id", "pass", "token"),
		rdsu,
		tmr,
		logger,
		shouldQuitChan)

	lenTestData := len(testData) - 1
	counter := 0
	expectedPass := testData[counter]
	for {
		// Poll for the password change
		time.Sleep(250 * time.Millisecond)

		// If the current pass does not match the last known password
		// then we should check the next password in the slice
		pass := GetCurrentPass()
		if pass != expectedPass {
			t.Logf("Counter %d, Current/expected password: %s, %s", counter, pass, expectedPass)
			counter++
			// Don't allow looking items that don't exist
			if counter > lenTestData {
				break
			}
			expectedPass = testData[counter]
		}

		// Check that the expected password is what we want
		t.Logf("Counter %d, Current/expected password: %s, %s", counter, pass, expectedPass)
		assert.Equal(expectedPass, pass)

		// Once the end has been reached then end the test
		if counter == lenTestData && expectedPass == testData[lenTestData] {
			break
		}
	}

	// Check that all the passwords have been checked
	assert.Equal(3, counter)
	shouldQuitChan <- true
	tmr.Stop()
}

func TestUpdateDSN(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		name        string
		passHolder  string
		pass        string
		dsn         string
		expectedDSN string
	}{
		{"simple test",
			"***",
			"PASSWORD",
			"postgres://db_user:***@host:5432/dev_db?sslmode=verify-full",
			"postgres://db_user:PASSWORD@host:5432/dev_db?sslmode=verify-full"},
		{"different password holder",
			"!!!",
			"PASSWORD",
			"postgres://db_user:!!!@host:5432/dev_db?sslmode=verify-full",
			"postgres://db_user:PASSWORD@host:5432/dev_db?sslmode=verify-full"},
		{"multiple occurrence of password holder, only first occurrence replaced",
			"***",
			"PASSWORD",
			"postgres://db_user:***@host:5432/dev_db?sslmode=***",
			"postgres://db_user:PASSWORD@host:5432/dev_db?sslmode=***"},
	}

	for _, tt := range tests {
		t.Logf("Running scenario: %s", tt.name)
		iamConfig.currentIamPass = tt.pass
		iamConfig.passHolder = tt.passHolder
		dsn, err := updateDSN(tt.dsn)
		assert.Equal(dsn, tt.expectedDSN)
		assert.Nil(err)
	}
}
