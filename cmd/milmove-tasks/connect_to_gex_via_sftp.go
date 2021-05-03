package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

// Call this from the command line with go run ./cmd/milmove-tasks connect-to-gex-via-sftp

func checkConnectToGEXViaSFTPConfig(v *viper.Viper, logger logger) error {
	logger.Debug("checking config")

	if err := cli.CheckGEX(v); err != nil {
		return err
	}

	return nil
}

func initConnectToGEXViaSFTPFlags(flag *pflag.FlagSet) {
	// Logging Levels
	cli.InitLoggingFlags(flag)

	// GEX
	cli.InitGEXFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

func connectToGEXViaSFTP(cmd *cobra.Command, args []string) error {
	v := viper.New()

	logger, err := logging.Config(
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

	err = checkConnectToGEXViaSFTPConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	port := os.Getenv("GEX_SFTP_PORT")
	if port == "" {
		return fmt.Errorf("Invalid credentials SFTP missing GEX_SFTP_PORT")
	}

	userID := os.Getenv("GEX_SFTP_USER_ID")
	if userID == "" {
		return fmt.Errorf("Invalid credentials SFTP missing GEX_SFTP_USER_ID")
	}

	remote := os.Getenv("GEX_SFTP_IP_ADDRESS")
	if remote == "" {
		return fmt.Errorf("Invalid credentials SFTP missing GEX_SFTP_IP_ADDRESS")
	}

	password := os.Getenv("GEX_SFTP_PASSWORD")
	if password == "" {
		return fmt.Errorf("Invalid credentials SFTP missing GEX_SFTP_PASSWORD")
	}

	hostKeyString := os.Getenv("GEX_SFTP_HOST_KEY")
	if hostKeyString == "" {
		return fmt.Errorf("Invalid credentials sftp missing GEX_SFTP_HOST_KEY")
	}
	hostKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(hostKeyString))
	if err != nil {
		return fmt.Errorf("Failed to parse host key %w", err)
	}

	config := &ssh.ClientConfig{
		User: userID,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// connect
	connection, err := ssh.Dial("tcp", remote+":"+port, config)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := connection.Close(); closeErr != nil {
			logger.Debug("Failed to close tcp connection", zap.Error(closeErr))
		}
	}()

	// create new SFTP client
	client, err := sftp.NewClient(connection)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			logger.Debug("Failed to close sftp client connection", zap.Error(closeErr))
		}
	}()

	pwd, err := client.Getwd()
	if err != nil {
		return err
	}

	fmt.Printf("Successfully connected via SFTP. The present working directory is %v", pwd)

	return nil
}
