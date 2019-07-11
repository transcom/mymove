package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"

	"github.com/gofrs/uuid"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/models"
)

const (
	// OfficeUsersFilenameFlag filename containing the details for new office users
	OfficeUsersFilenameFlag string = "office-users-filename"
	// OfficeUsersMigrationFile sql file containing the migration to add the new office users
	OfficeUsersMigrationFilenameFlag string = "migration-filename"
	// VersionTimeFormat is the Go time format for creating a version number.
	VersionTimeFormat string = "20060102150405"
)

const (
	// template for adding office users
	createOfficeUser string = `INSERT INTO public.office_users
(id, user_id, first_name, last_name, middle_initials, email, telephone, transportation_office_id, created_at, updated_at)
VALUES
{{- range $i, $e := $}}
{{ if $i}},{{end}}('{{$e.ID}}', NULL, '{{$e.FirstName}}', '{{$e.LastName}}', {{if .MiddleInitials}}'{{.MiddleInitials}}'{{else}}NULL{{end}}, '{{$e.Email}}', '{{$e.Telephone}}', '{{$e.TransportationOfficeID}}', now(), now())
{{- end}};
`

	// template to apply secure migration
	migration string = `exec("./apply-secure-migration.sh {{.}}")`
)

// OfficeUsersFilenameFlag initializes add_office_users command line flags
func InitAddOfficeUsersFlags(flag *pflag.FlagSet) {
	flag.StringP(OfficeUsersFilenameFlag, "f", "", "File name of csv file containing the new office users")
	flag.StringP(OfficeUsersMigrationFilenameFlag, "n", "", "File name of the migration files for the new office users")
}

// CheckAddOfficeUsers validates add_office_users command line flags
func CheckAddOfficeUsers(v *viper.Viper) error {
	officeUsersFileName := v.GetString(OfficeUsersFilenameFlag)
	if officeUsersFileName == "" {
		return fmt.Errorf("--office-users-filename is required")
	}
	officeUsersMigrationFilenameFlag := v.GetString(OfficeUsersMigrationFilenameFlag)
	if officeUsersMigrationFilenameFlag == "" {
		return fmt.Errorf("--migration-filename is required")
	}
	return nil
}

func initGenOfficeUserMigrationFlags(flag *pflag.FlagSet) {
	// Migration Config
	cli.InitMigrationFlags(flag)

	// Add Office Users
	InitAddOfficeUsersFlags(flag)

	// Sort command line flags
	flag.SortFlags = true
}

func ValidateOfficeUser(o *models.OfficeUser) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: o.LastName, Name: "LastName"},
		&validators.StringIsPresent{Field: o.FirstName, Name: "FirstName"},
		&validators.StringIsPresent{Field: o.Email, Name: "Email"},
		&validators.StringIsPresent{Field: o.Telephone, Name: "Telephone"},
		&validators.UUIDIsPresent{Field: o.TransportationOfficeID, Name: "TransportationOfficeID"},
	), nil
}

func readOfficeUsersCSV(fileName string) ([]models.OfficeUser, error) {
	csvFile, err := os.Open(fileName)
	if err != nil {
		return []models.OfficeUser{}, err
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	//skip header
	_, err = reader.Read()
	if err != nil {
		return []models.OfficeUser{}, err
	}
	var officeUsers []models.OfficeUser
	var id, transportOfficeUUID uuid.UUID
	var line []string
	for {
		line, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []models.OfficeUser{}, err
		}
		transportOfficeUUID, err = uuid.FromString(strings.TrimSpace(line[5]))
		if err != nil {
			return []models.OfficeUser{}, err
		}
		var middleInitials *string
		if line[1] != "" {
			mi := strings.TrimSpace(line[1])
			middleInitials = &mi
		}
		id, err = uuid.NewV4()
		officeUser := models.OfficeUser{
			ID:                     id,
			FirstName:              strings.TrimSpace(line[0]),
			MiddleInitials:         middleInitials,
			LastName:               strings.TrimSpace(line[2]),
			Email:                  strings.TrimSpace(line[3]),
			Telephone:              strings.TrimSpace(line[4]),
			TransportationOfficeID: transportOfficeUUID,
		}
		verrs, err := ValidateOfficeUser(&officeUser)
		if verrs.HasAny() {
			return []models.OfficeUser{}, fmt.Errorf("validation errors for office user %v: %v", officeUser, verrs)
		}
		if err != nil {
			return []models.OfficeUser{}, err
		}
		officeUsers = append(officeUsers, officeUser)
	}
	return officeUsers, nil
}

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
	log.Printf("new migration file created at:  %q\n", migrationPath)
	return nil
}

func genOfficeUserMigration(cmd *cobra.Command, args []string) error {
	err := cmd.ParseFlags(args)
	flag := cmd.Flags()
	err = flag.Parse(os.Args[1:])
	if err != nil {
		return errors.Wrap(err, "could not parse flags")
	}
	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		return errors.Wrap(err, "could not bind flags")
	}
	err = CheckAddOfficeUsers(v)
	if err != nil {
		return err
	}
	migrationsPath := v.GetString(cli.MigrationPathFlag)
	migrationManifest := v.GetString(cli.MigrationManifestFlag)
	officeUsersFileName := v.GetString(OfficeUsersFilenameFlag)
	migrationFileName := v.GetString(OfficeUsersMigrationFilenameFlag)

	officeUsers, err := readOfficeUsersCSV(officeUsersFileName)
	if err != nil {
		return errors.Wrap(err, "error reading csv file")
	}

	secureMigrationName := fmt.Sprintf("%s_%s.up.sql", time.Now().Format(VersionTimeFormat), migrationFileName)
	t1 := template.Must(template.New("add_office_user").Parse(createOfficeUser))
	err = createMigration("./tmp", secureMigrationName, t1, officeUsers)
	if err != nil {
		return err
	}
	localMigrationPath := filepath.Join("local_migrations", secureMigrationName)
	localMigrationFile, err := os.Create(localMigrationPath)
	defer closeFile(localMigrationFile)
	if err != nil {
		return errors.Wrapf(err, "error creating %s", localMigrationPath)
	}
	log.Printf("new migration file created at:  %q\n", localMigrationPath)

	migrationName := fmt.Sprintf("%s_%s.up.fizz", time.Now().Format(VersionTimeFormat), migrationFileName)
	t2 := template.Must(template.New("migration").Parse(migration))
	err = createMigration(migrationsPath, migrationName, t2, secureMigrationName)
	if err != nil {
		return err
	}

	err = addMigrationToManifest(migrationManifest, migrationName)
	if err != nil {
		return err
	}
	return nil
}
