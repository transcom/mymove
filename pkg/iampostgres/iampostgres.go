package iampostgres

// Custom IAM Postgres driver
// - https://stackoverflow.com/questions/56355577/using-database-sql-library-and-fetching-password-from-vault-when-a-new-connectio

import (
	"net/url"
	"strings"
	"time"

	"database/sql"
	"database/sql/driver"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"go.uber.org/zap"

	pg "github.com/lib/pq"
)

type config struct {
	useIAM         bool
	passHolder     string
	currentIamPass string
}

var iamConfig = config{false, "", ""}

// RDSPostgresDriver wrapper around postgres ddiver
type RDSPostgresDriver struct {
	*pg.Driver
}

// GetCurrentPass a helper function to get IAM password if needed outside of driver
func GetCurrentPass() string {
	// Blocks until the password from the dbConnectionDetails has a non blank password
	for {
		if iamConfig.currentIamPass == "" {
			time.Sleep(time.Millisecond * 250)
		} else {
			break
		}
	}

	return iamConfig.currentIamPass
}

func updateDSN(dsn string) string {
	dsn = strings.Replace(dsn, iamConfig.passHolder, GetCurrentPass(), 1)
	return dsn
}

// EnableIAM enables the use of IAM and pulls first credential set as a sanity check
func EnableIAM(host string, port string, region string, user string, passTemplate string, creds *credentials.Credentials, logger Logger) {
	// Lets enable and configure the DSN settings
	iamConfig.useIAM = true
	iamConfig.passHolder = passTemplate

	// GoRoutine to continually refresh the RDS IAM auth on a 10m interval.
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		// This for loop immediately runs the first tick then on interval
		for ; true; <-ticker.C {
			if creds == nil {
				logger.Error("IAM Credentials are missing")
				return
			}
			logger.Info("Using IAM Authentication")
			authToken, err := rdsutils.BuildAuthToken(host+":"+port, region, user, creds)
			if err != nil {
				logger.Error("Error building auth token", zap.Error(err))
				return
			}
			iamConfig.currentIamPass = url.QueryEscape(authToken)
			logger.Info("Successfully generated new IAM token")
		}
	}()

}

// Open wrapper around postgres Open func
func (d RDSPostgresDriver) Open(dsn string) (_ driver.Conn, err error) {
	if iamConfig.useIAM == true {
		dsn = updateDSN(dsn)
	}

	return d.Driver.Open(dsn)
}

func init() {
	sql.Register("iampostgres", &RDSPostgresDriver{&pg.Driver{}})
}
