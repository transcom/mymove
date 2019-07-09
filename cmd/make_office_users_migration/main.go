package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

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

const (
	//OfficeUsersFilenameFlag filename containing the details for new office users
	OfficeUsersFilenameFlag string = "office-users-filename"
	//OfficeUsersMigrationFile sql file containing the migration to add the new office users
	OfficeUsersMigrationFilenameFlag string = "migration-filename"
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

func readOfficeUsersCSV(v *viper.Viper) ([]models.OfficeUser, error) {
	csvFile, _ := os.Open(v.GetString(OfficeUsersFilenameFlag))
	reader := csv.NewReader(bufio.NewReader(csvFile))
	//skip header
	_, err := reader.Read()
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
	officeUsers, err := readOfficeUsersCSV(v)
	if err != nil {
		log.Fatal("error reading csv file: ", err)
	}
	outfile, err := os.Create(OfficeUsersMigrationFilenameFlag)
	defer closeFile(outfile)
	if err != nil {
		log.Fatalf("error creating %s: %v\n", OfficeUsersMigrationFilenameFlag, err)
	}
	t := template.Must(template.New("add_office_user").Parse(createOfficeUser))
	err = t.Execute(outfile, officeUsers)
	if err != nil {
		log.Println("error executing template: ", err)
	}
}
