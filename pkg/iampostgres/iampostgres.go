package iampostgres

// Custom IAM Postgres driver
// - https://stackoverflow.com/questions/56355577/using-database-sql-library-and-fetching-password-from-vault-when-a-new-connectio

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"sync"
	"time"

	"database/sql"
	"database/sql/driver"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"go.uber.org/zap"

	pg "github.com/lib/pq"
)

type config struct {
	useIAM           bool
	passHolder       string
	currentIamPass   string
	currentPassMutex sync.Mutex
	logger           Logger
}

var iamConfig = config{false, "", "", sync.Mutex{}, nil}

// RDSPostgresDriver wrapper around postgres driver
type RDSPostgresDriver struct {
	*pg.Driver
}

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

// EnableIAM enables the use of IAM and pulls first credential set as a sanity check
// Note: This method is intended to be non-blocking, so please add any changes to the goroutine
// Note: Ensure the timer is on an interval lower than 15 minutes (AWS RDS IAM auth limit)
func EnableIAM(host string, port string, region string, user string, passTemplate string, creds *credentials.Credentials, rus RDSUtilService, ticker *time.Ticker, logger Logger) {
	// Lets enable and configure the DSN settings
	iamConfig.useIAM = true
	iamConfig.passHolder = passTemplate
	iamConfig.logger = logger

	// GoRoutine to continually refresh the RDS IAM auth on a 10m interval.
	go func() {

		// Add some entropy to this value so all instances don't fire at the same time
		minDur := int64(100)
		maxDur := int64(5000)
		randDur := big.NewInt(maxDur - minDur)
		randInt, err := rand.Int(rand.Reader, randDur)
		if err != nil {
			logger.Error("Error building auth token", zap.Error(err))
			return
		}
		wait := time.Millisecond * time.Duration(randInt.Int64()+minDur)
		logger.Info(fmt.Sprintf("Waiting %v before enabling IAM access for entropy", wait))
		time.Sleep(wait)

		// This for loop immediately runs the first tick then on interval
		for ; true; <-ticker.C {
			if creds == nil {
				logger.Error("IAM Credentials are missing")
				return
			}
			logger.Info("Using IAM Authentication")
			authToken, err := rus.GetToken(host+":"+port, region, user, creds)
			if err != nil {
				logger.Error("Error building auth token", zap.Error(err))
				return
			}

			iamConfig.currentPassMutex.Lock()
			iamConfig.currentIamPass = url.QueryEscape(authToken)
			iamConfig.currentPassMutex.Unlock()
			logger.Info("Successfully generated new IAM token")
		}
	}()
}

// Open wrapper around postgres Open func
func (d RDSPostgresDriver) Open(dsn string) (_ driver.Conn, err error) {
	if iamConfig.useIAM == true {
		dsn, err = updateDSN(dsn)
		if err != nil {
			return nil, err
		}
	}

	return d.Driver.Open(dsn)
}

func init() {
	sql.Register("custompostgres", &RDSPostgresDriver{&pg.Driver{}})
}
