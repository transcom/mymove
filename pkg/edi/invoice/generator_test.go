package ediinvoice_test

import (
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/db/sequence"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type InvoiceSuite struct {
	testingsuite.PopTestSuite
	logger       ediinvoice.Logger
	Viper        *viper.Viper
	icnSequencer sequence.Sequencer
}

func TestInvoiceSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	flag := pflag.CommandLine
	// Flag to update the test EDI
	// Borrowed from https://about.sourcegraph.com/go/advanced-testing-in-go
	flag.Bool("update", false, "update .golden files")
	// Flag to toggle Invoice usage indicator from P>T (Production>Test)
	flag.Bool("send-prod-invoice", false, "Send Production Invoice")

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	hs := &InvoiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
		Viper:        v,
	}

	hs.icnSequencer = sequence.NewDatabaseSequencer(hs.DB(), ediinvoice.ICNSequenceName)

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
