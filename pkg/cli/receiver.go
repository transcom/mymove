package cli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// ReceiverBackendFlag is the Receiver Backend Flag
	ReceiverBackendFlag string = "receiver-backend"
	// SNSTagsUpdatedTopicFlag is the SNS Tags Updated Topic Flag
	SNSTagsUpdatedTopicFlag string = "sns-tags-updated-topic"
	// SNSRegionFlag is the SNS Region flag
	SNSRegionFlag string = "sns-region"
	// SNSAccountId is the application's AWS account id
	SNSAccountId string = "aws-account-id"
	// ReceiverCleanupOnStartFlag is the Receiver Cleanup On Start Flag
	ReceiverCleanupOnStartFlag string = "receiver-cleanup-on-start"
)

// InitReceiverFlags initializes Storage command line flags
func InitReceiverFlags(flag *pflag.FlagSet) {
	flag.String(ReceiverBackendFlag, "local", "Receiver backend to use, either local or sns_sqs.")
	flag.String(SNSTagsUpdatedTopicFlag, "", "SNS Topic for receiving event messages")
	flag.String(SNSRegionFlag, "", "Region used for SNS and SQS")
	flag.String(SNSAccountId, "", "SNS account Id")
	flag.Bool(ReceiverCleanupOnStartFlag, false, "Receiver will cleanup previous aws artifacts on start.")
}

// CheckReceiver validates Storage command line flags
func CheckReceiver(v *viper.Viper) error {

	receiverBackend := v.GetString(ReceiverBackendFlag)
	if !stringSliceContains([]string{"local", "sns_sqs"}, receiverBackend) {
		return fmt.Errorf("invalid receiver_backend %s, expecting local or sns_sqs", receiverBackend)
	}

	receiverCleanupOnStart := v.GetString(ReceiverCleanupOnStartFlag)
	if !stringSliceContains([]string{"true", "false"}, receiverCleanupOnStart) {
		return fmt.Errorf("invalid receiver_cleanup_on_start %s, expecting true or false", receiverCleanupOnStart)
	}

	if receiverBackend == "sns_sqs" {
		r := v.GetString(SNSRegionFlag)
		if r == "" {
			return fmt.Errorf("invalid value for %s: %s", SNSRegionFlag, r)
		}
		topic := v.GetString(SNSTagsUpdatedTopicFlag)
		if topic == "" {
			return fmt.Errorf("invalid value for %s: %s", SNSTagsUpdatedTopicFlag, topic)
		}
		accountId := v.GetString(SNSAccountId)
		if topic == "" {
			return fmt.Errorf("invalid value for %s: %s", SNSAccountId, accountId)
		}
	}

	return nil
}
