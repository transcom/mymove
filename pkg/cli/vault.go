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
	// VaultKeychainNameFlag is the aws-vault keychain name Flag
	VaultKeychainNameFlag string = "aws-vault-keychain-name"
	// VaultProfileFlag is the aws-vault profile name Flag
	VaultProfileFlag string = "aws-profile"

	// VaultKeychainNameDefault is the aws-vault default keychain name
	VaultKeychainNameDefault string = "login"
	// VaultProfileDefault is the aws-vault default profile name
	VaultProfileDefault string = "transcom-ppp"
)

type errInvalidKeychainName struct {
	KeychainName string
}

func (e *errInvalidKeychainName) Error() string {
	return fmt.Sprintf("invalid keychain name %s", e.KeychainName)
}

type errInvalidProfile struct {
	Profile string
}

func (e *errInvalidProfile) Error() string {
	return fmt.Sprintf("invalid profile %s", e.Profile)
}

// InitVaultFlags initializes Vault command line flags
func InitVaultFlags(flag *pflag.FlagSet) {
	flag.String(VaultKeychainNameFlag, VaultKeychainNameDefault, "The aws-vault keychain name")
	flag.String(VaultProfileFlag, VaultProfileDefault, "The aws-vault profile")
}

// CheckVault validates Vault command line flags
func CheckVault(v *viper.Viper) error {
	keychainName := v.GetString(VaultKeychainNameFlag)
	if len(keychainName) == 0 {
		return errors.Wrap(&errInvalidKeychainName{KeychainName: keychainName}, fmt.Sprintf("%q is invalid", VaultKeychainNameFlag))
	}

	keychainProfile := v.GetString(VaultProfileFlag)
	if len(keychainProfile) == 0 {
		return errors.Wrap(&errInvalidProfile{Profile: keychainName}, fmt.Sprintf("%q is invalid", VaultProfileFlag))
	}
	return nil
}

// GetAWSCredentials uses aws-vault to return AWS credentials
func GetAWSCredentials(keychainName string, keychainProfile string) (*credentials.Credentials, error) {

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
	provider, err := vault.NewVaultProvider(ring, keychainProfile, vOptions)
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

	keychainName := v.GetString(VaultKeychainNameFlag)
	keychainProfile := v.GetString(VaultProfileFlag)

	if len(keychainName) > 0 && len(keychainProfile) > 0 {
		creds, getAWSCredsErr := GetAWSCredentials(keychainName, keychainProfile)
		if getAWSCredsErr != nil {
			return nil, errors.Wrap(getAWSCredsErr, fmt.Sprintf("Unable to get AWS credentials from the keychain %s and profile %s", keychainName, keychainProfile))
		}
		awsConfig.CredentialsChainVerboseErrors = aws.Bool(verbose)
		awsConfig.Credentials = creds
	}
	return awsConfig, nil
}
