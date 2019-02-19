package testingsuite

import (
	"fmt"
	"github.com/codegangsta/envy/lib"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

// PopTestSuite is a suite for testing
type PopTestSuite struct {
	BaseTestSuite
	PackageName string
	db          *pop.Connection
}

func cloneDatabase(source, destination string) error {
	// #nosec G204
	drop := exec.Command("dropdb", "-U", "postgres", "-h", "localhost", "--if-exists", destination)

	if op, err := drop.CombinedOutput(); err != nil {
		return errors.Wrapf(err, "failed to drop the database %s: %s", destination, op)
	}
	// #nosec G204
	create := exec.Command("createdb", "-U", "postgres", "-h", "localhost", destination)

	if op, err := create.CombinedOutput(); err != nil {
		return errors.Wrapf(err, "failed to create the database %s: %s", destination, op)
	}

	// #nosec G204
	dump := exec.Command("pg_dump", "-U", "postgres", "-h", "localhost", "-F", "c", source)
	out, err := dump.StdoutPipe()
	if err != nil {
		return errors.Wrapf(err, "failed to dump the database %s", source)
	}

	// #nosec G204
	restore := exec.Command("pg_restore", "-U", "postgres", "-h", "localhost", "-d", destination)
	restore.Stdin = out

	err = dump.Start()
	if err != nil {
		return errors.Wrap(err, "failed to run the dump cmd")
	}

	if op, err := restore.CombinedOutput(); err != nil {
		return errors.Wrapf(err, "failed to start the restore cmd: %s", op)
	}

	err = dump.Wait()
	if err != nil {
		return errors.Wrap(err, "failed to wait for restore cmd")
	}

	return nil
}

// NewPopTestSuite returns a new PopTestSuite
func NewPopTestSuite() PopTestSuite {
	//pop.Debug = true

	// Find out what package is calling this function, and use that to build the database name.
	pc, _, _, _ := runtime.Caller(1)
	caller := runtime.FuncForPC(pc)

	fnName := strings.Replace(caller.Name(), "github.com/transcom/mymove/pkg/", "", 1)
	pkg := strings.Split(fnName, ".")[0]
	dbName := fmt.Sprintf("test_%s", strings.Replace(pkg, "/", "_", -1))

	if err := cloneDatabase("test_db", dbName); err != nil {
		log.Panicf("failed to clone database\n%#v", err)
	}

	db, err := pop.NewConnection(&pop.ConnectionDetails{
		Dialect:  "postgres",
		Database: dbName,
		Host:     envy.MustGet("DB_HOST"),
		Port:     envy.MustGet("DB_PORT"),
		User:     envy.MustGet("DB_USER"),
		Password: envy.MustGet("DB_PASSWORD"),
	})
	if err != nil {
		log.Panic(err)
	}
	err = db.Open()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(dbName)

	err = db.RawQuery("SELECT 1;").Exec()
	if err != nil {
		log.Panic(err)
	}

	return PopTestSuite{db: db, PackageName: pkg}
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

// NoVerrs prints any errors it receives
func (suite *PopTestSuite) NoVerrs(verrs *validate.Errors) bool {
	if !suite.False(verrs.HasAny()) {
		fmt.Println(verrs.String())
		return false
	}
	return true
}
