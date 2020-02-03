package main

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

	"github.com/gobuffalo/pop"

	nflect "github.com/gobuffalo/flect/name"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
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

	// Verbose
	cli.InitVerboseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

type Field struct {
	name    string
	dbName  string
	pointer bool
}

type Model struct {
	name   string
	dbName string
	fields []Field
}

type Column struct {
	Name     string `db:"column_name"`
	DataType string `db:"udt_name"`
	Nullable bool   `db:"is_nullable"`
}

// Use the Go parser to load all structs from the provide go file
func loadModelsFromFile(path string) ([]Model, error) {
	var models []Model

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return models, err
	}

	for _, decl := range node.Decls {
		if g, ok := decl.(*ast.GenDecl); ok {
			if g.Tok == token.TYPE {
				for _, spec := range g.Specs {
					typeSpec := spec.(*ast.TypeSpec)
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						m := Model{
							name:   typeSpec.Name.Name,
							dbName: nflect.Tableize(typeSpec.Name.Name),
						}
						for _, structField := range structType.Fields.List {
							if structField.Tag == nil {
								// If there's no struct tag, we're not using that column with Pop
								continue
							}
							tags := reflect.StructTag(strings.Trim(structField.Tag.Value, "`"))
							f := Field{
								name:   structField.Names[0].Name,
								dbName: tags.Get("db"),
							}
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
			if findErr.Error() == "sql: no rows in result set" {
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
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
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
