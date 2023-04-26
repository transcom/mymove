package iampostgres

// Custom IAM Postgres driver
// - https://stackoverflow.com/questions/56355577/using-database-sql-library-and-fetching-password-from-vault-when-a-new-connectio

import (
	"context"
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

var defaultInitialPassword = ""

const defaultPauseDuration = time.Millisecond * 250
const defaultMaxRetries = 120

var defaultPauseFn pauseFunc = func() { time.Sleep(defaultPauseDuration) }

// RDSPostgresDriver wrapper around postgres driver
type RDSPostgresDriver struct {
	*pg.Driver
}

// CustomPostgres is used to set the driverName to the custom postgres driver
const CustomPostgres string = "custompostgres"

type iamPostgresConfig struct {
	useIAM           bool
	maxRetries       int
	pauseFn          pauseFunc
	passHolder       string
	currentIamPass   string
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
func (i *iamPostgresConfig) getCurrentPass() string {
	// Blocks until the password from the dbConnectionDetails has a non blank password
	currentPass := ""

	counter := 0

	for {
		counter++

		i.currentPassMutex.Lock()
		currentPass = i.currentIamPass
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

	return currentPass
}

func (i *iamPostgresConfig) updateDSN(dsn string) (string, error) {
	if !strings.Contains(dsn, i.passHolder) {
		return "", errors.New("DSN does not contain password holder")
	}
	currentPass := i.getCurrentPass()

	dsn = strings.Replace(dsn, i.passHolder, currentPass, 1)
	return dsn, nil
}

func (i *iamPostgresConfig) generateNewIamPassword() {
	authToken, err := i.rus.GetToken(i.host+":"+i.port, i.region, i.user, i.creds)
	if err != nil {
		i.logger.Error("Error building IAM auth token", zap.Error(err))
	} else {
		i.currentPassMutex.Lock()
		i.currentIamPass = url.QueryEscape(authToken)
		i.currentPassMutex.Unlock()
	}
	i.logger.Info("Successfully generated new IAM auth token")
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
		return errors.New("waitDuration too short")
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

	// GoRoutine to continually refresh the RDS IAM auth on the given interval.
	go iamPostgres.refreshRDSIAM()
	return nil
}

type DriverConnWrapper struct {
	driver.Conn
	driver.Pinger
	driver.SessionResetter
	driver.Validator
	driver.ConnBeginTx
	driver.ConnPrepareContext
	driver.ExecerContext
	driver.QueryerContext
	// these interfaces are deprecated, so do not implement them
	// putting this info here in case someone in the future wonders
	//
	// driver.Execer
	// driver.Queryer
	driver.Tx
	createdAt time.Time
}

var errWrap = errors.New("Cannot wrap driver conn")

func (dcw DriverConnWrapper) Ping(ctx context.Context) error {

	err := dcw.Pinger.Ping(ctx)
	if err != nil {
		zap.L().Info("IAM iampostgres ping failed",
			zap.Any("createdAt", dcw.createdAt))
	}
	return err
}

func (dcw DriverConnWrapper) Close() error {
	err := dcw.Conn.Close()
	if err != nil {
		zap.L().Error("IAM iampostgres Close failed",
			zap.Any("createdAt", dcw.createdAt),
			zap.Error(err))
	}
	return err
}

func (dcw DriverConnWrapper) ResetSession(ctx context.Context) error {
	err := dcw.SessionResetter.ResetSession(ctx)
	if err != nil {
		zap.L().Error("IAM iampostgres reset session failed",
			zap.Any("createdAt", dcw.createdAt),
			zap.Error(err))
	}
	return err
}

func (d RDSPostgresDriver) newDriverConnWrapper(dsn string) (DriverConnWrapper, error) {
	conn, err := d.Driver.Open(dsn)
	if err != nil {
		return DriverConnWrapper{}, err
	}

	pinger, ok := conn.(driver.Pinger)
	if !ok {
		return DriverConnWrapper{}, errWrap
	}
	sessionResetter, ok := conn.(driver.SessionResetter)
	if !ok {
		return DriverConnWrapper{}, errWrap
	}
	validator, ok := conn.(driver.Validator)
	if !ok {
		return DriverConnWrapper{}, errWrap
	}
	connBeginTx, ok := conn.(driver.ConnBeginTx)
	if !ok {
		return DriverConnWrapper{}, errWrap
	}
	connPrepareContext, ok := conn.(driver.ConnPrepareContext)
	if !ok {
		return DriverConnWrapper{}, errWrap
	}
	execerContext, ok := conn.(driver.ExecerContext)
	if !ok {
		return DriverConnWrapper{}, errWrap
	}
	queryerContext, ok := conn.(driver.QueryerContext)
	if !ok {
		return DriverConnWrapper{}, errWrap
	}
	tx, ok := conn.(driver.Tx)
	if !ok {
		return DriverConnWrapper{}, errWrap
	}

	return DriverConnWrapper{
		Conn:               conn,
		Pinger:             pinger,
		SessionResetter:    sessionResetter,
		Validator:          validator,
		ConnBeginTx:        connBeginTx,
		ConnPrepareContext: connPrepareContext,
		ExecerContext:      execerContext,
		QueryerContext:     queryerContext,
		Tx:                 tx,
		createdAt:          time.Now(),
	}, nil
}

// Open wrapper around postgres Open func
func (d RDSPostgresDriver) Open(dsn string) (_ driver.Conn, err error) {
	useIAM := iamPostgres != nil && iamPostgres.useIAM
	// default to global logger
	if useIAM {
		dsn, err = iamPostgres.updateDSN(dsn)
		if err != nil {
			return nil, err
		}
	}
	return d.newDriverConnWrapper(dsn)
}

func init() {
	sql.Register(CustomPostgres, &RDSPostgresDriver{&pg.Driver{}})
	sqlx.BindDriver(CustomPostgres, sqlx.DOLLAR)
}
