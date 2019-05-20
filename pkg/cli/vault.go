package cli

import (
	"fmt"

	"github.com/99designs/aws-vault/prompt"
	"github.com/99designs/aws-vault/vault"
	"github.com/99designs/keyring"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// VaultAWSKeychainNameFlag is the aws-vault keychain name Flag
	VaultAWSKeychainNameFlag string = "aws-vault-keychain-name"
	// VaultAWSProfileFlag is the aws-vault profile name Flag
	VaultAWSProfileFlag string = "aws-profile"
	// VaultAWSVaultFlag is the aws-vault flag
	VaultAWSVaultFlag string = "aws-vault"
	// VaultAWSSessionTokenFlag is the AWS session token flag
	VaultAWSSessionTokenFlag string = "aws-session-token"

	// VaultAWSKeychainNameDefault is the aws-vault default keychain name
	VaultAWSKeychainNameDefault string = "login"
	// VaultAWSProfileDefault is the aws-vault default profile name
	VaultAWSProfileDefault string = "transcom-ppp"
)

type errInvalidKeychainName struct {
	KeychainName string
}

func (e *errInvalidKeychainName) Error() string {
	return fmt.Sprintf("invalid keychain name '%s'", e.KeychainName)
}

type errInvalidAWSProfile struct {
	Profile string
}

func (e *errInvalidAWSProfile) Error() string {
	return fmt.Sprintf("invalid aws profile '%s'", e.Profile)
}

type errInvalidVault struct {
	KeychainName string
	Profile      string
}

func (e *errInvalidVault) Error() string {
	return fmt.Sprintf("invalid keychain name '%s' or profile '%s'", e.KeychainName, e.Profile)
}

// InitVaultFlags initializes Vault command line flags
func InitVaultFlags(flag *pflag.FlagSet) {
	// Flags default to empty string to facilitate deploys from CircleCI
	flag.String(VaultAWSKeychainNameFlag, "", "The aws-vault keychain name")
	flag.String(VaultAWSProfileFlag, "", "The aws-vault profile")
}

// CheckVault validates Vault command line flags
func CheckVault(v *viper.Viper) error {
	if awsVault := v.GetString(VaultAWSVaultFlag); len(awsVault) > 0 {
		if sessionToken := v.GetString(VaultAWSSessionTokenFlag); len(sessionToken) == 0 {
			return errors.New("in aws-vault session, but missing aws-session-token")
		}
	}

	// Both keychain name and profile are required or both must be missing
	keychainName := v.GetString(VaultAWSKeychainNameFlag)
	keychainNames := []string{
		VaultAWSKeychainNameDefault,
	}
	if len(keychainName) > 0 && !stringSliceContains(keychainNames, keychainName) {
		return errors.Wrap(&errInvalidKeychainName{KeychainName: keychainName},
			fmt.Sprintf("%s is invalid, expected %v", VaultAWSKeychainNameFlag, keychainNames))
	}

	awsProfile := v.GetString(VaultAWSProfileFlag)
	awsProfiles := []string{
		VaultAWSProfileDefault,
	}
	if len(awsProfile) > 0 && !stringSliceContains(awsProfiles, awsProfile) {
		return errors.Wrap(&errInvalidAWSProfile{Profile: awsProfile},
			fmt.Sprintf("%s is invalid, expected %v", VaultAWSProfileFlag, awsProfiles))
	}

	// Require both are set or neither are set
	if (len(keychainName) != 0 && len(awsProfile) == 0) || (len(keychainName) == 0 && len(awsProfile) != 0) {
		return errors.Wrap(&errInvalidVault{KeychainName: keychainName, Profile: awsProfile},
			fmt.Sprintf("If either %s or %s is is set the other is required", VaultAWSKeychainNameFlag, VaultAWSProfileFlag))
	}
	return nil
}

// GetAWSCredentials uses aws-vault to return AWS credentials
func GetAWSCredentials(keychainName string, awsProfile string) (*credentials.Credentials, error) {

	// Open the keyring which holds the credentials
	ring, err := keyring.Open(keyring.Config{
		ServiceName:              "aws-vault",
		AllowedBackends:          []keyring.BackendType{keyring.KeychainBackend},
		KeychainName:             keychainName,
		KeychainTrustApplication: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Unable to configure and open keyring")
	}

	// Prepare options for the vault before creating the provider
	vConfig, err := vault.LoadConfigFromEnv()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to load AWS config from environment")
	}
	vOptions := vault.VaultOptions{
		Config:    vConfig,
		MfaPrompt: prompt.Method("terminal"),
	}
	vOptions = vOptions.ApplyDefaults()
	err = vOptions.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to validate aws-vault options")
	}

	// Get a new provider to retrieve the credentials
	provider, err := vault.NewVaultProvider(ring, awsProfile, vOptions)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create aws-vault provider")
	}
	credVals, err := provider.Retrieve()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to retrieve aws credentials from aws-vault")
	}
	return credentials.NewStaticCredentialsFromCreds(credVals), nil
}

// GetAWSConfig returns an AWS Config struct using aws-vault credentials for use in an AWS session
func GetAWSConfig(v *viper.Viper, verbose bool) (*aws.Config, error) {

	awsRegion := v.GetString(AWSRegionFlag)

	awsConfig := &aws.Config{
		Region: aws.String(awsRegion),
	}

	// If program is not wrapped in aws-vault wrapper then get credentials
	if awsVault := v.GetString(VaultAWSVaultFlag); len(awsVault) == 0 {
		keychainName := v.GetString(VaultAWSKeychainNameFlag)
		awsProfile := v.GetString(VaultAWSProfileFlag)
		if len(keychainName) > 0 && len(awsProfile) > 0 {
			creds, getAWSCredsErr := GetAWSCredentials(keychainName, awsProfile)
			if getAWSCredsErr != nil {
				return nil, errors.Wrap(getAWSCredsErr,
					fmt.Sprintf("Unable to get AWS credentials from the keychain %s and profile %s", keychainName, awsProfile))
			}
			awsConfig.CredentialsChainVerboseErrors = aws.Bool(verbose)
			awsConfig.Credentials = creds
		}
	}
	return awsConfig, nil
}
