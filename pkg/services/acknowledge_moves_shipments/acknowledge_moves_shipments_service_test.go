package acknowledgemovesshipments

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type AcknowledgeMovesAndShipmentsServiceSuite struct {
	*testingsuite.PopTestSuite
	fs *afero.Afero
}

func TestAcknowledgeMovesAndShipmentsServiceSuite(t *testing.T) {
	var f = afero.NewMemMapFs()
	file := &afero.Afero{Fs: f}
	ts := &AcknowledgeMovesAndShipmentsServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		fs:           file,
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
