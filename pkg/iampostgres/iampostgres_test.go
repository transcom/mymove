package iampostgres

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
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
	logger := zaptest.NewLogger(t)

	tmr := time.NewTicker(1 * time.Second)

	shouldQuitChan := make(chan bool)

	err := EnableIAM("server", "8080", "us-east-1", "dbuser", "***",
		nil,
		&rdsu,
		tmr,
		logger,
		shouldQuitChan)
	assert.Error(err, "Enable IAM with nil creds?")
}

func TestGetCurrentPassword(t *testing.T) {
	assert := assert.New(t)

	rdsu := RDSUTest{}
	// We need to put a blank token first to avoid race conditions.
	// See more below
	rdsu.passes = append(rdsu.passes, "")
	rdsu.passes = append(rdsu.passes, "abc")
	logger := zaptest.NewLogger(t)
	iamConfig.currentIamPass = "" // ensure iamConfig is in new state

	// use Millisecond so the tests run faster
	tickerDuration := 1 * time.Millisecond
	pauseCounter := 0
	iamConfig.pauseFn = func() {
		pauseCounter++
		// make sure we wait at least one ticker duration for the
		// password to be updated
		time.Sleep(tickerDuration)
	}

	tmr := time.NewTicker(tickerDuration)

	shouldQuitChan := make(chan bool)

	err := EnableIAM("server", "8080", "us-east-1", "dbuser", "***",
		credentials.NewStaticCredentials("id", "pass", "token"),
		&rdsu,
		tmr,
		logger,
		shouldQuitChan)
	assert.Nil(err, "Enable IAM error")

	// this should block once and then continue
	currentPass := getCurrentPass()
	assert.Equal(currentPass, "abc")
	shouldQuitChan <- true
	tmr.Stop()

	// This test wants to ensure that getCurrentPass does not return
	// an empty string and keeps looping until it does gets a non
	// empty string.
	//
	// If refreshRDSIAM runs first (called from EnableIAM), it will
	// set currentIamPass to the next token (which is set to the empty
	// string up above). Then getCurrentPass runs and it will see that
	// currentIamPass is the empty string and then loop, running
	// iamConfig.pauseFn (and incrementing pauseCounter). Then, the
	// refreshRDSIAM loop will run, getting the next token (set to
	// "abc" up above). Finally, the getCurrentPass loop will run
	// again (after the pauseFn finishes) and get the current
	// password.
	//
	// If getCurrentPass runs first, it will see the currentIamPass is
	// the empty string and run iamConfig.pauseFn (and increment
	// pauseCounter). Then it is the same as the refreshRDSIAM
	// scenario above.
	//
	// If the refreshRDSIAM go routine runs before getCurrentPass,
	// there would be only one pause by getCurrentPass. If
	// getCurrentPass runs before refreshRDSIAM, there would be two
	// pauses.
	//
	// If the first token is a non empty string, then when
	// refreshRDSIAM runs first it will set currentIamPass and then
	// getCurrentPassword will not loop at all. That means we are not
	// testing what we want to test, the looping property of
	// getCurrentPassword

	assert.True(1 <= pauseCounter, "expected pauseCounter to be greater than 1, was "+fmt.Sprintf("%d", pauseCounter))
}

func TestGetCurrentPasswordFail(t *testing.T) {
	// This tests when the timeout is hit

	assert := assert.New(t)

	rdsu := RDSUTest{}
	rdsu.passes = append(rdsu.passes, "") // set mocked pass to empty to simulate failed cred generation
	logger := zaptest.NewLogger(t)
	iamConfig.currentIamPass = ""

	// use Millisecond so the tests run faster
	tickerDuration := 1 * time.Millisecond
	pauseCounter := 0
	iamConfig.pauseFn = func() {
		pauseCounter++
		// make sure we wait at least one ticker duration for the
		// password to be updated
		time.Sleep(tickerDuration)
	}

	tmr := time.NewTicker(tickerDuration)

	shouldQuitChan := make(chan bool)

	err := EnableIAM("server", "8080", "us-east-1", "dbuser", "***",
		credentials.NewStaticCredentials("id", "pass", "token"),
		&rdsu,
		tmr,
		logger,
		shouldQuitChan)
	assert.Nil(err, "Enable IAM error")

	// this should block until maxRetries, then return empty string
	currentPass := getCurrentPass()
	assert.Equal(currentPass, "")
	shouldQuitChan <- true
	tmr.Stop()
	assert.Equal(int(iamConfig.maxRetries), pauseCounter)
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
	logger := zaptest.NewLogger(t)
	// Set the current password to something not in the above list of passwords
	// to cycle through.
	iamConfig.currentIamPass = "123"
	pauseCounter := 0
	iamConfig.pauseFn = func() { pauseCounter++ }

	// use Millisecond so the tests run faster
	tmr := time.NewTicker(1 * time.Millisecond)

	shouldQuitChan := make(chan bool)

	// Confirm that the password got set to what we initially set it to.
	pass := getCurrentPass()
	assert.Equal("123", pass)

	// Start cycling through the list of passwords.
	err := EnableIAM("server", "8080", "us-east-1", "dbuser", "***",
		credentials.NewStaticCredentials("id", "pass", "token"),
		&rdsu,
		tmr,
		logger,
		shouldQuitChan)
	assert.Nil(err, "Enable IAM error")

	// The sleep time should be greater than how often the password will cycle
	// so that the next time the password is fetched, it will have changed.
	// use Millisecond so the tests run faster
	time.Sleep(2 * time.Millisecond)

	// Confirm that the password has changed (it's no longer the initial
	// password) to the 1 password being cycled through.
	pass = getCurrentPass()
	assert.Equal("abc", pass)

	shouldQuitChan <- true
	tmr.Stop()

	// in this case, should never pause
	assert.Equal(0, pauseCounter)
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
