package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
)

const (
	// CertPathFlag is the path to the client mTLS certificate
	CertPathFlag string = "certpath"
	// KeyPathFlag is the path to the key mTLS certificate
	KeyPathFlag string = "keypath"
	// HostnameFlag is the hostname to connect to
	HostnameFlag string = "hostname"
	// PortFlag is the port to connect to
	PortFlag string = "port"
	// InsecureFlag indicates that TLS verification and validation can be skipped
	InsecureFlag string = "insecure"
)

// initRootFlags initializes flags relating to the notification api
func initRootFlags(flag *pflag.FlagSet) {
	cli.InitCACFlags(flag)
	cli.InitVerboseFlags(flag)

	flag.String(CertPathFlag, "./config/tls/devlocal-mtls.cer", "Path to the public cert")
	flag.String(KeyPathFlag, "./config/tls/devlocal-mtls.key", "Path to the private key")
	flag.String(HostnameFlag, cli.HTTPPrimeServerNameLocal, "The hostname to connect to")
	flag.Int(PortFlag, cli.MutualTLSPort, "The port to connect to")
	flag.Bool(InsecureFlag, false, "Skip TLS verification and validation")
}

// Debug prints helpful debugging information for requests
func Debug(data []byte, err error) {
	if err == nil {
		log.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

// Logger type exports the logger for use in the command files
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

// CheckRootConfig checks the validity of the notification api flags
func CheckRootConfig(v *viper.Viper) error {
	err := cli.CheckCAC(v)
	if err != nil {
		return err
	}

	err = cli.CheckVerbose(v)
	if err != nil {
		return err
	}

	if (v.GetString(CertPathFlag) != "" && v.GetString(KeyPathFlag) == "") || (v.GetString(CertPathFlag) == "" && v.GetString(KeyPathFlag) != "") {
		return fmt.Errorf("Both TLS certificate and key paths must be provided")
	}

	return nil
}

func main() {
	root := cobra.Command{
		Use:   "webhook-client [flags]",
		Short: "Webhook client",
		Long:  "Webhook client",
	}
	initRootFlags(root.PersistentFlags())

	postWebhookNotifyCommand := &cobra.Command{
		Use:          "post-webhook-notify",
		Short:        "Post Webhook Notify",
		Long:         "Post Webhook Notify",
		RunE:         postWebhookNotify,
		SilenceUsage: true,
	}
	initPostWebhookNotifyFlags(postWebhookNotifyCommand.Flags())
	root.AddCommand(postWebhookNotifyCommand)

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  "To install completion scripts run:\n\nwebhook-client completion > /usr/local/etc/bash_completion.d/webhook-client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(completionCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
