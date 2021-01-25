package testingsuite

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/transcom/mymove/pkg/random"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/flock"

	// Anonymously import lib/pq driver so it's available to Pop
	_ "github.com/lib/pq"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"0123456789"

//RA Summary: gosec - G404 - Insecure random number source (rand)
//RA: gosec detected use of the insecure package math/rand rather than the more secure cryptographically secure pseudo-random number generator crypto/rand.
//RA: This particular usage is mitigated by sourcing the seed from crypto/rand in order to create the new random number using math/rand.
//RA: Second, as part of the testing suite, the need for a secure random number here is not necessary.
//RA Developer Status: Mitigated
//RA Validator: jneuner@mitre.org
//RA Validator Status: Mitigated
//RA Modified Severity: CAT III
//  G404
var seededRand = rand.New(random.NewCryptoSeededSource())

var fileLock = flock.New(os.TempDir() + "/server-test-lock.lock")

// StringWithCharset returns a random string
// https://www.calhoun.io/creating-random-strings-in-go/
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// PopTestSuite is a suite for testing
type PopTestSuite struct {
	BaseTestSuite
	PackageName
	origConn            *pop.Connection
	lowPrivConn         *pop.Connection
	highPrivConn        *pop.Connection
	origConnDetails     *pop.ConnectionDetails
	lowPrivConnDetails  *pop.ConnectionDetails
	highPrivConnDetails *pop.ConnectionDetails

	// Enable this flag to avoid the use of the DB_USER_LOW_PRIV and DB_PASSWORD_LOW_PRIV
	// environment variables and instead fall back to the use of a single, high privileged
	// PostgreSQL database role. This role used to commonly be called the "migrations user".
	// However, this grants too many permissions to the database from the application. Therefore,
	// we created a new user with fewer permissions that are used by most tests.
	//
	// There is one type of situation where we still want to use the PostgreSQL role: areas of the
	// code that are testing migrations. In this situation, the following flag can be set to true
	// to enable the use of the role with elevated permissions.
	//
	// For more details, please see https://dp3.atlassian.net/browse/MB-5197
	useHighPrivsPSQLRole bool
}

func dropDB(conn *pop.Connection, destination string) error {
	dropQuery := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", destination)
	dropErr := conn.RawQuery(dropQuery).Exec()
	if dropErr != nil {
		return dropErr
	}
	return nil
}

func cloneDatabase(conn *pop.Connection, source, destination string) error {
	// Now that the lock is available clone the DB
	// Drop and then Create the DB
	if dropErr := dropDB(conn, destination); dropErr != nil {
		return dropErr
	}
	createErr := conn.RawQuery(fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s;", destination, source)).Exec()
	if createErr != nil {
		return createErr
	}

	return nil
}

// PackageName represents the project-relative name of a Go package.
type PackageName string

func (pn PackageName) String() string {
	return string(pn)
}

// Suffix returns a new PackageName with an underscore and the suffix appended to the end
// suffix should be a snake case string
func (pn PackageName) Suffix(suffix string) PackageName {
	return PackageName(pn.String() + "_" + suffix)
}

// CurrentPackage returns the project-relative name of the caller's package.
//
// "github.com/transcom/mymove/pkg/" is removed from the beginning of the absolute package name, so
// the return value will be e.g. "handlers/internalapi".
func CurrentPackage() PackageName {
	pc, _, _, _ := runtime.Caller(1)
	caller := runtime.FuncForPC(pc)

	fnName := strings.Replace(caller.Name(), "github.com/transcom/mymove/pkg/", "", 1)
	pkg := strings.Split(fnName, ".")[0]
	return PackageName(pkg)
}

// PopTestSuiteOption is type intended to be used to change a PopTestSuite object.
type PopTestSuiteOption func(*PopTestSuite)

// WithHighPrivPSQLRole is a functional option that can be passed into the NewPopTestSuite
// function to create a PopTestSuite that only uses the privileged SQL connection.
func WithHighPrivPSQLRole() PopTestSuiteOption {
	return func(pts *PopTestSuite) {
		// Mark a flag that indicates that we are only using a single privileged role.
		pts.useHighPrivsPSQLRole = true

		// Disconnect the low privileged connection and replace its connection and connection
		// details with those of the high privileged connection.
		if err := pts.lowPrivConn.Close(); err != nil {
			log.Panic(err)
		}

		pts.lowPrivConn = pts.highPrivConn
		pts.lowPrivConnDetails = pts.highPrivConnDetails
	}
}

// NewPopTestSuite returns a new PopTestSuite
func NewPopTestSuite(packageName PackageName, opts ...PopTestSuiteOption) PopTestSuite {
	// Try to obtain the lock in this method within 10 minutes
	lockCtx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	// Continually check if the lock is available
	_, lockErr := fileLock.TryLockContext(lockCtx, 678*time.Millisecond)
	if lockErr != nil {
		log.Panic(lockErr)
	}

	dbDialect := "postgres"
	dbName, dbNameErr := envy.MustGet("DB_NAME")
	if dbNameErr != nil {
		log.Panic(dbNameErr)
	}
	dbNameTest := envy.Get("DB_NAME_TEST", dbName)
	dbHost, dbHostErr := envy.MustGet("DB_HOST")
	if dbHostErr != nil {
		log.Panic(dbHostErr)
	}
	dbPort, dbPortErr := envy.MustGet("DB_PORT")
	if dbPortErr != nil {
		log.Panic(dbPortErr)
	}
	dbPortTest := envy.Get("DB_PORT_TEST", dbPort)
	dbUser, dbUserErr := envy.MustGet("DB_USER")
	if dbUserErr != nil {
		log.Panic(dbUserErr)
	}
	dbUserLowPriv, dbUserLowPrivErr := envy.MustGet("DB_USER_LOW_PRIV")
	if dbUserLowPrivErr != nil {
		log.Panic(dbUserLowPrivErr)
	}
	dbPassword, dbPasswordErr := envy.MustGet("DB_PASSWORD")
	if dbPasswordErr != nil {
		log.Panic(dbPasswordErr)
	}
	dbPasswordApp, dbPasswordAppErr := envy.MustGet("DB_PASSWORD_LOW_PRIV")
	if dbPasswordAppErr != nil {
		log.Panic(dbPasswordAppErr)
	}
	dbSSLMode := envy.Get("DB_SSL_MODE", "disable")

	dbOptions := map[string]string{
		"sslmode": dbSSLMode,
	}

	log.Printf("package %s is attempting to connect to database %s", packageName.String(), dbNameTest)
	origConnDetails := pop.ConnectionDetails{
		Dialect:  dbDialect,
		Driver:   "postgres",
		Database: dbNameTest,
		Host:     dbHost,
		Port:     dbPortTest,
		User:     dbUser,
		Password: dbPassword,
		Options:  dbOptions,
	}
	origConn, origConnErr := pop.NewConnection(&origConnDetails)
	if origConnErr != nil {
		log.Panic(origConnErr)
	}
	if openErr := origConn.Open(); openErr != nil {
		log.Panic(openErr)
	}

	// Doing this before cloning should pre-clean the DB for all tests
	log.Printf("attempting to truncate the database %s", dbNameTest)
	errTruncateAll := origConn.TruncateAll()
	if errTruncateAll != nil {
		log.Panicf("failed to truncate database '%s': %#v", dbNameTest, errTruncateAll)
	}

	uniq := StringWithCharset(6, charset)
	dbNamePackage := fmt.Sprintf("%s_%s_%s", dbNameTest, strings.Replace(packageName.String(), "/", "_", -1), uniq)
	log.Printf("attempting to clone database %s to %s... ", dbNameTest, dbNamePackage)
	if err := cloneDatabase(origConn, dbNameTest, dbNamePackage); err != nil {
		log.Panicf("failed to clone database '%s' to '%s': %#v", dbNameTest, dbNamePackage, err)
	}
	log.Println("success")

	// disconnect from the original DB
	if err := origConn.Close(); err != nil {
		log.Panic(err)
	}

	// Release the lock so other tests can clone the DB
	if err := fileLock.Unlock(); err != nil {
		log.Panic(err)
	}

	log.Printf("package %s is attempting to connect to database %s", packageName.String(), dbNamePackage)

	// Prepare a new connection to the temporary database with the same PostgreSQL role privileges
	// as what the migrations task will have when running migrations. These privileges will be
	// higher than the role that runs application.
	highPrivConnDetails := pop.ConnectionDetails{
		Dialect:  dbDialect,
		Driver:   "postgres",
		Database: dbNamePackage,
		Host:     dbHost,
		Port:     dbPortTest,
		User:     dbUser,
		Password: dbPassword,
		Options:  dbOptions,
	}
	highPrivsConn, highPrivsConnErr := pop.NewConnection(&highPrivConnDetails)
	if highPrivsConnErr != nil {
		log.Panic(highPrivsConnErr)
	}
	if openErr := highPrivsConn.Open(); openErr != nil {
		log.Panic(openErr)
	}

	// Prepare a new connection to the temporary database with the same PostgreSQL role privileges
	// as what the application will have when running the server. These privileges will be lower
	// than the role that runs database migrations.
	lowPrivConnDetails := pop.ConnectionDetails{
		Dialect:  dbDialect,
		Driver:   "postgres",
		Database: dbNamePackage,
		Host:     dbHost,
		Port:     dbPortTest,
		User:     dbUserLowPriv,
		Password: dbPasswordApp,
		Options:  dbOptions,
	}
	lowPrivsConn, lowPrivsConnErr := pop.NewConnection(&lowPrivConnDetails)
	if lowPrivsConnErr != nil {
		log.Panic(lowPrivsConnErr)
	}
	if openErr := lowPrivsConn.Open(); openErr != nil {
		log.Panic(openErr)
	}

	// Create a standardized PopTestSuite object.
	pts := &PopTestSuite{
		lowPrivConn:         lowPrivsConn,
		highPrivConn:        highPrivsConn,
		origConn:            origConn,
		lowPrivConnDetails:  &lowPrivConnDetails,
		highPrivConnDetails: &highPrivConnDetails,
		origConnDetails:     &origConnDetails,
		PackageName:         packageName,
	}

	// Apply the user-supplied options to the PopTestSuite object.
	for _, opt := range opts {
		opt(pts)
	}

	return *pts
}

// DB returns a db connection
func (suite *PopTestSuite) DB() *pop.Connection {
	return suite.lowPrivConn
}

// Truncate deletes all data from the specified tables.
func (suite *PopTestSuite) Truncate(tables []string) error {
	// Truncate the specified tables.
	for _, table := range tables {
		sql := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		if err := suite.highPrivConn.RawQuery(sql).Exec(); err != nil {
			return err
		}
	}
	return nil
}

// TruncateAll deletes all data from all tables that are owned by the user connected to the
// database.
func (suite *PopTestSuite) TruncateAll() error {
	return suite.highPrivConn.TruncateAll()
}

// MustSave requires saving without errors
func (suite *PopTestSuite) MustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.lowPrivConn.ValidateAndSave(model)
	if err != nil {
		suite.T().Errorf("Errors encountered saving %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Errorf("Validation errors encountered saving %v: %v", model, verrs)
	}
}

// MustCreate requires creating without errors
func (suite *PopTestSuite) MustCreate(db *pop.Connection, model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := db.ValidateAndCreate(model)
	if err != nil {
		suite.T().Errorf("Errors encountered creating %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Errorf("Validation errors encountered creating %v: %v", model, verrs)
	}
}

// MustDestroy requires deleting without errors
func (suite *PopTestSuite) MustDestroy(model interface{}) {
	t := suite.T()
	t.Helper()

	err := suite.lowPrivConn.Destroy(model)
	if err != nil {
		suite.T().Errorf("Errors encountered destroying %v: %v", model, err)
	}
}

// NoVerrs prints any errors it receives
func (suite *PopTestSuite) NoVerrs(verrs *validate.Errors) bool {
	if !suite.False(verrs.HasAny()) {
		fmt.Println(verrs.String())
		return false
	}
	return true
}

// TearDown runs the teardown for step for the suite
// Important steps are to close open DB connections and drop the DB
func (suite *PopTestSuite) TearDown() {
	// disconnect from the package DB connections
	if err := suite.lowPrivConn.Close(); err != nil {
		log.Panic(err)
	}
	if err := suite.highPrivConn.Close(); err != nil {
		log.Panic(err)
	}

	// Try to obtain the lock in this method within 10 minutes
	lockCtx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	// Continually check if the lock is available
	_, lockErr := fileLock.TryLockContext(lockCtx, 678*time.Millisecond)
	if lockErr != nil {
		log.Panic(lockErr)
	}

	// reconnect to the original DB
	origConn, origConnErr := pop.NewConnection(suite.origConnDetails)
	if origConnErr != nil {
		log.Panic(origConnErr)
	}
	if openErr := origConn.Open(); openErr != nil {
		log.Panic(openErr)
	}
	// Remove the package DB
	if err := dropDB(origConn, (*suite.lowPrivConnDetails).Database); err != nil {
		log.Panic(err)
	}
	// disconnect from the original DB
	if err := origConn.Close(); err != nil {
		log.Panic(err)
	}

	// Release the lock so other tests can clone the DB
	if err := fileLock.Unlock(); err != nil {
		log.Panic(err)
	}
}
