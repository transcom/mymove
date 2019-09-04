package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
)

const (
	// tempMigrationPath is the temporary path for generated migrations
	tempMigrationPath string = "./tmp"

	// localMigrationTemplate is the template for local migration files
	localMigrationTemplate string = `-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on prod/staging/experimental
-- DO NOT include any sensitive data.
`
)

// Close an open file or exit
func closeFile(outfile *os.File) {
	err := outfile.Close()
	if err != nil {
		log.Printf("error closing %s: %v\n", outfile.Name(), err)
		os.Exit(1)
	}
}

func createMigration(path string, filename string, t *template.Template, templateData interface{}) error {
	migrationPath := filepath.Join(path, filename)
	migrationFile, err := os.Create(migrationPath)
	defer closeFile(migrationFile)
	if err != nil {
		return errors.Wrapf(err, "error creating %s", migrationPath)
	}
	err = t.Execute(migrationFile, templateData)
	if err != nil {
		log.Println("error executing template: ", err)
	}
	log.Printf("new migration file created at: %q\n", migrationPath)
	return nil
}

func addMigrationToManifest(migrationManifest string, filename string) error {
	mmf, err := os.OpenFile(migrationManifest, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "could not open migration manifest")
	}
	defer mmf.Close()

	_, err = mmf.WriteString(filename + "\n")
	if err != nil {
		return errors.Wrap(err, "could not append to the migration manifest")
	}

	log.Printf("new migration appended to manifest at: %q\n", migrationManifest)
	return nil
}

func writeEmptyFile(migrationPath, filename string) error {
	path := filepath.Join(migrationPath, filename)

	err := ioutil.WriteFile(path, []byte{}, 0644)
	if err != nil {
		return errors.Wrap(err, "could not write new migration file")
	}

	log.Printf("new migration file created at: %q\n", path)
	return nil
}
