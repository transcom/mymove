package main

// This script audits the definitions of models in pkg/models
// and checks for mismatches between struct fields' types
// and column definitions in the database.

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/gobuffalo/pop/v5"

	nflect "github.com/gobuffalo/flect/name"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
)

type logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

func checkConfig(v *viper.Viper, logger logger) error {
	logger.Debug("checking config")

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	return nil
}

func initFlags(flag *pflag.FlagSet) {

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Logging Levels
	cli.InitLoggingFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

// Field represents a field in a DB
type Field struct {
	name    string
	dbName  string
	pointer bool
}

// Model represents a model in the code corresponding to a DB table
type Model struct {
	name   string
	dbName string
	fields []Field
}

// Column represents a column in a table in the DB
// The fields in this struct have to be public so that they
// can be set by Pop.
type Column struct {
	Name     string `db:"column_name"`
	DataType string `db:"udt_name"`
	Nullable bool   `db:"is_nullable"`
}

// Use the Go parser to load all structs from the provided go file
func loadModelsFromFile(path string) ([]Model, error) {
	var models []Model

	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return models, err
	}

	// Loop through all declarations in the parsed file
	for _, decl := range parsedFile.Decls {
		if g, ok := decl.(*ast.GenDecl); ok {

			// We are only interested in type declarations
			if g.Tok == token.TYPE {
				for _, spec := range g.Specs {
					typeSpec := spec.(*ast.TypeSpec)

					// Check to see if the type declaration is for a struct
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {

						// Build an internal representation of the model, which will
						// also house the fields when we get to them.
						m := Model{
							name:   typeSpec.Name.Name,                  // Name of the struct, e.g. ServiceMember
							dbName: nflect.Tableize(typeSpec.Name.Name), // e.g. service_members
						}

						// Loop through fields and add them to m
						for _, structField := range structType.Fields.List {
							if structField.Tag == nil {
								// If there's no struct tag, we're not using that column with Pop
								continue
							}

							// The struct tag string we get back from the parser includes backticks,
							// but reflect.StructTag expects only the string inside those backticks.
							trimmed := strings.Trim(structField.Tag.Value, "`")

							// Parse the field's tags, e.g. `db:"col_name" json:"col_name"`
							// into a data structure so we can fetch the values for specific tags.
							tags := reflect.StructTag(trimmed)

							f := Field{
								name:   structField.Names[0].Name,
								dbName: tags.Get("db"),
							}

							// Track if the field's type is a pointer type
							if _, ok := structField.Type.(*ast.StarExpr); ok {
								f.pointer = true
							}
							m.fields = append(m.fields, f)
						}
						models = append(models, m)
					}
				}
			}
		}
	}

	return models, nil
}

// Using a model definition, check that all matching columns in the database have compatible nullability
// Columns that aren't found are ignored.
func auditModel(db *pop.Connection, model Model) (bool, error) {
	printedModelName := false
	mismatch := false
	sql := "select column_name, udt_name, is_nullable::boolean from information_schema.columns where table_name=$1 AND column_name=$2"
	for _, field := range model.fields {
		var column Column
		query := db.RawQuery(sql, model.dbName, field.dbName)
		if findErr := query.First(&column); findErr != nil {
			if findErr.Error() == models.RecordNotFoundErrorString {
				continue
			} else {
				return false, findErr
			}
		}
		if field.pointer != column.Nullable {
			if !printedModelName {
				fmt.Printf("\n%s\n", model.name)
				printedModelName = true
			}
			mismatch = true
			var nullable string
			if column.Nullable {
				nullable = "NULL"
			} else {
				nullable = "NOT NULL"
			}
			pointer := ""
			if field.pointer {
				pointer = "*"
			}

			fmt.Printf("  %s%s : %s is %v\n", pointer, field.name, field.dbName, nullable)
		}
	}
	return mismatch, nil
}

func main() {
	flag := pflag.CommandLine
	initFlags(flag)
	err := flag.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("Could not parse flags: %v\n", err)
	}

	v := viper.New()
	bindErr := v.BindPFlags(flag)
	if bindErr != nil {
		log.Fatal("failed to bind flags", zap.Error(bindErr))
	}

	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	dbConnection, err := cli.InitDatabase(v, nil, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	// Track if we have found at least one issue
	fail := false

	files, dirErr := ioutil.ReadDir("./pkg/models")
	if dirErr != nil {
		logger.Fatal("reading directory", zap.Error(dirErr))
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if models, loadErr := loadModelsFromFile("./pkg/models/" + file.Name()); loadErr == nil {
			for _, model := range models {
				mismatch, auditErr := auditModel(dbConnection, model)
				if auditErr != nil {
					logger.Fatal("auditing model", zap.Error(auditErr))
				}
				if mismatch {
					fail = true
				}
			}
		} else {
			logger.Fatal("loading models", zap.Error(loadErr))
		}
	}

	if fail {
		// There was at least one mismatch, so exit non-zero
		os.Exit(1)
	}
}
