package iampostgres

// Custom IAM Postgres driver
// - https://stackoverflow.com/questions/56355577/using-database-sql-library-and-fetching-password-from-vault-when-a-new-connectio

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/jmoiron/sqlx"
	pg "github.com/lib/pq"
	"go.uber.org/zap"
)

type pauseFunc func()

var defaultInitialPassword = ""

const defaultPauseDuration = time.Millisecond * 250
const defaultMaxRetries = 120

var defaultPauseFn pauseFunc = func() { time.Sleep(defaultPauseDuration) }

// CustomPostgres is used to set the driverName to the custom postgres driver
const CustomPostgres string = "custompostgres"

type iamPostgresConfig struct {
	useIAM           bool
	maxRetries       int
	pauseFn          pauseFunc
	passHolder       string
	currentIamPass   string
	currentIamTime   time.Time
	currentPassMutex sync.Mutex
	logger           *zap.Logger
	host             string
	port             string
	region           string
	user             string
	passTemplate     string
	creds            *credentials.Credentials
	rus              RDSUtilService
	ticker           *time.Ticker
	shouldQuitChan   chan bool
}

var iamPostgres *iamPostgresConfig

// getCurrentPass gets IAM password if needed and will block till
// valid password is available. It is only called when opening a new
// database connection and thus only attempts to get the credentials a
// limited number of times so that open a connection does not block forever
func (i *iamPostgresConfig) getCurrentPass() (string, time.Time) {
	// Blocks until the password from the dbConnectionDetails has a non blank password
	currentPass := ""
	var currentPassTime time.Time

	counter := 0

	for {
		counter++

		i.currentPassMutex.Lock()
		currentPass = i.currentIamPass
		currentPassTime = i.currentIamTime
		i.currentPassMutex.Unlock()

		if currentPass == "" {
			i.logger.Info(fmt.Sprintf("Wait %d of %d, sleeping for IAM loop to populate RDS credentials.", counter, i.maxRetries))
		} else {
			break
		}
		if counter > i.maxRetries {
			i.logger.Error("Waited 30s for IAM creds to populate and giving up, returning empty password.")
			break
		}

		i.pauseFn()
	}

	return currentPass, currentPassTime
}

func (i *iamPostgresConfig) updateDSN(dsn string) (string, time.Time, error) {
	if !strings.Contains(dsn, i.passHolder) {
		return "", time.Now(), errors.New("DSN does not contain password holder")
	}
	currentPass, currentPassTime := i.getCurrentPass()

	// because this now uses the postgresql user=bob password=secret
	// style connection string, the password does not need to be query
	// escaped
	//
	// use the same escaper logic as in pg.ParseURL in case "'" or "\"
	// appears in the generated token
	escaper := strings.NewReplacer(`'`, `\'`, `\`, `\\`)
	dsn = strings.Replace(dsn, i.passHolder, escaper.Replace(currentPass), 1)
	return dsn, currentPassTime, nil
}

func (i *iamPostgresConfig) generateNewIamPassword() {
	authToken, err := i.rus.GetToken(i.host+":"+i.port, i.region, i.user, i.creds)
	if err != nil {
		i.logger.Error("Error building IAM auth token", zap.Error(err))
	} else {
		i.currentPassMutex.Lock()
		i.currentIamPass = authToken
		i.currentIamTime = time.Now()
		i.currentPassMutex.Unlock()
	}
	hash := sha256.Sum256([]byte(authToken))
	digest := hex.EncodeToString(hash[:])
	i.logger.Info("Successfully generated new IAM auth token",
		zap.String("pwDigest", digest))
}

// Refreshes the RDS IAM on the given interval.
func (i *iamPostgresConfig) refreshRDSIAM() {
	i.logger.Info("Starting refresh of RDS IAM")
	// This for loop immediately runs the first tick then on interval
	// This for loop will run indefinitely until true is passed to the
	// should quit channel.
	for {
		select {
		case <-i.shouldQuitChan:
			// Disable logging here as this goroutine is never
			// shutdown except in tests (as of 2023-01-19).
			// The logging causes a race condition in the tests
			// logger.Warn("Shutting down IAM credential refresh")
			i.ticker.Stop()
			return
		default:
			i.generateNewIamPassword()
			<-i.ticker.C
		}
	}
}

// EnableIAM enables the use of IAM and pulls first credential set as a sanity check
// Note: This method is intended to be non-blocking, so please add any changes to the goroutine
// Note: Ensure the timer is on an interval lower than 15 minutes (AWS RDS IAM auth limit)
func EnableIAM(host string, port string, region string, user string, passTemplate string, creds *credentials.Credentials, rus RDSUtilService, waitDuration time.Duration, logger *zap.Logger, shouldQuitChan chan bool) error {
	if creds == nil {
		return errors.New("IAM Credentials are missing")
	}

	minWaitDuration := 15 * time.Minute

	if waitDuration > minWaitDuration {
		return errors.New("waitDuration too long")
	}

	ticker := time.NewTicker(waitDuration)
	// Lets enable and configure the DSN settings
	iamPostgres = &iamPostgresConfig{
		useIAM:           true,
		maxRetries:       defaultMaxRetries,
		pauseFn:          defaultPauseFn,
		passHolder:       passTemplate,
		currentIamPass:   defaultInitialPassword,
		currentPassMutex: sync.Mutex{},
		logger:           logger,
		host:             host,
		port:             port,
		region:           region,
		user:             user,
		passTemplate:     passTemplate,
		creds:            creds,
		rus:              rus,
		ticker:           ticker,
		shouldQuitChan:   shouldQuitChan,
	}

	// ensure at least one token has been generated
	iamPostgres.generateNewIamPassword()

	// GoRoutine to continually refresh the RDS IAM auth on the given interval.
	go iamPostgres.refreshRDSIAM()
	return nil
}

// rdsPostgresConnector implements the database/sql/driver.Connector
// interface
type rdsPostgresConnector struct {
	dsn    string
	driver driver.Driver
}

func retryableConnect(ctx context.Context, originalDsn string) (driver.Conn, error) {
	useIAM := iamPostgres != nil && iamPostgres.useIAM
	dsn := originalDsn
	var dsnTime time.Time
	var err error
	if useIAM {
		dsn, dsnTime, err = iamPostgres.updateDSN(originalDsn)
		if err != nil {
			zap.L().Error("IAM iampostgres updateDSN failed", zap.Error(err))
			return nil, err
		}
	}

	connector, err := pg.NewConnector(dsn)
	if err != nil {
		zap.L().Error("IAM iampostgres NewConnector failed",
			zap.Any("useIAM", useIAM),
			zap.Any("dsnTime", dsnTime),
			zap.Any("diff", time.Now().Unix()-dsnTime.Unix()),
			zap.Error(err))
		return nil, err
	}

	conn, err := connector.Connect(ctx)
	if err != nil {
		parts := strings.Split(dsn, " ")
		pw := ""
		var b strings.Builder
		for i := range parts {
			if strings.HasPrefix(parts[i], "password") {
				pw = strings.TrimPrefix(parts[i], "password=")
			} else {
				b.WriteString(parts[i] + " ")
			}
		}
		hash := sha256.Sum256([]byte(pw))
		digest := hex.EncodeToString(hash[:])

		zap.L().Error("IAM iampostgres connector.Connect failed",
			zap.Bool("useIAM", useIAM),
			zap.Any("dsn", b.String()),
			zap.Any("dsnTime", dsnTime),
			zap.Any("diff", time.Now().Unix()-dsnTime.Unix()),
			zap.String("pwDigest", digest),
			zap.Error(err))
		return nil, err
	}

	return conn, nil
}

// Connect is called each time a new connection to the database is
// needed, so we can update the RDS IAM auth token immediately before
// connecting to the DB
func (c *rdsPostgresConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := retryableConnect(ctx, c.dsn)
	if err != nil {
		useIAM := iamPostgres != nil && iamPostgres.useIAM
		if useIAM {
			zap.L().Error("IAM iampostgres connect failed, retrying once", zap.Error(err))
			iamPostgres.generateNewIamPassword()
			conn, err = retryableConnect(ctx, c.dsn)
			if err != nil {
				zap.L().Error("IAM iampostgres retry failed", zap.Error(err))
			} else {
				zap.L().Error("IAM iampostgres retry ok")
			}
		}
	}

	return conn, err
}

func (c *rdsPostgresConnector) Driver() driver.Driver {
	return c.driver
}

// RDSPostgresDriver implemements driver.DriverContext
type RDSPostgresDriver struct {
}

// From go's documentation:
//
// If a Driver implements DriverContext, then sql.DB will call
// OpenConnector to obtain a Connector and then invoke that
// Connector's Connect method to obtain each needed connection,
// instead of invoking the Driver's Open method for each connection.
// The two-step sequence allows drivers to parse the name just once
// and also provides access to per-Conn contexts.
//
// Milmove wants this for IAM Authentication so we can update the auth
// token before connect
func (d *RDSPostgresDriver) OpenConnector(dsn string) (driver.Connector, error) {
	// convert to postgres style "username=foo password=bar" style so
	// we don't have to URL encode the password
	pgDsn, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, err
	}
	return &rdsPostgresConnector{pgDsn, &pg.Driver{}}, nil
}

var errNotImplemented = errors.New("Open Not Implemented")

func (d *RDSPostgresDriver) Open(_ string) (driver.Conn, error) {
	// will not be called because of OpenConnector
	return nil, errNotImplemented
}

func init() {
	sql.Register(CustomPostgres, &RDSPostgresDriver{})
	sqlx.BindDriver(CustomPostgres, sqlx.DOLLAR)
}
