package iampostgres

// Custom IAM Postgres driver
// - https://stackoverflow.com/questions/56355577/using-database-sql-library-and-fetching-password-from-vault-when-a-new-connectio

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/jmoiron/sqlx"
	pg "github.com/lib/pq"
	"go.uber.org/zap"
)

type pauseFunc func()

type config struct {
	useIAM           bool
	maxRetries       int
	pauseFn          pauseFunc
	passHolder       string
	currentIamPass   string
	currentPassMutex sync.Mutex
	logger           *zap.Logger
}

const defaultPauseDuration = time.Millisecond * 250
const defaultMaxRetries = 120

var defaultPauseFn pauseFunc = func() { time.Sleep(defaultPauseDuration) }
var iamConfig = config{
	useIAM:           false,
	maxRetries:       defaultMaxRetries,
	pauseFn:          defaultPauseFn,
	currentPassMutex: sync.Mutex{},
}

// RDSPostgresDriver wrapper around postgres driver
type RDSPostgresDriver struct {
	*pg.Driver
}

// CustomPostgres is used to set the driverName to the custom postgres driver
const CustomPostgres string = "custompostgres"

// getCurrentPass gets IAM password if needed and will block till
// valid password is available. It is only called when opening a new
// database connection and thus only attempts to get the credentials a
// limited number of times so that open a connection does not block forever
func getCurrentPass() string {
	// Blocks until the password from the dbConnectionDetails has a non blank password
	currentPass := ""

	counter := 0

	for {
		counter++

		iamConfig.currentPassMutex.Lock()
		currentPass = iamConfig.currentIamPass
		iamConfig.currentPassMutex.Unlock()

		if currentPass == "" {
			iamConfig.logger.Info(fmt.Sprintf("Wait %d of %d, sleeping for IAM loop to populate RDS credentials.", counter, iamConfig.maxRetries))
		} else {
			break
		}
		if counter > iamConfig.maxRetries {
			iamConfig.logger.Error("Waited 30s for IAM creds to populate and giving up, returning empty password.")
			break
		}

		iamConfig.pauseFn()
	}

	return currentPass
}

func updateDSN(dsn string) (string, error) {
	if !strings.Contains(dsn, iamConfig.passHolder) {
		return "", errors.New("DSN does not contain password holder")
	}

	dsn = strings.Replace(dsn, iamConfig.passHolder, getCurrentPass(), 1)
	return dsn, nil
}

// Refreshes the RDS IAM on the given interval.
func refreshRDSIAM(host string, port string, region string, user string, creds *credentials.Credentials, rus RDSUtilService, ticker *time.Ticker, logger *zap.Logger, shouldQuitChan chan bool) {
	logger.Info("Starting refresh of RDS IAM")
	// This for loop immediately runs the first tick then on interval
	// This for loop will run indefinitely until true is passed to the
	// should quit channel.
	for {
		select {
		case <-shouldQuitChan:
			logger.Warn("Shutting down IAM credential refresh")
			return
		default:
			authToken, err := rus.GetToken(host+":"+port, region, user, creds)
			if err != nil {
				logger.Error("Error building IAM auth token", zap.Error(err))
			} else {
				iamConfig.currentPassMutex.Lock()
				iamConfig.currentIamPass = url.QueryEscape(authToken)
				iamConfig.currentPassMutex.Unlock()
				logger.Info("Successfully generated new IAM auth token")
			}
			<-ticker.C
		}
	}
}

// EnableIAM enables the use of IAM and pulls first credential set as a sanity check
// Note: This method is intended to be non-blocking, so please add any changes to the goroutine
// Note: Ensure the timer is on an interval lower than 15 minutes (AWS RDS IAM auth limit)
func EnableIAM(host string, port string, region string, user string, passTemplate string, creds *credentials.Credentials, rus RDSUtilService, ticker *time.Ticker, logger *zap.Logger, shouldQuitChan chan bool) error {
	if creds == nil {
		return errors.New("IAM Credentials are missing")
	}
	// Lets enable and configure the DSN settings
	iamConfig.useIAM = true
	iamConfig.passHolder = passTemplate
	iamConfig.logger = logger

	// GoRoutine to continually refresh the RDS IAM auth on the given interval.
	go refreshRDSIAM(host, port, region, user, creds, rus, ticker, logger, shouldQuitChan)
	return nil
}

// Open wrapper around postgres Open func
func (d RDSPostgresDriver) Open(dsn string) (_ driver.Conn, err error) {
	if iamConfig.useIAM {
		dsn, err = updateDSN(dsn)
		if err != nil {
			return nil, err
		}
	}

	return d.Driver.Open(dsn)
}

func init() {
	sql.Register(CustomPostgres, &RDSPostgresDriver{&pg.Driver{}})
	sqlx.BindDriver(CustomPostgres, sqlx.DOLLAR)
}
