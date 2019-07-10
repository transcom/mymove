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

	"github.com/gofrs/uuid"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/models"
)

const createOfficeUser = `INSERT INTO public.office_users
(id, user_id, first_name, last_name, middle_initials, email, telephone, transportation_office_id, created_at, updated_at)
VALUES
{{- range $i, $e := $}}
{{ if $i}},{{end}}('{{$e.ID}}', NULL, '{{$e.FirstName}}', '{{$e.LastName}}', NULL, '{{$e.Email}}', '{{$e.Telephone}}', '{{$e.TransportationOfficeID}}', now(), now())
{{- end}};
`

const Migration = `exec("./apply-secure-migration.sh {{.}}")`

const (
	// OfficeUsersFilenameFlag filename containing the details for new office users
	OfficeUsersFilenameFlag string = "office-users-filename"
	// OfficeUsersMigrationFile sql file containing the migration to add the new office users
	OfficeUsersMigrationFilenameFlag string = "migration-filename"
	// VersionTimeFormat is the Go time format for creating a version number.
	VersionTimeFormat string = "20060102150405"
)

// OfficeUsersFilenameFlag initializes add_office_users command line flags
func InitAddOfficeUsersFlags(flag *pflag.FlagSet) {
	flag.StringP(OfficeUsersFilenameFlag, "f", "", "File name of csv file containing the new office users")
	flag.StringP(OfficeUsersMigrationFilenameFlag, "o", "", "File name of sql file containing the migration for the new office users")
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

func initFlags(flag *pflag.FlagSet) {
	// Verbose
	cli.InitVerboseFlags(flag)

	// Add Office Users
	InitAddOfficeUsersFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
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
		transportOfficeUUID, err = uuid.FromString(strings.TrimSpace(line[4]))
		if err != nil {
			return []models.OfficeUser{}, err
		}
		id, err = uuid.NewV4()
		officeUser := models.OfficeUser{
			ID:                     id,
			FirstName:              strings.TrimSpace(line[0]),
			LastName:               strings.TrimSpace(line[1]),
			Email:                  strings.TrimSpace(line[2]),
			Telephone:              strings.TrimSpace(line[3]),
			TransportationOfficeID: transportOfficeUUID,
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

func main() {
	flag := pflag.CommandLine
	initFlags(flag)
	err := flag.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		log.Fatal(err)
	}
	err = CheckAddOfficeUsers(v)
	if err != nil {
		log.Fatal(err)
	}
	officeUsersFileName := v.GetString(OfficeUsersFilenameFlag)
	migrationFileName := v.GetString(OfficeUsersMigrationFilenameFlag)

	officeUsers, err := readOfficeUsersCSV(officeUsersFileName)
	if err != nil {
		log.Fatal("error reading csv file: ", err)
	}

	secureMigrationName := fmt.Sprintf("%s_%s.up.sql", time.Now().Format(VersionTimeFormat), migrationFileName)
	secureMigrationPath := filepath.Join("./tmp", secureMigrationName)
	outfile1, err := os.Create(secureMigrationPath)
	defer closeFile(outfile1)
	if err != nil {
		log.Fatalf("error creating %s: %v\n", secureMigrationPath, err)
	}
	t1 := template.Must(template.New("add_office_user").Parse(createOfficeUser))
	err = t1.Execute(outfile1, officeUsers)
	if err != nil {
		log.Println("error executing template: ", err)
	}
	log.Printf("new secure migration file created at:  %q\n", secureMigrationPath)

	migrationName := fmt.Sprintf("%s_%s.up.fizz", time.Now().Format(VersionTimeFormat), migrationFileName)
	migrationPath := filepath.Join("./migrations", migrationName)
	outfile2, err := os.Create(migrationPath)
	defer closeFile(outfile2)
	if err != nil {
		log.Fatalf("error creating %s: %v\n", migrationPath, err)
	}
	t2 := template.Must(template.New("migration").Parse(Migration))
	err = t2.Execute(outfile2, secureMigrationName)
	if err != nil {
		log.Println("error executing template: ", err)
	}
	log.Printf("new migration file created at:  %q\n", migrationPath)
}
