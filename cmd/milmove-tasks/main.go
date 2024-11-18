package main

import (
	"os"

	"github.com/spf13/cobra"
)

// GitCommit is empty unless set as a build flag
// See https://blog.alexellis.io/inject-build-time-vars-golang/
var gitBranch string
var gitCommit string

func main() {

	root := cobra.Command{
		Use:   "milmove-tasks [flags]",
		Short: "MilMove tasks",
		Long:  "MilMove tasks",
	}

	root.AddCommand(&cobra.Command{
		Use:          "version",
		Short:        "Print version information to stdout",
		Long:         "Print version information to stdout",
		RunE:         versionFunction,
		SilenceUsage: true,
	})

	saveGHCFuelPriceDataCommand := &cobra.Command{
		Use:          "save-ghc-fuel-price-data",
		Short:        "saves GHC diesel fuel price data",
		Long:         "saves the national weekly average GHC diesel fuel price data from the EIA Open Data API",
		RunE:         saveGHCFuelPriceData,
		SilenceUsage: true,
	}
	initSaveGHCFuelPriceFlags(saveGHCFuelPriceDataCommand.Flags())
	root.AddCommand(saveGHCFuelPriceDataCommand)

	sendPaymentReminderCommand := &cobra.Command{
		Use:          "send-payment-reminder",
		Short:        "sends payment reminder email",
		Long:         "sends payment reminder email",
		RunE:         sendPaymentReminder,
		SilenceUsage: true,
	}
	initPaymentReminderFlags(sendPaymentReminderCommand.Flags())
	root.AddCommand(sendPaymentReminderCommand)

	postFileToGEXCommand := &cobra.Command{
		Use:          "post-file-to-gex",
		Short:        "posts a file to GEX",
		Long:         "posts a file to GEX",
		RunE:         postFileToGEX,
		SilenceUsage: true,
	}
	initPostFileToGEXFlags(postFileToGEXCommand.Flags())
	root.AddCommand(postFileToGEXCommand)

	connectToGEXViaSFTPCommand := &cobra.Command{
		Use:          "connect-to-gex-via-sftp",
		Short:        "connects to GEX via SFTP",
		Long:         "connects to GEX via SFTP",
		RunE:         connectToGEXViaSFTP,
		SilenceUsage: true,
	}
	initConnectToGEXViaSFTPFlags(connectToGEXViaSFTPCommand.Flags())
	root.AddCommand(connectToGEXViaSFTPCommand)

	processEDIsCommand := &cobra.Command{
		Use:          "process-edis",
		Short:        "process EDIs asynchrounously",
		Long:         "process EDIs asynchrounously",
		RunE:         processEDIs,
		SilenceUsage: true,
	}
	initConnectToGEXViaSFTPFlags(processEDIsCommand.Flags())
	root.AddCommand(processEDIsCommand)

	processTPPSCommand := &cobra.Command{
		Use:          "process-tpps",
		Short:        "process TPPS files asynchrounously",
		Long:         "process TPPS files asynchrounously",
		RunE:         processTPPS,
		SilenceUsage: true,
	}
	root.AddCommand(processTPPSCommand)

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  "To install completion scripts run:\n\nmilmove-tasks completion > /usr/local/etc/bash_completion.d/milmove-tasks",
		RunE: func(_ *cobra.Command, _ []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(completionCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
