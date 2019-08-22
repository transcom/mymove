package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/models"
)

const (
	// OfficeUsersFilenameFlag filename containing the details for new office users
	OfficeUsersFilenameFlag string = "office-users-filename"

	// template for adding office users
	createOfficeUser string = `INSERT INTO public.office_users
(id, user_id, first_name, last_name, middle_initials, email, telephone, transportation_office_id, created_at, updated_at)
VALUES
{{- range $i, $e := $}}
{{ if $i}},{{end}}('{{$e.ID}}', NULL, '{{$e.FirstName}}', '{{$e.LastName}}', {{if .MiddleInitials}}'{{.MiddleInitials}}'{{else}}NULL{{end}}, '{{$e.Email}}', '{{$e.Telephone}}', '{{$e.TransportationOfficeID}}', now(), now())
{{- end}};
`
)

// InitAddOfficeUsersFlags initializes command line flags
func InitAddOfficeUsersFlags(flag *pflag.FlagSet) {
	flag.StringP(OfficeUsersFilenameFlag, "f", "", "File name of csv file containing the new office users")
}

func initGenOfficeUserMigrationFlags(flag *pflag.FlagSet) {
	// Migration Config
	cli.InitMigrationFlags(flag)

	// Migration File Config
	cli.InitMigrationFileFlags(flag)

	// Add Office Users
	InitAddOfficeUsersFlags(flag)

	// Don't sort command line flags
	flag.SortFlags = false
}

// CheckAddOfficeUsersFlags validates add_office_users command line flags
func CheckAddOfficeUsersFlags(v *viper.Viper) error {
	if err := cli.CheckMigration(v); err != nil {
		return err
	}

	if err := cli.CheckMigrationFile(v); err != nil {
		return err
	}

	officeUsersFilename := v.GetString(OfficeUsersFilenameFlag)
	if officeUsersFilename == "" {
		return errors.Errorf("%s is missing", OfficeUsersFilenameFlag)
	}
	return nil
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
		if err != nil {
			return []models.OfficeUser{}, err
		}
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

func genOfficeUserMigration(cmd *cobra.Command, args []string) error {
	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Wrap(err, "could not ParseFlags on args")
	}

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
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	err = CheckAddOfficeUsersFlags(v)
	if err != nil {
		return err
	}

	migrationManifest := v.GetString(cli.MigrationManifestFlag)
	migrationName := v.GetString(cli.MigrationNameFlag)
	migrationVersion := v.GetString(cli.MigrationVersionFlag)
	officeUsersFilename := v.GetString(OfficeUsersFilenameFlag)

	officeUsers, err := readOfficeUsersCSV(officeUsersFilename)
	if err != nil {
		return errors.Wrap(err, "error reading csv file")
	}

	secureMigrationName := fmt.Sprintf("%s_%s.up.sql", migrationVersion, migrationName)
	t1 := template.Must(template.New("add_office_user").Parse(createOfficeUser))
	err = createMigration(tempMigrationPath, secureMigrationName, t1, officeUsers)
	if err != nil {
		return err
	}

	t2 := template.Must(template.New("local_migrations").Parse(localMigrationTemplate))
	err = createMigration("./local_migrations", secureMigrationName, t2, nil)
	if err != nil {
		return err
	}

	err = addMigrationToManifest(migrationManifest, secureMigrationName)
	if err != nil {
		return err
	}
	return nil
}
