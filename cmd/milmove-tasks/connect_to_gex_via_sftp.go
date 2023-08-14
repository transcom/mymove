package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

// Call this from the command line with go run ./cmd/milmove-tasks connect-to-gex-via-sftp

func initConnectToGEXViaSFTPFlags(flag *pflag.FlagSet) {
	// Logging Levels
	cli.InitLoggingFlags(flag)

	// GEX SFTP
	cli.InitGEXSFTPFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

func connectToGEXViaSFTP(_ *cobra.Command, _ []string) error {
	v := viper.New()

	logger, _, err := logging.Config(
		logging.WithEnvironment(v.GetString(cli.LoggingEnvFlag)),
		logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)),
		logging.WithStacktraceLength(v.GetInt(cli.StacktraceLengthFlag)),
	)
	if err != nil {
		logger.Fatal("Failed to initialize Zap logging", zap.Error(err))
	}
	zap.ReplaceGlobals(logger)

	flag := pflag.CommandLine
	initConnectToGEXViaSFTPFlags(flag)
	err = flag.Parse(os.Args[1:])
	if err != nil {
		log.Fatal("failed to parse flags", zap.Error(err))
	}

	err = v.BindPFlags(flag)
	if err != nil {
		log.Fatal("failed to bind flags", zap.Error(err))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	// SSH and SFTP Connection Setup
	sshClient, err := cli.InitGEXSSH(logger, v)
	if err != nil {
		logger.Error("couldn't initialize SSH client", zap.Error(err))
		return err
	}
	defer func() {
		if closeErr := sshClient.Close(); closeErr != nil {
			logger.Error("could not close SFTP client", zap.Error(closeErr))
		}
	}()

	sftpClient, err := cli.InitGEXSFTP(logger, sshClient)
	if err != nil {
		logger.Error("couldn't initialize SFTP client", zap.Error(err))
		return err
	}
	defer func() {
		if closeErr := sftpClient.Close(); closeErr != nil {
			logger.Error("could not close SFTP client", zap.Error(closeErr))
		}
	}()

	pwd, err := sftpClient.Getwd()
	if err != nil {
		return err
	}

	fmt.Printf("Successfully connected via SFTP. The present working directory is %v", pwd)

	return nil
}
