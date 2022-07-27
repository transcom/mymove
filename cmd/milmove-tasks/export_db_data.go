package main

import (
	"fmt"
	"os"
	"strconv"

	pg "github.com/habx/pg-commands"
	"github.com/spf13/cobra"
)

func exportDBData(cmd *cobra.Command, args []string) error {

	pgConfig := getConfig()
	dump := pg.NewDump(&pgConfig)
	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: false})
	if dumpExec.Error != nil {
		return dumpExec.Error.Err
	}

	fmt.Println("Dump success")
	fmt.Println(dumpExec.Output)

	// export to S3 bucket
	return nil
}

func getConfig() pg.Postgres {
	host := getEnvOrPanic("PGHOST")
	port, err := strconv.Atoi(getEnvOrPanic("PGPORT"))
	if err != nil {
		panic("PGPORT must be an integer")
	}
	db := getEnvOrPanic("PGDB")
	user := getEnvOrPanic("PGUSER")
	password := getEnvOrPanic("PGPASSWORD")

	return pg.Postgres{
		Host:     host,
		Port:     port,
		DB:       db,
		Username: user,
		Password: password,
	}
}

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if len(value) <= 0 {
		panic(fmt.Sprintf("config loading failed; required environment variable %s must be set", key))
	}
	return value
}
