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

type config struct {
	useIAM           bool
	passHolder       string
	currentIamPass   string
	currentPassMutex sync.Mutex
	logger           *zap.Logger
}

var iamConfig = config{false, "", "", sync.Mutex{}, nil}

// RDSPostgresDriver wrapper around postgres driver
type RDSPostgresDriver struct {
	*pg.Driver
}

// CustomPostgres is used to set the driverName to the custom postgres driver
const CustomPostgres string = "custompostgres"

// GetCurrentPass gets IAM password if needed and will block till valid password is available
func GetCurrentPass() string {
	// Blocks until the password from the dbConnectionDetails has a non blank password
	currentPass := ""

	counter := 0
	maxCount := 120 // pauses for 30s

	for {
		counter++

		iamConfig.currentPassMutex.Lock()
		currentPass = iamConfig.currentIamPass
		iamConfig.currentPassMutex.Unlock()

		if currentPass == "" {
			iamConfig.logger.Info(fmt.Sprintf("Wait %d of %d, sleeping for 250ms for IAM loop to populate RDS credentials.", counter, maxCount))
		} else {
			break
		}
		if counter > maxCount {
			iamConfig.logger.Error("Waited 30s for IAM creds to populate and giving up, returning empty password.")
			break
		}

		time.Sleep(time.Millisecond * 250)

	}

	return currentPass
}

func updateDSN(dsn string) (string, error) {
	if !strings.Contains(dsn, iamConfig.passHolder) {
		return "", errors.New("DSN does not contain password holder")
	}

	dsn = strings.Replace(dsn, iamConfig.passHolder, GetCurrentPass(), 1)
	return dsn, nil
}

// Refreshes the RDS IAM on the given interval.
func refreshRDSIAM(host string, port string, region string, user string, creds *credentials.Credentials, rus RDSUtilService, ticker *time.Ticker, logger *zap.Logger, errorMessagesChan chan error, shouldQuitChan chan bool) {
	logger.Info("Starting refresh of RDS IAM")
	// This for loop immediately runs the first tick then on interval
	// This for loop will run indefinitely until it either errors or true is
	// passed to the should quit channel.
	for {
		select {
		case <-shouldQuitChan:
			close(errorMessagesChan)
			return
		default:
			if creds == nil {
				logger.Error("IAM Credentials are missing")
				errorMessagesChan <- errors.New("IAM Credientials are missing")
				close(errorMessagesChan)
				return
			}
			logger.Info("Using IAM Authentication")
			authToken, err := rus.GetToken(host+":"+port, region, user, creds)
			if err != nil {
				logger.Error("Error building auth token", zap.Error(err))
				errorMessagesChan <- fmt.Errorf("Error building auth token %v", err)
				close(errorMessagesChan)
				return
			}

			iamConfig.currentPassMutex.Lock()
			iamConfig.currentIamPass = url.QueryEscape(authToken)
			iamConfig.currentPassMutex.Unlock()
			logger.Info("Successfully generated new IAM token")
			<-ticker.C
		}
	}
}

// EnableIAM enables the use of IAM and pulls first credential set as a sanity check
// Note: This method is intended to be non-blocking, so please add any changes to the goroutine
// Note: Ensure the timer is on an interval lower than 15 minutes (AWS RDS IAM auth limit)
func EnableIAM(host string, port string, region string, user string, passTemplate string, creds *credentials.Credentials, rus RDSUtilService, ticker *time.Ticker, logger *zap.Logger, shouldQuitChan chan bool) {
	// Lets enable and configure the DSN settings
	iamConfig.useIAM = true
	iamConfig.passHolder = passTemplate
	iamConfig.logger = logger

	errorMessagesChan := make(chan error)

	// GoRoutine to continually refresh the RDS IAM auth on the given interval.
	go refreshRDSIAM(host, port, region, user, creds, rus, ticker, logger, errorMessagesChan, shouldQuitChan)

	go logEnableIAMFailed(logger, errorMessagesChan)
}

func logEnableIAMFailed(logger *zap.Logger, errorMessagesChan chan error) {
	errorMessages := <-errorMessagesChan

	if errorMessages != nil {
		logger.Error("Refreshing RDS IAM failed")
	}
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
