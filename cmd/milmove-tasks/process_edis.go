package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/certs"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/db/sequence"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/services/invoice"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
)

const (
	// ProcessEDILastReadTimeFlag is the ENV var for the last read time
	ProcessEDILastReadTimeFlag string = "process-edi-last-read-time"
	// ProcessEDIDeleteFilesFlag is the ENV var for deleting SFTP files after they've been processed
	ProcessEDIDeleteFilesFlag string = "process-edi-delete-files"
)

// Call this from the command line with go run ./cmd/milmove-tasks process-edis
func checkProcessEDIsConfig(v *viper.Viper, logger *zap.Logger) error {
	logger.Debug("checking config")

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}

	err = cli.CheckLogging(v)
	if err != nil {
		return err
	}

	err = cli.CheckGEXSFTP(v)
	if err != nil {
		return err
	}

	if err := cli.CheckGEX(v); err != nil {
		return err
	}

	if err := cli.CheckCert(v); err != nil {
		return err
	}

	if err := cli.CheckEntrustCert(v); err != nil {
		return err
	}

	return nil
}

func initProcessEDIsFlags(flag *pflag.FlagSet) {
	// Logging Levels
	cli.InitLoggingFlags(flag)

	// DB Config
	cli.InitDatabaseFlags(flag)

	// GEX
	cli.InitGEXFlags(flag)

	// Certificate
	cli.InitCertFlags(flag)

	// Entrust Certificates
	cli.InitEntrustCertFlags(flag)

	// GEX SFTP Config
	cli.InitGEXSFTPFlags(flag)

	flag.String(ProcessEDILastReadTimeFlag, "", "Files older than this RFC3339 time will not be fetched.")
	flag.Bool(ProcessEDIDeleteFilesFlag, false, "If present, delete files on SFTP server that have been processed successfully")

	// Don't sort flags
	flag.SortFlags = false
}

func processEDIs(cmd *cobra.Command, args []string) error {
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

	startTime := time.Now()
	defer func() {
		elapsedTime := time.Since(startTime)
		logger.Info(fmt.Sprintf("Duration of processEDIs task: %v", elapsedTime))
	}()

	flag := pflag.CommandLine
	initProcessEDIsFlags(flag)
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

	err = checkProcessEDIsConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	var session *awssession.Session
	if v.GetBool(cli.DbIamFlag) {
		c := &aws.Config{
			Region: aws.String(v.GetString(cli.AWSRegionFlag)),
		}
		s, errorSession := awssession.NewSession(c)
		if errorSession != nil {
			logger.Fatal(errors.Wrap(errorSession, "error creating aws session").Error())
		}
		session = s
	}

	var dbCreds *credentials.Credentials
	if v.GetBool(cli.DbIamFlag) {
		if session != nil {
			// We want to get the credentials from the logged in AWS session rather than create directly,
			// because the session conflates the environment, shared, and container metadata config
			// within NewSession.  With stscreds, we use the Secure Token Service,
			// to assume the given role (that has rds db connect permissions).
			dbIamRole := v.GetString(cli.DbIamRoleFlag)
			logger.Info(fmt.Sprintf("assuming AWS role %q for db connection", dbIamRole))
			dbCreds = stscreds.NewCredentials(session, dbIamRole)
		}
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, dbCreds, logger)
	if err != nil {
		// No connection object means that the configuraton failed to validate and we should not startup
		// A valid connection object that still has an error indicates that the DB is not up and we should not startup
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	appCtx := appcontext.NewAppContext(dbConnection, logger, nil)
	dbEnv := v.GetString(cli.DbEnvFlag)
	gexURL := v.GetString(cli.GEXURLFlag)

	isDevOrTest := dbEnv == "development" || dbEnv == "test"
	if isDevOrTest {
		logger.Info(fmt.Sprintf("Starting in %s mode, which enables additional features", dbEnv))
	}

	sendToSyncada := v.GetBool(cli.SendToSyncada)
	logger.Info(fmt.Sprintf("SendToSyncada is %v", sendToSyncada))
	processEdiDeleteFiles := v.GetBool(ProcessEDIDeleteFilesFlag)
	logger.Info(fmt.Sprintf("ProcessEDIDeleteFiles is %v", processEdiDeleteFiles))

	// Set the ICNSequencer in the handler: if we are in dev/test mode and sending to a real
	// GEX URL, then we should use a random ICN number within a defined range to avoid duplicate
	// test ICNs in Syncada.
	var icnSequencer sequence.Sequencer
	if isDevOrTest && len(gexURL) > 0 {
		// ICNs are 9-digit numbers; reserve the ones in an upper range for development/testing.
		icnSequencer, err = sequence.NewRandomSequencer(ediinvoice.ICNRandomMin, ediinvoice.ICNRandomMax)
		if err != nil {
			logger.Fatal("Could not create random sequencer for ICN", zap.Error(err))
		}
	} else {
		icnSequencer = sequence.NewDatabaseSequencer(ediinvoice.ICNSequenceName)
	}

	// TODO I don't know why we need a separate logger for cert stuff
	certLogger, _, err := logging.Config(logging.WithEnvironment(dbEnv), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
	if err != nil {
		logger.Fatal("Failed to initialize Zap logging", zap.Error(err))
	}
	certificates, rootCAs, err := certs.InitDoDEntrustCertificates(v, certLogger)
	if certificates == nil || rootCAs == nil || err != nil {
		logger.Fatal("Error in getting tls certs", zap.Error(err))
	}
	tlsConfig := &tls.Config{Certificates: certificates, RootCAs: rootCAs, MinVersion: tls.VersionTLS12}

	gexSender := invoice.NewGexSenderHTTP(
		gexURL,
		true,
		tlsConfig,
		v.GetString(cli.GEXBasicAuthUsernameFlag),
		v.GetString(cli.GEXBasicAuthPasswordFlag))

	reviewedPaymentRequestProcessor, err := paymentrequest.InitNewPaymentRequestReviewedProcessor(appCtx, sendToSyncada, icnSequencer, gexSender)
	if err != nil {
		logger.Fatal("InitNewPaymentRequestReviewedProcessor failed", zap.Error(err))
	}

	// Process 858s
	reviewedPaymentRequestProcessor.ProcessReviewedPaymentRequest(appCtx)
	logger.Info("Finished processing reviewed payment requests")

	if !sendToSyncada {
		logger.Info("Skipping processing of response files EDI997 acknowledgement and EDI824 application advice responses")
		return nil
	}

	// SSH and SFTP Connection Setup
	sshClient, err := cli.InitGEXSSH(appCtx, v)
	if err != nil {
		logger.Fatal("couldn't initialize SSH client", zap.Error(err))
	}
	defer func() {
		if closeErr := sshClient.Close(); closeErr != nil {
			logger.Error("could not close SFTP client", zap.Error(closeErr))
		}
	}()

	sftpClient, err := cli.InitGEXSFTP(appCtx, sshClient)
	if err != nil {
		logger.Fatal("couldn't initialize SFTP client", zap.Error(err))
	}
	defer func() {
		if closeErr := sftpClient.Close(); closeErr != nil {
			logger.Error("could not close SFTP client", zap.Error(closeErr))
		}
	}()

	wrappedSFTPClient := invoice.NewSFTPClientWrapper(sftpClient)
	syncadaSFTPSession := invoice.NewSyncadaSFTPReaderSession(wrappedSFTPClient, v.GetBool(ProcessEDIDeleteFilesFlag))

	// Sample expected format: 2021-03-16T18:25:36Z
	lastReadTimeFlag := v.GetString(ProcessEDILastReadTimeFlag)
	var lastReadTime time.Time
	if lastReadTimeFlag != "" {
		lastReadTime, err = time.Parse(time.RFC3339, lastReadTimeFlag)
		if err != nil {
			logger.Error("couldn't parse last read time", zap.Error(err))
		}
	}
	logger.Info("lastRead", zap.String("lastReadTime", lastReadTime.String()))

	// Process 997s
	path997 := v.GetString(cli.GEXSFTP997PickupDirectory)
	_, err = syncadaSFTPSession.FetchAndProcessSyncadaFiles(appCtx, path997, lastReadTime, invoice.NewEDI997Processor())
	if err != nil {
		logger.Error("Error reading EDI997 acknowledgement responses", zap.Error(err))
	} else {
		logger.Info("Successfully processed EDI997 acknowledgement responses")
	}

	// Process 824s
	path824 := v.GetString(cli.GEXSFTP824PickupDirectory)
	_, err = syncadaSFTPSession.FetchAndProcessSyncadaFiles(appCtx, path824, lastReadTime, invoice.NewEDI824Processor())
	if err != nil {
		logger.Error("Error reading EDI824 application advice responses", zap.Error(err))
	} else {
		logger.Info("Successfully processed EDI824 application advice responses")
	}

	return nil
}
