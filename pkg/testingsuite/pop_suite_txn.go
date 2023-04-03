package testingsuite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/gobuffalo/pop/v6"
)

// per process tracking of which dbNum is used
// see findOrCreatePerTestTransactionDb
var lockedDbNum int

// use a driver that sqlx knows uses a dollar binding like postgres
// https://github.com/jmoiron/sqlx/blob/master/bind.go
const fakePopSqlxDriverName = "cloudsqlpostgres"

func containsString(slice []string, target string) bool {
	for _, d := range slice {
		if d == target {
			return true
		}
	}
	return false
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
	suite.preloadedTxnDb = suite.openTxnPopConnection("preload_data")
	// set lowPrivConn so suite.DB() will return the pop connection
	// inside the PreloadData `f`
	suite.lowPrivConn = suite.preloadedTxnDb
	defer func() {
		// ensure we unset the connection so that after preloading
		// data a new savepoint can be establised
		suite.lowPrivConn = nil
	}()
	suite.T().Cleanup(func() {
		// the test context is done, so we can close this db
		// connection (which will rollback the go-txdb transaction)
		err := suite.preloadedTxnDb.Close()
		if err != nil {
			log.Fatalf("Preload db close failed: %v", err)
		}
		suite.preloadedTxnDb = nil
	})

	// call the user provided preload function
	f()
}

// perTestTransactionDB opens a new db connection inside a transaction
// for a test or returns the already existing open one.
func (suite *PopTestSuite) perTestTransactionDB() *pop.Connection {
	// Create the db connection on demand for per test transactions.
	// This is necessary so that we know the current test name and can
	// create a new txdb db connection per test
	suite.txnMutex.Lock()
	defer suite.txnMutex.Unlock()
	testingName := suite.T().Name()
	// set up the pop connection for PreloadData that will
	// rollback at the end of the test
	if suite.preloadedTxnDb != nil {
		if suite.lowPrivConn == nil {
			// At this point, PreloadData has been called, so we want to
			// start a savepoint if this is a new connection. Pop does
			// not support nested transaction transparently, but we
			// know we are in a transaction because we are using
			// go-txdb, so we can manually create and rollback a
			// savepoint (aka a postgres nested transaction)
			err := suite.preloadedTxnDb.RawQuery("SAVEPOINT preload_data").Exec()
			if err != nil {
				log.Fatalf("Cannot start preload savepoint `preload_data`: %v", err)
			}

			// set lowPrivConn so that suite.DB() will use it
			suite.lowPrivConn = suite.preloadedTxnDb

			// ensure the savepoint is rolled back when this test
			// context (suite.T()) is finished
			suite.T().Cleanup(func() {
				// reset lowPrivConn now that this test is done
				suite.lowPrivConn = nil
				err := suite.preloadedTxnDb.RawQuery("ROLLBACK TO SAVEPOINT preload_data").Exec()
				if err != nil {
					log.Fatalf("Preload savepoint `preload_data` rollback failed: %v", err)
				}
			})
		}
		return suite.lowPrivConn
	}

	// At this point, we know this is a transactional test that is not
	// using PreloadData, so we can establish an entirely new
	// connection
	popConn, ok := suite.txnTestDb[testingName]
	if ok {
		return popConn
	}

	// openTxnPopConnection will wind up opening a go-txdb connection
	popConn = suite.openTxnPopConnection(testingName)
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

		// Delete the extra connection since we're about to panic.
		delete(suite.txnTestDb, testingName)
		err := popConn.Close()
		if err != nil {
			log.Panic(err)
		}

		log.Panic("Multiple test databases active simultaneously, use PreloadData: " + names)
	}
	suite.T().Cleanup(func() {
		// when this test context finishes, remove it from tracking
		delete(suite.txnTestDb, testingName)
		// ensure we close the db connection (which will roll back the
		// go-txdb transaction)
		err := popConn.Close()
		if err != nil {
			log.Panic(err)
		}
	})
	return popConn
}

// openTxnPopConnection opens a pop connection for the provided test name
func (suite *PopTestSuite) openTxnPopConnection(testNameOverride string) *pop.Connection {
	packageName := suite.PackageName.String()
	testName := testNameOverride
	if testName == "" {
		testName = suite.T().Name()
	}
	// when connecting using txdb, the URL needs to look like
	// something for pop, but does not have to point to a real db
	dbSanitizer := regexp.MustCompile("[^a-zA-Z0-9_]")
	txnFakeURL := "postgres:///" + dbSanitizer.ReplaceAllString(packageName+"_"+testName, "_")
	txnPopDeets := pop.ConnectionDetails{
		Driver:  fakePopSqlxDriverName,
		Dialect: suite.lowPrivConnDetails.Dialect,
		URL:     txnFakeURL,
	}
	popConn, err := pop.NewConnection(&txnPopDeets)
	if err != nil {
		log.Fatalf("Error creating new pop connection: %v", err)
	}
	err = popConn.Open()
	if err != nil {
		log.Fatalf("Error opening new pop connection: %v", err)
	}
	return popConn
}

func (suite *PopTestSuite) tryAdvisoryLock(dbNum int) bool {
	// lockStart is an arbitrary number, it could be anything
	const lockStart = 10000

	// Use an advisory lock to hold a database until the
	// connection is closed at the end of the package test suite
	// run.
	//
	// See https://www.postgresql.org/docs/current/functions-admin.html#FUNCTIONS-ADVISORY-LOCKS
	var lock bool
	lockQuery := fmt.Sprintf("SELECT pg_try_advisory_lock(%d)", lockStart+dbNum)
	err := suite.pgConn.RawQuery(lockQuery).First(&lock)
	if err != nil {
		log.Panicf("Error trying to get advisory lock: %v", err)
	}
	return lock
}

// findOrCreatePerTestTransactionDb tries to reuse a pool of databases so
// a clone doesn't have to be created, which greatly speeds up the
// tests. Because it is used for tests that rollback a transaction at
// the end of the test, the db can be reused.
func (suite *PopTestSuite) findOrCreatePerTestTransactionDb() {
	pid := os.Getpid()
	packageName := suite.PackageName.String()
	suite.pgConnDetails.Options["application_name"] = packageName
	suite.lowPrivConnDetails.Options["application_name"] = packageName
	suite.lowPrivConnDetails.Driver = fakePopSqlxDriverName
	// when doing per test transaction, high priv conn should never be used
	suite.highPrivConnDetails.Database = "UNUSED"
	if containsString(sql.Drivers(), fakePopSqlxDriverName) {
		// This is only called from NewPopTestSuite, which could be
		// be called multiple times per process if a package has
		// multiple test suites in it.
		//
		// Go runs each package's test in its own process
		//
		// We have already registered the fake txdb driver in this
		// process, so we have already been here before. Once the the
		// fake txdb driver has been registered for a db url, we
		// cannot change it, so we need to wait for that db to become
		// available again

		// use a non blocking advisory lock so we don't get into a
		// database deadlock situation
		for {
			if suite.tryAdvisoryLock(lockedDbNum) {
				break
			}
			// if we didn't get the lock, sleep until we do
			log.Printf("TXNDB: Waiting for locked db %s in pid %d",
				suite.lowPrivConnDetails.Database, pid)
			time.Sleep(time.Second)
		}

		testDbName := fmt.Sprintf("%s_%d", suite.dbNameTemplate, lockedDbNum)
		suite.lowPrivConnDetails.Database = testDbName
		log.Printf("TXNDB: package %s will RE-use database %s in pid %d", packageName,
			suite.lowPrivConnDetails.Database, pid)
		return
	}

	// we are not reusing a database, so get a lock on a new database
	dbNum := 1

	for {
		if suite.tryAdvisoryLock(dbNum) {
			break
		}
		dbNum++
	}

	lockedDbNum = dbNum
	// now we have a lock on dbNum until the pgConn closes
	// the test databases used here look like test_db_1, test_db_2, etc
	testDbName := fmt.Sprintf("%s_%d", suite.dbNameTemplate, lockedDbNum)
	suite.lowPrivConnDetails.Database = testDbName

	// Try to figure out if we need to recreate the test db from the
	// template db by looking at when each was modified
	// If the template db is newer, we need to recreate, otherwise we
	// can reuse
	mtimeQuery := "SELECT (pg_stat_file('base/'||oid ||'/PG_VERSION')).modification FROM pg_database WHERE datname = ?"
	var templateMtime time.Time
	err := suite.pgConn.RawQuery(mtimeQuery, suite.dbNameTemplate).First(&templateMtime)
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
		err = cloneDatabase(suite.pgConn, suite.dbNameTemplate, testDbName)
		if err != nil {
			log.Panic(err)
		}
	}

	// See https://github.com/DATA-DOG/go-txdb for more information
	// about how txdb works and why we need to register a fake driver
	// and then connect to a fake database name
	s := "postgres://%s:%s@%s:%s/%s?%s"
	dataSourceName := fmt.Sprintf(s,
		suite.lowPrivConnDetails.User,
		suite.lowPrivConnDetails.Password,
		suite.lowPrivConnDetails.Host,
		suite.lowPrivConnDetails.Port,
		suite.lowPrivConnDetails.Database,
		suite.lowPrivConnDetails.OptionsString(""))

	log.Printf("TXNDB: package %s will use database %s in pid %d", packageName,
		suite.lowPrivConnDetails.Database, pid)
	// Register will panic if the same driver is registered more than
	// once in the same process, so we've checked for that up above
	txdb.Register(fakePopSqlxDriverName, "postgres", dataSourceName)
}

// tearDownTxnTest closes the db connection established for this test
func (suite *PopTestSuite) tearDownTxnTest() {
	if !suite.usePerTestTransaction {
		return
	}
}
