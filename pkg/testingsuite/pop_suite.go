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

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/flock"

	// Anonymously import lib/pq driver so it's available to Pop
	_ "github.com/lib/pq"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

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
	db                 *pop.Connection
	dbConnDetails      *pop.ConnectionDetails
	primaryConnDetails *pop.ConnectionDetails
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

// NewPopTestSuite returns a new PopTestSuite
func NewPopTestSuite(packageName PackageName) PopTestSuite {
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
	dbPassword, dbPasswordErr := envy.MustGet("DB_PASSWORD")
	if dbPasswordErr != nil {
		log.Panic(dbPasswordErr)
	}
	dbSSLMode, dbSSLModeErr := envy.MustGet("DB_SSL_MODE")
	if dbSSLModeErr != nil {
		log.Panic(dbSSLModeErr)
	}

	dbOptions := map[string]string{
		"sslmode": dbSSLMode,
	}

	log.Printf("package %s is attempting to connect to database %s", packageName.String(), dbNameTest)
	primaryConnDetails := pop.ConnectionDetails{
		Dialect:  dbDialect,
		Driver:   "postgres",
		Database: dbNameTest,
		Host:     dbHost,
		Port:     dbPortTest,
		User:     dbUser,
		Password: dbPassword,
		Options:  dbOptions,
	}
	primaryConn, primaryConnErr := pop.NewConnection(&primaryConnDetails)
	if primaryConnErr != nil {
		log.Panic(primaryConnErr)
	}
	if openErr := primaryConn.Open(); openErr != nil {
		log.Panic(openErr)
	}

	// Doing this before cloning should pre-clean the DB for all tests
	log.Printf("attempting to truncate the database %s", dbNameTest)
	errTruncateAll := primaryConn.TruncateAll()
	if errTruncateAll != nil {
		log.Panicf("failed to truncate database '%s': %#v", dbNameTest, errTruncateAll)
	}

	uniq := StringWithCharset(6, charset)
	dbNamePackage := fmt.Sprintf("%s_%s_%s", dbNameTest, strings.Replace(packageName.String(), "/", "_", -1), uniq)
	fmt.Printf("attempting to clone database %s to %s... ", dbNameTest, dbNamePackage)
	if err := cloneDatabase(primaryConn, dbNameTest, dbNamePackage); err != nil {
		log.Panicf("failed to clone database '%s' to '%s': %#v", dbNameTest, dbNamePackage, err)
	}
	fmt.Println("success")

	// disconnect from the primary DB
	if err := primaryConn.Close(); err != nil {
		log.Panic(err)
	}

	// Release the lock so other tests can clone the DB
	if err := fileLock.Unlock(); err != nil {
		log.Panic(err)
	}

	log.Printf("package %s is attempting to connect to database %s", packageName.String(), dbNamePackage)

	packageConnDetails := pop.ConnectionDetails{
		Dialect:  dbDialect,
		Driver:   "postgres",
		Database: dbNamePackage,
		Host:     dbHost,
		Port:     dbPortTest,
		User:     dbUser,
		Password: dbPassword,
		Options:  dbOptions,
	}
	packageConn, packageConnErr := pop.NewConnection(&packageConnDetails)
	if packageConnErr != nil {
		log.Panic(packageConnErr)
	}

	if openErr := packageConn.Open(); openErr != nil {
		log.Panic(openErr)
	}

	return PopTestSuite{
		db:                 packageConn,
		dbConnDetails:      &packageConnDetails,
		primaryConnDetails: &primaryConnDetails,
		PackageName:        packageName}
}

// DB returns a db connection
func (suite *PopTestSuite) DB() *pop.Connection {
	return suite.db
}

// MustSave requires saving without errors
func (suite *PopTestSuite) MustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.db.ValidateAndSave(model)
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

	err := suite.db.Destroy(model)
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
	// disconnect from the package DB conn
	if err := suite.DB().Close(); err != nil {
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

	// reconnect to the primary DB
	primaryConn, primaryConnErr := pop.NewConnection(suite.primaryConnDetails)
	if primaryConnErr != nil {
		log.Panic(primaryConnErr)
	}
	if openErr := primaryConn.Open(); openErr != nil {
		log.Panic(openErr)
	}
	// Remove the package DB
	if err := dropDB(primaryConn, (*suite.dbConnDetails).Database); err != nil {
		log.Panic(err)
	}
	// disconnect from the primary DB
	if err := primaryConn.Close(); err != nil {
		log.Panic(err)
	}

	// Release the lock so other tests can clone the DB
	if err := fileLock.Unlock(); err != nil {
		log.Panic(err)
	}
}
