package cli

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/gobuffalo/pop/v6"
	"github.com/jmoiron/sqlx"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/XSAM/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	iampg "github.com/transcom/mymove/pkg/iampostgres"
)

const (
	// DbDebugFlag is the DB Debug flag
	DbDebugFlag string = "db-debug"
	// DbEnvFlag is the DB environment flag
	DbEnvFlag string = "db-env"
	// DbNameFlag is the DB name flag
	DbNameFlag string = "db-name"
	// DbHostFlag is the DB host flag
	DbHostFlag string = "db-host"
	// DbPortFlag is the DB port flag
	DbPortFlag string = "db-port"
	// DbUserFlag is the DB user flag
	DbUserFlag string = "db-user"
	// DbPasswordFlag is the DB password flag
	DbPasswordFlag string = "db-password"
	// DbPoolFlag is the DB pool flag
	DbPoolFlag string = "db-pool"
	// DbIdlePoolFlag is the DB idle pool flag
	DbIdlePoolFlag string = "db-idle-pool"
	// DbSSLModeFlag is the DB SSL Mode flag
	DbSSLModeFlag string = "db-ssl-mode"
	// DbSSLRootCertFlag is the DB SSL Root Cert flag
	DbSSLRootCertFlag string = "db-ssl-root-cert"
	// DbIamFlag is the DB IAM flag
	DbIamFlag string = "db-iam"
	// DbIamRoleFlag is the DB IAM Role flag
	DbIamRoleFlag string = "db-iam-role"
	// DbRegionFlag is the DB Region flag
	DbRegionFlag string = "db-region"
	// DbUseInstrumentedDriverFlag indicates if additional db
	// instrumentation should be done
	DbInstrumentedFlag = "db-instrumented"

	// DbEnvContainer is the Container DB Env name
	DbEnvContainer string = "container"
	// DbEnvTest is the Test DB Env name
	DbEnvTest string = "test"
	// DbEnvDevelopment is the Development DB Env name
	DbEnvDevelopment string = "development"

	// DbNameTest The name of the test database
	DbNameTest string = "test_db"

	// SSLModeDisable is the disable SSL Mode
	SSLModeDisable string = "disable"
	// SSLModeAllow is the allow SSL Mode
	SSLModeAllow string = "allow"
	// SSLModePrefer is the prefer SSL Mode
	SSLModePrefer string = "prefer"
	// SSLModeRequire is the require SSL Mode
	SSLModeRequire string = "require"
	// SSLModeVerifyCA is the verify-ca SSL Mode
	SSLModeVerifyCA string = "verify-ca"
	// SSLModeVerifyFull is the verify-full SSL Mode
	SSLModeVerifyFull string = "verify-full"

	// awsRdsT3SmallMaxConnections is the max connections to an RDS T3
	// Small instance
	//
	// The T3 small instance has 2 GB
	// https://aws.amazon.com/rds/instance-types/
	//
	// These docs say we can calculate the max connections
	// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Limits.html
	//
	// If correct it is
	//
	// LEAST({DBInstanceClassMemory/9531392}, 5000)
	//
	// DBInstanceClassMemory = 2147483648
	// so 2147483648 / 9531392 = 225.3 which is less than 5000
	//
	// we deploy two containers for the AWS service, so divide that in
	// half
	// 225 / 2 =~ 110
	awsRdsT3SmallMaxConnections = 110
	// DbPoolDefault is the default db pool connections
	DbPoolDefault = awsRdsT3SmallMaxConnections
	// DbIdlePoolDefault is the default db idle pool connections
	DbIdlePoolDefault = 2
	// DbPoolMax is the upper limit the db pool can use for connections which constrains the user input
	DbPoolMax int = awsRdsT3SmallMaxConnections
)

// The dependency https://github.com/lib/pq only supports a limited subset of SSL Modes and returns the error:
// pq: unsupported sslmode \"prefer\"; only \"require\" (default), \"verify-full\", \"verify-ca\", and \"disable\" supported
// - https://www.postgresql.org/docs/10/libpq-ssl.html
var allSSLModes = []string{
	SSLModeDisable,
	// SSLModeAllow,
	// SSLModePrefer,
	SSLModeRequire,
	SSLModeVerifyCA,
	SSLModeVerifyFull,
}

var containerSSLModes = []string{
	SSLModeRequire,
	SSLModeVerifyCA,
	SSLModeVerifyFull,
}

var allDbEnvs = []string{
	DbEnvContainer,
	DbEnvTest,
	DbEnvDevelopment,
}

type errInvalidDbPool struct {
	DbPool int
}

func (e *errInvalidDbPool) Error() string {
	return fmt.Sprintf("invalid db pool of %d. Pool must be greater than 0 and less than or equal to %d", e.DbPool, DbPoolMax)
}

type errInvalidDbIdlePool struct {
	DbPool     int
	DbIdlePool int
}

func (e *errInvalidDbIdlePool) Error() string {
	return fmt.Sprintf("invalid db idle pool of %d. Pool must be greater than 0 and less than or equal to %d", e.DbIdlePool, e.DbPool)
}

type errInvalidDbEnv struct {
	Value  string
	DbEnvs []string
}

func (e *errInvalidDbEnv) Error() string {
	return fmt.Sprintf("invalid db env %s, must be one of: ", e.Value) + strings.Join(e.DbEnvs, ", ")
}

type errInvalidSSLMode struct {
	Mode  string
	Modes []string
}

func (e *errInvalidSSLMode) Error() string {
	return fmt.Sprintf("invalid ssl mode %s, must be one of: "+strings.Join(e.Modes, ", "), e.Mode)
}

// InitDatabaseFlags initializes DB command line flags
func InitDatabaseFlags(flag *pflag.FlagSet) {
	flag.String(DbEnvFlag, DbEnvDevelopment, "database environment: "+strings.Join(allDbEnvs, ", "))
	flag.String(DbNameFlag, "dev_db", "Database Name")
	flag.String(DbHostFlag, "localhost", "Database Hostname")
	flag.Int(DbPortFlag, 5432, "Database Port")
	flag.String(DbUserFlag, "crud", "Database Username")
	flag.String(DbPasswordFlag, "", "Database Password")
	flag.Int(DbPoolFlag, DbPoolDefault, "Database Pool or max DB connections")
	flag.Int(DbIdlePoolFlag, DbIdlePoolDefault, "Database Idle Pool or max DB idle connections")
	flag.String(DbSSLModeFlag, SSLModeDisable, "Database SSL Mode: "+strings.Join(allSSLModes, ", "))
	flag.String(DbSSLRootCertFlag, "", "Path to the database root certificate file used for database connections")
	flag.Bool(DbDebugFlag, false, "Set Pop to debug mode")
	flag.Bool(DbIamFlag, false, "Use AWS IAM authentication")
	flag.String(DbIamRoleFlag, "", "The arn of the AWS IAM role to assume when connecting to the database.")
	// Required by https://docs.aws.amazon.com/sdk-for-go/api/service/rds/rdsutils/#BuildAuthToken
	flag.String(DbRegionFlag, "", "AWS Region of the database")
	flag.Bool(DbInstrumentedFlag, false, "Use instrumented db driver")
}

// CheckDatabase validates DB command line flags
func CheckDatabase(v *viper.Viper, logger *zap.Logger) error {

	if err := ValidateHost(v, DbHostFlag); err != nil {
		return err
	}

	if err := ValidatePort(v, DbPortFlag); err != nil {
		return err
	}

	dbPool := v.GetInt(DbPoolFlag)
	dbIdlePool := v.GetInt(DbIdlePoolFlag)
	if dbPool < 1 || dbPool > DbPoolMax {
		return &errInvalidDbPool{DbPool: dbPool}
	}

	if dbIdlePool > dbPool {
		return &errInvalidDbIdlePool{DbPool: dbPool, DbIdlePool: dbIdlePool}
	}

	dbEnv := v.GetString(DbEnvFlag)
	if !stringSliceContains(allDbEnvs, dbEnv) {
		return &errInvalidDbEnv{Value: dbEnv, DbEnvs: allDbEnvs}
	}

	sslMode := v.GetString(DbSSLModeFlag)
	if len(sslMode) == 0 || !stringSliceContains(allSSLModes, sslMode) {
		return &errInvalidSSLMode{Mode: sslMode, Modes: allSSLModes}
	}
	if dbEnv == DbEnvContainer && !stringSliceContains(containerSSLModes, sslMode) {
		return errors.Wrap(&errInvalidSSLMode{Mode: sslMode, Modes: containerSSLModes}, "container db env requires SSL connection to the database")
	} else if dbEnv != DbEnvContainer && !stringSliceContains(allSSLModes, sslMode) {
		return &errInvalidSSLMode{Mode: sslMode, Modes: allSSLModes}
	}

	if filename := v.GetString(DbSSLRootCertFlag); len(filename) > 0 {
		b, err := ioutil.ReadFile(filepath.Clean(filename))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error reading %s at %q", DbSSLRootCertFlag, filename))
		}
		tlsCerts := ParseCertificates(string(b))
		logger.Debug(fmt.Sprintf("certificate chain from %s parsed", DbSSLRootCertFlag), zap.Any("count", len(tlsCerts)))
	}

	// Check IAM Authentication
	if v.GetBool(DbIamFlag) {
		// DbRegionFlag must be set if IAM authentication is enabled.
		dbRegion := v.GetString(DbRegionFlag)
		if err := CheckAWSRegionForService(dbRegion, rds.ServiceName); err != nil {
			return errors.Wrap(err, fmt.Sprintf("'%q' is invalid for service %s", DbRegionFlag, rds.ServiceName))
		}

		dbIamRole := v.GetString(DbIamRoleFlag)
		if len(dbIamRole) == 0 {
			return errors.New("database IAM role not provided")
		}
	}

	return nil
}

// InitDatabase initializes a Pop connection from command line flags.
// v is the viper Configuration.
// creds must relate to an assumed role and can't point to a user or task role directly.
// logger is the application logger.
func InitDatabase(v *viper.Viper, creds *credentials.Credentials, logger *zap.Logger) (*pop.Connection, error) {

	dbEnv := v.GetString(DbEnvFlag)
	dbName := v.GetString(DbNameFlag)
	dbHost := v.GetString(DbHostFlag)
	dbPort := strconv.Itoa(v.GetInt(DbPortFlag))
	dbUser := v.GetString(DbUserFlag)
	dbPassword := v.GetString(DbPasswordFlag)
	dbPool := v.GetInt(DbPoolFlag)
	dbIdlePool := v.GetInt(DbIdlePoolFlag)
	dbUseInstrumentedDriver := v.GetBool(DbInstrumentedFlag)

	// Modify DB options by environment
	dbOptions := map[string]string{
		"sslmode": v.GetString(DbSSLModeFlag),
	}

	if dbEnv == DbEnvTest {
		// Leave the test database name hardcoded, since we run tests in the same
		// environment as development, and it's extra confusing to have to swap environment
		// variables before running tests.
		dbName = DbNameTest
	}

	if str := v.GetString(DbSSLRootCertFlag); len(str) > 0 {
		dbOptions["sslrootcert"] = str
	}

	// Construct a safe URL and log it
	s := "postgres://%s:%s@%s:%s/%s?sslmode=%s"
	dbURL := fmt.Sprintf(s, dbUser, "*****", dbHost, dbPort, dbName, dbOptions["sslmode"])
	logger.Info("Connecting to the database", zap.String("url", dbURL), zap.String(DbSSLRootCertFlag, v.GetString(DbSSLRootCertFlag)))

	// Configure DB connection details
	dbConnectionDetails := pop.ConnectionDetails{
		Dialect:  "postgres",
		Driver:   iampg.CustomPostgres,
		Database: dbName,
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		Options:  dbOptions,
		Pool:     dbPool,
		IdlePool: dbIdlePool,
	}

	if v.GetBool(DbIamFlag) {
		// Set a bogus password holder. It will be replaced with an RDS auth token as the password.
		passHolder := "*****"

		iampg.EnableIAM(dbConnectionDetails.Host,
			dbConnectionDetails.Port,
			v.GetString(DbRegionFlag),
			dbConnectionDetails.User,
			passHolder,
			creds,
			iampg.RDSU{},
			time.NewTicker(10*time.Minute), // Refresh every 10 minutes
			logger,
			make(chan bool))

		dbConnectionDetails.Password = passHolder
	}

	if dbUseInstrumentedDriver {
		// to fake pop out, we need to register the otelsql instrumented
		// driver under the driverName that pop would use. To do that,
		// we need to get the otelsql driver.Driver, which is easiest
		// to get from sql.DB.Driver()
		db, err := sql.Open(dbConnectionDetails.Driver, "")
		if err != nil {
			logger.Error("Failed opening uninstrumented connection", zap.Error(err))
			return nil, err
		}
		currentDriver := db.Driver()
		err = db.Close()
		if err != nil {
			logger.Error("Failed closing uninstrumented connection", zap.Error(err))
			return nil, err
		}

		// This is the name from pop's instrumented connection code
		// https://github.com/gobuffalo/pop/blob/master/connection_instrumented.go#L44
		popInstrumentedDriverName := "instrumented-sql-driver-postgres"
		// and we're going to fake out pop with the Driver so that the
		// driver name matches what pop is looking for, but it will
		// wind up using the desired driver under a wrapped otelsql connection
		dbConnectionDetails.Driver = "postgres"
		spanOptions := otelsql.SpanOptions{
			Ping:      true,
			AllowRoot: v.GetBool(DbDebugFlag),
			RowsNext:  v.GetBool(DbDebugFlag),
		}
		sql.Register(popInstrumentedDriverName,
			otelsql.WrapDriver(currentDriver,
				semconv.DBSystemPostgreSQL.Value.AsString(),
				otelsql.WithSpanOptions(spanOptions)))

		// now we can update the connection details to indicate we
		// want an instrumented connection
		dbConnectionDetails.UseInstrumentedDriver = true
		// pop expects at least one option when using instrumented
		// sql, but the options will be ignored since we are faking
		// things out
		dbConnectionDetails.InstrumentedDriverOptions = []instrumentedsql.Opt{
			instrumentedsql.WithOmitArgs(),
		}
		logger.Info("Using otelsql instrumented sql driver")
	}

	err := dbConnectionDetails.Finalize()
	if err != nil {
		logger.Error("Failed to finalize DB connection details", zap.Error(err))
		return nil, err
	}

	// Set up the connection
	connection, err := pop.NewConnection(&dbConnectionDetails)
	if err != nil {
		logger.Error("Failed create DB connection", zap.Error(err))
		return nil, err
	}

	// Open the connection - required
	err = connection.Open()
	if err != nil {
		logger.Error("Failed to open DB connection", zap.Error(err))
		return nil, err
	}

	err = testConnection(&dbConnectionDetails, v.GetBool(DbIamFlag), logger)
	if err != nil {
		logger.Error("Failed to ping database")
		return connection, err
	}

	// Return the open connection
	return connection, nil
}

//testConnection tests the connection to determine successful ping
func testConnection(dbConnDetails *pop.ConnectionDetails, useIam bool, logger *zap.Logger) error {
	// Copy connection info as we don't want to alter connection info
	dbConnectionDetails := pop.ConnectionDetails{
		Dialect:  "postgres",
		Driver:   dbConnDetails.Driver,
		Database: dbConnDetails.Database,
		Host:     dbConnDetails.Host,
		Port:     dbConnDetails.Port,
		User:     dbConnDetails.User,
		Password: dbConnDetails.Password,
		Options:  dbConnDetails.Options,
		Pool:     dbConnDetails.Pool,
		IdlePool: dbConnDetails.IdlePool,
	}

	if useIam {
		dbConnectionDetails.Password = iampg.GetCurrentPass()
	}

	// Set up the connection
	connection, err := pop.NewConnection(&dbConnectionDetails)
	if err != nil {
		logger.Error("Failed create DB connection", zap.Error(err))
		return err
	}

	// Open the connection
	err = connection.Open()
	if err != nil {
		logger.Error("Failed to open DB connection", zap.Error(err))
		return err
	}

	// Check the connection
	db, err := sqlx.Open(connection.Dialect.Details().Dialect, connection.Dialect.URL())
	if err != nil {
		logger.Warn("Failed to open DB by driver name", zap.Error(err))
		return err
	}

	// Make the db ping
	logger.Info("Starting database ping....")
	err = db.Ping()
	if err != nil {
		logger.Warn("Failed to ping DB connection", zap.Error(err))
		return err
	}

	logger.Info("...DB ping successful!")
	return nil
}
