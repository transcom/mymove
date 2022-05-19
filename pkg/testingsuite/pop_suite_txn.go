package testingsuite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"

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

// PopSuiteTxnStore is a pop.Store that uses go-txdb to wrap the
// connection in a transaction that will be rolled back when closed
type PopSuiteTxnStore struct {
	*sqlx.DB
	popConn *pop.Connection
	txList  []*pop.Tx
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
		for txLen := len(store.txList); txLen > 0; txLen-- {
			tx := store.txList[txLen-1]
			err := tx.Close()
			if err != nil {
				log.Fatalf("TX close failed: %v", err)
			}
		}
		store.txList = []*pop.Tx{}
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
	store.txList = append(store.txList, t)
	// Fake out POP!
	// Because we are using go-txdb to manage the transactions, we can
	// handle nested transactions. Setting TX on the pop connection
	// means the connection can only have a single TX, which breaks
	// some tests
	store.popConn.TX = nil
	return t, nil
}

// openTxnPopConnection sets up the pop Connection for this test
// suite using a per test connection that will run inside a transaction
func (suite *PopTestSuite) openTxnPopConnection() *pop.Connection {
	packageName := suite.PackageName.String()
	testName := ""
	if suite.T() != nil {
		testName = suite.T().Name()
	}

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
	fakeDbName := dbSanitizer.ReplaceAllString(packageName+"_"+testName, "_")
	db := sqlx.MustOpen(fakePopSqlxDriverName, fakeDbName)

	suite.lowPrivConnDetails.Driver = fakePopSqlxDriverName
	suite.lowPrivConnDetails.Options["application_name"] = fakeDbName
	suite.lowPrivConnDetails.Database = fakeDbName
	popConn, err := pop.NewConnection(suite.lowPrivConnDetails)
	if err != nil {
		log.Panic(err)
	}
	suiteStore := &PopSuiteTxnStore{
		db,
		popConn,
		[]*pop.Tx{},
	}
	popConn.Store = suiteStore

	return popConn
}

// tearDownTxnTest closes the db connection established for this test
func (suite *PopTestSuite) tearDownTxnTest() {
	if !suite.usePerTestTransaction {
		return
	}
	suite.perTestTxnMutex.Lock()
	defer suite.perTestTxnMutex.Unlock()
	t := suite.T()
	db, ok := suite.perTestTxnConn[t]
	if ok {
		delete(suite.perTestTxnConn, t)
		err := db.Close()
		if err != nil {
			log.Fatalf("Closing Subtest DB Failed!: %v", err)
		}
	}

}
