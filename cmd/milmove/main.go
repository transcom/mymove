package main

import (
	"os"

	"github.com/spf13/cobra"
)

// GitCommit is empty unless set as a build flag
// See https://blog.alexellis.io/inject-build-time-vars-golang/
var gitBranch string
var gitCommit string

func stringSliceContains(stringSlice []string, value string) bool {
	for _, x := range stringSlice {
		if value == x {
			return true
		}
	}
	return false
}

func main() {

	root := cobra.Command{
		Use:   "milmove [flags]",
		Short: "Webserver for MilMove",
		Long:  "Webserver for MilMove",
	}

	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information to stdout",
		Long:  "Print version information to stdout",
		RunE:  versionFunction,
	})

	serveCommand := &cobra.Command{
		Use:   "serve",
		Short: "Runs MilMove webserver",
		Long:  "Runs MilMove webserver",
		RunE:  serveFunction,
	}
	initServeFlags(serveCommand.Flags())
	root.AddCommand(serveCommand)

	migrateCommand := &cobra.Command{
		Use:           "migrate",
		Short:         "Runs MilMove migrations",
		Long:          "Runs MilMove migrations",
		RunE:          migrateFunction,
		SilenceUsage:  true, // not needed
		SilenceErrors: true, // not needed
	}
	initMigrateFlags(migrateCommand.Flags())
	root.AddCommand(migrateCommand)

	genCommand := &cobra.Command{
		Use:   "gen",
		Short: "Generate migrations and other objects",
		Long:  "Generate migrations and other objects",
		RunE:  nil,
	}
	root.AddCommand(genCommand)

	genMigrationCommand := &cobra.Command{
		Use:                   "migration -n NAME [-t TYPE]",
		Short:                 "Generate migrations and other objects",
		Long:                  "Generate migrations and other objects",
		RunE:                  genMigrationFunction,
		DisableFlagsInUseLine: true,
		SilenceErrors:         true, // not needed
	}
	initGenMigrationFlags(genMigrationCommand.Flags())
	genCommand.AddCommand(genMigrationCommand)

	genOfficeUserMigrationCommand := &cobra.Command{
		Use:                   "office-user-migration -f CSV_FILENAME -n MIGRATION_NAME",
		Short:                 "Generate migrations required for adding office users",
		Long:                  "Generate migrations required for adding office users",
		RunE:                  genOfficeUserMigration,
		DisableFlagsInUseLine: true,
		SilenceErrors:         true, // not needed
	}
	initGenOfficeUserMigrationFlags(genOfficeUserMigrationCommand.Flags())
	genCommand.AddCommand(genOfficeUserMigrationCommand)

	genDisableUserMigrationCommand := &cobra.Command{
		Use:                   "disable-user-migration -e EMAIL",
		Short:                 "Generate migrations required for disabling a user",
		Long:                  "Generate migrations required for disabling a user",
		RunE:                  genDisableUserMigration,
		DisableFlagsInUseLine: true,
	}
	initDisableUserMigrationFlags(genDisableUserMigrationCommand.Flags())
	genCommand.AddCommand(genDisableUserMigrationCommand)

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  "To install completion scripts run:\n\nmilmove completion > /usr/local/etc/bash_completion.d/milmove",
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(completionCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
