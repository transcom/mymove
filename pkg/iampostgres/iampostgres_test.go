package iampostgres

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type RDSUTest struct {
	passes []string
}

func (r *RDSUTest) GetToken(endpoint string, region string, user string, iamcreds *credentials.Credentials) (string, error) {
	if len(r.passes) == 0 {
		return "", errors.New("no passwords to rotate")
	}

	// Rotate the slice: first item goes to back of slice
	pass := r.passes[0]
	r.passes = append(r.passes[1:], pass)

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
		&rdsu,
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
		&rdsu,
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
		&rdsu,
		tmr,
		logger,
		shouldQuitChan)

	// this should block for 30s then return empty string
	currentPass := GetCurrentPass()
	assert.Equal(currentPass, "")
	shouldQuitChan <- true
	tmr.Stop()

}

/*
Test to see that the EnableIAM method is working.
It should be sufficient to check that the EnableIAM method is working by:
Setting an initial password, running EnableIAM, and then verifying that the
currentPassword is no longer the initial password, but is instead the password
that EnableIAM is cycling through.

While this wont completely eliminate the flakiness (since the test still needs
to wait for the password to change), it should significantly reduce the
flakiness by only having 1 sleep and 1 password to swtich to.

The test was written in this manner since testing to see that multiple
passwords are being cycled through will inherently make the test timing
dependent and thus, flaky.
i.e. If a new password gets cycled every 1 minute, the test would need to
sleep for a minute and then check to see that the current password is correct
one in the sequence. If anything falls out of sync, then the tests will fail
since what is expected will be dsycned from what is retrieved.
*/
func TestEnableIAMNormal(t *testing.T) {
	assert := assert.New(t)

	// Cycle through 1 password so that the test doesn't have to get too exact
	// about when the password changed and what password it is on.
	testData := []string{"abc"}
	rdsu := RDSUTest{}
	rdsu.passes = append(rdsu.passes, testData...)
	logger, _ := zap.NewProduction()
	// Set the current password to something not in the above list of passwords
	// to cycle through.
	iamConfig.currentIamPass = "123"

	tmr := time.NewTicker(1 * time.Second)

	shouldQuitChan := make(chan bool)

	// Confirm that the password got set to what we initially set it to.
	pass := GetCurrentPass()
	assert.Equal("123", pass)

	// Start cycling through the list of passwords.
	EnableIAM("server", "8080", "us-east-1", "dbuser", "***",
		credentials.NewStaticCredentials("id", "pass", "token"),
		&rdsu,
		tmr,
		logger,
		shouldQuitChan)

	// The sleep time should be greater than how often the password will cycle
	// so that the next time the password is fetched, it will have changed.
	time.Sleep(2 * time.Second)

	// Confirm that the password has changed (it's no longer the initial
	// password) to the 1 password being cycled through.
	pass = GetCurrentPass()
	assert.Equal("abc", pass)

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
