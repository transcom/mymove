package testingsuite

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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

// tearDownTxnTest closes the db connection established for this test
func (suite *PopTestSuite) tearDownTxnTest() {
	if !suite.usePerTestTransaction {
		return
	}
	// if suite.lowPrivConn == nil {
	// 	return
	// }
	// err := suite.lowPrivConn.Close()
	// if err != nil {
	// 	log.Panicf("Error closing lowPrivConn: %v", err)
	// }
	// suite.lowPrivConn = nil
}
