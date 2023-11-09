package trdm_test

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type TRDMSuite struct {
	*testingsuite.PopTestSuite
	viper  *viper.Viper
	logger *zap.Logger
	creds  aws.Credentials
}

type initFlags func(f *pflag.FlagSet)

func (suite *TRDMSuite) Setup(fn initFlags, flagSet []string) {
	suite.viper = nil

	flag := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	fn(flag)
	suite.NoError(flag.Parse(flagSet))

	v := viper.New()
	err := v.BindPFlags(flag)
	if err != nil {
		suite.logger.Fatal("could not bind flags", zap.Error(err))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	suite.SetViper(v)
}

func (suite *TRDMSuite) SetViper(v *viper.Viper) {
	suite.viper = v
}

func TestTRDMSuite(t *testing.T) {
	flag := pflag.CommandLine

	v := viper.New()
	bindErr := v.BindPFlags(flag)
	if bindErr != nil {
		log.Fatal("failed to bind flags", zap.Error(bindErr))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	logger := zaptest.NewLogger(t)

	// Setup mock creds
	mockCreds := aws.Credentials{
		AccessKeyID:     "mockAccessKeyID",
		SecretAccessKey: "mockSecretAccessKey",
		SessionToken:    "mockSessionToken",
		Source:          "mockProvider",
	}

	hs := &TRDMSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		viper:        v,
		logger:       logger,
		creds:        mockCreds,
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
