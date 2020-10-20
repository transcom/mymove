package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

// Call this from command line with go run ./cmd/send-to-syncada-via-sftp/ --local-file-path <localFilePath> --destination-file-name <destinationFileName>

func checkConfig(v *viper.Viper, logger logger) error {

	logger.Debug("checking config")

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	return nil
}

func initFlags(flag *pflag.FlagSet) {

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	// DB Config
	cli.InitSyncadaFlags(flag)

	flag.String("local-file-path", "", "The path where the file to be sent is located")
	flag.String("destination-file-name", "", "The name of the file to be stored in Syncada")

	// Don't sort flags
	flag.SortFlags = false
}

func main() {
	flag := pflag.CommandLine
	initFlags(flag)
	parseErr := flag.Parse(os.Args[1:])
	if parseErr != nil {
		log.Fatal("failed to parse flags", zap.Error(parseErr))
	}

	v := viper.New()
	bindErr := v.BindPFlags(flag)
	if bindErr != nil {
		log.Fatal("failed to bind flags", zap.Error(bindErr))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	fmt.Println("logger: ", logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	userID := v.GetString(cli.SyncadaSFTPUserIDFlag)
	password := v.GetString(cli.SyncadaSFTPPsswrdFlag)
	remote := v.GetString(cli.SyncadaSFTPIPAddressFlag)
	port := v.GetString(cli.SyncadaSFTPPortFlag)
	syncadaInboundDirectory := v.GetString(cli.SyncadaSFTPInboundDirectoryFlag)

	config := &ssh.ClientConfig{
		User: userID,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		/* #nosec */
		// The hostKey was removed because authentication is performed using a user ID and password
		// If hostKey configuration is needed, please see PR #5039: https://github.com/transcom/mymove/pull/5039
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// connect
	connection, err := ssh.Dial("tcp", remote+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	// create new SFTP client
	client, err := sftp.NewClient(connection)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// open local file
	localFile, err := os.Open(filepath.Clean(v.GetString("local-file-path")))
	if err != nil {
		log.Fatal(err)
	}

	// create destination file
	destinationFileName := v.GetString(("destination-file-name"))
	destinationFilePath := fmt.Sprintf("/%s/%s/%s", userID, syncadaInboundDirectory, destinationFileName)
	destinationFile, err := client.Create(destinationFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer destinationFile.Close()

	// copy source file to destination file
	bytes, err := io.Copy(destinationFile, localFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d bytes copied\n", bytes)
}
