package testingsuite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/gobuffalo/pop/v6"
	"github.com/jmoiron/sqlx"
)

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

// perTestTransactionDB opens a new db connection inside a transaction
// for a test or returns the already existing open one.
func (suite *PopTestSuite) perTestTransactionDB() *pop.Connection {
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

		// Delete the extra connection since we're about to panic.
		delete(suite.txnTestDb, testingName)
		err := popConn.Close()
		if err != nil {
			log.Panic(err)
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

// tearDownTxnTest closes the db connection established for this test
func (suite *PopTestSuite) tearDownTxnTest() {
	if !suite.usePerTestTransaction {
		return
	}
}

// PopSuiteTxnStore is a pop.Store that uses go-txdb to wrap the
// connection in a transaction that will be rolled back when closed
type PopSuiteTxnStore struct {
	*sqlx.DB
	popConn *pop.Connection
}

// Commit the transaction if it exists
// Needed to implement the pop.store interface, but won't actually be called
func (store *PopSuiteTxnStore) Commit() error {
	return nil
}

// Rollback the transaction if it exists
// Needed to implement the pop.store interface, but won't actually be called
func (store *PopSuiteTxnStore) Rollback() error {
	return nil
}

// Transaction starts a pop.Transaction
func (store *PopSuiteTxnStore) Transaction() (*pop.Tx, error) {
	return store.TransactionContext(context.Background())
}

// Close closes any open transactions and db connections
func (store *PopSuiteTxnStore) Close() error {
	if store.popConn != nil {
		err := store.DB.Close()
		if err != nil {
			log.Fatalf("DB close failed: %v", err)
		}
		store.DB = nil
		store.popConn = nil
	}
	return nil
}

// TransactionContext returns the current pop.Transaction
func (store *PopSuiteTxnStore) TransactionContext(ctx context.Context) (*pop.Tx, error) {
	return store.TransactionContextOptions(ctx, nil)
}

// TransactionContextOptions returns the current pop.Transaction
func (store *PopSuiteTxnStore) TransactionContextOptions(ctx context.Context, opts *sql.TxOptions) (*pop.Tx, error) {
	tx, err := store.DB.BeginTxx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("could not create new transaction %w", err)
	}
	t := &pop.Tx{
		ID: seededRand.Int(),
		Tx: tx,
	}
	// Fake out POP!
	// Because we are using go-txdb to manage the transactions, we can
	// handle nested transactions. Setting TX on the pop connection
	// means the connection can only have a single TX, which breaks
	// some tests
	store.popConn.TX = nil
	return t, nil
}
