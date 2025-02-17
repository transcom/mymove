package entitlements

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

// EntitlementsServiceSuite is a suite for testing entitlement service objects
type EntitlementsServiceSuite struct {
	*testingsuite.PopTestSuite
	fs *afero.Afero
}

func TestEntitlementsServiceSuite(t *testing.T) {
	var f = afero.NewMemMapFs()
	file := &afero.Afero{Fs: f}
	ts := &EntitlementsServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		fs:           file,
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
