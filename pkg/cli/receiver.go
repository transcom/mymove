package cli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// ReceiverBackend is the Receiver Backend Flag
	ReceiverBackendFlag string = "receiver-backend"
	// AWSSNSObjectTagsAddedTopic is the AWS SNS Object Tags Added Topic Flag
	AWSSNSObjectTagsAddedTopicFlag string = "aws-sns-object-tags-added-topic"
	// AWSS3RegionFlag is the AWS SNS Region Flag
	AWSSNSRegionFlag string = "aws-sns-region"
	// AWSSNSAccountId is the application's AWS account id
	AWSSNSAccountId string = "aws-account-id"
)

// InitReceiverFlags initializes Storage command line flags
func InitReceiverFlags(flag *pflag.FlagSet) {
	flag.String(ReceiverBackendFlag, "local", "Receiver backend to use, either local or sns&sqs.")
	flag.String(AWSSNSObjectTagsAddedTopicFlag, "", "SNS Topic for receiving event messages")
	flag.String(AWSSNSRegionFlag, "", "AWS region used for SNS and SQS")
	flag.String(AWSSNSAccountId, "", "AWS account Id")
}

// CheckReceiver validates Storage command line flags
func CheckReceiver(v *viper.Viper) error {

	receiverBackend := v.GetString(ReceiverBackendFlag)
	if !stringSliceContains([]string{"local", "sns&sqs"}, receiverBackend) {
		return fmt.Errorf("invalid receiver-backend %s, expecting local or sns&sqs", receiverBackend)
	}

	if receiverBackend == "sns&sqs" {
		r := v.GetString(AWSSNSRegionFlag)
		if r == "" {
			return fmt.Errorf("invalid value for %s: %s", AWSSNSRegionFlag, r)
		}
		topic := v.GetString(AWSSNSObjectTagsAddedTopicFlag)
		if topic == "" {
			return fmt.Errorf("invalid value for %s: %s", AWSSNSObjectTagsAddedTopicFlag, topic)
		}
		accountId := v.GetString(AWSSNSAccountId)
		if topic == "" {
			return fmt.Errorf("invalid value for %s: %s", AWSSNSAccountId, accountId)
		}
	}

	return nil
}
