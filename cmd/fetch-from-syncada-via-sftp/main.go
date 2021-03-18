package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/transcom/mymove/pkg/services"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

// Call this from command line with go run ./cmd/fetch-from-syncada-via-sftp/ --local-file-path <localFilePath> --syncada-file-name <syncadaFileName>

// TODO I'm not really sure where to put this. I think it should be usable for any SFTP connections.
// TODO And probably should be parameterized and not take everything from env vars
func createSyncadaSFTPClient() (services.SFTPClient, error) {
	port := os.Getenv("SYNCADA_SFTP_PORT")
	if port == "" {
		return nil, fmt.Errorf("Invalid credentials sftp missing SYNCADA_SFTP_PORT")
	}

	userID := os.Getenv("SYNCADA_SFTP_USER_ID")
	if userID == "" {
		return nil, fmt.Errorf("Invalid credentials sftp missing SYNCADA_SFTP_USER_ID")
	}

	remote := os.Getenv("SYNCADA_SFTP_IP_ADDRESS")
	if remote == "" {
		return nil, fmt.Errorf("Invalid credentials sftp missing SYNCADA_SFTP_IP_ADDRESS")
	}

	password := os.Getenv("SYNCADA_SFTP_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("Invalid credentials sftp missing SYNCADA_SFTP_PASSWORD")
	}

	hostKeyString := os.Getenv("SYNCADA_SFTP_HOST_KEY")
	if hostKeyString == "" {
		return nil, fmt.Errorf("Invalid credentials sftp missing SYNCADA_SFTP_HOST_KEY")
	}
	hostKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(hostKeyString))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse host key %w", err)
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
		return nil, err
	}
	//defer connection.Close()

	// create new SFTP client
	client, err := sftp.NewClient(connection)
	if err != nil {
		return nil, err
	}
	//defer client.Close()

	// TODO dont forget about the ssh client, it seems like probablty the sftp client has a connection to this that  aybe kt can close automatically

	return client, nil
}

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

	// Logging Levels
	cli.InitLoggingFlags(flag)

	flag.String("last-read-time", "", "Files older than this time will not be fetched.")
	flag.String("directory", "", "syncada path")

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

	logger, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		log.Fatalf("failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	client, err := createSyncadaSFTPClient()
	defer client.Close()
	if err != nil {
		logger.Fatal("couldn't initialize sftp session", zap.Error(err))
	}

	//2021-03-16T18:25:36Z
	t, err := time.Parse(time.RFC3339, v.GetString("last-read-time"))
	if err != nil {
		logger.Error("couldnt parse time", zap.Error(err))
	}
	//t := time.Now().Add(-1 * time.Hour)
	logger.Info("lastRead", zap.String("a", t.String()))
	//syncadaSFTPSession := invoice.InitNewSyncadaSFTPReaderSession(client, logger)
	//data, _, err := syncadaSFTPSession.FetchAndProcessSyncadaFiles(v.GetString("directory"), t, )
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Print(data)
}
