package testingsuite

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jmoiron/sqlx"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/random"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"

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
// #nosec G404
var seededRand = rand.New(random.NewCryptoSeededSource())

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
	dbNameTemplate      string
	logger              *zap.Logger
	pgConn              *pop.Connection
	lowPrivConn         *pop.Connection
	highPrivConn        *pop.Connection
	pgConnDetails       *pop.ConnectionDetails
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

	// Enable this flag to use per test transactions
	usePerTestTransaction bool
	// the preloadedTxnDb is used when PreloadData is invoked
	preloadedTxnDb *sqlx.DB
	// a mutex to make sure transactional tests aren't stepping on each other
	txnMutex sync.Mutex
	// tracking the *testing.T that are active to ensure multiple
	// databases are not open at the same time
	txnTestDb map[string]*pop.Connection
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
	}
}

// WithPerTestTransaction is a functional option that can be passed
// into the NewPopTestSuite function to create a PopTestSuite that
// runs each test inside a transaction and rolls back at the end of
// the test. See also PopTestSuite#PreloadData for how to share data
// between subtests. For implementation details see PopTestSuite#DB
//
// When WithPerTestTransaction is enabled, it means each subtest has a
// database that is isolated from other tests.

// If per test transaction is enabled and PreloadData has not been
// called, each subtest gets a new connection. This means subtests are
// really just like main tests, but subtests are a helpful way to
// group tests together, which can be useful.
//
// When using per test transactions, it works both with
//
// suite.Run("name", func( { ... })
//
// and
//
// suite.T().Run("name", func(t *testing.T) { ... })

func WithPerTestTransaction() PopTestSuiteOption {
	return func(pts *PopTestSuite) {
		pts.usePerTestTransaction = true
	}
}

// NewPopTestSuite returns a new PopTestSuite
func NewPopTestSuite(packageName PackageName, opts ...PopTestSuiteOption) *PopTestSuite {
	// Create a standardized PopTestSuite object.
	pts := &PopTestSuite{
		PackageName: packageName,
	}

	// Apply the user-supplied options to the PopTestSuite object.
	for _, opt := range opts {
		opt(pts)
	}

	if pts.useHighPrivsPSQLRole && pts.usePerTestTransaction {
		log.Fatal("Cannot use both high priv psql and per test transaction")
	}

	pts.getDbConnectionDetails()

	log.Printf("package %s is attempting to connect to database %s", packageName, pts.pgConnDetails.Database)
	pgConn, err := pop.NewConnection(pts.pgConnDetails)
	if err != nil {
		log.Panic(err)
	}
	if err = pgConn.Open(); err != nil {
		log.Panic(err)
	}
	pts.pgConn = pgConn

	if pts.usePerTestTransaction {
		pts.findOrCreatePerTestTransactionDb()
	} else {
		// set up database connections for non per test transactions
		// which may or may not be have useHighPrivsPSQLRole set
		pts.highPrivConn, err = pop.NewConnection(pts.highPrivConnDetails)
		if err != nil {
			log.Panic(err)
		}
		if err = pts.highPrivConn.Open(); err != nil {
			log.Panic(err)
		}

		pts.lowPrivConn, err = pop.NewConnection(pts.lowPrivConnDetails)
		if err != nil {
			log.Panic(err)
		}
		if err := pts.lowPrivConn.Open(); err != nil {
			log.Panic(err)
		}

		log.Printf("attempting to clone database %s to %s... ", pts.dbNameTemplate, pts.lowPrivConnDetails.Database)
		if err := cloneDatabase(pgConn, pts.dbNameTemplate, pts.lowPrivConnDetails.Database); err != nil {
			log.Panicf("failed to clone database '%s' to '%s': %#v", pts.dbNameTemplate, pts.lowPrivConnDetails.Database, err)
		}
		log.Println("success")

		// The db is already truncated as part of the test setup
	}

	if pts.useHighPrivsPSQLRole {
		// Disconnect the low privileged connection and replace its
		// connection and connection details with those of the high
		// privileged connection.
		if err := pts.lowPrivConn.Close(); err != nil {
			log.Panic(err)
		}

		pts.lowPrivConn = pts.highPrivConn
		pts.lowPrivConnDetails = pts.highPrivConnDetails
	}

	pts.txnTestDb = make(map[string]*pop.Connection)

	return pts
}

func (suite *PopTestSuite) getDbConnectionDetails() {
	dbDialect := "postgres"
	dbNameTest, err := envy.MustGet("DB_NAME_TEST")
	if err != nil {
		log.Panic(err)
	}
	dbHost, err := envy.MustGet("DB_HOST")
	if err != nil {
		log.Panic(err)
	}
	dbPort, err := envy.MustGet("DB_PORT")
	if err != nil {
		log.Panic(err)
	}
	dbPortTest := envy.Get("DB_PORT_TEST", dbPort)
	dbUser, err := envy.MustGet("DB_USER")
	if err != nil {
		log.Panic(err)
	}
	dbUserLowPriv, err := envy.MustGet("DB_USER_LOW_PRIV")
	if err != nil {
		log.Panic(err)
	}
	dbPassword, err := envy.MustGet("DB_PASSWORD")
	if err != nil {
		log.Panic(err)
	}
	dbPasswordApp, err := envy.MustGet("DB_PASSWORD_LOW_PRIV")
	if err != nil {
		log.Panic(err)
	}
	dbSSLMode := envy.Get("DB_SSL_MODE", "disable")

	dbOptions := map[string]string{
		"sslmode": dbSSLMode,
	}

	// Connect to postgres db to clone the test database
	//
	// The tests should never connect to dbNameTest directly because
	// postgres cannot use a database as a template if there is a
	// connection to it
	//
	// This way we don't need a lock to prevent simultaneous tests
	// from cloning
	suite.pgConnDetails = &pop.ConnectionDetails{
		Dialect:  dbDialect,
		Driver:   "postgres",
		Database: "postgres",
		Host:     dbHost,
		Port:     dbPortTest,
		User:     dbUser,
		Password: dbPassword,
		Options:  dbOptions,
	}

	uniq := StringWithCharset(6, charset)
	dbNamePackage := fmt.Sprintf("%s_%s_%s", dbNameTest, strings.Replace(suite.PackageName.String(), "/", "_", -1), uniq)

	// Prepare a new connection to the temporary database with the
	// same PostgreSQL role privileges as what the application will
	// have when running the server. These privileges will be lower
	// than the role that runs database migrations.
	suite.lowPrivConnDetails = &pop.ConnectionDetails{
		Dialect:  dbDialect,
		Driver:   "postgres",
		Database: dbNamePackage,
		Host:     dbHost,
		Port:     dbPortTest,
		User:     dbUserLowPriv,
		Password: dbPasswordApp,
		Options:  dbOptions,
	}
	// Prepare a new connection to the temporary database with the
	// same PostgreSQL role privileges as what the migrations task
	// will have when running migrations. These privileges will be
	// higher than the role that runs application.
	suite.highPrivConnDetails = &pop.ConnectionDetails{
		Dialect:  dbDialect,
		Driver:   "postgres",
		Database: dbNamePackage,
		Host:     dbHost,
		Port:     dbPortTest,
		User:     dbUser,
		Password: dbPassword,
		Options:  dbOptions,
	}

	suite.dbNameTemplate = dbNameTest
}

// findOrCreatePerTestTransactionDb tries to reuse a pool of databases so
// a clone doesn't have to be created, which greatly speeds up the
// tests. Because it is used for tests that rollback a transaction at
// the end of the test, the db can be reused.
func (suite *PopTestSuite) findOrCreatePerTestTransactionDb() {
	packageName := suite.PackageName.String()
	suite.pgConnDetails.Options["application_name"] = packageName
	suite.lowPrivConnDetails.Options["application_name"] = packageName
	// lockStart is an arbitrary number, it could be anything
	lockStart := 10000
	dbNum := 1

	// Use an advisory lock to hold a database until the
	// connection is closed at the end of the package test suite
	// run.
	//
	// See https://www.postgresql.org/docs/current/functions-admin.html#FUNCTIONS-ADVISORY-LOCKS
	var lock bool
	for {
		lockQuery := fmt.Sprintf("SELECT pg_try_advisory_lock(%d)", lockStart+dbNum)
		err := suite.pgConn.RawQuery(lockQuery).First(&lock)
		if err != nil {
			log.Panic(err)
		}
		if lock {
			break
		}
		dbNum++
	}
	// now we have a lock on dbNum until the pgConn closes
	templateDbName := suite.dbNameTemplate
	// the test databases used here look like test_db_1, test_db_2, etc
	testDbName := fmt.Sprintf("%s_%d", templateDbName, dbNum)
	// when doing per test transaction, high priv conn should never be used
	suite.highPrivConnDetails.Database = "UNUSED"
	suite.lowPrivConnDetails.Database = testDbName

	// Try to figure out if we need to recreate the test db from the
	// template db by looking at when each was modified
	// If the template db is newer, we need to recreate, otherwise we
	// can reuse
	mtimeQuery := "SELECT (pg_stat_file('base/'||oid ||'/PG_VERSION')).modification FROM pg_database WHERE datname = ?"
	var templateMtime time.Time
	err := suite.pgConn.RawQuery(mtimeQuery, templateDbName).First(&templateMtime)
	if err != nil {
		log.Panic(err)
	}
	var testDbMtime time.Time
	err = suite.pgConn.RawQuery(mtimeQuery, testDbName).First(&testDbMtime)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	}
	if testDbMtime.Unix() < templateMtime.Unix() {
		// If the testdb was modified before the source, we need to
		// recreate it
		err = dropDB(suite.pgConn, testDbName)
		if err != nil {
			log.Panic(err)
		}
		err = cloneDatabase(suite.pgConn, templateDbName, testDbName)
		if err != nil {
			log.Panic(err)
		}
	}
}

// openTxnDb sets up the sqlx.DB for this test that will rollback
// all changes when the db is closed
func (suite *PopTestSuite) openTxnDb(name string) *sqlx.DB {
	packageName := suite.PackageName.String()

	s := "postgres://%s:%s@%s:%s/%s?%s"
	dataSourceName := fmt.Sprintf(s, suite.lowPrivConnDetails.User,
		suite.lowPrivConnDetails.Password,
		suite.lowPrivConnDetails.Host,
		suite.lowPrivConnDetails.Port,
		suite.lowPrivConnDetails.Database,
		suite.lowPrivConnDetails.OptionsString(""))

	// See https://github.com/DATA-DOG/go-txdb for more information
	// about how txdb works and why we need to register a fake driver
	// and then connect to a fake database name
	if !containsString(sql.Drivers(), fakePopSqlxDriverName) {
		txdb.Register(fakePopSqlxDriverName, "postgres", dataSourceName)
	}

	dbSanitizer := regexp.MustCompile("[^a-zA-Z0-9_]")
	fakeDbName := dbSanitizer.ReplaceAllString(packageName+"_"+name, "_")

	suite.lowPrivConnDetails.Driver = fakePopSqlxDriverName
	suite.lowPrivConnDetails.Options["application_name"] = fakeDbName
	suite.lowPrivConnDetails.Database = fakeDbName
	return sqlx.MustOpen(fakePopSqlxDriverName, fakeDbName)
}

// txnPopConnection wraps the sqlx.DB in a *pop.Connection
func (suite *PopTestSuite) txnPopConnection(db *sqlx.DB) *pop.Connection {
	popConn, err := pop.NewConnection(suite.lowPrivConnDetails)
	if err != nil {
		log.Panic(err)
	}
	suiteStore := &PopSuiteTxnStore{
		db,
		popConn,
	}
	popConn.Store = suiteStore
	return popConn
}

// PreloadData sets up the data used by all subtests. It can only be
// called once. Once it is called, all subtests will use the same
// database, but will run inside a savepoint that are rolled back.
// This way each subtest can reuse the preloaded data, but each
// subtest cannot modify the data used by another test.
//
// Note that calling PreloadData automatically changes how suite.Run
// works
func (suite *PopTestSuite) PreloadData(f func()) {
	if !suite.usePerTestTransaction {
		log.Panic("Cannot use PreloadData without per test transaction")
	}
	if suite.preloadedTxnDb != nil {
		log.Panic("Cannot call PreloadData multiple times in the same test")
	}
	// establish a database connection for the suite that will
	// rollback when closed
	suite.preloadedTxnDb = suite.openTxnDb("preload_data")
	// set lowPrivConn so suite.DB() will return the pop connection
	suite.lowPrivConn = suite.txnPopConnection(suite.preloadedTxnDb)
	defer func() {
		// ensure we unset the connection so that after preloading
		// data a new pop connection will be established
		suite.lowPrivConn = nil
	}()
	suite.T().Cleanup(func() {
		err := suite.preloadedTxnDb.Close()
		if err != nil {
			log.Panic(err)
		}
		suite.preloadedTxnDb = nil
	})
	f()
}

// DB returns a db connection. It has logic to handle per test
// transactions, both when PreloadData is called and when it is not
func (suite *PopTestSuite) DB() *pop.Connection {
	if !suite.usePerTestTransaction {
		return suite.lowPrivConn
	}

	// Create the db connection on demand for per test transactions.
	// This is necessary so that we know the current test name and can
	// create a new txdb db connection per test
	suite.txnMutex.Lock()
	defer suite.txnMutex.Unlock()
	// set up the pop connection for PreloadData that will
	// rollback at the end of the test
	if suite.preloadedTxnDb != nil {
		// At this point, PreloadData has been called, so we want to
		// start a savepoint. This is what go-txdb does for us. See
		// https://github.com/DATA-DOG/go-txdb/blob/master/db.go#L195
		if suite.lowPrivConn == nil {
			tx, err := suite.preloadedTxnDb.Beginx()
			if err != nil {
				log.Panic(err)
			}
			suite.lowPrivConn = suite.txnPopConnection(suite.preloadedTxnDb)
			// ensure the savepoint is rolled back when this test
			// context (suite.T()) is finished
			suite.T().Cleanup(func() {
				suite.lowPrivConn = nil
				err = tx.Rollback()
				if err != nil {
					log.Panic(err)
				}
			})
		}
		return suite.lowPrivConn
	}

	// At this point, we know this is a transactional test that is not
	// using PreloadData, so we can establish a new connection using
	// go-txdb.
	testingName := suite.T().Name()
	popConn, ok := suite.txnTestDb[testingName]
	if ok {
		return popConn
	}
	// create a brand new database for this test so it is
	// completely isolated
	db := suite.openTxnDb(testingName)
	popConn = suite.txnPopConnection(db)
	suite.txnTestDb[testingName] = popConn

	// To prevent accidental dependencies between tests, check and
	// make sure that this test is only connecting to a single
	// database at a time.
	// This prevents situations like
	//
	// func (suite *MySuite) TopTest() {
	//   move = makeMove(suite.DB())
	//   suite.Run("subtest 1", func() {
	//     suite.NoError(doSomething(&move))
	//   })
	//
	//   modifyMove(&move)
	//   suite.Run("subtest 2", func( {
	//     suite.NoError(doSomethingElseWithModifiedMove(&move))
	//   })
	// }
	if len(suite.txnTestDb) > 1 {
		names := ""
		i := 0
		for k := range suite.txnTestDb {
			if i == 0 {
				names = k
			} else {
				names = names + "," + k
			}
			i++
		}
		log.Panic("Multiple test databases active simultaneously, use PreloadData: " + names)
	}
	suite.T().Cleanup(func() {
		delete(suite.txnTestDb, testingName)
		err := popConn.Close()
		if err != nil {
			log.Panic(err)
		}
	})
	return popConn
}

// Logger returns the logger for the test suite
func (suite *PopTestSuite) Logger() *zap.Logger {
	if suite.logger == nil {
		suite.logger = zaptest.NewLogger(suite.T())
	}
	return suite.logger
}

// AppContextForTest returns the AppContext for the test suite
func (suite *PopTestSuite) AppContextForTest() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.Logger(), nil)
}

// AppContextWithSessionForTest returns the AppContext for the test suite
func (suite *PopTestSuite) AppContextWithSessionForTest(session *auth.Session) appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.Logger(), session)
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

// TruncateAll deletes all data from all tables that are owned by the
// user connected to the database.
func (suite *PopTestSuite) TruncateAll() error {
	if suite.usePerTestTransaction {
		log.Fatal("Cannot TruncateAll with Per Test Transaction")
	}
	return suite.highPrivConn.TruncateAll()
}

// MustSave requires saving without errors
func (suite *PopTestSuite) MustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.DB().ValidateAndSave(model)
	if err != nil {
		suite.T().Errorf("Errors encountered saving %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Errorf("Validation errors encountered saving %v: %v", model, verrs)
	}
}

// MustCreate requires creating without errors
func (suite *PopTestSuite) MustCreate(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.DB().ValidateAndCreate(model)
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

	err := suite.DB().Destroy(model)
	if err != nil {
		suite.T().Errorf("Errors encountered destroying %v: %v", model, err)
	}
}

// NoVerrs prints any errors it receives
func (suite *PopTestSuite) NoVerrs(verrs *validate.Errors) bool {
	if !suite.False(verrs.HasAny()) {
		suite.Fail(verrs.String())
		return false
	}
	return true
}

// NilOrNoVerrs checks that an error is effecively nil
func (suite *PopTestSuite) NilOrNoVerrs(err error) {
	switch verr := err.(type) {
	case *validate.Errors:
		suite.False(verr.HasAny(), "non-empty validation errors: %v", verr)
	default:
		suite.Nil(err)
	}
}

// TearDownTest runs the teardown per test. It will only do something
// useful if per test transactions are enabled
func (suite *PopTestSuite) TearDownTest() {
	suite.tearDownTxnTest()
}

// TearDown runs the teardown for step for the suite
// Important steps are to close open DB connections and drop the DB
func (suite *PopTestSuite) TearDown() {
	// disconnect from the package DB connections
	if suite.lowPrivConn != nil {
		if err := suite.lowPrivConn.Close(); err != nil {
			log.Panic(err)
		}
	}
	if suite.highPrivConn != nil {
		if err := suite.highPrivConn.Close(); err != nil {
			log.Panic(err)
		}
	}

	// Remove the package DB if this isn't a per test transaction
	if !suite.usePerTestTransaction {
		if err := dropDB(suite.pgConn, (*suite.lowPrivConnDetails).Database); err != nil {
			log.Panic(err)
		}
	}
	err := suite.pgConn.Close()
	if err != nil {
		log.Panic(err)
	}
}
