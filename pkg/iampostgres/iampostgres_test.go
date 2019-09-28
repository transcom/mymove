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

	EnableIAM("server", "8080", "us-east-1", "dbuser", "***",
		nil,
		rdsu,
		time.NewTicker(1*time.Second),
		logger)
	time.Sleep(2 * time.Second)
	t.Logf("Current password: %s", iamConfig.currentIamPass)
	assert.Equal(iamConfig.currentIamPass, "")

}

func TestGetCurrentPass(t *testing.T) {
	assert := assert.New(t)

	rdsu := RDSUTest{}
	rdsu.passes = append(rdsu.passes, "abc")
	logger, _ := zap.NewProduction()
	iamConfig.currentIamPass = "" // ensure iamConfig is in new state

	EnableIAM("server", "8080", "us-east-1", "dbuser", "***",
		credentials.NewStaticCredentials("id", "pass", "token"),
		rdsu,
		time.NewTicker(2*time.Second),
		logger)

	// this should block for ~ 250ms and then continue
	currentPass := GetCurrentPass()
	assert.Equal(currentPass, "abc")

}

func TestEnableIAMNormal(t *testing.T) {
	assert := assert.New(t)

	rdsu := RDSUTest{}
	rdsu.passes = append(rdsu.passes, "abc")
	rdsu.passes = append(rdsu.passes, "123")
	rdsu.passes = append(rdsu.passes, "xyz")
	rdsu.passes = append(rdsu.passes, "999")
	logger, _ := zap.NewProduction()

	// We use 2 second timer since that builds in a buffer so we have stable tests
	EnableIAM("server", "8080", "us-east-1", "dbuser", "***",
		credentials.NewStaticCredentials("id", "pass", "token"),
		rdsu,
		time.NewTicker(2*time.Second),
		logger)

	time.Sleep(time.Second)
	assert.Equal(iamConfig.currentIamPass, "abc")
	t.Logf("Current password: %s", iamConfig.currentIamPass)

	time.Sleep(2 * time.Second)
	t.Logf("Current password: %s", iamConfig.currentIamPass)
	assert.Equal(iamConfig.currentIamPass, "123")

	time.Sleep(2 * time.Second)
	t.Logf("Current password: %s", iamConfig.currentIamPass)
	assert.Equal(iamConfig.currentIamPass, "xyz")

	time.Sleep(2 * time.Second)
	t.Logf("Current password: %s", iamConfig.currentIamPass)
	assert.Equal(iamConfig.currentIamPass, "999")

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
