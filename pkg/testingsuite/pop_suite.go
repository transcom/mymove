package testingsuite

import (
	"bytes"
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
	PackageName
	db *pop.Connection
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
	dumpErr := bytes.Buffer{}
	dump.Stderr = &dumpErr
	out, err := dump.Output()
	if err != nil {
		return errors.Wrapf(err, "failed to dump the database %s: %s %s", source, string(out), dumpErr.String())
	}

	// #nosec G204
	restore := exec.Command("pg_restore", "-U", "postgres", "-h", "localhost", "-d", destination)
	restore.Stdin = bytes.NewReader(out)

	if op, err := restore.CombinedOutput(); err != nil {
		return errors.Wrapf(err, "failed to run the restore cmd: %s", op)
	}

	return nil
}

// PackageName represents the project-relative name of a Go package.
type PackageName string

func (pn PackageName) String() string {
	return string(pn)
}

// CurrentPackage returns the project-relative name of the caller's package.
//
// "github.com/transcom/mymove/pkg/" is removed from the beginning of the absolute package name, so
// the return value will be e.g. "handlers/publicapi".
func CurrentPackage() PackageName {
	pc, _, _, _ := runtime.Caller(1)
	caller := runtime.FuncForPC(pc)

	fnName := strings.Replace(caller.Name(), "github.com/transcom/mymove/pkg/", "", 1)
	pkg := strings.Split(fnName, ".")[0]
	fmt.Println(caller.Name())
	fmt.Println(fnName)
	fmt.Println(pkg)
	return PackageName(pkg)
}

// NewPopTestSuite returns a new PopTestSuite
func NewPopTestSuite(packageName PackageName) PopTestSuite {
	//pop.Debug = true
	dbName := fmt.Sprintf("test_%s", strings.Replace(packageName.String(), "/", "_", -1))
	log.Printf("package %s is attempting to connect to database %s", packageName.String(), dbName)

	fmt.Printf("attempting to clone database %s to %s... ", "test_db", dbName)
	if err := cloneDatabase("test_db", dbName); err != nil {
		log.Panicf("failed to clone database '%s' to '%s': %#v", "testdb", dbName, err)
	}
	fmt.Println("success")
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

	err = db.RawQuery("SELECT 1;").Exec()
	if err != nil {
		log.Panic(err)
	}

	return PopTestSuite{db: db, PackageName: packageName}
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
